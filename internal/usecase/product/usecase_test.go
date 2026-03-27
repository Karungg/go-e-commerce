package product_test

import (
	"context"
	"errors"
	"testing"

	productDTO "go-e-commerce/internal/dto/product"
	"go-e-commerce/internal/entity"
	productMock "go-e-commerce/internal/mocks/product"
	"go-e-commerce/internal/pkg/apperror"
	productUC "go-e-commerce/internal/usecase/product"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

func TestCreateProduct_Success(t *testing.T) {
	mockRepo := new(productMock.ProductRepositoryMock)
	uc := productUC.NewProductUseCase(mockRepo)

	categoryID := uuid.New()
	req := &productDTO.CreateProductReq{
		Title:       "Phone",
		Description: "Smartphone",
		Price:       1000,
		Stock:       10,
		Image:       "img.jpg",
		CategoryID:  categoryID.String(),
		SKU:         "SKU-123",
	}

	mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*entity.Product")).Return(nil)

	res, err := uc.CreateProduct(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, req.Title, res.Title)
	assert.Equal(t, req.Price, res.Price)
	assert.NotEqual(t, uuid.Nil, res.ID)
	assert.Equal(t, "active", res.Status)

	mockRepo.AssertExpectations(t)
}

func TestGetAllProducts_Success(t *testing.T) {
	mockRepo := new(productMock.ProductRepositoryMock)
	uc := productUC.NewProductUseCase(mockRepo)

	mockProducts := []*entity.Product{
		{ID: uuid.New(), Title: "P1", Price: 100},
		{ID: uuid.New(), Title: "P2", Price: 200},
	}

	mockRepo.On("FindAll", mock.Anything).Return(mockProducts, nil)

	res, err := uc.GetAllProducts(context.Background())

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Len(t, res, 2)
	assert.Equal(t, mockProducts[0].Title, res[0].Title)

	mockRepo.AssertExpectations(t)
}

func TestGetProductByID_Success(t *testing.T) {
	mockRepo := new(productMock.ProductRepositoryMock)
	uc := productUC.NewProductUseCase(mockRepo)

	id := uuid.New()
	mockProduct := &entity.Product{ID: id, Title: "P1", Price: 100}

	mockRepo.On("FindByID", mock.Anything, id).Return(mockProduct, nil)

	res, err := uc.GetProductByID(context.Background(), id.String())

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, mockProduct.Title, res.Title)

	mockRepo.AssertExpectations(t)
}

func TestGetProductByID_NotFound(t *testing.T) {
	mockRepo := new(productMock.ProductRepositoryMock)
	uc := productUC.NewProductUseCase(mockRepo)

	id := uuid.New()

	mockRepo.On("FindByID", mock.Anything, id).Return(nil, gorm.ErrRecordNotFound)

	res, err := uc.GetProductByID(context.Background(), id.String())

	assert.Error(t, err)
	assert.Nil(t, res)

	var appErr *apperror.AppError
	assert.True(t, errors.As(err, &appErr))
	assert.Equal(t, apperror.CodeNotFound, appErr.Code)

	mockRepo.AssertExpectations(t)
}

func TestUpdateProduct_Success(t *testing.T) {
	mockRepo := new(productMock.ProductRepositoryMock)
	uc := productUC.NewProductUseCase(mockRepo)

	id := uuid.New()
	mockProduct := &entity.Product{ID: id, Title: "Old", Price: 100}

	req := &productDTO.UpdateProductReq{
		Title:       "New Phone",
		Description: "New Desc",
		Price:       1200,
		Stock:       20,
		Image:       "new.jpg",
		CategoryID:  uuid.New().String(),
		SKU:         "SKU-456",
		Status:      "inactive",
	}

	mockRepo.On("FindByID", mock.Anything, id).Return(mockProduct, nil)
	mockRepo.On("Update", mock.Anything, mock.AnythingOfType("*entity.Product")).Return(nil)

	res, err := uc.UpdateProduct(context.Background(), id.String(), req)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, req.Title, res.Title)
	assert.Equal(t, req.Price, res.Price)

	mockRepo.AssertExpectations(t)
}

func TestDeleteProduct_Success(t *testing.T) {
	mockRepo := new(productMock.ProductRepositoryMock)
	uc := productUC.NewProductUseCase(mockRepo)

	id := uuid.New()
	mockProduct := &entity.Product{ID: id, Title: "P1"}

	mockRepo.On("FindByID", mock.Anything, id).Return(mockProduct, nil)
	mockRepo.On("Delete", mock.Anything, id).Return(nil)

	err := uc.DeleteProduct(context.Background(), id.String())

	assert.NoError(t, err)

	mockRepo.AssertExpectations(t)
}
