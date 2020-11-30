package main

import (
	"fmt"
	"github.com/epiehl93/h24-notifier/cmd/commands"
	"github.com/spf13/cobra"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "h24",
	Short: "h24 binary to start server,aggregator or notificator",
	RunE: func(cmd *cobra.Command, args []string) error {
		err := cmd.Help()
		if err != nil {
			return err
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(
		commands.RunServer,
		commands.RunAggregator,
		commands.RunNotificator,
	)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
