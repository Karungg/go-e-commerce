package port

import (
	"context"
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

// CategoryRepository defines the contract for category data access
type CategoryRepository interface {
	Create(ctx context.Context, category *entity.Category) error
	FindAll(ctx context.Context) ([]*entity.Category, error)
	FindByID(ctx context.Context, id uuid.UUID) (*entity.Category, error)
	Update(ctx context.Context, category *entity.Category) error
	Delete(ctx context.Context, id uuid.UUID) error
}
