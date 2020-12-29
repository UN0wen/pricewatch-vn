package scraper

import (
	"encoding/json"
	"fmt"
	"html"
	"net/url"

	"github.com/UN0wen/pricewatch-vn/server/api/models"
	"github.com/pkg/errors"
)

// LazadaScraper is an empty struct to hold methods that implements Scraper
type LazadaScraper struct{}

// ScrapeInfo extracts the required information out of a page from the scraper config
func (s LazadaScraper) ScrapeInfo(path *url.URL) (item models.Item, err error) {
	// sanitized, err := url.Parse(path)
	// if err != nil {
	// 	err = errors.Wrapf(err, "Invalid URL provided")
	// 	return
	// }

	doc, err := GetDocument(path)

	if err != nil {
		return
	}

	// Name
	name, exists := doc.Find("meta[name=\"og:title\"]").Attr("content")

	if exists {
		name = html.UnescapeString(name)
		item.Name = name
	}

	// Description
	description, exists := doc.Find("meta[name=\"description\"]").Attr("content")

	if exists {
		description = html.UnescapeString(description)
		item.Description = description
	}

	// ImageURL
	imageURL, exists := doc.Find("meta[name=\"og:image\"]").Attr("content")

	if exists {
		var urlParsed *url.URL
		urlParsed, err = url.Parse(imageURL)
		if err != nil {
			return
		}
		item.ImageURL = urlParsed.Host + urlParsed.Path
	}

	// URL
	item.URL = path.Host + path.Path

	// Currency
	item.Currency = "VND"
	return
}

// ScrapePrice returns the current price for an item
func (s LazadaScraper) ScrapePrice(item models.Item) (itemPrice models.ItemPrice, err error) {
	var price int64
	var available bool

	sanitized, err := url.Parse(item.URL)
	if err != nil {
		err = errors.Wrapf(err, "Invalid URL provided")
		return
	}

	doc, err := GetDocument(sanitized)

	jsonData := doc.Find("script[type=\"application/ld+json\"]").First().Text()

	priceJSON := make(map[string]interface{})

	err = json.Unmarshal([]byte(jsonData), &priceJSON)

	if err != nil {
		err = errors.Wrapf(err, "Cannot parse json from Lazada from URL %s", item.URL)
		return
	}

	if offers, ok := priceJSON["offers"]; ok {
		offersMap := offers.(map[string]interface{})
		if low, ok := offersMap["lowPrice"]; ok {
			price = int64(low.(float64))
		} else if high, ok := offersMap["highPrice"]; ok {
			price = int64(high.(float64))
		}

		if avail, ok := offersMap["availability"]; ok {
			available = avail.(string) == InStockHTTPS
		}
	}

	if price == int64(0) {
		err = errors.New(fmt.Sprintf("Cannot parse price for Lazada with url %s", item.URL))
		return
	}
	itemPrice.Price = price
	itemPrice.Available = available
	return
}

// GetHost returns the host name for the scraper
func (s LazadaScraper) GetHost() (host string) {
	host = "www.lazada.vn"
	return
}
