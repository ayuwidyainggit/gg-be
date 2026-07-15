package repository

import (
	"strings"
	"testing"

	"master/entity"
)

func areaIntPtr(v int) *int { return &v }

func TestBuildAreaListQuery_SpecificAreaScopeUsesAreaMapping(t *testing.T) {
	filter := entity.AreaQueryFilter{
		CustId:        "C22001",
		ParentCustId:  "C22001",
		EmployeeId:    88,
		DistributorId: 0,
		Scope:         entity.EmployeeDropdownScope{AreaScope: "specific", RegionScope: "all"},
		Query:         "jak",
		IsActive:      areaIntPtr(1),
		Page:          1,
		Limit:         10,
	}

	countQuery, _, _, _ := buildAreaListQuery(filter, false)

	if !strings.Contains(countQuery, "mst.m_employee_area_mapping eam") {
		t.Fatalf("expected area mapping join, got %s", countQuery)
	}
	if !strings.Contains(countQuery, "COUNT(DISTINCT ma.area_id)") {
		t.Fatalf("expected distinct count, got %s", countQuery)
	}
	if !strings.Contains(countQuery, "ma.area_code ILIKE ?") {
		t.Fatalf("expected query alias ma in search clause, got %s", countQuery)
	}
}

func TestBuildAreaListQuery_AllAreaSpecificRegionUsesRegionMapping(t *testing.T) {
	filter := entity.AreaQueryFilter{
		CustId:        "C22001",
		ParentCustId:  "C22001",
		EmployeeId:    88,
		DistributorId: 0,
		Scope:         entity.EmployeeDropdownScope{AreaScope: "all", RegionScope: "specific"},
		RegionId:      []int{1},
		Page:          1,
		Limit:         10,
	}

	countQuery, _, _, _ := buildAreaListQuery(filter, false)

	if !strings.Contains(countQuery, "mst.m_employee_region_mapping erm") {
		t.Fatalf("expected region mapping join, got %s", countQuery)
	}
	if !strings.Contains(countQuery, "ma.region_id IN (?)") {
		t.Fatalf("expected explicit region intersection filter, got %s", countQuery)
	}
}
