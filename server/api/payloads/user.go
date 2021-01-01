package payloads

import (
	"errors"
	"net/http"

	"github.com/UN0wen/pricewatch-vn/server/api/models"
)

// UserRequest is the request payload for the User data model
type UserRequest struct {
	User *models.User `json:"user"`

	ProtectedID string `json:"id"` // override 'id' json to have more control
}

// Bind is the postprocessing for the UserRequest after the request is unmarshalled
func (a *UserRequest) Bind(r *http.Request) error {
	if a.User == nil {
		return errors.New("missing required User fields")
	}
	// just a post-process after a decode..
	a.ProtectedID = "" // unset the protected ID
	return nil
}

// UserResponse is the response payload for the User data model.
type UserResponse struct {
	User    *models.User `json:"user"`
	Elapsed int64        `json:"elapsed"`
}

// NewUserResponse generate a Response for User object
func NewUserResponse(user *models.User) *UserResponse {
	resp := &UserResponse{User: user}

	return resp
}

// Render is preprocessing before the response is marshalled
func (rd *UserResponse) Render(w http.ResponseWriter, r *http.Request) error {
	// Pre-processing before a response is marshalled and sent across the wire
	rd.Elapsed = 10
	return nil
}
