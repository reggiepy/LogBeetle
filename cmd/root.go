package cmd

import (
	"fmt"
	"github.com/reggiepy/LogBeetle/pkg/config"
	"github.com/reggiepy/LogBeetle/pkg/util/version"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

var (
	//配置文件
	configFile  string
	showVersion bool

	configFormat = NewEnum([]string{"humanReadable", "simple"}, "humanReadable")
)

func init() {
	// 使用Viper加载环境变量
	viper.SetEnvPrefix("lb") // 设置环境变量前缀
	viper.SetDefault("config", "./log-beetle.yaml")
	viper.AutomaticEnv() // 自动加载环境变量
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().BoolVarP(&showVersion, "version", "v", false, "show version information")
	rootCmd.Flags().StringP("config", "c", "", "config file")
	_ = viper.BindPFlag("config", rootCmd.Flags().Lookup("config"))
}

var rootCmd = cobra.Command{
	Use:     "LogBeetle",
	Short:   "LogBeetle help",
	Long:    `LogBeetle help`,
	Version: "",
	RunE: func(cmd *cobra.Command, args []string) error {
		if showVersion {
			fmt.Println(version.Full())
			return nil
		}
		_ = cmd.Help()
		return nil
	},
}

func initConfig() {
	configFile = viper.GetString("config")
	// 初始化配置
	config.Init(configFile)

	logPath := viper.GetString("log_path")
	if logPath != "" {
		config.Instance.LogConfig.LogFile = logPath
	}
	consumerLogPath := viper.GetString("consumer_log_path")
	if consumerLogPath != "" {
		config.Instance.ConsumerConfig.LogPath = consumerLogPath
	}
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
