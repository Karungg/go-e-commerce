package http

import (
	"errors"
	"go-e-commerce/internal/usecase"
	"net/http"

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
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := c.authUsecase.RegisterCustomer(ctx.Request.Context(), &req)
	if err != nil {
		if errors.Is(err, usecase.ErrEmailConflict) {
			ctx.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "customer registered successfully",
		"token":   token,
	})
}

func (c *AuthController) RegisterSeller(ctx *gin.Context) {
	var req usecase.RegisterSellerReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := c.authUsecase.RegisterSeller(ctx.Request.Context(), &req)
	if err != nil {
		if errors.Is(err, usecase.ErrEmailConflict) {
			ctx.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "seller registered successfully",
		"token":   token,
	})
}
