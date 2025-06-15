package userModel

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

func getUsernameError(err validator.FieldError) string {
	switch err.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", err.Field())
	case "min":
		return fmt.Sprintf("%s must be at least 8 characters long", err.Field())
	case "max":
		return fmt.Sprintf("%s must not exceed 16 characters", err.Field())
	case "alphanum":
		return fmt.Sprintf("%s can only contain letters and numbers", err.Field())
	default:
		return fmt.Sprintf("Invalid %s format", err.Field())
	}
}

func getAuthIdError(err validator.FieldError) string {
	switch err.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", err.Field())
	default:
		return fmt.Sprintf("Invalid %s format", err.Field())
	}
}

func getEmailError(err validator.FieldError) string {
	switch err.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", err.Field())
	default:
		return fmt.Sprintf("Invalid %s format", err.Field())
	}
}

func getTimezoneError(err validator.FieldError) string {
	switch err.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", err.Field())
	case "timezone":
		return fmt.Sprintf("%s must be a valid IANA timezone", err.Field())
	default:
		return fmt.Sprintf("Invalid %s format", err.Field())
	}
}

func getBioError(err validator.FieldError) string {
	switch err.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", err.Field())
	case "max":
		return fmt.Sprintf("%s must not exceed 500 characters", err.Field())
	default:
		return fmt.Sprintf("Invalid %s format", err.Field())
	}
}

func getProfilePictureError(err validator.FieldError) string {
	switch err.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", err.Field())
	case "base64":
		return fmt.Sprintf("%s must be a valid base64 encoded string", err.Field())
	default:
		return fmt.Sprintf("Invalid %s format", err.Field())
	}
}
