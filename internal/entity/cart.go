package entity

import (
	"time"

	"github.com/google/uuid"
)

type Cart struct {
	ID        uuid.UUID   `json:"id"`
	UserID    uuid.UUID   `json:"user_id"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
	Items     []CartItem  `json:"items,omitempty"`
}

type CartItem struct {
	ID        uuid.UUID  `json:"id"`
	CartID    uuid.UUID  `json:"cart_id"`
	ProductID uuid.UUID  `json:"product_id"`
	Quantity  int        `json:"quantity"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	Product   *Product   `json:"product,omitempty"` // For joining
}
