package cart

import (
	"context"

	"go-e-commerce/internal/entity"
	"github.com/google/uuid"
)

type CartRepository interface {
	GetCartByUserID(ctx context.Context, userID uuid.UUID) (*entity.Cart, error)
	CreateCart(ctx context.Context, cart *entity.Cart) error
	GetCartItem(ctx context.Context, cartID, productID uuid.UUID) (*entity.CartItem, error)
	AddCartItem(ctx context.Context, item *entity.CartItem) error
	UpdateCartItemQuantity(ctx context.Context, itemID uuid.UUID, quantity int) error
	DeleteCartItems(ctx context.Context, cartID uuid.UUID, itemIDs []uuid.UUID) error
}
