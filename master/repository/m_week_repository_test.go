package repository

import (
	"reflect"
	"strings"
	"testing"

	"master/entity"
)

func TestBuildMWeekListQueryFiltersByWorkingDayCalendarIDs(t *testing.T) {
	filter := entity.MWeekQueryFilter{
		WorkingDayCalendarID: []int{7, 9},
		Page:                 1,
		Limit:                10,
	}

	countQuery, countArgs, selectQuery, selectArgs := buildMWeekListQuery(filter, "P001", true)

	if !strings.Contains(countQuery, "working_day_calendar_id IN (?)") {
		t.Fatalf("expected calendar id filter in count query, got %s", countQuery)
	}
	if !strings.Contains(selectQuery, "working_day_calendar_id IN (?)") {
		t.Fatalf("expected calendar id filter in select query, got %s", selectQuery)
	}
	expectedArgs := []interface{}{[]int{7, 9}}
	if !reflect.DeepEqual(countArgs, expectedArgs) {
		t.Fatalf("expected count args %v, got %v", expectedArgs, countArgs)
	}
	if !reflect.DeepEqual(selectArgs, expectedArgs) {
		t.Fatalf("expected select args %v, got %v", expectedArgs, selectArgs)
	}
}

func TestBuildMWeekListQueryUsesParentCustIDForDistributorGeneratedRows(t *testing.T) {
	filter := entity.MWeekQueryFilter{
		ParentCustId: "P001",
		Page:         1,
		Limit:        10,
	}

	countQuery, _, selectQuery, _ := buildMWeekListQuery(filter, "D001", true)

	if !strings.Contains(countQuery, "cust_id = 'P001'") {
		t.Fatalf("expected count query to use parent cust id, got %s", countQuery)
	}
	if !strings.Contains(selectQuery, "cust_id = 'P001'") {
		t.Fatalf("expected select query to use parent cust id, got %s", selectQuery)
	}
}
