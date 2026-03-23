package repository

import (
	"context"
	"go-e-commerce/internal/entity"
	"go-e-commerce/internal/model"

	"gorm.io/gorm"
)

type CustomerRepository interface {
	CreateWithTx(ctx context.Context, tx *gorm.DB, customer *entity.Customer) error
}

type customerRepository struct {
	db *gorm.DB
}

func NewCustomerRepository(db *gorm.DB) CustomerRepository {
	return &customerRepository{db: db}
}

func (r *customerRepository) CreateWithTx(ctx context.Context, tx *gorm.DB, customer *entity.Customer) error {
	customerModel := &model.CustomerModel{
		ID:        customer.ID,
		UserID:    customer.UserID,
		FirstName: customer.FirstName,
		LastName:  customer.LastName,
		Phone:     customer.Phone,
		Address:   customer.Address,
	}

	if err := tx.WithContext(ctx).Create(customerModel).Error; err != nil {
		return err
	}

	customer.CreatedAt = customerModel.CreatedAt
	customer.UpdatedAt = customerModel.UpdatedAt
	return nil
}
