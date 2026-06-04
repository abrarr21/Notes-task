package utils

import "github.com/go-playground/validator/v10"

var validate = validator.New()

func Validate(input any) error {
	return validate.Struct(input)
}
