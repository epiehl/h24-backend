package aggregator

import (
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/epiehl93/h24-notifier/internal/adapter"
	"github.com/epiehl93/h24-notifier/internal/utils"
	"github.com/epiehl93/h24-notifier/pkg/models"
	"github.com/spf13/viper"
	"gorm.io/gorm"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type OutletAggregator interface {
	Run() error
	CheckUnavailability() error
}

type OutletAggregatorImpl struct {
	*adapter.Registry
}

func (o OutletAggregatorImpl) CheckUnavailability() error {
	items, err := o.Item.FindAvailableInOutlet()
	if err != nil {
		return err
	}

	lastCycle, err := o.Cycle.GetLastSuccessfulCycle(models.AggregationCycle)
	if err != nil {
		return err
	}

	for _, item := range items {
		// Check if it is older than last cycle
		if item.LastAggregatedAt.Before(lastCycle.At) {
			utils.Log.Infof("Item with sku %d: %s older than %s\n", item.SKU, item.LastAggregatedAt, lastCycle.At)
			err := o.Item.SetUnavailableInOutlet(item)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (o OutletAggregatorImpl) Run() error {
	now := time.Now()
	if err := o.AggregateItems(); err != nil {
		return o.ReturnErrorAndSetFailedCycle(err, now)
	}

	if err := o.SetSuccessfulCycle(now); err != nil {
		return err
	}

	if err := o.CheckUnavailability(); err != nil {
		return err
	}

	return nil
}

func (o OutletAggregatorImpl) AggregateItems() error {
	done := false
	siteIndex := 1

	for !done {

		url := fmt.Sprintf("%s://%s/%s?p=%d&is_ajax=1",
			viper.GetString("aggregator.outlet.endpoint.scheme"),
			viper.GetString("aggregator.outlet.endpoint.host"),
			viper.GetString("aggregator.outlet.endpoint.location"),
			siteIndex)

		utils.Log.Infof("Calling url %s \n", url)

		res, err := http.Get(url)
		if err != nil {
			return err
		}

		if res.StatusCode != 200 {
			return errors.New(fmt.Sprintf("request to %s terminated with status code %d", url, res.StatusCode))
		}

		// Load the doc
		doc, err := goquery.NewDocumentFromReader(res.Body)
		if err != nil {
			return err
		}

		selection := doc.Find("div.product-item")
		if selection.Length() == 0 {
			utils.Log.Infof("Response body seems empty. Finishing")
			done = true
			break
		}
		selection.Each(o.ProcessDocument)

		err = res.Body.Close()
		if err != nil {
			return err
		}

		time.Sleep(2000)
		siteIndex++
	}

	return nil
}

func (o OutletAggregatorImpl) ProcessDocument(i int, selection *goquery.Selection) {
	if err := o.ExtractItemFromSelection(selection); err != nil {
		utils.Log.Panic(err)
	}
	utils.Log.Infof("Successfully extracted data")
}

func (o OutletAggregatorImpl) ExtractItemFromSelection(selection *goquery.Selection) error {
	var eName string
	var eSku uint64
	var ePriceOld float64
	var ePriceNew float64
	var newItem *models.Item

	if len(selection.Nodes) > 1 {
		return errors.New("could not find nodes in selection")
	}

	for _, attr := range selection.Nodes[0].Attr {
		switch attr.Key {
		case "prod-name":
			eName = attr.Val
		case "sap-sku":
			if sku, err := ParseSkuFromString(attr.Val); err != nil {
				return err
			} else {
				eSku = sku
			}
		}
	}

	priceOld := selection.Find("div.product-item-price-old").Find("span.price")
	priceNew := selection.Find("div.product-item-price-new").Find("span.price")

	if price, err := ParsePriceFromString(priceOld.Text()); err != nil {
		return err
	} else {
		ePriceOld = price
	}

	if price, err := ParsePriceFromString(priceNew.Text()); err != nil {
		return err
	} else {
		ePriceNew = price
	}

	itemExists := true
	newItem, err := o.Item.GetBySKU(eSku)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.Log.Infof("Item with sku %d not found, creating\n", eSku)
			itemExists = false
		} else {
			return err
		}
	}

	if newItem == nil {
		newItem = &models.Item{}
	}

	newItem.Name = eName
	newItem.SKU = eSku
	newItem.OutletPrice = ePriceNew
	newItem.RetailPrice = ePriceOld
	newItem.LastAggregatedAt = time.Now()

	_, err = o.H24.EnrichItemWithRetailData(newItem)
	if err != nil {
		return err
	}

	if !itemExists {
		newItem.SKU = eSku
		err := o.Registry.Item.Create(newItem)
		if err != nil {
			return err
		}
	} else {
		if err := o.Registry.Item.Update(newItem); err != nil {
			return err
		}
	}

	err = o.Item.SetAvailableInOutlet(newItem)
	if err != nil {
		return err
	}

	return nil
}

func ParsePriceFromString(stringPrice string) (float64, error) {
	regexResult := regexp.MustCompile("(?P<Price>[\\d]+(?:,\\d\\d)?)")
	match := regexResult.FindStringSubmatch(stringPrice)

	s := strings.TrimSpace(strings.Replace(match[1], ",", ".", -1))

	if price, err := strconv.ParseFloat(s, 10); err != nil {
		return 0, err
	} else {
		return price, nil
	}
}

func ParseSkuFromString(sku string) (uint64, error) {
	if i, err := strconv.ParseUint(sku, 10, 64); err != nil {
		return 0, err
	} else {
		return i, nil
	}
}

func (o OutletAggregatorImpl) ReturnErrorAndSetFailedCycle(err error, at time.Time) error {
	e := o.Cycle.Create(&models.Cycle{Successful: false, At: at, Type: models.AggregationCycle})
	if e != nil {
		return e
	}
	utils.Log.Infof("Aggregation was not successful")
	return err
}

func (o OutletAggregatorImpl) SetSuccessfulCycle(at time.Time) error {
	err := o.Cycle.Create(&models.Cycle{Successful: true, At: at, Type: models.AggregationCycle})
	if err != nil {
		return err
	}
	utils.Log.Infof("Aggregation was successful")
	return nil
}

func NewOutletAggregator(r *adapter.Registry) OutletAggregator {
	return &OutletAggregatorImpl{r}
}
