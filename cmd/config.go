package cmd

import (
	"fmt"
	"github.com/reggiepy/LogBeetle/pkg/config"
	"github.com/spf13/cobra"
	"os"
)

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.AddCommand(configShowCmd)
	configShowCmd.Flags().Var(configFormat, "format", "humanReadable | simple")
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
	Short: "show config tools",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := config.ShowConfig(configFormat.Value); err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		return nil
	},
}
