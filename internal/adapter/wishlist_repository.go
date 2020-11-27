package adapter

import (
	"github.com/epiehl93/h24-notifier/pkg/models"
	"github.com/epiehl93/h24-notifier/pkg/repository"
	"gorm.io/gorm"
	"time"
)

type WishlistRepository struct {
	db *gorm.DB
}

func (w WishlistRepository) FindWishlistsWithAvailableInOutlet() ([]*models.Wishlist, error) {
	var wishlists []*models.Wishlist
	var returnWishlists []*models.Wishlist

	tx := w.db.Preload("Items").Find(&wishlists)
	if tx.Error != nil {
		return nil, tx.Error
	}

	for _, list := range wishlists {
		itemAvailable := false
		for _, item := range list.Items {
			if item.AvailableInOutlet {
				itemAvailable = true
				break
			}
		}
		if itemAvailable {
			returnWishlists = append(returnWishlists, list)
		}
	}

	return returnWishlists, nil
}

func (w WishlistRepository) RemoveItem(list *models.Wishlist, item *models.Item) error {
	err := list.RemoveItem(item)
	if err != nil {
		return err
	}

	err = w.db.Model(list).Association("Items").Delete(item)
	if err != nil {
		return err
	}

	return nil
}

func (w WishlistRepository) Create(list *models.Wishlist) error {
	tx := w.db.Create(&list)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

func (w WishlistRepository) Get(list *models.Wishlist) error {
	tx := w.db.Preload("Items").First(list)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

func (w WishlistRepository) GetAll(userSub string) ([]*models.Wishlist, error) {
	var lists []*models.Wishlist
	tx := w.db.Preload("Items").Where("user_sub = ?", userSub).Find(&lists)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return lists, nil
}

func (w WishlistRepository) Update(list *models.Wishlist) error {
	list.UpdatedAt = time.Now()
	tx := w.db.Save(list).Association("Items")
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

func (w WishlistRepository) Delete(list *models.Wishlist) error {
	tx := w.db.Delete(list)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

func NewWishlistRepository(db *gorm.DB) repository.WishlistRepository {
	return &WishlistRepository{db}
}
