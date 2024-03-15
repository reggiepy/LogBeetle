package web

import (
	"github.com/gin-gonic/gin"
	"github.com/reggiepy/LogBeetle/docs"
	"github.com/reggiepy/LogBeetle/middleware"
	"github.com/reggiepy/LogBeetle/pkg/logger"
	"github.com/reggiepy/LogBeetle/web/api"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"log"
)

// SetupRouter 设置路由
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
	r.GET("/log-beetle/v1/", api.HomeHandler)
	r.GET("/log-beetle/v1/ping", api.PingHandler)
	r.GET("/log-beetle/v1/about", api.AboutHandler)
	r.POST("/log-beetle/v1/send-message", api.SendMessageHandler)
	// use ginSwagger middleware to serve the API docs
	r.GET("/log-beetle/v1/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	return r
}
