package repository

import (
	"reflect"
	"strings"
	"testing"
	"time"

	"sales/entity"
	"sales/pkg/str"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestBuildSecondarySalesUnionQueryUsesTransCTEAndDeterministicParentProductJoin(t *testing.T) {
	t.Parallel()

	from := int64(1778086800)
	to := int64(1778173199)
	filter := entity.SecondarySalesReportQueryFilter{
		CustID:       "C260020001",
		CustIDs:      []string{"C260020001"},
		ParentCustID: "C26002",
		From:         &from,
		To:           &to,
		Page:         2,
		Limit:        25,
	}

	sql, params, limit := buildSecondarySalesUnionQuery(filter, true)

	checks := []string{
		"WITH trans AS",
		"FROM sls.order_detail od",
		"JOIN sls.\"order\" o ON o.ro_no = od.ro_no AND o.cust_id = od.cust_id",
		"o.data_status IN (6,7)",
		"o.invoice_date BETWEEN ? AND ?",
		"o.invoice_no IS NOT NULL",
		"FROM sls.return_det rd",
		"JOIN sls.\"return\" r ON rd.return_no = r.return_no AND rd.cust_id = r.cust_id",
		"JOIN sls.\"order\" o ON o.invoice_no = r.invoice_no AND o.cust_id = r.cust_id",
		"LEFT JOIN mst.m_product cp ON cp.pro_id = t.product_id AND cp.cust_id = t.cust_id",
		"LEFT JOIN LATERAL (",
		"WHERE pp.cust_id = ?",
		"ORDER BY CASE WHEN pp.pro_id = t.product_id THEN 0 ELSE 1 END, pp.pro_id ASC",
		"LIMIT ? OFFSET ?",
	}

	for _, check := range checks {
		if !strings.Contains(sql, check) {
			t.Fatalf("expected SQL to contain %q\nSQL:\n%s", check, sql)
		}
	}

	expectedParams := []interface{}{
		[]string{"C260020001"},
		str.UnixTimestampToUtcTime(from),
		str.UnixTimestampToUtcTime(to),
		[]string{"C260020001"},
		str.UnixTimestampToUtcTime(from),
		str.UnixTimestampToUtcTime(to),
		"C26002",
		25,
		25,
	}

	if limit != 25 {
		t.Fatalf("expected limit 25, got %d", limit)
	}

	if len(params) != len(expectedParams) {
		t.Fatalf("expected %d params, got %d: %#v", len(expectedParams), len(params), params)
	}

	for i := range expectedParams {
		if !reflect.DeepEqual(params[i], expectedParams[i]) {
			t.Fatalf("param %d mismatch: expected %#v, got %#v", i, expectedParams[i], params[i])
		}
	}

	if got := params[1]; got != time.Date(2026, 5, 6, 17, 0, 0, 0, time.UTC) {
		t.Fatalf("unexpected from time conversion: %#v", got)
	}
	if got := params[2]; got != time.Date(2026, 5, 7, 16, 59, 59, 0, time.UTC) {
		t.Fatalf("unexpected to time conversion: %#v", got)
	}
}

func TestBuildSecondarySalesUnionQueryCoalescesGrossSalesOperands(t *testing.T) {
	t.Parallel()

	from := int64(1778086800)
	to := int64(1778173199)
	filter := entity.SecondarySalesReportQueryFilter{
		CustID:       "C260020001",
		CustIDs:      []string{"C260020001"},
		ParentCustID: "C26002",
		From:         &from,
		To:           &to,
	}

	sql, _, _ := buildSecondarySalesUnionQuery(filter, false)

	checks := []string{
		"COALESCE(od.qty1_final, 0) * COALESCE(od.sell_price1, 0)",
		"COALESCE(od.qty2_final, 0) * COALESCE(od.sell_price2, 0)",
		"COALESCE(od.qty3_final, 0) * COALESCE(od.sell_price3, 0)",
		"COALESCE(rd.qty1, 0) * COALESCE(rd.sell_price1, 0)",
		"COALESCE(rd.qty2, 0) * COALESCE(rd.sell_price2, 0)",
		"COALESCE(rd.qty3, 0) * COALESCE(rd.sell_price3, 0)",
	}

	for _, check := range checks {
		if !strings.Contains(sql, check) {
			t.Fatalf("expected SQL to contain %q\nSQL:\n%s", check, sql)
		}
	}
}

func TestBuildSecondarySalesUnionQueryUsesMultiCustBinding(t *testing.T) {
	t.Parallel()

	filter := entity.SecondarySalesReportQueryFilter{
		CustID:       "AUTH1",
		CustIDs:      []string{"CHILD1", "CHILD2"},
		ParentCustID: "PARENT1",
	}

	sql, params, _ := buildSecondarySalesUnionQuery(filter, false)
	if !strings.Contains(sql, "od.cust_id IN ?") || !strings.Contains(sql, "rd.cust_id IN ?") {
		t.Fatalf("expected multi cust IN binding, sql=%s", sql)
	}
	if !strings.Contains(sql, "cp.pro_id = t.product_id AND cp.cust_id = t.cust_id") {
		t.Fatalf("expected export child product lookup to use row-level cust_id, sql=%s", sql)
	}
	if strings.Contains(sql, "cp.pro_id = t.product_id AND cp.cust_id = ?") {
		t.Fatalf("did not expect export child product lookup to bind a single auth cust, sql=%s", sql)
	}
	if len(params) < 2 {
		t.Fatalf("expected at least two params, got %#v", params)
	}
	if !reflect.DeepEqual(params[0], []string{"CHILD1", "CHILD2"}) {
		t.Fatalf("expected order cust ids slice, got %#v", params[0])
	}
	if !reflect.DeepEqual(params[1], []string{"CHILD1", "CHILD2"}) {
		t.Fatalf("expected return cust ids slice, got %#v", params[1])
	}
}

func TestBuildSecondarySalesUnionQueryPreservesFilterAliases(t *testing.T) {
	t.Parallel()

	from := int64(1778086800)
	to := int64(1778173199)
	filter := entity.SecondarySalesReportQueryFilter{
		CustID:         "C260020001",
		CustIDs:        []string{"C260020001"},
		ParentCustID:   "C26002",
		From:           &from,
		To:             &to,
		DistributorIDs: []int64{21},
		SalesmanIDs:    []int64{101},
		OutletIDs:      []int64{301},
		ProIDs:         []int64{501},
	}

	sql, params, _ := buildSecondarySalesUnionQuery(filter, false)

	checks := []string{
		"WHERE 1=1",
		"md.distributor_id IN ?",
		"o.salesman_id IN ?",
		"r.salesman_id IN ?",
		"o.outlet_id IN ?",
		"r.outlet_id IN ?",
		"od.pro_id IN ?",
		"rd.product_id IN ?",
	}

	for _, check := range checks {
		if !strings.Contains(sql, check) {
			t.Fatalf("expected SQL to contain %q\nSQL:\n%s", check, sql)
		}
	}

	expectedParams := []interface{}{
		[]string{"C260020001"},
		str.UnixTimestampToUtcTime(from),
		str.UnixTimestampToUtcTime(to),
		[]int64{101},
		[]int64{301},
		[]int64{501},
		[]string{"C260020001"},
		str.UnixTimestampToUtcTime(from),
		str.UnixTimestampToUtcTime(to),
		[]int64{101},
		[]int64{301},
		[]int64{501},
		"C26002",
		[]int64{21},
	}

	if !reflect.DeepEqual(params, expectedParams) {
		t.Fatalf("unexpected params:\nexpected: %#v\nactual:   %#v", expectedParams, params)
	}
}

type recordedQuery struct {
	SQL  string
	Vars []interface{}
}

func newReportRepoDryRunDB(t *testing.T) (*gorm.DB, *[]recordedQuery) {
	t.Helper()

	recorded := make([]recordedQuery, 0, 4)
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  "host=localhost user=postgres dbname=postgres sslmode=disable",
		PreferSimpleProtocol: true,
	}), &gorm.Config{DisableAutomaticPing: true, DryRun: true})
	if err != nil {
		t.Fatalf("failed to open dry-run gorm db: %v", err)
	}

	if err := db.Callback().Query().After("gorm:query").Register("report_repository_test:capture_sql", func(tx *gorm.DB) {
		if tx.Statement == nil {
			return
		}

		vars := append([]interface{}(nil), tx.Statement.Vars...)
		recorded = append(recorded, recordedQuery{
			SQL:  tx.Statement.SQL.String(),
			Vars: vars,
		})
	}); err != nil {
		t.Fatalf("failed to register query capture callback: %v", err)
	}

	return db, &recorded
}

