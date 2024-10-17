package config

import (
	"encoding/json"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Env            string         `json:"env" yaml:"Env"`          // 环境：prod、dev
	BaseUrl        string         `json:"base_url" yaml:"BaseUrl"` // base url
	Port           int            `json:"port" yaml:"Port"`        // 端口
	LogPath        string         `json:"log_path" yaml:"LogPath"` // 日志文件名
	LogConfig      LogConfig      `json:"log_config" yaml:"LogConfig"`
	ConsumerConfig ConsumerConfig `json:"consumer_config" yaml:"ConsumerConfig"`
	NSQConfig      NSQConfig      `json:"nsq_config" yaml:"NSQConfig"`
	Store          Store          `json:"store" yaml:"Store"`
	Search         Search         `json:"search" yaml:"Search"`
}

type LogConfig struct {
	File       string `json:"file" yaml:"File"`              // 日志文件名
	MaxSize    int    `json:"max_size" yaml:"MaxSize"`       // 日志文件大小限制（单位：MB）
	MaxBackups int    `json:"max_backups" yaml:"MaxBackups"` // 最大保留的旧日志文件数量
	MaxAge     int    `json:"max_age" yaml:"MaxAge"`         // 旧日志文件保留天数
	Compress   bool   `json:"compress" yaml:"Compress"`      // 是否压缩旧日志文件
	Level      string `json:"log_level" yaml:"LogLevel"`     // 日志级别
	Format     string `json:"log_format" yaml:"LogFormat"`   // 日志格式（如：json、logfmt）
}

type NSQConfig struct {
	AuthSecret  string `json:"auth_secret" yaml:"AuthSecret"`   // 权限
	NSQDAddress string `json:"nsqd_address" yaml:"NSQDAddress"` // 地址
}

type ConsumerConfig struct {
	LogPath      string         `json:"log_path" yaml:"LogPath"` // 日志输出路径
	NSQConsumers []NSQConsumers `json:"nsq_consumers" yaml:"NSQConsumers"`
}

type NSQConsumers struct {
	Name     string `json:"name" yaml:"Name"`
	Topic    string `json:"topic" yaml:"Topic"`
	FileName string `json:"file_name" yaml:"FileName"`
}

type Store struct {
	Root            string `json:"root" yaml:"Root"`
	ChanLength      int    `json:"chan_length" yaml:"ChanLength"`
	AutoAddDate     bool   `json:"auto_add_date" yaml:"AutoAddDate"`
	SaveDays        int    `json:"save_days" yaml:"SaveDays"`
	MaxIdleTime     int    `json:"max_idle_time" yaml:"MaxIdleTime"`
	GoMaxProcessIdx int    `json:"go_max_process_idx" yaml:"GoMaxProcessIdx"`
}

type Search struct {
	PageSize        int  `json:"page_size" yaml:"PageSize"`
	NearSearchSize  int  `json:"near_search_size" yaml:"NearSearchSize"`
	MultiLineSearch bool `json:"multi_line_search" yaml:"MultiLineSearch"`
}

// ToJson 将配置转换为JSON格式
func (c *Config) ToJson() string {
	jsonData, _ := json.MarshalIndent(c, "", "  ")
	return string(jsonData)
}

// LoadJson 从JSON文件加载配置
func (c *Config) LoadJson(data string) error {
	if err := json.Unmarshal([]byte(data), c); err != nil {
		return err
	}
	return nil
}

// ToYaml 将配置转换为YAML格式
func (c *Config) ToYaml() (string, error) {
	yamlData, err := yaml.Marshal(c)
	if err != nil {
		return "", err
	}
	return string(yamlData), nil
}

// LoadYaml 从YAML文件加载配置
func (c *Config) LoadYaml(data string) error {
	if err := yaml.Unmarshal([]byte(data), c); err != nil {
		return err
	}
	return nil
}
