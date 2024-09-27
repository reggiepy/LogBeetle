package zapLogger

import (
	"os"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// NewLogger 初始化Logger
func NewLogger(config *LoggerConfig) (*zap.Logger, func()) {
	var cores []zapcore.Core

	// 支持的日志格式
	supportedLogFormats := map[string]struct{}{
		"json":   {},
		"logfmt": {},
	}

	// 设置日志格式，支持 "json" 和 "logfmt"
	logFormat := "json"
	if _, exists := supportedLogFormats[config.Format]; exists {
		logFormat = config.Format
	}

	// 设置日志级别
	level, err := zapcore.ParseLevel(config.Level)
	if err != nil {
		level = zapcore.InfoLevel
	}

	var lj *lumberjack.Logger
	// 创建文件日志输出 core
	if config.InFile {
		lj = &lumberjack.Logger{
			Filename:   config.File,
			MaxSize:    config.MaxSize,    // 文件大小 (MB)
			MaxBackups: config.MaxBackups, // 保留的最大文件个数
			MaxAge:     config.MaxAge,     // 保留的最大天数
			Compress:   config.Compress,   // 是否压缩
		}

		fileCore := zapcore.NewCore(
			NewEncoder(logFormat),
			zapcore.AddSync(lj), // 将日志写入文件
			level,
		)
		cores = append(cores, fileCore)
	}

	// 创建终端日志输出 core
	if config.InConsole {
		consoleCore := zapcore.NewCore(
			NewEncoder(logFormat),
			zapcore.AddSync(os.Stdout), // 输出到终端
			level,
		)
		cores = append(cores, consoleCore)
	}

	// 创建最终的 Logger Core
	core := zapcore.NewTee(cores...)

	// 构建 Logger 并添加调用者信息
	// zap.Addcaller() 输出日志打印文件和行数如： logger/logger_test.go:33
	// zap.AddCallerSkip(1) 会记录调用日志的文件名和行号，但是有时候我们可能封装了日志方法，
	// 并希望跳过这些封装函数的调用栈，直接定位到业务代码的调用位置。这时就需要使用
	logger := zap.New(core, zap.AddCaller())

	// 可选：替换全局的 zap logger
	zap.ReplaceGlobals(logger)

	// 1. zap.ReplaceGlobals 函数将当前初始化的 logger 替换到全局的 logger,
	// 2. 使用 logger 的时候 直接通过 zap.S().Debugf("xxx") or zap.L().Debug("xxx")
	// 3. 使用 zap.S() 和 zap.L() 提供全局锁，保证一个全局的安全访问logger的方式
	// zap.ReplaceGlobals(logger) // 替换zap包中全局的logger实例，后续在其他包中只需使用zap.L()调用即可
	// zap.L().Debug("")
	// zap.S().Debugf("")
	return logger, func() {
		_ = logger.Sync() // 确保所有日志写入
		if lj != nil {
			_ = lj.Close() // 关闭 lumberjack logger
		}
	}
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
