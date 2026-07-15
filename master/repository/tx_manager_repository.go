package repository

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
)

// TransactionManager defines transaction operations
type TransactionManager interface {
	WithinTransaction(ctx context.Context, fn func(context.Context) error) error
}

// transactionManagerImpl implements TransactionManager with sqlx
type transactionManagerImpl struct {
	db *sqlx.DB
}

// NewTransactionManager creates a new transaction manager
func NewTransactionManager(db *sqlx.DB) TransactionManager {
	return &transactionManagerImpl{db: db}
}

// WithinTransaction wraps operations in a database transaction
// If fn returns an error, the transaction is rolled back
// If fn succeeds, the transaction is committed
func (tm *transactionManagerImpl) WithinTransaction(ctx context.Context, fn func(context.Context) error) error {
	tx, err := tm.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	// Create a context with the transaction embedded
	txCtx := context.WithValue(ctx, "tx", tx)

	// Execute the function
	if err := fn(txCtx); err != nil {
		// Rollback on error
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return fmt.Errorf("transaction failed (%w) and rollback failed (%w)", err, rollbackErr)
		}
		return err
	}

	// Commit on success
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// GetTxFromContext extracts the transaction from context
func GetTxFromContext(ctx context.Context) *sqlx.Tx {
	tx, ok := ctx.Value("tx").(*sqlx.Tx)
	if !ok {
		return nil
	}
	return tx
}
