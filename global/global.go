package global

import (
	"github.com/reggiepy/LogBeetle/config"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var (
	LbConfig config.Config
	LbViper  *viper.Viper
	LbLogger *zap.Logger
)
