package global

import (
	"github.com/nsqio/go-nsq"
	"github.com/reggiepy/LogBeetle/config"
	"github.com/reggiepy/LogBeetle/pkg/consumer/manager"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"time"
)

var (
	LbStartTime = time.Now().Local()
	LbConfig    config.Config
	LbViper     *viper.Viper
	LbLogger    *zap.Logger

	// NSQ
	LbNsqProducer     *nsq.Producer
	LBConsumerManager *manager.Manager

	LbExitFunc []func()

	//	系统信息
	StartTime = time.Now().Format("2006-01-02 15:04:05")
)

func OnExit(f ...func()) {
	LbExitFunc = append(LbExitFunc, f...)
}

func Exit() {
	for _, fnExit := range LbExitFunc {
		fnExit()
	}
}
