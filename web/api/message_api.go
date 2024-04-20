package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/reggiepy/LogBeetle/pkg/producer/nsqproducer"
)

// @Summary 发送消息
// @Description 发送消息到nsq
// @Tags 消息管理
// @Accept x-www-form-urlencoded
// @Produce json
//
//	@Param			message			formData		string		true	"message"
//	@Param			project_name	formData		string		true	"project_name"
//
// @Success      200  {object}   model.JSONResult
// @router      /log-beetle/v1/send-message   [post]
func SendMessageHandler(c *gin.Context) {
	// 从请求体中获取消息内容
	message := c.DefaultPostForm("message", "")
	if message == "" {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Message cannot be empty",
		})
		return
	}
	// 从请求体中获取消息内容
	projectName := c.DefaultPostForm("project_name", "test")
	if projectName == "" {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "project_name cannot be empty",
		})
		return
	}
	start := time.Now()
	// 向 NSQ 发送消息
	err := nsqproducer.Producer.Publish(projectName, []byte(message))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to send message",
		})
		return
	}
	end := time.Now()
	elapsed := end.Sub(start)
	fmt.Printf("topic【%s】消息写入时间：%s\n", projectName, elapsed)

	c.JSON(http.StatusOK, gin.H{
		"message": "Message send successfully",
	})
}
