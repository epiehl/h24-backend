package web

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/epiehl93/h24-notifier/config"
	"github.com/epiehl93/h24-notifier/internal/adapter"
	"github.com/epiehl93/h24-notifier/pkg/models"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"gopkg.in/square/go-jose.v2"
	"gorm.io/gorm"
	"log"
	"net/http"
	"strconv"
	"strings"
)

const (
	// Key in the Headers hashmap of the token that points to the key ID
	keyIdTokenHeaderKey = "kid"
)

type JWK struct {
	Algorithm string `json:"alg"`
	E         string `json:"e"`
	KeyID     string `json:"kid"`
	Kty       string `json:"kty"`
	N         string `json:"n"`
	Use       string `json:"use"`
}

type JWKS struct {
	Keys []JWK `json:"keys"`
}

type CognitoClaims struct {
	Scope string `json:"scope"`
	Sub   string `json:"sub"`
}

func (cc CognitoClaims) Valid() error {
	return nil
}

type ApplicationContext interface {
	WishlistCtx(next http.Handler) http.Handler
	ItemCtx(next http.Handler) http.Handler
	AuthCtx(next http.Handler) http.Handler
}

type applicationContext struct {
	adapter.Registry
}

// WishlistCtx tries to pull the wishlist object from the repository by querying with the supplied ID
func (c applicationContext) WishlistCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		userUUID := ctx.Value("userUUID").(string)

		wishlistID := chi.URLParam(r, "wishlistID")
		id, err := strconv.ParseUint(wishlistID, 10, 64)
		if err != nil {
			_ = render.Render(w, r, ErrRender(err))
			log.Println(err)
			return
		}

		wlist := &models.Wishlist{ID: id}
		err = c.WishlistRepository.Get(wlist)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				_ = render.Render(w, r, ErrNotFound())
				log.Println(err)
				return
			} else {
				_ = render.Render(w, r, ErrRender(err))
				return
			}
		}

		if wlist.UserSub != userUUID {
			_ = render.Render(w, r, ErrNotFound())
			return
		}

		newCtx := context.WithValue(ctx, "wishlist", wlist)
		next.ServeHTTP(w, r.WithContext(newCtx))
	})
}

// ItemCtx pulls an item from the repository via sku
func (c applicationContext) ItemCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		skuParam := chi.URLParam(r, "itemSKU")
		sku, err := strconv.ParseUint(skuParam, 10, 64)
		if err != nil {
			_ = render.Render(w, r, ErrRender(err))
			return
		}

		item, err := c.ItemRepository.GetBySKU(sku)
		if err != nil {
			// Try to fetch article from retail if not found in db
			if errors.Is(err, gorm.ErrRecordNotFound) {
				item, err = c.H24Connector.GetBySKU(sku)
				if err := c.ItemRepository.Create(item); err != nil {
					_ = render.Render(w, r, ErrRender(err))
					return
				}
			} else {
				_ = render.Render(w, r, ErrRender(err))
				return
			}
		}

		ctx := context.WithValue(r.Context(), "item", item)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (c applicationContext) AuthCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userUUID := ""

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		w.Header().Set("Access-Control-Allow-Methods", "*")

		// Handle pre-flight requests
		if r.Method == http.MethodOptions {
			w.WriteHeader(200)
			return
		}

		accessHeader := r.Header.Get("Authorization")
		if accessHeader == "" {
			_ = render.Render(w, r, ErrUnauthorized(errors.New("authorization header is empty")))
			return
		}
		accessToken := strings.Split(accessHeader, " ")[1]
		verifiedToken, err := c.RetrieveAndValidateToken(accessToken)
		if err != nil {
			_ = render.Render(w, r, ErrUnauthorized(err))
			return
		}

		userUUID = verifiedToken.Claims.(*jwt.StandardClaims).Subject

		if userUUID == "" {
			_ = render.Render(w, r, ErrRender(errors.New("could not find user")))
		}

		ctx := context.WithValue(r.Context(), "userUUID", userUUID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (c applicationContext) RetrieveAndValidateToken(token string) (*jwt.Token, error) {
	// todo: eventually cache that
	jwks := &jose.JSONWebKeySet{}

	// get jwks from url
	resp, err := http.Get(config.C.Auth.JWksUrl)
	defer func(resp *http.Response) {
		_ = resp.Body.Close()
	}(resp)
	if err != nil {
		return nil, err
	}

	if err := json.NewDecoder(resp.Body).Decode(jwks); err != nil {
		return nil, err
	}

	verifiedToken, err := new(jwt.Parser).ParseWithClaims(
		token,
		&jwt.StandardClaims{},
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
				return nil, errors.New(fmt.Sprintf("invalid signing method, got %s", token.Header))
			}

			untypedKeyId, found := token.Header[keyIdTokenHeaderKey]
			if !found {
				return nil, errors.New("no key id")
			}
			keyId, ok := untypedKeyId.(string)
			if !ok {
				return nil, errors.New("found key id but value was not a string")
			}

			keys := jwks.Key(keyId)
			if len(keys) == 0 {
				return nil, errors.New("could not find key")
			}

			key := keys[0]

			if !key.IsPublic() {
				return nil, errors.New("key is not public")
			}

			return key.Key, nil
		},
	)

	if err != nil {
		return verifiedToken, err
	}

	return verifiedToken, nil
}

func NewApplicationContext(registry adapter.Registry) ApplicationContext {
	return applicationContext{registry}
}
