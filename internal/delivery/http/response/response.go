package response

import "github.com/gin-gonic/gin"

// WebResponse represents a standardized API payload structure widely recognized by frontend clients
type WebResponse struct {
	Code    int         `json:"code"`
	Status  string      `json:"status"`   // "success", "error"
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Errors  interface{} `json:"errors,omitempty"`
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

// Error serves a canonical HTTP 40X or 50X payload with "status": "error"
func Error(ctx *gin.Context, code int, message string, errors interface{}) {
	ctx.JSON(code, WebResponse{
		Code:    code,
		Status:  "error",
		Message: message,
		Errors:  errors,
	})
}
