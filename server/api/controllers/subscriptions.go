package controllers

import (
	"net/http"

	"github.com/UN0wen/pricewatch-vn/server/api/models"
	"github.com/UN0wen/pricewatch-vn/server/api/payloads"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/google/uuid"
)

// Subscribe subscribes an user to email notifications
// for when the items price goes below a certain price
func Subscribe(w http.ResponseWriter, r *http.Request) {
	data := &payloads.SubscriptionRequest{}
	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, payloads.ErrInvalidRequest(err))
		return
	}
	inSub := data.Subscription

	itemIDParam := chi.URLParam(r, "itemID")
	itemID, err := uuid.Parse(itemIDParam)

	if err != nil {
		render.Render(w, r, payloads.ErrNotFound)
		return
	}

	userID := r.Context().Value("userID").(uuid.UUID)

	sub := models.Subscription{
		UserID:      userID,
		ItemID:      itemID,
		Email:       inSub.Email,
		TargetPrice: inSub.TargetPrice,
	}

	_, err = models.LayerInstance().Subscription.Insert(sub)
	if err != nil {
		render.Render(w, r, payloads.ErrInvalidRequest(err))
		return
	}

	render.Status(r, 200)
}

// Unsubscribe unsubscribes an user from email notifications
func Unsubscribe(w http.ResponseWriter, r *http.Request) {
	itemIDParam := chi.URLParam(r, "itemID")
	itemID, err := uuid.Parse(itemIDParam)

	if err != nil {
		render.Render(w, r, payloads.ErrNotFound)
		return
	}

	userID := r.Context().Value("userID").(uuid.UUID)

	err = models.LayerInstance().Subscription.Delete(userID, itemID)
	if err != nil {
		render.Render(w, r, payloads.ErrInvalidRequest(err))
		return
	}

	render.Status(r, 200)
}
