package repository

import (
	"testing"

	"master/entity"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
)

func setupProductRepositoryTest(t *testing.T) (ProductRepository, sqlmock.Sqlmock, func()) {
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

// newProductDistPriceRows returns a Rows with the full column set expected by ProductDistPrice scan.
func newProductDistPriceRows() *sqlmock.Rows {
	return sqlmock.NewRows([]string{
		"pro_id", "pro_code", "pro_name",
		"unit_id1", "unit_id2", "unit_id3", "unit_id4", "unit_id5",
		"conv_unit2", "conv_unit3", "conv_unit4", "conv_unit5",
		"purch_price1", "purch_price2", "purch_price3", "purch_price4", "purch_price5",
		"sell_price1", "sell_price2", "sell_price3", "sell_price4", "sell_price5",
		"vat",
		"brand_id", "brand_code", "brand_name",
		"pl_id", "pl_code", "pl_name",
		"sbrand1_id", "sbrand1_code", "sbrand1_name",
		"dist_price_id",
	})
}

// addProductDistPriceRow appends a single product row with the given sell_price1.
func addProductDistPriceRow(rows *sqlmock.Rows, sellPrice1 float64) *sqlmock.Rows {
	return rows.AddRow(
		int64(501), "PRO-501", "Product Test",
		"PCS", "BOX", "CTN", nil, nil,
		float32(12), float32(24), nil, nil,
		500.0, 0.0, 0.0, 0.0, 0.0,
		sellPrice1, 0.0, 0.0, 0.0, 0.0,
		nil,
		nil, nil, nil,
		nil, nil, nil,
		nil, nil, nil,
		nil,
	)
}

// expectResolveDistPriceGroupID sets up the mock for the resolveDistPriceGroupID helper query.
func expectResolveDistPriceGroupID(mock sqlmock.Sqlmock, parentCustID string, distributorID int64, distPriceGroupID int) {
	mock.ExpectQuery(`(?s)SELECT dist_price_grp_id\s+FROM mst\.m_distributor\s+WHERE cust_id = \$1\s+AND distributor_id = \$2`).
		WithArgs(parentCustID, distributorID).
		WillReturnRows(sqlmock.NewRows([]string{"dist_price_grp_id"}).AddRow(distPriceGroupID))
}

// expectCountQuery sets up the mock for the COUNT(*) query.
func expectCountQuery(mock sqlmock.Sqlmock, total int) {
	mock.ExpectQuery(`(?s)SELECT COUNT\(\*\) AS total.*FROM mst\.m_product p`).
		WillReturnRows(sqlmock.NewRows([]string{"total"}).AddRow(total))
}

// TestFindAllByDistributorLookupDistPrice_HealthyParentLink_ReadsParentTransactionPrice verifies
// that when a distributor product has a parent_pro_id, the pricing lateral join uses
// pricing_lookup_pro_id (the parent) and the UNION ALL includes the parent generic candidate.
func TestFindAllByDistributorLookupDistPrice_HealthyParentLink_ReadsParentTransactionPrice(t *testing.T) {
	repo, mock, cleanup := setupProductRepositoryTest(t)
	defer cleanup()

	dataFilter := entity.ProductQueryFilter{
		CustId:           "C26002001",
		ParentCustId:     "C26002",
		DistributorID:    102,
		DistPriceGroupId: 2,
		OutletId:         1843,
		OrderDate:        "2026-05-16",
		Page:             1,
		Limit:            10,
	}

	// 1. resolveDistPriceGroupID
	expectResolveDistPriceGroupID(mock, "C26002", int64(102), 2)

	// 2. COUNT
	expectCountQuery(mock, 1)

	// 3. SELECT — regex checks that the UNION ALL lateral join contains both
	//    the parent generic candidate (cust_id = 'C26002', pricing_lookup_pro_id)
	//    and the UNION ALL keyword, and the ORDER BY scope_priority clause.
	sellPrice1 := 999.0
	mock.ExpectQuery(`(?s).*pricing_lookup_pro_id.*UNION ALL.*cust_id = 'C26002'.*paged_products\.pricing_lookup_pro_id.*ORDER BY start_date DESC, scope_priority ASC LIMIT 1.*`).
		WillReturnRows(addProductDistPriceRow(newProductDistPriceRows(), sellPrice1))

	products, total, _, err := repo.FindAllByDistributorLookupDistPrice(dataFilter)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(products) != 1 {
		t.Fatalf("expected 1 product, got %d", len(products))
	}
	if products[0].SellPrice1 == nil || *products[0].SellPrice1 != 999.0 {
		t.Fatalf("expected SellPrice1=999, got %v", products[0].SellPrice1)
	}
	if total != 1 {
		t.Fatalf("expected total=1, got %d", total)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

// TestFindAllByDistributorLookupDistPrice_EmptyResult_WhenNoChildRows verifies that
// when the count query returns 0, the function returns an empty slice without error.
// The SELECT query still runs (repository always executes it regardless of count).
func TestFindAllByDistributorLookupDistPrice_EmptyResult_WhenNoChildRows(t *testing.T) {
	repo, mock, cleanup := setupProductRepositoryTest(t)
	defer cleanup()

	dataFilter := entity.ProductQueryFilter{
		CustId:           "C26002001",
		ParentCustId:     "C26002",
		DistributorID:    102,
		DistPriceGroupId: 2,
		OutletId:         1843,
		OrderDate:        "2026-05-16",
		Page:             1,
		Limit:            10,
	}

	// 1. resolveDistPriceGroupID
	expectResolveDistPriceGroupID(mock, "C26002", int64(102), 2)

	// 2. COUNT returns 0
	expectCountQuery(mock, 0)

	// 3. SELECT still runs, returns no rows
	mock.ExpectQuery(`(?s).*pricing_lookup_pro_id.*UNION ALL.*`).
		WillReturnRows(newProductDistPriceRows()) // no rows added

	products, total, _, err := repo.FindAllByDistributorLookupDistPrice(dataFilter)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(products) != 0 {
		t.Fatalf("expected 0 products, got %d", len(products))
	}
	if total != 0 {
		t.Fatalf("expected total=0, got %d", total)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

// TestFindAllByDistributorLookupDistPrice_QueryIncludesChildAndParentGenericUnion verifies
// the root-cause fix: the SELECT query must include UNION ALL of child generic rows
// (cust_id = child, pro_id = paged_products.pro_id, COALESCE(outlet_id,0)=0) and
// parent generic rows (cust_id = parent, pro_id = paged_products.pricing_lookup_pro_id),
// ordered by start_date DESC, scope_priority ASC so child overrides parent on equal dates.
func TestFindAllByDistributorLookupDistPrice_QueryIncludesChildAndParentGenericUnion(t *testing.T) {
	repo, mock, cleanup := setupProductRepositoryTest(t)
	defer cleanup()

	dataFilter := entity.ProductQueryFilter{
		CustId:           "C260020001",
		ParentCustId:     "C26002",
		DistributorID:    102,
		DistPriceGroupId: 2,
		OutletId:         1843,
		OrderDate:        "2026-05-16",
		Page:             1,
		Limit:            10,
	}

	// 1. resolveDistPriceGroupID
	expectResolveDistPriceGroupID(mock, "C26002", int64(102), 2)

	// 2. COUNT
	expectCountQuery(mock, 1)

	// 3. SELECT — assert all key clauses from the UNION ALL fix are present:
	//    a) child generic candidate: cust_id = 'C260020001', pro_id = paged_products.pro_id, COALESCE(outlet_id, 0) = 0
	//    b) UNION ALL keyword
	//    c) parent generic candidate: cust_id = 'C26002', pro_id = paged_products.pricing_lookup_pro_id
	//    d) ORDER BY start_date DESC, scope_priority ASC LIMIT 1
	mock.ExpectQuery(
		`(?s)` +
			`.*cust_id = 'C260020001'` + // child generic cust_id
			`.*pro_id = paged_products\.pro_id` + // child generic pro_id
			`.*COALESCE\(outlet_id, 0\) = 0` + // child generic outlet filter
			`.*UNION ALL` + // union keyword
			`.*cust_id = 'C26002'` + // parent generic cust_id
			`.*paged_products\.pricing_lookup_pro_id` + // parent generic pro_id
			`.*ORDER BY start_date DESC, scope_priority ASC LIMIT 1` + // precedence order
			`.*`,
	).WillReturnRows(addProductDistPriceRow(newProductDistPriceRows(), 150000.0))

	products, total, _, err := repo.FindAllByDistributorLookupDistPrice(dataFilter)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(products) != 1 {
		t.Fatalf("expected 1 product, got %d", len(products))
	}
	if products[0].SellPrice1 == nil || *products[0].SellPrice1 != 150000.0 {
		t.Fatalf("expected SellPrice1=150000, got %v", products[0].SellPrice1)
	}
	if total != 1 {
		t.Fatalf("expected total=1, got %d", total)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}
