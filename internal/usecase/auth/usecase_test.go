package auth_test

import (
	"context"
	"io"
	"log/slog"
	"testing"

	"go-e-commerce/internal/entity"
	authMock "go-e-commerce/internal/mocks/auth"
	"go-e-commerce/internal/pkg/apperror"
	authDTO "go-e-commerce/internal/dto/auth"
	"go-e-commerce/internal/repository"
	"go-e-commerce/internal/security"
	authUC "go-e-commerce/internal/usecase/auth"
	"golang.org/x/crypto/bcrypt"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
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
	userRepo := new(authMock.UserRepositoryMock)
	customerRepo := new(authMock.CustomerRepositoryMock)
	sellerRepo := new(authMock.SellerRepositoryMock)
	jwtAuth := security.NewJWTAuth("secret", 24)

	txManager := repository.NewTransactionManager(db)
	uc := authUC.NewAuthUseCase(txManager, getDiscardLogger(), userRepo, customerRepo, sellerRepo, jwtAuth)

	req := &authDTO.RegisterCustomerReq{
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

	userRepo.On("Create", mock.Anything, mock.AnythingOfType("*entity.User")).Return(nil)
	customerRepo.On("Create", mock.Anything, mock.AnythingOfType("*entity.Customer")).Return(nil)

	token, err := uc.RegisterCustomer(context.Background(), req)

	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	userRepo.AssertExpectations(t)
	customerRepo.AssertExpectations(t)
	assert.NoError(t, sqlMock.ExpectationsWereMet())
}

func TestRegisterCustomer_EmailExists(t *testing.T) {
	db, _ := setupTestDB(t)
	userRepo := new(authMock.UserRepositoryMock)
	customerRepo := new(authMock.CustomerRepositoryMock)
	sellerRepo := new(authMock.SellerRepositoryMock)
	jwtAuth := security.NewJWTAuth("secret", 24)

	txManager := repository.NewTransactionManager(db)
	uc := authUC.NewAuthUseCase(txManager, getDiscardLogger(), userRepo, customerRepo, sellerRepo, jwtAuth)

	req := &authDTO.RegisterCustomerReq{
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
	userRepo := new(authMock.UserRepositoryMock)
	customerRepo := new(authMock.CustomerRepositoryMock)
	sellerRepo := new(authMock.SellerRepositoryMock)
	jwtAuth := security.NewJWTAuth("secret", 24)

	txManager := repository.NewTransactionManager(db)
	uc := authUC.NewAuthUseCase(txManager, getDiscardLogger(), userRepo, customerRepo, sellerRepo, jwtAuth)

	req := &authDTO.RegisterCustomerReq{
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
	userRepo := new(authMock.UserRepositoryMock)
	customerRepo := new(authMock.CustomerRepositoryMock)
	sellerRepo := new(authMock.SellerRepositoryMock)
	jwtAuth := security.NewJWTAuth("secret", 24)

	txManager := repository.NewTransactionManager(db)
	uc := authUC.NewAuthUseCase(txManager, getDiscardLogger(), userRepo, customerRepo, sellerRepo, jwtAuth)

	req := &authDTO.RegisterSellerReq{
		Email:            "seller@example.com",
		Password:         "password123",
		StoreName:        "Super Store",
		StoreDescription: "Best store ever",
	}

	userRepo.On("FindByEmail", mock.Anything, req.Email).Return(nil, nil)
	sellerRepo.On("FindByStoreName", mock.Anything, req.StoreName).Return(nil, nil)

	sqlMock.ExpectBegin()
	sqlMock.ExpectCommit()

	userRepo.On("Create", mock.Anything, mock.AnythingOfType("*entity.User")).Return(nil)
	sellerRepo.On("Create", mock.Anything, mock.AnythingOfType("*entity.Seller")).Return(nil)

	token, err := uc.RegisterSeller(context.Background(), req)

	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	userRepo.AssertExpectations(t)
	sellerRepo.AssertExpectations(t)
	assert.NoError(t, sqlMock.ExpectationsWereMet())
}

func TestRegisterSeller_StoreNameExists(t *testing.T) {
	db, _ := setupTestDB(t)
	userRepo := new(authMock.UserRepositoryMock)
	customerRepo := new(authMock.CustomerRepositoryMock)
	sellerRepo := new(authMock.SellerRepositoryMock)
	jwtAuth := security.NewJWTAuth("secret", 24)

	txManager := repository.NewTransactionManager(db)
	uc := authUC.NewAuthUseCase(txManager, getDiscardLogger(), userRepo, customerRepo, sellerRepo, jwtAuth)

	req := &authDTO.RegisterSellerReq{
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

func TestLogin_Success(t *testing.T) {
	db, _ := setupTestDB(t)
	userRepo := new(authMock.UserRepositoryMock)
	customerRepo := new(authMock.CustomerRepositoryMock)
	sellerRepo := new(authMock.SellerRepositoryMock)
	jwtAuth := security.NewJWTAuth("secret", 24)

	txManager := repository.NewTransactionManager(db)
	uc := authUC.NewAuthUseCase(txManager, getDiscardLogger(), userRepo, customerRepo, sellerRepo, jwtAuth)

	req := &authDTO.LoginReq{
		Email:    "test@example.com",
		Password: "password123",
	}

	hashedBytes, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.MinCost)
	hashedPassword := string(hashedBytes)

	mockUser := &entity.User{
		ID:       uuid.New(), 
		Email:    req.Email,
		Password: hashedPassword,
		Role:     entity.RoleCustomer,
	}

	userRepo.On("FindByEmail", mock.Anything, req.Email).Return(mockUser, nil)

	token, err := uc.Login(context.Background(), req)

	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	userRepo.AssertExpectations(t)
}

func TestLogin_UserNotFound(t *testing.T) {
	db, _ := setupTestDB(t)
	userRepo := new(authMock.UserRepositoryMock)
	customerRepo := new(authMock.CustomerRepositoryMock)
	sellerRepo := new(authMock.SellerRepositoryMock)
	jwtAuth := security.NewJWTAuth("secret", 24)

	txManager := repository.NewTransactionManager(db)
	uc := authUC.NewAuthUseCase(txManager, getDiscardLogger(), userRepo, customerRepo, sellerRepo, jwtAuth)

	req := &authDTO.LoginReq{
		Email:    "notfound@example.com",
		Password: "password123",
	}

	userRepo.On("FindByEmail", mock.Anything, req.Email).Return(nil, nil)

	token, err := uc.Login(context.Background(), req)

	assert.Error(t, err)
	assert.ErrorIs(t, err, apperror.ErrInvalidPassword)
	assert.Empty(t, token)

	userRepo.AssertExpectations(t)
}

func TestLogin_InvalidPassword(t *testing.T) {
	db, _ := setupTestDB(t)
	userRepo := new(authMock.UserRepositoryMock)
	customerRepo := new(authMock.CustomerRepositoryMock)
	sellerRepo := new(authMock.SellerRepositoryMock)
	jwtAuth := security.NewJWTAuth("secret", 24)

	txManager := repository.NewTransactionManager(db)
	uc := authUC.NewAuthUseCase(txManager, getDiscardLogger(), userRepo, customerRepo, sellerRepo, jwtAuth)

	req := &authDTO.LoginReq{
		Email:    "test@example.com",
		Password: "wrongpassword",
	}

	hashedBytes, _ := bcrypt.GenerateFromPassword([]byte("correctpassword"), bcrypt.MinCost)
	hashedPassword := string(hashedBytes)

	mockUser := &entity.User{
		Email:    req.Email,
		Password: hashedPassword,
		Role:     entity.RoleCustomer,
	}

	userRepo.On("FindByEmail", mock.Anything, req.Email).Return(mockUser, nil)

	token, err := uc.Login(context.Background(), req)

	assert.Error(t, err)
	assert.ErrorIs(t, err, apperror.ErrInvalidPassword)
	assert.Empty(t, token)

	userRepo.AssertExpectations(t)
}

func TestUpdateCustomer_Success(t *testing.T) {
	db, _ := setupTestDB(t)
	userRepo := new(authMock.UserRepositoryMock)
	customerRepo := new(authMock.CustomerRepositoryMock)
	sellerRepo := new(authMock.SellerRepositoryMock)
	jwtAuth := security.NewJWTAuth("secret", 24)

	txManager := repository.NewTransactionManager(db)
	uc := authUC.NewAuthUseCase(txManager, getDiscardLogger(), userRepo, customerRepo, sellerRepo, jwtAuth)

	userID := uuid.New()
	customerID := uuid.New()
	req := &authDTO.UpdateCustomerReq{
		FirstName: "John Updated",
		LastName:  "Doe Updated",
		Phone:     "0987654321",
		Address:   "456 New Street",
	}

	existingCustomer := &entity.Customer{
		ID:        customerID,
		UserID:    userID,
		FirstName: "John",
		LastName:  "Doe",
		Phone:     "1234567890",
		Address:   "123 Old Street",
	}

	customerRepo.On("FindByUserID", mock.Anything, userID).Return(existingCustomer, nil)
	customerRepo.On("FindByPhone", mock.Anything, req.Phone).Return(nil, nil)
	customerRepo.On("Update", mock.Anything, mock.AnythingOfType("*entity.Customer")).Return(nil)

	err := uc.UpdateCustomer(context.Background(), userID, req)

	assert.NoError(t, err)
	customerRepo.AssertExpectations(t)
}

func TestUpdateCustomer_UserNotFound(t *testing.T) {
	db, _ := setupTestDB(t)
	userRepo := new(authMock.UserRepositoryMock)
	customerRepo := new(authMock.CustomerRepositoryMock)
	sellerRepo := new(authMock.SellerRepositoryMock)
	jwtAuth := security.NewJWTAuth("secret", 24)

	txManager := repository.NewTransactionManager(db)
	uc := authUC.NewAuthUseCase(txManager, getDiscardLogger(), userRepo, customerRepo, sellerRepo, jwtAuth)

	userID := uuid.New()
	req := &authDTO.UpdateCustomerReq{
		FirstName: "John Updated",
	}

	customerRepo.On("FindByUserID", mock.Anything, userID).Return(nil, nil)

	err := uc.UpdateCustomer(context.Background(), userID, req)

	assert.Error(t, err)
	assert.ErrorIs(t, err, apperror.ErrUserNotFound)
	customerRepo.AssertExpectations(t)
}

func TestUpdateCustomer_PhoneConflict(t *testing.T) {
	db, _ := setupTestDB(t)
	userRepo := new(authMock.UserRepositoryMock)
	customerRepo := new(authMock.CustomerRepositoryMock)
	sellerRepo := new(authMock.SellerRepositoryMock)
	jwtAuth := security.NewJWTAuth("secret", 24)

	txManager := repository.NewTransactionManager(db)
	uc := authUC.NewAuthUseCase(txManager, getDiscardLogger(), userRepo, customerRepo, sellerRepo, jwtAuth)

	userID := uuid.New()
	customerID := uuid.New()
	anotherCustomerID := uuid.New()
	req := &authDTO.UpdateCustomerReq{
		Phone: "0987654321",
	}

	existingCustomer := &entity.Customer{
		ID:        customerID,
		UserID:    userID,
		Phone:     "1234567890",
	}

	anotherCustomer := &entity.Customer{
		ID:    anotherCustomerID,
		Phone: "0987654321",
	}

	customerRepo.On("FindByUserID", mock.Anything, userID).Return(existingCustomer, nil)
	customerRepo.On("FindByPhone", mock.Anything, req.Phone).Return(anotherCustomer, nil)

	err := uc.UpdateCustomer(context.Background(), userID, req)

	assert.Error(t, err)
	assert.ErrorIs(t, err, apperror.ErrPhoneConflict)
	customerRepo.AssertExpectations(t)
}
