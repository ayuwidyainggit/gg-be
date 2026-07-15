package repository

import (
	"strings"
	"testing"
)

func TestActivitySalesReportPaymentDataUsesDepositDate(t *testing.T) {
	p := activityReportSQLParams{
		CustID:       "C260020001",
		ParentCustID: "C26002",
		FromDate:     "2026-05-01",
		ToDate:       "2026-05-31",
		SalesmanIDs:  []int{471},
	}

	distSQL, _ := buildActivitySalesReportDistributorSQL(activityReportPJPConfig{
		Schema:             "pjp",
		FilterOutletByCust: true,
	}, p)
	assertPaymentDataUsesDeposit(t, distSQL)

	principalSQL, _ := buildActivitySalesReportPrincipalSQL(p)
	assertPaymentDataUsesDeposit(t, principalSQL)
}

func assertPaymentDataUsesDeposit(t *testing.T, fullSQL string) {
	t.Helper()
	start := strings.Index(fullSQL, "payment_data AS (")
	if start < 0 {
		t.Fatal("payment_data CTE not found")
	}
	rest := fullSQL[start:]
	end := strings.Index(rest, "\n)\nSELECT")
	if end < 0 {
		t.Fatal("payment_data CTE end not found")
	}
	paymentCTE := rest[:end]
	if !strings.Contains(paymentCTE, "d.deposit_date::date") {
		t.Fatal("payment_data must filter by acf.deposit.deposit_date")
	}
	if strings.Contains(paymentCTE, "o.ro_date::date") {
		t.Fatal("payment_data must not filter by sls.order.ro_date")
	}
	if strings.Contains(paymentCTE, "acf.collection") {
		t.Fatal("payment_data must not use acf.collection")
	}
	if !strings.Contains(paymentCTE, "acf.deposit_detail") {
		t.Fatal("payment_data must join acf.deposit_detail")
	}
}

func TestActivityReportCustIDPredicateSingleAndMulti(t *testing.T) {
	singleSQL, singleArgs := activityReportCustIDPredicate("p.cust_id", []string{"C26002"})
	if singleSQL != "p.cust_id = ?" || len(singleArgs) != 1 || singleArgs[0] != "C26002" {
		t.Fatalf("unexpected single predicate: sql=%q args=%v", singleSQL, singleArgs)
	}

	multiSQL, multiArgs := activityReportCustIDPredicate("p.cust_id", []string{"C26002", "C26003"})
	if multiSQL != "p.cust_id IN (?,?)" || len(multiArgs) != 2 {
		t.Fatalf("unexpected multi predicate: sql=%q args=%v", multiSQL, multiArgs)
	}
}

func TestBuildActivitySalesReportDistributorSQLUsesEmployeeNameForSalesmanName(t *testing.T) {
	p := activityReportSQLParams{
		CustID:       "C260020001",
		ParentCustID: "C26002",
		FromDate:     "2026-06-01",
		ToDate:       "2026-06-15",
	}
	sql, _ := buildActivitySalesReportDistributorSQL(activityReportPJPConfig{
		Schema:             "pjp",
		FilterOutletByCust: true,
	}, p)
	for _, check := range []string{
		`SELECT emp_id, emp_code, emp_name`,
		`COALESCE(NULLIF(TRIM(e.emp_name), ''), p.salesman_name) AS salesman_name`,
		`ORDER BY salesman_name ASC, visit_date ASC`,
	} {
		if !strings.Contains(sql, check) {
			t.Fatalf("expected distributor SQL to contain %q\nSQL:\n%s", check, sql)
		}
	}
}

func TestBuildActivitySalesReportPrincipalSQLUsesEmployeeNameForSalesmanName(t *testing.T) {
	p := activityReportSQLParams{
		CustIDs:      []string{"C26002"},
		ParentCustID: "C26002",
		FromDate:     "2026-06-01",
		ToDate:       "2026-06-15",
	}
	sql, _ := buildActivitySalesReportPrincipalSQL(p)
	for _, check := range []string{
		`SELECT emp_id, emp_code, emp_name`,
		`COALESCE(NULLIF(TRIM(e.emp_name), ''), p.salesman_name) AS salesman_name`,
		`ORDER BY salesman_name ASC, visit_date ASC`,
	} {
		if !strings.Contains(sql, check) {
			t.Fatalf("expected principal SQL to contain %q\nSQL:\n%s", check, sql)
		}
	}
}

