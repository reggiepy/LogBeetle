package web

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"log"
	"net/http"
)

// DddRequestID 是一个中间件函数，用于给每个请求添加一个request id
func DddRequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 生成一个唯一的request id，并将其存储到Gin的上下文中
		requestID := uuid.New().String()
		c.Set("requestID", requestID)

		// 将request id 添加到HTTP请求的header中
		c.Header("X-Request-ID", requestID)

		// 执行下一个中间件或请求处理函数
		c.Next()
	}
}

// GetRequestID 是一个helper函数，用于从Gin的上下文中获取request id
func GetRequestID(c *gin.Context) string {
	requestID, _ := c.Get("requestID")
	return requestID.(string)
}

// ErrorHandler 是一个中间件函数，用于捕获异常并返回适当的错误响应
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// 捕获异常并记录日志
				log.Println("Request processing error:", err)

				// 返回 500 内部服务器错误响应
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "Internal Server Error",
				})
			}
		}()

		// 继续执行下一个中间件或请求处理函数
		c.Next()
	}
}