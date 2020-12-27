package payloads

import (
	"net/http"
	"time"

	"github.com/UN0wen/pricewatch-vn/server/utils"
	"github.com/google/uuid"
)

// SessionResponse is the response payload for the Session data model.
type SessionResponse struct {
	JWT string
}

// NewSessionResponse generate a JWT token for a session ID, then send it
func NewSessionResponse(id uuid.UUID) *SessionResponse {
	resp := &SessionResponse{JWT: utils.GenerateJWT(id)}

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
