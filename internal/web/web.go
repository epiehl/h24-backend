package web

import (
	"fmt"
	"github.com/766b/chi-logger"
	"github.com/epiehl93/h24-notifier/internal/adapter"
	"github.com/epiehl93/h24-notifier/internal/utils"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/spf13/viper"
	httpSwagger "github.com/swaggo/http-swagger"
	"log"
	"net/http"
	"os"
)

type App interface {
	Run() error
}

type AppImpl struct {
	*adapter.Registry
	WishlistController
	HealthController
	ItemController
	ApplicationContext
	router *chi.Mux
}

func (a AppImpl) Run() error {
	APIVersion = "0.0.2"

	utils.Log.Infof("adding middlewares...")

	a.router.Use(middleware.RequestID)
	a.router.Use(middleware.RealIP)
	a.router.Use(chilogger.NewZapMiddleware("router", utils.LLogger))
	a.router.Use(middleware.Recoverer)

	utils.Log.Infof("adding routes...")
	// Health route
	a.router.Get("/health", a.GetHealth)

	// Api V1
	a.router.Route("/api/v1", func(r chi.Router) {
		// Wishlist Endpoint
		r.Route("/wishlist", func(r chi.Router) {
			r.Use(a.AuthCtx)
			r.Get("/", a.GetAllWishlists)
			r.Post("/", a.CreateWishlist)

			r.Route("/{wishlistID}", func(r chi.Router) {
				r.Use(a.WishlistCtx)
				r.Get("/", a.GetWishlist)
				r.Delete("/", a.DeleteWishlist)
				r.Route("/item/{itemSKU}", func(r chi.Router) {
					r.Use(a.ItemCtx)
					r.Post("/", a.AddItemToWishlist)
					r.Delete("/", a.DeleteItemFromWishlist)
				})
			})
		})
		// Item Endpoint
		r.Route("/item", func(r chi.Router) {
			r.Get("/", a.PaginatedListItems)
			r.Get("/search", a.SearchItems)
			r.Route("/{itemSKU}", func(r chi.Router) {
				r.Use(a.ItemCtx)
				r.Get("/", a.GetItem)
			})
		})
	})

	if !viper.GetBool("server.production") {
		// Serve swagger.json
		var cwd string
		var err error
		if cwd, err = os.Getwd(); err != nil {
			log.Panicln(err)
		}

		rootDocs := cwd + "/assets/swagger"
		fs := http.FileServer(http.Dir(rootDocs))

		a.router.Get("/swagger.json", func(w http.ResponseWriter, r *http.Request) {
			if _, err := os.Stat(rootDocs + r.RequestURI); err != nil {
				utils.Log.Error(err)
				render.Status(r, 404)
				return
			} else {
				fs.ServeHTTP(w, r)
			}
		})

		a.router.Get("/swagger/*", httpSwagger.Handler(
			httpSwagger.URL("http://"+viper.GetString("server.host")+":"+viper.GetString("server.port")+"/swagger.json"),
		))
	}

	utils.Log.Infof("running server at %s:%s",
		viper.GetString("server.host"),
		viper.GetString("server.port"),
	)

	err := http.ListenAndServe(fmt.Sprintf("%s:%s",
		viper.GetString("server.host"),
		viper.GetString("server.port")),
		a.router,
	)

	if err != nil {
		return err
	}
	return nil
}

func NewApp(r *adapter.Registry, wController WishlistController, hController HealthController, iController ItemController, ctx ApplicationContext) (App, error) {
	return &AppImpl{
		r,
		wController,
		hController,
		iController,
		ctx,
		chi.NewRouter(),
	}, nil
}
