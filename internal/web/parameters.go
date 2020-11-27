package web

/*
These are structs only used for generating swagger documentation.
They are not used anywhere in the code else than here.
*/

// swagger:parameters GetWishlist DeleteWishlist AddItemToWishlist DeleteItemFromWishlist
type IDParam struct {
	// in: path
	ID uint64 `json:"id"`
}

// swagger:parameters AddItemToWishlist DeleteItemFromWishlist GetItem
type SKUParam struct {
	// in: path
	SKU uint64 `json:"sku"`
}

// swagger:parameters PaginatedListItems
type LimitParam struct {
	// in: query
	// required: false
	// schema:
	//   type: integer
	//   default: 10
	Limit int `json:"limit"`
}

// swagger:parameters PaginatedListItems
type PageParam struct {
	// in: query
	// required: false
	// schema:
	//   type: integer
	//   default: 1
	Page int `json:"page"`
}

// swagger:parameters PaginatedListItems
type AvailableInOutlet struct {
	// in: query
	// required: false
	// schema:
	//   type: boolean
	AvailableInOutlet bool `json:"available_in_outlet"`
}

// swagger:parameters SearchItems
type SearchTerm struct {
	// in: query
	// required: true
	// schema:
	//   type: string
	Term string `json:"term"`
}
