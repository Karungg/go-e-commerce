package auth

import (
	"context"
	"go-e-commerce/internal/entity"

	"github.com/stretchr/testify/mock"
)

type SellerRepositoryMock struct {
	mock.Mock
}

func (m *SellerRepositoryMock) Create(ctx context.Context, seller *entity.Seller) error {
	args := m.Called(ctx, seller)
	return args.Error(0)
}

func (m *SellerRepositoryMock) FindByStoreName(ctx context.Context, storeName string) (*entity.Seller, error) {
	args := m.Called(ctx, storeName)
	if args.Get(0) != nil {
		return args.Get(0).(*entity.Seller), args.Error(1)
	}
	return nil, args.Error(1)
}
