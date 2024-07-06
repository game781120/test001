package router

import (
	"github.com/gin-gonic/gin"
	"thundersoft.com/brain/DigitalVisitor/controller"
	"thundersoft.com/brain/DigitalVisitor/service"
)

func SetRouter(router *gin.Engine) {

	// 设置上传文件的大小限制
	router.MaxMultipartMemory = 8 << 20 // 8 MiB

	greetingService := service.SimpleGreetingService{}
	// 创建控制器实例
	greetingController := controller.NewGreetingController(&greetingService)

	relayV1Router := router.Group("/v1")
	{
		relayV1Router.POST("/visitor/appointment", greetingController.Appointment)
		relayV1Router.POST("/visitor/delRecord", greetingController.DelRecord)
		relayV1Router.POST("/visitor/getQaList", greetingController.GetQaList)
		relayV1Router.POST("/visitor/chat", greetingController.Chat)
		relayV1Router.POST("/visitor/signIn", greetingController.CustomerSignIn)
		relayV1Router.POST("/visitor/getRecord", greetingController.GetRecord)
		relayV1Router.POST("/visitor/history", greetingController.History)
		relayV1Router.GET("/visitor/export", greetingController.Export)
		relayV1Router.GET("/visitor/appId", greetingController.GetAppId)
		relayV1Router.GET("/visitor/flyBookCall", greetingController.FlyBookCall)
		relayV1Router.GET("/visitor/isManager", greetingController.IsManger)
		relayV1Router.POST("/visitor/addManager", greetingController.AddManger)
		relayV1Router.GET("/visitor/queryManager", greetingController.QueryManager)
		relayV1Router.POST("/visitor/delManager", greetingController.DelManager)

		relayV1Router.POST("/chat/completions", greetingController.ChatHandler)

	}
}
