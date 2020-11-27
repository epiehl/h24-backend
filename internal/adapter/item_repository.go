package adapter

import (
	"github.com/epiehl93/h24-notifier/pkg/models"
	"github.com/epiehl93/h24-notifier/pkg/repository"
	"gorm.io/gorm"
	"time"
)

type ItemRepository struct {
	db *gorm.DB
}

func (i ItemRepository) FindAvailableInOutlet() ([]*models.Item, error) {
	var items []*models.Item

	tx := i.db.Where("available_in_outlet = ?", true).Find(&items)
	if tx.Error != nil {
		return nil, tx.Error
	}

	if tx.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	return items, nil
}

func (i ItemRepository) SetAvailableInOutlet(item *models.Item) error {
	item.AvailableInOutlet = true

	err := i.Update(item)
	if err != nil {
		return err
	}

	return nil
}

func (i ItemRepository) SetUnavailableInOutlet(item *models.Item) error {
	item.AvailableInOutlet = false
	item.OutletPrice = 0
	item.AvailableInOutletSince = time.Time{}

	err := i.Update(item)
	if err != nil {
		return err
	}

	return nil
}

func (i ItemRepository) GetBySKU(sku uint64) (*models.Item, error) {
	item := &models.Item{}
	tx := i.db.Where("sku = ?", sku).First(item)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return item, nil
}

func (i ItemRepository) Create(item *models.Item) error {
	tx := i.db.Create(item)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

func (i ItemRepository) GetAll() ([]*models.Item, error) {
	var items []*models.Item

	tx := i.db.Find(&items)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return items, nil
}

func (i ItemRepository) GetAllPaginated(limit int, page int) ([]*models.Item, error) {
	var items []*models.Item

	if limit > 100 {
		limit = 100
	}

	offset := (page - 1) * limit
	tx := i.db.Scopes(func(db *gorm.DB) *gorm.DB {
		return db.Offset(offset).Limit(limit)
	}).Find(&items)

	if tx.Error != nil {
		return nil, tx.Error
	}

	return items, nil
}

func (i ItemRepository) Get(item *models.Item) error {
	tx := i.db.First(item)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

func (i ItemRepository) Update(item *models.Item) error {
	if item.AvailableInOutlet && item.AvailableInOutletSince.IsZero() {
		item.AvailableInOutletSince = time.Now()
	}
	item.UpdatedAt = time.Now()
	tx := i.db.Save(item)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

func (i ItemRepository) Delete(item *models.Item) error {
	tx := i.db.Delete(item)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

func NewItemRepository(db *gorm.DB) repository.ItemRepository {
	return &ItemRepository{db}
}
