package port

import "github.com/google/uuid"

// TokenGenerator defines the contract for generating authentication credentials
type TokenGenerator interface {
	GenerateToken(userID uuid.UUID, role string) (string, error)
}
