package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/UN0wen/pricewatch-vn/server/api/models"
	"github.com/UN0wen/pricewatch-vn/server/router"
	"github.com/UN0wen/pricewatch-vn/server/scraper"
	"github.com/UN0wen/pricewatch-vn/server/services"
	"github.com/UN0wen/pricewatch-vn/server/utils"
)

func main() {
	// Setup DB

	if layer := models.LayerInstance(); layer == nil {
		utils.Sugar.Fatalf("Cannot connect to database at %s:%s", utils.DBHost, utils.DBPort)
	}

	// Setup the Scraper
	if scraper := scraper.Instance(); scraper == nil {
		utils.Sugar.Fatalf("Cannot initialize the scrapers")
	}

	// Setup Routes
	router := router.NewRouter()

	addr := fmt.Sprintf(":%s", utils.ServerPort)
	server := &http.Server{
		Addr:         addr,
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	err := services.UpdateAll()
	utils.CheckError(err)
	utils.Sugar.Infof("Started server on port %s", utils.ServerPort)
	utils.Sugar.Fatal(server.ListenAndServe())
}
