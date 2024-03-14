package main

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

func main() {
	// 创建文件输出器
	file, err := os.Create("app.log")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// 设置文件输出配置
	encoderConfig := zapcore.EncoderConfig{
		MessageKey:     "msg",
		LevelKey:       "level",
		TimeKey:        "time",
		NameKey:        "logger",
		CallerKey:      "caller",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder, // 使用ISO8601格式的时间
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.FullCallerEncoder,
	}

	// 创建一个encoder
	encoder := zapcore.NewJSONEncoder(encoderConfig)

	// 设置日志级别
	atomicLevel := zap.NewAtomicLevel()
	atomicLevel.SetLevel(zapcore.InfoLevel)

	// 创建一个文件写入器
	core := zapcore.NewCore(encoder, zapcore.AddSync(file), atomicLevel)

	// 创建Logger
	logger := zap.New(core)

	// 记录日志
	logger.Info("这是一条日志消息")

	// 刷新缓冲区并关闭日志
	logger.Sync()
}
