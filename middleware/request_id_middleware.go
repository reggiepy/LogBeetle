package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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
