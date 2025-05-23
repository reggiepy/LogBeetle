package public

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/reggiepy/LogBeetle/com"
	"github.com/reggiepy/LogBeetle/ldb/search"
	"github.com/reggiepy/LogBeetle/ldb/storage/logdata"
	"github.com/reggiepy/LogBeetle/model"
	"net/http"
	"strings"
)

type ApiLog struct{}

// @Summary 搜索日志
// @Description 日志搜索
// @Tags 日志管理
// @Accept json
// @Produce json
//
// @Param request body model.SearchRequest true "请求参数"
//
// @Success      200  {object}   model.JSONResult
// @router      /log-beetle/v1/log/search   [post]
func (a *ApiLog) Search(c *gin.Context) {
	var l model.SearchRequest
	if err := model.RequestShouldBindJSON(c, &l); err != nil {
		return
	}

	cond := &search.SearchCondition{
		StoreName:        l.StoreName,                              // 日志仓条件
		SearchKey:        l.SearchKey,                              // 输入的查询关键词
		CurrentStoreName: l.CurrentStoreName,                       // 滚动查询时定位用日志仓
		CurrentId:        com.StringToUint32(l.CurrentId, 0),       // 滚动查询时定位用ID
		Forward:          com.StringToBool(l.Forward, true),        // 是否向下滚动查询
		OldNearId:        com.StringToUint32(l.OldNearId, 0),       // 相邻检索旧ID
		NewNearId:        com.StringToUint32(l.NewNearId, 0),       // 相邻检索新ID
		NearStoreName:    l.NearStoreName,                          // 相邻检索时新ID对应的日志仓
		DatetimeFrom:     l.DatetimeFrom,                           // 日期范围（From）
		DatetimeTo:       l.DatetimeTo,                             // 日期范围（To）
		OrgSystem:        com.ToLower(strings.TrimSpace(l.System)), // 系统
		User:             com.ToLower(strings.TrimSpace(l.User)),   // 用户
		Loglevel:         com.ToLower(l.LogLevel),                  // 日志级别
	}
	cond.Loglevels = com.Split(cond.Loglevel, ",") // 多选条件
	data := serverPublic.LogService.Search(cond)
	model.ResponseSuccessData(data, c)
}

// @Summary 添加测试日志
// @Description 添加测试日志
// @Tags 日志管理
// @Accept json
// @Produce json
//
// @Param request body string false "请求参数"
//
// @Success      200  {object}   model.JSONResult
// @router      /log-beetle/v1/log/addTestData   [post]
func (a *ApiLog) AddTestData(c *gin.Context) {
	err := serverPublic.LogService.AddTestData()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": fmt.Sprintf("%v", err),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Message send successfully",
	})
}

// @Summary 添加日志
// @Description 添加日志
// @Tags 日志管理
// @Accept json
// @Produce json
//
// @Param request body model.AddLogRequest true "请求参数"
//
// @Success      200  {object}   model.JSONResult
// @router      /log-beetle/v1/log/add   [post]
func (a *ApiLog) JsonLogAdd(c *gin.Context) {
	var l logdata.LogDataModel
	if err := model.RequestShouldBindJSON(c, &l); err != nil {
		return
	}
	err := serverPublic.LogService.JsonLogAdd(&l)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": fmt.Sprintf("%v", err),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Message send successfully",
	})
}
