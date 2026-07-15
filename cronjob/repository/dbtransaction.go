package repository

import (
	"context"

	"github.com/gofiber/fiber/v2/log"

	"gorm.io/gorm"
)

type Dbtransaction interface {
	WithinTransaction(ctx context.Context, tFunc func(ctx context.Context) error) error
}

type (
	RepositoryDbtransaction struct {
		*gorm.DB
	}
)

func NewDbtransactionRepo(db *gorm.DB) *RepositoryDbtransaction {
	return &RepositoryDbtransaction{db}
}

var txKey *gorm.DB

// extractTx extracts transaction from context
func extractTx(ctx context.Context) *gorm.DB {
	if tx, ok := ctx.Value(txKey).(*gorm.DB); ok {

		return tx
	}
	return nil
}

// injectTx injects transaction to context
func injectTx(ctx context.Context, tx *gorm.DB) context.Context {
	return context.WithValue(ctx, txKey, tx)
}

func (repo *RepositoryDbtransaction) WithinTransaction(ctx context.Context, tFunc func(ctx context.Context) error) error {

	tx := repo.DB.Begin()

	// run callback
	err := tFunc(injectTx(ctx, tx))
	if err != nil {
		if errRollback := tx.Rollback(); errRollback != nil {
			log.Errorf("rollback transaction: %v", errRollback)
		}
		return err
	}

	tx.Commit();
	return nil
}
