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
)
