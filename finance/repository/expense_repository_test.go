package repository

import (
	"context"
	"finance/entity"
	"strings"
	"testing"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func newExpenseTypeDryRunDB(t *testing.T) *gorm.DB {
	t.Helper()

	dsn := "host=localhost user=postgres password=postgres dbname=postgres port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{DryRun: true})
	if err != nil {
		t.Fatalf("failed to init dry-run db: %v", err)
	}

	return db
}

func TestExpenseRepository_BuildExpenseTypeListBaseQuery_UsesParentCustomerScopeAndSearch(t *testing.T) {
	db := newExpenseTypeDryRunDB(t)
	repo := &expenseRepositoryImpl{DB: db}
	filter := entity.ExpenseQueryFilter{
		Q:            "Uang bensin",
		Page:         1,
		Limit:        10,
		CustId:       "C220010001",
		ParentCustId: "C22001",
	}

	tx := repo.buildExpenseTypeListBaseQuery(context.Background(), filter).Find(&[]map[string]interface{}{})
	sql := tx.Statement.SQL.String()

	if !strings.Contains(sql, "acf.expense_type.cust_id =") {
		t.Fatalf("expected SQL to include expense type tenant scope, got %s", sql)
	}
	if !strings.Contains(sql, "acf.expense_type.is_del = false") {
		t.Fatalf("expected SQL to include soft-delete filter, got %s", sql)
	}
	if !strings.Contains(sql, "acf.expense_type.expense_type_code LIKE") || !strings.Contains(sql, "acf.expense_type.expense_type_name LIKE") {
		t.Fatalf("expected SQL to keep search scoped by code/name, got %s", sql)
	}

	vars := tx.Statement.Vars
	if len(vars) != 3 {
		t.Fatalf("expected parent customer and search vars, got %#v", vars)
	}
	if vars[0] != "C22001" {
		t.Fatalf("expected parent customer scope C22001, got %#v", vars[0])
	}
	if vars[1] != "%Uang bensin%" || vars[2] != "%Uang bensin%" {
		t.Fatalf("expected search pattern vars, got %#v", vars)
	}
}
