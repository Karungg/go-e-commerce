package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserModel struct {
	ID        uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Email     string         `gorm:"type:varchar(255);uniqueIndex;not null"`
	Password  string         `gorm:"type:varchar(255);not null"`
	Role      string         `gorm:"type:varchar(50);not null;default:'customer'"` // admin, customer, seller
	IsActive  bool           `gorm:"type:boolean;default:true"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (UserModel) TableName() string {
	return "users"
}
