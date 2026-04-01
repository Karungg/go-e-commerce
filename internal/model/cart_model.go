package model

import (
	"time"

	"github.com/google/uuid"
)

type CartModel struct {
	ID        uuid.UUID       `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UserID    uuid.UUID       `gorm:"type:uuid;uniqueIndex;not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Items     []CartItemModel `gorm:"foreignKey:CartID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

func (CartModel) TableName() string {
	return "carts"
}

type CartItemModel struct {
	ID        uuid.UUID    `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	CartID    uuid.UUID    `gorm:"type:uuid;not null;uniqueIndex:idx_cart_items_cart_product"`
	ProductID uuid.UUID    `gorm:"type:uuid;not null;uniqueIndex:idx_cart_items_cart_product"`
	Quantity  int          `gorm:"type:integer;not null;default:1"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Product   ProductModel `gorm:"foreignKey:ProductID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

func (CartItemModel) TableName() string {
	return "cart_items"
}