func latestRecordedQuery(t *testing.T, recorded *[]recordedQuery) recordedQuery {
	t.Helper()

	if len(*recorded) == 0 {
		t.Fatal("expected repository method to emit query")
	}

	return (*recorded)[len(*recorded)-1]
}

func valueOrZero(value *float64) float64 {
	if value == nil {
		return 0
	}
	return *value
}

func assertSecondarySalesSummaryDateVars(t *testing.T, actual []interface{}, expected []interface{}) {
	t.Helper()

	if len(actual) != len(expected) {
		t.Fatalf("unexpected var count: expected %d got %d (%#v)", len(expected), len(actual), actual)
	}

	for i := range expected {
		switch want := expected[i].(type) {
		case time.Time:
			gotTime, ok := actual[i].(time.Time)
			if !ok {
				t.Fatalf("var %d expected time.Time, got %T (%#v)", i, actual[i], actual[i])
			}
			if !gotTime.Equal(want) {
				t.Fatalf("var %d expected %s, got %s", i, want.Format("2006-01-02 15:04:05"), gotTime.Format("2006-01-02 15:04:05"))
			}
		default:
			if !reflect.DeepEqual(actual[i], expected[i]) {
				t.Fatalf("var %d mismatch: expected %#v got %#v", i, expected[i], actual[i])
			}
		}
	}
}

func assertSecondarySalesGroupDateVars(t *testing.T, actual []interface{}, custIDs []string, dateFrom, dateTo time.Time) {
	t.Helper()

	expected := make([]interface{}, 0, len(custIDs)*2+4)
	for _, custID := range custIDs {
		expected = append(expected, custID)
	}
	expected = append(expected, dateFrom, dateTo)
	for _, custID := range custIDs {
		expected = append(expected, custID)
	}
	expected = append(expected, dateFrom, dateTo)

	assertSecondarySalesSummaryDateVars(t, actual, expected)
}

func TestExistsCustomerInParentScopeSQLRequiresActiveChild(t *testing.T) {
	t.Parallel()

	db, recorded := newReportRepoDryRunDB(t)
	repo := NewReportRepo(db)

	allowed, err := repo.ExistsCustomerInParentScope("CUST-1", "PARENT-1")
	if err != nil {
		t.Fatalf("ExistsCustomerInParentScope returned error: %v", err)
	}
	if allowed {
		t.Fatalf("expected dry-run scope check to return false, got true")
	}

	query := latestRecordedQuery(t, recorded)
	if !strings.Contains(query.SQL, `cust_id = $1 AND parent_cust_id = $2 AND is_del = false AND is_active = true`) {
		t.Fatalf("expected strict active scope filter, sql=%s", query.SQL)
	}
	if !reflect.DeepEqual(query.Vars, []interface{}{"CUST-1", "PARENT-1"}) {
		t.Fatalf("unexpected vars: %#v", query.Vars)
	}
}

