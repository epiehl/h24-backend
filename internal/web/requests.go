package web

import "net/http"

// swagger:model CreateWishlistRequest
type CreateWishlistRequest struct {
	Name string `json:"name"`
}

func (req CreateWishlistRequest) Bind(r *http.Request) error {
	return nil
}
