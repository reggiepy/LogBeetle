package sub

import (
	"fmt"
	"github.com/gookit/goutil/fsutil"
	"github.com/gookit/goutil/jsonutil"
	"github.com/reggiepy/LogBeetle/boot"
	"github.com/reggiepy/LogBeetle/config"
	"github.com/reggiepy/LogBeetle/global"
	"github.com/reggiepy/LogBeetle/pkg/goutils/enumUtils"
	"github.com/reggiepy/LogBeetle/pkg/goutils/yamlutil"
	"github.com/spf13/viper"
	"os"

	"github.com/spf13/cobra"
)

var (
	format  *enumUtils.Enum = enumUtils.NewEnum([]string{"humanReadable", "simple"}, "humanReadable")
	force   bool
	cfgFile string
)

func init() {
	configShowCmd.Flags().Var(format, "format", "humanReadable | simple")
	configShowCmd.Flags().StringVarP(&cfgFile, "config", "c", "", "config file")
	_ = viper.BindPFlag("config", configGenerateCmd.PersistentFlags().Lookup("config"))
	configCmd.AddCommand(configShowCmd)

	configGenerateCmd.Flags().BoolVarP(&force, "force", "f", false, "Generate configuration forces")
	configGenerateCmd.Flags().StringVarP(&cfgFile, "config", "c", "", "config file")
	_ = viper.BindPFlag("config", configGenerateCmd.PersistentFlags().Lookup("config"))
	configCmd.AddCommand(configGenerateCmd)
	rootCmd.AddCommand(configCmd)
}

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "config tools",
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

var configShowCmd = &cobra.Command{
	Use:   "show",
	Short: "show config",
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"humanReadable", "simple"}, cobra.ShellCompDirectiveNoFileComp
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		global.LbViper = boot.Viper()
		configFormat := format.String()

		var data string
		switch configFormat {
		case "humanReadable":
			data, _ = jsonutil.EncodeString(global.LbConfig)
		case "simple":
			dataBytes, _ := jsonutil.Encode(global.LbConfig)
			data = string(dataBytes)
		}
		fmt.Println(data)
		return nil
	},
}

var configGenerateCmd = &cobra.Command{
	Use:   "generate",
	Short: "generate default config",
	RunE: func(cmd *cobra.Command, args []string) error {
		var err error
		configFile := viper.GetString("config")
		if configFile == "" {
			return fmt.Errorf("config file not specified")
		}

		configFileExt := fsutil.Extname(configFile)

		defaultConfig := config.DefaultConfig()
		configString := ""
		switch configFileExt {
		case "yaml":
			configString, _ = yamlutil.EncodeString(defaultConfig)
		case "json":
			configString, _ = jsonutil.EncodeString(defaultConfig)
		default:
			return fmt.Errorf("unsupported config file extension: %s", configFileExt)
		}
		flags := os.O_CREATE | os.O_TRUNC | os.O_WRONLY
		if !force {
			flags |= os.O_EXCL
		}
		err = fsutil.WriteFile(configFile, configString, os.ModePerm, flags)
		if err != nil {
			return fmt.Errorf("write config to file failed: %v", err)
		}
		return nil
	},
}
