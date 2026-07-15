package repository

import (
	"strings"
	"testing"

	"master/entity"
)

func regionIntPtr(v int) *int { return &v }

func TestBuildRegionListQuery_SpecificScopeUsesEmployeeRegionMapping(t *testing.T) {
	filter := entity.RegionQueryFilter{
		CustId:        "C22001",
		ParentCustId:  "C22001",
		EmployeeId:    77,
		DistributorId: 0,
		Scope:         entity.EmployeeDropdownScope{RegionScope: "specific"},
		RegionId:      []int{1, 2},
		Query:         "jak",
		IsActive:      regionIntPtr(1),
		Sort:          "region_id:asc",
		Page:          1,
		Limit:         10,
	}

	countQuery, countArgs, selectQuery, _ := buildRegionListQuery(filter, false)

	if !strings.Contains(countQuery, "mst.m_employee_region_mapping erm") {
		t.Fatalf("expected employee region mapping join, got %s", countQuery)
	}
	if !strings.Contains(countQuery, "erm.emp_id = ?") {
		t.Fatalf("expected employee id filter, got %s", countQuery)
	}
	if !strings.Contains(countQuery, "COUNT(DISTINCT a.region_id)") {
		t.Fatalf("expected distinct count, got %s", countQuery)
	}
	if !strings.Contains(selectQuery, "SELECT DISTINCT") {
		t.Fatalf("expected distinct select, got %s", selectQuery)
	}
	if !strings.Contains(countQuery, "a.region_id IN (?)") {
		t.Fatalf("expected region intersection filter, got %s", countQuery)
	}
	if len(countArgs) < 6 {
		t.Fatalf("expected args for cust, emp, cust scope, region ids, query, active. got %d", len(countArgs))
	}
}

func TestBuildRegionListQuery_AllScopeUsesCustID(t *testing.T) {
	filter := entity.RegionQueryFilter{
		CustId:        "C22001",
		ParentCustId:  "C22001",
		EmployeeId:    77,
		DistributorId: 0,
		Scope:         entity.EmployeeDropdownScope{RegionScope: "all"},
		Page:          1,
		Limit:         10,
	}

	countQuery, _, _, _ := buildRegionListQuery(filter, false)

	if strings.Contains(countQuery, "mst.m_employee_region_mapping") {
		t.Fatalf("did not expect mapping join for all scope, got %s", countQuery)
	}
	if !strings.Contains(countQuery, "a.cust_id = ?") {
		t.Fatalf("expected cust scope filter, got %s", countQuery)
	}
}

func TestBuildRegionListQuery_NonPrincipalUsesParentCustID(t *testing.T) {
	filter := entity.RegionQueryFilter{
		CustId:        "C220010001",
		ParentCustId:  "C22001",
		DistributorId: 99,
		Page:          1,
		Limit:         10,
	}

	countQuery, args, _, _ := buildRegionListQuery(filter, false)
	if strings.Contains(countQuery, "mst.m_employee_region_mapping") {
		t.Fatalf("did not expect principal mapping join for non-principal, got %s", countQuery)
	}
	if len(args) == 0 || args[0] != "C22001" {
		t.Fatalf("expected first scope arg parent cust id C22001, got %+v", args)
	}
}
