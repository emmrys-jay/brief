package user

import (
	"brief/utility"

	"github.com/go-playground/validator/v10"
)

type Controller struct {
	Validate *validator.Validate
	Logger   *utility.Logger
}
