package port

import (
	"context"
	"go-e-commerce/internal/entity"
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
