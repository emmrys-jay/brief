package router

import (
	"brief/pkg/handler/url"
	mdw "brief/pkg/middleware"
	"brief/utility"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

// Url registers url paths with router 'r'
func Url(r chi.Router, validate *validator.Validate, logger *utility.Logger) chi.Router {

	urlCtrl := url.Controller{Validate: validate, Logger: logger}

	// Free Endpoint
	r.Group(func(r chi.Router) {
		r.Get("/{hash}", urlCtrl.Redirect)
	})

	// Shorten endpoint
	r.Group(func(r chi.Router) {
		r.Use(mdw.Shorten)
		r.Post("/url/shorten", urlCtrl.Shorten)
	})

	// User endpoints
	r.Group(func(r chi.Router) {
		r.Use(mdw.Me) // user middleware

		r.Get("/url", urlCtrl.GetUrls)
		r.Delete("/url", urlCtrl.Delete)
	})

	// Admin endpoints
	r.Group(func(r chi.Router) {
		r.Use(mdw.Admin) // admin middleware

		r.Get("/url/get-all", urlCtrl.GetAll)
		r.Get("/url/{user-id}", urlCtrl.GetUrlsByUserID)
		r.Delete("/url/{id}", urlCtrl.DeleteUrlByID)
	})

	return r
}
