package sub

import (
	"fmt"
	"github.com/reggiepy/LogBeetle/boot"
	"os"
	"os/signal"
	"syscall"

	"github.com/reggiepy/LogBeetle/global"
	"github.com/reggiepy/LogBeetle/goutils/enumUtils"
	"github.com/reggiepy/LogBeetle/version"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	showVersion  bool
	configFormat = enumUtils.NewEnum([]string{"humanReadable", "simple"}, "humanReadable")
)

func init() {
	viper.SetDefault("config", "./log-beetle.yaml")
	// 设置全局标志
	rootCmd.PersistentFlags().BoolVarP(&showVersion, "version", "v", false, "show version information")
	rootCmd.PersistentFlags().StringP("config", "c", "", "config file")

	// 添加命令行参数
	rootCmd.Flags().String("log-file", "", "file to log")
	rootCmd.Flags().String("consumer-log-path", "", "Path to consumer log")
	rootCmd.Flags().String("nsq-address", "", "Address of NSQ")
	rootCmd.Flags().String("nsq-auth-secret", "", "NSQ auth secret")

	// 绑定命令行参数到Viper
	_ = viper.BindPFlag("log_file", rootCmd.Flags().Lookup("log-file"))
	_ = viper.BindPFlag("nsq_address", rootCmd.Flags().Lookup("nsq-address"))
	_ = viper.BindPFlag("nsq_auth_secret", rootCmd.Flags().Lookup("nsq-auth-secret"))
	_ = viper.BindPFlag("consumer_log_path", rootCmd.Flags().Lookup("consumer-log-path"))
	_ = viper.BindPFlag("config", rootCmd.PersistentFlags().Lookup("config"))
}

var rootCmd = cobra.Command{
	Use:   "LogBeetle",
	Short: "LogBeetle help",
	Long:  `LogBeetle help`,
	Args: func(cmd *cobra.Command, args []string) error {
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		if showVersion {
			fmt.Println(version.Full())
			return nil
		}
		StartServer()
		return nil
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func StartServer() {
	// 初始化全局组件
	global.LbViper = boot.Viper()
	global.LbLogger = boot.Log()
	global.LbNsqProducer = boot.NsqProducer()
	boot.Consumer()
	boot.Boot()

	// 捕获信号，以优雅地退出程序
	waitForShutdown()

	// 关闭资源
	cleanup()
	fmt.Println("Main program stopped")
}

func waitForShutdown() {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh
}

func cleanup() {
	global.LbNsqProducer.Stop()
	global.LbLogger.Info("NSQ producer stopped")
	global.LBConsumerManager.Stop()
	global.LbLogger.Info("NSQ consumer stopped")
	_ = global.LbLogger.Sync() // 确保在程序退出时刷新日志缓冲区
}
