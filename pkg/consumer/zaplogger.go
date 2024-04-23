package consumer

import (
	"github.com/natefinch/lumberjack"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io"
	"time"
)

// NewZAPLogger 创建一个新的 zap.Logger 实例，使用指定的写入器。
// 它配置编码器以使用 ISO8601TimeEncoder 格式的时间和 CapitalLevelEncoder 格式的日志级别。
// 日志记录器的核心设置为使用 JSON 编码器和记录 InfoLevel 消息。
func NewZAPLogger(w io.Writer) *zap.Logger {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		zapcore.AddSync(w),
		zapcore.InfoLevel,
	)
	return zap.New(core)
}

// NewZEROLogger 创建一个新的 zerolog.Logger 实例，使用指定的写入器。
// 它创建一个 ConsoleWriter 以确保在 Windows 平台下的正确编码。
// 日志记录器的输出设置为写入器，时间格式设置为 RFC3339，并禁用颜色。
func NewZEROLogger(w io.Writer) zerolog.Logger {
	writer := zerolog.ConsoleWriter{
		Out:        w,
		TimeFormat: time.RFC3339,
		NoColor:    true,
	}
	return log.Output(writer)
}

// NewLJLoggerWriteCloser 创建一个新的 lumberjack.Logger 实例，使用指定的文件路径。
// 它配置日志记录器的文件路径、最大大小、最大备份数量、最大保存天数和是否压缩。
// 返回日志记录器实例的指针。
func NewLJLoggerWriteCloser(filePath string) *lumberjack.Logger {
	return &lumberjack.Logger{
		Filename:   filePath, // 日志文件名
		MaxSize:    1,        // 日志文件大小限制，单位为 MB
		MaxBackups: 5,        // 最大保留的旧日志文件数量
		MaxAge:     30,       // 旧日志文件保留天数
		Compress:   false,    // 是否压缩旧日志文件
	}
}
