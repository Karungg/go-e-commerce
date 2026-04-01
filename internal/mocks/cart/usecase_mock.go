package cart

import (
	"context"

	"go-e-commerce/internal/dto"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type CartUseCaseMock struct {
	mock.Mock
}

func (m *CartUseCaseMock) GetCart(ctx context.Context, userID uuid.UUID) (*dto.CartResponse, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) != nil {
		return args.Get(0).(*dto.CartResponse), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *CartUseCaseMock) AddToCart(ctx context.Context, userID uuid.UUID, req *dto.AddCartItemRequest) error {
	args := m.Called(ctx, userID, req)
	return args.Error(0)
}

func (m *CartUseCaseMock) UpdateCartItem(ctx context.Context, userID uuid.UUID, itemID uuid.UUID, req *dto.UpdateCartItemRequest) error {
	args := m.Called(ctx, userID, itemID, req)
	return args.Error(0)
}

func (m *CartUseCaseMock) BatchDeleteCartItems(ctx context.Context, userID uuid.UUID, req *dto.BatchDeleteCartItemsRequest) error {
	args := m.Called(ctx, userID, req)
	return args.Error(0)
}
