package usecase

import (
	"context"
	"fmt"
	"log/slog"
	
	"go-e-commerce/internal/entity"
	"go-e-commerce/internal/pkg/apperror"
	"go-e-commerce/internal/repository"
	"go-e-commerce/internal/security"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthUseCase interface {
	RegisterCustomer(ctx context.Context, req *RegisterCustomerReq) (string, error)
	RegisterSeller(ctx context.Context, req *RegisterSellerReq) (string, error)
}

type authUseCase struct {
	db           *gorm.DB
	logger       *slog.Logger
	userRepo     repository.UserRepository
	customerRepo repository.CustomerRepository
	sellerRepo   repository.SellerRepository
	jwtAuth      *security.JWTAuth
}

// Request DTOs
type RegisterCustomerReq struct {
	Email     string `json:"email" binding:"required,email,max=100"`
	Password  string `json:"password" binding:"required,min=8,max=64"`
	FirstName string `json:"first_name" binding:"required,alpha,min=2,max=50"`
	LastName  string `json:"last_name" binding:"required,alpha,min=2,max=50"`
	Phone     string `json:"phone" binding:"omitempty,numeric,min=10,max=15"` 
	Address   string `json:"address" binding:"omitempty,max=255"`
}

type RegisterSellerReq struct {
	Email            string `json:"email" binding:"required,email,max=100"`
	Password         string `json:"password" binding:"required,min=8,max=64"`
	StoreName        string `json:"store_name" binding:"required,min=3,max=100"`
	StoreDescription string `json:"store_description" binding:"omitempty,max=500"`
	LogoUrl          string `json:"logo_url" binding:"omitempty,url,max=255"`
}

func NewAuthUseCase(
	db *gorm.DB,
	logger *slog.Logger,
	userRepo repository.UserRepository,
	customerRepo repository.CustomerRepository,
	sellerRepo repository.SellerRepository,
	jwtAuth *security.JWTAuth,
) AuthUseCase {
	return &authUseCase{
		db:           db,
		logger:       logger,
		userRepo:     userRepo,
		customerRepo: customerRepo,
		sellerRepo:   sellerRepo,
		jwtAuth:      jwtAuth,
	}
}

func (u *authUseCase) hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	return string(bytes), nil
}

func (u *authUseCase) RegisterCustomer(ctx context.Context, req *RegisterCustomerReq) (string, error) {
	existingUser, _ := u.userRepo.FindByEmail(ctx, req.Email)
	if existingUser != nil {
		u.logger.WarnContext(ctx, "Registration failed due to email conflict", slog.String("email", req.Email))
		return "", apperror.ErrEmailConflict
	}

	hashedPwd, err := u.hashPassword(req.Password)
	if err != nil {
		u.logger.ErrorContext(ctx, "Password hashing failed", slog.Any("error", err))
		return "", err
	}

	userID := uuid.New()
	user := &entity.User{
		ID:       userID,
		Email:    req.Email,
		Password: hashedPwd,
		Role:     entity.RoleCustomer,
		IsActive: true,
	}

	customer := &entity.Customer{
		ID:        uuid.New(),
		UserID:    userID,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Phone:     req.Phone,
		Address:   req.Address,
	}

	err = u.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := u.userRepo.CreateWithTx(ctx, tx, user); err != nil {
			return fmt.Errorf("failed to insert user record: %w", err)
		}
		if err := u.customerRepo.CreateWithTx(ctx, tx, customer); err != nil {
			return fmt.Errorf("failed to insert customer record: %w", err)
		}
		return nil
	})

	if err != nil {
		u.logger.ErrorContext(ctx, "Transaction failed during customer registration", slog.Any("error", err))
		return "", err
	}

	token, err := u.jwtAuth.GenerateToken(userID, string(entity.RoleCustomer))
	if err != nil {
		u.logger.ErrorContext(ctx, "Failed to generate JWT", slog.Any("error", err))
		return "", fmt.Errorf("failed to generate authentication token: %w", err)
	}

	u.logger.InfoContext(ctx, "Customer registered successfully", slog.String("user_id", userID.String()))
	return token, nil
}

func (u *authUseCase) RegisterSeller(ctx context.Context, req *RegisterSellerReq) (string, error) {
	existingUser, _ := u.userRepo.FindByEmail(ctx, req.Email)
	if existingUser != nil {
		u.logger.WarnContext(ctx, "Registration failed due to email conflict", slog.String("email", req.Email))
		return "", apperror.ErrEmailConflict
	}

	hashedPwd, err := u.hashPassword(req.Password)
	if err != nil {
		u.logger.ErrorContext(ctx, "Password hashing failed", slog.Any("error", err))
		return "", err
	}

	userID := uuid.New()
	user := &entity.User{
		ID:       userID,
		Email:    req.Email,
		Password: hashedPwd,
		Role:     entity.RoleSeller,
		IsActive: true,
	}

	seller := &entity.Seller{
		ID:               uuid.New(),
		UserID:           userID,
		StoreName:        req.StoreName,
		StoreDescription: req.StoreDescription,
		LogoUrl:          req.LogoUrl,
		IsVerified:       false,
	}

	err = u.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := u.userRepo.CreateWithTx(ctx, tx, user); err != nil {
			return fmt.Errorf("failed to insert user record: %w", err)
		}
		if err := u.sellerRepo.CreateWithTx(ctx, tx, seller); err != nil {
			return fmt.Errorf("failed to insert seller record: %w", err)
		}
		return nil
	})

	if err != nil {
		u.logger.ErrorContext(ctx, "Transaction failed during seller registration", slog.Any("error", err))
		return "", err
	}

	token, err := u.jwtAuth.GenerateToken(userID, string(entity.RoleSeller))
	if err != nil {
		u.logger.ErrorContext(ctx, "Failed to generate JWT", slog.Any("error", err))
		return "", fmt.Errorf("failed to generate authentication token: %w", err)
	}

	u.logger.InfoContext(ctx, "Seller registered successfully", slog.String("user_id", userID.String()))
	return token, nil
}
