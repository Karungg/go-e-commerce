package product_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	productCtrl "go-e-commerce/internal/delivery/http/product"
	productDTO "go-e-commerce/internal/dto/product"
	"go-e-commerce/internal/pkg/apperror"
	productMock "go-e-commerce/internal/mocks/product"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupProductRouter(mockUseCase *productMock.ProductUseCaseMock) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	api := router.Group("/api")

	controller := productCtrl.NewProductController(mockUseCase)
	products := api.Group("/products")
	{
		products.POST("", controller.Create)
		products.GET("", controller.GetAll)
		products.GET("/:id", controller.GetByID)
		products.PUT("/:id", controller.Update)
		products.DELETE("/:id", controller.Delete)
	}

	return router
}

func TestCreateProduct_Success(t *testing.T) {
	mockUC := new(productMock.ProductUseCaseMock)

	id := uuid.New()
	mockUC.On("CreateProduct", mock.Anything, mock.AnythingOfType("*product.CreateProductReq")).Return(&productDTO.ProductRes{
		ID:          id,
		Title:       "Phone",
		Description: "Good Phone",
		Price:       1000,
		Stock:       10,
	}, nil)

	router := setupProductRouter(mockUC)

	reqBody := productDTO.CreateProductReq{
		Title:       "Phone",
		Description: "Good Phone",
		Price:       1000,
		Stock:       10,
		Image:       "img.jpg",
		CategoryID:  uuid.New().String(),
		SKU:         "SKU-11",
	}
	body, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest(http.MethodPost, "/api/products", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var res map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &res)
	assert.Equal(t, "success", res["status"])

	mockUC.AssertExpectations(t)
}

func TestGetAllProducts_Success(t *testing.T) {
	mockUC := new(productMock.ProductUseCaseMock)

	mockUC.On("GetAllProducts", mock.Anything).Return([]*productDTO.ProductRes{
		{ID: uuid.New(), Title: "P1", Price: 100},
	}, nil)

	router := setupProductRouter(mockUC)

	req, _ := http.NewRequest(http.MethodGet, "/api/products", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockUC.AssertExpectations(t)
}

func TestGetProductByID_Success(t *testing.T) {
	mockUC := new(productMock.ProductUseCaseMock)
	id := uuid.New()

	mockUC.On("GetProductByID", mock.Anything, id.String()).Return(&productDTO.ProductRes{
		ID: id, Title: "P1", Price: 100,
	}, nil)

	router := setupProductRouter(mockUC)

	req, _ := http.NewRequest(http.MethodGet, "/api/products/"+id.String(), nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockUC.AssertExpectations(t)
}

func TestGetProductByID_NotFound(t *testing.T) {
	mockUC := new(productMock.ProductUseCaseMock)
	id := uuid.New()

	mockUC.On("GetProductByID", mock.Anything, id.String()).Return(nil, apperror.ErrProductNotFound)

	router := setupProductRouter(mockUC)

	req, _ := http.NewRequest(http.MethodGet, "/api/products/"+id.String(), nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	mockUC.AssertExpectations(t)
}

func TestDeleteProduct_Success(t *testing.T) {
	mockUC := new(productMock.ProductUseCaseMock)
	id := uuid.New()

	mockUC.On("DeleteProduct", mock.Anything, id.String()).Return(nil)

	router := setupProductRouter(mockUC)

	req, _ := http.NewRequest(http.MethodDelete, "/api/products/"+id.String(), nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockUC.AssertExpectations(t)
}
