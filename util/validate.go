package util

import (
	"github.com/go-playground/validator/v10"
)

func Validate(err error) []string {
	var errors []string
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, fieldErr := range validationErrors {
			errors = append(errors, fieldErr.Field()+" is required.")
		}
	} else {
		errors = append(errors, err.Error())
	}
	return errors
}