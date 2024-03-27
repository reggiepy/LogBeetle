package config

import (
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v2"
	"os"
	"reflect"
)

var Instance *Config

type Consumer struct {
	Name     string `yaml:"Name"`
	Topic    string `yaml:"Topic"`
	Channel  string `yaml:"Channel"`
	FileName string `yaml:"FileName"`
}

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
		LogPath   string     `yaml:"LogPath"` //日志输出路径
		Consumers []Consumer `yaml:"Consumers"`
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
	config.ConsumerConfig.Consumers = make([]Consumer, 0)

	config.NSQConfig.AuthSecret = ""
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

// setStructDefaults 从默认配置中复制非零值字段到目标配置中
// 如果目标配置中的字段是切片类型，则递归处理每个元素
// 如果目标配置中的字段是结构体类型，则递归调用该方法处理嵌套结构体
// 如果目标配置中的字段是零值，则使用默认配置中对应字段的值进行填充
func setStructDefaults(config, defaultConfig interface{}) {
	// 获取目标配置和默认配置的反射值
	configValue := reflect.ValueOf(config).Elem()
	defaultConfigValue := reflect.ValueOf(defaultConfig).Elem()

	// 遍历目标配置的字段
	for i := 0; i < configValue.NumField(); i++ {
		field := configValue.Field(i)

		// 如果字段是切片类型，则递归处理每个元素
		if configValue.Type().Field(i).Name == "Consumers" && field.Kind() == reflect.Slice {
			// 删除目标配置中的空切片以及空元素
			var updatedSlice reflect.Value
			updatedSlice = reflect.MakeSlice(field.Type(), 0, 0)
			for j := 0; j < field.Len(); j++ {
				if field.Index(j).Interface() == reflect.Zero(field.Index(j).Type()).Interface() {
					// 如果元素是零值，则不追加到新切片中
				} else {
					// 追加非零值元素到新切片中
					updatedSlice = reflect.Append(updatedSlice, field.Index(j))
				}
			}
			// 将更新后的切片替换原始切片
			field.Set(updatedSlice)
		} else if field.Kind() == reflect.Struct {
			// 如果字段是结构体类型，则递归调用该方法处理嵌套结构体
			setStructDefaults(field.Addr().Interface(), defaultConfigValue.Field(i).Addr().Interface())
		} else {
			// 如果字段是零值，则使用默认配置中对应字段的值进行填充
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
		_, _ = fmt.Fprintf(os.Stderr, err.Error())

	} else if err = yaml.Unmarshal(yamlFile, Instance); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, err.Error())
	}
	setDefaults(Instance)
	return Instance
}

func ShowConfig(format string) error {
	var jsonData []byte
	var err error
	switch format {
	case "simple":
		// 将结构体序列化为 JSON
		jsonData, err = json.Marshal(Instance)
	case "humanReadable":
		// 将结构体序列化为 JSON
		jsonData, err = json.MarshalIndent(Instance, "", "    ")
	default:
		return fmt.Errorf("不支持的format类型: %s", format)
	}
	if err != nil {
		return fmt.Errorf("配置信息json序列化失败: %v", err)
	}
	fmt.Println(string(jsonData))
	return nil
}
