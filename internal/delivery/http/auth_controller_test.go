package http_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	delivery "go-e-commerce/internal/delivery/http"
	"go-e-commerce/internal/dto"
	"go-e-commerce/internal/mocks"
	"go-e-commerce/internal/pkg/apperror"
	"go-e-commerce/internal/port"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupRouter(authUsecase port.AuthUseCase) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	api := router.Group("/api")
	
	authController := delivery.NewAuthController(authUsecase)
	auth := api.Group("/auth")
	{
		auth.POST("/register/customer", authController.RegisterCustomer)
		auth.POST("/register/seller", authController.RegisterSeller)
		auth.POST("/login", authController.Login)
		auth.POST("/logout", authController.Logout)
	}
	
	return router
}

func TestRegisterCustomer_Success(t *testing.T) {
	mockUsecase := new(mocks.AuthUseCaseMock)
	mockUsecase.On("RegisterCustomer", mock.Anything, mock.AnythingOfType("*dto.RegisterCustomerReq")).Return("mock.jwt.token", nil)

	router := setupRouter(mockUsecase)

	reqBody := dto.RegisterCustomerReq{
		Email:     "cust@example.com",
		Password:  "password",
		FirstName: "John",
		LastName:  "Doe",
	}
	body, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest(http.MethodPost, "/api/auth/register/customer", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	
	var res map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &res)
	
	assert.Equal(t, "success", res["status"])
	assert.Equal(t, "customer registered successfully", res["message"])
	
	data := res["data"].(map[string]interface{})
	assert.Equal(t, "mock.jwt.token", data["token"])

	mockUsecase.AssertExpectations(t)
}

func TestRegisterCustomer_InvalidJSON(t *testing.T) {
	mockUsecase := new(mocks.AuthUseCaseMock)
	router := setupRouter(mockUsecase)

	req, _ := http.NewRequest(http.MethodPost, "/api/auth/register/customer", bytes.NewBuffer([]byte("{invalid json")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var res map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &res)
	
	assert.Equal(t, "error", res["status"])
	assert.Equal(t, "Invalid request payload", res["message"])
	assert.NotNil(t, res["errors"])
}

func TestRegisterCustomer_UsecaseErrorConflict(t *testing.T) {
	mockUsecase := new(mocks.AuthUseCaseMock)
	// Inject the strictly handled Domain Error dictating 409 boundaries natively
	mockUsecase.On("RegisterCustomer", mock.Anything, mock.AnythingOfType("*dto.RegisterCustomerReq")).Return("", apperror.ErrEmailConflict)

	router := setupRouter(mockUsecase)

	reqBody := dto.RegisterCustomerReq{
		Email:     "cust@example.com",
		Password:  "password",
		FirstName: "John",
		LastName:  "Doe",
	}
	body, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest(http.MethodPost, "/api/auth/register/customer", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusConflict, w.Code)
	
	var res map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &res)
	
	assert.Equal(t, "error", res["status"])
	assert.Equal(t, "Registration failed", res["message"])
	assert.Equal(t, apperror.ErrEmailConflict.Message, res["errors"])
	mockUsecase.AssertExpectations(t)
}

func TestRegisterSeller_Success(t *testing.T) {
	mockUsecase := new(mocks.AuthUseCaseMock)
	mockUsecase.On("RegisterSeller", mock.Anything, mock.AnythingOfType("*dto.RegisterSellerReq")).Return("mock.jwt.token", nil)

	router := setupRouter(mockUsecase)

	reqBody := dto.RegisterSellerReq{
		Email:            "seller@example.com",
		Password:         "password",
		StoreName:        "Awesome Store",
		StoreDescription: "Selling things",
	}
	body, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest(http.MethodPost, "/api/auth/register/seller", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	
	var res map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &res)
	
	assert.Equal(t, "success", res["status"])
	assert.Equal(t, "seller registered successfully", res["message"])
	
	data := res["data"].(map[string]interface{})
	assert.Equal(t, "mock.jwt.token", data["token"])

	mockUsecase.AssertExpectations(t)
}

func TestLogin_Success(t *testing.T) {
	mockUsecase := new(mocks.AuthUseCaseMock)
	mockUsecase.On("Login", mock.Anything, mock.AnythingOfType("*dto.LoginReq")).Return("mock.login.token", nil)

	router := setupRouter(mockUsecase)

	reqBody := dto.LoginReq{
		Email:    "test@example.com",
		Password: "password123",
	}
	body, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	
	var res map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &res)
	
	assert.Equal(t, "success", res["status"])
	assert.Equal(t, "login successful", res["message"])
	
	data := res["data"].(map[string]interface{})
	assert.Equal(t, "mock.login.token", data["token"])

	mockUsecase.AssertExpectations(t)
}

func TestLogin_InvalidJSON(t *testing.T) {
	mockUsecase := new(mocks.AuthUseCaseMock)
	router := setupRouter(mockUsecase)

	req, _ := http.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewBuffer([]byte(`{invalid json`)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var res map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &res)
	
	assert.Equal(t, "error", res["status"])
	assert.Equal(t, "Invalid request payload", res["message"])
	assert.NotNil(t, res["errors"])
}

func TestLogin_UsecaseError(t *testing.T) {
	mockUsecase := new(mocks.AuthUseCaseMock)
	mockUsecase.On("Login", mock.Anything, mock.AnythingOfType("*dto.LoginReq")).Return("", apperror.ErrInvalidPassword)

	router := setupRouter(mockUsecase)

	reqBody := dto.LoginReq{
		Email:    "test@example.com",
		Password: "wrongpassword",
	}
	body, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	
	var res map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &res)
	
	assert.Equal(t, "error", res["status"])
	assert.Equal(t, "Login failed", res["message"])
	assert.Equal(t, apperror.ErrInvalidPassword.Message, res["errors"])

	mockUsecase.AssertExpectations(t)
}

func TestLogout_Success(t *testing.T) {
	mockUsecase := new(mocks.AuthUseCaseMock)
	mockUsecase.On("Logout", mock.Anything).Return(nil)

	router := setupRouter(mockUsecase)

	req, _ := http.NewRequest(http.MethodPost, "/api/auth/logout", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	
	var res map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &res)
	
	assert.Equal(t, "success", res["status"])
	assert.Equal(t, "logout successful", res["message"])
	
	mockUsecase.AssertExpectations(t)
}

func TestLogout_UsecaseError(t *testing.T) {
	mockUsecase := new(mocks.AuthUseCaseMock)
	mockUsecase.On("Logout", mock.Anything).Return(errors.New("db disconnect"))

	router := setupRouter(mockUsecase)

	req, _ := http.NewRequest(http.MethodPost, "/api/auth/logout", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	
	mockUsecase.AssertExpectations(t)
}
