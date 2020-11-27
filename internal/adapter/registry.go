package adapter

import (
	"github.com/epiehl93/h24-notifier/pkg/repository"
)

type Registry struct {
	Wishlist repository.WishlistRepository
	Item     repository.ItemRepository
	H24      repository.H24Connector
	Cycle    repository.CycleRepository
}

func NewRegistry(wishlist repository.WishlistRepository, item repository.ItemRepository, h24 repository.H24Connector, cycle repository.CycleRepository) Registry {
	return Registry{
		Wishlist: wishlist,
		Item:     item,
		H24:      h24,
		Cycle:    cycle,
	}
}
