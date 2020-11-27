package repository

import (
	"github.com/epiehl93/h24-notifier/pkg/models"
)

type ItemRepository interface {
	Create(item *models.Item) error
	GetAll() ([]*models.Item, error)
	GetAllPaginated(limit int, page int, availableInOutlet *bool) ([]*models.Item, int64, int64, error)
	Get(item *models.Item) error
	GetBySKU(sku uint64) (*models.Item, error)
	Update(item *models.Item) error
	Delete(item *models.Item) error
	SetUnavailableInOutlet(item *models.Item) error
	SetAvailableInOutlet(item *models.Item) error
	FindAvailableInOutlet() ([]*models.Item, error)
	SearchInName(name string) ([]*models.Item, error)
}
