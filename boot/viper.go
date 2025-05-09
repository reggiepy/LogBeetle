package boot

import (
	"errors"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/reggiepy/LogBeetle/global"
	"github.com/spf13/viper"
)

func Viper() *viper.Viper {
	configFile := viper.GetString("config")
	v := viper.New()
	if configFile == "" {
		v.SetConfigName("config") // 不带扩展名
		v.SetConfigType("yaml")   // 指定文件类型
		v.AddConfigPath(".")      // 当前目录查找
	} else {
		v.SetConfigFile(configFile) // 传入完整路径，如 "./config.yaml"
	}
	v.SetEnvPrefix("LB") // 设置环境变量前缀
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		var configFileNotFoundError viper.ConfigFileNotFoundError
		if errors.As(err, &configFileNotFoundError) {
			panic(fmt.Errorf("Config file not found：%s \n", err.Error()))
		}
	}
	v.WatchConfig()
	v.OnConfigChange(func(e fsnotify.Event) {
		//fmt.Println("Config file changed: ", e.String())
		BindConfig(v)
		//configString, _ := jsonutil.EncodeString(global.LbConfig)
		//fmt.Println("Config file changed: ", configString)
	})
	BindConfig(v)
	SetupCombaConfig()
	return v
}

func BindConfig(v *viper.Viper) {
	if err := v.Unmarshal(&global.LbConfig); err != nil {
		panic(fmt.Errorf("Failed to bind config file：%s \n", err))
	}
}

func SetupCombaConfig() {
	// 从Viper获取日志文件路径
	logPath := viper.GetString("log_file")
	// 如果日志文件路径不为空，则更新配置中的日志文件路径
	if logPath != "" {
		global.LbConfig.LogConfig.File = logPath
	}

	// 从Viper获取消费者日志路径
	consumerLogPath := viper.GetString("consumer_log_path")
	// 如果消费者日志路径不为空，则更新配置中的消费者日志路径
	if consumerLogPath != "" {
		global.LbConfig.ConsumerConfig.LogPath = consumerLogPath
	}

	// 从Viper获取NSQ地址
	nsqAddress := viper.GetString("nsq_address")
	// 如果NSQ地址不为空，则更新配置中的NSQ地址
	if nsqAddress != "" {
		global.LbConfig.NSQConfig.NSQDAddress = nsqAddress
	}

	// 从Viper获取NSQ认证密钥
	nsqAuthSecret := viper.GetString("nsq_auth_secret")
	// 如果NSQ认证密钥不为空，则更新配置中的NSQ认证密钥
	if nsqAuthSecret != "" {
		global.LbConfig.NSQConfig.AuthSecret = nsqAuthSecret
	}
}
