package router

import (
	"brief/pkg/handler/health"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	log "github.com/sirupsen/logrus"
)

func Health(r chi.Router, validate *validator.Validate, logger *log.Logger) chi.Router {

	health := health.Controller{Validate: validate, Logger: logger}

	r.Post("/health", health.Post)
	r.Get("/health", health.Get)

	return r
}
