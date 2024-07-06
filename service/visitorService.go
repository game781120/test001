package service

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"thundersoft.com/brain/DigitalVisitor/mydb"
	"thundersoft.com/brain/DigitalVisitor/utils"
	"thundersoft.com/brainos/openai"
	"time"
)

func (s *SimpleGreetingService) DealDataStreamLLm(c *gin.Context, content string, baseURL, authToken, modelName string) {
	c.Status(http.StatusOK)
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Writer.Flush()

	client := openai.NewClient(baseURL, authToken)
	ctx := context.Background()
	req := openai.ChatCompletionRequest{
		Model: modelName,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleUser,
				Content: content,
			},
		},
		Stream: true,
	}
	stream, err := client.CreateChatCompletionStream(ctx, req)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return

	}
	defer stream.Close()

	//由于sdk 返回过来的数据已经没有 data 头了，所以此处需要重新拼装
	for {
		response, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			byte1 := []byte("data: DONE\n\n")
			c.Writer.Write(byte1)
			c.Writer.Flush()

			return
		}
		// 将数据转换为JSON格式
		jsonData, _ := json.Marshal(response)
		// 发送数据给客户端
		byte1 := []byte("data: ")
		byte2 := []byte("\n\n")
		slice3 := append(append(byte1, jsonData...), byte2...)
		c.Writer.Write(slice3)
		c.Writer.Flush()
		time.Sleep(8 * time.Millisecond)
	}
}

func (s SimpleGreetingService) StreamSendCommon(c *gin.Context, texts string, isStream, isEnd bool) {
	c.Status(http.StatusOK)
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Writer.Flush()

	var res openai.ChatCompletionStreamResponse
	res.Choices = make([]openai.ChatCompletionStreamChoice, 1)
	if isStream {
		runes := []rune(texts)
		for _, v := range runes {
			res.Choices[0].Delta.Content = string(v)

			jsonData, _ := json.Marshal(res)
			// 发送数据给客户端
			byte1 := []byte("data: ")
			byte2 := []byte("\n\n")
			slice3 := append(append(byte1, jsonData...), byte2...)
			c.Writer.Write(slice3)
			c.Writer.Flush()
			time.Sleep(8 * time.Millisecond)
		}
	} else {
		res.Choices[0].Delta.Content = texts
		jsonData, _ := json.Marshal(res)
		// 发送数据给客户端
		byte1 := []byte("data: ")
		byte2 := []byte("\n\n")
		slice3 := append(append(byte1, jsonData...), byte2...)
		c.Writer.Write(slice3)
		c.Writer.Flush()
		time.Sleep(8 * time.Millisecond)
	}
	if isEnd {
		byte1 := []byte("data: DONE\n\n")
		c.Writer.Write(byte1)
		c.Writer.Flush()
	}
}

func (s *SimpleGreetingService) DealDataLLm(c *gin.Context, content string, baseURL, authToken, modelName string) {
	client := openai.NewClient(baseURL, authToken)
	res, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: modelName,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: content,
				},
			},
		},
	)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": err.Error(),
			"context": "",
			"success": false,
			"error":   true,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "OK",
		"context": res.Choices[0].Message.Content,
		"success": true,
		"error":   false,
	})
}

// Appointment 访客预约
func (s *SimpleGreetingService) Appointment(req *utils.AppointmentReq) utils.AppointmentRes {
	var res utils.AppointmentRes
	id := mydb.InsertData(&req.Appointment)
	if len(id) > 0 {
		res.Code = 200
		res.Msg = "success"
		res.Appointment = req.Appointment
	} else {
		res.Code = http.StatusInternalServerError
		res.Msg = "插入到数据库失败"
		res.Appointment = req.Appointment
	}

	return res
}

// SignIn 访客签到
func (s *SimpleGreetingService) CustomerSignIn(req utils.CustomerSignInReq) utils.CustomerSignInRes {
	var res utils.CustomerSignInRes
	err := mydb.UpdateSingIn(req.AppointmentId)
	if err == nil {
		res.Code = 200
	} else {
		res.Code = http.StatusInternalServerError
	}
	return res
}

// GetRecord 访客签到
func (s *SimpleGreetingService) GetRecord(req *utils.QueryRecordReq) ([]utils.Appointment, error) {
	dataList, err := mydb.Query(req)
	if err == nil {
		if req.PageSize != 0 {
			if len(dataList) < req.PageSize {
				req.PageSize = len(dataList)
			} else if len(dataList) < req.PageSize*(req.Current-1) {
				return []utils.Appointment{}, nil
			}
			dataList = dataList[(req.Current-1)*req.PageSize : (req.Current)*req.PageSize]
		}
	} else {
		return nil, err
	}
	return dataList, nil
}

func (s *SimpleGreetingService) DelRecord(req *utils.DelRecordReq) utils.DelRecordRes {
	var res utils.DelRecordRes
	res.AppointmentIdList = make([]string, 0)
	for _, id := range req.AppointmentIdList {
		err := mydb.DelData(id)
		if err == nil {
			res.AppointmentIdList = append(res.AppointmentIdList, id)
		}
	}
	res.AppointmentIdList = req.AppointmentIdList
	res.Code = 200
	return res
}
