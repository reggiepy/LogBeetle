package sub

import (
	"github.com/reggiepy/LogBeetle/pkg/config"
	"github.com/spf13/cobra"
	"os"
)

var (
	//配置文件
	configFile string
)

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.Flags().StringVarP(&configFile, "config", "c", "./log-beetle.yaml", "config file")
	configCmd.ValidArgsFunction = func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"json", "simple"}, cobra.ShellCompDirectiveNoFileComp
	}
}

var rootCmd = cobra.Command{
	Use:     "cli",
	Short:   "cli help",
	Long:    `cli help`,
	Version: "",
	RunE: func(cmd *cobra.Command, args []string) error {
		_ = cmd.Help()
		return nil
	},
}

func initConfig() {
	// 初始化配置
	config.Init(configFile)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
