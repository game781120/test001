package controller

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"thundersoft.com/brain/DigitalVisitor/conf"
	"thundersoft.com/brain/DigitalVisitor/feishu"
	"thundersoft.com/brain/DigitalVisitor/files"
	"thundersoft.com/brain/DigitalVisitor/logs"
	"thundersoft.com/brain/DigitalVisitor/myHttp"
	"thundersoft.com/brain/DigitalVisitor/myMilvus"
	"thundersoft.com/brain/DigitalVisitor/mydb"
	"thundersoft.com/brain/DigitalVisitor/service"
	"thundersoft.com/brain/DigitalVisitor/utils"
	"thundersoft.com/brainos/openai"
)

func (ct *Controller) DelRecord(c *gin.Context) {
	var req utils.DelRecordReq
	err := json.NewDecoder(c.Request.Body).Decode(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, nil)
		return
	}
	res := ct.service.DelRecord(&req)
	c.JSON(res.Code, res)
}

func (ct *Controller) GetQaList(c *gin.Context) {
	var req utils.QaReq
	err := json.NewDecoder(c.Request.Body).Decode(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, nil)
		return
	}
	var comRes utils.BaseResponse
	if len(req.KeyWords) > 0 {
		comRes.Data, err = myMilvus.QueryQaListKeyWords(req)
		if err != nil {
			comRes.Code = http.StatusInternalServerError
			c.JSON(comRes.Code, comRes)
			return
		}

	} else {
		comRes.Data, err = myMilvus.QueryQaListAll(req)
		if err != nil {
			comRes.Code = http.StatusInternalServerError
			c.JSON(comRes.Code, comRes)
			return
		}
		comRes.Code = 200
		c.JSON(comRes.Code, comRes)

	}
	comRes.Code = 200
	c.JSON(comRes.Code, comRes)

}

func render(data string, c *gin.Context) {
	c.Render(-1, utils.StringRender{
		Data: data,
	})
}
func streamOut(data string, done bool, c *gin.Context) {
	render(service.GetCompletionStreamChoice(data), c)
	if done {
		render(fmt.Sprintf("%s[DONE]\n\n", service.HeaderData), c)
	}
}

