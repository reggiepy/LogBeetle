package zapLogger

import "encoding/json"

type LoggerConfig struct {
	InFile         bool `json:"InFile" yaml:"InFile"`                 // 是否输出到终端
	InConsole      bool `json:"InConsole" yaml:"LogInFile"`           // 是否输出到文件
	ReplaceGlobals bool `json:"ReplaceGlobals" yaml:"ReplaceGlobals"` // 是否替换全局日志记录器

	File       string `json:"File" yaml:"File"`             // 日志文件名
	MaxSize    int    `json:"MaxSize" yaml:"MaxSize"`       // 日志文件大小限制（单位：MB）
	MaxBackups int    `json:"MaxBackups" yaml:"MaxBackups"` // 最大保留的旧日志文件数量
	MaxAge     int    `json:"MaxAge" yaml:"MaxAge"`         // 旧日志文件保留天数
	Compress   bool   `json:"Compress" yaml:"Compress"`     // 是否压缩旧日志文件
	Level      string `json:"LogLevel" yaml:"LogLevel"`     // 日志级别
	Format     string `json:"LogFormat" yaml:"LogFormat"`   // 日志格式（如：json、logfmt）
}

func (l *LoggerConfig) clone() *LoggerConfig {
	clone := *l
	return &clone
}
func (l *LoggerConfig) ToJSON() string {
	jsonStr, err := json.Marshal(l)
	if err != nil {
		return ""
	}
	return string(jsonStr)
}

func (l *LoggerConfig) LoadJSON(jsonStr string) error {
	return json.Unmarshal([]byte(jsonStr), l)
}

func (l *LoggerConfig) WithOptions(options ...Option) *LoggerConfig {
	c := l.clone()
	for _, opt := range options {
		opt.apply(l)
	}
	return c
}

// NewConfig 创建默认配置
func NewLoggerConfig(opts ...Option) *LoggerConfig {
	config := &LoggerConfig{
		InFile:         true, // 默认同时输出到终端和文件
		InConsole:      false,
		File:           "app.log",
		MaxSize:        1,
		MaxBackups:     5,
		MaxAge:         30,
		Compress:       false,
		Level:          "info",
		Format:         "json", // 默认使用 JSON 格式
		ReplaceGlobals: false,  // 默认不替换全局日志记录器
	}
	return config.WithOptions(opts...)
}
