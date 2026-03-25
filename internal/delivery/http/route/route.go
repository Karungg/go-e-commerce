package route

import (
	authCtrl "go-e-commerce/internal/delivery/http/auth"
	categoryCtrl "go-e-commerce/internal/delivery/http/category"
	"go-e-commerce/internal/delivery/http/middleware"
	authPort "go-e-commerce/internal/port/auth"

	"github.com/gin-gonic/gin"
)

// SetupRoutes centralizes all API route registrations
func SetupRoutes(
	api *gin.RouterGroup,
	authController *authCtrl.AuthController,
	categoryController *categoryCtrl.CategoryController,
	jwtAuth authPort.TokenValidator,
) {
	// Public Routes
	auth := api.Group("/auth")
	{
		auth.POST("/register/customer", authController.RegisterCustomer)
		auth.POST("/register/seller", authController.RegisterSeller)
		auth.POST("/login", authController.Login)
	}

	categoriesPublic := api.Group("/categories")
	{
		categoriesPublic.GET("", categoryController.GetAll)
		categoriesPublic.GET("/:id", categoryController.GetByID)
	}

	authProtected := api.Group("/auth")
	authProtected.Use(middleware.RequireAuth(jwtAuth))
	{
		authProtected.POST("/logout", authController.Logout)
	}

	categoriesProtected := api.Group("/categories")
	categoriesProtected.Use(middleware.RequireAuth(jwtAuth))
	{
		categoriesProtected.POST("", categoryController.Create)
		categoriesProtected.PUT("/:id", categoryController.Update)
		categoriesProtected.DELETE("/:id", categoryController.Delete)
	}
}
