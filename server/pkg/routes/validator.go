package routes

import (
	"gopkg.in/go-playground/validator.v9"
	"strings"
)

func newValidator() *Validator {
	validate := validator.New()

	// Passwords must have 8 or more characters
	validate.RegisterValidation("isValidPassword", func(fl validator.FieldLevel) bool {
		password := strings.TrimSpace(fl.Field().String())

		if len(password) < 8  {
			return false
		}

		return true
	})

	return &Validator{
		Validator: validate,
	}
}

type Validator struct {
	Validator *validator.Validate
}

func (cv *Validator) Validate(i interface{}) error {
	return cv.Validator.Struct(i)
}