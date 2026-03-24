package mocks

import (
	"context"
	"go-e-commerce/internal/entity"

	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

type CustomerRepositoryMock struct {
	mock.Mock
}

func (m *CustomerRepositoryMock) CreateWithTx(ctx context.Context, tx *gorm.DB, customer *entity.Customer) error {
	args := m.Called(ctx, tx, customer)
	return args.Error(0)
}

func (m *CustomerRepositoryMock) FindByPhone(ctx context.Context, phone string) (*entity.Customer, error) {
	args := m.Called(ctx, phone)
	var cust *entity.Customer
	if args.Get(0) != nil {
		cust = args.Get(0).(*entity.Customer)
	}
	return cust, args.Error(1)
}
