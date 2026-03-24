package mocks

import (
	"context"
	"go-e-commerce/internal/entity"

	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

type SellerRepositoryMock struct {
	mock.Mock
}

func (m *SellerRepositoryMock) CreateWithTx(ctx context.Context, tx *gorm.DB, seller *entity.Seller) error {
	args := m.Called(ctx, tx, seller)
	return args.Error(0)
}

func (m *SellerRepositoryMock) FindByStoreName(ctx context.Context, storeName string) (*entity.Seller, error) {
	args := m.Called(ctx, storeName)
	var seller *entity.Seller
	if args.Get(0) != nil {
		seller = args.Get(0).(*entity.Seller)
	}
	return seller, args.Error(1)
}
