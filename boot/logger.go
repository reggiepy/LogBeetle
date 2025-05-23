package boot

import (
	"encoding/json"
	"fmt"
	"github.com/reggiepy/LogBeetle/global"
	"github.com/reggiepy/goutils/logutil/zapLogger"
	"go.uber.org/zap"
)

func Logger() *zap.Logger {
	logConfig := zapLogger.NewLoggerConfig(
		zapLogger.WithInConsole(true),
		zapLogger.WithReplaceGlobals(true),
	)
	jsonBytes, _ := json.Marshal(global.LbConfig.LogConfig)
	err := logConfig.LoadJSON(string(jsonBytes))
	if err != nil {
		fmt.Printf("Error marshalling log config, use default config: %v\n", err)
	}
	//fmt.Println("Log Config: ", logConfig.ToJSON())
	logger, cleanup := zapLogger.NewLogger(logConfig)
	global.LbLoggerClearup = cleanup
	return logger
}
