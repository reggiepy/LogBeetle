package web

import (
	"github.com/gin-gonic/gin"
	"github.com/reggiepy/LogBeetle/docs"
	"github.com/reggiepy/LogBeetle/middleware"
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
	router := gin.Default()
	gin.DebugPrintRouteFunc = func(httpMethod, absolutePath, handlerName string, nuHandlers int) {
		log.Printf("endpoint %v %v %v %v\n", httpMethod, absolutePath, handlerName, nuHandlers)
	}
	// 注册全局中间件
	router.Use(middleware.ErrorHandler())
	router.Use(middleware.DddRequestID())
	router.Use(middleware.Cors())

	// 定义路由和处理程序
	router.GET("/log-beetle/v1/", api.HomeHandler)
	router.GET("/log-beetle/v1/ping", api.PingHandler)
	router.GET("/log-beetle/v1/about", api.AboutHandler)
	router.POST("/log-beetle/v1/send-message", api.SendMessageHandler)
	// use ginSwagger middleware to serve the API docs
	router.GET("/log-beetle/v1/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	return router
}
