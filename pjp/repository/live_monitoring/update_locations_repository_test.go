package live_monitoring

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestGetUpdateLocations_ReturnsTimelineForEmployee_PJP(t *testing.T) {
	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New() error = %v", err)
	}
	defer sqlDB.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: sqlDB, PreferSimpleProtocol: true}), &gorm.Config{})
	if err != nil {
		t.Fatalf("gorm.Open() error = %v", err)
	}

	repo := &liveMonitoringRepository{}
	empID := 479
	date := "2026-07-08"
	jwtCust := "C220010001"
	branch := "pjp"

	queryPattern := `(?s).*SELECT type, latitude, longitude, destination_id, destination_type, destination_name, recorded_at FROM.*WHERE ov\.arrive_at IS NOT NULL.*WHERE ov\.leave_at IS NOT NULL.*`

	rows := sqlmock.NewRows([]string{"type", "latitude", "longitude", "destination_id", "destination_type", "destination_name", "recorded_at"}).
		AddRow("clock_in", 107.0, -6.0, nil, nil, nil, "2026-07-08T06:30:00+07:00").
		AddRow("gps", 107.1, -6.1, nil, nil, nil, "2026-07-08T09:15:00+07:00").
		AddRow("arrive", 107.2, -6.2, int64(101), nil, "Outlet A", "2026-07-08T10:00:00+07:00")

	mock.ExpectQuery(queryPattern).
		WithArgs(empID, date, date, jwtCust, jwtCust, empID, date, date, jwtCust, jwtCust, empID, date, date, jwtCust, jwtCust).
		WillReturnRows(rows)

	results, err := repo.GetUpdateLocations(context.Background(), gormDB, empID, date, jwtCust, branch)
	if err != nil {
		t.Fatalf("GetUpdateLocations() error = %v", err)
	}
	if len(results) != 3 {
		t.Fatalf("GetUpdateLocations() rows = %d, want 3", len(results))
	}
	if results[0].Type != "clock_in" || results[0].Latitude != 107.0 {
		t.Fatalf("row[0] = %+v, want type clock_in and latitude 107.0", results[0])
	}
	if results[2].DestinationType != nil {
		t.Fatalf("row[2].DestinationType = %v, want nil for pjp branch", results[2].DestinationType)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("sqlmock expectations not met: %v", err)
	}
}

func TestGetUpdateLocations_ReturnsTimelineWithNullableDestinationType_PJPPrinciples(t *testing.T) {
	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New() error = %v", err)
	}
	defer sqlDB.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: sqlDB, PreferSimpleProtocol: true}), &gorm.Config{})
	if err != nil {
		t.Fatalf("gorm.Open() error = %v", err)
	}

	repo := &liveMonitoringRepository{}
	empID := 479
	date := "2026-07-08"
	jwtCust := "C220010001"
	branch := "pjp_principles"

	queryPattern := `(?s).*SELECT type, latitude, longitude, destination_id, destination_type, destination_name, recorded_at FROM.*WHERE ov\.arrive_at IS NOT NULL.*WHERE ov\.leave_at IS NOT NULL.*`

	rows := sqlmock.NewRows([]string{"type", "latitude", "longitude", "destination_id", "destination_type", "destination_name", "recorded_at"}).
		AddRow("clock_in", 107.0, -6.0, nil, nil, nil, "2026-07-08T06:30:00+07:00").
		AddRow("arrive", 107.2, -6.2, int64(101), nil, "Outlet A", "2026-07-08T10:00:00+07:00")

	mock.ExpectQuery(queryPattern).
		WithArgs(empID, date, date, jwtCust, jwtCust, empID, date, date, jwtCust, jwtCust, empID, date, date, jwtCust, jwtCust).
		WillReturnRows(rows)

	results, err := repo.GetUpdateLocations(context.Background(), gormDB, empID, date, jwtCust, branch)
	if err != nil {
		t.Fatalf("GetUpdateLocations() error = %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("GetUpdateLocations() rows = %d, want 2", len(results))
	}
	if results[1].Type != "arrive" || results[1].DestinationType != nil {
		t.Fatalf("row[1] = %+v, want type arrive and destination_type nil for both branches", results[1])
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("sqlmock expectations not met: %v", err)
	}
}

func TestGetEmployeeRole_ReturnsEmployeeCustIDWhenFound(t *testing.T) {
	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New() error = %v", err)
	}
	defer sqlDB.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: sqlDB, PreferSimpleProtocol: true}), &gorm.Config{})
	if err != nil {
		t.Fatalf("gorm.Open() error = %v", err)
	}

	repo := &liveMonitoringRepository{}
	empID := 479
	jwtCust := "C220010001"

	queryPattern := `(?s).*SELECT cust_id FROM mst\.m_employee WHERE emp_id = \$1 AND cust_id IN.*`

	rows := sqlmock.NewRows([]string{"cust_id"}).AddRow("C2200100010001")

	mock.ExpectQuery(queryPattern).
		WithArgs(empID, jwtCust, jwtCust).
		WillReturnRows(rows)

	custID, err := repo.GetEmployeeRole(context.Background(), gormDB, empID, jwtCust)
	if err != nil {
		t.Fatalf("GetEmployeeRole() error = %v", err)
	}
	if custID != "C2200100010001" {
		t.Fatalf("GetEmployeeRole() custID = %q, want C2200100010001", custID)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("sqlmock expectations not met: %v", err)
	}
}

func TestGetEmployeeRole_ReturnsNotFoundWhenEmployeeOutsideTenant(t *testing.T) {
	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New() error = %v", err)
	}
	defer sqlDB.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: sqlDB, PreferSimpleProtocol: true}), &gorm.Config{})
	if err != nil {
		t.Fatalf("gorm.Open() error = %v", err)
	}

	repo := &liveMonitoringRepository{}
	empID := 999
	jwtCust := "C220010001"

	queryPattern := `(?s).*SELECT cust_id FROM mst\.m_employee WHERE emp_id = \$1 AND cust_id IN.*`

	rows := sqlmock.NewRows([]string{"cust_id"})

	mock.ExpectQuery(queryPattern).
		WithArgs(empID, jwtCust, jwtCust).
		WillReturnRows(rows)

	_, err = repo.GetEmployeeRole(context.Background(), gormDB, empID, jwtCust)
	if err == nil {
		t.Fatalf("GetEmployeeRole() error = nil, want ErrRecordNotFound")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("sqlmock expectations not met: %v", err)
	}
}
