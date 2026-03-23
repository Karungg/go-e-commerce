package http_test

import (
	"bytes"
	"encoding/json"
	delivery "go-e-commerce/internal/delivery/http"
	"go-e-commerce/internal/mocks"
	"go-e-commerce/internal/pkg/apperror"
	"go-e-commerce/internal/usecase"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupRouter(authUsecase usecase.AuthUseCase) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	api := router.Group("/api")
	delivery.NewAuthController(api, authUsecase)
	return router
}

func TestRegisterCustomer_Success(t *testing.T) {
	mockUsecase := new(mocks.AuthUseCaseMock)
	mockUsecase.On("RegisterCustomer", mock.Anything, mock.AnythingOfType("*usecase.RegisterCustomerReq")).Return("mock.jwt.token", nil)

	router := setupRouter(mockUsecase)

	reqBody := usecase.RegisterCustomerReq{
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
	mockUsecase.On("RegisterCustomer", mock.Anything, mock.AnythingOfType("*usecase.RegisterCustomerReq")).Return("", apperror.ErrEmailConflict)

	router := setupRouter(mockUsecase)

	reqBody := usecase.RegisterCustomerReq{
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
	mockUsecase.On("RegisterSeller", mock.Anything, mock.AnythingOfType("*usecase.RegisterSellerReq")).Return("mock.jwt.token", nil)

	router := setupRouter(mockUsecase)

	reqBody := usecase.RegisterSellerReq{
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
