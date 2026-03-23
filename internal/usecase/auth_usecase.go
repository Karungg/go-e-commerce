package usecase

import (
	"errors"
	"go-e-commerce/internal/entity"
	"go-e-commerce/internal/repository"
	"go-e-commerce/internal/security"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthUseCase interface {
	RegisterCustomer(req *RegisterCustomerReq) (string, error)
	RegisterSeller(req *RegisterSellerReq) (string, error)
}

type authUseCase struct {
	db           *gorm.DB
	userRepo     repository.UserRepository
	customerRepo repository.CustomerRepository
	sellerRepo   repository.SellerRepository
	jwtAuth      *security.JWTAuth
}

// Request DTOs
type RegisterCustomerReq struct {
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=6"`
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
	Phone     string `json:"phone"`
	Address   string `json:"address"`
}

type RegisterSellerReq struct {
	Email            string `json:"email" binding:"required,email"`
	Password         string `json:"password" binding:"required,min=6"`
	StoreName        string `json:"store_name" binding:"required"`
	StoreDescription string `json:"store_description"`
	LogoUrl          string `json:"logo_url"`
}

func NewAuthUseCase(
	db *gorm.DB,
	userRepo repository.UserRepository,
	customerRepo repository.CustomerRepository,
	sellerRepo repository.SellerRepository,
	jwtAuth *security.JWTAuth,
) AuthUseCase {
	return &authUseCase{
		db:           db,
		userRepo:     userRepo,
		customerRepo: customerRepo,
		sellerRepo:   sellerRepo,
		jwtAuth:      jwtAuth,
	}
}

func (u *authUseCase) hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func (u *authUseCase) RegisterCustomer(req *RegisterCustomerReq) (string, error) {
	existingUser, _ := u.userRepo.FindByEmail(req.Email)
	if existingUser != nil {
		return "", errors.New("email is already registered")
	}

	hashedPwd, err := u.hashPassword(req.Password)
	if err != nil {
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

	// Run within a database transaction
	err = u.db.Transaction(func(tx *gorm.DB) error {
		if err := u.userRepo.CreateWithTx(tx, user); err != nil {
			return err
		}
		if err := u.customerRepo.CreateWithTx(tx, customer); err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return "", err
	}

	return u.jwtAuth.GenerateToken(userID, string(entity.RoleCustomer))
}

func (u *authUseCase) RegisterSeller(req *RegisterSellerReq) (string, error) {
	existingUser, _ := u.userRepo.FindByEmail(req.Email)
	if existingUser != nil {
		return "", errors.New("email is already registered")
	}

	hashedPwd, err := u.hashPassword(req.Password)
	if err != nil {
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

	// Run within a database transaction
	err = u.db.Transaction(func(tx *gorm.DB) error {
		if err := u.userRepo.CreateWithTx(tx, user); err != nil {
			return err
		}
		if err := u.sellerRepo.CreateWithTx(tx, seller); err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return "", err
	}

	return u.jwtAuth.GenerateToken(userID, string(entity.RoleSeller))
}
