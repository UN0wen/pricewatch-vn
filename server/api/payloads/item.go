package payloads

import (
	"errors"
	"net/http"

	"github.com/UN0wen/pricewatch-vn/server/api/models"
)

// ItemRequest is the request payload for the Item data model
type ItemRequest struct {
	*models.Item
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
	*models.Item
}

// NewItemResponse generate a Response for User object
func NewItemResponse(item *models.Item) *ItemResponse {
	resp := &ItemResponse{Item: item}

	return resp
}

// Render is preprocessing before the response is marshalled
func (rd *ItemResponse) Render(w http.ResponseWriter, r *http.Request) error {
	// Pre-processing before a response is marshalled and sent across the wire
	return nil
}
