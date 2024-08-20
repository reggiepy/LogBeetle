package public

import (
	"github.com/reggiepy/LogBeetle/global"
	"github.com/reggiepy/LogBeetle/model"
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
	total := len(global.LBConsumerManager.Topics())
	dataList := make([]string, total)
	copy(dataList, global.LBConsumerManager.Topics())
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
