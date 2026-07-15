package repository

import (
	"strings"
	"testing"
)

func TestAppendIntInFilter(t *testing.T) {
	baseQuery := " WHERE o.is_del = false "
	updatedQuery := appendIntInFilter(baseQuery, "o.ot_class_id", []int{123, 456})

	if !strings.Contains(updatedQuery, "o.ot_class_id IN (123,456)") {
		t.Fatalf("expected ot_class_id filter in query, got %s", updatedQuery)
	}
}

func TestAppendIntInFilter_EmptyValues(t *testing.T) {
	baseQuery := " WHERE o.is_del = false "
	updatedQuery := appendIntInFilter(baseQuery, "o.ot_class_id", nil)

	if updatedQuery != baseQuery {
		t.Fatalf("expected query to stay unchanged, got %s", updatedQuery)
	}
}

func TestAppendIntInFilter_VerificationStatusApproved(t *testing.T) {
	baseQuery := " WHERE o.is_del = false "
	updatedQuery := appendIntInFilter(baseQuery, "o.verification_status", []int{1})

	if !strings.Contains(updatedQuery, "o.verification_status IN (1)") {
		t.Fatalf("expected verification_status filter in query, got %s", updatedQuery)
	}
}

func TestAppendIntInFilter_OutletStatusValidSurveyValues(t *testing.T) {
	baseQuery := " WHERE o.is_del = false "
	updatedQuery := appendIntInFilter(baseQuery, "o.outlet_status", []int{1, 5, 6, 7})

	if !strings.Contains(updatedQuery, "o.outlet_status IN (1,5,6,7)") {
		t.Fatalf("expected outlet_status multi-value filter in query, got %s", updatedQuery)
	}
}

func TestBuildOutletCustScopeWhere_WithResolvedCustIDs(t *testing.T) {
	query := buildOutletCustScopeWhere("C26002", "C26002", []string{"C260020001", "C260020002"})

	if !strings.Contains(query, "o.cust_id IN ('C260020001','C260020002')") {
		t.Fatalf("expected resolved cust ids filter, got %s", query)
	}
}

func TestBuildOutletCustScopeWhere_PrincipalScopeIncludesChildren(t *testing.T) {
	query := buildOutletCustScopeWhere("C26002", "C26002", nil)

	if !strings.Contains(query, "o.cust_id = 'C26002'") {
		t.Fatalf("expected principal cust id scope, got %s", query)
	}
}

func TestBuildOutletCustScopeWhere_DistributorScopeUsesSingleCustID(t *testing.T) {
	query := buildOutletCustScopeWhere("C260020001", "C26002", nil)
	expected := " WHERE o.is_del = false and o.cust_id = 'C260020001' "

	if query != expected {
		t.Fatalf("expected %s, got %s", expected, query)
	}
}

func TestShouldApplyOutletIsActiveFilter_DefaultAndPrecedence(t *testing.T) {
	if !shouldApplyOutletIsActiveFilter(nil) {
		t.Fatalf("expected nil include_inactive to keep applying is_active filter")
	}

	includeInactiveOne := 1
	if shouldApplyOutletIsActiveFilter(&includeInactiveOne) {
		t.Fatalf("expected include_inactive=1 to bypass is_active filter")
	}

	includeInactiveZero := 0
	if !shouldApplyOutletIsActiveFilter(&includeInactiveZero) {
		t.Fatalf("expected include_inactive=0 to keep applying is_active filter")
	}
}
