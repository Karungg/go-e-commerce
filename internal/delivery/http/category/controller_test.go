package category_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	categoryCtrl "go-e-commerce/internal/delivery/http/category"
	categoryDTO "go-e-commerce/internal/dto/category"
	"go-e-commerce/internal/pkg/apperror"
	categoryMock "go-e-commerce/internal/mocks/category"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupCategoryRouter(mockUseCase *categoryMock.CategoryUseCaseMock) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	api := router.Group("/api")

	controller := categoryCtrl.NewCategoryController(mockUseCase)
	categories := api.Group("/categories")
	{
		categories.POST("", controller.Create)
		categories.GET("", controller.GetAll)
		categories.GET("/:id", controller.GetByID)
		categories.PUT("/:id", controller.Update)
		categories.DELETE("/:id", controller.Delete)
	}

	return router
}

func TestCreateCategory_Success(t *testing.T) {
	mockUC := new(categoryMock.CategoryUseCaseMock)

	id := uuid.New()
	mockUC.On("CreateCategory", mock.Anything, mock.AnythingOfType("*category.CreateCategoryReq")).Return(&categoryDTO.CategoryRes{
		ID:          id,
		Title:       "Electronics",
		Description: "Gadgets",
	}, nil)

	router := setupCategoryRouter(mockUC)

	reqBody := categoryDTO.CreateCategoryReq{
		Title:       "Electronics",
		Description: "Gadgets",
	}
	body, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest(http.MethodPost, "/api/categories", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var res map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &res)
	assert.Equal(t, "success", res["status"])

	mockUC.AssertExpectations(t)
}

func TestGetAllCategories_Success(t *testing.T) {
	mockUC := new(categoryMock.CategoryUseCaseMock)

	mockUC.On("GetAllCategories", mock.Anything).Return([]*categoryDTO.CategoryRes{
		{ID: uuid.New(), Title: "C1", Description: "D1"},
	}, nil)

	router := setupCategoryRouter(mockUC)

	req, _ := http.NewRequest(http.MethodGet, "/api/categories", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockUC.AssertExpectations(t)
}

func TestGetCategoryByID_Success(t *testing.T) {
	mockUC := new(categoryMock.CategoryUseCaseMock)
	id := uuid.New()

	mockUC.On("GetCategoryByID", mock.Anything, id.String()).Return(&categoryDTO.CategoryRes{
		ID: id, Title: "C1", Description: "D1",
	}, nil)

	router := setupCategoryRouter(mockUC)

	req, _ := http.NewRequest(http.MethodGet, "/api/categories/"+id.String(), nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockUC.AssertExpectations(t)
}

func TestGetCategoryByID_NotFound(t *testing.T) {
	mockUC := new(categoryMock.CategoryUseCaseMock)
	id := uuid.New()

	mockUC.On("GetCategoryByID", mock.Anything, id.String()).Return(nil, apperror.ErrCategoryNotFound)

	router := setupCategoryRouter(mockUC)

	req, _ := http.NewRequest(http.MethodGet, "/api/categories/"+id.String(), nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	mockUC.AssertExpectations(t)
}

func TestDeleteCategory_Success(t *testing.T) {
	mockUC := new(categoryMock.CategoryUseCaseMock)
	id := uuid.New()

	mockUC.On("DeleteCategory", mock.Anything, id.String()).Return(nil)

	router := setupCategoryRouter(mockUC)

	req, _ := http.NewRequest(http.MethodDelete, "/api/categories/"+id.String(), nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockUC.AssertExpectations(t)
}
