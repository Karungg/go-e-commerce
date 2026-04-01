package cart

import (
	"net/http"

	"go-e-commerce/internal/delivery/http/response"
	"go-e-commerce/internal/dto"
	cartPort "go-e-commerce/internal/port/cart"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CartController struct {
	useCase cartPort.CartUseCase
}

func NewCartController(useCase cartPort.CartUseCase) *CartController {
	return &CartController{
		useCase: useCase,
	}
}

func (c *CartController) GetCart(ctx *gin.Context) {
	strUserID, exists := ctx.Get("userID")
	if !exists {
		response.Error(ctx, http.StatusUnauthorized, "Unauthorized", "User ID not found in context")
		return
	}

	userID, err := uuid.Parse(strUserID.(string))
	if err != nil {
		response.Error(ctx, http.StatusUnauthorized, "Unauthorized", "Invalid User ID format")
		return
	}

	res, err := c.useCase.GetCart(ctx.Request.Context(), userID)
	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, "Failed to get cart", err.Error())
		return
	}

	response.Success(ctx, http.StatusOK, "Cart retrieved successfully", res)
}

func (c *CartController) AddToCart(ctx *gin.Context) {
	strUserID, exists := ctx.Get("userID")
	if !exists {
		response.Error(ctx, http.StatusUnauthorized, "Unauthorized", "User ID not found in context")
		return
	}

	userID, err := uuid.Parse(strUserID.(string))
	if err != nil {
		response.Error(ctx, http.StatusUnauthorized, "Unauthorized", "Invalid User ID format")
		return
	}

	var req dto.AddCartItemRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, http.StatusBadRequest, "Invalid request", err.Error())
		return
	}

	if err := c.useCase.AddToCart(ctx.Request.Context(), userID, &req); err != nil {
		if err.Error() == "insufficient stock" || err.Error() == "product not found" {
			response.Error(ctx, http.StatusBadRequest, "Failed to add to cart", err.Error())
			return
		}
		response.Error(ctx, http.StatusInternalServerError, "Failed to add to cart", err.Error())
		return
	}

	response.Success(ctx, http.StatusOK, "Item added to cart successfully", nil)
}

func (c *CartController) UpdateCartItem(ctx *gin.Context) {
	strUserID, exists := ctx.Get("userID")
	if !exists {
		response.Error(ctx, http.StatusUnauthorized, "Unauthorized", "User ID not found in context")
		return
	}

	userID, err := uuid.Parse(strUserID.(string))
	if err != nil {
		response.Error(ctx, http.StatusUnauthorized, "Unauthorized", "Invalid User ID format")
		return
	}

	itemIDStr := ctx.Param("id")
	itemID, err := uuid.Parse(itemIDStr)
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, "Invalid request", "Invalid item ID format")
		return
	}

	var req dto.UpdateCartItemRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, http.StatusBadRequest, "Invalid request", err.Error())
		return
	}

	if err := c.useCase.UpdateCartItem(ctx.Request.Context(), userID, itemID, &req); err != nil {
		response.Error(ctx, http.StatusBadRequest, "Failed to update item", err.Error())
		return
	}

	response.Success(ctx, http.StatusOK, "Item quantity updated successfully", nil)
}

func (c *CartController) BatchDeleteCartItems(ctx *gin.Context) {
	strUserID, exists := ctx.Get("userID")
	if !exists {
		response.Error(ctx, http.StatusUnauthorized, "Unauthorized", "User ID not found in context")
		return
	}

	userID, err := uuid.Parse(strUserID.(string))
	if err != nil {
		response.Error(ctx, http.StatusUnauthorized, "Unauthorized", "Invalid User ID format")
		return
	}

	var req dto.BatchDeleteCartItemsRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, http.StatusBadRequest, "Invalid request", err.Error())
		return
	}

	if err := c.useCase.BatchDeleteCartItems(ctx.Request.Context(), userID, &req); err != nil {
		response.Error(ctx, http.StatusInternalServerError, "Failed to delete items", err.Error())
		return
	}

	response.Success(ctx, http.StatusOK, "Items removed from cart successfully", nil)
}
