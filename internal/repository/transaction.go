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

func NewTransactionManager(db *gorm.DB) authPort.TransactionManager {
	return &transactionManager{db: db}
}

func (t *transactionManager) RunInTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	return t.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		ctxWithTx := context.WithValue(ctx, txKey{}, tx)
		return fn(ctxWithTx)
	})
}

func ExtractTx(ctx context.Context, db *gorm.DB) *gorm.DB {
	if tx, ok := ctx.Value(txKey{}).(*gorm.DB); ok {
		return tx
	}
	return db
}
