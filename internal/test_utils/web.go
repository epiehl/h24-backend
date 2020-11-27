package test_utils

import (
	"github.com/epiehl93/h24-notifier/config"
	"github.com/epiehl93/h24-notifier/internal/web"
	"github.com/shurcooL/graphql"
)

func StartTestServer() (web.App, error) {
	config.ReadConfig()

	db, err := SetupTestDb()
	if err != nil {
		return nil, err
	}

	gql := graphql.NewClient(config.C.H24Connector.Endpoint, nil)
	app := web.NewApp(db, gql)

	if err := app.Run(); err != nil {
		return nil, err
	}

	return app, nil
}
