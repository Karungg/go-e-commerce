package usecase_test

import (
	"context"
	"errors"
	"testing"

	"go-e-commerce/internal/dto"
	"go-e-commerce/internal/entity"
	"go-e-commerce/internal/mocks"
	"go-e-commerce/internal/pkg/apperror"
	"go-e-commerce/internal/usecase"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

func TestCreateCategory_Success(t *testing.T) {
	mockRepo := new(mocks.CategoryRepositoryMock)
	uc := usecase.NewCategoryUseCase(mockRepo)

	req := &dto.CreateCategoryReq{
		Title:       "Electronics",
		Description: "Gadgets",
	}

	mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*entity.Category")).Return(nil)

	res, err := uc.CreateCategory(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, req.Title, res.Title)
	assert.Equal(t, req.Description, res.Description)
	assert.NotEqual(t, uuid.Nil, res.ID)

	mockRepo.AssertExpectations(t)
}

func TestGetAllCategories_Success(t *testing.T) {
	mockRepo := new(mocks.CategoryRepositoryMock)
	uc := usecase.NewCategoryUseCase(mockRepo)

	mockCategories := []*entity.Category{
		{ID: uuid.New(), Title: "C1", Description: "D1"},
		{ID: uuid.New(), Title: "C2", Description: "D2"},
	}

	mockRepo.On("FindAll", mock.Anything).Return(mockCategories, nil)

	res, err := uc.GetAllCategories(context.Background())

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Len(t, res, 2)
	assert.Equal(t, mockCategories[0].Title, res[0].Title)

	mockRepo.AssertExpectations(t)
}

func TestGetCategoryByID_Success(t *testing.T) {
	mockRepo := new(mocks.CategoryRepositoryMock)
	uc := usecase.NewCategoryUseCase(mockRepo)

	id := uuid.New()
	mockCategory := &entity.Category{ID: id, Title: "C1", Description: "D1"}

	mockRepo.On("FindByID", mock.Anything, id).Return(mockCategory, nil)

	res, err := uc.GetCategoryByID(context.Background(), id.String())

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, mockCategory.Title, res.Title)

	mockRepo.AssertExpectations(t)
}

func TestGetCategoryByID_NotFound(t *testing.T) {
	mockRepo := new(mocks.CategoryRepositoryMock)
	uc := usecase.NewCategoryUseCase(mockRepo)

	id := uuid.New()

	mockRepo.On("FindByID", mock.Anything, id).Return(nil, gorm.ErrRecordNotFound)

	res, err := uc.GetCategoryByID(context.Background(), id.String())

	assert.Error(t, err)
	assert.Nil(t, res)
	
	var appErr *apperror.AppError
	assert.True(t, errors.As(err, &appErr))
	assert.Equal(t, apperror.CodeNotFound, appErr.Code)

	mockRepo.AssertExpectations(t)
}

func TestUpdateCategory_Success(t *testing.T) {
	mockRepo := new(mocks.CategoryRepositoryMock)
	uc := usecase.NewCategoryUseCase(mockRepo)

	id := uuid.New()
	mockCategory := &entity.Category{ID: id, Title: "Old", Description: "Old Desc"}

	req := &dto.UpdateCategoryReq{
		Title:       "New Title",
		Description: "New Desc",
	}

	mockRepo.On("FindByID", mock.Anything, id).Return(mockCategory, nil)
	mockRepo.On("Update", mock.Anything, mock.AnythingOfType("*entity.Category")).Return(nil)

	res, err := uc.UpdateCategory(context.Background(), id.String(), req)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, req.Title, res.Title)

	mockRepo.AssertExpectations(t)
}

func TestDeleteCategory_Success(t *testing.T) {
	mockRepo := new(mocks.CategoryRepositoryMock)
	uc := usecase.NewCategoryUseCase(mockRepo)

	id := uuid.New()
	mockCategory := &entity.Category{ID: id, Title: "C1", Description: "D1"}

	mockRepo.On("FindByID", mock.Anything, id).Return(mockCategory, nil)
	mockRepo.On("Delete", mock.Anything, id).Return(nil)

	err := uc.DeleteCategory(context.Background(), id.String())

	assert.NoError(t, err)

	mockRepo.AssertExpectations(t)
}
