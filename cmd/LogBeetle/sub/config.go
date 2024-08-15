package sub

import (
	"fmt"
	"os"

	"github.com/reggiepy/LogBeetle/boot"
	"github.com/reggiepy/LogBeetle/config"
	"github.com/reggiepy/LogBeetle/global"
	"github.com/reggiepy/LogBeetle/goutils/fileUtils"
	"github.com/reggiepy/LogBeetle/goutils/jsonUtils"
	"github.com/spf13/viper"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.AddCommand(configShowCmd)
	configShowCmd.Flags().Var(configFormat, "format", "humanReadable | simple")
	configCmd.ValidArgsFunction = func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"json", "simple"}, cobra.ShellCompDirectiveNoFileComp
	}

	configCmd.AddCommand(configGenerateCmd)
	configCmd.ValidArgsFunction = func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"json", "simple"}, cobra.ShellCompDirectiveNoFileComp
	}
}

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "config tools",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println(cmd.Help())
		return nil
	},
}

var configShowCmd = &cobra.Command{
	Use:   "show",
	Short: "show config",
	RunE: func(cmd *cobra.Command, args []string) error {
		global.LbViper = boot.Viper()
		data, err := jsonUtils.AnyToJson(global.LbConfig, configFormat.String())
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		fmt.Println(data)
		return nil
	},
}

var configGenerateCmd = &cobra.Command{
	Use:   "generate",
	Short: "generate default config",
	RunE: func(cmd *cobra.Command, args []string) error {
		configFile := viper.GetString("config")
		if fileUtils.FileExists(configFile) {
			return fmt.Errorf("config file already exists. please remove it before running this command")
		}
		defaultConfig := config.DefaultConfig()

		err := config.SaveConfigToFile(defaultConfig, configFile)
		if err != nil {
			return err
		}
		return nil
	},
}
