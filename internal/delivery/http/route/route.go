package route

import (
	deliveryHttp "go-e-commerce/internal/delivery/http"
	"go-e-commerce/internal/port"

	"github.com/gin-gonic/gin"
)

// SetupRoutes centralizes all API route registrations
func SetupRoutes(
	api *gin.RouterGroup,
	authController *deliveryHttp.AuthController,
	jwtAuth port.TokenValidator,
) {
	// Public Routes
	auth := api.Group("/auth")
	{
		auth.POST("/register/customer", authController.RegisterCustomer)
		auth.POST("/register/seller", authController.RegisterSeller)
	}

}
