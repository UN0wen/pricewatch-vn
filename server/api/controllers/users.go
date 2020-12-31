package controllers

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/UN0wen/pricewatch-vn/server/api/models"
	"github.com/UN0wen/pricewatch-vn/server/api/payloads"
	"github.com/UN0wen/pricewatch-vn/server/utils"
	"github.com/go-chi/render"
	"github.com/google/uuid"
)

// getClaims will extract the authorization token from a request and get the associated claims for that id.
func getClaims(r *http.Request) (claimString string, err error) {
	tokenString := r.Header.Get("Authorization")

	if tokenString == "" {
		err = errors.New("No Authorization token")
		return
	}
	claims := utils.ExtractClaims(tokenString[len("Bearer "):])
	claimString = fmt.Sprintf("%v", claims["id"])
	return
}

// SessionCtx middleware is used to load an Session object from
// the authorization token passed through as the request. In case
// the User could not be found, we stop here and return a 404.
func SessionCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var session models.Session
		var err error
		sessionID, err := getClaims(r)
		if err != nil {
			render.Render(w, r, payloads.ErrUnauthorized(err))
			return
		}

		if sessionID != "" {
			var id uuid.UUID
			id, err = uuid.Parse(sessionID)
			if err != nil {
				render.Render(w, r, payloads.ErrNotFound)
				return
			}
			session, err = models.LayerInstance().Session.GetByID(id)
		} else {
			render.Render(w, r, payloads.ErrUnauthorized(err))
			return
		}

		if err != nil {
			render.Render(w, r, payloads.ErrUnauthorized(err))
			return
		}

		ctx := context.WithValue(r.Context(), "userID", session.UserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetUser returns a specific User.
func GetUser(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(uuid.UUID)

	user, err := models.LayerInstance().User.GetByID(userID)

	if err != nil {
		render.Render(w, r, payloads.ErrInternalError(err))
		return
	}

	if err := render.Render(w, r, payloads.NewUserResponse(&user)); err != nil {
		render.Render(w, r, payloads.ErrRender(err))
		return
	}
}

// UpdateUser updates an user password/username
// It returns the updated user
func UpdateUser(w http.ResponseWriter, r *http.Request) {
	data := &payloads.UserRequest{}
	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, payloads.ErrInvalidRequest(err))
		return
	}
	user := data.User

	userID := r.Context().Value("userID").(uuid.UUID)
	updatedUser, err := models.LayerInstance().User.Update(userID, *user)

	if err != nil {
		render.Render(w, r, payloads.ErrInternalError(err))
		return
	}

	if err := render.Render(w, r, payloads.NewUserResponse(&updatedUser)); err != nil {
		render.Render(w, r, payloads.ErrRender(err))
		return
	}
}

// DeleteUser deletes a user from the database
func DeleteUser(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(uuid.UUID)
	err := models.LayerInstance().User.DeleteByID(userID)

	if err != nil {
		render.Render(w, r, payloads.ErrInternalError(err))
		return
	}

	render.Status(r, http.StatusOK)
}

// CreateUser creates a new User and returns it
// back to the client as an acknowledgement. It generates a new
// uuid when called.
func CreateUser(w http.ResponseWriter, r *http.Request) {
	data := &payloads.UserRequest{}
	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, payloads.ErrInvalidRequest(err))
		return
	}

	// Creating a new User
	user := data.User
	userID, err := models.LayerInstance().User.Insert(*user)

	if err != nil {
		render.Render(w, r, payloads.ErrInvalidRequest(err))
		return
	}
	user.ID = userID
	render.Status(r, http.StatusCreated)
	render.Render(w, r, payloads.NewUserResponse(user))
}

// LoginUser logs in a new user based on the provided credentials
func LoginUser(w http.ResponseWriter, r *http.Request) {
	data := &payloads.UserRequest{}

	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, payloads.ErrInvalidRequest(err))
		return
	}

	// Find the user and return the found user
	user := data.User
	found, err := models.LayerInstance().User.Login(*user)

	// Failed login
	if err != nil {
		render.Render(w, r, payloads.ErrUnauthorized(err))
		return
	}

	// Creates new session for User or return the current session

	var sessionID uuid.UUID
	// Try to get current session
	// TODO: cache
	sessions, err := models.LayerInstance().Session.Get(models.SessionQuery{UserID: found.ID})

	// Session not found, create new session
	if err != nil || len(sessions) == 0 {
		sessionID, err = models.LayerInstance().Session.Insert(models.Session{UserID: found.ID})
	} else {
		sessionID = sessions[0].ID
	}
	render.Status(r, http.StatusOK)
	render.Render(w, r, payloads.NewSessionResponse(sessionID, &found))
}
