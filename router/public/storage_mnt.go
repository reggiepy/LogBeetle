package public

import "github.com/gin-gonic/gin"

type RouterStorageMnt struct{}

func (r *RouterStorageMnt) InitRouterStorageMntRouter(publicGroup *gin.RouterGroup) (R gin.IRoutes) {
	publicGroup = publicGroup.Group("store")
	{
		publicGroup.POST("/names", apiPublic.ApiStorageMnt.Names)
		publicGroup.POST("/list", apiPublic.ApiStorageMnt.List)
		publicGroup.POST("/delete", apiPublic.ApiStorageMnt.Delete)
	}
	return publicGroup
}
