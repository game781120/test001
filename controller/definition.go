package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"thundersoft.com/brain/DigitalVisitor/service"
)

// GreetingController 定义了问候的控制器
type Controller struct {
	BaseController
	service service.GreetingService
}

// NewGreetingController 创建一个新的 GreetingController 实例
func NewGreetingController(s service.GreetingService) *Controller {
	return &Controller{service: s}
}

type BaseController struct {
}

func (con BaseController) Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "OK",
		"context": data,
		"success": true,
		"error":   false,
	})
}

func (con BaseController) Error(c *gin.Context, msg string) {
	c.JSON(http.StatusInternalServerError, gin.H{
		"code":    500,
		"message": msg,
		"context": nil,
		"success": false,
		"error":   true,
	})
}
