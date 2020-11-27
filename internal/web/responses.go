package web

import (
	"github.com/epiehl93/h24-notifier/pkg/models"
	"github.com/go-chi/render"
	"net/http"
	"strconv"
)

// WishlistResponse contains response information that will be returned to the request issuer
//
// swagger:response WishlistResponse
type WishlistResponse struct {
	Wishlist
}

func (resp WishlistResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func NewWishlistResponse(wishlist *models.Wishlist) *WishlistResponse {
	skus := make([]uint64, wishlist.Length())
	var itemUrl string

	for index, item := range wishlist.Items {
		if index == 0 {
			itemUrl = item.ImageUrl
		}
		skus[index] = item.SKU
	}

	respList := Wishlist{
		wishlist.ID,
		wishlist.Name,
		itemUrl,
		skus,
	}

	resp := &WishlistResponse{
		respList,
	}
	return resp
}

func NewWishlistListResponse(w []*models.Wishlist) []render.Renderer {
	var l []render.Renderer

	// required to not have null in our response body ever
	if len(w) == 0 {
		l = make([]render.Renderer, 0)
	}
	for _, list := range w {
		l = append(l, NewWishlistResponse(list))
	}
	return l
}

//  DeleteResponse contains information about the deletion process
//
//  swagger:response DeleteResponse
type DeleteResponse struct {
	// in: body
	Success string `json:"success"`
}

func (d DeleteResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func NewDeleteResponse(success bool) *DeleteResponse {
	return &DeleteResponse{Success: strconv.FormatBool(success)}
}

// ItemResponse contains response information that will be returned to the request issuer
//
// swagger:response ItemResponse
type ItemResponse struct {
	Item
}

func (resp ItemResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func NewItemResponse(item *models.Item) *ItemResponse {
	newItem := Item{
		item.ID,
		item.Name,
		item.SKU,
		item.ImageUrl,
		item.RetailUrl,
		item.RetailPrice,
		item.RetailDiscount,
		item.RetailDiscountPrice,
		item.AvailableInRetail,
		item.OutletPrice,
		item.AvailableInOutlet,
	}
	return &ItemResponse{newItem}
}

func NewItemListResponse(list []*models.Item) []render.Renderer {
	var respList []render.Renderer
	if len(list) == 0 {
		respList = make([]render.Renderer, 0)
	} else {
		for _, item := range list {
			respList = append(respList, NewItemResponse(item))
		}
	}

	return respList
}

// PaginatedItemListResponse contains response information that will be returned to the request issuer
type PaginatedItemListResponse struct {
	Items    []*Item `json:"items"`
	MaxPages int64   `json:"max_pages"`
	NumItems int64   `json:"num_items"`
}

func (p PaginatedItemListResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func NewPaginatedItemListResponse(items []*models.Item, maxPages int64, numItems int64) render.Renderer {
	var dataItems []*Item
	for _, item := range items {
		newItem := &Item{}
		newItem = newItem.FromModelItem(item)
		dataItems = append(dataItems, newItem)
	}
	return PaginatedItemListResponse{
		dataItems,
		maxPages,
		numItems,
	}
}
