package cart

import (
	"context"

	cartDTO "go-e-commerce/internal/dto/cart"
	"github.com/google/uuid"
)

type CartUseCase interface {
	GetCart(ctx context.Context, userID uuid.UUID) (*cartDTO.CartResponse, error)
	AddToCart(ctx context.Context, userID uuid.UUID, req *cartDTO.AddCartItemRequest) error
	UpdateCartItem(ctx context.Context, userID uuid.UUID, itemID uuid.UUID, req *cartDTO.UpdateCartItemRequest) error
	BatchDeleteCartItems(ctx context.Context, userID uuid.UUID, req *cartDTO.BatchDeleteCartItemsRequest) error
}
