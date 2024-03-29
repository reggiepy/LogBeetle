package logger

import (
	"github.com/natefinch/lumberjack"
	"github.com/reggiepy/LogBeetle/pkg/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger *zap.Logger

// InitLogger 初始化Logger
func InitLogger(cfg *config.Config) (err error) {
	logLevel := map[string]zapcore.Level{
		"debug": zapcore.DebugLevel,
		"info":  zapcore.InfoLevel,
		"warn":  zapcore.WarnLevel,
		"error": zapcore.ErrorLevel,
	}
	logConfig := cfg.LogConfig
	writeSyncer := getLogWriter(logConfig.LogFile, logConfig.MaxSize, logConfig.MaxBackups, logConfig.MaxAge, logConfig.Compress)

	logFormats := []string{"json", "logfmt"}
	logFormat := "json"
	for _, format := range logFormats {
		if format == logConfig.LogFormat {
			logFormat = format
		}
	}
	encoder := getEncoder(logFormat)
	l, ok := logLevel[logConfig.LogLevel] // 日志打印级别
	if !ok {
		l = logLevel["info"]
	}
	//l, err := zapcore.ParseLevel(logConfig.LogLevel)
	//if err != nil {
	//	l = zapcore.InfoLevel
	//}
	core := zapcore.NewCore(encoder, writeSyncer, l)

	Logger = zap.New(core, zap.AddCaller()) // zap.Addcaller() 输出日志打印文件和行数如： logger/logger_test.go:33
	// 1. zap.ReplaceGlobals 函数将当前初始化的 logger 替换到全局的 logger,
	// 2. 使用 logger 的时候 直接通过 zap.S().Debugf("xxx") or zap.L().Debug("xxx")
	// 3. 使用 zap.S() 和 zap.L() 提供全局锁，保证一个全局的安全访问logger的方式
	zap.ReplaceGlobals(Logger) // 替换zap包中全局的logger实例，后续在其他包中只需使用zap.L()调用即可
	//zap.L().Debug("")
	//zap.S().Debugf("")
	return
}

func getEncoder(logFormat string) zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.TimeKey = "time"
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoderConfig.EncodeDuration = zapcore.SecondsDurationEncoder
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	if logFormat == "json" {
		return zapcore.NewJSONEncoder(encoderConfig) // 以json格式写入
	}
	return zapcore.NewConsoleEncoder(encoderConfig) // 以logfmt格式写入
}

func getLogWriter(filename string, maxSize, maxBackup, maxAge int, compress bool) zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   filename,
		MaxSize:    maxSize,
		MaxBackups: maxBackup,
		MaxAge:     maxAge,
		Compress:   compress,
	}
	return zapcore.AddSync(lumberJackLogger)
}
