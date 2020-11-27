package web

import (
	"errors"
	"github.com/epiehl93/h24-notifier/internal/adapter"
	"github.com/epiehl93/h24-notifier/pkg/models"
	"github.com/go-chi/render"
	"log"
	"net/http"
)

type WishlistController interface {
	GetWishlist(w http.ResponseWriter, r *http.Request)
	GetAllWishlists(w http.ResponseWriter, r *http.Request)
	CreateWishlist(w http.ResponseWriter, r *http.Request)
	DeleteWishlist(w http.ResponseWriter, r *http.Request)
	AddItemToWishlist(w http.ResponseWriter, r *http.Request)
	DeleteItemFromWishlist(w http.ResponseWriter, r *http.Request)
}

type wishlistController struct {
	adapter.Registry
}

//  swagger:operation GET /wishlist GetAllWishlists
//  ---
//  summary: Gets all wishlists associated with the user
//  responses:
//    default:
//      "$ref": "#/responses/ErrResponse"
//    200:
//      description: wishlist response
//      schema:
//        type: array
//        items:
//          "$ref": "#/definitions/Wishlist"
func (c wishlistController) GetAllWishlists(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userUUID := ctx.Value("userUUID").(string)

	lists, err := c.WishlistRepository.GetAll(userUUID)
	if err != nil {
		_ = render.Render(w, r, ErrRender(err))
	}

	if err := render.RenderList(w, r, NewWishlistListResponse(lists)); err != nil {
		_ = render.Render(w, r, ErrRender(err))
		return
	}

}

//  swagger:operation POST /wishlist/{id}/item/{sku} AddItemToWishlist
//  ---
//  summary: Adds an item to our wishlist
//  responses:
//    default:
//      "$ref": "#/responses/ErrResponse"
//    200:
//      "$ref": "#/responses/WishlistResponse"
func (c wishlistController) AddItemToWishlist(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	list, ok := ctx.Value("wishlist").(*models.Wishlist)
	if !ok {
		_ = render.Render(w, r, ErrRender(errors.New("error fetching wishlist")))
		return
	}

	item, ok := ctx.Value("item").(*models.Item)
	if !ok {
		_ = render.Render(w, r, ErrRender(errors.New("error fetching item")))
		return
	}

	if err := list.AddItem(item); err != nil {
		_ = render.Render(w, r, ErrRender(err))
		return
	}

	if err := c.WishlistRepository.Update(list); err != nil {
		_ = render.Render(w, r, ErrRender(err))
		return
	}

	render.Status(r, http.StatusCreated)
	err := render.Render(w, r, NewWishlistResponse(list))
	if err != nil {
		log.Println(err)
	}
}

//  swagger:operation DELETE /wishlist/{id}/item/{sku} DeleteItemFromWishlist
//  ---
//  summary: Delete an item to our wishlist
//  responses:
//    default:
//      "$ref": "#/responses/ErrResponse"
//    200:
//      "$ref": "#/responses/WishlistResponse"
func (c wishlistController) DeleteItemFromWishlist(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	list, ok := ctx.Value("wishlist").(*models.Wishlist)
	if !ok {
		_ = render.Render(w, r, ErrRender(errors.New("error fetching wishlist")))
		return
	}

	item, ok := ctx.Value("item").(*models.Item)
	if !ok {
		_ = render.Render(w, r, ErrRender(errors.New("error fetching item")))
		return
	}

	if err := c.WishlistRepository.RemoveItem(list, item); err != nil {
		_ = render.Render(w, r, ErrRender(errors.New("error deleting item from wishlist")))
		return
	}

	err := render.Render(w, r, NewWishlistResponse(list))
	if err != nil {
		log.Println(err)
	}
}

//  swagger:operation POST /wishlist CreateWishlist
//  ---
//  summary: Create a new wishlist
//
//  parameters:
//  - name: wishlist
//    in: body
//    description: wishlist parameter
//    schema:
//      "$ref": "#/definitions/CreateWishlistRequest"
//    required: true
//  responses:
//    default:
//      "$ref": "#/responses/ErrResponse"
//    200:
//      "$ref": "#/responses/WishlistResponse"
func (c wishlistController) CreateWishlist(w http.ResponseWriter, r *http.Request) {
	data := &CreateWishlistRequest{}
	userUUID := r.Context().Value("userUUID").(string)
	if err := render.Bind(r, data); err != nil {
		_ = render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	list := &models.Wishlist{Name: data.Name, UserSub: userUUID}
	if err := c.WishlistRepository.Create(list); err != nil {
		log.Println(err)
		_ = render.Render(w, r, ErrInternalServerError(err))
		return
	}

	render.Status(r, http.StatusCreated)
	_ = render.Render(w, r, NewWishlistResponse(list))
}

//  swagger:operation GET /wishlist/{id} GetWishlist
//  ---
//  summary: Get a wishlist by ID
//  responses:
//    default:
//      "$ref": "#/responses/ErrResponse"
//    200:
//      "$ref": "#/responses/WishlistResponse"
func (c wishlistController) GetWishlist(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	wlist, ok := ctx.Value("wishlist").(*models.Wishlist)
	if !ok {
		_ = render.Render(w, r, ErrInvalidRequest(errors.New("error fetching wishlist")))
		return
	}

	if err := render.Render(w, r, NewWishlistResponse(wlist)); err != nil {
		_ = render.Render(w, r, ErrRender(err))
		return
	}
}

//  swagger:operation DELETE /wishlist/{id} DeleteWishlist
//  ---
//  summary: Delete a wishlist by ID
//  responses:
//    default:
//      "$ref": "#/responses/ErrResponse"
//    200:
//      "$ref": "#/responses/DeleteResponse"
func (c wishlistController) DeleteWishlist(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	wlist, ok := ctx.Value("wishlist").(*models.Wishlist)
	if !ok {
		_ = render.Render(w, r, ErrInvalidRequest(errors.New("error fetching wishlist")))
		return
	}

	err := c.WishlistRepository.Delete(wlist)
	if err != nil {
		_ = render.Render(w, r, ErrRender(errors.New("unable to delete wishlist")))
	}

	if err := render.Render(w, r, NewDeleteResponse(true)); err != nil {
		_ = render.Render(w, r, ErrRender(err))
		return
	}
}

func NewWishlistController(r adapter.Registry) WishlistController {
	return &wishlistController{r}
}
