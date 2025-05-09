package boot

import (
	"fmt"
	"github.com/nsqio/go-nsq"
	"github.com/reggiepy/LogBeetle/config"
	"github.com/reggiepy/LogBeetle/global"
	"github.com/reggiepy/LogBeetle/pkg/goutils/signailUtils"
)

func NsqProducer(nsqConfig config.NSQConfig) *nsq.Producer {
	var (
		err error
	)
	cfg := nsq.NewConfig()
	err = cfg.Set("auth_secret", nsqConfig.AuthSecret)
	if err != nil {
		global.LbLogger.Info(fmt.Sprintf("Failed to set auth_secret %v", err))
	}
	producer, err := nsq.NewProducer(nsqConfig.NSQDAddress, cfg)
	if err != nil {
		global.LbLogger.Fatal(fmt.Sprintf("Failed to create producer: %v", err))
	}

	signailUtils.OnExit(func() {
		producer.Stop()
		global.LbLogger.Info("NSQ producer stopped")
	})
	return producer
}
