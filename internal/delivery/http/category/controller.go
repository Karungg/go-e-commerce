package category

import (
	"errors"
	"net/http"

	"go-e-commerce/internal/delivery/http/response"
	categoryDTO "go-e-commerce/internal/dto/category"
	"go-e-commerce/internal/pkg/apperror"
	categoryPort "go-e-commerce/internal/port/category"

	"github.com/gin-gonic/gin"
)

type CategoryController struct {
	categoryUseCase categoryPort.CategoryUseCase
}

func NewCategoryController(categoryUseCase categoryPort.CategoryUseCase) *CategoryController {
	return &CategoryController{
		categoryUseCase: categoryUseCase,
	}
}

func (c *CategoryController) Create(ctx *gin.Context) {
	var req categoryDTO.CreateCategoryReq
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
	var req categoryDTO.UpdateCategoryReq

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
