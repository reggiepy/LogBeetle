package boot

import (
	"encoding/json"
	"fmt"
	"github.com/reggiepy/LogBeetle/global"
	"github.com/reggiepy/LogBeetle/goutils/logUtils"
	"go.uber.org/zap"
	"os"
)

func Log() *zap.Logger {
	logConfig := logUtils.NewDefaultConfig()
	jsonBytes, _ := json.Marshal(global.LbConfig.LogConfig)
	err := json.Unmarshal(jsonBytes, logConfig)
	if err != nil {
		fmt.Printf("Error marshalling log config: %v\n", err)
		os.Exit(1)
	}
	logger, err := logConfig.NewLogger()
	if err != nil {
		fmt.Printf("Error creating logger: %v\n", err)
		os.Exit(1)
	}
	return logger
}