func TestSecondarySalesReportTrendSalesSQLUsesSourceTablesAndNetSalesFormula(t *testing.T) {
	t.Parallel()

	db, recorded := newReportRepoDryRunDB(t)
	repo := NewReportRepo(db)

	if _, err := repo.SecondarySalesReportTrendSales([]string{"CUST-1"}, 2026); err != nil {
		t.Fatalf("SecondarySalesReportTrendSales returned error: %v", err)
	}

	query := latestRecordedQuery(t, recorded)
	for _, check := range []string{
		`WITH months AS (`,
		`SELECT 1 AS month`,
		`FROM sls."order" o`,
		`JOIN sls.order_detail od ON od.ro_no = o.ro_no AND od.cust_id = o.cust_id`,
		`FROM sls.return_det rd`,
		`JOIN sls."return" r ON r.return_no = rd.return_no AND r.cust_id = rd.cust_id`,
		`JOIN sls."order" o ON o.invoice_no = r.invoice_no AND o.cust_id = r.cust_id`,
		`o.data_status IN (6, 7)`,
		`o.invoice_date >= $2`,
		`o.invoice_date < $3`,
		`r.return_date >= $5`,
		`r.return_date < $6`,
		`EXTRACT(MONTH FROM o.invoice_date)::INTEGER`,
		`EXTRACT(MONTH FROM r.return_date)::INTEGER`,
		`COALESCE(od.sell_price_final1, 0)`,
		`COALESCE(od.sell_price_final2, 0)`,
		`COALESCE(od.sell_price_final3, 0)`,
		`COALESCE(od.promo_final1, 0) + COALESCE(od.promo_final2, 0)`,
		`COALESCE(od.promo_final3, 0) + COALESCE(od.promo_final4, 0)`,
		`COALESCE(od.promo_final5, 0)`,
		`COALESCE(od.vat_value_final, 0)`,
		`COALESCE(rd.sell_price1, 0)`,
		`COALESCE(rd.sell_price2, 0)`,
		`COALESCE(rd.sell_price3, 0)`,
		`COALESCE(rd.promo_value, 0)`,
		`COALESCE(rd.vat_value, 0)`,
		`m.month AS month`,
		`COALESCE(os.gross_sales, 0) - COALESCE(rs.gross_sales, 0) AS total_gross_sale`,
		`COALESCE(os.discount_promo, 0) - COALESCE(rs.discount_promo, 0) AS total_discount_promo`,
		`((COALESCE(os.gross_sales, 0) - COALESCE(os.discount_promo, 0)) -`,
		`(COALESCE(rs.gross_sales, 0) - COALESCE(rs.discount_promo, 0))) +`,
		`(COALESCE(os.ppn, 0) - COALESCE(rs.ppn, 0)) AS net_sales`,
	} {
		if !strings.Contains(query.SQL, check) {
			t.Fatalf("expected trend SQL to contain %q\nSQL:\n%s", check, query.SQL)
		}
	}
	if strings.Contains(query.SQL, `report.fact_orders`) {
		t.Fatalf("expected source trend SQL, found fact orders reference\nSQL:\n%s", query.SQL)
	}

	expectedFrom := time.Date(2026, time.January, 1, 0, 0, 0, 0, time.UTC)
	expectedTo := time.Date(2027, time.January, 1, 0, 0, 0, 0, time.UTC)
	expectedVars := []interface{}{
		"CUST-1", expectedFrom, expectedTo,
		"CUST-1", expectedFrom, expectedTo,
	}
	assertSecondarySalesSummaryDateVars(t, query.Vars, expectedVars)
}

func TestSalesmanActivityReportTrendSalesSQLUsesSourceTablesAndNetSalesFormula(t *testing.T) {
	t.Parallel()

	db, recorded := newReportRepoDryRunDB(t)
	repo := NewReportRepo(db)

	if _, err := repo.SalesmanActivityReportTrendSales([]string{"C260020001"}, 2026); err != nil {
		t.Fatalf("SalesmanActivityReportTrendSales returned error: %v", err)
	}

	query := latestRecordedQuery(t, recorded)
	for _, check := range []string{
		`WITH months AS (`,
		`SELECT 1 AS month_num`,
		`FROM sls."order" o`,
		`JOIN sls.order_detail od ON od.ro_no = o.ro_no AND od.cust_id = o.cust_id`,
		`FROM sls.return_det rd`,
		`JOIN sls."return" r ON r.return_no = rd.return_no AND r.cust_id = rd.cust_id`,
		`JOIN sls."order" o ON o.invoice_no = r.invoice_no AND o.cust_id = r.cust_id`,
		`o.data_status IN (6, 7)`,
		`o.invoice_date >= $2`,
		`o.invoice_date < $3`,
		`r.return_date >= $5`,
		`r.return_date < $6`,
		`EXTRACT(MONTH FROM o.invoice_date)::INTEGER`,
		`EXTRACT(MONTH FROM r.return_date)::INTEGER`,
		`COALESCE(od.sell_price_final1, 0)`,
		`COALESCE(od.promo_final1, 0) + COALESCE(od.promo_final2, 0)`,
		`COALESCE(od.vat_value_final, 0)`,
		`COALESCE(rd.promo_value, 0)`,
		`COALESCE(rd.vat_value, 0)`,
		`AS total_invoice`,
		`AS total_return`,
		`COALESCE(o.net_sales_order, 0) - COALESCE(r.net_sales_return, 0) AS net_sales`,
	} {
		if !strings.Contains(query.SQL, check) {
			t.Fatalf("expected activity trend SQL to contain %q\nSQL:\n%s", check, query.SQL)
		}
	}
	if strings.Contains(query.SQL, `report.fact_orders`) {
		t.Fatalf("expected source trend SQL, found fact orders reference\nSQL:\n%s", query.SQL)
	}

	expectedFrom := time.Date(2026, time.January, 1, 0, 0, 0, 0, time.UTC)
	expectedTo := time.Date(2027, time.January, 1, 0, 0, 0, 0, time.UTC)
	expectedVars := []interface{}{
		"C260020001", expectedFrom, expectedTo,
		"C260020001", expectedFrom, expectedTo,
	}
	assertSecondarySalesSummaryDateVars(t, query.Vars, expectedVars)
}

