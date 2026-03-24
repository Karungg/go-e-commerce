package port

import (
	"context"
	"go-e-commerce/internal/entity"

	"gorm.io/gorm"
)

// UserRepository defines the contract for user data access
type UserRepository interface {
	CreateWithTx(ctx context.Context, tx *gorm.DB, user *entity.User) error
	FindByEmail(ctx context.Context, email string) (*entity.User, error)
}

// CustomerRepository defines the contract for customer data access
type CustomerRepository interface {
	CreateWithTx(ctx context.Context, tx *gorm.DB, customer *entity.Customer) error
	FindByPhone(ctx context.Context, phone string) (*entity.Customer, error)
}

// SellerRepository defines the contract for seller data access
type SellerRepository interface {
	CreateWithTx(ctx context.Context, tx *gorm.DB, seller *entity.Seller) error
	FindByStoreName(ctx context.Context, storeName string) (*entity.Seller, error)
}
