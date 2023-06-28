package url

import (
	"github.com/go-playground/validator/v10"
	log "github.com/sirupsen/logrus"
)

type Controller struct {
	Validate *validator.Validate
	Logger   *log.Logger
}
