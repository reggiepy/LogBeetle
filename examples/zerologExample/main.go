package main

import (
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	// 打开日志文件，如果不存在则创建
	file, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal().Err(err).Msg("无法打开日志文件")
	}
	defer file.Close()

	// 创建 ConsoleWriter 以确保在 Windows 平台下的正确编码
	writer := zerolog.ConsoleWriter{
		Out:        file,
		TimeFormat: time.RFC3339,
		NoColor:    true,
	}

	// 将文件输出器附加到 zerolog
	log.Logger = log.Output(writer)

	// 写入日志
	log.Info().Msg("这是一条日志消息")
}
