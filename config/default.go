package config

import (
	"fmt"
	"os"
	"reflect"

	"gopkg.in/yaml.v2"
)

func DefaultConfig() *Config {
	return &Config{
		Env:     "dev",
		BaseUrl: "http://127.0.0.1:1233",
		Port:    1233,
		LogPath: "logs",
		LogConfig: LogConfig{
			File:       "./logs/log-beetle.log", // 默认日志文件位置
			MaxSize:    5,                       // 最大日志文件大小限制为 5MB
			MaxBackups: 3,                       // 保留最多 3 个旧日志
			MaxAge:     7,                       // 旧日志文件保留 7 天
			Compress:   true,                    // 启用日志压缩
			Level:      "info",                  // 默认日志级别为 info
			Format:     "json",                  // 默认日志格式为 JSON
		},
		ConsumerConfig: ConsumerConfig{
			LogPath:      "./logs",
			NSQConsumers: []NSQConsumers{}, // 初始化为空切片
		},
		NSQConfig: NSQConfig{
			AuthSecret:  "",
			NSQDAddress: "127.0.0.1:4150",
		},
		Store: Store{
			Root:            "/store",
			ChanLength:      64,
			AutoAddDate:     true,
			SaveDays:        180,
			MaxIdleTime:     300,
			GoMaxProcessIdx: -1,
		},
		Search: Search{
			PageSize:        100,
			NearSearchSize:  200,
			MultiLineSearch: false,
		},
	}
}

func SetDefaults(config *Config) {
	defaultConfig := DefaultConfig()
	SetStructDefaults(config, defaultConfig)
}

// SetStructDefaults 从默认配置中复制非零值字段到目标配置中
// 如果目标配置中的字段是切片类型，则递归处理每个元素
// 如果目标配置中的字段是结构体类型，则递归调用该方法处理嵌套结构体
// 如果目标配置中的字段是零值，则使用默认配置中对应字段的值进行填充
func SetStructDefaults(config, defaultConfig interface{}) {
	// 获取目标配置和默认配置的反射值
	configValue := reflect.ValueOf(config).Elem()
	defaultConfigValue := reflect.ValueOf(defaultConfig).Elem()

	// 遍历目标配置的字段
	for i := 0; i < configValue.NumField(); i++ {
		field := configValue.Field(i)

		// 如果字段是切片类型，则递归处理每个元素
		if configValue.Type().Field(i).Name == "NSQConsumers" && field.Kind() == reflect.Slice {
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
			SetStructDefaults(field.Addr().Interface(), defaultConfigValue.Field(i).Addr().Interface())
		} else {
			// 如果字段是零值，则使用默认配置中对应字段的值进行填充
			if field.Interface() == reflect.Zero(field.Type()).Interface() {
				defaultFieldValue := defaultConfigValue.Field(i)
				field.Set(defaultFieldValue)
			}
		}
	}
}

func IsEmptyValue(value reflect.Value) bool {
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

// SaveConfigToFile 将配置保存到 YAML 文件
func SaveConfigToFile(config *Config, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("无法创建文件: %v", err)
	}
	defer file.Close()

	encoder := yaml.NewEncoder(file)
	defer encoder.Close()

	if err := encoder.Encode(config); err != nil {
		return fmt.Errorf("无法编码配置到 YAML: %v", err)
	}

	return nil
}
