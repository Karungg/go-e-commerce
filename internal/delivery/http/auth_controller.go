package http

import (
	"errors"
	"net/http"

	"go-e-commerce/internal/delivery/http/response"
	"go-e-commerce/internal/dto"
	"go-e-commerce/internal/pkg/apperror"
	"go-e-commerce/internal/port"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	authUsecase port.AuthUseCase
}

func NewAuthController(authUsecase port.AuthUseCase) *AuthController {
	return &AuthController{
		authUsecase: authUsecase,
	}
}

func (c *AuthController) RegisterCustomer(ctx *gin.Context) {
	var req dto.RegisterCustomerReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, http.StatusBadRequest, "Invalid request payload", apperror.FormatValidationError(err))
		return
	}

	token, err := c.authUsecase.RegisterCustomer(ctx.Request.Context(), &req)
	if err != nil {
		var appErr *apperror.AppError
		if errors.As(err, &appErr) {
			status := response.MapAppErrorToHTTPStatus(appErr)
			response.Error(ctx, status, "Registration failed", appErr.Message)
			return
		}
		
		response.Error(ctx, http.StatusInternalServerError, "Internal server error", err.Error())
		return
	}

	response.Success(ctx, http.StatusCreated, "customer registered successfully", gin.H{"token": token})
}

func (c *AuthController) RegisterSeller(ctx *gin.Context) {
	var req dto.RegisterSellerReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, http.StatusBadRequest, "Invalid request payload", apperror.FormatValidationError(err))
		return
	}

	token, err := c.authUsecase.RegisterSeller(ctx.Request.Context(), &req)
	if err != nil {
		var appErr *apperror.AppError
		if errors.As(err, &appErr) {
			status := response.MapAppErrorToHTTPStatus(appErr)
			response.Error(ctx, status, "Registration failed", appErr.Message)
			return
		}
		
		response.Error(ctx, http.StatusInternalServerError, "Internal server error", err.Error())
		return
	}

	response.Success(ctx, http.StatusCreated, "seller registered successfully", gin.H{"token": token})
}

func (c *AuthController) Login(ctx *gin.Context) {
	var req dto.LoginReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, http.StatusBadRequest, "Invalid request payload", apperror.FormatValidationError(err))
		return
	}

	token, err := c.authUsecase.Login(ctx.Request.Context(), &req)
	if err != nil {
		var appErr *apperror.AppError
		if errors.As(err, &appErr) {
			status := response.MapAppErrorToHTTPStatus(appErr)
			response.Error(ctx, status, "Login failed", appErr.Message)
			return
		}

		response.Error(ctx, http.StatusInternalServerError, "Internal server error", err.Error())
		return
	}

	response.Success(ctx, http.StatusOK, "login successful", gin.H{"token": token})
}