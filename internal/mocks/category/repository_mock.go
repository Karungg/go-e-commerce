package category

import (
	"context"
	"go-e-commerce/internal/entity"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type CategoryRepositoryMock struct {
	mock.Mock
}

func (m *CategoryRepositoryMock) Create(ctx context.Context, category *entity.Category) error {
	args := m.Called(ctx, category)
	return args.Error(0)
}

func (m *CategoryRepositoryMock) FindAll(ctx context.Context) ([]*entity.Category, error) {
	args := m.Called(ctx)
	if args.Get(0) != nil {
		return args.Get(0).([]*entity.Category), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *CategoryRepositoryMock) FindByID(ctx context.Context, id uuid.UUID) (*entity.Category, error) {
	args := m.Called(ctx, id)
	if args.Get(0) != nil {
		return args.Get(0).(*entity.Category), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *CategoryRepositoryMock) Update(ctx context.Context, category *entity.Category) error {
	args := m.Called(ctx, category)
	return args.Error(0)
}

func (m *CategoryRepositoryMock) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
