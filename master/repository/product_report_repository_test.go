package repository

import (
	"testing"

	"master/entity"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
)

func setupProductReportRepositoryTest(t *testing.T) (ProductRepository, sqlmock.Sqlmock, func()) {
	t.Helper()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := NewProductRepository(sqlxDB)

	cleanup := func() {
		_ = db.Close()
	}

	return repo, mock, cleanup
}

func TestProductReportRepository_CountAndDataQuery(t *testing.T) {
	repo, mock, cleanup := setupProductReportRepositoryTest(t)
	defer cleanup()

	filter := entity.ProductReportQueryFilter{
		CustIDs:   []string{"C26002", "C260020001"},
		Query:     "ABC",
		Page:      1,
		Limit:     20,
		SortBy:    "pro_name",
		SortOrder: "asc",
	}

	// Count query
	mock.ExpectQuery(`(?s)SELECT COUNT\(\*\).*FROM mst\.m_product mp.*LEFT JOIN \(.*SELECT cust_id, BOOL_OR\(COALESCE\(allow_upload_secondary_sales, false\)\) AS allow_upload_secondary_sales.*FROM mst\.m_distributor.*GROUP BY cust_id.*\) md.*LEFT JOIN mst\.m_product parent.*WHERE mp\.is_del = false AND mp\.is_active = true AND mp\.cust_id IN \(\?, \?\).*AND \(mp\.pro_name ILIKE \? OR mp\.pro_code ILIKE \?\)`).
		WithArgs("C26002", "C260020001", "%ABC%", "%ABC%").
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(5))

	// Data query
	mock.ExpectQuery(`(?s)SELECT.*CASE WHEN.*parent\.cust_id.*END AS cust_id.*FROM mst\.m_product mp.*LEFT JOIN \(.*SELECT cust_id, BOOL_OR\(COALESCE\(allow_upload_secondary_sales, false\)\) AS allow_upload_secondary_sales.*FROM mst\.m_distributor.*GROUP BY cust_id.*\) md.*LEFT JOIN mst\.m_product parent.*WHERE mp\.is_del = false AND mp\.is_active = true AND mp\.cust_id IN \(\?, \?\).*AND \(mp\.pro_name ILIKE \? OR mp\.pro_code ILIKE \?\).*ORDER BY pro_name ASC LIMIT \? OFFSET \?`).
		WithArgs("C26002", "C260020001", "%ABC%", "%ABC%", 20, 0).
		WillReturnRows(sqlmock.NewRows([]string{
			"cust_id", "pro_id", "pro_code", "pro_name",
			"original_cust_id", "original_pro_id", "original_pro_code", "original_parent_pro_id",
			"type",
		}))

	rows, total, lastPage, err := repo.ReportList(filter)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if total != 5 {
		t.Fatalf("expected total=5, got %d", total)
	}
	if lastPage != 1 {
		t.Fatalf("expected lastPage=1, got %d", lastPage)
	}
	if len(rows) != 0 {
		t.Fatalf("expected 0 rows, got %d", len(rows))
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestProductReportRepository_ParentJoinSQL(t *testing.T) {
	repo, mock, cleanup := setupProductReportRepositoryTest(t)
	defer cleanup()

	filter := entity.ProductReportQueryFilter{
		CustIDs:   []string{"C260020001"},
		Page:      1,
		Limit:     20,
		SortBy:    "pro_name",
		SortOrder: "asc",
	}

	mock.ExpectQuery(`(?s)SELECT COUNT\(\*\).*FROM mst\.m_product mp.*LEFT JOIN \(.*SELECT cust_id, BOOL_OR\(COALESCE\(allow_upload_secondary_sales, false\)\) AS allow_upload_secondary_sales.*FROM mst\.m_distributor.*GROUP BY cust_id.*\) md.*LEFT JOIN mst\.m_product parent.*WHERE mp\.is_del = false AND mp\.is_active = true AND mp\.cust_id IN \(\?\)`).
		WithArgs("C260020001").
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

	mock.ExpectQuery(`(?s)SELECT.*parent\.pro_id IS NOT NULL.*parent\.pro_id = mp\.parent_pro_id.*parent\.cust_id = LEFT\(mp\.cust_id, 6\)`).
		WithArgs("C260020001", 20, 0).
		WillReturnRows(sqlmock.NewRows([]string{
			"cust_id", "pro_id", "pro_code", "pro_name",
			"original_cust_id", "original_pro_id", "original_pro_code", "original_parent_pro_id",
			"type",
		}))

	_, _, _, err := repo.ReportList(filter)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestProductReportRepository_NormalizationBranches(t *testing.T) {
	repo, mock, cleanup := setupProductReportRepositoryTest(t)
	defer cleanup()

	filter := entity.ProductReportQueryFilter{CustIDs: []string{"C260020001"}, Page: 1, Limit: 20, SortBy: "pro_id", SortOrder: "asc"}
	mock.ExpectQuery(`SELECT COUNT\(\*\).*FROM mst\.m_product mp`).WithArgs("C260020001").WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(2))
	mock.ExpectQuery(`(?s)SELECT.*parent\.pro_id IS NOT NULL.*ORDER BY pro_id ASC LIMIT \? OFFSET \?`).
		WithArgs("C260020001", 20, 0).
		WillReturnRows(sqlmock.NewRows([]string{
			"cust_id", "pro_id", "pro_code", "pro_name", "original_cust_id", "original_pro_id", "original_pro_code", "original_parent_pro_id", "type",
		}).AddRow("C260020001", int64(11), "CHILD", "Child", nil, nil, nil, nil, "Product Mapping").
			AddRow("C26002", int64(7), "PARENT", "Parent", "C260020001", int64(12), "MAPPED", int64(7), "Product Mapping"))

	rows, _, _, err := repo.ReportList(filter)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(rows) != 2 {
		t.Fatalf("expected 2 rows, got %d", len(rows))
	}
	if rows[0].ProductID != 11 || rows[0].OriginalProductID != nil || rows[0].OriginalCustID != nil || rows[0].OriginalProductCode != nil || rows[0].OriginalParentID != nil {
		t.Fatalf("missing-parent branch mismatch: %+v", rows[0])
	}
	if rows[1].CustID != "C26002" || rows[1].ProductID != 7 || rows[1].ProductCode != "PARENT" || rows[1].ProductName != "Parent" || rows[1].OriginalCustID == nil || *rows[1].OriginalCustID != "C260020001" || rows[1].OriginalProductID == nil || *rows[1].OriginalProductID != 12 || rows[1].OriginalProductCode == nil || *rows[1].OriginalProductCode != "MAPPED" || rows[1].OriginalParentID == nil || *rows[1].OriginalParentID != 7 {
		t.Fatalf("eligible-parent branch mismatch: %+v", rows[1])
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestProductReportRepository_ScanCategories(t *testing.T) {
	tests := []struct {
		name string
		custID string
		proID int64
		proCode string
		proName string
		origCustID *string
		origProID *int64
		origProCode *string
		origParentID *int64
		expectedType string
	}{
		{name: "principal_own_products", custID: "C26002", proID: 1, proCode: "P001", proName: "Principal Product", expectedType: "Own Products"},
		{name: "distributor_own_product", custID: "C260020001", proID: 2, proCode: "D001", proName: "Distributor Own", expectedType: "Own Products"},
		{name: "product_assigned", custID: "C260020001", proID: 3, proCode: "A001", proName: "Assigned Product", expectedType: "Product Assigned"},
		{name: "mapping_enabled", custID: "C26002", proID: 10, proCode: "M001", proName: "Mapped Product", origCustID: reportStrPtr("C260020001"), origProID: reportInt64Ptr(3), origProCode: reportStrPtr("A001"), origParentID: reportInt64Ptr(1), expectedType: "Product Mapping"},
		{name: "mapping_disabled", custID: "C260020001", proID: 4, proCode: "MD001", proName: "Mapping Disabled", expectedType: "Product Mapping"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo, mock, cleanup := setupProductReportRepositoryTest(t)
			defer cleanup()
			filter := entity.ProductReportQueryFilter{CustIDs: []string{tt.custID}, Page: 1, Limit: 20, SortBy: "pro_name", SortOrder: "asc"}
			mock.ExpectQuery(`(?s)SELECT COUNT\(\*\).*FROM mst\.m_product mp.*WHERE mp\.is_del = false AND mp\.is_active = true AND mp\.cust_id IN \(\?\)`).WithArgs(tt.custID).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))
			mock.ExpectQuery(`(?s)SELECT.*FROM mst\.m_product mp.*ORDER BY pro_name ASC LIMIT \? OFFSET \?`).WithArgs(tt.custID, 20, 0).WillReturnRows(sqlmock.NewRows([]string{"cust_id", "pro_id", "pro_code", "pro_name", "original_cust_id", "original_pro_id", "original_pro_code", "original_parent_pro_id", "type"}).AddRow(tt.custID, tt.proID, tt.proCode, tt.proName, tt.origCustID, tt.origProID, tt.origProCode, tt.origParentID, tt.expectedType))
			rows, total, lastPage, err := repo.ReportList(filter)
			if err != nil { t.Fatalf("expected no error, got %v", err) }
			if total != 1 || lastPage != 1 || len(rows) != 1 { t.Fatalf("unexpected result: total=%d lastPage=%d rows=%d", total, lastPage, len(rows)) }
			if err := mock.ExpectationsWereMet(); err != nil { t.Fatalf("unmet expectations: %v", err) }
		})
	}
}

