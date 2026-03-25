package http_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	delivery "go-e-commerce/internal/delivery/http"
	"go-e-commerce/internal/dto"
	"go-e-commerce/internal/port"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// CategoryUseCaseMock represents a mocked port.CategoryUseCase
type CategoryUseCaseMock struct {
	mock.Mock
}

func (m *CategoryUseCaseMock) CreateCategory(ctx context.Context, req *dto.CreateCategoryReq) (*dto.CategoryRes, error) {
	args := m.Called(ctx, req)
	if args.Get(0) != nil {
		return args.Get(0).(*dto.CategoryRes), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *CategoryUseCaseMock) GetAllCategories(ctx context.Context) ([]*dto.CategoryRes, error) {
	args := m.Called(ctx)
	if args.Get(0) != nil {
		return args.Get(0).([]*dto.CategoryRes), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *CategoryUseCaseMock) GetCategoryByID(ctx context.Context, id string) (*dto.CategoryRes, error) {
	args := m.Called(ctx, id)
	if args.Get(0) != nil {
		return args.Get(0).(*dto.CategoryRes), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *CategoryUseCaseMock) UpdateCategory(ctx context.Context, id string, req *dto.UpdateCategoryReq) (*dto.CategoryRes, error) {
	args := m.Called(ctx, id, req)
	if args.Get(0) != nil {
		return args.Get(0).(*dto.CategoryRes), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *CategoryUseCaseMock) DeleteCategory(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func setupCategoryRouter(uc port.CategoryUseCase) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	ctrl := delivery.NewCategoryController(uc)

	router.POST("/categories", ctrl.Create)
	router.GET("/categories", ctrl.GetAll)
	router.GET("/categories/:id", ctrl.GetByID)
	router.PUT("/categories/:id", ctrl.Update)
	router.DELETE("/categories/:id", ctrl.Delete)

	return router
}

func TestCategoryController_GetAll(t *testing.T) {
	ucMock := new(CategoryUseCaseMock)
	router := setupCategoryRouter(ucMock)

	mockData := []*dto.CategoryRes{
		{ID: uuid.New(), Title: "Cat1", Description: "Desc1"},
	}

	ucMock.On("GetAllCategories", mock.Anything).Return(mockData, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/categories", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	ucMock.AssertExpectations(t)
}

func TestCategoryController_Create(t *testing.T) {
	ucMock := new(CategoryUseCaseMock)
	router := setupCategoryRouter(ucMock)

	reqDto := &dto.CreateCategoryReq{Title: "New", Description: "Desc"}
	resDto := &dto.CategoryRes{ID: uuid.New(), Title: "New", Description: "Desc"}

	ucMock.On("CreateCategory", mock.Anything, reqDto).Return(resDto, nil)

	body, _ := json.Marshal(reqDto)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/categories", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	ucMock.AssertExpectations(t)
}
