package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProductModel struct {
	ID          uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Title       string         `gorm:"type:varchar(255);not null"`
	Description string         `gorm:"type:text"`
	Price       float64        `gorm:"type:decimal(12,2);not null"`
	Stock       int            `gorm:"type:integer;not null;default:0"`
	Image       string         `gorm:"type:varchar(255)"`
	CategoryID  uuid.UUID      `gorm:"type:uuid;not null"`
	Category    CategoryModel  `gorm:"foreignKey:CategoryID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`

	SKU         string         `gorm:"type:varchar(100);uniqueIndex"`
	Status      string         `gorm:"type:varchar(50);default:'active'"`

	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

func (ProductModel) TableName() string {
	return "products"
}