func TestProductReportRepository_DistributorAggregationPreservesCardinality(t *testing.T) {
	repo, mock, cleanup := setupProductReportRepositoryTest(t)
	defer cleanup()
	filter := entity.ProductReportQueryFilter{CustIDs: []string{"C22001"}, Page: 1, Limit: 20, SortBy: "pro_id", SortOrder: "asc"}
	derivedJoin := `(?s)LEFT JOIN \(\s*SELECT cust_id, BOOL_OR\(COALESCE\(allow_upload_secondary_sales, false\)\) AS allow_upload_secondary_sales\s*FROM mst\.m_distributor\s*GROUP BY cust_id\s*\) md ON md\.cust_id = mp\.cust_id`
	mock.ExpectQuery(`SELECT COUNT\(\*\).*` + derivedJoin + `.*WHERE mp\.is_del = false AND mp\.is_active = true AND mp\.cust_id IN \(\?\)`).WithArgs("C22001").WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))
	mock.ExpectQuery(`SELECT.*` + derivedJoin + `.*WHERE mp\.is_del = false AND mp\.is_active = true AND mp\.cust_id IN \(\?\).*ORDER BY pro_id ASC LIMIT \? OFFSET \?`).WithArgs("C22001", 20, 0).WillReturnRows(sqlmock.NewRows([]string{"cust_id", "pro_id", "pro_code", "pro_name", "original_cust_id", "original_pro_id", "original_pro_code", "original_parent_pro_id", "type"}).AddRow("C22001", 495, "P495", "Product 495", nil, nil, nil, nil, "Own Products"))
	rows, total, _, err := repo.ReportList(filter)
	if err != nil { t.Fatalf("expected no error, got %v", err) }
	if total != 1 || len(rows) != 1 { t.Fatalf("expected one count and one row, got total=%d rows=%d", total, len(rows)) }
	if err := mock.ExpectationsWereMet(); err != nil { t.Fatalf("unmet expectations: %v", err) }
}

func reportStrPtr(s string) *string { return &s }
func reportInt64Ptr(i int64) *int64 { return &i }
