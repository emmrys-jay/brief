package router

import (
	"brief/pkg/handler/user"
	mdw "brief/pkg/middleware"
	"brief/utility"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

// User registers user paths with router 'e'
func User(r chi.Router, validate *validator.Validate, ApiVersion string, logger *utility.Logger) chi.Router {

	userCtrl := user.Controller{Validate: validate, Logger: logger}

	// Free endpoints
	r.Group(func(r chi.Router) {
		r.Post("/users", userCtrl.Register)
		r.Post("/users/login", userCtrl.Login)
		r.Post("/users/forgot-password", userCtrl.ForgotPassword)
	})

	// User endpoints
	r.Group(func(r chi.Router) {
		r.Use(mdw.Me) // user middleware

		r.Get("/users", userCtrl.GetMe)
		r.Patch("/users", userCtrl.UpdateMe)
		r.Patch("/users/verify", userCtrl.VerifyMe)
		r.Patch("/users/reset-password", userCtrl.ResetPassword)
	})

	// Admin endpoints
	r.Group(func(r chi.Router) {
		r.Use(mdw.Admin) // admin middleware

		r.Get("/users/get-all", userCtrl.GetAll)
		r.Get("/users/{idOrEmail}", userCtrl.GetUserByIdOrEmail)
		r.Patch("/users/lock/{idOrEmail}", userCtrl.LockUser)
		r.Patch("/users/lock/{idOrEmail}", userCtrl.UnlockUser)
	})

	return r
}
