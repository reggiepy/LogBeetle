package api

import (
	"fmt"
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
	err = config.Set("auth_secret", "%n&yFA2JD85z^g")
	if err != nil {
		fmt.Printf("Failed to set auth_secret %v", err)
	}
	producer, err = nsq.NewProducer("192.168.1.110:4150", config)
	if err != nil {
		log.Fatalf("Failed to create NSQ producer: %v", err)
	}
}

// @Summary 发送消息
// @Description 发送消息到nsq
// @Tags 消息管理
// @Accept x-www-form-urlencoded
// @Produce json
//
//	@Param			message	formData		string		true	"message"
//
// @Success      200  {object}   model.JSONResult
// @router      /send-message   [post]
func SendMessageHandler(c *gin.Context) {
	// 从请求体中获取消息内容
	message := c.DefaultPostForm("message", "")
	if message == "" {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Message cannot be empty",
		})
		return
	}
	// 向 NSQ 发送消息
	err := producer.Publish("test", []byte(message))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to send message to NSQ",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Message sent to NSQ successfully",
	})
}
