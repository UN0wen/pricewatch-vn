package main

import (
	"net/http"
	"os"
	"time"

	"github.com/UN0wen/pricewatch-vn/server/api/models"
	"github.com/UN0wen/pricewatch-vn/server/utils"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/google/uuid"
)

func main() {
	layer := models.LayerInstance()

	user := models.User{
		ID:       uuid.New(),
		Username: "admin",
		Email:    "admin@gmail.com",
		Password: "password",
	}

	err := layer.User.Insert(user)
	utils.CheckError(err)

	uuid, err := uuid.Parse("05d50b0b-9e3e-405e-8c4b-ce9ce4e6f1aa")
	utils.CheckError(err)
	user2, err := layer.User.Get(models.UserQuery{
		ID: uuid,
	}, "", "=")
	utils.CheckError(err)
	utils.Sugar.Infof("%s", user2)

	router := chi.NewRouter()
	router.Use(middleware.Logger)
	FileServer(router, "./frontend")

	server := &http.Server{
		Addr:         ":8080",
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	server.ListenAndServe()
}

// FileServer is serving static files.
func FileServer(router *chi.Mux, root string) {
	fs := http.FileServer(http.Dir(root))

	router.Get("/*", func(w http.ResponseWriter, r *http.Request) {
		if _, err := os.Stat(root + r.RequestURI); os.IsNotExist(err) {
			http.StripPrefix(r.RequestURI, fs).ServeHTTP(w, r)
		} else {
			fs.ServeHTTP(w, r)
		}
	})
}
