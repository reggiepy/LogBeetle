package main

import (
	"os"

	"github.com/rs/zerolog"
)

func main() {
	// 创建一个新的 zerolog 日志记录器，并设置输出到标准输出
	logger := zerolog.New(os.Stdout)

	// 配置日志的输出格式为 JSON 格式
	logger = logger.Output(zerolog.ConsoleWriter{Out: os.Stdout})

	// 输出日志
	logger.Info().Str("key", "value").Msg("logging with zerolog")
}
