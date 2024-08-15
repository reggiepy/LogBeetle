package public

import "github.com/gin-gonic/gin"

type RouterNsq struct {
}

func (r *RouterNsq) InitRouterNsq(publicGroup *gin.RouterGroup) (R gin.IRoutes) {
	publicGroup = publicGroup.Group("nsq")
	{
		publicGroup.POST("/register-topic", apiPublic.ApiNsq.RegisterTopic)
	}
	return publicGroup
}
