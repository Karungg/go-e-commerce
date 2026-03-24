package port

import "github.com/google/uuid"

// TokenGenerator defines the contract for generating authentication credentials
type TokenGenerator interface {
	GenerateToken(userID uuid.UUID, role string) (string, error)
}

// TokenPayload holds the essential identity claims from an authenticated token
type TokenPayload struct {
	UserID uuid.UUID
	Role   string
}

// TokenValidator defines the contract for validating authentication credentials
type TokenValidator interface {
	ValidateToken(tokenString string) (*TokenPayload, error)
}