func streamOutEx(data string, done bool, c *gin.Context) {
	render(service.GetCompletionStreamChoiceEx(data), c)
	if done {
		render(fmt.Sprintf("%s[DONE]\n\n", service.HeaderData), c)
	}
}
func (ct *Controller) CreateMarkDownData(strType, content string) string {
	if strType == "video" {
		str1 := "<li><div><i class=\"video-play\"></i><i class=\"video-plus\"></i><i class=\"video-full\"></i></div><video controls disablePictureInPicture=\"true\" controlslist=\"nodownload noplaybackrate nofullscreen noremoteplayback\" src=\"%s\"></video></li>"
		return fmt.Sprintf(str1, content)
	} else if strType == "img" {
		str1 := "<li><img src=\"%s\" /><div><i class=\"img-full\"></i></div></li>"
		return fmt.Sprintf(str1, content)
	}
	return ""

}
func (ct *Controller) ChatHandler(c *gin.Context) {
	fmt.Println("-------ChatHandler")
	var chatRequest openai.ChatCompletionRequest
	err := json.NewDecoder(c.Request.Body).Decode(&chatRequest)
	if err != nil {
		streamOut(err.Error(), true, c)
		return
	}
	fmt.Println("chatRequest", chatRequest)
	datasetIDs := make([]string, 0)
	for k, v := range chatRequest.Custom.(map[string]interface{}) {
		if k == "knowledge_ids" {
			ids, ok := v.([]interface{})
			if ok {
				for _, id := range ids {
					datasetIDs = append(datasetIDs, id.(string))
				}
			}
		}
	}
	fmt.Println("datasetIDs", datasetIDs)

	messages := chatRequest.Messages
	question := messages[len(messages)-1].Content
	fmt.Println("question=", question)
	var junShiReq utils.JunShiReq
	junShiReq.Query = question

	res, err := myHttp.JuShiQuestion(junShiReq)
	fmt.Println("res=", res)
	if res == nil || err != nil {
		streamOut(err.Error(), true, c)
		return
	}
	if res.IsPass == "False" {
		streamOutEx("拒绝", true, c)
		return
	} else {
		render(service.GetCompletionStreamChoice(""), c)
	}
	nluRes, err := myHttp.NluQuestion(question)
	fmt.Println("--------nluRes=", nluRes)
	if nluRes.Code == 200 {
		if len(nluRes.Data.Nulagent.Reserve.Name) > 0 {
			nluByte, err := json.Marshal(nluRes.Data.Nulagent.Reserve)
			if err == nil {
				str1 := "\n\n```execute\n" + `{"nulagent":{"reserve":` + string(nluByte) + "}}" + "\n```\n\n"
				streamOut(str1, true, c)
			}
		}

		if len(nluRes.Data.Nulagent.SignIn.Name) > 0 {
			nluByte, err := json.Marshal(nluRes.Data.Nulagent.SignIn)
			if err == nil {
				str1 := "\n\n```execute\n" + `{"nulagent":{"signIn":` + string(nluByte) + "}}" + "\n```\n\n"
				streamOut(str1, true, c)
			}
		}
	}
	callData := ""
	queryRes, err := myMilvus.Query(question, datasetIDs)

	if err == nil && len(queryRes) > 0 {
		for _, v := range queryRes {
			if v.Qa == 1 {
				if v.Scores > 0.0 {
					continue
				}
				if len(v.Answer) > 0 {
					answer := fmt.Sprintf("### %s\n", v.Answer)
					ct.service.StreamSendCommon(c, answer, true, false)
				}
				callData += `<ul class="digi-more">`
				if len(v.ImageUrls) > 0 {
					imageUrlsList := strings.Split(v.ImageUrls, "###")
					for _, url := range imageUrlsList {
						callData += ct.CreateMarkDownData("img", url)
					}

				}
				if len(v.VideoUrls) > 0 {
					videoUrlsList := strings.Split(v.VideoUrls, "###")
					for _, url := range videoUrlsList {
						callData += ct.CreateMarkDownData("video", url)
					}

				}
				callData += `</ul>`
				ct.service.StreamSendCommon(c, callData, false, true)
				return
			} else {
				if v.Scores > 0.2 {
					continue
				}
				callData += v.Content + "\n"
			}
		}
	}

	fmt.Println("callData=", callData)

	content := ""
	if len(callData) > 0 {
		content = callData + "\n请根据上面的数据回答下面的问题。问题如下：\n" + question
	} else {
		content = question
	}
	fmt.Println("ChatHandler=", content)
	if chatRequest.Stream {
		ct.service.DealDataStreamLLm(c, content, conf.ConfigInfo.Model.LlmHost,
			conf.ConfigInfo.Model.LlmAuthToken, conf.ConfigInfo.Model.Name)
	} else {
		ct.service.DealDataLLm(c, content, conf.ConfigInfo.Model.LlmHost,
			conf.ConfigInfo.Model.LlmAuthToken, conf.ConfigInfo.Model.Name)
	}

}

// Appointment 访客预约
func (ct *Controller) Chat(c *gin.Context) {
	var req utils.ChatRequest
	err := json.NewDecoder(c.Request.Body).Decode(&req)
	if err != nil {
		return
	}

	DatasetList := strings.Split(conf.ConfigInfo.Milvus.DatasetIds, ",")
	callData := ""
	res, err := myMilvus.Query(req.Message, DatasetList)
	qa := 0
	if err == nil && len(res) > 0 {
		for _, v := range res {
			callData += v.Content + "\n"
			if v.Qa == 1 {
				qa = 1
			}
		}
	}
	if qa == 1 {
		ct.service.StreamSendCommon(c, callData, true, true)
		return
	}
	content := ""
	if len(callData) > 0 {
		content = callData + "\n请根据上面的数据回答下面的问题。问题如下：\n" + req.Message
	} else {
		content = req.Message
	}

	if req.Stream {
		ct.service.DealDataStreamLLm(c, content, conf.ConfigInfo.Model.LlmHost,
			conf.ConfigInfo.Model.LlmAuthToken, conf.ConfigInfo.Model.Name)
	} else {
		ct.service.DealDataLLm(c, content, conf.ConfigInfo.Model.LlmHost,
			conf.ConfigInfo.Model.LlmAuthToken, conf.ConfigInfo.Model.Name)
	}

}

// Appointment 访客预约
func (ct *Controller) Appointment(c *gin.Context) {
	var req utils.AppointmentReq
	err := json.NewDecoder(c.Request.Body).Decode(&req)
	if err != nil {
		logs.Info("Appointment", err)
		return
	}
	logs.Info("Appointment req", req)
	res := ct.service.Appointment(&req)
	c.JSON(res.Code, res)
}