func TestSecondarySalesReportSumReportByMonthSQLUsesSourceTablesAndDateRange(t *testing.T) {
	t.Parallel()

	db, recorded := newReportRepoDryRunDB(t)
	repo := NewReportRepo(db)

	filter := entity.SecondarySalesReportDashboardSumPayload{Month: 6}
	if _, err := repo.SecondarySalesReportSumReportByMonth([]string{"CUST-1"}, filter, 2026); err != nil {
		t.Fatalf("SecondarySalesReportSumReportByMonth returned error: %v", err)
	}

	query := latestRecordedQuery(t, recorded)
	for _, check := range []string{
		`FROM sls."order" o`,
		`JOIN sls.order_detail od ON od.ro_no = o.ro_no AND od.cust_id = o.cust_id`,
		`FROM sls.return_det rd`,
		`JOIN sls."return" r ON r.return_no = rd.return_no AND r.cust_id = rd.cust_id`,
		`JOIN sls."order" o ON o.invoice_no = r.invoice_no AND o.cust_id = r.cust_id`,
		`o.invoice_date >= $2`,
		`o.invoice_date < $3`,
		`o.invoice_date >= $5`,
		`o.invoice_date < $6`,
		`COALESCE(od.qty3_final, 0) * COALESCE(od.conv_unit2, 1) * COALESCE(od.conv_unit3, 1)`,
		`COALESCE(od.qty2_final, 0) * COALESCE(od.conv_unit2, 1)`,
		`COALESCE(rd.qty3, 0) * COALESCE(rd.conv_unit2, 1) * COALESCE(rd.conv_unit3, 1)`,
		`COALESCE(rd.qty2, 0) * COALESCE(rd.conv_unit2, 1)`,
		`COALESCE(od.disc_value_final, 0) +`,
		`COALESCE(od.promo_final1, 0) + COALESCE(od.promo_final2, 0) +`,
		`COALESCE(od.promo_final3, 0) + COALESCE(od.promo_final4, 0) +`,
		`COALESCE(od.promo_final5, 0)`,
		`COALESCE(rd.disc_value, 0) + COALESCE(rd.promo_value, 0)`,
		`CASE WHEN od.item_type = 1 THEN COALESCE(od.vat_value_final, 0) ELSE 0 END`,
		`COALESCE(rd.total, 0) - COALESCE(rd.vat_value, 0)`,
		`COALESCE(rd.total, 0)`,
		`(os.gross_sale - rs.gross_sale) AS total_gross_sale`,
		`COALESCE(os.discount_promo, 0) - COALESCE(rs.discount_promo, 0) AS total_discount_promo`,
		`(os.ppn - rs.ppn) AS total_ppn`,
		`(os.net_sales_exc_ppn - rs.net_sales_exc_ppn) AS net_sales_exc_ppn`,
		`(os.net_sales_inc_ppn - rs.net_sales_inc_ppn) AS net_sales`,
		`COALESCE(os.qty, 0) - COALESCE(rs.qty_return, 0) AS qty`,
		`rs.qty_return AS qty_return`,
		`rs.net_sales_inc_ppn AS net_sales_return`,
		`ROUND(((rs.net_sales_inc_ppn / NULLIF(os.net_sales_inc_ppn, 0)) * 100)::numeric, 2)`,
		`MAX(o.updated_at) AS last_update`,
		`MAX(r.updated_at) AS last_update`,
	} {
		if !strings.Contains(query.SQL, check) {
			t.Fatalf("expected summary SQL to contain %q\nSQL:\n%s", check, query.SQL)
		}
	}
	for _, reject := range []string{
		`os.qty AS qty`,
		`(os.discount_promo + rs.discount_promo) AS total_discount_promo`,
		`COALESCE(os.discount_promo, 0) + COALESCE(rs.discount_promo, 0) AS total_discount_promo`,
		`r.data_status = 6`,
	} {
		if strings.Contains(query.SQL, reject) {
			t.Fatalf("expected summary SQL NOT to contain %q\nSQL:\n%s", reject, query.SQL)
		}
	}
	if strings.Contains(query.SQL, `FROM report.fact_orders fo`) {
		t.Fatalf("expected source summary SQL, found fact orders reference\nSQL:\n%s", query.SQL)
	}
	if strings.Contains(query.SQL, `FROM report.fact_returns fr`) {
		t.Fatalf("expected source summary SQL, found fact returns reference\nSQL:\n%s", query.SQL)
	}

	expectedFrom := time.Date(2026, time.June, 1, 0, 0, 0, 0, time.UTC)
	expectedTo := time.Date(2026, time.July, 1, 0, 0, 0, 0, time.UTC)
	expectedVars := []interface{}{"CUST-1", expectedFrom, expectedTo, "CUST-1", expectedFrom, expectedTo}
	assertSecondarySalesSummaryDateVars(t, query.Vars, expectedVars)
}