func TestBuildActivitySalesReportDistributorSQLUsesDistributorCodeAsBusinessUnit(t *testing.T) {
	p := activityReportSQLParams{
		CustID:       "C260020001",
		ParentCustID: "C26002",
		FromDate:     "2026-06-01",
		ToDate:       "2026-06-15",
	}
	sql, _ := buildActivitySalesReportDistributorSQL(activityReportPJPConfig{
		Schema:             "pjp",
		FilterOutletByCust: true,
	}, p)
	for _, check := range []string{
		`JOIN mst.m_distributor md ON md.cust_id = mc.cust_id`,
		`c.distributor_code AS business_unit_code`,
		`c.cust_name AS business_unit_name`,
		`'' AS distributor_code`,
		`'' AS distributor_name`,
	} {
		if !strings.Contains(sql, check) {
			t.Fatalf("expected distributor SQL to contain %q\nSQL:\n%s", check, sql)
		}
	}
	if strings.Contains(sql, `c.cust_id AS business_unit_code`) {
		t.Fatal("business_unit_code must not use cust_id")
	}
	if !strings.Contains(sql, `WHERE mc.cust_id = ?`) {
		t.Fatal("cust_data must filter by mc.cust_id")
	}
}

func TestBuildActivitySalesReportPrincipalSQLUsesDistributorCodeAsBusinessUnit(t *testing.T) {
	p := activityReportSQLParams{
		CustIDs:      []string{"C26002", "C26003"},
		ParentCustID: "C26002",
		FromDate:     "2026-06-01",
		ToDate:       "2026-06-01",
	}
	sql, _ := buildActivitySalesReportPrincipalSQL(p)
	for _, check := range []string{
		`LEFT JOIN mst.m_distributor md ON md.cust_id = mc.cust_id`,
		`'' AS business_unit_code`,
		`c.cust_name AS business_unit_name`,
		`outlet_dist_data`,
		`COALESCE(NULLIF(TRIM(d.distributor_code), ''), NULLIF(TRIM(o_dist.distributor_code), '')) AS distributor_code`,
		`COALESCE(NULLIF(TRIM(d.distributor_name), ''), NULLIF(TRIM(o_dist.distributor_name), '')) AS distributor_name`,
	} {
		if !strings.Contains(sql, check) {
			t.Fatalf("expected principal SQL to contain %q\nSQL:\n%s", check, sql)
		}
	}
	if strings.Contains(sql, `c.cust_id AS business_unit_code`) {
		t.Fatal("business_unit_code must not use cust_id directly")
	}
	if strings.Contains(sql, `CASE WHEN c.distributor_code <> '' THEN c.distributor_code ELSE c.cust_id END AS business_unit_code`) {
		t.Fatal("principal business_unit_code must be empty")
	}
	if !strings.Contains(sql, `WHERE mc.cust_id IN (?,?)`) {
		t.Fatal("cust_data must filter by mc.cust_id")
	}
}

func TestBuildActivitySalesReportPrincipalSQLSupportsMultiCustID(t *testing.T) {
	p := activityReportSQLParams{
		CustIDs:      []string{"C26002", "C26003"},
		ParentCustID: "C26002",
		FromDate:     "2026-06-01",
		ToDate:       "2026-06-01",
	}
	sql, _ := buildActivitySalesReportPrincipalSQL(p)
	if !strings.Contains(sql, "p.cust_id IN (?,?)") {
		t.Fatalf("expected principal pjp filter to use IN clause, got fragment around pjp_data")
	}
	if strings.Count(sql, "cust_id IN (?,?)") < 3 {
		t.Fatalf("expected multiple cust_id IN clauses in principal SQL")
	}
}

