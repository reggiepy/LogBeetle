package public

import "github.com/gin-gonic/gin"

type RouterSystem struct {}


func (r *RouterSystem) InitRouterSystem (publicGroup *gin.RouterGroup) (R gin.IRoutes) {
	publicGroup = publicGroup.Group("system")
	{
		publicGroup.POST("/system-info", apiPublic.ApiSystem.GetSystemInfo)
	}
	return publicGroup
}