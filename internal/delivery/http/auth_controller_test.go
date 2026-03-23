package http_test

import (
	"bytes"
	"encoding/json"
	"errors"
	delivery "go-e-commerce/internal/delivery/http"
	"go-e-commerce/internal/mocks"
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
	mockUsecase.On("RegisterCustomer", mock.AnythingOfType("*usecase.RegisterCustomerReq")).Return("mock.jwt.token", nil)

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
	
	assert.Equal(t, "customer registered successfully", res["message"])
	assert.Equal(t, "mock.jwt.token", res["token"])

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
}

func TestRegisterCustomer_UsecaseError(t *testing.T) {
	mockUsecase := new(mocks.AuthUseCaseMock)
	mockUsecase.On("RegisterCustomer", mock.AnythingOfType("*usecase.RegisterCustomerReq")).Return("", errors.New("email is already registered"))

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

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	
	var res map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &res)
	
	assert.Equal(t, "email is already registered", res["error"])
	mockUsecase.AssertExpectations(t)
}

func TestRegisterSeller_Success(t *testing.T) {
	mockUsecase := new(mocks.AuthUseCaseMock)
	mockUsecase.On("RegisterSeller", mock.AnythingOfType("*usecase.RegisterSellerReq")).Return("mock.jwt.token", nil)

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
	
	assert.Equal(t, "seller registered successfully", res["message"])
	assert.Equal(t, "mock.jwt.token", res["token"])

	mockUsecase.AssertExpectations(t)
}
