package repository

import (
	"testing"
	"time"

	"master/entity"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
)

func setupProductAssignmentRepositoryTest(t *testing.T) (ProductAssignmentRepository, sqlmock.Sqlmock, func()) {
	t.Helper()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := NewProductAssignmentRepository(sqlxDB)

	cleanup := func() {
		_ = db.Close()
	}

	return repo, mock, cleanup
}

func TestProductAssignmentRepository_FindAll_CoalescesMissingCreatedByName(t *testing.T) {
	repo, mock, cleanup := setupProductAssignmentRepositoryTest(t)
	defer cleanup()

	filter := entity.ProductAssignmentQueryFilter{
		CustId: "C26009",
		Page:   1,
		Limit:  10,
	}

	mock.ExpectQuery(`(?s)SELECT COUNT\(\*\).*FROM mst\.m_product_assignment pa.*WHERE pa\.cust_id = \$1`).
		WithArgs("C26009").
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

	now := time.Date(2026, 6, 4, 10, 0, 0, 0, time.UTC)
	mock.ExpectQuery(`(?s)COALESCE\(NULLIF\(u\.user_fullname, ''\), u\.user_name, ''\) AS created_by_name.*FROM mst\.m_product_assignment pa.*LIMIT 10 OFFSET 0`).
		WithArgs("C26009").
		WillReturnRows(sqlmock.NewRows([]string{
			"id",
			"cust_id",
			"action_date",
			"pro_id",
			"pro_code",
			"pro_name",
			"distributor_id",
			"distributor_code",
			"distributor_name",
			"assignment_type",
			"created_by",
			"created_by_name",
			"created_at",
		}).AddRow(
			int64(1),
			"C26009",
			now,
			int64(101),
			"PRO-101",
			"Product 101",
			int64(201),
			"DIST-201",
			"Distributor 201",
			"assignment",
			int64(379),
			"",
			now,
		))

	rows, total, lastPage, err := repo.FindAll(filter)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if total != 1 || lastPage != 1 {
		t.Fatalf("expected total=1 lastPage=1, got total=%d lastPage=%d", total, lastPage)
	}
	if len(rows) != 1 {
		t.Fatalf("expected 1 row, got %d", len(rows))
	}
	if rows[0].CreatedByName != "" {
		t.Fatalf("expected empty created_by_name fallback, got %q", rows[0].CreatedByName)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}
