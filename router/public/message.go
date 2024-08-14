package public

import "github.com/gin-gonic/gin"

type RouterMessage struct{}

func (r *RouterMessage) InitRouterMessage(publicGroup *gin.RouterGroup) (R gin.IRoutes) {
	{
		publicGroup.POST("/send-message", apiPublic.ApiMessage.SendMessageHandler)
	}
	return publicGroup
}
