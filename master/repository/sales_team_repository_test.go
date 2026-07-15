package repository

import (
	"strings"
	"testing"
)

func TestBuildSalesTeamCustScopeCondition_WithoutDistributorFilterUsesCustID(t *testing.T) {
	condition := buildSalesTeamCustScopeCondition(nil, "C22001", "C220010001")

	if condition != " a.cust_id = 'C220010001' " {
		t.Fatalf("expected direct cust_id scope when distributor filter empty, got: %s", condition)
	}
}

func TestBuildSalesTeamCustScopeCondition_WithDistributorFilter(t *testing.T) {
	condition := buildSalesTeamCustScopeCondition([]int{67, 68}, "C22001", "C220010001")

	if !strings.Contains(condition, "FROM smc.m_customer mc") {
		t.Fatalf("expected distributor condition to map through smc.m_customer, got: %s", condition)
	}

	if !strings.Contains(condition, "mc.parent_cust_id = 'C22001'") {
		t.Fatalf("expected parent_cust_id scope in distributor condition, got: %s", condition)
	}

	if !strings.Contains(condition, "mc.distributor_id IN (67,68)") {
		t.Fatalf("expected distributor IN clause, got: %s", condition)
	}
}

func TestBuildSalesTeamCustScopeCondition_WithPrincipalScopeOnly(t *testing.T) {
	condition := buildSalesTeamCustScopeCondition([]int{0}, "C22001", "C220010001")

	if !strings.Contains(condition, "a.cust_id = 'C22001'") {
		t.Fatalf("expected principal cust_id clause, got: %s", condition)
	}

	if strings.Contains(condition, "mc.distributor_id IN") {
		t.Fatalf("expected principal-only scope without distributor IN clause, got: %s", condition)
	}
}

func TestBuildSalesTeamCustScopeCondition_WithPrincipalAndDistributorScope(t *testing.T) {
	condition := buildSalesTeamCustScopeCondition([]int{0, 67, 67}, "C22001", "C220010001")

	if !strings.Contains(condition, "a.cust_id = 'C22001'") {
		t.Fatalf("expected principal cust_id clause, got: %s", condition)
	}

	if !strings.Contains(condition, "mc.distributor_id IN (67)") {
		t.Fatalf("expected distributor IN clause, got: %s", condition)
	}

	if !strings.Contains(condition, " OR ") {
		t.Fatalf("expected combined principal and distributor scope, got: %s", condition)
	}
}
