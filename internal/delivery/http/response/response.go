package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go-e-commerce/internal/pkg/apperror"
)

type PaginationMeta struct {
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	TotalItems int `json:"total_items"`
	TotalPages int `json:"total_pages"`
}

type WebResponse struct {
	Code    int             `json:"code"`
	Status  string          `json:"status"` 
	Message string          `json:"message"`
	Data    interface{}     `json:"data,omitempty"`
	Meta    *PaginationMeta `json:"meta,omitempty"`
	Errors  interface{}     `json:"errors,omitempty"`
}

func Success(ctx *gin.Context, code int, message string, data interface{}) {
	ctx.JSON(code, WebResponse{
		Code:    code,
		Status:  "success",
		Message: message,
		Data:    data,
	})
}

func SuccessWithMeta(ctx *gin.Context, code int, message string, data interface{}, meta *PaginationMeta) {
	ctx.JSON(code, WebResponse{
		Code:    code,
		Status:  "success",
		Message: message,
		Data:    data,
		Meta:    meta,
	})
}

func Error(ctx *gin.Context, code int, message string, errors interface{}) {
	ctx.JSON(code, WebResponse{
		Code:    code,
		Status:  "error",
		Message: message,
		Errors:  errors,
	})
}

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
