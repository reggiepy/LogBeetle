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

	lj := &lumberjack.Logger{
		Filename:   c.LogFile,
		MaxSize:    c.MaxSize,
		MaxBackups: c.MaxBackups,
		MaxAge:     c.MaxAge,
		Compress:   c.Compress,
	}

	logFormat := "json"
	if logFormats.Has(c.LogFormat) {
		logFormat = c.LogFormat
	}
	// 日志打印级别
	l, err := zapcore.ParseLevel(c.LogLevel)
	if err != nil {
		l = zapcore.InfoLevel
	}
	core := zapcore.NewCore(
		NewEncoder(logFormat),
		zapcore.AddSync(lj),
		l,
	)

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
	// 创建一个自定义的 EncoderConfig 实例
	encoderConfig := zap.NewProductionEncoderConfig()

	// 设置时间编码格式为 ISO 8601
	// 这将以标准的 ISO 8601 格式输出时间，例如 "2024-08-15T10:00:00Z"
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	// 设置日志级别编码格式为大写
	// 这将把日志级别输出为大写字母，例如 INFO、ERROR
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder

	// 设置时间字段的键名为 "time"
	// 在日志输出中，时间字段将使用 "time" 作为键名
	encoderConfig.TimeKey = "time"

	// 设置持续时间的编码格式为秒
	// 这将把持续时间格式化为秒数
	encoderConfig.EncodeDuration = zapcore.SecondsDurationEncoder

	// 设置调用者信息的编码格式为简短格式
	// 这将以相对路径和行号输出调用者信息，例如 "main.go:42"
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	if logFormat == "json" {
		return zapcore.NewJSONEncoder(encoderConfig) // 以json格式写入
	}
	return zapcore.NewConsoleEncoder(encoderConfig) // 以logfmt格式写入
}
