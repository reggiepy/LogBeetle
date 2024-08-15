package public

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ApiMessage struct{}

// @Summary 发送消息
// @Description 发送消息到nsq
// @Tags 消息管理
// @Accept x-www-form-urlencoded
// @Produce json
//
// @Param			message			formData		string		true	"message"
// @Param			project_name	formData		string		true	"project_name"
//
// @Success      200  {object}   model.JSONResult
// @router      /log-beetle/v1/send-message   [post]
func (a *ApiMessage) SendMessageHandler(c *gin.Context) {
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

	err := serverPublic.ServiceMessage.SendMessage(projectName, message)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": fmt.Sprintf("%v", err),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Message send successfully",
	})
}
