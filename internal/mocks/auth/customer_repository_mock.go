package auth

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
	if args.Get(0) != nil {
		return args.Get(0).(*entity.Customer), args.Error(1)
	}
	return nil, args.Error(1)
}
