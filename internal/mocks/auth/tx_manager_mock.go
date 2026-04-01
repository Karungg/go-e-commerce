package auth

import (
	"context"
)

type TransactionManagerMock struct {
}

func (m *TransactionManagerMock) RunInTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	// For testing purposes, we just execute the function directly without a real DB transaction
	return fn(ctx)
}
