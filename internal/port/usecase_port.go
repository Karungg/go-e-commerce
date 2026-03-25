package port

import (
	"context"
	"go-e-commerce/internal/dto"
)

// AuthUseCase defines the contract for authentication business logic
type AuthUseCase interface {
	RegisterCustomer(ctx context.Context, req *dto.RegisterCustomerReq) (string, error)
	RegisterSeller(ctx context.Context, req *dto.RegisterSellerReq) (string, error)
	Login(ctx context.Context, req *dto.LoginReq) (string, error)
	Logout(ctx context.Context) error
}
