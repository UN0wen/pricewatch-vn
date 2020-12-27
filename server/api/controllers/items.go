package controllers

import (
	"errors"
	"net/http"

	"github.com/UN0wen/pricewatch-vn/server/api/models"
	"github.com/UN0wen/pricewatch-vn/server/api/payloads"
	"github.com/go-chi/render"
	"github.com/google/uuid"
)

// GetItem returns a specific Item with id.
func GetItem(w http.ResponseWriter, r *http.Request) {
	data := &payloads.ItemRequest{}
	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, payloads.ErrInvalidRequest(err))
		return
	}
	item := data.Item
	if item.ID == uuid.Nil {
		render.Render(w, r, payloads.ErrInvalidRequest(errors.New("No item ID provided")))
		return
	}

	returnedItem, err := models.LayerInstance().Item.GetByID(item.ID)

	if err != nil {
		render.Render(w, r, payloads.ErrInternalError(err))
		return
	}

	if err := render.Render(w, r, payloads.NewItemResponse(&returnedItem)); err != nil {
		render.Render(w, r, payloads.ErrRender(err))
		return
	}
}

// CreateItem creates a new item and then return the item if it is successful
func CreateItem(w http.ResponseWriter, r *http.Request) {
	data := &payloads.ItemRequest{}
	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, payloads.ErrInvalidRequest(err))
		return
	}
	item := data.Item

	itemID, err := models.LayerInstance().Item.Insert(*item)
	if err != nil {
		render.Render(w, r, payloads.ErrInternalError(err))
		return
	}

	item.ID = itemID
	if err := render.Render(w, r, payloads.NewItemResponse(item)); err != nil {
		render.Render(w, r, payloads.ErrRender(err))
		return
	}
}
