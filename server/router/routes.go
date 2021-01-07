package router

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/UN0wen/pricewatch-vn/server/api/controllers"
	"github.com/UN0wen/pricewatch-vn/server/middleware"
	"github.com/UN0wen/pricewatch-vn/server/utils"
	"github.com/go-chi/chi"
	chimiddleware "github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"
)

func createUserRoutes(r *chi.Mux) {
	r.Route("/api/user", func(r chi.Router) {
		r.With(middleware.Authenticate).With(controllers.SessionCtx).Get("/", controllers.GetUser)       // Get /users
		r.With(middleware.Authenticate).With(controllers.SessionCtx).Put("/", controllers.UpdateUser)    // Update
		r.With(middleware.Authenticate).With(controllers.SessionCtx).Delete("/", controllers.DeleteUser) // Delete

		// UserItems
		r.With(middleware.Authenticate).With(controllers.SessionCtx).Get("/item/{itemID}", controllers.GetUserItem) //
		r.With(middleware.Authenticate).With(controllers.SessionCtx).Get("/items", controllers.GetUserItems)        //
		r.With(middleware.Authenticate).With(controllers.SessionCtx).Post("/item", controllers.CreateUserItem)      //
	})
}

func createItemRoutes(r *chi.Mux) {
	r.Route("/api/items", func(r chi.Router) {
		r.Get("/", controllers.GetItems)
		r.Get("/prices", controllers.GetItemsWithPrice)
		r.Get("/search", controllers.SearchItems)
	})
	r.Route("/api/item", func(r chi.Router) {
		r.With(middleware.Authenticate).With(controllers.SessionCtx).Post("/", controllers.CreateItem) // Create
		r.Get("/{itemID}", controllers.GetItemWithPrice)                                               // Get /users
		r.Get("/{itemID}/price", controllers.GetPrice)
		r.Get("/{itemID}/prices", controllers.GetPrices)
		r.Post("/validate", controllers.ValidateURL)
	})
}

func createAuthRoutes(r *chi.Mux) {
	r.Post("/api/signup", controllers.CreateUser)
	r.Post("/api/login", controllers.LoginUser)
}

// NewRouter creates a chi Router with all routes and middleware configured
func NewRouter() *chi.Mux {
	router := chi.NewRouter()

	router.Use(chimiddleware.RequestID)
	router.Use(chimiddleware.Logger)
	router.Use(chimiddleware.Recoverer)
	router.Use(chimiddleware.URLFormat)
	router.Use(render.SetContentType(render.ContentTypeJSON))
	router.Use(cors.Handler(cors.Options{
		// AllowedOrigins: []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	// Create API routes
	createUserRoutes(router)
	createItemRoutes(router)
	createAuthRoutes(router)

	spa := spaHandler{staticPath: "build", indexPath: "index.html"}
	router.Handle("/*", spa)
	return router
}

// spaHandler implements the http.Handler interface, so we can use it
// to respond to HTTP requests. The path to the static directory and
// path to the index file within that static directory are used to
// serve the SPA in the given static directory.
type spaHandler struct {
	staticPath string
	indexPath  string
}

// ServeHTTP inspects the URL path to locate a file within the static dir
// on the SPA handler. If a file is found, it will be served. If not, the
// file located at the index path on the SPA handler will be served. This
// is suitable behavior for serving an SPA (single page application).
func (h spaHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// get the absolute path to prevent directory traversal
	path, err := filepath.Abs(r.URL.Path)
	if err != nil {
		// if we failed to get the absolute path respond with a 400 bad request
		// and stop
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	utils.Sugar.Infof("%s", path)
	// prepend the path with the path to the static directory
	path = filepath.Join(h.staticPath, path)

	// check whether a file exists at the given path
	_, err = os.Stat(path)
	if os.IsNotExist(err) {
		// file does not exist, serve index.html
		http.ServeFile(w, r, filepath.Join(h.staticPath, h.indexPath))
		return
	} else if err != nil {
		// if we got an error (that wasn't that the file doesn't exist) stating the
		// file, return a 500 internal server error and stop
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// otherwise, use http.FileServer to serve the static dir
	http.FileServer(http.Dir(h.staticPath)).ServeHTTP(w, r)
}
