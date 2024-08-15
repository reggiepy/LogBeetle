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
// @Accept json
// @Produce json
//
// @Param request body model.RequestGetNsqTopicList true "请求参数"
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
