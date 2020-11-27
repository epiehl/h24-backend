package web

// Wishlist is a json representation of a wishlist
// swagger:model
type Wishlist struct {
	// ID is the unique id of the wishlist
	ID uint64 `json:"id"`
	// Name for the wishlist
	Name string `json:"name"`
	// Items added to this wishlist
	ItemSkus []uint64 `json:"item_skus"`
}

// Item is a json representation of an item
// swagger:model
type Item struct {
	ID                  uint64  `json:"id"`
	Name                string  `json:"name"`
	SKU                 uint64  `json:"sku"`
	RetailPrice         float64 `json:"retail_price"`
	RetailDiscount      float64 `json:"retail_discount"`
	RetailDiscountPrice float64 `json:"retail_discount_price"`
	OutletPrice         float64 `json:"outlet_price"`
	AvailableInOutlet   bool    `json:"available_in_outlet"`
}
