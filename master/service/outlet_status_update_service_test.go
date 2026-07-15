package service

import (
	"errors"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"

	"master/repository"
)

func setupOutletStatusUpdateServiceTest(t *testing.T) (*outletServiceImpl, sqlmock.Sqlmock, func()) {
	t.Helper()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := repository.NewOutletRepository(sqlxDB)
	svc := NewOutletService(repo, nil, nil, nil)

	cleanup := func() {
		_ = db.Close()
	}
	return svc, mock, cleanup
}

func TestOutletService_UpdateStatuses_Success(t *testing.T) {
	svc, mock, cleanup := setupOutletStatusUpdateServiceTest(t)
	defer cleanup()

	mock.ExpectQuery(`(?s)WITH base AS`).
		WillReturnRows(sqlmock.NewRows([]string{"rows_affected"}).AddRow(int64(7)))
	mock.ExpectExec(`(?s)UPDATE mst.m_outlet o`).
		WillReturnResult(sqlmock.NewResult(0, 3))

	rows, err := svc.UpdateStatuses()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if rows != 10 {
		t.Fatalf("expected rows_affected 10, got %d", rows)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestOutletService_UpdateStatuses_NoChange(t *testing.T) {
	svc, mock, cleanup := setupOutletStatusUpdateServiceTest(t)
	defer cleanup()

	mock.ExpectQuery(`(?s)WITH base AS`).
		WillReturnRows(sqlmock.NewRows([]string{"rows_affected"}).AddRow(int64(0)))
	mock.ExpectExec(`(?s)UPDATE mst.m_outlet o`).
		WillReturnResult(sqlmock.NewResult(0, 0))

	rows, err := svc.UpdateStatuses()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if rows != 0 {
		t.Fatalf("expected rows_affected 0, got %d", rows)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestOutletService_UpdateStatuses_Error(t *testing.T) {
	svc, mock, cleanup := setupOutletStatusUpdateServiceTest(t)
	defer cleanup()

	mock.ExpectQuery(`(?s)WITH base AS`).
		WillReturnError(errors.New("db exploded"))

	_, err := svc.UpdateStatuses()
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestOutletService_UpdateStatuses_NooPromotionError(t *testing.T) {
	svc, mock, cleanup := setupOutletStatusUpdateServiceTest(t)
	defer cleanup()

	mock.ExpectQuery(`(?s)WITH base AS`).
		WillReturnRows(sqlmock.NewRows([]string{"rows_affected"}).AddRow(int64(2)))
	mock.ExpectExec(`(?s)UPDATE mst.m_outlet o`).
		WillReturnError(errors.New("db exploded"))

	_, err := svc.UpdateStatuses()
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}
