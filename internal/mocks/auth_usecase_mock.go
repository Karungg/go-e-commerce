package mocks

import (
	"context"
	"go-e-commerce/internal/dto"

	"github.com/stretchr/testify/mock"
)

type AuthUseCaseMock struct {
	mock.Mock
}

func (m *AuthUseCaseMock) RegisterCustomer(ctx context.Context, req *dto.RegisterCustomerReq) (string, error) {
	args := m.Called(ctx, req)
	return args.String(0), args.Error(1)
}

func (m *AuthUseCaseMock) RegisterSeller(ctx context.Context, req *dto.RegisterSellerReq) (string, error) {
	args := m.Called(ctx, req)
	return args.String(0), args.Error(1)
}

func (m *AuthUseCaseMock) Login(ctx context.Context, req *dto.LoginReq) (string, error) {
	args := m.Called(ctx, req)
	return args.String(0), args.Error(1)
}

func (m *AuthUseCaseMock) Logout(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

