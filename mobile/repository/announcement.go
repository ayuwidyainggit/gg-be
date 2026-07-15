package repository

import (
	"context"

	"gorm.io/gorm"
)

type (
	RepositoryAnnouncementsImpl struct {
		*gorm.DB
	}
)

type AnnouncementsRepository interface {
}

func NewAnnouncementsRepository(db *gorm.DB) *RepositoryAnnouncementsImpl {
	return &RepositoryAnnouncementsImpl{db}
}

func (repo *RepositoryAnnouncementsImpl) model(ctx context.Context) *gorm.DB {
	tx := extractTx(ctx)
	if tx != nil {
		return tx.WithContext(ctx)
	}
	return repo.WithContext(ctx)
}
