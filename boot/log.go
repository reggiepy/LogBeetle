package boot

import (
	"encoding/json"
	"fmt"
	"github.com/reggiepy/LogBeetle/global"
	"github.com/reggiepy/LogBeetle/pkg/logger"
	"os"
)

func Log() {
	logConfig := &logger.Config{}
	jsonBytes, _ := json.Marshal(global.LbConfig.LogConfig)
	err := json.Unmarshal(jsonBytes, logConfig)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	err = logger.InitLogger(*logConfig)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
