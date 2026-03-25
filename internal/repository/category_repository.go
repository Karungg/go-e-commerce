package repository

import (
	"context"
	"go-e-commerce/internal/entity"
	"go-e-commerce/internal/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CategoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) *CategoryRepository {
	return &CategoryRepository{db: db}
}

func (r *CategoryRepository) Create(ctx context.Context, category *entity.Category) error {
	categoryModel := &model.CategoryModel{
		ID:          category.ID,
		Title:       category.Title,
		Description: category.Description,
	}

	db := ExtractTx(ctx, r.db)
	if err := db.WithContext(ctx).Create(categoryModel).Error; err != nil {
		return err
	}

	return nil
}

func (r *CategoryRepository) FindAll(ctx context.Context) ([]*entity.Category, error) {
	var categoryModels []model.CategoryModel
	if err := r.db.WithContext(ctx).Find(&categoryModels).Error; err != nil {
		return nil, err
	}

	categories := make([]*entity.Category, len(categoryModels))
	for i, m := range categoryModels {
		categories[i] = &entity.Category{
			ID:          m.ID,
			Title:       m.Title,
			Description: m.Description,
		}
	}
	return categories, nil
}

func (r *CategoryRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.Category, error) {
	var categoryModel model.CategoryModel
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&categoryModel).Error; err != nil {
		return nil, err
	}

	return &entity.Category{
		ID:          categoryModel.ID,
		Title:       categoryModel.Title,
		Description: categoryModel.Description,
	}, nil
}

func (r *CategoryRepository) Update(ctx context.Context, category *entity.Category) error {
	categoryModel := &model.CategoryModel{
		ID:          category.ID,
		Title:       category.Title,
		Description: category.Description,
	}

	db := ExtractTx(ctx, r.db)
	if err := db.WithContext(ctx).Model(&model.CategoryModel{}).Where("id = ?", category.ID).Updates(categoryModel).Error; err != nil {
		return err
	}

	return nil
}

func (r *CategoryRepository) Delete(ctx context.Context, id uuid.UUID) error {
	db := ExtractTx(ctx, r.db)
	if err := db.WithContext(ctx).Where("id = ?", id).Delete(&model.CategoryModel{}).Error; err != nil {
		return err
	}
	return nil
}
