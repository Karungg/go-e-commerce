package category

import (
	"context"

	categoryDTO "go-e-commerce/internal/dto/category"

	"github.com/stretchr/testify/mock"
)

type CategoryUseCaseMock struct {
	mock.Mock
}

func (m *CategoryUseCaseMock) CreateCategory(ctx context.Context, req *categoryDTO.CreateCategoryReq) (*categoryDTO.CategoryRes, error) {
	args := m.Called(ctx, req)
	if args.Get(0) != nil {
		return args.Get(0).(*categoryDTO.CategoryRes), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *CategoryUseCaseMock) GetAllCategories(ctx context.Context) ([]*categoryDTO.CategoryRes, error) {
	args := m.Called(ctx)
	if args.Get(0) != nil {
		return args.Get(0).([]*categoryDTO.CategoryRes), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *CategoryUseCaseMock) GetCategoryByID(ctx context.Context, id string) (*categoryDTO.CategoryRes, error) {
	args := m.Called(ctx, id)
	if args.Get(0) != nil {
		return args.Get(0).(*categoryDTO.CategoryRes), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *CategoryUseCaseMock) UpdateCategory(ctx context.Context, id string, req *categoryDTO.UpdateCategoryReq) (*categoryDTO.CategoryRes, error) {
	args := m.Called(ctx, id, req)
	if args.Get(0) != nil {
		return args.Get(0).(*categoryDTO.CategoryRes), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *CategoryUseCaseMock) DeleteCategory(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
