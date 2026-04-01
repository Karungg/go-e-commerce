package cart

import (
	"context"

	"go-e-commerce/internal/entity"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type CartRepositoryMock struct {
	mock.Mock
}

func (m *CartRepositoryMock) GetCartByUserID(ctx context.Context, userID uuid.UUID) (*entity.Cart, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) != nil {
		return args.Get(0).(*entity.Cart), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *CartRepositoryMock) CreateCart(ctx context.Context, cart *entity.Cart) error {
	args := m.Called(ctx, cart)
	return args.Error(0)
}

func (m *CartRepositoryMock) GetCartItem(ctx context.Context, cartID, productID uuid.UUID) (*entity.CartItem, error) {
	args := m.Called(ctx, cartID, productID)
	if args.Get(0) != nil {
		return args.Get(0).(*entity.CartItem), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *CartRepositoryMock) AddCartItem(ctx context.Context, item *entity.CartItem) error {
	args := m.Called(ctx, item)
	return args.Error(0)
}

func (m *CartRepositoryMock) UpdateCartItemQuantity(ctx context.Context, itemID uuid.UUID, quantity int) error {
	args := m.Called(ctx, itemID, quantity)
	return args.Error(0)
}

func (m *CartRepositoryMock) DeleteCartItems(ctx context.Context, cartID uuid.UUID, itemIDs []uuid.UUID) error {
	args := m.Called(ctx, cartID, itemIDs)
	return args.Error(0)
}
