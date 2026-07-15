package repository

import (
	"strings"
	"testing"

	"master/entity"
)

func TestBuildSalesmanCustScopeCondition_WithDistributorFilter(t *testing.T) {
	condition := buildSalesmanCustScopeCondition(nil, []int{10, 20}, "C22001", "C220010001")

	if !strings.Contains(condition, "FROM smc.m_customer mc") {
		t.Fatalf("expected distributor condition to map through smc.m_customer, got: %s", condition)
	}

	if !strings.Contains(condition, "mc.parent_cust_id = 'C22001'") {
		t.Fatalf("expected parent_cust_id scope in distributor condition, got: %s", condition)
	}

	if !strings.Contains(condition, "mc.distributor_id IN (10,20)") {
		t.Fatalf("expected distributor IN clause, got: %s", condition)
	}
}

func TestBuildSalesmanCustScopeCondition_WithPrincipalAndDistributorScope(t *testing.T) {
	condition := buildSalesmanCustScopeCondition(nil, []int{0, 67, 67, 68}, "C22001", "C220010001")

	if !strings.Contains(condition, "s.cust_id = 'C22001'") {
		t.Fatalf("expected principal cust_id clause to be included, got: %s", condition)
	}

	if !strings.Contains(condition, "mc.distributor_id IN (67,68)") {
		t.Fatalf("expected distributor IN clause, got: %s", condition)
	}

	if !strings.Contains(condition, " OR ") {
		t.Fatalf("expected grouped OR union between principal and distributor scope, got: %s", condition)
	}

	if !strings.HasPrefix(condition, "( ") || !strings.HasSuffix(condition, " )") {
		t.Fatalf("expected grouped condition with parentheses, got: %s", condition)
	}
}

func TestBuildSalesmanCustScopeCondition_FallbackToCustIDWhenParentEmpty(t *testing.T) {
	condition := buildSalesmanCustScopeCondition(nil, []int{10}, "", "C22001")

	if !strings.Contains(condition, "mc.parent_cust_id = 'C22001'") {
		t.Fatalf("expected cust_id fallback when parent_cust_id empty, got: %s", condition)
	}
}

func TestBuildSalesmanCustScopeCondition_PrincipalOnlyUsesParentCustID(t *testing.T) {
	condition := buildSalesmanCustScopeCondition(nil, []int{0}, "C22001", "C220010001")

	if condition != "( s.cust_id = 'C22001' )" {
		t.Fatalf("expected principal-only grouped scope, got: %s", condition)
	}

	if strings.Contains(condition, "mc.distributor_id IN") {
		t.Fatalf("expected principal-only scope to exclude distributor branch, got: %s", condition)
	}
}

func TestBuildSalesmanCustScopeCondition_DistributorOnlyUsesMappedScope(t *testing.T) {
	condition := buildSalesmanCustScopeCondition(nil, []int{10, 20}, "C22001", "C220010001")

	if strings.Contains(condition, "s.cust_id = 'C22001'") {
		t.Fatalf("expected distributor-only scope to exclude principal branch, got: %s", condition)
	}

	if !strings.Contains(condition, "mc.distributor_id IN (10,20)") {
		t.Fatalf("expected distributor-only scope to include mapped distributor ids, got: %s", condition)
	}

	if !strings.HasPrefix(condition, "( ") || !strings.HasSuffix(condition, " )") {
		t.Fatalf("expected distributor-only grouped condition, got: %s", condition)
	}
}

func TestBuildSalesmanCustScopeCondition_DeduplicatesAndIgnoresNegativeIDs(t *testing.T) {
	condition := buildSalesmanCustScopeCondition(nil, []int{0, 120, -1, 120}, "C22001", "C220010001")

	if !strings.Contains(condition, "s.cust_id = 'C22001'") {
		t.Fatalf("expected principal scope to be included, got: %s", condition)
	}

	if !strings.Contains(condition, "mc.distributor_id IN (120)") {
		t.Fatalf("expected deduplicated positive distributor ids only, got: %s", condition)
	}

	if strings.Contains(condition, "-1") {
		t.Fatalf("expected negative distributor ids to be ignored, got: %s", condition)
	}
}

func TestBuildSalesmanCustScopeCondition_WithoutDistributorFilterUsesCustID(t *testing.T) {
	condition := buildSalesmanCustScopeCondition(nil, nil, "C22001", "C220010001")

	if condition != " s.cust_id = 'C220010001' " {
		t.Fatalf("expected direct cust_id scope when distributor filter empty, got: %s", condition)
	}
}

