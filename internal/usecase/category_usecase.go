package usecase

import (
	"context"
	"errors"

	"go-e-commerce/internal/dto"
	"go-e-commerce/internal/entity"
	"go-e-commerce/internal/pkg/apperror"
	"go-e-commerce/internal/port"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type categoryUseCase struct {
	categoryRepo port.CategoryRepository
}

func NewCategoryUseCase(categoryRepo port.CategoryRepository) port.CategoryUseCase {
	return &categoryUseCase{
		categoryRepo: categoryRepo,
	}
}

func (u *categoryUseCase) CreateCategory(ctx context.Context, req *dto.CreateCategoryReq) (*dto.CategoryRes, error) {
	category := &entity.Category{
		ID:          uuid.New(),
		Title:       req.Title,
		Description: req.Description,
	}

	if err := u.categoryRepo.Create(ctx, category); err != nil {
		return nil, err
	}

	return &dto.CategoryRes{
		ID:          category.ID,
		Title:       category.Title,
		Description: category.Description,
	}, nil
}

func (u *categoryUseCase) GetAllCategories(ctx context.Context) ([]*dto.CategoryRes, error) {
	categories, err := u.categoryRepo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	var res []*dto.CategoryRes
	for _, c := range categories {
		res = append(res, &dto.CategoryRes{
			ID:          c.ID, Title:       c.Title,
			Description: c.Description,
		})
	}
	return res, nil
}

func (u *categoryUseCase) GetCategoryByID(ctx context.Context, id string) (*dto.CategoryRes, error) {
	categoryID, err := uuid.Parse(id)
	if err != nil {
		return nil, apperror.ErrBadRequest
	}

	category, err := u.categoryRepo.FindByID(ctx, categoryID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.ErrCategoryNotFound
		}
		return nil, err
	}

	return &dto.CategoryRes{
		ID:          category.ID,
		Title:       category.Title,
		Description: category.Description,
	}, nil
}

func (u *categoryUseCase) UpdateCategory(ctx context.Context, id string, req *dto.UpdateCategoryReq) (*dto.CategoryRes, error) {
	categoryID, err := uuid.Parse(id)
	if err != nil {
		return nil, apperror.ErrBadRequest
	}

	category, err := u.categoryRepo.FindByID(ctx, categoryID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.ErrCategoryNotFound
		}
		return nil, err
	}

	category.Title = req.Title
	category.Description = req.Description

	if err := u.categoryRepo.Update(ctx, category); err != nil {
		return nil, err
	}

	return &dto.CategoryRes{
		ID:          category.ID,
		Title:       category.Title,
		Description: category.Description,
	}, nil
}

func (u *categoryUseCase) DeleteCategory(ctx context.Context, id string) error {
	categoryID, err := uuid.Parse(id)
	if err != nil {
		return apperror.ErrBadRequest
	}

	_, err = u.categoryRepo.FindByID(ctx, categoryID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperror.ErrCategoryNotFound
		}
		return err
	}

	return u.categoryRepo.Delete(ctx, categoryID)
}
