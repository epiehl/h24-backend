package repository

import "github.com/epiehl93/h24-notifier/pkg/models"

type H24Connector interface {
	GetBySKU(sku uint64) (*models.Item, error)
}
