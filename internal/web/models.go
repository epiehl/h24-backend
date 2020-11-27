package web

import "github.com/epiehl93/h24-notifier/pkg/models"

// Wishlist is a json representation of a wishlist
// swagger:model
type Wishlist struct {
	// ID is the unique id of the wishlist
	ID uint64 `json:"id"`
	// Name for the wishlist
	Name         string `json:"name"`
	ExampleImage string `json:"example_image"`
	// Items added to this wishlist
	ItemSkus []uint64 `json:"item_skus"`
}

// Item is a json representation of an item
// swagger:model
type Item struct {
	ID                  uint64  `json:"id"`
	Name                string  `json:"name"`
	SKU                 uint64  `json:"sku"`
	ImageUrl            string  `json:"image_url"`
	RetailUrl           string  `json:"retail_url"`
	RetailPrice         float64 `json:"retail_price"`
	RetailDiscount      float64 `json:"retail_discount"`
	RetailDiscountPrice float64 `json:"retail_discount_price"`
	AvailableInRetail   bool    `json:"available_in_retail"`
	OutletPrice         float64 `json:"outlet_price"`
	AvailableInOutlet   bool    `json:"available_in_outlet"`
}

func (i Item) FromModelItem(item *models.Item) *Item {
	i.ID = item.ID
	i.Name = item.Name
	i.SKU = item.SKU
	i.ImageUrl = item.ImageUrl
	i.RetailUrl = item.RetailUrl
	i.RetailPrice = item.RetailPrice
	i.RetailDiscount = item.RetailDiscount
	i.RetailDiscountPrice = item.RetailDiscountPrice
	i.AvailableInRetail = item.AvailableInRetail
	i.OutletPrice = item.OutletPrice
	i.AvailableInOutlet = item.AvailableInOutlet
	return &i
}