func TestBuildSalesmanCustScopeCondition_WithCustIDFilter(t *testing.T) {
	condition := buildSalesmanCustScopeCondition([]string{"C260020002"}, nil, "C26002", "C26002")

	if condition != ` s.cust_id = 'C260020002' ` {
		t.Fatalf("expected single cust_id equality, got: %s", condition)
	}
}

func TestBuildSalesmanCustScopeCondition_WithMultipleCustIDs(t *testing.T) {
	condition := buildSalesmanCustScopeCondition([]string{"C26002", "C260020002"}, nil, "C26002", "C26002")

	if condition != ` s.cust_id IN ('C26002','C260020002') ` {
		t.Fatalf("expected cust_id IN clause, got: %s", condition)
	}
}

func TestBuildSalesmanCustScopeCondition_CustIDTakesPrecedenceOverDistributor(t *testing.T) {
	condition := buildSalesmanCustScopeCondition([]string{"C260020002"}, []int{10, 20}, "C26002", "C26002")

	if strings.Contains(condition, "distributor_id") {
		t.Fatalf("cust_id filter should take precedence over distributor_id, got: %s", condition)
	}

	if condition != ` s.cust_id = 'C260020002' ` {
		t.Fatalf("expected cust_id filter to be applied directly, got: %s", condition)
	}
}

func TestFindAllByCustId_AppliesDistributorFilterToQuery(t *testing.T) {
	filter := entity.SalesmanQueryFilter{
		Page:          1,
		Limit:         10,
		Sort:          "sales_name:asc",
		SalesTeamId:   "24,1",
		DistributorID: []int{10, 20},
	}

	condition := buildSalesmanCustScopeCondition(filter.CustIds, filter.DistributorID, "C22001", "C220010001")

	if !strings.Contains(condition, "mc.distributor_id IN (10,20)") {
		t.Fatalf("expected distributor condition to include selected distributors, got: %s", condition)
	}

	if strings.Contains(condition, "mc.distributor_id IN ()") {
		t.Fatalf("expected distributor condition without empty IN clause, got: %s", condition)
	}
}

func TestBuildSalesmanEmployeeJoin_UsesSalesmanCustID(t *testing.T) {
	join := buildSalesmanEmployeeJoin()

	if !strings.Contains(join, "emp.cust_id = s.cust_id") {
		t.Fatalf("expected employee join to follow salesman tenant, got: %s", join)
	}
}

func TestBuildSalesmanSalesTeamJoin_UsesSalesmanCustID(t *testing.T) {
	join := buildSalesmanSalesTeamJoin()

	if !strings.Contains(join, "st.cust_id = s.cust_id") {
		t.Fatalf("expected sales team join to follow salesman tenant, got: %s", join)
	}
}

func TestBuildSalesmanWarehouseJoin_UsesSalesmanCustID(t *testing.T) {
	join := buildSalesmanWarehouseJoin("wh", "s.wh_id")

	if !strings.Contains(join, "wh.wh_id = s.wh_id") {
		t.Fatalf("expected warehouse join to use selected source field, got: %s", join)
	}

	if !strings.Contains(join, "wh.cust_id = s.cust_id") {
		t.Fatalf("expected warehouse join to follow salesman tenant, got: %s", join)
	}
}

func TestBuildSalesmanCanvasJoin_UsesSalesmanCustID(t *testing.T) {
	join := buildSalesmanCanvasJoin()

	if !strings.Contains(join, "msc.cust_id = s.cust_id") {
		t.Fatalf("expected canvas join to follow salesman tenant, got: %s", join)
	}
}

func TestBuildSalesmanVehicleJoin_UsesSalesmanCustID(t *testing.T) {
	join := buildSalesmanVehicleJoin()

	if !strings.Contains(join, "mv.cust_id = s.cust_id") {
		t.Fatalf("expected vehicle join to follow salesman tenant, got: %s", join)
	}
}

func TestBuildSalesmanDriverJoin_UsesSalesmanCustID(t *testing.T) {
	join := buildSalesmanDriverJoin()

	if !strings.Contains(join, "me.cust_id = s.cust_id") {
		t.Fatalf("expected driver join to follow salesman tenant, got: %s", join)
	}
}
