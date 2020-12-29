package scraper

import (
	"html"
	"net/url"
	"strconv"

	"github.com/UN0wen/pricewatch-vn/server/api/models"
	"github.com/UN0wen/pricewatch-vn/server/utils"
	"github.com/pkg/errors"
)

// TikiScraper is an empty struct to hold methods that implements Scraper
type TikiScraper struct{}

// ScrapeInfo extracts the required information out of a page from the scraper config
func (s TikiScraper) ScrapeInfo(path *url.URL) (item models.Item, err error) {
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
	name, exists := doc.Find("meta[property=\"og:title\"]").Attr("content")

	if exists {
		name = html.UnescapeString(name)
		item.Name = name
	}

	// Description
	description, exists := doc.Find("meta[property=\"og:description\"]").Attr("content")

	if exists {
		description = html.UnescapeString(description)
		item.Description = description
	}

	// ImageURL
	imageURL, exists := doc.Find("meta[property=\"og:image\"]").Attr("content")

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
func (s TikiScraper) ScrapePrice(item models.Item) (itemPrice models.ItemPrice, err error) {
	var price int64
	var available bool
	sanitized, err := url.Parse(item.URL)
	if err != nil {
		err = errors.Wrapf(err, "Invalid URL provided")
		return
	}

	doc, err := GetDocument(sanitized)

	priceString, exists := doc.Find("meta[itemProp=\"price\"]").Attr("content")

	if exists {
		price, err = strconv.ParseInt(priceString, 10, 64)

		if err != nil {
			return
		}
	}

	availableString, exists := doc.Find("link[itemProp=\"availability\"]").Attr("href")

	utils.Sugar.Infof("%s", availableString)
	if exists {
		available = availableString == InStockHTTP
	}

	itemPrice.Price = price
	itemPrice.Available = available
	return
}

// GetHost returns the host name for the scraper
func (s TikiScraper) GetHost() (host string) {
	host = "tiki.vn"
	return
}
