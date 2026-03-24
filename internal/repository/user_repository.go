package repository

import (
	"context"
	"go-e-commerce/internal/entity"
	"go-e-commerce/internal/model"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, user *entity.User) error {
	userModel := &model.UserModel{
		ID:       user.ID,
		Email:    user.Email,
		Password: user.Password,
		Role:     string(user.Role),
		IsActive: user.IsActive,
	}

	db := ExtractTx(ctx, r.db)
	if err := db.WithContext(ctx).Create(userModel).Error; err != nil {
		return err
	}
	
	user.CreatedAt = userModel.CreatedAt
	user.UpdatedAt = userModel.UpdatedAt
	return nil
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	var userModel model.UserModel
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&userModel).Error; err != nil {
		return nil, err
	}

	return &entity.User{
		ID:        userModel.ID,
		Email:     userModel.Email,
		Password:  userModel.Password,
		Role:      entity.Role(userModel.Role),
		IsActive:  userModel.IsActive,
		CreatedAt: userModel.CreatedAt,
		UpdatedAt: userModel.UpdatedAt,
	}, nil
}
