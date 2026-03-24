package route

import (
	deliveryHttp "go-e-commerce/internal/delivery/http"
	"go-e-commerce/internal/delivery/http/middleware"
	"go-e-commerce/internal/security"

	"github.com/gin-gonic/gin"
)

// SetupRoutes centralizes all API route registrations
func SetupRoutes(
	api *gin.RouterGroup,
	authController *deliveryHttp.AuthController,
	jwtAuth *security.JWTAuth,
) {
	// Public Routes
	auth := api.Group("/auth")
	{
		auth.POST("/register/customer", authController.RegisterCustomer)
		auth.POST("/register/seller", authController.RegisterSeller)
	}

	// Example of Protected Routes (to be populated later)
	// protected := api.Group("/protected")
	// protected.Use(middleware.RequireAuth(jwtAuth))
	// {
	//     // Example Role requirement
	//     // adminRoutes := protected.Group("/admin")
	//     // adminRoutes.Use(middleware.RequireRole(string(entity.RoleAdmin)))
	//	   // { ... }
	// }
	_ = middleware.RequireAuth // Just to ensure import is kept if unused initially
}
