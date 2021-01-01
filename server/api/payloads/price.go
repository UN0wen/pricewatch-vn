package payloads

import (
	"net/http"

	"github.com/UN0wen/pricewatch-vn/server/api/models"
	"github.com/go-chi/render"
)

// ItemPriceResponse is the response payload for the ItemPrice data model.
type ItemPriceResponse struct {
	ItemPrice *models.ItemPrice `json:"price"`
}

// NewItemPriceResponse generate a Response for Item object
func NewItemPriceResponse(itemPrice *models.ItemPrice) *ItemPriceResponse {
	resp := &ItemPriceResponse{ItemPrice: itemPrice}

	return resp
}

// NewItemPriceListResponse generates a list of renders for Items
func NewItemPriceListResponse(itemPrices []models.ItemPrice) []render.Renderer {
	list := []render.Renderer{}
	for i := range itemPrices {
		list = append(list, NewItemPriceResponse(&itemPrices[i]))
	}

	return list
}

// Render is preprocessing before the response is marshalled
func (rd *ItemPriceResponse) Render(w http.ResponseWriter, r *http.Request) error {
	// Pre-processing before a response is marshalled and sent across the wire
	return nil
}