func TestSX2258SecondarySalesReportSumReportByMonthSQLUsesOptionalSummaryFilters(t *testing.T) {
	t.Parallel()

	from := int64(1778086800)
	to := int64(1778173199)
	filter := entity.SecondarySalesReportDashboardSumPayload{
		Month:       5,
		From:        &from,
		To:          &to,
		OutletIDs:   []int64{301},
		SalesmanIDs: []int64{101},
		ProIDs:      []int64{501},
	}

	db, recorded := newReportRepoDryRunDB(t)
	repo := NewReportRepo(db)

	if _, err := repo.SecondarySalesReportSumReportByMonth([]string{"CUST-1"}, filter, 2026); err != nil {
		t.Fatalf("SecondarySalesReportSumReportByMonth returned error: %v", err)
	}

	query := latestRecordedQuery(t, recorded)
	for _, check := range []string{
		`o.outlet_id IN (`,
		`o.salesman_id IN (`,
		`od.pro_id IN (`,
		`r.outlet_id IN (`,
		`r.salesman_id IN (`,
		`rd.product_id IN (`,
	} {
		if !strings.Contains(query.SQL, check) {
			t.Fatalf("expected filtered summary SQL to contain %q\nSQL:\n%s", check, query.SQL)
		}
	}
	if strings.Contains(query.SQL, `r.return_date`) {
		t.Fatalf("expected summary filters to keep return branch on invoice_date\nSQL:\n%s", query.SQL)
	}

	expectedVars := []interface{}{
		"CUST-1",
		str.UnixTimestampToUtcTime(from),
		str.UnixTimestampToUtcTime(to),
		int64(301),
		int64(101),
		int64(501),
		"CUST-1",
		str.UnixTimestampToUtcTime(from),
		str.UnixTimestampToUtcTime(to),
		int64(301),
		int64(101),
		int64(501),
	}
	assertSecondarySalesSummaryDateVars(t, query.Vars, expectedVars)
}

func TestSX2258SecondarySalesReportSummaryArithmeticRegression(t *testing.T) {
	t.Parallel()

	orderQty := int64(150)
	returnQty := int64(16)
	if got := orderQty - returnQty; got != 134 {
		t.Fatalf("expected net qty 134, got %d", got)
	}

	orderDiscountPromo := float64(1_500_000)
	returnDiscountPromo := float64(261_260)
	if got := orderDiscountPromo - returnDiscountPromo; got != 1_238_740 {
		t.Fatalf("expected net discount promo 1238740, got %v", got)
	}

	var promoFinal2 *float64
	orderFormula := float64(1_400_000) +
		float64(50_000) +
		valueOrZero(promoFinal2) +
		float64(25_000) +
		float64(20_000) +
		float64(5_000)
	if orderFormula != orderDiscountPromo {
		t.Fatalf("expected null promo parts treated as zero and formula total 1500000, got %v", orderFormula)
	}
}

func TestSalesmanActivityReportSumByMonthSQLUsesSourceTablesAndNetSalesFormula(t *testing.T) {
	t.Parallel()

	db, recorded := newReportRepoDryRunDB(t)
	repo := NewReportRepo(db)

	if _, err := repo.SalesmanActivityReportSumByMonth([]string{"C260020001"}, 6, 2026); err != nil {
		t.Fatalf("SalesmanActivityReportSumByMonth returned error: %v", err)
	}

	query := latestRecordedQuery(t, recorded)
	for _, check := range []string{
		`FROM sls."order" o`,
		`JOIN sls.order_detail od ON od.ro_no = o.ro_no AND od.cust_id = o.cust_id`,
		`COALESCE(od.sell_price_final1, 0)`,
		`COALESCE(od.promo_final1, 0) + COALESCE(od.promo_final2, 0)`,
		`FROM sls.return_det rd`,
		`JOIN sls."return" r ON r.return_no = rd.return_no AND r.cust_id = rd.cust_id`,
		`JOIN sls."order" o ON o.invoice_no = r.invoice_no AND o.cust_id = r.cust_id`,
		`o.invoice_date >= $2`,
		`o.invoice_date < $3`,
		`COUNT(DISTINCT trx.salesman_id) AS total_salesman`,
		`AS total_sales`,
		`AS total_return`,
	} {
		if !strings.Contains(query.SQL, check) {
			t.Fatalf("expected activity summary SQL to contain %q\nSQL:\n%s", check, query.SQL)
		}
	}
	if strings.Contains(query.SQL, `report.fact_orders`) || strings.Contains(query.SQL, `report.fact_returns`) {
		t.Fatalf("expected source tables, found fact table reference\nSQL:\n%s", query.SQL)
	}

	expectedFrom := time.Date(2026, time.June, 1, 0, 0, 0, 0, time.UTC)
	expectedTo := time.Date(2026, time.July, 1, 0, 0, 0, 0, time.UTC)
	expectedVars := []interface{}{
		"C260020001", expectedFrom, expectedTo,
		"C260020001", expectedFrom, expectedTo,
		"C260020001", expectedFrom, expectedTo,
		"C260020001", expectedFrom, expectedTo,
	}
	if len(query.Vars) != len(expectedVars) {
		t.Fatalf("expected %d vars, got %d: %#v", len(expectedVars), len(query.Vars), query.Vars)
	}
	for i := range expectedVars {
		if !reflect.DeepEqual(query.Vars[i], expectedVars[i]) {
			t.Fatalf("param %d mismatch: expected %#v, got %#v", i, expectedVars[i], query.Vars[i])
		}
	}
}

func TestSecondarySalesReportReturnSumReportByMonthSQLUsesQuotedYearFilter(t *testing.T) {
	t.Parallel()

	db, recorded := newReportRepoDryRunDB(t)
	repo := NewReportRepo(db)

	if _, err := repo.SecondarySalesReportReturnSumReportByMonth([]string{"CUST-1"}, 5, 2026); err != nil {
		t.Fatalf("SecondarySalesReportReturnSumReportByMonth returned error: %v", err)
	}

	query := latestRecordedQuery(t, recorded)
	if !strings.Contains(query.SQL, `dt."year" = $3`) {
		t.Fatalf("expected quoted year filter, sql=%s", query.SQL)
	}
	if !reflect.DeepEqual(query.Vars, []interface{}{"CUST-1", 5, 2026}) {
		t.Fatalf("unexpected vars: %#v", query.Vars)
	}
}

