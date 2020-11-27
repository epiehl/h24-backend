//+build wireinject

package internal

import (
	"github.com/epiehl93/h24-notifier/internal/adapter"
	"github.com/epiehl93/h24-notifier/internal/aggregator"
	"github.com/epiehl93/h24-notifier/internal/notificator"
	"github.com/epiehl93/h24-notifier/internal/web"
	"github.com/google/wire"
	"github.com/shurcooL/graphql"
	"gorm.io/gorm"
)

func InitializeRegistry(db *gorm.DB, gql *graphql.Client) adapter.Registry {
	wire.Build(adapter.NewCycleRepository, adapter.NewH24Connector, adapter.NewWishlistRepository, adapter.NewItemRepository, adapter.NewRegistry)
	return adapter.Registry{}
}

func InitializeApp(r *adapter.Registry) (web.App, error) {
	wire.Build(web.NewApp, web.NewApplicationContext, web.NewHealthController, web.NewItemController, web.NewWishlistController)
	return &web.AppImpl{}, nil
}

func InitializeOutletAggregator(r *adapter.Registry) aggregator.OutletAggregator {
	wire.Build(aggregator.NewOutletAggregator)
	return &aggregator.OutletAggregatorImpl{}
}

func InitializeNotificator(r *adapter.Registry) notificator.Notificator {
	wire.Build(notificator.NewNotificator)
	return &notificator.NotificatorImpl{}
}
