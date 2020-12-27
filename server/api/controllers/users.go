package controllers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/UN0wen/pricewatch-vn/server/api/models"
	"github.com/UN0wen/pricewatch-vn/server/api/payloads"
	"github.com/UN0wen/pricewatch-vn/server/utils"
	"github.com/go-chi/render"
	"github.com/google/uuid"
)

// getClaims will extract the authorization token from a request and get the associated claims for that id.
func getClaims(r *http.Request) string {
	tokenString := r.Header.Get("Authorization")[len("Bearer "):]
	claims := utils.ExtractClaims(tokenString)
	return fmt.Sprintf("%v", claims["id"])
}

// UserCtx middleware is used to load an User object from
// the URL parameters passed through as the request. In case
// the User could not be found, we stop here and return a 404.
func UserCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var user models.User
		var err error
		if userID := getClaims(r); userID != "" {
			id, err := uuid.Parse(userID)
			if err != nil {
				render.Render(w, r, payloads.ErrNotFound)
				return
			}
			user, err = models.LayerInstance().User.GetByID(id)
		} else {
			render.Render(w, r, payloads.ErrNotFound)
			return
		}
		if err != nil {
			render.Render(w, r, payloads.ErrNotFound)
			return
		}

		ctx := context.WithValue(r.Context(), "user", &user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetUser returns a specific User.
func GetUser(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(*models.User)

	if err := render.Render(w, r, payloads.NewUserResponse(user)); err != nil {
		render.Render(w, r, payloads.ErrRender(err))
		return
	}
}

// CreateUser creates a new User and returns it
// back to the client as an acknowledgement.
func CreateUser(w http.ResponseWriter, r *http.Request) {
	data := &payloads.UserRequest{}
	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, payloads.ErrInvalidRequest(err))
		return
	}

	user := data.User
	err := models.LayerInstance().User.Insert(*user)

	if err != nil {
		render.Render(w, r, payloads.ErrInvalidRequest(err))
		return
	}

	render.Status(r, http.StatusCreated)
	render.Render(w, r, payloads.NewUserResponse(user))
}
