package cart

import "github.com/google/uuid"

type AddCartItemRequest struct {
	ProductID uuid.UUID `json:"product_id" binding:"required"`
	Quantity  int       `json:"quantity" binding:"required,min=1"`
}

type UpdateCartItemRequest struct {
	Quantity int `json:"quantity" binding:"required,min=1"`
}

type BatchDeleteCartItemsRequest struct {
	CartItemIDs []uuid.UUID `json:"cart_item_ids" binding:"required,min=1"`
}
