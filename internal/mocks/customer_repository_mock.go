package mocks

import (
	"go-e-commerce/internal/entity"

	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

type CustomerRepositoryMock struct {
	mock.Mock
}

func (m *CustomerRepositoryMock) CreateWithTx(tx *gorm.DB, customer *entity.Customer) error {
	args := m.Called(tx, customer)
	return args.Error(0)
}