func TestSplitActivityReportCustIDs(t *testing.T) {
	principal, distributor := splitActivityReportCustIDs("C26002", []string{
		"C26002", "C260020001", "C260020002",
	})
	if len(principal) != 1 || principal[0] != "C26002" {
		t.Fatalf("expected principal IDs [C26002], got %v", principal)
	}
	if len(distributor) != 2 {
		t.Fatalf("expected 2 distributor IDs, got %v", distributor)
	}
}

func TestBuildActivitySalesReportCombinedSQLUsesUnionForMixedCustIDs(t *testing.T) {
	principalIDs, distributorIDs := splitActivityReportCustIDs("C26002", []string{"C26002", "C260020001"})
	if len(principalIDs) != 1 || principalIDs[0] != "C26002" {
		t.Fatalf("unexpected principal IDs: %v", principalIDs)
	}
	if len(distributorIDs) != 1 || distributorIDs[0] != "C260020001" {
		t.Fatalf("unexpected distributor IDs: %v", distributorIDs)
	}

	base := activityReportSQLParams{
		ParentCustID: "C26002",
		FromDate:     "2026-06-17",
		ToDate:       "2026-06-18",
		SalesmanIDs:  []int{491},
	}

	principalParams := base
	principalParams.CustIDs = principalIDs
	principalSQL, _ := buildActivitySalesReportPrincipalQuery(principalParams)

	distributorParams := base
	distributorParams.CustIDs = distributorIDs
	distributorParams.SkipSalesmanFilter = true
	distributorSQL, _ := buildActivitySalesReportDistributorQuery(
		activityReportPJPConfig{Schema: "pjp", FilterOutletByCust: true, IsPrincipal: false},
		distributorParams,
	)

	combined := "(" + principalSQL + ") UNION ALL (" + distributorSQL + ")"
	if strings.Contains(combined, "ORDER BY salesman_name ASC UNION") {
		t.Fatal("union parts must not append row ordering before UNION ALL")
	}
	if strings.Contains(distributorSQL, "p.salesman_id IN") {
		t.Fatal("distributor branch in mixed multi-BU query must not filter by salesman_id")
	}
	if !strings.Contains(combined, "pjp_principles.outlet_visit_list") {
		t.Fatal("expected principal pjp_principles query in union")
	}
	if !strings.Contains(combined, "pjp.outlet_visit_list") {
		t.Fatal("expected distributor pjp query in union")
	}
}

func TestAppendActivityReportSalesmanFilterSkippedWhenFlagSet(t *testing.T) {
	p := activityReportSQLParams{
		CustIDs:            []string{"C260020001"},
		SalesmanIDs:        []int{491},
		SkipSalesmanFilter: true,
	}
	sql, args := appendActivityReportSalesmanFilter(p)
	if sql != "" || len(args) != 0 {
		t.Fatalf("expected empty salesman filter, got sql=%q args=%v", sql, args)
	}
}

func TestBuildActivitySalesReportSalesmanFilterUsesEmpCodeAcrossBusinessUnits(t *testing.T) {
	p := activityReportSQLParams{
		CustIDs:                  []string{"C260020001"},
		SalesmanReferenceCustIDs: []string{"C26002", "C260020001"},
		SalesmanIDs:              []int{491},
	}
	sql, _ := appendActivityReportSalesmanFilterForBranch(p, []string{"C260020001"})
	for _, check := range []string{
		`e.emp_code IN (`,
		`e2.emp_id IN (?)`,
		`e2.cust_id IN (?,?)`,
		`e.cust_id IN (?)`,
	} {
		if !strings.Contains(sql, check) {
			t.Fatalf("expected salesman filter to contain %q\nSQL:%s", check, sql)
		}
	}
}

