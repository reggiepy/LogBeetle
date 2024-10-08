package boot

import (
	"fmt"
	"github.com/nsqio/go-nsq"
	"github.com/reggiepy/LogBeetle/global"
)

func NsqProducer() {
	var (
		err error
	)
	nsqConfig := global.LbConfig.NSQConfig
	config := nsq.NewConfig()
	err = config.Set("auth_secret", nsqConfig.AuthSecret)
	if err != nil {
		global.LbLogger.Info(fmt.Sprintf("Failed to set auth_secret %v", err))
	}
	producer, err := nsq.NewProducer(nsqConfig.NSQDAddress, config)
	if err != nil {
		global.LbLogger.Fatal(fmt.Sprintf("Failed to create producer: %v", err))
	}

	global.OnExit(func() {
		global.LbNsqProducer.Stop()
		global.LbLogger.Info("NSQ producer stopped")
	})

	global.LbNsqProducer = producer
}
