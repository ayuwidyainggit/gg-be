package repository

import (
	"context"

	"gorm.io/gorm"
)

type (
	RepositoryLeaderboardsImpl struct {
		*gorm.DB
	}
)

type LeaderboardsRepository interface {
}

func NewLeaderboardsRepository(db *gorm.DB) *RepositoryLeaderboardsImpl {
	return &RepositoryLeaderboardsImpl{db}
}

func (repo *RepositoryLeaderboardsImpl) model(ctx context.Context) *gorm.DB {
	tx := extractTx(ctx)
	if tx != nil {
		return tx.WithContext(ctx)
	}
	return repo.WithContext(ctx)
}
