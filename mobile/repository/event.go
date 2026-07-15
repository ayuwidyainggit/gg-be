package repository

import (
	"context"

	"gorm.io/gorm"
)

type (
	RepositoryEventsImpl struct {
		*gorm.DB
	}
)

type EventsRepository interface {
}

func NewEventsRepository(db *gorm.DB) *RepositoryEventsImpl {
	return &RepositoryEventsImpl{db}
}

func (repo *RepositoryEventsImpl) model(ctx context.Context) *gorm.DB {
	tx := extractTx(ctx)
	if tx != nil {
		return tx.WithContext(ctx)
	}
	return repo.WithContext(ctx)
}
