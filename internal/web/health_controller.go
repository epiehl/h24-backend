package web

import (
	"github.com/epiehl93/h24-notifier/internal/adapter"
	"github.com/go-chi/render"
	"net/http"
)

type HealthController interface {
	GetHealth(w http.ResponseWriter, r *http.Request)
}

type healthController struct {
	adapter.Registry
}

func (h healthController) GetHealth(w http.ResponseWriter, r *http.Request) {
	render.Status(r, http.StatusOK)
}

func NewHealthController(r adapter.Registry) HealthController {
	return &healthController{r}
}
