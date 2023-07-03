package router

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"brief/internal/config"
	"brief/pkg/handler/url"

	log "github.com/sirupsen/logrus"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-playground/validator/v10"
	httpSwagger "github.com/swaggo/http-swagger/v2"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"

	_ "brief/docs"
)

func Setup(validate *validator.Validate, logger *log.Logger) chi.Router {
	r := chi.NewRouter()

	// Middlewares
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Compress(5))
	r.Use(middleware.Timeout(60 * time.Second))
	r.Use(middleware.SetHeader("Content-Type", "application/json"))

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "PATCH"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Token", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	ApiVersion := "v1"

	// Redirect Endpoint
	urlCtrl := url.Controller{Validate: validate, Logger: logger}
	r.Group(func(r chi.Router) {
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Server is running"))
		})
		r.Get("/{hash}", urlCtrl.Redirect)
	})

	// Endpoints starting with "/api/v1"
	r.Route(fmt.Sprintf("/api/%s", ApiVersion), func(r chi.Router) {
		Health(r, validate, logger)
		User(r, validate, logger)
		Url(r, validate, logger)
	})

	// Swagger endpoint
	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://0.0.0.0:"+config.GetConfig().ServerPort+"/swagger/doc.json"), //The url pointing to API definition
	))

	// Not found
	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		res := map[string]interface{}{
			"name":    "Not Found",
			"message": "Page not found.",
			"code":    404,
			"status":  http.StatusNotFound,
		}
		w.WriteHeader(http.StatusNotFound)
		resV, _ := json.Marshal(res)
		w.Write(resV)
	})

	return r
}
