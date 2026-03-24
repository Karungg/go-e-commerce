package repository

import (
	"context"
	"go-e-commerce/internal/entity"
	"go-e-commerce/internal/model"

	"gorm.io/gorm"
)

type SellerRepository struct {
	db *gorm.DB
}

func NewSellerRepository(db *gorm.DB) *SellerRepository {
	return &SellerRepository{db: db}
}

func (r *SellerRepository) CreateWithTx(ctx context.Context, tx *gorm.DB, seller *entity.Seller) error {
	sellerModel := &model.SellerModel{
		ID:               seller.ID,
		UserID:           seller.UserID,
		StoreName:        seller.StoreName,
		StoreDescription: seller.StoreDescription,
		LogoUrl:          seller.LogoUrl,
		IsVerified:       seller.IsVerified,
	}

	if err := tx.WithContext(ctx).Create(sellerModel).Error; err != nil {
		return err
	}

	seller.CreatedAt = sellerModel.CreatedAt
	seller.UpdatedAt = sellerModel.UpdatedAt
	return nil
}

func (r *SellerRepository) FindByStoreName(ctx context.Context, storeName string) (*entity.Seller, error) {
	var sellerModel model.SellerModel
	if err := r.db.WithContext(ctx).Where("store_name = ?", storeName).First(&sellerModel).Error; err != nil {
		return nil, err
	}

	return &entity.Seller{
		ID:               sellerModel.ID,
		UserID:           sellerModel.UserID,
		StoreName:        sellerModel.StoreName,
		StoreDescription: sellerModel.StoreDescription,
		LogoUrl:          sellerModel.LogoUrl,
		IsVerified:       sellerModel.IsVerified,
		CreatedAt:        sellerModel.CreatedAt,
		UpdatedAt:        sellerModel.UpdatedAt,
	}, nil
}
