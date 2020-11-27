package web

import (
	"errors"
	"github.com/epiehl93/h24-notifier/internal/adapter"
	"github.com/epiehl93/h24-notifier/internal/utils"
	"github.com/epiehl93/h24-notifier/pkg/models"
	"github.com/go-chi/render"
	"gorm.io/gorm"
	"net/http"
	"regexp"
	"strconv"
)

type ItemController interface {
	GetItem(w http.ResponseWriter, r *http.Request)
	PaginatedListItems(w http.ResponseWriter, r *http.Request)
	SearchItems(w http.ResponseWriter, r *http.Request)
}

type itemController struct {
	*adapter.Registry
}

//  swagger:operation GET /item/search SearchItems
//  ---
//  summary: Search an item by Searchterm
//  responses:
//    default:
//      "$ref": "#/responses/ErrResponse"
//    200:
//      description: item response
//      schema:
//        type: array
//        items:
//          "$ref": "#/definitions/Item"
func (i itemController) SearchItems(w http.ResponseWriter, r *http.Request) {
	term := r.URL.Query().Get("term")
	var items []*models.Item
	if term == "" {
		_ = render.Render(w, r, ErrRender(errors.New("please specify term in query")))
		return
	}

	reg, err := regexp.Compile("^\\d{10,18}$")
	if err != nil {
		_ = render.Render(w, r, ErrRender(err))
		return
	}
	// This looks like an sku, fetch it and terminate
	ok := reg.MatchString(term)
	if ok {
		sku, err := strconv.ParseUint(term, 10, 64)
		if err != nil {
			_ = render.Render(w, r, ErrRender(err))
			return
		}
		item, err := i.Item.GetBySKU(sku)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				_ = render.Render(w, r, ErrNotFound())
				return
			} else {
				_ = render.Render(w, r, ErrRender(err))
				return
			}
		}
		items = append(items, item)
	} else {
		items, err = i.Item.SearchInName(term)
		if err != nil {
			_ = render.Render(w, r, ErrRender(err))
			return
		}
	}

	if err := render.RenderList(w, r, NewItemListResponse(items)); err != nil {
		_ = render.Render(w, r, ErrRender(err))
	}
}

//  swagger:operation GET /item/{sku} GetItem
//  ---
//  summary: Get an item by SKU
//  responses:
//    default:
//      "$ref": "#/responses/ErrResponse"
//    200:
//      "$ref": "#/responses/ItemResponse"
func (i itemController) GetItem(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	item, ok := ctx.Value("item").(*models.Item)
	if !ok {
		_ = render.Render(w, r, ErrRender(errors.New("could not fetch item")))
		return
	}

	if err := render.Render(w, r, NewItemResponse(item)); err != nil {
		_ = render.Render(w, r, ErrRender(err))
		return
	}
}

//  swagger:operation GET /item PaginatedListItems
//  ---
//  summary: List items
//  security:
//    - Bearer: []
//  responses:
//    default:
//      "$ref": "#/responses/ErrResponse"
//    200:
//      description: paginated list response
//      schema:
//        type: object
//		  properties:
//          items:
//            type: string
//            description: This should be a list of items but openapi 2 sucks
//          maxPages:
//            type: integer
//            description: the number of pages
func (i itemController) PaginatedListItems(w http.ResponseWriter, r *http.Request) {
	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil {
		_ = render.Render(w, r, ErrRender(err))
	}

	if limit == 0 {
		limit = 100
	}

	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		_ = render.Render(w, r, ErrRender(err))
	}

	if page <= 0 {
		page = 1
	}

	var availableInOutlet *bool
	queryOutletAvailable := r.URL.Query().Get("available_in_outlet")
	if queryOutletAvailable != "" {
		if b, err := strconv.ParseBool(queryOutletAvailable); err != nil {
			utils.Log.Error(err)
			availableInOutlet = nil
		} else {
			availableInOutlet = &b
		}
	}

	items, maxPages, numItems, err := i.Item.GetAllPaginated(limit, page, availableInOutlet)
	if err != nil {
		_ = render.Render(w, r, ErrRender(err))
		return
	}

	// Enrich items
	for _, item := range items {
		changed, err := i.H24.EnrichItemWithRetailData(item)
		if err != nil {
			_ = render.Render(w, r, ErrRender(err))
			return
		}
		if changed {
			err := i.Item.Update(item)
			if err != nil {
				_ = render.Render(w, r, ErrRender(err))
				return
			}
		}
	}

	if err := render.Render(w, r, NewPaginatedItemListResponse(items, maxPages, numItems)); err != nil {
		_ = render.Render(w, r, ErrRender(err))
		return
	}
}

func NewItemController(registry *adapter.Registry) ItemController {
	return &itemController{registry}
}
