package cart

import (
	"time"

	"github.com/google/uuid"
)

type ProductResponse struct {
	ID          uuid.UUID `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	Stock       int       `json:"stock"`
	Image       string    `json:"image"`
}

type CartItemResponse struct {
	ID        uuid.UUID       `json:"id"`
	CartID    uuid.UUID       `json:"cart_id"`
	ProductID uuid.UUID       `json:"product_id"`
	Quantity  int             `json:"quantity"`
	Product   ProductResponse `json:"product"`
	Subtotal  float64         `json:"subtotal"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
}

type CartResponse struct {
	ID         uuid.UUID          `json:"id"`
	UserID     uuid.UUID          `json:"user_id"`
	Items      []CartItemResponse `json:"items"`
	TotalPrice float64            `json:"total_price"`
	CreatedAt  time.Time          `json:"created_at"`
	UpdatedAt  time.Time          `json:"updated_at"`
}
