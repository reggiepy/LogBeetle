package consumer

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io"
)

// NewZAPLogger 创建一个新的 zap.Logger 实例，使用指定的写入器。
// 它配置编码器以使用 ISO8601TimeEncoder 格式的时间和 CapitalLevelEncoder 格式的日志级别。
// 日志记录器的核心设置为使用 JSON 编码器和记录 InfoLevel 消息。
func NewZAPLogger(w io.Writer) *zap.Logger {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoderConfig.TimeKey = "time"
	encoderConfig.EncodeDuration = zapcore.SecondsDurationEncoder
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		zapcore.AddSync(w),
		zapcore.InfoLevel,
	)
	return zap.New(core)
}
