package main

import (
	"time"

	"github.com/natefinch/lumberjack"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	// 创建 lumberjack.Logger 实例用于日志切割
	logFile := &lumberjack.Logger{
		Filename:   "app.log", // 日志文件名
		MaxSize:    10,        // 日志文件大小限制，单位为 MB
		MaxBackups: 5,         // 最大保留的旧日志文件数量
		MaxAge:     30,        // 旧日志文件保留天数
		Compress:   true,      // 是否压缩旧日志文件
	}

	// 创建 ConsoleWriter 以确保在 Windows 平台下的正确编码
	writer := zerolog.ConsoleWriter{
		Out:        logFile,
		TimeFormat: time.RFC3339,
		NoColor:    true,
	}

	// 将文件输出器附加到 zerolog
	log.Logger = log.Output(writer)

	// 写入日志
	log.Info().Msg("这是一条日志消息")
}
