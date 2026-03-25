package auth

import (
	"context"
	"fmt"
	"log/slog"

	authDTO "go-e-commerce/internal/dto/auth"
	"go-e-commerce/internal/entity"
	"go-e-commerce/internal/pkg/apperror"
	authPort "go-e-commerce/internal/port/auth"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type authUseCase struct {
	txManager      authPort.TransactionManager
	logger         *slog.Logger
	userRepo       authPort.UserRepository
	customerRepo   authPort.CustomerRepository
	sellerRepo     authPort.SellerRepository
	tokenGenerator authPort.TokenGenerator
}

func NewAuthUseCase(
	txManager authPort.TransactionManager,
	logger *slog.Logger,
	userRepo authPort.UserRepository,
	customerRepo authPort.CustomerRepository,
	sellerRepo authPort.SellerRepository,
	tokenGenerator authPort.TokenGenerator,
) authPort.AuthUseCase {
	return &authUseCase{
		txManager:      txManager,
		logger:         logger,
		userRepo:       userRepo,
		customerRepo:   customerRepo,
		sellerRepo:     sellerRepo,
		tokenGenerator: tokenGenerator,
	}
}

func (u *authUseCase) hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	return string(bytes), nil
}

func (u *authUseCase) comparePassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func (u *authUseCase) Login(ctx context.Context, req *authDTO.LoginReq) (string, error) {
	user, err := u.userRepo.FindByEmail(ctx, req.Email)
	if err != nil || user == nil {
		u.logger.WarnContext(ctx, "Login failed: user not found", slog.String("email", req.Email))
		return "", apperror.ErrInvalidPassword
	}

	if err := u.comparePassword(user.Password, req.Password); err != nil {
		u.logger.WarnContext(ctx, "Login failed: invalid password", slog.String("email", req.Email))
		return "", apperror.ErrInvalidPassword
	}

	token, err := u.tokenGenerator.GenerateToken(user.ID, string(user.Role))
	if err != nil {
		u.logger.ErrorContext(ctx, "Failed to generate JWT", slog.Any("error", err))
		return "", fmt.Errorf("failed to generate authentication token: %w", err)
	}

	u.logger.InfoContext(ctx, "User logged in successfully", slog.String("user_id", user.ID.String()))
	return token, nil
}

func (u *authUseCase) Logout(ctx context.Context) error {
	u.logger.InfoContext(ctx, "User logged out successfully")
	return nil
}

func (u *authUseCase) RegisterCustomer(ctx context.Context, req *authDTO.RegisterCustomerReq) (string, error) {
	existingUser, _ := u.userRepo.FindByEmail(ctx, req.Email)
	if existingUser != nil {
		u.logger.WarnContext(ctx, "Registration failed due to email conflict", slog.String("email", req.Email))
		return "", apperror.ErrEmailConflict
	}

	if req.Phone != "" {
		existingCustomer, _ := u.customerRepo.FindByPhone(ctx, req.Phone)
		if existingCustomer != nil {
			u.logger.WarnContext(ctx, "Registration failed due to phone conflict", slog.String("phone", req.Phone))
			return "", apperror.ErrPhoneConflict
		}
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

	err = u.txManager.RunInTransaction(ctx, func(ctx context.Context) error {
		if err := u.userRepo.Create(ctx, user); err != nil {
			return fmt.Errorf("failed to insert user record: %w", err)
		}
		if err := u.customerRepo.Create(ctx, customer); err != nil {
			return fmt.Errorf("failed to insert customer record: %w", err)
		}
		return nil
	})

	if err != nil {
		u.logger.ErrorContext(ctx, "Transaction failed during customer registration", slog.Any("error", err))
		return "", err
	}

	token, err := u.tokenGenerator.GenerateToken(userID, string(entity.RoleCustomer))
	if err != nil {
		u.logger.ErrorContext(ctx, "Failed to generate JWT", slog.Any("error", err))
		return "", fmt.Errorf("failed to generate authentication token: %w", err)
	}

	u.logger.InfoContext(ctx, "Customer registered successfully", slog.String("user_id", userID.String()))
	return token, nil
}

func (u *authUseCase) RegisterSeller(ctx context.Context, req *authDTO.RegisterSellerReq) (string, error) {
	existingUser, _ := u.userRepo.FindByEmail(ctx, req.Email)
	if existingUser != nil {
		u.logger.WarnContext(ctx, "Registration failed due to email conflict", slog.String("email", req.Email))
		return "", apperror.ErrEmailConflict
	}

	existingSeller, _ := u.sellerRepo.FindByStoreName(ctx, req.StoreName)
	if existingSeller != nil {
		u.logger.WarnContext(ctx, "Registration failed due to store name conflict", slog.String("store_name", req.StoreName))
		return "", apperror.ErrStoreNameConflict
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

	err = u.txManager.RunInTransaction(ctx, func(ctx context.Context) error {
		if err := u.userRepo.Create(ctx, user); err != nil {
			return fmt.Errorf("failed to insert user record: %w", err)
		}
		if err := u.sellerRepo.Create(ctx, seller); err != nil {
			return fmt.Errorf("failed to insert seller record: %w", err)
		}
		return nil
	})

	if err != nil {
		u.logger.ErrorContext(ctx, "Transaction failed during seller registration", slog.Any("error", err))
		return "", err
	}

	token, err := u.tokenGenerator.GenerateToken(userID, string(entity.RoleSeller))
	if err != nil {
		u.logger.ErrorContext(ctx, "Failed to generate JWT", slog.Any("error", err))
		return "", fmt.Errorf("failed to generate authentication token: %w", err)
	}

	u.logger.InfoContext(ctx, "Seller registered successfully", slog.String("user_id", userID.String()))
	return token, nil
}
