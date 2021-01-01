package payloads

import (
	"net/http"
	"time"

	"github.com/UN0wen/pricewatch-vn/server/api/models"
)

// SessionResponse is the response payload for the Session data model.
type SessionResponse struct {
	User *models.User `json:"user"`
	JWT  string       `json:"jwt"`
}

// NewSessionResponse sends a user as well as the jwt token
func NewSessionResponse(jwt string, user *models.User) *SessionResponse {
	resp := &SessionResponse{User: user, JWT: jwt}

	return resp
}

// Render is preprocessing before the response is marshalled
func (rd *SessionResponse) Render(w http.ResponseWriter, r *http.Request) error {
	jwtCookie := http.Cookie{
		Name:    "jwt",
		Value:   rd.JWT,
		Expires: time.Now().Add(time.Hour * 24),
	}

	http.SetCookie(w, &jwtCookie)
	return nil
}
