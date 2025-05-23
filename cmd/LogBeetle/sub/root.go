package sub

import (
	"fmt"
	"github.com/reggiepy/LogBeetle/boot"
	"github.com/reggiepy/LogBeetle/global"
	"github.com/reggiepy/LogBeetle/pkg/version"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	showVersion bool
)

func init() {
	//cobra.OnInitialize(initConfig)
	// 设置全局标志
	rootCmd.PersistentFlags().BoolVarP(&showVersion, "version", "v", false, "show version information")
	rootCmd.PersistentFlags().StringP("config", "c", "", "config file (default: config.yaml)")

	// 添加命令行参数
	rootCmd.Flags().String("log-file", "", "file to log")
	rootCmd.Flags().String("consumer-log-path", "", "Path to consumer log")
	rootCmd.Flags().String("nsq-address", "", "Address of NSQ")
	rootCmd.Flags().String("nsq-auth-secret", "", "NSQ auth secret")

	// 绑定命令行参数到Viper
	_ = viper.BindPFlag("config", rootCmd.PersistentFlags().Lookup("config"))
	_ = viper.BindPFlag("log_file", rootCmd.Flags().Lookup("log-file"))
	_ = viper.BindPFlag("nsq_address", rootCmd.Flags().Lookup("nsq-address"))
	_ = viper.BindPFlag("nsq_auth_secret", rootCmd.Flags().Lookup("nsq-auth-secret"))
	_ = viper.BindPFlag("consumer_log_path", rootCmd.Flags().Lookup("consumer-log-path"))
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
		global.LbViper = boot.Viper()
		//configString, _ := jsonutil.EncodeString(global.LbConfig)
		//fmt.Println("Config: ", configString)
		global.LbLogger = boot.Logger()
		global.LbNsqProducer = boot.NsqProducer(global.LbConfig.NSQConfig)
		boot.Ldb()
		global.LBConsumerManager = boot.Consumer()
		boot.Boot()
		return nil
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

//func initConfig() {
//
//}
