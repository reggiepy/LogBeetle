package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// 处理首页请求
func HomeHandler(c *gin.Context) {
	c.Header("Content-Type", "text/html; charset=utf-8")
	c.String(http.StatusOK, "<h1>hello world.</h1>")
}

// ping
func PingHandler(c *gin.Context) {
	c.String(http.StatusOK, "pong")
}

// 处理关于页面请求
func AboutHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"title": "About Us",
	})
}
