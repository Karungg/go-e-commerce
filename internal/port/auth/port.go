package auth

import (
	"context"

	authDTO "go-e-commerce/internal/dto/auth"
	"go-e-commerce/internal/entity"

	"github.com/google/uuid"
)

type UserRepository interface {
	Create(ctx context.Context, user *entity.User) error
	FindByEmail(ctx context.Context, email string) (*entity.User, error)
}

type CustomerRepository interface {
	Create(ctx context.Context, customer *entity.Customer) error
	FindByPhone(ctx context.Context, phone string) (*entity.Customer, error)
	FindByUserID(ctx context.Context, userID uuid.UUID) (*entity.Customer, error)
	Update(ctx context.Context, customer *entity.Customer) error
}

type SellerRepository interface {
	Create(ctx context.Context, seller *entity.Seller) error
	FindByStoreName(ctx context.Context, storeName string) (*entity.Seller, error)
}

type TokenGenerator interface {
	GenerateToken(userID uuid.UUID, role string) (string, error)
}

type TokenPayload struct {
	UserID uuid.UUID
	Role   string
}

type TokenValidator interface {
	ValidateToken(tokenString string) (*TokenPayload, error)
}
type TransactionManager interface {
	RunInTransaction(ctx context.Context, fn func(ctx context.Context) error) error
}

type AuthUseCase interface {
	RegisterCustomer(ctx context.Context, req *authDTO.RegisterCustomerReq) (string, error)
	RegisterSeller(ctx context.Context, req *authDTO.RegisterSellerReq) (string, error)
	Login(ctx context.Context, req *authDTO.LoginReq) (string, error)
	Logout(ctx context.Context) error
	UpdateCustomer(ctx context.Context, userID uuid.UUID, req *authDTO.UpdateCustomerReq) error
}
