package commands

import (
	"github.com/epiehl93/h24-notifier/config"
	"github.com/epiehl93/h24-notifier/internal/utils"
	"github.com/epiehl93/h24-notifier/internal/web"
	"github.com/shurcooL/graphql"
	"github.com/spf13/cobra"
)

var RunServer = &cobra.Command{
	Use:   "server",
	Short: "start the server",
	RunE: func(cmd *cobra.Command, args []string) error {
		config.ReadConfig()

		db, err := utils.SetupDB()
		if err != nil {
			return err
		}

		gql := graphql.NewClient(config.C.H24Connector.Endpoint, nil)

		if err := utils.MigrateTables(db); err != nil {
			return err
		}

		app, err := web.NewApp(db, gql)
		if err != nil {
			return err
		}
		if err := app.Run(); err != nil {
			return err
		}
		return nil
	},
}
