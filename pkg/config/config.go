package config

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"os"
	"reflect"
)

var Instance *Config

type Config struct {
	Env     string `yaml:"Env"`     // 环境：prod、dev
	BaseUrl string `yaml:"BaseUrl"` // base url
	Port    int    `yaml:"Port"`    // 端口
	LogPath string `yaml:"LogPath"` // 日志文件名

	LogConfig struct {
		LogFile    string `yaml:"LogFile"`    // 日志文件名
		MaxSize    int    `yaml:"MaxSize"`    // 日志文件大小限制，单位为 MB
		MaxBackups int    `yaml:"MaxBackups"` // 最大保留的旧日志文件数量
		MaxAge     int    `yaml:"MaxAge"`     // 旧日志文件保留天数
		Compress   bool   `yaml:"Compress"`   // 是否压缩旧日志文件
		LogLevel   string `yaml:"LogLevel"`   // 日志等级
		LogFormat  string `yaml:"LogFormat"`  // 日志等级
	} `yaml:"LogConfig"`

	ConsumerConfig struct {
		LogPath string `yaml:"LogPath"` //日志输出路径
	} `yaml:"ConsumerConfig"`

	NSQConfig struct {
		AuthSecret  string `yaml:"AuthSecret"`  // 权限
		NSQDAddress string `yaml:"NSQDAddress"` // 地址
	} `yaml:"NSQConfig"`
}

func DefaultConfig() *Config {
	config := new(Config)
	config.Env = "dev"
	config.BaseUrl = "http://127.0.0.1:1233"
	config.Port = 1233
	config.LogPath = "logs"

	config.LogConfig.LogFile = "app.log"
	config.LogConfig.MaxSize = 1    // 默认值为 100 MB
	config.LogConfig.MaxBackups = 5 // 默认值为 10
	config.LogConfig.MaxAge = 30    // 默认值为 30 天
	config.LogConfig.Compress = true
	config.LogConfig.LogLevel = "info"
	config.LogConfig.LogFormat = "json"

	config.ConsumerConfig.LogPath = "./logs"

	config.NSQConfig.AuthSecret = "%n&yFA2JD85z^g"
	config.NSQConfig.NSQDAddress = "127.0.0.1:4150"

	return config
}

//func setDefaults(config *Config) {
//	defaultConfig := DefaultConfig()
//
//	if config.Env == "" {
//		config.Env = defaultConfig.Env
//	}
//
//	if config.BaseUrl == "" {
//		config.BaseUrl = defaultConfig.BaseUrl
//	}
//
//	if config.Port == 0 {
//		config.Port = defaultConfig.Port
//	}
//
//	if config.LogPath == "" {
//		config.LogPath = defaultConfig.LogPath
//	}
//
//	if config.LogConfig.LogFile == "" {
//		config.LogConfig.LogFile = defaultConfig.LogConfig.LogFile
//	}
//
//	if config.LogConfig.MaxSize == 0 {
//		config.LogConfig.MaxSize = defaultConfig.LogConfig.MaxSize
//	}
//
//	if config.LogConfig.MaxBackups == 0 {
//		config.LogConfig.MaxBackups = defaultConfig.LogConfig.MaxBackups
//	}
//
//	if config.LogConfig.MaxAge == 0 {
//		config.LogConfig.MaxAge = defaultConfig.LogConfig.MaxAge
//	}
//
//	if config.LogConfig.Compress == false {
//		config.LogConfig.Compress = defaultConfig.LogConfig.Compress
//	}
//
//	if config.LogConfig.LogLevel == "" {
//		config.LogConfig.LogLevel = defaultConfig.LogConfig.LogLevel
//	}
//
//	if config.LogConfig.LogFormat == "" {
//		config.LogConfig.LogFormat = defaultConfig.LogConfig.LogFormat
//	}
//
//	if config.ConsumerConfig.LogPath == "" {
//		config.ConsumerConfig.LogPath = defaultConfig.ConsumerConfig.LogPath
//	}
//}

func setDefaults(config *Config) {
	defaultConfig := DefaultConfig()
	setStructDefaults(config, defaultConfig)
}

func setStructDefaults(config, defaultConfig interface{}) {
	configValue := reflect.ValueOf(config).Elem()
	defaultConfigValue := reflect.ValueOf(defaultConfig).Elem()
	for i := 0; i < configValue.NumField(); i++ {
		field := configValue.Field(i)
		if field.Kind() == reflect.Struct {
			setStructDefaults(field.Addr().Interface(), defaultConfigValue.Field(i).Addr().Interface())
		} else {
			//isEmptyValue是一个私有方法，用于检查字段的值是否为空。
			//在示例中，我并没有使用这个方法，而是使用了reflect.Zero函数来检查字段的零值。
			//reflect.Zero函数会返回给定类型的零值，然后与字段的当前值进行比较，以确定字段是否为零值。
			if field.Interface() == reflect.Zero(field.Type()).Interface() {
				defaultFieldValue := defaultConfigValue.Field(i)
				field.Set(defaultFieldValue)
			}
		}
	}
}

func isEmptyValue(value reflect.Value) bool {
	switch value.Kind() {
	case reflect.String:
		return value.String() == ""
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return value.Int() == 0
	case reflect.Bool:
		return !value.Bool()
	default:
		return false
	}
}

func Init(filename string) *Config {
	Instance = &Config{}
	if yamlFile, err := os.ReadFile(filename); err != nil {
		fmt.Println(err)
	} else if err = yaml.Unmarshal(yamlFile, Instance); err != nil {
		fmt.Println(err)
	}
	setDefaults(Instance)
	return Instance
}
