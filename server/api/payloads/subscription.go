package payloads

import (
	"errors"
	"net/http"

	"github.com/UN0wen/pricewatch-vn/server/api/models"
)

// SubscriptionRequest is the request payload for the Subscription data model
type SubscriptionRequest struct {
	Subscription *models.Subscription `json:"subscription"`
}

// Bind is the postprocessing for the ItemRequest after the request is unmarshalled
func (a *SubscriptionRequest) Bind(r *http.Request) error {
	if a.Subscription == nil {
		return errors.New("missing required Subscription fields")
	}
	return nil
}
