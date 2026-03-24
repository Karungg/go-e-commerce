package port

import "context"

// TransactionManager defines the contract for running operations within a database transaction
type TransactionManager interface {
	RunInTransaction(ctx context.Context, fn func(ctx context.Context) error) error
}
