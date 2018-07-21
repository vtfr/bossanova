package model

import (
	"regexp"

	"gopkg.in/go-playground/validator.v9"
)

var validate *validator.Validate

func init() {
	boardRx := regexp.MustCompile("^[a-zA-Z][a-zA-Z0-9]{0,15}$")

	validate = validator.New()
	validate.RegisterValidation("board", func(fl validator.FieldLevel) bool {
		return boardRx.MatchString(fl.Field().String())
	})
}

// Validate validates a struct based on a validator tags
func Validate(v interface{}) error {
	return validate.Struct(v)
}
