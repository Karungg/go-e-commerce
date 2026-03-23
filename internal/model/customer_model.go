package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CustomerModel struct {
	ID        uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UserID    uuid.UUID      `gorm:"type:uuid;uniqueIndex;not null"`
	User      UserModel      `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	FirstName string         `gorm:"type:varchar(100);not null"`
	LastName  string         `gorm:"type:varchar(100);not null"`
	Phone     string         `gorm:"type:varchar(30)"`
	Address   string         `gorm:"type:text"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (CustomerModel) TableName() string {
	return "customers"
}
