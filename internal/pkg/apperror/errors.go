package apperror

import "net/http"

// AppError represents a predefined domain error combining HTTP Codes and user-facing Messages natively
type AppError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// Error strictly implements the Go error interface natively
func (e *AppError) Error() string {
	return e.Message
}

// ---------------------------------------------------------
// CENTRAL ERROR DICTIONARY
// ---------------------------------------------------------

var (
	// Auth Errors
	ErrEmailConflict   = &AppError{http.StatusConflict, "email is already strictly registered"}
	ErrInvalidPassword = &AppError{http.StatusUnauthorized, "invalid email or password"}
	ErrUserNotFound    = &AppError{http.StatusNotFound, "user profile could not be found"}

	// System Errors
	ErrInternal = &AppError{http.StatusInternalServerError, "internal server error"}
)