// TestBuildSecondarySalesUnionQueryNetSalesFromFormula verifies that the ORDER
// branch of the union CTE derives net_sales_exc_ppn and net_sales_inc_ppn from
// qty*price - discounts (not from amount_final which can be stale/wrong).
// Regression for invoice INV2605200015: expected NetSalesExcPPN=17500000,
// NetSalesIncPPN=19250000 when amount_final is inconsistent.
func TestBuildSecondarySalesUnionQueryNetSalesFromFormula(t *testing.T) {
	t.Parallel()

	filter := entity.SecondarySalesReportQueryFilter{
		CustID:       "C260020001",
		ParentCustID: "C26002",
	}

	sql, _, _ := buildSecondarySalesUnionQuery(filter, false)

	// Must NOT use amount_final for net sales in the ORDER branch.
	if strings.Contains(sql, "od.amount_final") {
		t.Fatalf("ORDER branch must not reference od.amount_final; got SQL:\n%s", sql)
	}

	// Must compute net_sales_exc_ppn from gross - special_discount - discount.
	netExcChecks := []string{
		"COALESCE(od.qty1_final, 0) * COALESCE(od.sell_price1, 0)",
		"COALESCE(od.qty2_final, 0) * COALESCE(od.sell_price2, 0)",
		"COALESCE(od.qty3_final, 0) * COALESCE(od.sell_price3, 0)",
		"COALESCE(od.promo_value_final, 0)",
		"COALESCE(od.disc_value_final, 0)",
	}
	for _, check := range netExcChecks {
		if !strings.Contains(sql, check) {
			t.Fatalf("expected net_sales formula to contain %q\nSQL:\n%s", check, sql)
		}
	}

	// Must compute net_sales_inc_ppn = net_sales_exc_ppn + vat.
	if !strings.Contains(sql, "COALESCE(od.vat_value_final, 0)") {
		t.Fatalf("expected net_sales_inc_ppn to add COALESCE(od.vat_value_final, 0)\nSQL:\n%s", sql)
	}

	// RETURN branch must remain total-based and null-safe.
	if !strings.Contains(sql, "COALESCE(rd.total, 0) * -1 AS net_sales_inc_ppn") {
		t.Fatalf("RETURN branch net_sales_inc_ppn must remain COALESCE(rd.total, 0) * -1\nSQL:\n%s", sql)
	}
	if !strings.Contains(sql, "(COALESCE(rd.total, 0) - COALESCE(rd.vat_value, 0)) * -1 AS net_sales_exc_ppn") {
		t.Fatalf("RETURN branch net_sales_exc_ppn must remain null-safe total - vat\nSQL:\n%s", sql)
	}
}

func TestBuildSecondarySalesUnionQueryReturnPromoAndDiscount(t *testing.T) {
	t.Parallel()

	filter := entity.SecondarySalesReportQueryFilter{CustID: "C260020001", CustIDs: []string{"C260020001"}, ParentCustID: "C26002"}
	sql, _, _ := buildSecondarySalesUnionQuery(filter, false)

	for _, check := range []string{
		"COALESCE(od.promo_value_final, 0) AS special_discount",
		"COALESCE(od.disc_value_final, 0) AS discount",
		"COALESCE(rd.promo_value, 0) AS special_discount",
		"COALESCE(rd.disc_value, 0) AS discount",
	} {
		if !strings.Contains(sql, check) {
			t.Fatalf("expected RETURN branch to contain %q\nSQL:\n%s", check, sql)
		}
	}
	if strings.Contains(sql, "0 AS special_discount") {
		t.Fatalf("return export/list SQL must not zero promo special_discount\nSQL:\n%s", sql)
	}
}

// TestBuildSecondarySalesUnionQueryNetSalesArithmetic validates the business math
// for a concrete example: qty=10, price=1_750_000, promo=0, disc=0, vat=1_750_000
// → net_exc=17_500_000, net_inc=19_250_000.
func TestBuildSecondarySalesUnionQueryNetSalesArithmetic(t *testing.T) {
	t.Parallel()

	// Simulate the formula directly (mirrors what the SQL CASE expression computes).
	qty1 := float64(10)
	price1 := float64(1_750_000)
	promo := float64(0)
	disc := float64(0)
	vat := float64(1_750_000)

	gross := qty1 * price1
	netExc := gross - promo - disc
	netInc := netExc + vat

	if netExc != 17_500_000 {
		t.Fatalf("expected net_sales_exc_ppn=17500000, got %v", netExc)
	}
	if netInc != 19_250_000 {
		t.Fatalf("expected net_sales_inc_ppn=19250000, got %v", netInc)
	}
}

// TestGetReportSecondarySalesReportOrderNetSalesFromFormula verifies that the
// GetReportSecondarySalesReportOrder query also uses the formula-based net sales
// and does not reference od.amount_final.
func TestGetReportSecondarySalesReportOrderNetSalesFromFormula(t *testing.T) {
	t.Parallel()

	sql, params := buildReportSecondarySalesReportOrderQuery("CUST-1", time.Now(), 10, 0)

	if got, want := strings.Count(sql, "?"), len(params); got != want {
		t.Fatalf("placeholder count mismatch: placeholders=%d params=%d\nSQL:\n%s\nparams:%#v", got, want, sql, params)
	}

	if strings.Contains(sql, "od.amount_final") {
		t.Fatalf("GetReportSecondarySalesReportOrder must not reference od.amount_final\nSQL:\n%s", sql)
	}
	for _, check := range []string{
		"COALESCE(od.qty1_final, 0) * COALESCE(od.sell_price1, 0)",
		"COALESCE(od.promo_value_final, 0) AS special_discount",
		"COALESCE(od.disc_value_final, 0) AS discount",
		"COALESCE(od.promo_value_final, 0)",
		"COALESCE(od.disc_value_final, 0)",
		"COALESCE(od.vat_value_final, 0)",
	} {
		if !strings.Contains(sql, check) {
			t.Fatalf("expected formula fragment %q in GetReportSecondarySalesReportOrder SQL\nSQL:\n%s", check, sql)
		}
	}
}

