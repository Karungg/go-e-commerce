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
