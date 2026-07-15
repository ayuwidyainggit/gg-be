package repository

import (
	"finance/entity"
	"strings"
	"testing"
)

func TestBuildPaymentDepositQuery_ARUsesDepositAndOptionalEmpID(t *testing.T) {
	repo := &RepositoryPaymentDepositReportImpl{}
	filter := entity.PaymentDepositReportQueryFilter{
		DepositType: []string{"AR"},
		EmpID:       []string{"421", "415", "381"},
		DepositNo:   []string{"DP1", "DP2"},
		StartDate:   "2026-04-24",
		EndDate:     "2026-04-27",
	}

	sql, _ := repo.buildQuery(filter, "C260020001")

	assertContains(t, sql, "FROM acf.deposit d")
	assertContains(t, sql, "d.cust_id = ?")
	assertContains(t, sql, "d.deleted_at IS NULL")
	assertContains(t, sql, "d.deposit_date BETWEEN ? AND ?")
	assertContains(t, sql, "d.emp_id IN ?")
	assertContains(t, sql, "d.deposit_no IN ?")
	assertNotContains(t, sql, "account_payable_payment")
}

func TestBuildPaymentDepositQuery_ARWithoutEmpIDDoesNotFilterCollector(t *testing.T) {
	repo := &RepositoryPaymentDepositReportImpl{}
	filter := entity.PaymentDepositReportQueryFilter{
		DepositType: []string{"AR"},
		StartDate:   "2026-04-24",
		EndDate:     "2026-04-27",
	}

	sql, _ := repo.buildQuery(filter, "C260020001")

	assertContains(t, sql, "FROM acf.deposit d")
	assertNotContains(t, sql, "d.emp_id IN ?")
}

func TestBuildPaymentDepositQuery_APUsesAccountPayableAndNoEmpID(t *testing.T) {
	repo := &RepositoryPaymentDepositReportImpl{}
	filter := entity.PaymentDepositReportQueryFilter{
		DepositType: []string{"AP"},
		EmpID:       []string{"421"},
		DepositNo:   []string{"PY1"},
		StartDate:   "2026-04-24",
		EndDate:     "2026-04-27",
	}

	sql, _ := repo.buildQuery(filter, "C260020001")

	assertContains(t, sql, "FROM acf.account_payable_payment app")
	assertContains(t, sql, "app.cust_id = ?")
	assertContains(t, sql, "app.deleted_by IS NULL")
	assertContains(t, sql, "app.account_payable_payment_date BETWEEN ? AND ?")
	assertContains(t, sql, "app.account_payable_payment_no IN ?")
	assertNotContains(t, sql, "emp_id IN ?")
}

func TestBuildPaymentDepositQuery_ARAndAPUsesUnionAll(t *testing.T) {
	repo := &RepositoryPaymentDepositReportImpl{}
	filter := entity.PaymentDepositReportQueryFilter{
		DepositType: []string{"AR", "AP"},
		StartDate:   "2026-04-24",
		EndDate:     "2026-04-27",
	}

	sql, _ := repo.buildQuery(filter, "C260020001")

	assertContains(t, sql, "FROM acf.deposit d")
	assertContains(t, sql, "FROM acf.account_payable_payment app")
	assertContains(t, sql, "UNION ALL")
}

func TestBuildDownloadARQuery_ExpenseIsNegative(t *testing.T) {
	repo := &RepositoryPaymentDepositReportImpl{}
	sql, _ := repo.buildDownloadARQuery("C1", "P1", "2026-04-24", "2026-04-27", nil, nil)
	assertContains(t, sql, "-ABS(COALESCE(SUM(de.payment_amount), 0)) AS expense")
}

func TestBuildSafeSortAcceptsCreatedDateAlias(t *testing.T) {
	repo := &RepositoryPaymentDepositReportImpl{}

	if got, want := repo.buildSafeSort("created_date:desc"), "t.deposit_date DESC"; got != want {
		t.Fatalf("buildSafeSort() = %q, want %q", got, want)
	}
}

func TestPaymentDepositReport_BuildDownloadARQuery_ExpenseTypeJoinByPKOnly(t *testing.T) {
	repo := &RepositoryPaymentDepositReportImpl{}
	sql, args := repo.buildDownloadARQuery("C001", "P001", "2026-06-01", "2026-06-30", nil, nil)

	// Join must be PK-only — no etr.cust_id scope
	assertNotContains(t, sql, "etr.cust_id")
	assertContains(t, sql, "etr.expense_type_id = ex.expense_type_id")

	// SELECT expense_name must include both code and name for "code - name" format
	assertContains(t, sql, "expense_type_code")
	assertContains(t, sql, "expense_type_name")

	// GROUP BY in the expense branch must include etr.expense_type_code
	groupByIdx := strings.LastIndex(sql, "GROUP BY")
	if groupByIdx == -1 {
		t.Fatal("expected GROUP BY clause in AR expense query")
	}
	assertContains(t, sql[groupByIdx:], "etr.expense_type_code")

	// Semantic check (position-independent): parentCustId must no longer be bound
	// by the expense_type join. parentCustId is still used once for the m_outlet
	// join in the first AR subquery, so it must appear exactly once in args.
	const custId, parentCustId = "C001", "P001"
	parentCount := 0
	for _, a := range args {
		if a == parentCustId {
			parentCount++
		}
	}
	if parentCount != 1 {
		t.Fatalf("expected parentCustId %q to appear exactly once in args (m_outlet join only, not expense_type join), got %d: %#v", parentCustId, parentCount, args)
	}
	if !containsArg(args, custId) {
		t.Fatalf("expected custId %q to be present in args, got %#v", custId, args)
	}
}

func containsArg(args []interface{}, want interface{}) bool {
	for _, a := range args {
		if a == want {
			return true
		}
	}
	return false
}

func assertContains(t *testing.T, haystack, needle string) {
	t.Helper()
	if !strings.Contains(haystack, needle) {
		t.Fatalf("expected SQL to contain %q, got %s", needle, haystack)
	}
}

func assertNotContains(t *testing.T, haystack, needle string) {
	t.Helper()
	if strings.Contains(haystack, needle) {
		t.Fatalf("expected SQL not to contain %q, got %s", needle, haystack)
	}
}
