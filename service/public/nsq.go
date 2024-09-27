package public

import (
	"github.com/reggiepy/LogBeetle/global"
	"github.com/reggiepy/LogBeetle/model"
	"github.com/reggiepy/LogBeetle/pkg/consumer"
	"sort"
)

type ServiceNsq struct {
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
func (s *ServiceNsq) RegisterTopic(requestNsqTopicList model.RequestGetNsqTopicList) (error, []string, int64) {
	pageSize := requestNsqTopicList.PageSize
	offset := requestNsqTopicList.PageSize * max(requestNsqTopicList.Page-1, 0)
	dataList := global.LBConsumerManager.GetNamesByType(consumer.NSQConsumer.String())
	total := len(dataList)
	if requestNsqTopicList.Desc {
		sort.Sort(sort.Reverse(sort.StringSlice(dataList)))
	} else {
		sort.Strings(dataList)
	}
	// 计算分页的切片范围
	if offset >= total {
		dataList = []string{}
	} else if offset+pageSize > total {
		dataList = dataList[offset:]
	} else {
		dataList = dataList[offset : offset+pageSize]
	}
	return nil, dataList, int64(total)
}
