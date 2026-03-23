package mocks

import (
	"go-e-commerce/internal/usecase"

	"github.com/stretchr/testify/mock"
)

type AuthUseCaseMock struct {
	mock.Mock
}

func (m *AuthUseCaseMock) RegisterCustomer(req *usecase.RegisterCustomerReq) (string, error) {
	args := m.Called(req)
	return args.String(0), args.Error(1)
}

func (m *AuthUseCaseMock) RegisterSeller(req *usecase.RegisterSellerReq) (string, error) {
	args := m.Called(req)
	return args.String(0), args.Error(1)
}
