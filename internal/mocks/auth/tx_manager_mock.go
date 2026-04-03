package auth

import (
	"context"
)

type TransactionManagerMock struct {
}

func (m *TransactionManagerMock) RunInTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	return fn(ctx)
}
