package route

import (
	authCtrl "go-e-commerce/internal/delivery/http/auth"
	categoryCtrl "go-e-commerce/internal/delivery/http/category"
	"go-e-commerce/internal/delivery/http/middleware"
	productCtrl "go-e-commerce/internal/delivery/http/product"
	authPort "go-e-commerce/internal/port/auth"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(
	api *gin.RouterGroup,
	authController *authCtrl.AuthController,
	categoryController *categoryCtrl.CategoryController,
	productController *productCtrl.ProductController,
	jwtAuth authPort.TokenValidator,
) {
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

	productsPublic := api.Group("/products")
	{
		productsPublic.GET("", productController.GetAll)
		productsPublic.GET("/:id", productController.GetByID)
	}

	authProtected := api.Group("/auth")
	authProtected.Use(middleware.RequireAuth(jwtAuth))
	{
		authProtected.POST("/logout", authController.Logout)
		authProtected.PUT("/customer", authController.UpdateCustomer)
	}

	categoriesProtected := api.Group("/categories")
	categoriesProtected.Use(middleware.RequireAuth(jwtAuth))
	{
		categoriesProtected.POST("", categoryController.Create)
		categoriesProtected.PUT("/:id", categoryController.Update)
		categoriesProtected.DELETE("/:id", categoryController.Delete)
	}

	productsProtected := api.Group("/products")
	productsProtected.Use(middleware.RequireAuth(jwtAuth))
	{
		productsProtected.POST("", productController.Create)
		productsProtected.PUT("/:id", productController.Update)
		productsProtected.DELETE("/:id", productController.Delete)
	}
}
