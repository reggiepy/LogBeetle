package main

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

func main() {
	// 创建lumberjack日志切割器
	lumberjackLogger := &lumberjack.Logger{
		Filename:   "app.log",
		MaxSize:    500,  // 每个日志文件的最大大小，单位为MB
		MaxBackups: 3,    // 保留的旧日志文件的最大数量
		MaxAge:     28,   // 保留的旧日志文件的最大天数
		Compress:   true, // 是否压缩旧日志文件
	}

	// 创建Zap核心
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
		zapcore.AddSync(lumberjackLogger),
		zapcore.InfoLevel,
	)

	// 创建Logger
	logger := zap.New(core)

	// 记录日志
	logger.Info("这是一条日志消息")

	// 刷新缓冲区并关闭日志
	logger.Sync()
}
