package entity

import (
	"time"

	"github.com/google/uuid"
)

type Seller struct {
	ID               uuid.UUID
	UserID           uuid.UUID
	StoreName        string
	StoreDescription string
	LogoUrl          string
	IsVerified       bool
	CreatedAt        time.Time
	UpdatedAt        time.Time
}
