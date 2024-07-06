package feishu

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"io"
	"net/http"
	"strings"
	"thundersoft.com/brain/DigitalVisitor/conf"
	"thundersoft.com/brain/DigitalVisitor/utils"
)

type Response struct {
	Code              int    `json:"code"`
	Expire            int    `json:"expire"`
	Msg               string `json:"msg"`
	TenantAccessToken string `json:"tenant_access_token"`
}

func decodeString(body io.Reader, output *string) error {
	b, err := io.ReadAll(body)
	if err != nil {
		return err
	}
	*output = string(b)
	return nil
}
func decodeResponse(body io.Reader, v any) error {
	if v == nil {
		return nil
	}

	switch o := v.(type) {
	case *string:
		return decodeString(body, o)
	default:
		return json.NewDecoder(body).Decode(v)
	}
}

// GetTenantAccessToken 获取 Tenant Access Token
func GetTenantAccessToken() (*utils.AccessTokenResponse, error) {
	data := map[string]interface{}{
		"app_id":     conf.ConfigInfo.FlyBook.AppId,
		"app_secret": conf.ConfigInfo.FlyBook.AppSecret,
	}
	// 将数据转换为JSON格式
	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Printf("转换JSON数据出错：%s\n", err)
		return nil, err
	}
	payload := strings.NewReader(string(jsonData))
	req, _ := http.NewRequest("POST", conf.ConfigInfo.FlyBook.TokenUrl, payload)
	req.Header.Add("Content-Type", "application/json; charset=utf-8")

	res, _ := http.DefaultClient.Do(req)
	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)

	var response utils.AccessTokenResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		fmt.Println("解析JSON数据失败:", err)
		return nil, err
	}
	return &response, nil
}

// GetDepartmentId 获取 获取部门id
func GetDepartmentId(openId, userAccessToken string) (*utils.DepartmentIdRes, error) {
	query := fmt.Sprintf("?user_id_type=open_id&department_id_type=open_department_id")
	url := conf.ConfigInfo.FlyBook.DepartmentIdUrl + "/" + openId + query

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", userAccessToken))
	req.Header.Add("Content-Type", "application/json; charset=utf-8")
	fmt.Println("userAccessToken=", userAccessToken)
	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)

	//resMap := make(map[string]interface{}, 0)
	var resMap utils.DepartmentIdRes
	err := json.Unmarshal(body, &resMap)
	if err != nil {
		fmt.Println("解析JSON数据失败:", err)
		return nil, err
	}
	fmt.Println("GetDepartmentId resMap=", resMap)
	return &resMap, nil

}

// GetUserInfo 获取 Tenant Access Token
func GetUserInfo(code string, accessToken utils.AccessTokenResponse) (*utils.UserInfoRes, error) {

	data := map[string]interface{}{
		"grant_type": "authorization_code",
		"code":       code,
	}
	// 将数据转换为JSON格式
	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Printf("转换JSON数据出错：%s\n", err)
		return nil, err
	}
	payload := strings.NewReader(string(jsonData))
	req, _ := http.NewRequest("POST", conf.ConfigInfo.FlyBook.UserUrl, payload)
	req.Header.Add("Content-Type", "application/json; charset=utf-8")
	req.Header.Add("Authorization", "Bearer "+accessToken.AppAccessToken)

	fmt.Println("url=", conf.ConfigInfo.FlyBook.UserUrl)
	fmt.Println("payload=", string(jsonData))
	fmt.Println("Authorization=", accessToken.AppAccessToken)

	res, _ := http.DefaultClient.Do(req)
	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)

	//resMap := make(map[string]interface{}, 0)
	var resMap utils.UserInfoRes
	err = json.Unmarshal(body, &resMap)
	if err != nil {
		fmt.Println("解析JSON数据失败:", err)
		return nil, err
	}
	fmt.Println("GetUserInfo resMap", resMap)
	return &resMap, nil
}

// GetUserToken 获取 User Access Token
func GetUserToken(code string, accessToken utils.AccessTokenResponse) (*utils.UserTokenRes, error) {

	data := map[string]interface{}{
		"grant_type": "authorization_code",
		"code":       code,
	}
	// 将数据转换为JSON格式
	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Printf("转换JSON数据出错：%s\n", err)
		return nil, err
	}
	payload := strings.NewReader(string(jsonData))
	req, _ := http.NewRequest("POST", conf.ConfigInfo.FlyBook.UserTokenUrl, payload)
	req.Header.Add("Content-Type", "application/json; charset=utf-8")
	req.Header.Add("Authorization", "Bearer "+accessToken.AppAccessToken)

	res, _ := http.DefaultClient.Do(req)
	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)

	var resMap utils.UserTokenRes
	err = json.Unmarshal(body, &resMap)
	if err != nil {
		fmt.Println("解析JSON数据失败:", err)
		return nil, err
	}
	return &resMap, nil
}

// GetUserDepartment 根据name 进行模糊查询 获取用户列表
func GetUserDepartment(departmentId, accessToken string) (*utils.DepartmentInfoRes, error) {
	query := fmt.Sprintf("?user_id_type=open_id&department_id_type=open_department_id")
	url := conf.ConfigInfo.FlyBook.DepartmentUrl + "/" + departmentId + query

	fmt.Println("-----url=", url)

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	req.Header.Add("Content-Type", "application/json; charset=utf-8")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)

	//resMap := make(map[string]interface{}, 0)
	var resMap utils.DepartmentInfoRes
	err := json.Unmarshal(body, &resMap)
	if err != nil {
		fmt.Println("解析JSON数据失败:", err)
		return nil, err
	}
	fmt.Println("GetUserDepartment resMap=", resMap)
	return &resMap, nil

}

// SendNessages 发送机器人消息 用于预定会议成功后的通知
func SendNessages(userid, content string) error {

	tokenResponse, err := GetTenantAccessToken()
	if err != nil {
		return err
	}
	if 0 != tokenResponse.Code {
		return nil
	}

	var messageRequest utils.FlyMessageRequest
	messageRequest.Content = `{"text":"` + content + `"}`

	messageRequest.MsgType = "text"

	messageRequest.ReceiveID = userid
	u := uuid.New()
	messageRequest.UUID = u.String()
	jsonData, err := json.Marshal(messageRequest)
	if err != nil {
		fmt.Println("转换为JSON字符串失败:", err)
		return err
	}
	fmt.Println("jsonData=", string(jsonData))
	payload := strings.NewReader(string(jsonData))

	tenantAccessToken := fmt.Sprintf("Bearer %s", tokenResponse.TenantAccessToken)
	req, err0 := http.NewRequest("POST", conf.ConfigInfo.FlyBook.MessagesUrl, payload)
	if err0 != nil {
		fmt.Println("创建请求失败:", err0)
		return nil
	}
	req.Header.Add("Authorization", tenantAccessToken)
	req.Header.Add("Content-Type", "application/json; charset=utf-8")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)
	var jsonMap map[string]interface{}
	err = json.Unmarshal(body, &jsonMap)
	if err != nil {
		fmt.Println("解析JSON数据失败:", err)
		return err
	}
	fmt.Println("SendNessages jsonMap=", jsonMap)

	return nil

}
