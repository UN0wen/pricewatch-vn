package controllers

import (
	"net/http"

	"github.com/UN0wen/pricewatch-vn/server/api/models"
	"github.com/UN0wen/pricewatch-vn/server/api/payloads"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/google/uuid"
)

// GetPrice returns the most current price for item with id.
func GetPrice(w http.ResponseWriter, r *http.Request) {
	var itemPrice models.ItemPrice
	var err error
	itemIDParam := chi.URLParam(r, "itemID")
	itemID, err := uuid.Parse(itemIDParam)

	if err != nil {
		render.Render(w, r, payloads.ErrNotFound)
		return
	}

	itemPrice, err = models.LayerInstance().ItemPrice.GetPrice(itemID)

	if err != nil {
		render.Render(w, r, payloads.ErrInternalError(err))
		return
	}

	if err := render.Render(w, r, payloads.NewItemPriceResponse(&itemPrice)); err != nil {
		render.Render(w, r, payloads.ErrRender(err))
		return
	}
}

// GetPrices returns all prices for item with id
func GetPrices(w http.ResponseWriter, r *http.Request) {
	var err error
	var itemPrices []models.ItemPrice
	itemIDParam := chi.URLParam(r, "itemID")
	itemID, err := uuid.Parse(itemIDParam)

	if err != nil {
		render.Render(w, r, payloads.ErrNotFound)
		return
	}

	itemPrices, err = models.LayerInstance().ItemPrice.GetAllPrices(itemID)

	if err != nil {
		render.Render(w, r, payloads.ErrInternalError(err))
		return
	}

	if err := render.RenderList(w, r, payloads.NewItemPriceListResponse(itemPrices)); err != nil {
		render.Render(w, r, payloads.ErrRender(err))
		return
	}
}
