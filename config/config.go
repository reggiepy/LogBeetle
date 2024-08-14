package config

type Config struct {
	Env            string         `yaml:"Env"`     // 环境：prod、dev
	BaseUrl        string         `yaml:"BaseUrl"` // base url
	Port           int            `yaml:"Port"`    // 端口
	LogPath        string         `yaml:"LogPath"` // 日志文件名
	LogConfig      LogConfig      `yaml:"LogConfig"`
	ConsumerConfig ConsumerConfig `yaml:"ConsumerConfig"`
	NSQConfig      NSQConfig      `yaml:"NSQConfig"`
}
type LogConfig struct {
	LogFile    string `yaml:"LogFile"`    // 日志文件名
	MaxSize    int    `yaml:"MaxSize"`    // 日志文件大小限制，单位为 MB
	MaxBackups int    `yaml:"MaxBackups"` // 最大保留的旧日志文件数量
	MaxAge     int    `yaml:"MaxAge"`     // 旧日志文件保留天数
	Compress   bool   `yaml:"Compress"`   // 是否压缩旧日志文件
	LogLevel   string `yaml:"LogLevel"`   // 日志等级
	LogFormat  string `yaml:"LogFormat"`  // 日志等级
}
type NSQConfig struct {
	AuthSecret  string `yaml:"AuthSecret"`  // 权限
	NSQDAddress string `yaml:"NSQDAddress"` // 地址
}

type ConsumerConfig struct {
	LogPath      string         `yaml:"LogPath"` // 日志输出路径
	NSQConsumers []NSQConsumers `yaml:"NSQConsumers"`
}
type NSQConsumers struct {
	Name     string `yaml:"Name"`
	Topic    string `yaml:"Topic"`
	FileName string `yaml:"FileName"`
}
