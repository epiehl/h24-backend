package main

import (
	"github.com/epiehl93/h24-notifier/cmd/commands"
	"github.com/epiehl93/h24-notifier/internal/utils"
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
	err := utils.InitLogger()
	if err != nil {
		utils.Log.Panic(err)
	}
	if err := rootCmd.Execute(); err != nil {
		utils.Log.Error(err)
		os.Exit(1)
	}
}
