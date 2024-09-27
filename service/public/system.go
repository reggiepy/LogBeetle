package public

import (
	"github.com/reggiepy/LogBeetle/global"
	"github.com/reggiepy/LogBeetle/model"
	"github.com/reggiepy/LogBeetle/version"
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
		Version:         version.Full(),
	}
	return nil, ret
}
