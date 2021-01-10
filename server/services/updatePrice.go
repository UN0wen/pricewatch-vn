package services

import (
	"net/url"
	"sync"

	"github.com/UN0wen/pricewatch-vn/server/api/models"
	"github.com/UN0wen/pricewatch-vn/server/scraper"
	"github.com/UN0wen/pricewatch-vn/server/utils"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

// TODO: service to update the price of items every ~30mins

// PriceRise is returned if the item's price rose
// PriceFall is returned if the item's price fell
// PriceUnchanged is returned if the price is unchanged
const (
	PriceUnchanged = iota
	PriceFall
	PriceRise
)

// result of the update
type result struct {
	itemID      uuid.UUID
	priceChange int
	err         error
}

// UpdateOne takes an item, then scrapes the URL and return an updated variable
func UpdateOne(item models.Item) (updated int, err error) {
	path, err := url.Parse(item.URL)
	s := scraper.Instance().Scrapers[path.Host]

	oldItemPrices, err := models.LayerInstance().ItemPrice.GetAllPrices(item.ID)

	if err != nil {
		err = errors.Wrapf(err, "Could not find the current price for item with url %s", item.URL)
		return
	}

	itemPrice, err := s.ScrapePrice(item)
	if err != nil {
		err = errors.Wrapf(err, "Could not scrape the price for item with url %s", item.URL)
		return
	}

	if len(oldItemPrices) == 0 {
		updated = PriceRise
	} else {
		oldItemPrice := oldItemPrices[0]
		switch {
		case itemPrice.Price == oldItemPrice.Price:
			updated = PriceUnchanged
		case itemPrice.Price > oldItemPrice.Price:
			updated = PriceRise
		case itemPrice.Price < oldItemPrice.Price:
			updated = PriceFall
		}
	}

	utils.Sugar.Infof("%v", itemPrice)
	if updated > 0 {
		itemPrice.ItemID = item.ID
		_, err = models.LayerInstance().ItemPrice.Insert(itemPrice)

		if err != nil {
			err = errors.Wrapf(err, "Could not insert new item price for item with url %s", item.URL)
		}
	}
	return
}

// UpdateAll tries to go over every item in the database and update them
func UpdateAll() (err error) {
	items, err := models.LayerInstance().Item.GetAll()

	if err != nil {
		err = errors.Wrap(err, "Could not get all items for UpdateAll")
		return
	}

	var wg sync.WaitGroup
	var results []result
	ch := make(chan result)

	for _, item := range items {
		wg.Add(1)
		go produce(ch, &wg, item)
	}

	go func() {
		for v := range ch {
			results = append(results, v)
		}
	}()

	wg.Wait()
	close(ch)

	utils.Sugar.Infof("UpdateAll finished with results:")

	for _, res := range results {
		utils.Sugar.Infof("%s: %d, %s", res.itemID, res.priceChange, res.err)
		// Send email
		if err == nil && res.priceChange == PriceFall {

		}
	}

	return
}

// concurrency functions

func produce(ch chan result, wg *sync.WaitGroup, item models.Item) {
	defer wg.Done()

	updated, err := UpdateOne(item)

	ch <- result{
		itemID:      item.ID,
		priceChange: updated,
		err:         err,
	}
}
