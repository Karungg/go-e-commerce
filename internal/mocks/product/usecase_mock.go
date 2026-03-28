package product

import (
	"context"

	productDTO "go-e-commerce/internal/dto/product"

	"github.com/stretchr/testify/mock"
)

type ProductUseCaseMock struct {
	mock.Mock
}

func (m *ProductUseCaseMock) CreateProduct(ctx context.Context, req *productDTO.CreateProductReq) (*productDTO.ProductRes, error) {
	args := m.Called(ctx, req)
	if args.Get(0) != nil {
		return args.Get(0).(*productDTO.ProductRes), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *ProductUseCaseMock) GetAllProducts(ctx context.Context, page, limit int) ([]*productDTO.ProductRes, int64, error) {
	args := m.Called(ctx, page, limit)
	if args.Get(0) != nil {
		return args.Get(0).([]*productDTO.ProductRes), args.Get(1).(int64), args.Error(2)
	}
	return nil, args.Get(1).(int64), args.Error(2)
}

func (m *ProductUseCaseMock) GetProductByID(ctx context.Context, id string) (*productDTO.ProductRes, error) {
	args := m.Called(ctx, id)
	if args.Get(0) != nil {
		return args.Get(0).(*productDTO.ProductRes), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *ProductUseCaseMock) UpdateProduct(ctx context.Context, id string, req *productDTO.UpdateProductReq) (*productDTO.ProductRes, error) {
	args := m.Called(ctx, id, req)
	if args.Get(0) != nil {
		return args.Get(0).(*productDTO.ProductRes), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *ProductUseCaseMock) DeleteProduct(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
