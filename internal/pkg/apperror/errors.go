package apperror

// ErrorCode is a domain-specific string code for an error
type ErrorCode string

const (
	CodeConflict           ErrorCode = "CONFLICT"
	CodeUnauthorized       ErrorCode = "UNAUTHORIZED"
	CodeNotFound           ErrorCode = "NOT_FOUND"
	CodeInternal           ErrorCode = "INTERNAL_ERROR"
)

// AppError represents a predefined domain error combining ErrorCodes and user-facing Messages natively
type AppError struct {
	Code    ErrorCode `json:"code"`
	Message string    `json:"message"`
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
	ErrEmailConflict     = &AppError{CodeConflict, "email is already strictly registered"}
	ErrPhoneConflict     = &AppError{CodeConflict, "phone number is already associated with an account"}
	ErrStoreNameConflict = &AppError{CodeConflict, "store name is already taken"}
	ErrInvalidPassword   = &AppError{CodeUnauthorized, "invalid email or password"}
	ErrUserNotFound      = &AppError{CodeNotFound, "user profile could not be found"}

	// System Errors
	ErrInternal = &AppError{CodeInternal, "internal server error"}
)
