package controllers

import (
	"net/http"

	"github.com/UN0wen/pricewatch-vn/server/api/models"
	"github.com/UN0wen/pricewatch-vn/server/api/payloads"
	"github.com/go-chi/render"
	"github.com/google/uuid"
)

// GetUserItems gets all items for a certain user
func GetUserItems(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(uuid.UUID)
	userItems, err := models.LayerInstance().UserItem.GetByUser(userID)
	if err != nil {
		render.Render(w, r, payloads.ErrInternalError(err))
		return
	}

	if err := render.RenderList(w, r, payloads.NewUserItemListResponse(userItems)); err != nil {
		render.Render(w, r, payloads.ErrRender(err))
		return
	}
}

// // GetUserItem returns the itemID if the given itemID is in the user's list and nothing otherwise
// func GetUserItem(w http.ResponseWriter, r *http.Request) {
// 	itemIDParam := chi.URLParam(r, "itemID")
// 	itemID, err := uuid.Parse(itemIDParam)

// 	if err != nil {
// 		render.Render(w, r, payloads.ErrNotFound)
// 		return
// 	}

// 	userID := r.Context().Value("userID").(uuid.UUID)
// 	userItems, err := models.LayerInstance().UserItem.Get(models.UserItemQuery{UserID: userID, ItemID: itemID})
// 	if err != nil {
// 		render.Render(w, r, payloads.ErrNotFound)
// 		return
// 	}

// 	if len(userItems) > 0 {
// 		if err = render.Render(w, r, payloads.NewUserItemResponse(&userItems[0])); err != nil {
// 			render.Render(w, r, payloads.ErrRender(err))
// 			return
// 		}
// 	} else {
// 		if err = render.Render(w, r, payloads.NewUserItemResponse(&models.UserItem{})); err != nil {
// 			render.Render(w, r, payloads.ErrRender(err))
// 			return
// 		}
// 	}
// }

// CreateUserItem adds an item to an user
func CreateUserItem(w http.ResponseWriter, r *http.Request) {
	data := &payloads.ItemRequest{}
	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, payloads.ErrInvalidRequest(err))
		return
	}
	item := data.Item

	// add item to userItems
	userID := r.Context().Value("userID").(uuid.UUID)
	_, err := models.LayerInstance().UserItem.Insert(models.UserItem{UserID: userID, ItemID: item.ID})
	if err != nil {
		render.Render(w, r, payloads.ErrInternalError(err))
		return
	}

	render.Status(r, http.StatusOK)
}
