package models

import (
	"errors"
	"fmt"
	"time"
)

type AuthProvider string

type Wishlist struct {
	ID        uint64    `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Items     []*Item   `json:"items" gorm:"many2many:wishlist_items;constraint:OnDelete:CASCADE;"`
	UserSub   string    `json:"user_sub"`
}

type Item struct {
	ID                     uint64    `json:"id"`
	Name                   string    `json:"name"`
	SKU                    uint64    `json:"sku" gorm:"unique;not null"`
	ImageUrl               string    `json:"image_url"`
	RetailUrl              string    `json:"retail_url"`
	RetailPrice            float64   `json:"retail_price"`
	RetailDiscount         float64   `json:"retail_discount"`
	RetailDiscountPrice    float64   `json:"retail_discount_price"`
	OutletPrice            float64   `json:"outlet_price"`
	CreatedAt              time.Time `json:"created_at"`
	UpdatedAt              time.Time `json:"updated_at"`
	AvailableInOutlet      bool      `json:"available_in_outlet"`
	AvailableInOutletSince time.Time `json:"available_in_outlet_since"`
	LastAggregatedAt       time.Time `json:"last_aggregated_at"`
}

func (w *Wishlist) AddItem(i *Item) error {
	var itemLimit int

	// check if item exists already
	found := false
	for _, item := range w.Items {
		if item.ID == i.ID {
			found = true
			break
		}
	}

	if found {
		return nil
	}

	if len(w.Items) > 10 {
		return errors.New(fmt.Sprintf("can't have more than %d items in a wishlist", itemLimit))
	}

	w.Items = append(w.Items, i)
	return nil
}

// RemoveItem removes an item from the wishlist
func (w *Wishlist) RemoveItem(i *Item) error {
	var itemIndex int

	found := false
	for index, item := range w.Items {
		if item.ID == i.ID {
			itemIndex = index
			found = true
			break
		}
	}

	if !found {
		return errors.New("could not find item in wishlist")
	}

	// remove item by copying the last item onto it
	w.Items[itemIndex] = w.Items[w.Length()-1]
	w.Items = w.Items[:w.Length()-1]

	return nil
}

func (w *Wishlist) Length() int {
	return len(w.Items)
}

type CycleType string

const (
	AggregationCycle  = "aggregation"
	NotificationCycle = "notification"
	UnsetCycle        = ""
)

type Cycle struct {
	ID         uint64    `json:"id"`
	At         time.Time `json:"at"`
	Successful bool      `json:"successful"`
	Type       CycleType `json:"type"`
}
