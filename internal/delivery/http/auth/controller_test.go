package auth_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	authCtrl "go-e-commerce/internal/delivery/http/auth"
	authDTO "go-e-commerce/internal/dto/auth"
	authMock "go-e-commerce/internal/mocks/auth"
	"go-e-commerce/internal/pkg/apperror"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupRouter(authUsecase *authMock.AuthUseCaseMock) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	api := router.Group("/api")

	authController := authCtrl.NewAuthController(authUsecase)
	auth := api.Group("/auth")
	{
		auth.POST("/register/customer", authController.RegisterCustomer)
		auth.POST("/register/seller", authController.RegisterSeller)
		auth.POST("/login", authController.Login)
		auth.POST("/logout", authController.Logout)
		
		// Map UpdateCustomer to PUT profile
		auth.PUT("/customer", func(c *gin.Context) {
			// Mocking the authorization middleware setting userID
			testUserID, _ := c.GetQuery("test_user_id")
			if testUserID != "" {
				c.Set("userID", testUserID)
			}
			authController.UpdateCustomer(c)
		})
	}

	return router
}

func TestRegisterCustomer_Success(t *testing.T) {
	mockUsecase := new(authMock.AuthUseCaseMock)
	mockUsecase.On("RegisterCustomer", mock.Anything, mock.AnythingOfType("*auth.RegisterCustomerReq")).Return("mock.jwt.token", nil)

	router := setupRouter(mockUsecase)

	reqBody := authDTO.RegisterCustomerReq{
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
	mockUsecase := new(authMock.AuthUseCaseMock)
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
	mockUsecase := new(authMock.AuthUseCaseMock)
	mockUsecase.On("RegisterCustomer", mock.Anything, mock.AnythingOfType("*auth.RegisterCustomerReq")).Return("", apperror.ErrEmailConflict)

	router := setupRouter(mockUsecase)

	reqBody := authDTO.RegisterCustomerReq{
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
	mockUsecase := new(authMock.AuthUseCaseMock)
	mockUsecase.On("RegisterSeller", mock.Anything, mock.AnythingOfType("*auth.RegisterSellerReq")).Return("mock.jwt.token", nil)

	router := setupRouter(mockUsecase)

	reqBody := authDTO.RegisterSellerReq{
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
	mockUsecase := new(authMock.AuthUseCaseMock)
	mockUsecase.On("Login", mock.Anything, mock.AnythingOfType("*auth.LoginReq")).Return("mock.login.token", nil)

	router := setupRouter(mockUsecase)

	reqBody := authDTO.LoginReq{
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
	mockUsecase := new(authMock.AuthUseCaseMock)
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
	mockUsecase := new(authMock.AuthUseCaseMock)
	mockUsecase.On("Login", mock.Anything, mock.AnythingOfType("*auth.LoginReq")).Return("", apperror.ErrInvalidPassword)

	router := setupRouter(mockUsecase)

	reqBody := authDTO.LoginReq{
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
	mockUsecase := new(authMock.AuthUseCaseMock)
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
	mockUsecase := new(authMock.AuthUseCaseMock)
	mockUsecase.On("Logout", mock.Anything).Return(errors.New("db disconnect"))

	router := setupRouter(mockUsecase)

	req, _ := http.NewRequest(http.MethodPost, "/api/auth/logout", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	mockUsecase.AssertExpectations(t)
}

func TestUpdateCustomer_Success(t *testing.T) {
	mockUsecase := new(authMock.AuthUseCaseMock)
	mockUsecase.On("UpdateCustomer", mock.Anything, mock.AnythingOfType("uuid.UUID"), mock.AnythingOfType("*auth.UpdateCustomerReq")).Return(nil)

	router := setupRouter(mockUsecase)

	reqBody := authDTO.UpdateCustomerReq{
		Email:     "test@example.com",
		FirstName: "John",
		LastName:  "Doe",
	}
	body, _ := json.Marshal(reqBody)

	userID := uuid.New().String()
	req, _ := http.NewRequest(http.MethodPut, "/api/auth/customer?test_user_id="+userID, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var res map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &res)

	assert.Equal(t, "success", res["status"])
	assert.Equal(t, "customer profile updated successfully", res["message"])

	mockUsecase.AssertExpectations(t)
}

func TestUpdateCustomer_UserNotFound(t *testing.T) {
	mockUsecase := new(authMock.AuthUseCaseMock)
	mockUsecase.On("UpdateCustomer", mock.Anything, mock.AnythingOfType("uuid.UUID"), mock.AnythingOfType("*auth.UpdateCustomerReq")).Return(apperror.ErrUserNotFound)

	router := setupRouter(mockUsecase)

	reqBody := authDTO.UpdateCustomerReq{
		Email:     "test@example.com",
		FirstName: "John",
		LastName:  "Doe",
	}
	body, _ := json.Marshal(reqBody)

	userID := uuid.New().String()
	req, _ := http.NewRequest(http.MethodPut, "/api/auth/customer?test_user_id="+userID, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)

	var res map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &res)

	assert.Equal(t, "error", res["status"])
	assert.Equal(t, "Update failed", res["message"])
	assert.Equal(t, apperror.ErrUserNotFound.Message, res["errors"])

	mockUsecase.AssertExpectations(t)
}

