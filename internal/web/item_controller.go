package web

import (
	"errors"
	"github.com/epiehl93/h24-notifier/internal/adapter"
	"github.com/epiehl93/h24-notifier/pkg/models"
	"github.com/go-chi/render"
	"net/http"
	"strconv"
)

type ItemController interface {
	GetItem(w http.ResponseWriter, r *http.Request)
	ListItems(w http.ResponseWriter, r *http.Request)
}

type itemController struct {
	adapter.Registry
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

//  swagger:operation GET /item ListItems
//  ---
//  summary: List items
//  security:
//    - Bearer: []
//  responses:
//    default:
//      "$ref": "#/responses/ErrResponse"
//    200:
//      description: item response
//      schema:
//        type: array
//        items:
//          "$ref": "#/definitions/Item"
func (i itemController) ListItems(w http.ResponseWriter, r *http.Request) {
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

	if page == 0 {
		page = 1
	}

	items, err := i.ItemRepository.GetAllPaginated(limit, page)
	if err != nil {
		_ = render.Render(w, r, ErrRender(err))
		return
	}

	if err := render.RenderList(w, r, NewItemListResponse(items)); err != nil {
		_ = render.Render(w, r, ErrRender(err))
		return
	}
}

func NewItemController(registry adapter.Registry) ItemController {
	return &itemController{registry}
}
