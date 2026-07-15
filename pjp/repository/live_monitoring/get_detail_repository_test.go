package live_monitoring

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestGetSubmittedSurveyData_ReturnsSubmittedSurveyGroupedByTitleAndOutlet(t *testing.T) {
	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New() error = %v", err)
	}
	defer sqlDB.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: sqlDB, PreferSimpleProtocol: true}), &gorm.Config{})
	if err != nil {
		t.Fatalf("gorm.Open() error = %v", err)
	}

	repo := NewLiveMonitoringRepository()
	custIDs := []string{"CUST-1", "CUST-2"}
	date := "2026-05-28"
	empID := 210

	queryPattern := `(?s).*FROM mst\.survey_answer sa.*JOIN mst\.m_survey ms ON ms\.survey_id = sa\.survey_id.*JOIN mst\.m_outlet mo ON mo\.outlet_id = sa\.outlet_id AND mo\.cust_id = sa\.cust_id.*sa\.cust_id IN \(\$1,\$2\).*DATE\(sa\.answer_date\) = \$3.*sa\.emp_id = \$4.*sa\.status = \$5.*GROUP BY ms\.survey_title, mo\.outlet_code, mo\.outlet_name.*ORDER BY ms\.survey_title ASC, mo\.outlet_code ASC, mo\.outlet_name ASC.*`

	rows := sqlmock.NewRows([]string{"submission", "survey_title", "outlet_code", "outlet_name"}).
		AddRow(int64(3), "Store Audit", "OUT-01", "Outlet A")

	mock.ExpectQuery(queryPattern).
		WithArgs(custIDs[0], custIDs[1], date, empID, "Submitted").
		WillReturnRows(rows)

	results, err := repo.GetSubmittedSurveyData(context.Background(), gormDB, custIDs, date, empID)
	if err != nil {
		t.Fatalf("GetSubmittedSurveyData() error = %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("GetSubmittedSurveyData() rows = %d, want 1", len(results))
	}
	if results[0].Submission != 3 || results[0].OutletCode != "OUT-01" {
		t.Fatalf("row[0] = %+v, want submission 3 and outlet OUT-01", results[0])
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("sqlmock expectations not met: %v", err)
	}
}

func TestGetCollections_AllocatesCollectionPerInvoicePerOutlet(t *testing.T) {
	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New() error = %v", err)
	}
	defer sqlDB.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn:                 sqlDB,
		PreferSimpleProtocol: true,
	}), &gorm.Config{})
	if err != nil {
		t.Fatalf("gorm.Open() error = %v", err)
	}

	repo := NewLiveMonitoringRepository()
	custIDs := []string{"CUST-1", "CUST-2"}
	date := "2026-05-01"
	empID := 77

	queryPattern := `(?s).*GROUP BY deposit_no, cust_id, invoice_no.*ppi.invoice_no = dd.invoice_no.*GROUP BY mo.outlet_id, mo.outlet_code, mo.outlet_name.*`

	rows := sqlmock.NewRows([]string{"outlet_id", "outlet_code", "outlet_name", "collection_total"}).
		AddRow(101, "TOKO-A", "Toko A", 1000000.0).
		AddRow(102, "TOKO-B", "Toko B", 500000.0)

	mock.ExpectQuery(queryPattern).
		WithArgs(custIDs[0], custIDs[1], date, empID).
		WillReturnRows(rows)

	results, err := repo.GetCollections(context.Background(), gormDB, custIDs, date, empID)
	if err != nil {
		t.Fatalf("GetCollections() error = %v", err)
	}

	if len(results) != 2 {
		t.Fatalf("GetCollections() rows = %d, want 2", len(results))
	}

	if results[0].OutletName != "Toko A" || results[0].CollectionTotal != 1000000 {
		t.Fatalf("row[0] = %+v, want outlet Toko A with total 1000000", results[0])
	}

	if results[1].OutletName != "Toko B" || results[1].CollectionTotal != 500000 {
		t.Fatalf("row[1] = %+v, want outlet Toko B with total 500000", results[1])
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("sqlmock expectations not met: %v", err)
	}
}
