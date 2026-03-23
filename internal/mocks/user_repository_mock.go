package mocks

import (
	"go-e-commerce/internal/entity"

	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

type UserRepositoryMock struct {
	mock.Mock
}

func (m *UserRepositoryMock) CreateWithTx(tx *gorm.DB, user *entity.User) error {
	args := m.Called(tx, user)
	return args.Error(0)
}

func (m *UserRepositoryMock) FindByEmail(email string) (*entity.User, error) {
	args := m.Called(email)
	if args.Get(0) != nil {
		return args.Get(0).(*entity.User), args.Error(1)
	}
	return nil, args.Error(1)
}
