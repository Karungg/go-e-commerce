package repository

import (
	"context"
	"go-e-commerce/internal/entity"
	"go-e-commerce/internal/model"

	"gorm.io/gorm"
)

type UserRepository interface {
	CreateWithTx(ctx context.Context, tx *gorm.DB, user *entity.User) error
	FindByEmail(ctx context.Context, email string) (*entity.User, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) CreateWithTx(ctx context.Context, tx *gorm.DB, user *entity.User) error {
	userModel := &model.UserModel{
		ID:       user.ID,
		Email:    user.Email,
		Password: user.Password,
		Role:     string(user.Role),
		IsActive: user.IsActive,
	}

	if err := tx.WithContext(ctx).Create(userModel).Error; err != nil {
		return err
	}
	
	user.CreatedAt = userModel.CreatedAt
	user.UpdatedAt = userModel.UpdatedAt
	return nil
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
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
