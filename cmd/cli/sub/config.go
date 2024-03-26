package sub

import (
	"fmt"
	"github.com/reggiepy/LogBeetle/pkg/config"
	"github.com/spf13/cobra"
	"os"
)

var (
	configFormat = NewEnum([]string{"humanReadable", "simple"}, "humanReadable")
)

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.AddCommand(configShowCmd)
	configShowCmd.Flags().Var(configFormat, "format", "humanReadable | simple")
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
	Short: "show config tools",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := config.ShowConfig(configFormat.Value); err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		return nil
	},
}
