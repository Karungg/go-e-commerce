package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SellerModel struct {
	ID               uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UserID           uuid.UUID      `gorm:"type:uuid;uniqueIndex;not null"`
	User             UserModel      `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	StoreName        string         `gorm:"type:varchar(255);not null"`
	StoreDescription string         `gorm:"type:text"`
	LogoUrl          string         `gorm:"type:varchar(255)"`
	IsVerified       bool           `gorm:"type:boolean;default:false"`

	CreatedAt        time.Time
	UpdatedAt        time.Time
	DeletedAt        gorm.DeletedAt `gorm:"index"`
}

func (SellerModel) TableName() string {
	return "sellers"
}