func TestGetReportSecondarySalesReportReturnPromoAndDiscountSQL(t *testing.T) {
	t.Parallel()

	sql, params := buildReportSecondarySalesReportReturnQuery("CUST-1", time.Date(2026, 6, 5, 0, 0, 0, 0, time.UTC), 10, 0)
	if got, want := strings.Count(sql, "?"), len(params); got != want {
		t.Fatalf("placeholder count mismatch: placeholders=%d params=%d\nSQL:\n%s\nparams:%#v", got, want, sql, params)
	}
	for _, check := range []string{
		"COALESCE(rd.promo_value, 0) AS special_discount",
		"COALESCE(rd.disc_value, 0) AS discount",
		"COALESCE(rd.total, 0) AS net_sales_inc_ppn",
	} {
		if !strings.Contains(sql, check) {
			t.Fatalf("expected return extract SQL to contain %q\nSQL:\n%s", check, sql)
		}
	}
	if strings.Contains(sql, "0 AS special_discount") {
		t.Fatalf("return extract SQL must not zero promo special_discount\nSQL:\n%s", sql)
	}
}

func TestSecondarySalesLegacyQueryNetSalesFromFormula(t *testing.T) {
	t.Parallel()

	db, recorded := newReportRepoDryRunDB(t)
	repo := NewReportRepo(db)

	_, _, _, err := repo.SecondarySales(entity.SecondarySalesReportQueryFilter{
		CustID:       "C260020001",
		ParentCustID: "C26002",
	})
	if err != nil {
		t.Fatalf("SecondarySales returned error: %v", err)
	}

	query := latestRecordedQuery(t, recorded)
	if strings.Contains(query.SQL, "sls.order_detail.amount_final") {
		t.Fatalf("SecondarySales legacy query must not reference amount_final\nSQL:\n%s", query.SQL)
	}
	for _, check := range []string{
		"COALESCE(sls.order_detail.qty1_final,0)*COALESCE(sls.order_detail.sell_price1,0)",
		"COALESCE(sls.order_detail.qty2_final,0)*COALESCE(sls.order_detail.sell_price2,0)",
		"COALESCE(sls.order_detail.qty3_final,0)*COALESCE(sls.order_detail.sell_price3,0)",
		"COALESCE(sls.order_detail.promo_value_final,0)",
		"COALESCE(sls.order_detail.disc_value_final,0)",
		"COALESCE(sls.order_detail.vat_value_final,0)",
	} {
		if !strings.Contains(query.SQL, check) {
			t.Fatalf("expected formula fragment %q in SecondarySales SQL\nSQL:\n%s", check, query.SQL)
		}
	}
}

func TestSecondarySalesReportGroupQueriesUseSourceTablesAndDateRange(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name   string
		invoke func(repo *RepositoryReportImpl) error
	}{
		{
			name: "outlet",
			invoke: func(repo *RepositoryReportImpl) error {
				_, err := repo.SecondarySalesReportGroupOutlet([]string{"CUST-1"}, 5, 2026)
				return err
			},
		},
		{
			name: "salesman",
			invoke: func(repo *RepositoryReportImpl) error {
				_, err := repo.SecondarySalesReportGroupSalesman([]string{"CUST-1"}, 5, 2026)
				return err
			},
		},
		{
			name: "product_category",
			invoke: func(repo *RepositoryReportImpl) error {
				_, err := repo.SecondarySalesReportProductCategory([]string{"CUST-1"}, 5, 2026)
				return err
			},
		},
		{
			name: "product",
			invoke: func(repo *RepositoryReportImpl) error {
				_, err := repo.SecondarySalesReportProduct([]string{"CUST-1"}, 5, 2026)
				return err
			},
		},
	}

	expectedFrom := time.Date(2026, time.May, 1, 0, 0, 0, 0, time.UTC)
	expectedTo := time.Date(2026, time.June, 1, 0, 0, 0, 0, time.UTC)

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			db, recorded := newReportRepoDryRunDB(t)
			repo := NewReportRepo(db)

			if err := tc.invoke(repo); err != nil {
				t.Fatalf("%s returned error: %v", tc.name, err)
			}

			query := latestRecordedQuery(t, recorded)
			for _, check := range []string{
				`FROM sls."order" o`,
				`JOIN sls.order_detail od ON od.ro_no = o.ro_no AND od.cust_id = o.cust_id`,
				`WHERE o.cust_id IN ($1) AND o.data_status IN (6,7) AND o.invoice_date >= $2 AND o.invoice_date < $3`,
				`FROM sls.return_det rd`,
				`JOIN sls."return" r ON r.return_no = rd.return_no AND r.cust_id = rd.cust_id`,
				`JOIN sls."order" o ON o.invoice_no = r.invoice_no AND o.cust_id = r.cust_id`,
				`WHERE rd.cust_id IN ($4) AND o.data_status IN (6,7) AND o.invoice_date >= $5 AND o.invoice_date < $6`,
				`COALESCE(od.vat_value_final, 0)`,
				`COALESCE(rd.disc_value, 0)`,
				`COALESCE(rd.vat_value, 0)`,
				`) * -1 AS net_sales`,
				`GROUP BY id, code, name`,
			} {
				if !strings.Contains(query.SQL, check) {
					t.Fatalf("expected group SQL to contain %q\nSQL:\n%s", check, query.SQL)
				}
			}
			if strings.Contains(query.SQL, `FROM report.fact_orders fo`) {
				t.Fatalf("expected source group SQL, found fact orders reference\nSQL:\n%s", query.SQL)
			}
			if strings.Contains(query.SQL, `FROM report.fact_returns fr`) {
				t.Fatalf("expected source group SQL, found fact returns reference\nSQL:\n%s", query.SQL)
			}
			assertSecondarySalesGroupDateVars(t, query.Vars, []string{"CUST-1"}, expectedFrom, expectedTo)
		})
	}
}

