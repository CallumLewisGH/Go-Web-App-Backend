package validation

import (
	"errors"

	"github.com/go-playground/validator/v10"
)

type ValidationMessenger interface {
	GetValidationError(err validator.FieldError) string
}

var (
	validate *validator.Validate = validator.New()
)

// Validate performs validation and returns custom error messages
func ValidateBody[T ValidationMessenger](input T) []string {
	var validationErrors []string

	err := validate.Struct(input)
	if err == nil {
		return nil
	}

	validatorErrs := err.(validator.ValidationErrors)
	for _, err := range validatorErrs {
		outputErrorMessage := errors.New(input.GetValidationError(err))
		validationErrors = append(validationErrors, outputErrorMessage.Error())
	}

	return validationErrors
}
