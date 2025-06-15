package userModel

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

type CreateUserRequest struct {
	Username string `validate:"required,min=8,max=16,alphanum"`
	Email    string `validate:"required,email"`
	AuthId   string `validate:"required"`
}

func (r CreateUserRequest) GetValidationError(err validator.FieldError) string {
	switch err.Field() {
	case "Username":
		return getUsernameError(err)
	case "Email":
		return getEmailError(err)
	case "AuthId":
		return getAuthIdError(err)
	default:
		return fmt.Sprintf("Invalid value for %s", err.Field())
	}
}
