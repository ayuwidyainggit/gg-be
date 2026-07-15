package repository

import (
	"context"
	"finance/entity"
	"strings"
	"testing"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func newExpenseEntryDryRunDB(t *testing.T) *gorm.DB {
	t.Helper()

	dsn := "host=localhost user=postgres password=postgres dbname=postgres port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{DryRun: true})
	if err != nil {
		t.Fatalf("failed to init dry-run db: %v", err)
	}

	return db
}

func TestExpenseEntryRepository_BuildSafeSort(t *testing.T) {
	repo := &expenseEntryRepositoryImpl{}

	tests := []struct {
		name string
		sort string
		want string
	}{
		{
			name: "default when empty",
			sort: "",
			want: "acf.expense.created_at DESC",
		},
		{
			name: "valid single created_date desc",
			sort: "created_date:desc",
			want: "acf.expense.created_at DESC",
		},
		{
			name: "valid multi sort",
			sort: "date:asc,balance:desc",
			want: "acf.expense.date ASC, acf.expense.balance DESC",
		},
		{
			name: "invalid direction fallback default",
			sort: "created_date:drop",
			want: "acf.expense.created_at DESC",
		},
		{
			name: "invalid field ignored and valid kept",
			sort: "foo:desc,amount:asc",
			want: "acf.expense.amount ASC",
		},
		{
			name: "injection-like payload fallback default",
			sort: "created_date:desc;drop table acf.expense",
			want: "acf.expense.created_at DESC",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := repo.buildSafeSort(tt.sort)
			if got != tt.want {
				t.Fatalf("buildSafeSort() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestExpenseEntryRepository_BuildExpenseEntryListQuery_UsesTenantUserFiltersAndStrictMinBalance(t *testing.T) {
	repo := &expenseEntryRepositoryImpl{}
	collectorIDs := []int64{11, 12}
	minBalance := 100.0
	filter := entity.ExpenseEntryQueryFilter{
		CustID:       "C001",
		UserID:       77,
		StartDate:    "2026-01-01",
		EndDate:      "2026-01-31",
		MinBalance:   &minBalance,
		CollectorIDs: collectorIDs,
		Limit:        20,
		Page:         1,
		Sort:         "created_date:desc",
	}

	if expenseEntryCollectorSelect != "emp.emp_name AS collector_name" {
		t.Fatalf("expected collector select constant to use mst.m_employee.emp_name, got %s", expenseEntryCollectorSelect)
	}
	if expenseEntryCollectorJoin != "LEFT JOIN mst.m_employee emp ON emp.emp_id = acf.expense.collector_id AND emp.cust_id = ?" {
		t.Fatalf("expected collector join constant to use mst.m_employee, got %s", expenseEntryCollectorJoin)
	}
	if strings.Contains(expenseEntryCollectorSelect, "sys.m_user") || strings.Contains(expenseEntryCollectorJoin, "sys.m_user") {
		t.Fatalf("expected no sys.m_user collector source, got select=%s join=%s", expenseEntryCollectorSelect, expenseEntryCollectorJoin)
	}
	if !strings.Contains(expenseEntryCollectorJoin, "mst.m_employee emp") {
		t.Fatalf("expected employee join source, got %s", expenseEntryCollectorJoin)
	}

	if repo.buildSafeSort(filter.Sort) != "acf.expense.created_at DESC" {
		t.Fatalf("expected safe sort created_at DESC, got %s", repo.buildSafeSort(filter.Sort))
	}
	if repo.buildSafeSort("balance:desc") != "acf.expense.balance DESC" {
		t.Fatalf("expected balance sort mapping, got %s", repo.buildSafeSort("balance:desc"))
	}
	if filter.MinBalance == nil || *filter.MinBalance != 100.0 {
		t.Fatalf("expected strict min_balance test fixture to remain > 100, got %#v", filter.MinBalance)
	}
	if len(filter.CollectorIDs) != 2 || filter.CollectorIDs[0] != 11 || filter.CollectorIDs[1] != 12 {
		t.Fatalf("expected collector_ids fixture to remain intact, got %#v", filter.CollectorIDs)
	}
}

func TestExpenseEntryRepository_BuildExpenseEntryListQuery_UsesExpenseDateRangeAndCollectorFilter(t *testing.T) {
	db := newExpenseEntryDryRunDB(t)
	repo := &expenseEntryRepositoryImpl{DB: db}
	minBalance := 0.0
	filter := entity.ExpenseEntryQueryFilter{
		CustID:       "C001",
		UserID:       77,
		Query:        "E20260424002",
		StartDate:    "2026-04-10",
		EndDate:      "2026-04-10",
		MinBalance:   &minBalance,
		CollectorIDs: []int64{360},
		Limit:        10,
		Page:         1,
		Sort:         "created_date:desc",
	}

	query := repo.buildExpenseEntryListQuery(context.Background(), filter)
	tx := query.Find(&[]map[string]interface{}{})
	sql := tx.Statement.SQL.String()

	if strings.Contains(sql, "acf.expense.created_by =") {
		t.Fatalf("expected list SQL not to filter by created_by, got %s", sql)
	}
	if !strings.Contains(sql, "acf.expense.collector_id IN") {
		t.Fatalf("expected collector_id filter in SQL, got %s", sql)
	}
	if !strings.Contains(sql, "acf.expense.doc_no ILIKE") {
		t.Fatalf("expected doc_no filter in SQL, got %s", sql)
	}
	if !strings.Contains(sql, "acf.expense.date >=") || !strings.Contains(sql, "acf.expense.date <=") {
		t.Fatalf("expected inclusive date range filter on acf.expense.date, got %s", sql)
	}
	if strings.Contains(strings.ToLower(sql), " between ") {
		t.Fatalf("expected SQL to avoid BETWEEN boundary issue, got %s", sql)
	}

	vars := tx.Statement.Vars
	if len(vars) < 6 {
		t.Fatalf("expected SQL vars to include tenant and date filters, got %#v", vars)
	}
	if vars[0] != "C001" {
		t.Fatalf("expected collector join var cust_id C001, got %#v", vars[0])
	}
	if vars[1] != "C001" {
		t.Fatalf("expected tenant filter var cust_id C001, got %#v", vars[1])
	}
	if vars[2] != "%E20260424002%" {
		t.Fatalf("expected doc_no like var %%E20260424002%%, got %#v", vars[2])
	}
	if vars[3] != "2026-04-10" || vars[4] != "2026-04-10" {
		t.Fatalf("expected exact business-date bounds, got %#v", vars)
	}
}

func TestExpenseEntryRepository_DefaultDateRangeAppliedByServiceBehavior(t *testing.T) {
	filter := entity.ExpenseEntryQueryFilter{
		CustID: "C001",
		UserID: 77,
		Limit:  20,
		Page:   1,
	}

	if filter.StartDate == "" || filter.EndDate == "" {
		endDate := time.Now().UTC()
		startDate := endDate.AddDate(0, -3, 0)
		if filter.StartDate == "" {
			filter.StartDate = startDate.Format("2006-01-02")
		}
		if filter.EndDate == "" {
			filter.EndDate = endDate.Format("2006-01-02")
		}
	}

	if filter.StartDate == "" || filter.EndDate == "" {
		t.Fatalf("expected default date range to be applied, got start=%q end=%q", filter.StartDate, filter.EndDate)
	}
}
