package payloads

import (
	"errors"
	"net/http"

	"github.com/UN0wen/pricewatch-vn/server/api/models"
	"github.com/go-chi/render"
)

// ItemRequest is the request payload for the Item data model
type ItemRequest struct {
	Item *models.Item `json:"item"`
}

// Bind is the postprocessing for the ItemRequest after the request is unmarshalled
func (a *ItemRequest) Bind(r *http.Request) error {
	if a.Item == nil {
		return errors.New("missing required Item fields")
	}
	return nil
}

// ItemResponse is the response payload for the Item data model.
type ItemResponse struct {
	Item *models.Item `json:"item"`
}

// NewItemResponse generate a Response for Item object
func NewItemResponse(item *models.Item) *ItemResponse {
	resp := &ItemResponse{Item: item}

	return resp
}

// NewItemListResponse generates a list of renders for Items
func NewItemListResponse(items []models.Item) []render.Renderer {
	list := []render.Renderer{}
	for i := range items {
		list = append(list, NewItemResponse(&items[i]))
	}

	return list
}

// Render is preprocessing before the response is marshalled
func (rd *ItemResponse) Render(w http.ResponseWriter, r *http.Request) error {
	// Pre-processing before a response is marshalled and sent across the wire
	return nil
}

// ItemWithPriceResponse is the response payload for the ItemWithPrice data model.
type ItemWithPriceResponse struct {
	ItemWithPrice *models.ItemWithPrice `json:"item_with_price"`
}

// NewItemWithPriceResponse generate a Response for ItemWithPrice object
func NewItemWithPriceResponse(item *models.ItemWithPrice) *ItemWithPriceResponse {
	resp := &ItemWithPriceResponse{ItemWithPrice: item}

	return resp
}

// NewItemWithPriceListResponse generates a list of renders for ItemWithPrices
func NewItemWithPriceListResponse(items []models.ItemWithPrice) []render.Renderer {
	list := []render.Renderer{}
	for i := range items {
		list = append(list, NewItemWithPriceResponse(&items[i]))
	}

	return list
}

// Render is preprocessing before the response is marshalled
func (rd *ItemWithPriceResponse) Render(w http.ResponseWriter, r *http.Request) error {
	// Pre-processing before a response is marshalled and sent across the wire
	return nil
}
