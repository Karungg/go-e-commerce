package product

import (
	"errors"
	"math"
	"net/http"
	"strconv"

	"go-e-commerce/internal/delivery/http/response"
	productDTO "go-e-commerce/internal/dto/product"
	"go-e-commerce/internal/pkg/apperror"
	productPort "go-e-commerce/internal/port/product"

	"github.com/gin-gonic/gin"
)

type ProductController struct {
	productUseCase productPort.ProductUseCase
}

func NewProductController(productUseCase productPort.ProductUseCase) *ProductController {
	return &ProductController{
		productUseCase: productUseCase,
	}
}

func (c *ProductController) Create(ctx *gin.Context) {
	var req productDTO.CreateProductReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, http.StatusBadRequest, "Invalid request payload", apperror.FormatValidationError(err))
		return
	}

	res, err := c.productUseCase.CreateProduct(ctx.Request.Context(), &req)
	if err != nil {
		c.handleError(ctx, err)
		return
	}

	response.Success(ctx, http.StatusCreated, "Product created successfully", res)
}

func (c *ProductController) GetAll(ctx *gin.Context) {
	page, err := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(ctx.DefaultQuery("limit", "10"))
	if err != nil || limit < 1 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	res, total, err := c.productUseCase.GetAllProducts(ctx.Request.Context(), page, limit)
	if err != nil {
		c.handleError(ctx, err)
		return
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))
	meta := &response.PaginationMeta{
		Page:       page,
		Limit:      limit,
		TotalItems: int(total),
		TotalPages: totalPages,
	}

	response.SuccessWithMeta(ctx, http.StatusOK, "Products fetched successfully", res, meta)
}

func (c *ProductController) GetByID(ctx *gin.Context) {
	id := ctx.Param("id")

	res, err := c.productUseCase.GetProductByID(ctx.Request.Context(), id)
	if err != nil {
		c.handleError(ctx, err)
		return
	}

	response.Success(ctx, http.StatusOK, "Product fetched successfully", res)
}

func (c *ProductController) Update(ctx *gin.Context) {
	id := ctx.Param("id")
	var req productDTO.UpdateProductReq

	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, http.StatusBadRequest, "Invalid request payload", apperror.FormatValidationError(err))
		return
	}

	res, err := c.productUseCase.UpdateProduct(ctx.Request.Context(), id, &req)
	if err != nil {
		c.handleError(ctx, err)
		return
	}

	response.Success(ctx, http.StatusOK, "Product updated successfully", res)
}

func (c *ProductController) Delete(ctx *gin.Context) {
	id := ctx.Param("id")

	err := c.productUseCase.DeleteProduct(ctx.Request.Context(), id)
	if err != nil {
		c.handleError(ctx, err)
		return
	}

	response.Success(ctx, http.StatusOK, "Product deleted successfully", nil)
}

func (c *ProductController) handleError(ctx *gin.Context, err error) {
	var appErr *apperror.AppError
	if errors.As(err, &appErr) {
		status := response.MapAppErrorToHTTPStatus(appErr)
		response.Error(ctx, status, "Request failed", appErr.Message)
		return
	}

	response.Error(ctx, http.StatusInternalServerError, "Internal server error", err.Error())
}
