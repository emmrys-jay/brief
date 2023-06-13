package router

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"brief/utility"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-playground/validator/v10"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

func Setup(validate *validator.Validate, logger *utility.Logger) chi.Router {
	r := chi.NewRouter()

	// Middlewares
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Compress(5))
	r.Use(middleware.Timeout(60 * time.Second))

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "PATCH"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Token", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	ApiVersion := "v1"
	r.Route(fmt.Sprintf("/api/%s", ApiVersion), func(r chi.Router) {
		Health(r, validate, ApiVersion, logger)
		User(r, validate, ApiVersion, logger)
	})

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		res := map[string]interface{}{
			"name":    "Not Found",
			"message": "Page not found.",
			"code":    404,
			"status":  http.StatusNotFound,
		}
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(res)
	})

	return r
}
