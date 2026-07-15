package repository

import (
	"strings"
	"testing"
	"time"
)

func TestBuildActivitySalesmanGroupSalesSQLUsesTransactionalNetSalesIncVAT(t *testing.T) {
	t.Parallel()

	sql := buildActivitySalesmanGroupSalesSQL()
	for _, check := range []string{
		`COALESCE(od.sell_price_final1, 0)`,
		`COALESCE(od.promo_final1, 0)`,
		`COALESCE(rd.promo_value, 0)`,
		`-1 AS multiplier`,
		`UNION ALL`,
		`SUM((`,
		`) * t.multiplier) AS net_sales`,
		`me.emp_code`,
		`ms.sales_name`,
		`o.invoice_date >= ?`,
		`o.data_status IN (6, 7)`,
	} {
		if !strings.Contains(sql, check) {
			t.Fatalf("expected sales group SQL to contain %q\nSQL:\n%s", check, sql)
		}
	}
	if strings.Contains(sql, `report.fact_orders`) {
		t.Fatal("expected source tables, found fact_orders")
	}
}

func TestBuildActivitySalesmanGroupReturnSQLUsesTransactionalNetSalesIncVAT(t *testing.T) {
	t.Parallel()

	sql := buildActivitySalesmanGroupReturnSQL()
	for _, check := range []string{
		`FROM sls.return_det rd`,
		`COALESCE(rd.sell_price1, 0)`,
		`SUM(`,
		`) AS net_sales`,
		`me.emp_code`,
		`ms.sales_name`,
	} {
		if !strings.Contains(sql, check) {
			t.Fatalf("expected return group SQL to contain %q\nSQL:\n%s", check, sql)
		}
	}
	if strings.Contains(sql, `report.fact_returns`) || strings.Contains(sql, `UNION ALL`) {
		t.Fatal("return group should only use return source data")
	}
}

func TestActivitySalesmanReportGroupSalesmanSQLUsesMonthDateRange(t *testing.T) {
	t.Parallel()

	db, recorded := newReportRepoDryRunDB(t)
	repo := NewReportRepo(db)

	if _, err := repo.ActivitySalesmanReportGroupSalesman([]string{"C260020001"}, 6, 2026); err != nil {
		t.Fatalf("ActivitySalesmanReportGroupSalesman returned error: %v", err)
	}

	query := latestRecordedQuery(t, recorded)
	expectedFrom := time.Date(2026, time.June, 1, 0, 0, 0, 0, time.UTC)
	expectedTo := time.Date(2026, time.July, 1, 0, 0, 0, 0, time.UTC)
	if !strings.Contains(query.SQL, `order_data AS`) {
		t.Fatalf("expected order_data CTE in SQL:\n%s", query.SQL)
	}
	if !strings.Contains(query.SQL, `cust_id IN`) {
		t.Fatalf("expected cust_id IN clause in SQL:\n%s", query.SQL)
	}
	if len(query.Vars) != 6 {
		t.Fatalf("expected 6 vars, got %d: %#v", len(query.Vars), query.Vars)
	}
	for _, idx := range []int{0, 3} {
		switch v := query.Vars[idx].(type) {
		case string:
			if v != "C260020001" {
				t.Fatalf("param %d mismatch: expected C260020001, got %#v", idx, v)
			}
		case []string:
			if len(v) != 1 || v[0] != "C260020001" {
				t.Fatalf("param %d mismatch: expected [C260020001], got %#v", idx, v)
			}
		default:
			t.Fatalf("param %d has unexpected type %T: %#v", idx, v, v)
		}
	}
	for _, idx := range []int{1, 2, 4, 5} {
		got, ok := query.Vars[idx].(time.Time)
		if !ok {
			t.Fatalf("param %d expected time.Time, got %#v", idx, query.Vars[idx])
		}
		want := expectedFrom
		if idx == 2 || idx == 5 {
			want = expectedTo
		}
		if !got.Equal(want) {
			t.Fatalf("param %d mismatch: expected %v, got %v", idx, want, got)
		}
	}
}
