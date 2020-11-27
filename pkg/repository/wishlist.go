package repository

import (
	"github.com/epiehl93/h24-notifier/pkg/models"
)

type WishlistRepository interface {
	Create(list *models.Wishlist) error
	Get(list *models.Wishlist) error
	GetAll(userSub string) ([]*models.Wishlist, error)
	Update(list *models.Wishlist) error
	Delete(list *models.Wishlist) error
	RemoveItem(list *models.Wishlist, item *models.Item) error
	FindWishlistsWithAvailableInOutlet() ([]*models.Wishlist, error)
}
