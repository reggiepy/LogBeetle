package logUtils

import (
	"github.com/natefinch/lumberjack"
	"github.com/reggiepy/LogBeetle/goutils/arrayUtils"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Config struct {
	LogFile    string `json:"LogFile" yaml:"LogFile"`       // 日志文件名
	MaxSize    int    `json:"MaxSize" yaml:"MaxSize"`       // 日志文件大小限制，单位为 MB
	MaxBackups int    `json:"MaxBackups" yaml:"MaxBackups"` // 最大保留的旧日志文件数量
	MaxAge     int    `json:"MaxAge" yaml:"MaxAge"`         // 旧日志文件保留天数
	Compress   bool   `json:"Compress" yaml:"Compress"`     // 是否压缩旧日志文件
	LogLevel   string `json:"LogLevel" yaml:"LogLevel"`     // 日志等级
	LogFormat  string `json:"LogFormat" yaml:"LogFormat"`   // 日志等级
}

// NewDefaultConfig 创建默认配置
func NewDefaultConfig() *Config {
	return &Config{
		LogFile:    "app.log",
		MaxSize:    1,
		MaxBackups: 5,
		MaxAge:     30,
		Compress:   false,
		LogLevel:   "info",
		LogFormat:  "json",
	}
}

// NewLogger 初始化Logger
func (c *Config) NewLogger(opts ...Option) (*zap.Logger, error) {
	var (
		logger     *zap.Logger
		logFormats = arrayUtils.NewSet("json", "logfmt")
	)
	for _, opt := range opts {
		opt.apply(c)
	}

	lumberJackLogger := &lumberjack.Logger{
		Filename:   c.LogFile,
		MaxSize:    c.MaxSize,
		MaxBackups: c.MaxBackups,
		MaxAge:     c.MaxAge,
		Compress:   c.Compress,
	}
	writeSyncer := zapcore.AddSync(lumberJackLogger)

	logFormat := "json"
	for logFormats.Has(c.LogFormat) {
		logFormat = c.LogFormat
	}
	encoder := NewEncoder(logFormat)

	// 日志打印级别
	l, err := zapcore.ParseLevel(c.LogLevel)
	if err != nil {
		l = zapcore.InfoLevel
	}
	core := zapcore.NewCore(encoder, writeSyncer, l)

	logger = zap.New(core, zap.AddCaller()) // zap.Addcaller() 输出日志打印文件和行数如： logger/logger_test.go:33
	// 1. zap.ReplaceGlobals 函数将当前初始化的 logger 替换到全局的 logger,
	// 2. 使用 logger 的时候 直接通过 zap.S().Debugf("xxx") or zap.L().Debug("xxx")
	// 3. 使用 zap.S() 和 zap.L() 提供全局锁，保证一个全局的安全访问logger的方式
	// zap.ReplaceGlobals(logger) // 替换zap包中全局的logger实例，后续在其他包中只需使用zap.L()调用即可
	// zap.L().Debug("")
	// zap.S().Debugf("")
	return logger, nil
}

func NewEncoder(logFormat string) zapcore.Encoder {
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
