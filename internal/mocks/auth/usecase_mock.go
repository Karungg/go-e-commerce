package auth

import (
	"context"

	authDTO "go-e-commerce/internal/dto/auth"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type AuthUseCaseMock struct {
	mock.Mock
}

func (m *AuthUseCaseMock) RegisterCustomer(ctx context.Context, req *authDTO.RegisterCustomerReq) (string, error) {
	args := m.Called(ctx, req)
	return args.String(0), args.Error(1)
}

func (m *AuthUseCaseMock) RegisterSeller(ctx context.Context, req *authDTO.RegisterSellerReq) (string, error) {
	args := m.Called(ctx, req)
	return args.String(0), args.Error(1)
}

func (m *AuthUseCaseMock) Login(ctx context.Context, req *authDTO.LoginReq) (string, error) {
	args := m.Called(ctx, req)
	return args.String(0), args.Error(1)
}

func (m *AuthUseCaseMock) Logout(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *AuthUseCaseMock) UpdateCustomer(ctx context.Context, userID uuid.UUID, req *authDTO.UpdateCustomerReq) error {
	args := m.Called(ctx, userID, req)
	return args.Error(0)
}
