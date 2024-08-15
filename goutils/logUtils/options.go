package logUtils

type Option interface {
	apply(config *Config)
}

type optionFunc func(config *Config)

func (o optionFunc) apply(config *Config) {
	o(config)
}

// WithLogFile 设置日志文件名
func WithLogFile(logFile string) Option {
	return optionFunc(func(config *Config) {
		config.LogFile = logFile
		return
	})
}

// WithMaxSize 设置日志文件大小限制，单位为 MB
func WithMaxSize(maxSize int) Option {
	return optionFunc(func(loggerConfig *Config) {
		loggerConfig.MaxSize = maxSize
	})
}

// WithMaxBackups 设置最大保留的旧日志文件数量
func WithMaxBackups(maxBackups int) Option {
	return optionFunc(func(loggerConfig *Config) {
		loggerConfig.MaxBackups = maxBackups
	})
}

// WithMaxAge 设置旧日志文件保留天数
func WithMaxAge(maxAge int) Option {
	return optionFunc(func(loggerConfig *Config) {
		loggerConfig.MaxAge = maxAge
	})
}

// WithCompress 设置是否压缩旧日志文件
func WithCompress(compress bool) Option {
	return optionFunc(func(loggerConfig *Config) {
		loggerConfig.Compress = compress
	})
}

// WithLogLevel 设置日志等级
func WithLogLevel(logLevel string) Option {
	return optionFunc(func(loggerConfig *Config) {
		loggerConfig.LogLevel = logLevel
	})
}

// WithLogFormat 设置日志格式
func WithLogFormat(logFormat string) Option {
	return optionFunc(func(loggerConfig *Config) {
		loggerConfig.LogFormat = logFormat
	})
}
