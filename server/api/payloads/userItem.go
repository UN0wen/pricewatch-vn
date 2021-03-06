package payloads

import (
	"net/http"

	"github.com/UN0wen/pricewatch-vn/server/api/models"
	"github.com/go-chi/render"
	"github.com/google/uuid"
)

// UserItemResponse is the response payload for the ItemPrice data model.
type UserItemResponse struct {
	UserItem *models.UserItem `json:"user_item"`
	Valid    bool             `json:"valid"`
}

// NewUserItemResponse generate a Response for Item object
func NewUserItemResponse(userItem *models.UserItem) *UserItemResponse {
	resp := &UserItemResponse{UserItem: userItem}

	return resp
}

// NewUserItemListResponse generates a list of renders for Items
func NewUserItemListResponse(userItems []models.UserItem) []render.Renderer {
	list := []render.Renderer{}
	for i := range userItems {
		list = append(list, NewUserItemResponse(&userItems[i]))
	}

	return list
}

// Render is preprocessing before the response is marshalled
func (rd *UserItemResponse) Render(w http.ResponseWriter, r *http.Request) error {

	// Valid is true if the itemID is not nil
	rd.Valid = rd.UserItem.ItemID != uuid.Nil
	return nil
}
