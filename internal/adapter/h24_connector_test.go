package adapter

import (
	"github.com/shurcooL/graphql"
	"github.com/stretchr/testify/assert"
	"testing"
)

const graphQlEndpoint = "https://www.home24.de/graphql"

func TestH24Connector_GetBySKU(t *testing.T) {
	client := graphql.NewClient(graphQlEndpoint, nil)
	conn := NewH24Connector(client)

	item, err := conn.GetBySKU(1000091730)

	assert.NoError(t, err)
	assert.NotNil(t, item)

	t.Logf("Name is: %s \n", item.Name)
	t.Logf("Sku is: %d \n", item.SKU)
	t.Logf("Amount is: %f \n", item.RetailPrice)
	t.Logf("RetailDiscountPrice is: %f \n", item.RetailDiscountPrice)
	t.Logf("RetailDiscount is: %f \n", item.RetailDiscount)

	assert.NotNil(t, item.SKU)
	assert.NotEqual(t, item.SKU, 0)

	assert.NotNil(t, item.Name)
	assert.NotEqual(t, item.RetailPrice, 0)

	item, err = conn.GetBySKU(1000000566)

	assert.NoError(t, err)
	assert.NotNil(t, item)

	assert.NotNil(t, item.SKU)
	assert.NotEqual(t, item.SKU, 0)

	assert.NotNil(t, item.Name)
	assert.NotEqual(t, item.RetailPrice, 0)

	t.Logf("Name is: %s \n", item.Name)
	t.Logf("Sku is: %d \n", item.SKU)
	t.Logf("Amount is: %f \n", item.RetailPrice)
	t.Logf("RetailDiscountPrice is: %f \n", item.RetailDiscountPrice)
	t.Logf("RetailDiscount is: %f \n", item.RetailDiscount)
}
