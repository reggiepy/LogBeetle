package public

import (
	"github.com/gin-gonic/gin"
	"github.com/reggiepy/LogBeetle/model"
)

type ApiStorageMnt struct{}

// @Summary 查询日志仓名称列表
// @Description 查询日志仓名称列表
// @Tags 日志存储管理
// @Accept json
// @Produce json
//
// @Success      200  {object}   model.JSONResult
// @router      /log-beetle/v1/store/names   [post]
func (a *ApiStorageMnt) Names(c *gin.Context) {
	data := serverPublic.StorageMntService.Names()
	model.ResponseSuccessData(data, c)
}

// @Summary 查询日志仓信息列表
// @Description 查询日志仓信息列表
// @Tags 日志存储管理
// @Accept json
// @Produce json
//
// @Success      200  {object}   model.JSONResult
// @router      /log-beetle/v1/store/list   [post]
func (a *ApiStorageMnt) List(c *gin.Context) {
	data := serverPublic.StorageMntService.List()
	model.ResponseSuccessData(data, c)
}

// @Summary 删除指定日志仓
// @Description 删除指定日志仓
// @Tags 日志存储管理
// @Accept json
// @Produce json
//
// @Param request body model.DeleteStoreRequest true "请求参数"
//
// @Success      200  {object}   model.JSONResult
// @router      /log-beetle/v1/store/delete   [post]
func (a *ApiStorageMnt) Delete(c *gin.Context) {
	var l model.DeleteStoreRequest
	if err := model.RequestShouldBindJSON(c, &l); err != nil {
		return
	}
	data := serverPublic.StorageMntService.Delete(l.StoreName)
	model.ResponseSuccessData(data, c)
}
