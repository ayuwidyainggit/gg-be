package repository

import (
	"strings"
	"testing"

	"master/entity"
)

func TestBuildFindDistributorsByCustIDQuery_UsesInForMultiValueRegionAndArea(t *testing.T) {
	filter := entity.BusinessUnitQueryFilter{
		CustId:       "C22001",
		ParentCustId: "C22001",
		EmployeeId:   77,
		Scope:        entity.EmployeeDropdownScope{DistributorScope: "all", RegionScope: "all", AreaScope: "all"},
		RegionId:     []int{1, 2, 3},
		AreaId:       []int{10, 20},
		IsActive:     []int{1},
		Query:        "dist",
		Sort:         "area_id:asc",
		Page:         1,
		Limit:        10,
	}

	countQuery, countArgs, selectQuery, selectArgs := buildFindDistributorsByCustIDQuery(filter)

	if !strings.Contains(countQuery, "md.parent_cust_id = ?") {
		t.Fatalf("expected parent_cust_id filter in principal mapping query, got query: %s", countQuery)
	}
	if !strings.Contains(countQuery, "COUNT(DISTINCT md.distributor_id)") {
		t.Fatalf("expected distinct distributor count, got query: %s", countQuery)
	}
	if !strings.Contains(countQuery, "md.region_id IN") {
		t.Fatalf("expected region filter to use IN clause, got query: %s", countQuery)
	}
	if !strings.Contains(countQuery, "md.area_id IN") {
		t.Fatalf("expected area filter to use IN clause, got query: %s", countQuery)
	}
	if !strings.Contains(countQuery, "md.distributor_code ILIKE") || !strings.Contains(countQuery, "md.distributor_name ILIKE") {
		t.Fatalf("expected query search clause for q filter, got query: %s", countQuery)
	}
	if !strings.Contains(selectQuery, "SELECT DISTINCT") {
		t.Fatalf("expected distinct select query, got query: %s", selectQuery)
	}
	if !strings.Contains(selectQuery, "ORDER BY area_id asc") {
		t.Fatalf("expected sort clause to be preserved, got query: %s", selectQuery)
	}
	if !strings.Contains(selectQuery, "LIMIT 10 OFFSET 0") {
		t.Fatalf("expected pagination clause LIMIT/OFFSET, got query: %s", selectQuery)
	}
	if len(countArgs) != len(selectArgs) {
		t.Fatalf("expected count and select args length to match, got count=%d select=%d", len(countArgs), len(selectArgs))
	}
	if len(countArgs) < 6 {
		t.Fatalf("expected query args include cust_id, region ids, area ids, and search terms, got %d args", len(countArgs))
	}
}

func TestBuildFindDistributorsByCustIDQuery_SpecificDistributorUsesMapping(t *testing.T) {
	filter := entity.BusinessUnitQueryFilter{
		CustId:       "C22001",
		ParentCustId: "C22001",
		EmployeeId:   77,
		Scope:        entity.EmployeeDropdownScope{DistributorScope: "specific", RegionScope: "all", AreaScope: "all"},
		Page:         1,
		Limit:        10,
	}

	countQuery, _, _, _ := buildFindDistributorsByCustIDQuery(filter)
	if !strings.Contains(countQuery, "mst.m_employee_distributor_mapping edm") {
		t.Fatalf("expected distributor mapping join, got query: %s", countQuery)
	}
	if strings.Contains(countQuery, "md.parent_cust_id = ?") {
		t.Fatalf("did not expect parent scope fallback for specific distributor scope, got query: %s", countQuery)
	}
}

func TestBuildFindDistributorsByCustIDQuery_SpecificRegionAndAreaUsesBothMappings(t *testing.T) {
	filter := entity.BusinessUnitQueryFilter{
		CustId:       "C22001",
		ParentCustId: "C22001",
		EmployeeId:   77,
		Scope:        entity.EmployeeDropdownScope{DistributorScope: "all", RegionScope: "specific", AreaScope: "specific"},
		Page:         1,
		Limit:        10,
	}

	countQuery, _, _, _ := buildFindDistributorsByCustIDQuery(filter)
	if !strings.Contains(countQuery, "mst.m_employee_area_mapping eam") {
		t.Fatalf("expected area mapping join, got query: %s", countQuery)
	}
	if !strings.Contains(countQuery, "mst.m_employee_region_mapping erm") {
		t.Fatalf("expected region mapping join, got query: %s", countQuery)
	}
}
