package entity

import (
	"time"

	"github.com/google/uuid"
)

type Role string

const (
	RoleAdmin    Role = "admin"
	RoleCustomer Role = "customer"
	RoleSeller   Role = "seller"
)

type User struct {
	ID        uuid.UUID
	Email     string
	Password  string
	Role      Role
	IsActive  bool
	CreatedAt time.Time
	UpdatedAt time.Time
}
