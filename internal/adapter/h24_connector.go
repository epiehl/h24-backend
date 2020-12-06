package adapter

import (
	"context"
	"fmt"
	"github.com/epiehl93/h24-notifier/internal/utils"
	"github.com/epiehl93/h24-notifier/pkg/models"
	"github.com/epiehl93/h24-notifier/pkg/repository"
	"github.com/shurcooL/graphql"
)

type H24Connector struct {
	*graphql.Client
}

func (h H24Connector) GetBySKU(sku uint64) (*models.Item, error) {
	stringSku := fmt.Sprintf("%018d", sku)
	utils.Log.Infof("Fetching sku %s from retail", stringSku)
	variables := map[string]interface{}{
		"search_sku": graphql.String(stringSku),
	}

	var query struct {
		Articles []struct {
			Name   graphql.String `json:"name"`
			Url    graphql.String `json:"url"`
			Sku    graphql.String `json:"sku"`
			Images []struct {
				Path graphql.String `json:"path"`
			} `json:"images" graphql:"images(limit: 1)"`
			Prices struct {
				Currency graphql.String `json:"currency"`
				Regular  struct {
					Amount graphql.Float `json:"amount" graphql:"amount: value"`
				} `json:"regular"`
				Special struct {
					Amount   graphql.Float `json:"amount" graphql:"amount: value"`
					Discount graphql.Float `json:"discount"`
				} `json:"special"`
			}
		} `graphql:"articles(skus: [$search_sku], locale: de_DE)" json:"articles"`
	}

	err := h.Query(context.Background(), &query, variables)
	if err != nil {
		return nil, err
	}

	var retailPrice float64
	var retailDiscount float64
	var retailDiscountPrice float64

	switch query.Articles[0].Prices.Special.Amount {
	case 0:
		retailPrice = float64(query.Articles[0].Prices.Regular.Amount)
		retailDiscount = 0
		retailDiscountPrice = 0
	default:
		retailDiscountPrice = float64(query.Articles[0].Prices.Special.Amount)
		retailDiscount = float64(query.Articles[0].Prices.Special.Discount)
		retailPrice = retailDiscountPrice + retailDiscount
	}

	return &models.Item{
		Name:                string(query.Articles[0].Name),
		SKU:                 sku,
		ImageUrl:            string(query.Articles[0].Images[0].Path),
		RetailUrl:           string(query.Articles[0].Url),
		RetailPrice:         retailPrice / 100,
		RetailDiscount:      retailDiscount / 100,
		RetailDiscountPrice: retailDiscountPrice / 100,
	}, nil
}

func NewH24Connector(gql *graphql.Client) repository.H24Connector {
	return H24Connector{gql}
}