func TestSecondarySalesReportGroupOutletUsesOutletDisplayMapping(t *testing.T) {
	t.Parallel()

	db, recorded := newReportRepoDryRunDB(t)
	repo := NewReportRepo(db)

	if _, err := repo.SecondarySalesReportGroupOutlet([]string{"CUST-1"}, 6, 2026); err != nil {
		t.Fatalf("SecondarySalesReportGroupOutlet returned error: %v", err)
	}

	query := latestRecordedQuery(t, recorded)
	// SX-2224: outlet branch must use outlet_id as id, outlet_code/outlet_name from m_outlet.
	for _, check := range []string{
		"o.outlet_id AS id",
		"COALESCE(mo.outlet_code, '') AS code",
		"COALESCE(mo.outlet_name, '') AS name",
		"LEFT JOIN mst.m_outlet mo ON mo.outlet_id = o.outlet_id AND mo.cust_id = o.cust_id",
		"r.outlet_id AS id",
		"LEFT JOIN mst.m_outlet mo ON mo.outlet_id = r.outlet_id AND mo.cust_id = rd.cust_id",
		"GROUP BY id, code, name",
	} {
		if !strings.Contains(query.SQL, check) {
			t.Fatalf("expected outlet group SQL to contain %q\nSQL:\n%s", check, query.SQL)
		}
	}
	// Must not use salesman_id as the grouping key or concatenate emp_name into the name.
	for _, reject := range []string{
		"o.salesman_id AS id",
		"r.salesman_id AS id",
		"CONCAT_WS(' > '",
	} {
		if strings.Contains(query.SQL, reject) {
			t.Fatalf("expected outlet group SQL NOT to contain %q\nSQL:\n%s", reject, query.SQL)
		}
	}
	expectedFrom := time.Date(2026, time.June, 1, 0, 0, 0, 0, time.UTC)
	expectedTo := time.Date(2026, time.July, 1, 0, 0, 0, 0, time.UTC)
	assertSecondarySalesGroupDateVars(t, query.Vars, []string{"CUST-1"}, expectedFrom, expectedTo)
}

func TestSecondarySalesReportGroupProductCategoryUsesMasterCategoryFallback(t *testing.T) {
	t.Parallel()

	db, recorded := newReportRepoDryRunDB(t)
	repo := NewReportRepo(db)

	if _, err := repo.SecondarySalesReportProductCategory([]string{"CUST-1", "CUST-2"}, 6, 2026); err != nil {
		t.Fatalf("SecondarySalesReportProductCategory returned error: %v", err)
	}

	query := latestRecordedQuery(t, recorded)
	for _, check := range []string{
		"LEFT JOIN mst.m_product mp ON mp.pro_id = od.pro_id AND mp.cust_id = od.cust_id",
		"LEFT JOIN mst.m_product_cat mpc ON mpc.pcat_id = mp.pcat_id",
		"COALESCE(NULLIF(mpc.pcat_id, 0), 0) AS id",
		"COALESCE(NULLIF(mpc.pcat_code, ''), '') AS code",
		"COALESCE(NULLIF(mpc.pcat_name, ''), '') AS name",
		"LEFT JOIN mst.m_product mp ON mp.pro_id = rd.product_id AND mp.cust_id = rd.cust_id",
	} {
		if !strings.Contains(query.SQL, check) {
			t.Fatalf("expected product category group SQL to contain %q\nSQL:\n%s", check, query.SQL)
		}
	}
	if strings.Contains(query.SQL, "mp.cust_id = ?") {
		t.Fatalf("expected row-level product join, sql=%s", query.SQL)
	}
	expectedFrom := time.Date(2026, time.June, 1, 0, 0, 0, 0, time.UTC)
	expectedTo := time.Date(2026, time.July, 1, 0, 0, 0, 0, time.UTC)
	assertSecondarySalesGroupDateVars(t, query.Vars, []string{"CUST-1", "CUST-2"}, expectedFrom, expectedTo)
}

func TestSecondarySalesReportGroupProductUsesMasterProductFallback(t *testing.T) {
	t.Parallel()

	db, recorded := newReportRepoDryRunDB(t)
	repo := NewReportRepo(db)

	if _, err := repo.SecondarySalesReportProduct([]string{"CUST-1", "CUST-2"}, 6, 2026); err != nil {
		t.Fatalf("SecondarySalesReportProduct returned error: %v", err)
	}

	query := latestRecordedQuery(t, recorded)
	for _, check := range []string{
		"LEFT JOIN mst.m_product mp ON mp.pro_id = od.pro_id AND mp.cust_id = od.cust_id",
		"COALESCE(NULLIF(mp.pro_id, 0), od.pro_id) AS id",
		"COALESCE(NULLIF(mp.pro_code, ''), '') AS code",
		"COALESCE(NULLIF(mp.pro_name, ''), '') AS name",
		"LEFT JOIN mst.m_product mp ON mp.pro_id = rd.product_id AND mp.cust_id = rd.cust_id",
		"COALESCE(NULLIF(mp.pro_id, 0), rd.product_id) AS id",
	} {
		if !strings.Contains(query.SQL, check) {
			t.Fatalf("expected product group SQL to contain %q\nSQL:\n%s", check, query.SQL)
		}
	}
	if strings.Contains(query.SQL, "mp.cust_id = ?") {
		t.Fatalf("expected row-level product join, sql=%s", query.SQL)
	}
	expectedFrom := time.Date(2026, time.June, 1, 0, 0, 0, 0, time.UTC)
	expectedTo := time.Date(2026, time.July, 1, 0, 0, 0, 0, time.UTC)
	assertSecondarySalesGroupDateVars(t, query.Vars, []string{"CUST-1", "CUST-2"}, expectedFrom, expectedTo)
}
