package web

import (
	"fmt"
	"github.com/epiehl93/h24-notifier/config"
	"github.com/epiehl93/h24-notifier/internal/adapter"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/gorilla/sessions"
	"github.com/markbates/goth/gothic"
	"github.com/shurcooL/graphql"
	httpSwagger "github.com/swaggo/http-swagger"
	"gorm.io/gorm"
	"log"
	"net/http"
	"os"
)

type App interface {
	Run() error
}

type app struct {
	adapter.Registry
	WishlistController
	HealthController
	ItemController
	ApplicationContext
	router *chi.Mux
}

func (a app) Run() error {
	APIVersion = "0.0.2"

	log.Println("configuring auth...")

	gothic.Store = sessions.NewCookieStore([]byte("mysuperhardcodedsecret"))
	log.Println("adding middlewares...")

	a.router.Use(middleware.RequestID)
	a.router.Use(middleware.RealIP)
	a.router.Use(middleware.Logger)
	a.router.Use(middleware.Recoverer)

	log.Println("adding routes...")
	// Health route
	a.router.Get("/health", a.GetHealth)

	// Api V1
	a.router.Route("/api/v1", func(r chi.Router) {
		r.Use(a.AuthCtx)
		// Wishlist Endpoint
		r.Route("/wishlist", func(r chi.Router) {
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
			r.Get("/", a.ListItems)
			r.Route("/{itemSKU}", func(r chi.Router) {
				r.Use(a.ItemCtx)
				r.Get("/", a.GetItem)
			})
		})
	})

	if !config.C.Server.Production {
		// Serve swagger.json
		var cwd string
		var err error
		if cwd, err = os.Getwd(); err != nil {
			log.Panicln(err)
		}

		rootDocs := cwd + "/assets/swagger"
		fs := http.FileServer(http.Dir(rootDocs))
		fmt.Println(fs)

		a.router.Get("/swagger.json", func(w http.ResponseWriter, r *http.Request) {
			if _, err := os.Stat(rootDocs + r.RequestURI); err != nil {
				fmt.Println(err)
				render.Status(r, 404)
				return
			} else {
				fs.ServeHTTP(w, r)
			}
		})

		a.router.Get("/swagger/*", httpSwagger.Handler(
			httpSwagger.URL("http://"+config.C.Server.Host+":"+config.C.Server.Port+"/swagger.json"),
		))
	}

	log.Println("starting server...")
	err := http.ListenAndServe(fmt.Sprintf("%s:%s", config.C.Server.Host, config.C.Server.Port), a.router)
	if err != nil {
		return err
	}

	log.Println("done...")
	return nil
}

func NewApp(db *gorm.DB, gql *graphql.Client) (App, error) {
	r := adapter.NewRegistry(db, gql)
	return &app{
		r,
		NewWishlistController(r),
		NewHealthController(r),
		NewItemController(r),
		NewApplicationContext(r),
		chi.NewRouter(),
	}, nil
}
