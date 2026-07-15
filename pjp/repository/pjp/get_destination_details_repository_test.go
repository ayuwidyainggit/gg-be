package pjp

import (
	"context"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestGetDestinationDetailsPrincipalScopesByPjpDateCustomerAndSorts(t *testing.T) {
	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New() error = %v", err)
	}
	defer sqlDB.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: sqlDB, PreferSimpleProtocol: true}), &gorm.Config{})
	if err != nil {
		t.Fatalf("gorm.Open() error = %v", err)
	}

	repo := NewPjpRepository()

	mock.ExpectQuery(`(?s)SELECT count\(\*\) FROM pjp_principles\.destinations_history dh WHERE DATE\(dh\.date\) = \$1 AND dh\.pjp_id = \$2 AND dh\.cust_id IN \(SELECT cust_id FROM "smc"\."m_customer" WHERE cust_id = \$3 OR parent_cust_id = \$4\)`).
		WithArgs("2026-06-29", 62, "CUST-1", "CUST-1").
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(3))

	rows := sqlmock.NewRows([]string{
		"route_code", "route_name", "week", "year", "date", "destination_id", "destination_code",
		"destination_type", "destination_name", "longitude", "latitude", "destination_status", "destination_address",
	}).AddRow(3883, "Route 1", 84, 2026, time.Date(2026, 6, 29, 0, 0, 0, 0, time.UTC), 1722, "BMI260003", "outlet", "Toko merah", "1", "2", "1", "Jalan Bangka 1")

	mock.ExpectQuery(`(?s).*FROM pjp_principles\.destinations_history dh.*LEFT JOIN pjp_principles\.routes r ON r\.route_code = dh\.route_code.*WHERE DATE\(dh\.date\) = \$1 AND dh\.pjp_id = \$2 AND dh\.cust_id IN \(SELECT cust_id FROM "smc"\."m_customer" WHERE cust_id = \$3 OR parent_cust_id = \$4\).*ORDER BY dh\.destination_id desc LIMIT 20`).
		WithArgs("2026-06-29", 62, "CUST-1", "CUST-1").
		WillReturnRows(rows)

	result, total, err := repo.GetDestinationDetails(context.Background(), gormDB, 62, "2026-06-29", 20, 1, "desc", "CUST-1", true)
	if err != nil {
		t.Fatalf("GetDestinationDetails() error = %v", err)
	}
	if total != 3 || len(result) != 1 {
		t.Fatalf("total/result = %d/%d, want 3/1", total, len(result))
	}
	if result[0].DestinationCode != "BMI260003" {
		t.Fatalf("destination code = %q, want BMI260003", result[0].DestinationCode)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("sqlmock expectations not met: %v", err)
	}
}

func TestGetDestinationDetailsDistributorUsesRouteOutletHistory(t *testing.T) {
	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New() error = %v", err)
	}
	defer sqlDB.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: sqlDB, PreferSimpleProtocol: true}), &gorm.Config{})
	if err != nil {
		t.Fatalf("gorm.Open() error = %v", err)
	}

	repo := NewPjpRepository()

	mock.ExpectQuery(`SELECT count\(\*\) FROM pjp\.route_outlet_history roh WHERE roh\.date = \$1 AND roh\.pjp_id = \$2 AND roh\.cust_id = \$3`).
		WithArgs("2026-06-01", 204, "DIST-1").
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

	rows := sqlmock.NewRows([]string{
		"route_code", "route_name", "week", "year", "date", "destination_id", "destination_code",
		"destination_type", "destination_name", "longitude", "latitude", "destination_status", "destination_address",
	}).AddRow(3883, "Route 1", 84, 2026, time.Date(2026, 6, 1, 0, 0, 0, 0, time.UTC), 1722, "BMI260003", "outlet", "Toko merah", "1", "2", "1", "Jalan Bangka 1")

	mock.ExpectQuery(`(?s).*FROM pjp\.route_outlet_history roh WHERE roh\.date = \$1 AND roh\.pjp_id = \$2 AND roh\.cust_id = \$3.*ORDER BY roh\.outlet_id asc LIMIT 20`).
		WithArgs("2026-06-01", 204, "DIST-1").
		WillReturnRows(rows)

	result, total, err := repo.GetDestinationDetails(context.Background(), gormDB, 204, "2026-06-01", 20, 1, "asc", "DIST-1", false)
	if err != nil {
		t.Fatalf("GetDestinationDetails() error = %v", err)
	}
	if total != 1 || len(result) != 1 {
		t.Fatalf("total/result = %d/%d, want 1/1", total, len(result))
	}
	if result[0].DestinationID != 1722 || result[0].DestinationType != "outlet" {
		t.Fatalf("row = %+v, want outlet 1722", result[0])
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("sqlmock expectations not met: %v", err)
	}
}
