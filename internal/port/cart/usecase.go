package cart

import (
	"context"

	"go-e-commerce/internal/dto"
	"github.com/google/uuid"
)

type CartUseCase interface {
	GetCart(ctx context.Context, userID uuid.UUID) (*dto.CartResponse, error)
	AddToCart(ctx context.Context, userID uuid.UUID, req *dto.AddCartItemRequest) error
	UpdateCartItem(ctx context.Context, userID uuid.UUID, itemID uuid.UUID, req *dto.UpdateCartItemRequest) error
	BatchDeleteCartItems(ctx context.Context, userID uuid.UUID, req *dto.BatchDeleteCartItemsRequest) error
}
