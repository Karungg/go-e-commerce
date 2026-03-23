package apperror

import (
	"errors"
	"fmt"

	"github.com/go-playground/validator/v10"
)

type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// FormatValidationError elegantly iterates dynamically across complex Gin payload arrays mapping frontend variables correctly.
func FormatValidationError(err error) interface{} {
	var ve validator.ValidationErrors
	if errors.As(err, &ve) {
		out := make([]FieldError, len(ve))
		for i, fe := range ve {
			out[i] = FieldError{
				Field:   fe.Field(),
				Message: getErrorMsg(fe),
			}
		}
		return out
	}
	// Fallback natively to Raw Error strings globally avoiding panics
	return err.Error()
}

func getErrorMsg(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "This field is required"
	case "email":
		return "Must be a valid email address format"
	case "min":
		return fmt.Sprintf("Must be strictly at least %s characters long", fe.Param())
	case "max":
		return fmt.Sprintf("Must be exactly at most %s characters long", fe.Param())
	case "alpha":
		return "Must contain purely strict alphabetic natively characters"
	case "numeric":
		return "Must uniquely represent purely numeric strings"
	case "url":
		return "Must inherently emulate standard URL patterns"
	case "e164":
		return "Must securely match strict e164 (+...) format standards"
	}
	return fmt.Sprintf("Data structurally failed validation on '%s' tag globally", fe.Tag())
}
