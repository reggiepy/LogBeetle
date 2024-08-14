package public

import "github.com/gin-gonic/gin"

type RouterIndex struct{}

func (r *RouterIndex) InitIndexRouter(publicGroup *gin.RouterGroup) (R gin.IRoutes) {
	{
		publicGroup.GET("/", apiPublic.ApiIndex.HomeHandler)
		publicGroup.GET("/ping", apiPublic.ApiIndex.PingHandler)
		publicGroup.GET("/about", apiPublic.ApiIndex.AboutHandler)
	}
	return publicGroup
}
