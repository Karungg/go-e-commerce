package mocks

import (
	"context"
	"go-e-commerce/internal/entity"

	"github.com/stretchr/testify/mock"
)

type CustomerRepositoryMock struct {
	mock.Mock
}

func (m *CustomerRepositoryMock) Create(ctx context.Context, customer *entity.Customer) error {
	args := m.Called(ctx, customer)
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
