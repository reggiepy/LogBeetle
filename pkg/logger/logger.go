package logger

import (
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"reflect"
)

var Logger *zap.Logger

type Config struct {
	LogFile    string `json:"LogFile" yaml:"LogFile"`       // 日志文件名
	MaxSize    int    `json:"MaxSize" yaml:"MaxSize"`       // 日志文件大小限制，单位为 MB
	MaxBackups int    `json:"MaxBackups" yaml:"MaxBackups"` // 最大保留的旧日志文件数量
	MaxAge     int    `json:"MaxAge" yaml:"MaxAge"`         // 旧日志文件保留天数
	Compress   bool   `json:"Compress" yaml:"Compress"`     // 是否压缩旧日志文件
	LogLevel   string `json:"LogLevel" yaml:"LogLevel"`     // 日志等级
	LogFormat  string `json:"LogFormat" yaml:"LogFormat"`   // 日志等级
}

type ConfigOption func(loggerConfig *Config) error

// WithLogFile 设置日志文件名
func WithLogFile(logFile string) ConfigOption {
	return func(loggerConfig *Config) error {
		loggerConfig.LogFile = logFile
		return nil
	}
}

// WithMaxSize 设置日志文件大小限制，单位为 MB
func WithMaxSize(maxSize int) ConfigOption {
	return func(loggerConfig *Config) error {
		loggerConfig.MaxSize = maxSize
		return nil
	}
}

// WithMaxBackups 设置最大保留的旧日志文件数量
func WithMaxBackups(maxBackups int) ConfigOption {
	return func(loggerConfig *Config) error {
		loggerConfig.MaxBackups = maxBackups
		return nil
	}
}

// WithMaxAge 设置旧日志文件保留天数
func WithMaxAge(maxAge int) ConfigOption {
	return func(loggerConfig *Config) error {
		loggerConfig.MaxAge = maxAge
		return nil
	}
}

// WithCompress 设置是否压缩旧日志文件
func WithCompress(compress bool) ConfigOption {
	return func(loggerConfig *Config) error {
		loggerConfig.Compress = compress
		return nil
	}
}

// WithLogLevel 设置日志等级
func WithLogLevel(logLevel string) ConfigOption {
	return func(loggerConfig *Config) error {
		loggerConfig.LogLevel = logLevel
		return nil
	}
}

// WithLogFormat 设置日志格式
func WithLogFormat(logFormat string) ConfigOption {
	return func(loggerConfig *Config) error {
		loggerConfig.LogFormat = logFormat
		return nil
	}
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

// InitLogger 初始化Logger
func InitLogger(logConfig Config, options ...ConfigOption) (err error) {
	if reflect.DeepEqual(logConfig, Config{}) {
		logConfig = *NewDefaultConfig()
	}
	for _, option := range options {
		err := option(&logConfig)
		if err != nil {
			panic(err)
		}
	}
	logLevel := map[string]zapcore.Level{
		"debug": zapcore.DebugLevel,
		"info":  zapcore.InfoLevel,
		"warn":  zapcore.WarnLevel,
		"error": zapcore.ErrorLevel,
	}
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
