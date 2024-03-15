package config

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"os"
)

var Instance *Config

type Config struct {
	Env     string `yaml:"Env"`     // 环境：prod、dev
	BaseUrl string `yaml:"BaseUrl"` // base url
	Port    string `yaml:"Port"`    // 端口
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
}

func Init(filename string) *Config {
	Instance = &Config{}
	if yamlFile, err := os.ReadFile(filename); err != nil {
		fmt.Println(err)
	} else if err = yaml.Unmarshal(yamlFile, Instance); err != nil {
		fmt.Println(err)
	}
	return Instance
}
