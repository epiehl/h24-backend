package adapter

import (
	"github.com/epiehl93/h24-notifier/pkg/repository"
	"github.com/shurcooL/graphql"
	"gorm.io/gorm"
)

type Registry struct {
	repository.WishlistRepository
	repository.ItemRepository
	repository.H24Connector
	repository.CycleRepository
}

func NewRegistry(db *gorm.DB, gql *graphql.Client) Registry {
	return Registry{
		NewWishlistRepository(db),
		NewItemRepository(db),
		NewH24Connector(gql),
		NewCycleRepository(db),
	}
}
