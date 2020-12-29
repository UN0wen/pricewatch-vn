package controllers

import (
	"net/http"
	"net/url"

	"github.com/UN0wen/pricewatch-vn/server/api/models"
	"github.com/UN0wen/pricewatch-vn/server/api/payloads"
	"github.com/UN0wen/pricewatch-vn/server/scraper"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/google/uuid"
)

// GetItem returns a specific Item with id.
func GetItem(w http.ResponseWriter, r *http.Request) {
	var item models.Item
	var err error
	itemIDParam := chi.URLParam(r, "itemID")
	itemID, err := uuid.Parse(itemIDParam)

	if err != nil {
		item, err = models.LayerInstance().Item.GetByID(itemID)
	} else {
		render.Render(w, r, payloads.ErrNotFound)
		return
	}
	if err != nil {
		render.Render(w, r, payloads.ErrNotFound)
		return
	}

	if err := render.Render(w, r, payloads.NewItemResponse(&item)); err != nil {
		render.Render(w, r, payloads.ErrRender(err))
		return
	}
}

// GetItems returns all items.
func GetItems(w http.ResponseWriter, r *http.Request) {
	items, err := models.LayerInstance().Item.Get(models.ItemQuery{})

	if err != nil {
		render.Render(w, r, payloads.ErrInternalError(err))
		return
	}

	if err := render.RenderList(w, r, payloads.NewItemListResponse(items)); err != nil {
		render.Render(w, r, payloads.ErrRender(err))
		return
	}
}

// CreateItem creates a new item and then return the item if it is successful
// It expects an URL that it can use to parse into an item object
// TODO: Add item to useritems
func CreateItem(w http.ResponseWriter, r *http.Request) {
	data := &payloads.ItemRequest{}
	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, payloads.ErrInvalidRequest(err))
		return
	}
	item := data.Item

	path, err := url.Parse(item.URL)
	if err != nil {
		render.Render(w, r, payloads.ErrInvalidRequest(err))
		return
	}

	// Get new Item
	if correspondingScraper, ok := scraper.Instance().Scrapers[path.Host]; ok {
		*item, err = correspondingScraper.ScrapeInfo(path)
		if err != nil {
			render.Render(w, r, payloads.ErrInternalError(err))
			return
		}
	} else {
		render.Render(w, r, payloads.ErrNotImplemented)
		return
	}

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
