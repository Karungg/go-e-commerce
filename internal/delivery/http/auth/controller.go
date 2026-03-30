package auth

import (
	"errors"
	"net/http"

	authDTO "go-e-commerce/internal/dto/auth"
	"go-e-commerce/internal/delivery/http/response"
	"go-e-commerce/internal/pkg/apperror"
	authPort "go-e-commerce/internal/port/auth"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AuthController struct {
	authUsecase authPort.AuthUseCase
}

func NewAuthController(authUsecase authPort.AuthUseCase) *AuthController {
	return &AuthController{
		authUsecase: authUsecase,
	}
}

func (c *AuthController) RegisterCustomer(ctx *gin.Context) {
	var req authDTO.RegisterCustomerReq
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
	var req authDTO.RegisterSellerReq
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
	var req authDTO.LoginReq
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

func (c *AuthController) Logout(ctx *gin.Context) {
	err := c.authUsecase.Logout(ctx.Request.Context())
	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, "Internal server error", err.Error())
		return
	}

	response.Success(ctx, http.StatusOK, "logout successful", nil)
}

func (c *AuthController) UpdateCustomer(ctx *gin.Context) {
	var req authDTO.UpdateCustomerReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, http.StatusBadRequest, "Invalid request payload", apperror.FormatValidationError(err))
		return
	}

	userIDRaw, exists := ctx.Get("userID")
	if !exists {
		response.Error(ctx, http.StatusUnauthorized, "Authorization failed", "User ID not found in context")
		return
	}

	userID, ok := userIDRaw.(string)
	if !ok {
		response.Error(ctx, http.StatusUnauthorized, "Authorization failed", "Invalid user ID format in context")
		return
	}

	parsedUserID, err := uuid.Parse(userID)
	if err != nil {
		response.Error(ctx, http.StatusUnauthorized, "Authorization failed", "Invalid user ID string")
		return
	}

	err = c.authUsecase.UpdateCustomer(ctx.Request.Context(), parsedUserID, &req)
	if err != nil {
		var appErr *apperror.AppError
		if errors.As(err, &appErr) {
			status := response.MapAppErrorToHTTPStatus(appErr)
			response.Error(ctx, status, "Update failed", appErr.Message)
			return
		}

		response.Error(ctx, http.StatusInternalServerError, "Internal server error", err.Error())
		return
	}

	response.Success(ctx, http.StatusOK, "customer profile updated successfully", nil)
}
