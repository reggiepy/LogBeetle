package middleware

import (
	"fmt"
	"go.uber.org/zap"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ErrorHandler 是一个中间件函数，用于捕获异常并返回适当的错误响应
func ErrorHandler(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// 捕获异常并记录日志
				logger.Info(fmt.Sprintf("Request processing error: %v", err))

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
