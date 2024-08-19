package global

import (
	"github.com/nsqio/go-nsq"
	"github.com/reggiepy/LogBeetle/config"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"time"
)

var (
	LbStartTime   = time.Now().Local()
	LbConfig      config.Config
	LbViper       *viper.Viper
	LbLogger      *zap.Logger
	LbNsqProducer *nsq.Producer

	//	注册的topic
	LbRegisterTopic []string
)
