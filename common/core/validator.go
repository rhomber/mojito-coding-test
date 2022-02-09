package core

import (
	"gopkg.in/go-playground/validator.v9"
)

func newValidator() *validator.Validate {
	validate := validator.New()

	// Add custom validators.

	return validate
}
