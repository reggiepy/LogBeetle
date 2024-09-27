package config

type Config struct {
	Env            string         `yaml:"Env"`     // 环境：prod、dev
	BaseUrl        string         `yaml:"BaseUrl"` // base url
	Port           int            `yaml:"Port"`    // 端口
	LogPath        string         `yaml:"LogPath"` // 日志文件名
	LogConfig      LogConfig      `yaml:"LogConfig"`
	ConsumerConfig ConsumerConfig `yaml:"ConsumerConfig"`
	NSQConfig      NSQConfig      `yaml:"NSQConfig"`
	Store          Store          `yaml:"Store"`
	Search         Search         `yaml:"Search"`
}
type LogConfig struct {
	File       string `json:"File" yaml:"File"`             // 日志文件名
	MaxSize    int    `json:"MaxSize" yaml:"MaxSize"`       // 日志文件大小限制（单位：MB）
	MaxBackups int    `json:"MaxBackups" yaml:"MaxBackups"` // 最大保留的旧日志文件数量
	MaxAge     int    `json:"MaxAge" yaml:"MaxAge"`         // 旧日志文件保留天数
	Compress   bool   `json:"Compress" yaml:"Compress"`     // 是否压缩旧日志文件
	Level      string `json:"LogLevel" yaml:"LogLevel"`     // 日志级别
	Format     string `json:"LogFormat" yaml:"LogFormat"`   // 日志格式（如：json、logfmt）
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

type Store struct {
	Root                 string `yaml:"Root"`
	ChanLength           int    `yaml:"ChanLength"`
	AutoAddDate          bool   `yaml:"AutoAddDate"`
	SaveDays             int    `yaml:"SaveDays"`
	MaxIdleTime          int    `yaml:"MaxIdleTime"`
	GoMaxProcessIdx      int    `yaml:"GoMaxProcessIdx"`
}

type Search struct {
	PageSize  int  `yaml:"PageSize"`
	NearSearchSize  int  `yaml:"NearSearchSize"`
	MultiLineSearch bool `yaml:"MultiLineSearch"`
}
