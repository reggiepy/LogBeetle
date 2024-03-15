package main

import (
	"fmt"
	"github.com/reggiepy/LogBeetle/pkg/config"
	"github.com/spf13/cobra"
)

var (
	configFile string
)

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.Flags().StringVarP(&configFile, "config", "c", "./log-beetle.yaml", "config file")
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
	conf := config.Init(configFile)
	fmt.Println(conf)
}
