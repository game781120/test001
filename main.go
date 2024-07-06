package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"log/slog"
	"thundersoft.com/brain/DigitalVisitor/conf"
	"thundersoft.com/brain/DigitalVisitor/elastics"
	"thundersoft.com/brain/DigitalVisitor/logs"
	"thundersoft.com/brain/DigitalVisitor/myMilvus"
	"thundersoft.com/brain/DigitalVisitor/mydb"
	"thundersoft.com/brain/DigitalVisitor/router"
)

func main() {
	conf.LoadYaml()
	logs.InitLog()
	myMilvus.Init()
	mydb.Start()
	esInstance := &elastics.EsData{
		Username: conf.ConfigInfo.Elastic.Username,
		Password: conf.ConfigInfo.Elastic.Password,
		Host:     conf.ConfigInfo.Elastic.Host,
		Port:     conf.ConfigInfo.Elastic.Port,
		Index:    conf.ConfigInfo.Elastic.Index,
		ES:       nil,
		Mapping:  elastics.Mapping001,
	}
	if b, err := elastics.Start(esInstance); !b || err != nil {
		slog.Error("init esInstance failed ", "esInstance", esInstance, "error", err)
	}
	// 创建服务实例
	//greetingService := service.SimpleGreetingService{}
	//// 创建控制器实例
	//greetingController := controller.NewGreetingController(&greetingService)
	// 设置路由
	//http.HandleFunc("/greet", greetingController.SayHello)
	//// 启动HTTP服务器
	//log.Println("Server starting on port 8080...")
	//if err := http.ListenAndServe(":8080", nil); err != nil {
	//	log.Fatal(err)
	//}

	// Initialize HTTP server
	server := gin.New()
	server.Use(gin.Recovery())
	// This will cause SSE not to work!!!
	// server.Use(gzip.Gzip(gzip.DefaultCompression))
	router.SetRouter(server) // Set the router for the server
	portStr := fmt.Sprintf(":%s", conf.ConfigInfo.Server.Port)
	err := server.Run(portStr)
	if err != nil {
		log.Fatalf("failed to start HTTP server: " + err.Error())
	}
}
