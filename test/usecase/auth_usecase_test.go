package usecase_test

import (
	"context"
	"io"
	"log/slog"
	"testing"

	"go-e-commerce/internal/dto"
	"go-e-commerce/internal/entity"
	"go-e-commerce/internal/mocks"
	"go-e-commerce/internal/pkg/apperror"
	"go-e-commerce/internal/security"
	"go-e-commerce/internal/usecase"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
	sqlDB, smock, err := sqlmock.New()
	assert.NoError(t, err)

	grmDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: sqlDB,
	}), &gorm.Config{
		SkipDefaultTransaction: true,
	})
	assert.NoError(t, err)

	return grmDB, smock
}

func getDiscardLogger() *slog.Logger {
	return slog.New(slog.NewJSONHandler(io.Discard, nil))
}

func TestRegisterCustomer_Success(t *testing.T) {
	db, sqlMock := setupTestDB(t)
	userRepo := new(mocks.UserRepositoryMock)
	customerRepo := new(mocks.CustomerRepositoryMock)
	sellerRepo := new(mocks.SellerRepositoryMock)
	jwtAuth := security.NewJWTAuth("secret", 24)

	uc := usecase.NewAuthUseCase(db, getDiscardLogger(), userRepo, customerRepo, sellerRepo, jwtAuth)

	req := &dto.RegisterCustomerReq{
		Email:     "test@example.com",
		Password:  "password123",
		FirstName: "John",
		LastName:  "Doe",
		Phone:     "123456789",
		Address:   "123 Street",
	}

	userRepo.On("FindByEmail", mock.Anything, req.Email).Return(nil, nil)
	customerRepo.On("FindByPhone", mock.Anything, req.Phone).Return(nil, nil)

	sqlMock.ExpectBegin()
	sqlMock.ExpectCommit()

	userRepo.On("CreateWithTx", mock.Anything, mock.Anything, mock.AnythingOfType("*entity.User")).Return(nil)
	customerRepo.On("CreateWithTx", mock.Anything, mock.Anything, mock.AnythingOfType("*entity.Customer")).Return(nil)

	token, err := uc.RegisterCustomer(context.Background(), req)

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

	uc := usecase.NewAuthUseCase(db, getDiscardLogger(), userRepo, customerRepo, sellerRepo, jwtAuth)

	req := &dto.RegisterCustomerReq{
		Email: "test@example.com",
	}

	existingUser := &entity.User{Email: req.Email}
	userRepo.On("FindByEmail", mock.Anything, req.Email).Return(existingUser, nil)

	token, err := uc.RegisterCustomer(context.Background(), req)

	assert.Error(t, err)
	assert.ErrorIs(t, err, apperror.ErrEmailConflict)
	assert.Empty(t, token)
}

func TestRegisterCustomer_PhoneExists(t *testing.T) {
	db, _ := setupTestDB(t)
	userRepo := new(mocks.UserRepositoryMock)
	customerRepo := new(mocks.CustomerRepositoryMock)
	sellerRepo := new(mocks.SellerRepositoryMock)
	jwtAuth := security.NewJWTAuth("secret", 24)

	uc := usecase.NewAuthUseCase(db, getDiscardLogger(), userRepo, customerRepo, sellerRepo, jwtAuth)

	req := &dto.RegisterCustomerReq{
		Email: "new@example.com",
		Phone: "123456789",
	}

	userRepo.On("FindByEmail", mock.Anything, req.Email).Return(nil, nil)
	
	existingCustomer := &entity.Customer{Phone: req.Phone}
	customerRepo.On("FindByPhone", mock.Anything, req.Phone).Return(existingCustomer, nil)

	token, err := uc.RegisterCustomer(context.Background(), req)

	assert.Error(t, err)
	assert.ErrorIs(t, err, apperror.ErrPhoneConflict)
	assert.Empty(t, token)
}

func TestRegisterSeller_Success(t *testing.T) {
	db, sqlMock := setupTestDB(t)
	userRepo := new(mocks.UserRepositoryMock)
	customerRepo := new(mocks.CustomerRepositoryMock)
	sellerRepo := new(mocks.SellerRepositoryMock)
	jwtAuth := security.NewJWTAuth("secret", 24)

	uc := usecase.NewAuthUseCase(db, getDiscardLogger(), userRepo, customerRepo, sellerRepo, jwtAuth)

	req := &dto.RegisterSellerReq{
		Email:            "seller@example.com",
		Password:         "password123",
		StoreName:        "Super Store",
		StoreDescription: "Best store ever",
	}

	userRepo.On("FindByEmail", mock.Anything, req.Email).Return(nil, nil)
	sellerRepo.On("FindByStoreName", mock.Anything, req.StoreName).Return(nil, nil)

	sqlMock.ExpectBegin()
	sqlMock.ExpectCommit()

	userRepo.On("CreateWithTx", mock.Anything, mock.Anything, mock.AnythingOfType("*entity.User")).Return(nil)
	sellerRepo.On("CreateWithTx", mock.Anything, mock.Anything, mock.AnythingOfType("*entity.Seller")).Return(nil)

	token, err := uc.RegisterSeller(context.Background(), req)

	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	userRepo.AssertExpectations(t)
	sellerRepo.AssertExpectations(t)
	assert.NoError(t, sqlMock.ExpectationsWereMet())
}

func TestRegisterSeller_StoreNameExists(t *testing.T) {
	db, _ := setupTestDB(t)
	userRepo := new(mocks.UserRepositoryMock)
	customerRepo := new(mocks.CustomerRepositoryMock)
	sellerRepo := new(mocks.SellerRepositoryMock)
	jwtAuth := security.NewJWTAuth("secret", 24)

	uc := usecase.NewAuthUseCase(db, getDiscardLogger(), userRepo, customerRepo, sellerRepo, jwtAuth)

	req := &dto.RegisterSellerReq{
		Email:     "new_seller@example.com",
		StoreName: "Super Store",
	}

	userRepo.On("FindByEmail", mock.Anything, req.Email).Return(nil, nil)

	existingSeller := &entity.Seller{StoreName: req.StoreName}
	sellerRepo.On("FindByStoreName", mock.Anything, req.StoreName).Return(existingSeller, nil)

	token, err := uc.RegisterSeller(context.Background(), req)

	assert.Error(t, err)
	assert.ErrorIs(t, err, apperror.ErrStoreNameConflict)
	assert.Empty(t, token)
}
