package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go-e-commerce/internal/pkg/apperror"
)

// PaginationMeta contains pagination details
type PaginationMeta struct {
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	TotalItems int `json:"total_items"`
	TotalPages int `json:"total_pages"`
}

// WebResponse represents a standardized API payload structure widely recognized by frontend clients
type WebResponse struct {
	Code    int             `json:"code"`
	Status  string          `json:"status"`   // "success", "error"
	Message string          `json:"message"`
	Data    interface{}     `json:"data,omitempty"`
	Meta    *PaginationMeta `json:"meta,omitempty"`
	Errors  interface{}     `json:"errors,omitempty"`
}

// Success serves a canonical HTTP 20X payload with "status": "success"
func Success(ctx *gin.Context, code int, message string, data interface{}) {
	ctx.JSON(code, WebResponse{
		Code:    code,
		Status:  "success",
		Message: message,
		Data:    data,
	})
}

// SuccessWithMeta serves a canonical HTTP 20X payload with "status": "success" and pagination metadata
func SuccessWithMeta(ctx *gin.Context, code int, message string, data interface{}, meta *PaginationMeta) {
	ctx.JSON(code, WebResponse{
		Code:    code,
		Status:  "success",
		Message: message,
		Data:    data,
		Meta:    meta,
	})
}

// Error serves a canonical HTTP 40X or 50X payload with "status": "error"
func Error(ctx *gin.Context, code int, message string, errors interface{}) {
	ctx.JSON(code, WebResponse{
		Code:    code,
		Status:  "error",
		Message: message,
		Errors:  errors,
	})
}

// MapAppErrorToHTTPStatus converts an app domain error code to an HTTP status code
func MapAppErrorToHTTPStatus(appErr *apperror.AppError) int {
	switch appErr.Code {
	case apperror.CodeConflict:
		return http.StatusConflict
	case apperror.CodeUnauthorized:
		return http.StatusUnauthorized
	case apperror.CodeNotFound:
		return http.StatusNotFound
	default:
		return http.StatusInternalServerError
	}
}
