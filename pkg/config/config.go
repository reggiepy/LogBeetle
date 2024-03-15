package config

import (
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

var Instance *Config

type Config struct {
	Env        string `yaml:"Env"`        // 环境：prod、dev
	BaseUrl    string `yaml:"BaseUrl"`    // base url
	Port       string `yaml:"Port"`       // 端口
	LogFile    string `yaml:"LogFile"`    // 日志文件
}

func Init(filename string) *Config {
	Instance = &Config{}
	if yamlFile, err := ioutil.ReadFile(filename); err != nil {
		log.Error().Msg(err)
	} else if err = yaml.Unmarshal(yamlFile, Instance); err != nil {
		log.Error().Msg(err)
	}
	return Instance
}
