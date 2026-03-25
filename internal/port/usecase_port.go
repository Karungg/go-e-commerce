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

// CategoryUseCase defines the contract for category business logic
type CategoryUseCase interface {
	CreateCategory(ctx context.Context, req *dto.CreateCategoryReq) (*dto.CategoryRes, error)
	GetAllCategories(ctx context.Context) ([]*dto.CategoryRes, error)
	GetCategoryByID(ctx context.Context, id string) (*dto.CategoryRes, error)
	UpdateCategory(ctx context.Context, id string, req *dto.UpdateCategoryReq) (*dto.CategoryRes, error)
	DeleteCategory(ctx context.Context, id string) error
}
