package public

import (
	"github.com/gin-gonic/gin"
	"github.com/reggiepy/LogBeetle/model"
)

type ApiSystem struct {
}

// @Summary 获取系统信息
// @Description 获取系统信息
// @Tags 系统
// @Accept json
// @Produce json
//
// @Param request body model.SystemInfoRequest true "请求参数"
//
// @Success      200  {object}   model.JSONResult
// @Router /log-beetle/v1/system/system-info [post]
func (s *ApiSystem) GetSystemInfo(c *gin.Context) {
	err, data := serverPublic.ServiceSystem.SystemInfo()
	if err != nil {
		return
	}
	model.ResponseSuccessData(data, c)
}
