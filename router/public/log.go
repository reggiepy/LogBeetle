package public

import "github.com/gin-gonic/gin"

type RouterLog struct{}

func (r *RouterLog) InitRouterLog(publicGroup *gin.RouterGroup) (R gin.IRoutes) {
	publicGroup = publicGroup.Group("log")
	{
		publicGroup.POST("/search", apiPublic.ApiLog.Search)
		publicGroup.POST("/addTestData", apiPublic.ApiLog.AddTestData)
		publicGroup.POST("/add", apiPublic.ApiLog.JsonLogAdd)
	}
	return publicGroup
}
