package entity

import (
	"time"

	"github.com/google/uuid"
)

type Product struct {
	ID          uuid.UUID
	Title       string
	Description string
	Price       float64
	Stock       int
	Image       string
	CategoryID  uuid.UUID
	SKU         string
	Status      string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
