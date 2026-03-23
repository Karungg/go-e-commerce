package repository

import (
	"go-e-commerce/internal/entity"
	"go-e-commerce/internal/model"

	"gorm.io/gorm"
)

type SellerRepository interface {
	CreateWithTx(tx *gorm.DB, seller *entity.Seller) error
}

type sellerRepository struct {
	db *gorm.DB
}

func NewSellerRepository(db *gorm.DB) SellerRepository {
	return &sellerRepository{db: db}
}

func (r *sellerRepository) CreateWithTx(tx *gorm.DB, seller *entity.Seller) error {
	sellerModel := &model.SellerModel{
		ID:               seller.ID,
		UserID:           seller.UserID,
		StoreName:        seller.StoreName,
		StoreDescription: seller.StoreDescription,
		LogoUrl:          seller.LogoUrl,
		IsVerified:       seller.IsVerified,
	}

	if err := tx.Create(sellerModel).Error; err != nil {
		return err
	}

	seller.CreatedAt = sellerModel.CreatedAt
	seller.UpdatedAt = sellerModel.UpdatedAt
	return nil
}
