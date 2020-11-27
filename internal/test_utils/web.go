package test_utils

import (
	"github.com/epiehl93/h24-notifier/config"
	"github.com/epiehl93/h24-notifier/internal/utils"
	"github.com/epiehl93/h24-notifier/internal/web"
	"github.com/shurcooL/graphql"
	"github.com/spf13/viper"
)

func StartTestServer() (web.App, error) {
	config.ReadConfig()

	db, err := SetupTestDb()
	if err != nil {
		return nil, err
	}

	gql := graphql.NewClient(viper.GetString("h24connector.endpoint"), nil)
	app, err := web.NewApp(db, gql)
	if err != nil {
		utils.Log.Panic(err)
	}

	if err := app.Run(); err != nil {
		return nil, err
	}

	return app, nil
}
