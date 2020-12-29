package main

import (
	"net/url"

	"github.com/UN0wen/pricewatch-vn/server/scraper"
	"github.com/UN0wen/pricewatch-vn/server/utils"
	"github.com/pkg/errors"
)

func main() {
	// Setup DB

	// if layer := models.LayerInstance(); layer == nil {
	// 	utils.Sugar.Fatalf("Cannot connect to database at %s:%s", utils.DBHost, utils.DBPort)
	// }

	// // Setup Routes
	// router := router.NewRouter()

	// addr := fmt.Sprintf(":%s", utils.ServerPort)
	// server := &http.Server{
	// 	Addr:         addr,
	// 	Handler:      router,
	// 	ReadTimeout:  10 * time.Second,
	// 	WriteTimeout: 10 * time.Second,
	// }

	// utils.Sugar.Infof("Started server on port %s", utils.ServerPort)
	// utils.Sugar.Fatal(server.ListenAndServe())

	scraper.Instance()
	sanitized, err := url.Parse("https://tiki.vn/bo-doi-kem-duong-chong-lao-hoa-da-ngay-va-dem-pond-s-age-miracle-50g-hu-p26709162.html?spid=26709163&src=lp-79")
	if err != nil {
		err = errors.Wrapf(err, "Invalid URL provided")
		return
	}
	item, err := scraper.Instance().Scrapers[sanitized.Host].ScrapeInfo(sanitized)
	utils.CheckError(err)

	utils.Sugar.Infof("%s", item)

	itemPrice, err := scraper.Instance().Scrapers[sanitized.Host].ScrapePrice(item)
	utils.CheckError(err)

	utils.Sugar.Infof("%s", itemPrice)
}
