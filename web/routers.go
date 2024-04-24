package web

import (
	"github.com/gin-gonic/gin"
	"github.com/reggiepy/LogBeetle/docs"
	"github.com/reggiepy/LogBeetle/middleware"
	"github.com/reggiepy/LogBeetle/pkg/config"
	"github.com/reggiepy/LogBeetle/pkg/logger"
	"github.com/reggiepy/LogBeetle/web/api"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"log"
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
func SetupRouter() *gin.Engine {
	// programatically set swagger info
	docs.SwaggerInfo.Title = "Log API"
	docs.SwaggerInfo.Description = "This is a sample server Petstore server."
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "127.0.0.1:1233"
	docs.SwaggerInfo.BasePath = "/"
	docs.SwaggerInfo.Schemes = []string{"http", "https"}
	// 创建路由引擎
	//r := gin.Default()
	if config.Instance.Env == "dev" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.New()
	r.Use(middleware.GinLogger(logger.Logger), middleware.GinRecovery(logger.Logger, true))
	gin.DebugPrintRouteFunc = func(httpMethod, absolutePath, handlerName string, nuHandlers int) {
		log.Printf("endpoint %v %v %v %v\n", httpMethod, absolutePath, handlerName, nuHandlers)
	}
	// 注册全局中间件
	r.Use(middleware.ErrorHandler())
	r.Use(middleware.DddRequestID())
	r.Use(middleware.Cors())

	// 定义路由和处理程序
	r.GET("/log-beetle/", api.HomeHandler)
	r.GET("/log-beetle/ping", api.PingHandler)
	r.GET("/log-beetle/about", api.AboutHandler)
	r.POST("/log-beetle/v1/send-message", api.SendMessageHandler)
	if config.Instance.Env == "dev" {
		// use ginSwagger middleware to serve the API docs
		r.GET("/log-beetle/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}
	return r
}
