package public

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ApiIndex struct{}

// @Summary 首页
// @Description 首页
// @Tags 管理
// @Accept plain
// @Produce plain
// @Success      200  {object}   model.AboutResponse
// @router      /log-beetle/v1/   [get]
func (a *ApiIndex) HomeHandler(c *gin.Context) {
	c.Header("Content-Type", "text/html; charset=utf-8")
	c.String(http.StatusOK, "<h1>hello world.</h1>")
}

// @Summary ping
// @Description ping
// @Tags 管理
// @Accept plain
// @Produce plain
// @Success      200  {string}   string
// @router      /log-beetle/v1/ping   [get]
func (a *ApiIndex) PingHandler(c *gin.Context) {
	c.String(http.StatusOK, "pong")
}

// @Summary 关于
// @Description 关于
// @Tags 管理
// @Accept plain
// @Produce json
// @Success      200  {string}   model.AboutResponse
// @router      /log-beetle/v1/about   [get]
func (a *ApiIndex) AboutHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"title": "About Us",
	})
}
