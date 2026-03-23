package http

import (
	"errors"
	"net/http"

	"go-e-commerce/internal/delivery/http/response"
	"go-e-commerce/internal/pkg/apperror"
	"go-e-commerce/internal/usecase"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	authUsecase usecase.AuthUseCase
}

func NewAuthController(router *gin.RouterGroup, authUsecase usecase.AuthUseCase) {
	controller := &AuthController{
		authUsecase: authUsecase,
	}

	authRoutes := router.Group("/auth")
	{
		authRoutes.POST("/register/customer", controller.RegisterCustomer)
		authRoutes.POST("/register/seller", controller.RegisterSeller)
	}
}

func (c *AuthController) RegisterCustomer(ctx *gin.Context) {
	var req usecase.RegisterCustomerReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, http.StatusBadRequest, "Invalid request payload", err.Error())
		return
	}

	token, err := c.authUsecase.RegisterCustomer(ctx.Request.Context(), &req)
	if err != nil {
		// Strictly intercept pre-defined Domain Errors
		var appErr *apperror.AppError
		if errors.As(err, &appErr) {
			response.Error(ctx, appErr.Code, "Registration failed", appErr.Message)
			return
		}
		
		response.Error(ctx, http.StatusInternalServerError, "Internal server error", err.Error())
		return
	}

	response.Success(ctx, http.StatusCreated, "customer registered successfully", gin.H{"token": token})
}

func (c *AuthController) RegisterSeller(ctx *gin.Context) {
	var req usecase.RegisterSellerReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, http.StatusBadRequest, "Invalid request payload", err.Error())
		return
	}

	token, err := c.authUsecase.RegisterSeller(ctx.Request.Context(), &req)
	if err != nil {
		var appErr *apperror.AppError
		if errors.As(err, &appErr) {
			response.Error(ctx, appErr.Code, "Registration failed", appErr.Message)
			return
		}
		
		response.Error(ctx, http.StatusInternalServerError, "Internal server error", err.Error())
		return
	}

	response.Success(ctx, http.StatusCreated, "seller registered successfully", gin.H{"token": token})
}
