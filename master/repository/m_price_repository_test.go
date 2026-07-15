package repository

import (
	"reflect"
	"testing"
	"time"

	"master/entity"
	"master/model"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
)

func setupMPriceRepositoryTest(t *testing.T) (MPriceRepository, sqlmock.Sqlmock, func()) {
	t.Helper()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := NewMPriceRepository(sqlxDB)

	cleanup := func() {
		_ = db.Close()
	}

	return repo, mock, cleanup
}

func TestMPriceRepository_FindAllByCustID_AppliesSpecFiltersAndScansDistributorIDs(t *testing.T) {
	repo, mock, cleanup := setupMPriceRepositoryTest(t)
	defer cleanup()

	filter := entity.MPriceQueryFilter{
		Page:               1,
		Limit:              5,
		Query:              "MESSI",
		Sort:               "created_date:desc",
		Status:             []int{1, 10},
		EffectiveDateStart: "2026-05-01",
		EffectiveDateEnd:   "2026-05-31",
		DistributorIDs:     []int64{11, 22},
	}

	mock.ExpectQuery(`(?s)SELECT COUNT\(prc\.price_id\) AS total.*WHERE prc\.cust_id = 'CUST001' AND prc\.is_del IS FALSE AND prc\.effective_date::date BETWEEN '2026-05-01' AND '2026-05-31' AND \(prd\.pro_name ILIKE '%MESSI%' OR prd\.pro_code ILIKE '%MESSI%'\) AND prc\.status IN \(1,10\) AND prc\.distributor_ids && ARRAY\[11,22\]::bigint\[\]`).
		WillReturnRows(sqlmock.NewRows([]string{"total"}).AddRow(6))

	effectiveDate := time.Date(2026, 5, 5, 0, 0, 0, 0, time.UTC)
	now := time.Date(2026, 5, 4, 10, 0, 0, 0, time.UTC)
	mock.ExpectQuery(`(?s)SELECT .*prc\.distributor_ids.*WHERE prc\.cust_id = 'CUST001' AND prc\.is_del IS FALSE AND prc\.effective_date::date BETWEEN '2026-05-01' AND '2026-05-31' AND \(prd\.pro_name ILIKE '%MESSI%' OR prd\.pro_code ILIKE '%MESSI%'\) AND prc\.status IN \(1,10\) AND prc\.distributor_ids && ARRAY\[11,22\]::bigint\[\].*ORDER BY prc\.created_at DESC LIMIT 5 OFFSET 0`).
		WillReturnRows(sqlmock.NewRows([]string{
			"cust_id", "price_id", "status", "pro_id", "pro_code", "pro_name",
			"unit_name1", "unit_name2", "unit_name3", "coverage", "effective_date",
			"unit_id1", "unit_id2", "unit_id3", "conv_unit2", "conv_unit3",
			"purch_price1", "purch_price2", "purch_price3",
			"sell_price1", "sell_price2", "sell_price3",
			"new_purch_price1", "new_purch_price2", "new_purch_price3",
			"new_sell_price1", "new_sell_price2", "new_sell_price3",
			"created_by_id", "created_by", "created_at", "updated_by_id", "updated_by", "updated_at",
			"distributor_ids",
		}).AddRow(
			"CUST001", "PRICE-001", 1, int64(77), "PRO-77", "Action Figure Messi",
			"PCS", "BOX", "CRT", "D", effectiveDate,
			"PCS", "BOX", "CRT", 2, 4,
			10.0, 20.0, 30.0,
			40.0, 50.0, 60.0,
			100.0, 200.0, 300.0,
			400.0, 500.0, 600.0,
			int64(7), "Creator", now, int64(8), "Updater", now,
			"{11,22}",
		))

	rows, total, lastPage, err := repo.FindAllByCustID(filter, "CUST001")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if total != 6 || lastPage != 2 {
		t.Fatalf("expected total=6 lastPage=2, got total=%d lastPage=%d", total, lastPage)
	}
	if len(rows) != 1 {
		t.Fatalf("expected 1 row, got %d", len(rows))
	}
	if !reflect.DeepEqual([]int64(rows[0].DistributorIDs), []int64{11, 22}) {
		t.Fatalf("expected distributor_ids [11 22], got %+v", rows[0].DistributorIDs)
	}
	if rows[0].ProductCode != "PRO-77" {
		t.Fatalf("expected product code PRO-77, got %q", rows[0].ProductCode)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestMPriceRepository_Delete_SoftDeletesUsingAuditColumns(t *testing.T) {
	repo, mock, cleanup := setupMPriceRepositoryTest(t)
	defer cleanup()

	mock.ExpectExec(`(?s)UPDATE mst\.m_price.*SET is_del = true,.*deleted_at = CURRENT_TIMESTAMP,.*deleted_by = .*WHERE is_del IS FALSE.*cust_id = .*price_id = .*`).
		WithArgs(int64(99), "CUST001", "PRICE-009").
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := repo.Delete("CUST001", "PRICE-009", 99)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestMPriceRepository_FindOneByMPriceIDAndCustID_AllowsNullUnitNames(t *testing.T) {
	repo, mock, cleanup := setupMPriceRepositoryTest(t)
	defer cleanup()

	effectiveDate := time.Date(2026, 5, 5, 0, 0, 0, 0, time.UTC)
	createdAt := time.Date(2026, 5, 4, 10, 0, 0, 0, time.UTC)

	mock.ExpectQuery(`(?s)SELECT .*COALESCE\(un1\.unit_name, ''\) AS unit_name1.*COALESCE\(un3\.unit_name, ''\) AS unit_name3.*WHERE prc\.price_id = \$1.*prc\.cust_id = \$2`).
		WithArgs("PRICE-001", "CUST001").
		WillReturnRows(sqlmock.NewRows([]string{
			"cust_id", "price_id", "coverage", "effective_date", "pro_id",
			"unit_id1", "unit_id2", "unit_id3", "conv_unit2", "conv_unit3",
			"purch_price1", "purch_price2", "purch_price3",
			"sell_price1", "sell_price2", "sell_price3",
			"new_purch_price1", "new_purch_price2", "new_purch_price3",
			"new_sell_price1", "new_sell_price2", "new_sell_price3",
			"status", "created_by_id", "created_by", "created_at",
			"updated_by_id", "updated_by", "updated_at", "distributor_ids",
			"pro_code", "pro_name", "unit_name1", "unit_name2", "unit_name3",
		}).AddRow(
			"CUST001", "PRICE-001", "D", effectiveDate, int64(77),
			"PCS", "BOX", "CRT", 2, 4,
			10.0, 20.0, 30.0,
			40.0, 50.0, 60.0,
			100.0, 200.0, 300.0,
			400.0, 500.0, 600.0,
			1, int64(9), "Creator", createdAt,
			int64(10), "Updater", createdAt, "{11}",
			"PRO-77", "Produk 77", "", "", "",
		))

	detail, err := repo.FindOneByMPriceIDAndCustID(entity.DetailMPriceParams{
		CustID:  "CUST001",
		PriceID: "PRICE-001",
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if detail.UnitName1 != "" || detail.UnitName2 != "" || detail.UnitName3 != "" {
		t.Fatalf("expected empty unit names when joined values are null, got %+v", []string{detail.UnitName1, detail.UnitName2, detail.UnitName3})
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestMPriceRepository_FindAffectedDistributorProductIDs_PrincipalScope(t *testing.T) {
	repo, mock, cleanup := setupMPriceRepositoryTest(t)
	defer cleanup()

	mock.ExpectQuery(`(?s)SELECT DISTINCT distributor_id.*FROM mst\.m_product.*WHERE is_del IS FALSE.*distributor_id IS NOT NULL.*parent_pro_id = \?`).
		WithArgs(int64(77)).
		WillReturnRows(sqlmock.NewRows([]string{"distributor_id"}).
			AddRow(int64(11)).
			AddRow(int64(22)))

	distributorIDs, err := repo.FindAffectedDistributorProductIDs(model.MPriceDetail{
		CustID: "PARENT001",
		ProID:  int64(77),
	}, "PARENT001")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !reflect.DeepEqual(distributorIDs, []int64{11, 22}) {
		t.Fatalf("expected affected distributor IDs [11 22], got %+v", distributorIDs)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestMPriceRepository_FindUpdatedDistributorProductIDs_MatchesPublishedPrices(t *testing.T) {
	repo, mock, cleanup := setupMPriceRepositoryTest(t)
	defer cleanup()

	mock.ExpectQuery(`(?s)SELECT DISTINCT distributor_id.*FROM mst\.m_product.*distributor_id IN \(\?, \?\).*purch_price1 = \?.*sell_price3 = \?.*parent_pro_id = \?`).
		WithArgs(int64(11), int64(22), 101.0, 202.0, 303.0, 404.0, 505.0, 606.0, int64(77)).
		WillReturnRows(sqlmock.NewRows([]string{"distributor_id"}).
			AddRow(int64(11)))

	distributorIDs, err := repo.FindUpdatedDistributorProductIDs(model.MPriceDetail{
		CustID:         "PARENT001",
		ProID:          int64(77),
		NewPurchPrice1: 101,
		NewPurchPrice2: 202,
		NewPurchPrice3: 303,
		NewSellPrice1:  404,
		NewSellPrice2:  505,
		NewSellPrice3:  606,
	}, "PARENT001", []int64{11, 22})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !reflect.DeepEqual(distributorIDs, []int64{11}) {
		t.Fatalf("expected updated distributor IDs [11], got %+v", distributorIDs)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestMPriceRepository_FindBrokenDistributorChildLinks_ReturnsBrokenDistributors(t *testing.T) {
	repo, mock, cleanup := setupMPriceRepositoryTest(t)
	defer cleanup()

	// distributor 11: parent_pro_id = 77 (healthy, not returned)
	// distributor 22: parent_pro_id = 0  (broken, returned)
	// distributor 33: no row at all       (not returned)
	// JOIN mst.m_distributor filters by parent_cust_id = "PARENT001"
	mock.ExpectQuery(`(?s)SELECT DISTINCT p\.distributor_id.*FROM mst\.m_product p.*JOIN mst\.m_distributor d.*d\.parent_cust_id = \?.*p\.distributor_id IN \(\?, \?, \?\).*p\.is_del = false.*COALESCE\(p\.parent_pro_id, 0\) != \?`).
		WithArgs("PARENT001", int64(11), int64(22), int64(33), int64(77)).
		WillReturnRows(sqlmock.NewRows([]string{"distributor_id"}).
			AddRow(int64(22)))

	result, err := repo.FindBrokenDistributorChildLinks(77, "PARENT001", []int64{11, 22, 33})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !reflect.DeepEqual(result, []int64{22}) {
		t.Fatalf("expected broken distributor IDs [22], got %+v", result)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestMPriceRepository_FindBrokenDistributorChildLinks_EmptyWhenNoDistributorIDs(t *testing.T) {
	repo, _, cleanup := setupMPriceRepositoryTest(t)
	defer cleanup()

	// Empty distributorIDs must short-circuit before hitting the DB.
	result, err := repo.FindBrokenDistributorChildLinks(77, "PARENT001", []int64{})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(result) != 0 {
		t.Fatalf("expected empty result, got %+v", result)
	}
}
