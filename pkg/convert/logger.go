package convert

import (
	"encoding/json"
	"github.com/reggiepy/LogBeetle/pkg/config"
	"github.com/reggiepy/LogBeetle/pkg/logger"
)

func ConfigToLoggerConfig(config *config.Config) (*logger.Config, error) {
	var ret = &logger.Config{}
	jsonBytes, _ := json.Marshal(config.LogConfig)
	err := json.Unmarshal(jsonBytes, ret)
	if err != nil {
		return nil, err
	}
	return ret, nil
}
