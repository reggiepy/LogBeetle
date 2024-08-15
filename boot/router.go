package boot

import (
	"github.com/gin-gonic/gin"
	"github.com/reggiepy/LogBeetle/docs"
	"github.com/reggiepy/LogBeetle/global"
	"github.com/reggiepy/LogBeetle/middleware"
	"github.com/reggiepy/LogBeetle/router"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// SetupRouter 设置路由
// @title           Swagger Example API
// @version         1.0
// @description     This is a sample server celler server.
// @termsOfService  http://swagger.io/terms/
// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io
// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html
// @host      localhost:8080
// @BasePath  /api/v1
// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/
func Router() *gin.Engine {
	// programatically set swagger info
	docs.SwaggerInfo.Title = "Log API"
	docs.SwaggerInfo.Description = "This is a sample server Petstore server."
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "127.0.0.1:1233"
	docs.SwaggerInfo.BasePath = "/"
	docs.SwaggerInfo.Schemes = []string{"http", "https"}
	if global.LbConfig.Env == "dev" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	// r := gin.Default()
	r := gin.New()

	// 注册全局中间件
	r.Use(middleware.ErrorHandler())
	r.Use(middleware.RequestID())
	r.Use(middleware.Cors())

	group := r.Group("log-beetle")
	// Public route: starts with "public". There is no need to regroup or authenticate within the route
	PublicGroup := group.Group("v1")
	RouterPublic(PublicGroup)

	// Private route：starts with "". The route is grouped according to the actual performance, and authentication is required
	PrivateGroup := group.Group("")
	// 注册全局中间件
	RouterPrivate(PrivateGroup)
	if global.LbConfig.Env == "dev" {
		// use ginSwagger middleware to serve the API docs
		r.GET("/log-beetle/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}
	return r
}

func RouterPublic(PublicGroup *gin.RouterGroup) {
	routerPublic := router.LbRouter.RouterPublic
	{
		routerPublic.RouterMessage.InitRouterMessage(PublicGroup)
		routerPublic.RouterIndex.InitIndexRouter(PublicGroup)
		routerPublic.RouterNsq.InitRouterNsq(PublicGroup)
	}
}

func RouterPrivate(PrivateGroup *gin.RouterGroup) {
	_ = router.LbRouter.RouterPrivate
	//{
	//}
}
