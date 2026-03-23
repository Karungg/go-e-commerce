package mocks

import (
	"context"
	"go-e-commerce/internal/usecase"

	"github.com/stretchr/testify/mock"
)

type AuthUseCaseMock struct {
	mock.Mock
}

func (m *AuthUseCaseMock) RegisterCustomer(ctx context.Context, req *usecase.RegisterCustomerReq) (string, error) {
	args := m.Called(ctx, req)
	return args.String(0), args.Error(1)
}

func (m *AuthUseCaseMock) RegisterSeller(ctx context.Context, req *usecase.RegisterSellerReq) (string, error) {
	args := m.Called(ctx, req)
	return args.String(0), args.Error(1)
}
