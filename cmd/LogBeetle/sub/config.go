package sub

import (
	"fmt"
	"github.com/reggiepy/LogBeetle/pkg/config"
	"github.com/spf13/cobra"
	"os"
)

var (
	configFormat *Enum
)

func init() {
	rootCmd.AddCommand(configCmd)

	configFormat = NewEnum([]string{"humanReadable", "simple"}, "humanReadable")
	configCmd.Flags().Var(configFormat, "format", "humanReadable | simple")
}

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "config tools",
	RunE: func(cmd *cobra.Command, args []string) error {

		if err := config.ShowConfig(configFormat.Value); err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		return nil
	},
}
