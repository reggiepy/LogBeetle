package sub

import (
	"fmt"
	"github.com/reggiepy/LogBeetle/goutils/enumUtils"
	"os"

	"github.com/reggiepy/LogBeetle/boot"
	"github.com/reggiepy/LogBeetle/global"
	"github.com/reggiepy/LogBeetle/goutils/jsonUtils"
	"github.com/spf13/viper"

	"github.com/spf13/cobra"
)

type ConfigConfig struct {
	Format *enumUtils.Enum
	Force  bool
	Config string
}

var (
	configConfig = ConfigConfig{
		Format: enumUtils.NewEnum([]string{"humanReadable", "simple"}, "humanReadable"),
	}
)

func init() {
	configShowCmd.Flags().Var(configConfig.Format, "format", "humanReadable | simple")
	configShowCmd.Flags().StringVarP(&configConfig.Config, "config", "c", "", "config file")
	_ = viper.BindPFlag("config", configGenerateCmd.PersistentFlags().Lookup("config"))
	configCmd.AddCommand(configShowCmd)

	configGenerateCmd.Flags().BoolVarP(&configConfig.Force, "force", "f", false, "Generate configuration forces")
	configGenerateCmd.Flags().StringVarP(&configConfig.Config, "config", "c", "", "config file")
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
		return []string{"json", "simple"}, cobra.ShellCompDirectiveNoFileComp
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		global.LbViper = boot.Viper()
		data, err := jsonUtils.AnyToJson(global.LbConfig, configConfig.Format.String())
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
		global.LbViper = boot.Viper()
		if !configConfig.Force {
			if err := global.LbViper.SafeWriteConfig(); err != nil {
				return err
			}
		} else {
			if err := viper.WriteConfig(); err != nil {
				return err
			}
		}
		return nil
	},
}
