package controllers

import (
	"errors"
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
		render.Render(w, r, payloads.ErrNotFound)
		return
	}

	item, err = models.LayerInstance().Item.GetByID(itemID)

	if err != nil {
		render.Render(w, r, payloads.ErrNotFound)
		return
	}

	if err := render.Render(w, r, payloads.NewItemResponse(&item)); err != nil {
		render.Render(w, r, payloads.ErrRender(err))
		return
	}
}

// GetItemWithPrice returns a specific Item with its latest price.
func GetItemWithPrice(w http.ResponseWriter, r *http.Request) {
	var err error
	itemIDParam := chi.URLParam(r, "itemID")
	itemID, err := uuid.Parse(itemIDParam)

	if err != nil {
		render.Render(w, r, payloads.ErrNotFound)
		return
	}

	itemWithPrice, err := models.LayerInstance().Item.GetWithPrice(itemID)

	if err != nil {
		render.Render(w, r, payloads.ErrNotFound)
		return
	}

	if err := render.Render(w, r, payloads.NewItemWithPriceResponse(&itemWithPrice)); err != nil {
		render.Render(w, r, payloads.ErrRender(err))
		return
	}
}

// SearchItems returns items with prices matching the search query
func SearchItems(w http.ResponseWriter, r *http.Request) {
	var err error
	searchQuery := r.URL.Query().Get("q")

	if err != nil {
		render.Render(w, r, payloads.ErrNotFound)
		return
	} else if searchQuery == "" {
		render.Render(w, r, payloads.ErrInvalidRequest(errors.New("Query must exists")))
		return
	}

	itemsWithPrice, err := models.LayerInstance().Item.Search(searchQuery)

	if err != nil {
		render.Render(w, r, payloads.ErrNotFound)
		return
	}

	if err := render.RenderList(w, r, payloads.NewItemWithPriceListResponse(itemsWithPrice)); err != nil {
		render.Render(w, r, payloads.ErrRender(err))
		return
	}
}

// GetItems returns all items.
func GetItems(w http.ResponseWriter, r *http.Request) {
	items, err := models.LayerInstance().Item.GetAll()

	if err != nil {
		render.Render(w, r, payloads.ErrInternalError(err))
		return
	}

	if err := render.RenderList(w, r, payloads.NewItemListResponse(items)); err != nil {
		render.Render(w, r, payloads.ErrRender(err))
		return
	}
}

// GetItemsWithPrice returns all items with prices.
func GetItemsWithPrice(w http.ResponseWriter, r *http.Request) {
	items, err := models.LayerInstance().Item.GetAllWithPrice()

	if err != nil {
		render.Render(w, r, payloads.ErrInternalError(err))
		return
	}

	if err := render.RenderList(w, r, payloads.NewItemWithPriceListResponse(items)); err != nil {
		render.Render(w, r, payloads.ErrRender(err))
		return
	}
}

// CreateItem creates a new item and then return the item if it is successful
// It expects an URL that it can use to parse into an item object
func CreateItem(w http.ResponseWriter, r *http.Request) {
	data := &payloads.ItemRequest{}
	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, payloads.ErrInvalidRequest(err))
		return
	}
	item := data.Item

	// insert new item
	returnedItem, err := models.LayerInstance().Item.Insert(*item)
	if err != nil {
		render.Render(w, r, payloads.ErrInternalError(err))
		return
	}

	// add item to userItems
	userID := r.Context().Value("userID").(uuid.UUID)
	_, err = models.LayerInstance().UserItem.Insert(models.UserItem{UserID: userID, ItemID: returnedItem.ID})
	if err != nil {
		render.Render(w, r, payloads.ErrInternalError(err))
		return
	}

	item.ID = returnedItem.ID
	if err := render.Render(w, r, payloads.NewItemResponse(item)); err != nil {
		render.Render(w, r, payloads.ErrRender(err))
		return
	}
}

// CreateItemFromURL creates a new item and then return the item if it is successful
// It expects an URL that it can use to parse into an item object
func CreateItemFromURL(w http.ResponseWriter, r *http.Request) {
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

	if err := render.Render(w, r, payloads.NewItemResponse(item)); err != nil {
		render.Render(w, r, payloads.ErrRender(err))
		return
	}
}

// ValidateURL checks if an URL is supported
func ValidateURL(w http.ResponseWriter, r *http.Request) {
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

	// Check if corresponding scraper exists
	if _, ok := scraper.Instance().Scrapers[path.Host]; ok {
		render.Status(r, http.StatusOK)
	} else {
		render.Render(w, r, payloads.ErrNotImplemented)
		return
	}

}
