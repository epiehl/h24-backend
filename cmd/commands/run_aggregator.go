package commands

import (
	"github.com/epiehl93/h24-notifier/config"
	"github.com/epiehl93/h24-notifier/internal/adapter"
	"github.com/epiehl93/h24-notifier/internal/aggregator"
	"github.com/epiehl93/h24-notifier/internal/utils"
	"github.com/shurcooL/graphql"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var RunAggregator = &cobra.Command{
	Use:   "aggregator",
	Short: "start the aggregator",
	RunE: func(cmd *cobra.Command, args []string) error {
		config.ReadConfig()

		db, err := utils.SetupDB()
		if err != nil {
			return err
		}

		gql := graphql.NewClient(viper.GetString("h24connector.endpoint"), nil)

		if err := utils.MigrateTables(db); err != nil {
			return err
		}

		reg := adapter.NewRegistry(db, gql)

		a := aggregator.NewOutletAggregator(reg)
		if err := a.Run(); err != nil {
			return err
		}
		return nil
	},
}
