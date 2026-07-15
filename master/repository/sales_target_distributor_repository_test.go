package repository

import (
	"strings"
	"testing"
	"time"

	"master/entity"
)

func TestBuildSalesTargetDistributorStatusClause_MultiValueStatus(t *testing.T) {
	currentYear := time.Now().Year()
	statuses := []int{int(entity.SALES_TARGET_STATUS_ACTIVE), int(entity.SALES_TARGET_STATUS_INACTIVE)}

	clause := buildSalesTargetDistributorStatusClause(statuses, currentYear)

	if clause == "" {
		t.Fatalf("expected non-empty clause")
	}

	if !strings.Contains(clause, "std.year <=") || !strings.Contains(clause, "std.is_active = true") {
		t.Fatalf("expected active condition in clause, got: %s", clause)
	}

	if !strings.Contains(clause, "std.year >") || !strings.Contains(clause, "std.is_active = false") {
		t.Fatalf("expected inactive condition in clause, got: %s", clause)
	}

	if !strings.Contains(clause, " OR ") {
		t.Fatalf("expected OR combination for multi-value status, got: %s", clause)
	}
}

func TestBuildSalesTargetDistributorStatusClause_EmptyStatus(t *testing.T) {
	clause := buildSalesTargetDistributorStatusClause([]int{}, time.Now().Year())
	if clause != "" {
		t.Fatalf("expected empty clause, got: %s", clause)
	}
}
