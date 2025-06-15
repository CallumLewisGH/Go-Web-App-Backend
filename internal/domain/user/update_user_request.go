package userModel

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

type UpdateUserRequest struct {
	Username       *string `validate:"omitempty,min=8,max=16,alphanum"`
	Email          *string `validate:"omitempty,email"`
	Timezone       *string `validate:"omitempty,timezone"`
	ProfilePicture *string `validate:"omitempty,base64"`
	Bio            *string `validate:"omitempty,max=500"`
}

func (r UpdateUserRequest) GetValidationError(err validator.FieldError) string {
	switch err.Field() {
	case "Username":
		return getUsernameError(err)
	case "Email":
		return getEmailError(err)
	case "Timezone":
		return getTimezoneError(err)
	case "Bio":
		return getBioError(err)
	case "ProfilePicture":
		return getProfilePictureError(err)
	default:
		return fmt.Sprintf("Invalid value for %s", err.Field())
	}
}
