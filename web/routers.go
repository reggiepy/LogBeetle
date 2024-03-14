package web

import (
	"github.com/gin-gonic/gin"
	"github.com/nsqio/go-nsq"
	"log"
	"net/http"
)

var producer *nsq.Producer

func init() {
	initNSQProducer()
}

func initNSQProducer() {
	config := nsq.NewConfig()
	var err error
	producer, err = nsq.NewProducer("192.168.1.110:4150", config)
	if err != nil {
		log.Fatalf("Failed to create NSQ producer: %v", err)
	}
}

// SetupRouter 设置路由
func SetupRouter() *gin.Engine {
	// 创建路由引擎
	router := gin.Default()
	gin.DebugPrintRouteFunc = func(httpMethod, absolutePath, handlerName string, nuHandlers int) {
		log.Printf("endpoint %v %v %v %v\n", httpMethod, absolutePath, handlerName, nuHandlers)
	}
	// 注册全局中间件
	router.Use(ErrorHandler())
	router.Use(DddRequestID())

	// 定义路由和处理程序
	router.GET("/", homeHandler)
	router.GET("/ping", pingHandler)
	router.GET("/about", aboutHandler)
	router.GET("/send-message", sendMessageHandler)

	return router
}

// 处理首页请求
func homeHandler(c *gin.Context) {
	c.JSON(http.StatusOK, "hello world.")
}

// ping
func pingHandler(c *gin.Context) {
	c.JSON(http.StatusOK, "pong")
}

// 处理关于页面请求
func aboutHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "about.html", gin.H{
		"title": "About Us",
	})
}

// @Summary 发送消息
// @Description 发送消息到nsq
// @Tags 消息管理
// @Accept json
// @Produce json
// @Param id path int true "用户ID"
// @Success 200 {object} UserResponse
// @Router /user/{id} [get]
func sendMessageHandler(c *gin.Context) {
	// 从请求体中获取消息内容
	message := c.PostForm("message")

	// 向 NSQ 发送消息
	err := producer.Publish("test", []byte(message))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to send message to NSQ",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Message sent to NSQ successfully",
	})
}
