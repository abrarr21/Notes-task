package utils

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

func Validate(input any) map[string]string {
	err := validate.Struct(input)
	if err == nil {
		return nil
	}

	errors := map[string]string{}

	for _, e := range err.(validator.ValidationErrors) {
		field := strings.ToLower(e.Field())
		errors[field] = formatError(e)
	}

	return errors
}

func formatError(e validator.FieldError) string {
	switch e.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", strings.ToLower(e.Field()))
	case "min":
		if e.Param() == "1" {
			return fmt.Sprintf("%s cannot be empty", strings.ToLower(e.Field()))
		}
		return fmt.Sprintf("%s must be at least %s characters", strings.ToLower(e.Field()), e.Param())
	case "max":
		return fmt.Sprintf("%s must be at most %s characters", strings.ToLower(e.Field()), e.Param())
	case "email":
		return fmt.Sprintf("%s must be a valid email address", strings.ToLower(e.Field()))
	default:
		return fmt.Sprintf("%s is invalid", strings.ToLower(e.Field()))
	}
}
