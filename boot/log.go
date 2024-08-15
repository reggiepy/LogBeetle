package boot

import (
	"encoding/json"
	"fmt"
	"github.com/reggiepy/LogBeetle/global"
	"github.com/reggiepy/LogBeetle/goutils/logUtils"
	"os"
)

func Log() {
	logConfig := &logUtils.Config{}
	jsonBytes, _ := json.Marshal(global.LbConfig.LogConfig)
	err := json.Unmarshal(jsonBytes, logConfig)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	err = logUtils.InitLogger(*logConfig)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
