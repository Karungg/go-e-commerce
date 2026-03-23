package mocks

import (
	"go-e-commerce/internal/entity"

	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

type SellerRepositoryMock struct {
	mock.Mock
}

func (m *SellerRepositoryMock) CreateWithTx(tx *gorm.DB, seller *entity.Seller) error {
	args := m.Called(tx, seller)
	return args.Error(0)
}
