package scraper

import (
	"net/http"
	"net/url"
	"sync"

	"github.com/PuerkitoBio/goquery"
	"github.com/UN0wen/pricewatch-vn/server/api/models"
	"github.com/pkg/errors"
)

// ConfigRoot is the folder where scraper configs are stored
const ConfigRoot = "./scraper/configs/"

// InStockHTTPS is the standard in stock enum in HTTPS
const InStockHTTPS = "https://schema.org/InStock"

// InStockHTTP is the standard in stock enum in HTTP
const InStockHTTP = "http://schema.org/InStock"

var client = &http.Client{}

type scraper struct {
	Scrapers map[string]Scraper
}

// Scraper is an interface implemented by all Scrapers
type Scraper interface {
	ScrapeInfo(path *url.URL) (item models.Item, err error)
	ScrapePrice(item models.Item) (itemPrice models.ItemPrice, err error)
	GetHost() (host string)
}

// Register scraper singletons here
var scrapers = []Scraper{
	LazadaScraper{},
	TikiScraper{},
}

// Singleton reference to the model layer.
var instance *scraper

// Lock for running only once.
var once sync.Once

// Instance gets the static singleton reference
// using double check synchronization.
// It returns the reference to the scraper instance.
func Instance() *scraper {
	once.Do(func() {
		var err error
		scraperMap := make(map[string]Scraper)
		if err != nil {
			err = errors.Wrapf(err, "Could not set up Scraper configs")
			return
		}

		for _, scraper := range scrapers {
			scraperMap[scraper.GetHost()] = scraper
		}

		instance = &scraper{Scrapers: scraperMap}
	})

	return instance
}

// GetDocument returns the goquery document from an URL
func GetDocument(sanitized *url.URL) (doc *goquery.Document, err error) {
	req, err := http.NewRequest("GET", "https://"+sanitized.Host+sanitized.Path, nil)
	if err != nil {
		errors.Wrapf(err, "The external server can't be reached")
		return
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)")
	resp, err := client.Do(req)

	if err != nil {
		err = errors.Wrapf(err, "The external server can't be reached")
		return
	}

	defer resp.Body.Close()

	doc, err = goquery.NewDocumentFromReader(resp.Body)

	if err != nil {
		err = errors.Wrapf(err, "Cannot parse shopping site's response HTML")
	}

	return
}