// CustomerSignIn 访客签到
func (ct *Controller) CustomerSignIn(c *gin.Context) {
	var req utils.CustomerSignInReq
	err := json.NewDecoder(c.Request.Body).Decode(&req)
	if err != nil {
		return
	}
	res := ct.service.CustomerSignIn(req)
	if res.Code == 200 {
		userOneId, name, com, err := mydb.QueryAppointments(req.AppointmentId)
		if err == nil && len(userOneId) > 0 {
			content := fmt.Sprintf("%s 的  %s 已经签到 ！", com, name)
			feishu.SendNessages(userOneId, content)
		}

	}

	c.JSON(200, res)
}

// GetRecord 来访记录
func (ct *Controller) GetRecord(c *gin.Context) {
	var res utils.CommonResponse
	var req utils.QueryRecordReq
	err := json.NewDecoder(c.Request.Body).Decode(&req)
	if err != nil {
		res.Code = http.StatusBadRequest
		res.Msg = err.Error()
		c.JSON(res.Code, res)
		return
	}
	fmt.Println("GetRecord req=", req)
	if req.Type == 0 {
		if req.Current == 0 {
			req.Current = 1
		}
	}
	if len(req.ReceptionName) > 0 {
		var mData utils.ManagerData
		mData.Name = req.ReceptionName
		b, err := mydb.IsManager(&mData)
		if err == nil || b {
			req.ReceptionName = ""
		}
	}

	dataList, err := ct.service.GetRecord(&req)
	if err != nil {
		res.Code = http.StatusBadRequest
		res.Msg = err.Error()
		c.JSON(res.Code, res)
		return
	}
	res.Code = http.StatusOK
	res.Total = len(dataList)
	res.CurrentPage = req.Current
	res.Data = dataList
	c.JSON(res.Code, res)
}

func (ct *Controller) History(c *gin.Context) {
	var res utils.CommonResponse
	var req utils.QueryRecordReq
	err := json.NewDecoder(c.Request.Body).Decode(&req)
	if err != nil && err.Error() != "EOF" {
		res.Code = http.StatusBadRequest
		res.Msg = err.Error()
		c.JSON(res.Code, res)
		return
	}
	if req.Current == 0 {
		req.Current = 1
	}

	fmt.Println("History req=", req)
	dataList, err := ct.service.GetRecord(&req)
	if err != nil {
		res.Code = http.StatusBadRequest
		res.Msg = err.Error()
		c.JSON(res.Code, res)
		return
	}
	res.Code = http.StatusOK
	res.Total = len(dataList)
	res.CurrentPage = req.Current
	res.Data = dataList
	fmt.Println("History res=", res)
	c.JSON(res.Code, res)

}

func (ct *Controller) Export(c *gin.Context) {
	var req utils.QueryRecordReq
	err := json.NewDecoder(c.Request.Body).Decode(&req)
	if err != nil && err.Error() != "EOF" {
		c.JSON(200, nil)
		return
	}
	req.PageSize = 0
	res, err := ct.service.GetRecord(&req)
	filePath := files.WriteFile(res)
	// 打开 XLSX 文件
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		c.Status(http.StatusNotFound)
		return
	}
	defer f.Close()

	// 获取文件信息
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		fmt.Println(err)
		return
	}
	// 获取文件大小
	fileSize := fileInfo.Size()

	fileName := filepath.Base(filePath)
	str := fmt.Sprintf("attachment; filename=%s", fileName)
	// 设置 HTTP 头信息
	c.Header("Content-Disposition", str)
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Length", strconv.FormatInt(fileSize, 10))

	// 读取 XLSX 文件内容并发送
	if _, err := f.WriteTo(c.Writer); err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

}

func (ct *Controller) GetAppId(c *gin.Context) {
	flyBook := map[string]string{
		"appId":     conf.ConfigInfo.FlyBook.AppId,
		"appSecret": conf.ConfigInfo.FlyBook.AppSecret,
	}
	fmt.Println("---- GetAppId", flyBook)
	c.JSON(http.StatusOK, flyBook)
}

