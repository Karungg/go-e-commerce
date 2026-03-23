package usecase_test

import (
	"errors"
	"go-e-commerce/internal/entity"
	"go-e-commerce/internal/mocks"
	"go-e-commerce/internal/security"
	"go-e-commerce/internal/usecase"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
	sqlDB, mock, err := sqlmock.New()
	assert.NoError(t, err)

	grmDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: sqlDB,
	}), &gorm.Config{
		SkipDefaultTransaction: true,
	})
	assert.NoError(t, err)

	return grmDB, mock
}

func TestRegisterCustomer_Success(t *testing.T) {
	db, sqlMock := setupTestDB(t)
	userRepo := new(mocks.UserRepositoryMock)
	customerRepo := new(mocks.CustomerRepositoryMock)
	sellerRepo := new(mocks.SellerRepositoryMock)
	jwtAuth := security.NewJWTAuth("secret", 24)

	uc := usecase.NewAuthUseCase(db, userRepo, customerRepo, sellerRepo, jwtAuth)

	req := &usecase.RegisterCustomerReq{
		Email:     "test@example.com",
		Password:  "password123",
		FirstName: "John",
		LastName:  "Doe",
		Phone:     "123456789",
		Address:   "123 Street",
	}

	userRepo.On("FindByEmail", req.Email).Return(nil, errors.New("not found"))

	// Expect GORM Transaction boundaries
	sqlMock.ExpectBegin()
	sqlMock.ExpectCommit()

	userRepo.On("CreateWithTx", mock.Anything, mock.AnythingOfType("*entity.User")).Return(nil)
	customerRepo.On("CreateWithTx", mock.Anything, mock.AnythingOfType("*entity.Customer")).Return(nil)

	token, err := uc.RegisterCustomer(req)

	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	userRepo.AssertExpectations(t)
	customerRepo.AssertExpectations(t)
	assert.NoError(t, sqlMock.ExpectationsWereMet())
}

func TestRegisterCustomer_EmailExists(t *testing.T) {
	db, _ := setupTestDB(t)
	userRepo := new(mocks.UserRepositoryMock)
	customerRepo := new(mocks.CustomerRepositoryMock)
	sellerRepo := new(mocks.SellerRepositoryMock)
	jwtAuth := security.NewJWTAuth("secret", 24)

	uc := usecase.NewAuthUseCase(db, userRepo, customerRepo, sellerRepo, jwtAuth)

	req := &usecase.RegisterCustomerReq{
		Email: "test@example.com",
	}

	existingUser := &entity.User{Email: req.Email}
	userRepo.On("FindByEmail", req.Email).Return(existingUser, nil)

	token, err := uc.RegisterCustomer(req)

	assert.Error(t, err)
	assert.Equal(t, "email is already registered", err.Error())
	assert.Empty(t, token)
}

func TestRegisterSeller_Success(t *testing.T) {
	db, sqlMock := setupTestDB(t)
	userRepo := new(mocks.UserRepositoryMock)
	customerRepo := new(mocks.CustomerRepositoryMock)
	sellerRepo := new(mocks.SellerRepositoryMock)
	jwtAuth := security.NewJWTAuth("secret", 24)

	uc := usecase.NewAuthUseCase(db, userRepo, customerRepo, sellerRepo, jwtAuth)

	req := &usecase.RegisterSellerReq{
		Email:            "seller@example.com",
		Password:         "password123",
		StoreName:        "Super Store",
		StoreDescription: "Best store ever",
	}

	userRepo.On("FindByEmail", req.Email).Return(nil, errors.New("not found"))

	sqlMock.ExpectBegin()
	sqlMock.ExpectCommit()

	userRepo.On("CreateWithTx", mock.Anything, mock.AnythingOfType("*entity.User")).Return(nil)
	sellerRepo.On("CreateWithTx", mock.Anything, mock.AnythingOfType("*entity.Seller")).Return(nil)

	token, err := uc.RegisterSeller(req)

	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	userRepo.AssertExpectations(t)
	sellerRepo.AssertExpectations(t)
	assert.NoError(t, sqlMock.ExpectationsWereMet())
}
