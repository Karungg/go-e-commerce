package repository

import (
	"context"

	authPort "go-e-commerce/internal/port/auth"

	"gorm.io/gorm"
)

type txKey struct{}

type transactionManager struct {
	db *gorm.DB
}

// NewTransactionManager creates a new transaction manager that uses GORM
func NewTransactionManager(db *gorm.DB) authPort.TransactionManager {
	return &transactionManager{db: db}
}

// RunInTransaction executes the given function inside a database transaction
func (t *transactionManager) RunInTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	return t.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Inject the transaction context into the standard context
		ctxWithTx := context.WithValue(ctx, txKey{}, tx)
		return fn(ctxWithTx)
	})
}

// ExtractTx is a helper function to extract the GORM transaction from the context,
// or return the default DB if no transaction is found.
func ExtractTx(ctx context.Context, db *gorm.DB) *gorm.DB {
	if tx, ok := ctx.Value(txKey{}).(*gorm.DB); ok {
		return tx
	}
	return db
}