func (ct *Controller) FlyBookCall(c *gin.Context) {
	fmt.Println("---- FlyBookCall")
	code := c.DefaultQuery("code", "")
	fmt.Println("---- code=", code)
	if len(code) > 0 {
		accessToken, err := feishu.GetTenantAccessToken()
		if err != nil {
			fmt.Println("GetTenantAccessToken =", err)
			c.JSON(http.StatusBadRequest, nil)
			return
		}
		//获取到用户的 open_id,user_access_token,用户名称
		userRes, err := feishu.GetUserInfo(code, *accessToken)
		if err != nil {
			fmt.Println("GetUserInfo err=", err)
			c.JSON(http.StatusBadRequest, nil)
			return
		}
		//根据用户的 open_id,user_access_token 获取用户部门id
		depRes, err := feishu.GetDepartmentId(userRes.Data.OpenID, userRes.Data.AccessToken)
		if len(depRes.Data.User.DepartmentIDs) > 0 {
			//根据用户的部门id 和 user_access_token 获取部门名称
			departmentInfo, err := feishu.GetUserDepartment(depRes.Data.User.DepartmentIDs[0], userRes.Data.AccessToken)
			if err == nil && departmentInfo.Code == 0 {
				userRes.Data.DepartmentName = departmentInfo.Data.Department.Name
				fmt.Println("----- departmentInfo=", departmentInfo)
			}
		}

		c.JSON(http.StatusOK, userRes)
	} else {
		c.JSON(http.StatusBadRequest, nil)
	}

}

func (ct *Controller) IsManger(c *gin.Context) {
	resData := map[string]bool{
		"IsManger": false,
	}
	userName := c.DefaultQuery("userName", "")
	if len(userName) > 0 {
		var md utils.ManagerData
		md.Name = userName
		b, _ := mydb.IsManager(&md)
		if b {
			resData["IsManger"] = true
			c.JSON(http.StatusOK, resData)
			return
		}
	}
	c.JSON(http.StatusOK, resData)
}

func (ct *Controller) AddManger(c *gin.Context) {
	var res utils.CommonResponse
	var req utils.ManagerData
	err := json.NewDecoder(c.Request.Body).Decode(&req)
	if err != nil {
		res.Code = http.StatusBadRequest
		res.Msg = "添加管理员失败"
		c.JSON(res.Code, res)
		return
	}
	b, err := mydb.InsertManager(req)
	if err != nil || !b {
		res.Code = http.StatusBadRequest
		res.Msg = "添加管理员失败"
		c.JSON(res.Code, res)
		return
	}
	res.Code = http.StatusOK
	res.Msg = "添加管理员成功"
	c.JSON(res.Code, res)

}

func (ct *Controller) QueryManager(c *gin.Context) {
	pCurrent := 1
	pSize := 20
	keyWords := c.DefaultQuery("keyWords", "")
	current := c.DefaultQuery("current", "1")
	pageSize := c.DefaultQuery("pageSize", "20")
	if len(current) > 0 {
		pCurrent, _ = strconv.Atoi(current)
	}
	if len(pageSize) > 0 {
		pSize, _ = strconv.Atoi(pageSize)
	}
	var mRes utils.CommonResponse
	fmt.Printf("QueryManager keyWords=%s current=%d pageSize=%d\n", keyWords, pCurrent, pSize)
	res, err := mydb.QueryManager(keyWords, pCurrent, pSize)
	if err != nil {
		mRes.Code = http.StatusBadRequest
		mRes.Msg = err.Error()
		fmt.Println("QueryManager err=", err)
		c.JSON(mRes.Code, mRes)
		return
	}
	mRes.CurrentPage = pCurrent
	mRes.Total = len(res)
	mRes.Data = res
	mRes.Code = http.StatusOK
	fmt.Println("QueryManager=", mRes)
	c.JSON(mRes.Code, mRes)
}

func (ct *Controller) DelManager(c *gin.Context) {
	var res utils.CommonResponse
	var req utils.ManagerData
	err := json.NewDecoder(c.Request.Body).Decode(&req)
	if err != nil && err.Error() != "EOF" {
		res.Code = http.StatusBadRequest
		res.Msg = "删除管理员失败"
		c.JSON(res.Code, res)
		return
	}
	b, err := mydb.DelManager(req)
	if err != nil || !b {
		res.Code = http.StatusBadRequest
		res.Msg = "删除管理员失败"
		c.JSON(res.Code, res)
		return
	}
	res.Code = http.StatusOK
	res.Msg = "删除管理员成功"
	c.JSON(res.Code, res)

}
