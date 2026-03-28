package product

import (
	"context"
	"go-e-commerce/internal/entity"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type ProductRepositoryMock struct {
	mock.Mock
}

func (m *ProductRepositoryMock) Create(ctx context.Context, product *entity.Product) error {
	args := m.Called(ctx, product)
	return args.Error(0)
}

func (m *ProductRepositoryMock) FindAll(ctx context.Context, limit, offset int) ([]*entity.Product, int64, error) {
	args := m.Called(ctx, limit, offset)
	if args.Get(0) != nil {
		return args.Get(0).([]*entity.Product), args.Get(1).(int64), args.Error(2)
	}
	return nil, args.Get(1).(int64), args.Error(2)
}

func (m *ProductRepositoryMock) FindByID(ctx context.Context, id uuid.UUID) (*entity.Product, error) {
	args := m.Called(ctx, id)
	if args.Get(0) != nil {
		return args.Get(0).(*entity.Product), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *ProductRepositoryMock) Update(ctx context.Context, product *entity.Product) error {
	args := m.Called(ctx, product)
	return args.Error(0)
}

func (m *ProductRepositoryMock) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