func TestFinalizeActivityReportSQLFiltersDistributorCode(t *testing.T) {
	base := `SELECT '' AS business_unit_code, 'DIST-A' AS distributor_code`
	p := activityReportSQLParams{
		DistributorCodes: []string{"DIST-A", "DIST-B"},
		Limit:            10,
	}

	sql, args := finalizeActivityReportSQL(base, p, nil)
	for _, check := range []string{
		`activity_report_src`,
		activityReportEffectiveDistributorCodeExpr,
		`IN (?,?)`,
		`ORDER BY salesman_name ASC`,
		`LIMIT ?`,
	} {
		if !strings.Contains(sql, check) {
			t.Fatalf("expected SQL to contain %q\nSQL:%s", check, sql)
		}
	}
	if len(args) != 3 || args[0] != "DIST-A" || args[1] != "DIST-B" || args[2] != 10 {
		t.Fatalf("unexpected args: %#v", args)
	}

	countSQL, countArgs := finalizeActivityReportSQL(base, activityReportSQLParams{
		DistributorCodes: []string{"DIST-A"},
		ForCount:         true,
	}, nil)
	if !strings.Contains(countSQL, `SELECT COUNT(*) FROM`) || !strings.Contains(countSQL, activityReportEffectiveDistributorCodeExpr) {
		t.Fatalf("unexpected count SQL: %s", countSQL)
	}
	if len(countArgs) != 1 || countArgs[0] != "DIST-A" {
		t.Fatalf("unexpected count args: %#v", countArgs)
	}
}

func TestIsActivityReportPrincipalCust(t *testing.T) {
	tests := []struct {
		name         string
		custID       string
		parentCustID string
		want         bool
	}{
		{"principal C26004", "C26004", "C26004", true},
		{"distributor child C260040005", "C260040005", "C26004", false},
		{"principal C26002", "C26002", "C26002", true},
		{"child C260020001", "C260020001", "C26002", false},
		{"fallback short id without parent", "C26004", "", true},
		{"fallback long id without parent", "C260040005", "", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isActivityReportPrincipalCust(tt.custID, tt.parentCustID); got != tt.want {
				t.Fatalf("isActivityReportPrincipalCust(%q, %q) = %v, want %v", tt.custID, tt.parentCustID, got, tt.want)
			}
		})
	}
}

func TestActivityReportClockTimeSelectDoesNotDoubleApplyWIBOffset(t *testing.T) {
	t.Parallel()

	if strings.Contains(activityReportClockTimeWIBSelect, `+ INTERVAL '7 hour'`) {
		t.Fatal("clock in/out must not add +7h; attendances.created_at is already in DB session timezone (WIB)")
	}
	if !strings.Contains(activityReportClockTimeWIBSelect, `TO_CHAR(att.clock_in_time, 'YYYY-MM-DD HH24:MI')`) {
		t.Fatalf("expected direct TO_CHAR on attendances timestamp:\n%s", activityReportClockTimeWIBSelect)
	}
	if !strings.Contains(activityReportVisitTimeWIBSelect, `+ INTERVAL '7 hour'`) {
		t.Fatal("visit check-in/out must still convert unix epoch to WIB")
	}
}

func TestActivityReportRemarksSelectUsesApprovedLeaveRequest(t *testing.T) {
	t.Parallel()

	for _, check := range []string{
		`mobile.leave_request`,
		`lr.emp_id = p.salesman_id`,
		`lr.cust_id = p.cust_id`,
		`LOWER(TRIM(lr.approval)) IN ('approve', 'approved')`,
		`p.date::date BETWEEN lr.start_date AND lr.end_date`,
		`'On Leave'`,
		`ELSE '-'`,
	} {
		if !strings.Contains(activityReportRemarksSelect, check) {
			t.Fatalf("expected remarks select to contain %q\nSQL:\n%s", check, activityReportRemarksSelect)
		}
	}

	p := activityReportSQLParams{
		CustID:   "C260020001",
		FromDate: "2026-07-01",
		ToDate:   "2026-07-14",
	}
	principalSQL, _ := buildActivitySalesReportPrincipalSQL(p)
	distributorSQL, _ := buildActivitySalesReportDistributorSQL(activityReportPJPConfig{
		Schema:             "pjp",
		FilterOutletByCust: true,
	}, p)
	for _, sql := range []string{principalSQL, distributorSQL} {
		if !strings.Contains(sql, "mobile.leave_request") {
			t.Fatalf("expected leave_request remarks in activity report SQL:\n%s", sql)
		}
		if !strings.Contains(sql, "AS remarks") {
			t.Fatalf("expected remarks column alias in activity report SQL:\n%s", sql)
		}
	}
}
