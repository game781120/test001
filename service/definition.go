package service

import (
	"github.com/gin-gonic/gin"
	"thundersoft.com/brain/DigitalVisitor/utils"
)

// GreetingService 定义了服务的接口
type GreetingService interface {
	StreamSendCommon(c *gin.Context, content string, isStream, isEnd bool)
	DealDataStreamLLm(c *gin.Context, content string, baseURL, authToken, modelName string)
	DealDataLLm(c *gin.Context, content string, baseURL, authToken, modelName string)
	Appointment(req *utils.AppointmentReq) utils.AppointmentRes
	CustomerSignIn(req utils.CustomerSignInReq) utils.CustomerSignInRes
	GetRecord(req *utils.QueryRecordReq) ([]utils.Appointment, error)
	DelRecord(req *utils.DelRecordReq) utils.DelRecordRes
}

// SimpleGreetingService 实现了 GreetingService 接口
type SimpleGreetingService struct{}
