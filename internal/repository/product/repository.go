package product

import (
	"context"
	"go-e-commerce/internal/entity"
	"go-e-commerce/internal/model"
	"go-e-commerce/internal/repository"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProductRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

func (r *ProductRepository) Create(ctx context.Context, product *entity.Product) error {
	productModel := &model.ProductModel{
		ID:          product.ID,
		Title:       product.Title,
		Description: product.Description,
		Price:       product.Price,
		Stock:       product.Stock,
		Image:       product.Image,
		CategoryID:  product.CategoryID,
		SKU:         product.SKU,
		Status:      product.Status,
	}

	db := repository.ExtractTx(ctx, r.db)
	if err := db.WithContext(ctx).Create(productModel).Error; err != nil {
		return err
	}

	return nil
}

func (r *ProductRepository) FindAll(ctx context.Context, limit, offset int) ([]*entity.Product, int64, error) {
	var total int64
	if err := r.db.WithContext(ctx).Model(&model.ProductModel{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var productModels []model.ProductModel
	if err := r.db.WithContext(ctx).Limit(limit).Offset(offset).Find(&productModels).Error; err != nil {
		return nil, 0, err
	}

	products := make([]*entity.Product, len(productModels))
	for i, m := range productModels {
		products[i] = &entity.Product{
			ID:          m.ID,
			Title:       m.Title,
			Description: m.Description,
			Price:       m.Price,
			Stock:       m.Stock,
			Image:       m.Image,
			CategoryID:  m.CategoryID,
			SKU:         m.SKU,
			Status:      m.Status,
		}
	}
	return products, total, nil
}

func (r *ProductRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.Product, error) {
	var productModel model.ProductModel
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&productModel).Error; err != nil {
		return nil, err
	}

	return &entity.Product{
		ID:          productModel.ID,
		Title:       productModel.Title,
		Description: productModel.Description,
		Price:       productModel.Price,
		Stock:       productModel.Stock,
		Image:       productModel.Image,
		CategoryID:  productModel.CategoryID,
		SKU:         productModel.SKU,
		Status:      productModel.Status,
	}, nil
}

func (r *ProductRepository) Update(ctx context.Context, product *entity.Product) error {
	productModel := &model.ProductModel{
		ID:          product.ID,
		Title:       product.Title,
		Description: product.Description,
		Price:       product.Price,
		Stock:       product.Stock,
		Image:       product.Image,
		CategoryID:  product.CategoryID,
		SKU:         product.SKU,
		Status:      product.Status,
	}

	db := repository.ExtractTx(ctx, r.db)
	if err := db.WithContext(ctx).Model(&model.ProductModel{}).Where("id = ?", product.ID).Updates(productModel).Error; err != nil {
		return err
	}

	return nil
}

func (r *ProductRepository) Delete(ctx context.Context, id uuid.UUID) error {
	db := repository.ExtractTx(ctx, r.db)
	if err := db.WithContext(ctx).Where("id = ?", id).Delete(&model.ProductModel{}).Error; err != nil {
		return err
	}
	return nil
}
