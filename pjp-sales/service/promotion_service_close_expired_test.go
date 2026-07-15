package service

import (
	"errors"
	"testing"
	"time"

	"sales/repository"
)

type mockPromotionV2RepoForCloseExpired struct {
	repository.PromotionV2Repository
	closeExpiredFn func(expiredBefore time.Time) (int64, error)
}

func (m *mockPromotionV2RepoForCloseExpired) CloseExpiredPromotionStatuses(expiredBefore time.Time) (int64, error) {
	if m.closeExpiredFn != nil {
		return m.closeExpiredFn(expiredBefore)
	}
	return 0, nil
}

func TestCloseExpiredPromotionsUsesStartOfCurrentWIBDay(t *testing.T) {
	var gotExpiredBefore time.Time
	repo := &mockPromotionV2RepoForCloseExpired{
		closeExpiredFn: func(expiredBefore time.Time) (int64, error) {
			gotExpiredBefore = expiredBefore
			return 2, nil
		},
	}
	service := &promotionServiceImpl{PromotionV2Repository: repo}

	if err := service.CloseExpiredPromotions(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	loc, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		t.Fatalf("unexpected error loading timezone: %v", err)
	}

	now := time.Now().In(loc)
	expected := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc)
	if !gotExpiredBefore.Equal(expected) {
		t.Fatalf("expected expiredBefore %v, got %v", expected, gotExpiredBefore)
	}
}

func TestCloseExpiredPromotionsReturnsRepositoryError(t *testing.T) {
	repo := &mockPromotionV2RepoForCloseExpired{
		closeExpiredFn: func(expiredBefore time.Time) (int64, error) {
			return 0, errors.New("db error")
		},
	}
	service := &promotionServiceImpl{PromotionV2Repository: repo}

	if err := service.CloseExpiredPromotions(); err == nil {
		t.Fatal("expected repository error")
	}
}
