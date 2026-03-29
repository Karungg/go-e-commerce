package apperror

type ErrorCode string

const (
	CodeBadRequest         ErrorCode = "BAD_REQUEST"
	CodeConflict           ErrorCode = "CONFLICT"
	CodeUnauthorized       ErrorCode = "UNAUTHORIZED"
	CodeNotFound           ErrorCode = "NOT_FOUND"
	CodeInternal           ErrorCode = "INTERNAL_ERROR"
)

type AppError struct {
	Code    ErrorCode `json:"code"`
	Message string    `json:"message"`
}

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

	// Product Errors
	ErrProductNotFound   = &AppError{CodeNotFound, "product not found"}

	// System Errors
	ErrBadRequest = &AppError{CodeBadRequest, "invalid request"}
	ErrInternal   = &AppError{CodeInternal, "internal server error"}
)
