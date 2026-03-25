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

type CategoryController struct {
	categoryUseCase port.CategoryUseCase
}

func NewCategoryController(categoryUseCase port.CategoryUseCase) *CategoryController {
	return &CategoryController{
		categoryUseCase: categoryUseCase,
	}
}

func (c *CategoryController) Create(ctx *gin.Context) {
	var req dto.CreateCategoryReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, http.StatusBadRequest, "Invalid request payload", apperror.FormatValidationError(err))
		return
	}

	res, err := c.categoryUseCase.CreateCategory(ctx.Request.Context(), &req)
	if err != nil {
		c.handleError(ctx, err)
		return
	}

	response.Success(ctx, http.StatusCreated, "Category created successfully", res)
}

func (c *CategoryController) GetAll(ctx *gin.Context) {
	res, err := c.categoryUseCase.GetAllCategories(ctx.Request.Context())
	if err != nil {
		c.handleError(ctx, err)
		return
	}

	response.Success(ctx, http.StatusOK, "Categories fetched successfully", res)
}

func (c *CategoryController) GetByID(ctx *gin.Context) {
	id := ctx.Param("id")

	res, err := c.categoryUseCase.GetCategoryByID(ctx.Request.Context(), id)
	if err != nil {
		c.handleError(ctx, err)
		return
	}

	response.Success(ctx, http.StatusOK, "Category fetched successfully", res)
}

func (c *CategoryController) Update(ctx *gin.Context) {
	id := ctx.Param("id")
	var req dto.UpdateCategoryReq

	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, http.StatusBadRequest, "Invalid request payload", apperror.FormatValidationError(err))
		return
	}

	res, err := c.categoryUseCase.UpdateCategory(ctx.Request.Context(), id, &req)
	if err != nil {
		c.handleError(ctx, err)
		return
	}

	response.Success(ctx, http.StatusOK, "Category updated successfully", res)
}

func (c *CategoryController) Delete(ctx *gin.Context) {
	id := ctx.Param("id")

	err := c.categoryUseCase.DeleteCategory(ctx.Request.Context(), id)
	if err != nil {
		c.handleError(ctx, err)
		return
	}

	response.Success(ctx, http.StatusOK, "Category deleted successfully", nil)
}

func (c *CategoryController) handleError(ctx *gin.Context, err error) {
	var appErr *apperror.AppError
	if errors.As(err, &appErr) {
		status := response.MapAppErrorToHTTPStatus(appErr)
		response.Error(ctx, status, "Request failed", appErr.Message)
		return
	}

	response.Error(ctx, http.StatusInternalServerError, "Internal server error", err.Error())
}
