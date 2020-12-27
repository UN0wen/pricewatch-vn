package router

import (
	"net/http"
	"os"

	"github.com/UN0wen/pricewatch-vn/server/api/controllers"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
)

func createUserRoutes(r *chi.Mux) {
	r.Route("/api/users", func(r chi.Router) {
		r.With(controllers.UserCtx).Get("/", controllers.GetUser) // Get /users
	})
}

// NewRouter creates a chi Router with all routes and middleware configured
func NewRouter() *chi.Mux {
	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)
	router.Use(render.SetContentType(render.ContentTypeJSON))
	FileServer(router, "./frontend")

	createUserRoutes(router)
	return router
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
