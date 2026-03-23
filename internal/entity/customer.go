package entity

import (
	"time"

	"github.com/google/uuid"
)

type Customer struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	FirstName string
	LastName  string
	Phone     string
	Address   string
	CreatedAt time.Time
	UpdatedAt time.Time
}
