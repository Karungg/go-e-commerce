package auth

import (
	"context"

	authDTO "go-e-commerce/internal/dto/auth"
	"go-e-commerce/internal/entity"

	"github.com/google/uuid"
)

// UserRepository defines the contract for user data access
type UserRepository interface {
	Create(ctx context.Context, user *entity.User) error
	FindByEmail(ctx context.Context, email string) (*entity.User, error)
}

// CustomerRepository defines the contract for customer data access
type CustomerRepository interface {
	Create(ctx context.Context, customer *entity.Customer) error
	FindByPhone(ctx context.Context, phone string) (*entity.Customer, error)
}

// SellerRepository defines the contract for seller data access
type SellerRepository interface {
	Create(ctx context.Context, seller *entity.Seller) error
	FindByStoreName(ctx context.Context, storeName string) (*entity.Seller, error)
}

// TokenGenerator defines the contract for generating authentication credentials
type TokenGenerator interface {
	GenerateToken(userID uuid.UUID, role string) (string, error)
}

// TokenPayload holds the essential identity claims from an authenticated token
type TokenPayload struct {
	UserID uuid.UUID
	Role   string
}

// TokenValidator defines the contract for validating authentication credentials
type TokenValidator interface {
	ValidateToken(tokenString string) (*TokenPayload, error)
}

// TransactionManager defines the contract for running operations within a database transaction
type TransactionManager interface {
	RunInTransaction(ctx context.Context, fn func(ctx context.Context) error) error
}

// AuthUseCase defines the contract for authentication business logic
type AuthUseCase interface {
	RegisterCustomer(ctx context.Context, req *authDTO.RegisterCustomerReq) (string, error)
	RegisterSeller(ctx context.Context, req *authDTO.RegisterSellerReq) (string, error)
	Login(ctx context.Context, req *authDTO.LoginReq) (string, error)
	Logout(ctx context.Context) error
}
