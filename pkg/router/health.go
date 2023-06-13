package router

import (
	"brief/pkg/handler/health"

	"brief/utility"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

func Health(r chi.Router, validate *validator.Validate, ApiVersion string, logger *utility.Logger) chi.Router {

	health := health.Controller{Validate: validate, Logger: logger}

	r.Post("/health", health.Post)
	r.Get("/health", health.Get)

	return r
}
