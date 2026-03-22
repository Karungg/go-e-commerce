package entity

import "github.com/google/uuid"

type Category struct {
	ID          uuid.UUID
	Title       string
	Description string
}
