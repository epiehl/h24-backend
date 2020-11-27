package main

import (
	"github.com/epiehl93/h24-notifier/config"
	"github.com/epiehl93/h24-notifier/internal/notificator"
	"github.com/epiehl93/h24-notifier/internal/utils"
	"github.com/shurcooL/graphql"
	"log"
)

func main() {
	config.ReadConfig()

	db, err := utils.SetupDB()
	if err != nil {
		log.Panicln(err)
	}

	gql := graphql.NewClient(config.C.H24Connector.Endpoint, nil)

	if err := utils.MigrateTables(db); err != nil {
		log.Panicln(err)
	}

	app := notificator.NewNotificator(db, gql)
	if err := app.Run(); err != nil {
		log.Panicln(err)
	}
}
