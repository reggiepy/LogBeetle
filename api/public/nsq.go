package public

import (
	"github.com/gin-gonic/gin"
	"github.com/reggiepy/LogBeetle/model"
)

type ApiNsq struct {
}

// @Summary 列出 NSQ 注册的 topic
// @Description 发送消息到 NSQ
// @Tags NSQ
// @Accept x-www-form-urlencoded
// @Produce json
//
// @Param page formData int false "page" default(1)
// @Param pageSize formData int false "pageSize" default(10)
// @Param sortBy formData string false "sortBy"
// @Param desc formData bool false "desc" default(false)
//
// @Success      200  {object}   model.JSONResult
// @Router /log-beetle/v1/nsq/register-topic [post]
func (a *ApiNsq) RegisterTopic(c *gin.Context) {
	var toGetDataList model.RequestGetNsqTopicList
	if err := model.RequestShouldBindJSON(c, &toGetDataList); err != nil {
		return
	}
	err, dataList, total := serverPublic.ServiceNsq.RegisterTopic(toGetDataList)
	if err != nil {
		return
	}
	model.ResponseSuccessData(model.ResponsePageWithParentId{
		Records:  dataList,
		Page:     toGetDataList.Page,
		PageSize: toGetDataList.PageSize,
		Total:    total,
	}, c)
}
