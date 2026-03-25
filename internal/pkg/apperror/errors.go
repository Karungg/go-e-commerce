package apperror

// ErrorCode is a domain-specific string code for an error
type ErrorCode string

const (
	CodeBadRequest         ErrorCode = "BAD_REQUEST"
	CodeConflict           ErrorCode = "CONFLICT"
	CodeUnauthorized       ErrorCode = "UNAUTHORIZED"
	CodeNotFound           ErrorCode = "NOT_FOUND"
	CodeInternal           ErrorCode = "INTERNAL_ERROR"
)

// AppError represents a predefined domain error
type AppError struct {
	Code    ErrorCode `json:"code"`
	Message string    `json:"message"`
}

// Error implements the error interface
func (e *AppError) Error() string {
	return e.Message
}


var (
	// Auth Errors
	ErrEmailConflict     = &AppError{CodeConflict, "email is already strictly registered"}
	ErrPhoneConflict     = &AppError{CodeConflict, "phone number is already associated with an account"}
	ErrStoreNameConflict = &AppError{CodeConflict, "store name is already taken"}
	ErrInvalidPassword   = &AppError{CodeUnauthorized, "invalid email or password"}
	ErrUserNotFound      = &AppError{CodeNotFound, "user profile could not be found"}

	// Category Errors
	ErrCategoryNotFound  = &AppError{CodeNotFound, "category not found"}

	// System Errors
	ErrBadRequest = &AppError{CodeBadRequest, "invalid request"}
	ErrInternal   = &AppError{CodeInternal, "internal server error"}
)
