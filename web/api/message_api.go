package api

import (
	"github.com/gin-gonic/gin"
	"github.com/reggiepy/LogBeetle/nsqworker"
	"net/http"
)

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
	err := nsqworker.Producer.Publish("test", []byte(message))
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
