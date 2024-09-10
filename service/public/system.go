package public

import (
	"github.com/reggiepy/LogBeetle/global"
	"github.com/reggiepy/LogBeetle/model"
	"runtime"
)

type ServiceSystem struct {
}

func (s *ServiceSystem) SystemInfo() (error, model.SystemInfoResponse) {
	ret := model.SystemInfoResponse{
		ConsumerInfo: model.ConsumerInfo{
			ConsumerCount: global.LBConsumerManager.Count(),
		},
		StartTime:       global.StartTime,
		GoroutineNumber: runtime.NumGoroutine(),
	}
	return nil, ret
}
