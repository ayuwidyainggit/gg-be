package repository

import (
	"strings"
	"testing"
	"time"

	"sales/model"
)

func TestBuildActivityReportGeotagSQLUsesPJPSchema(t *testing.T) {
	t.Parallel()

	for _, schema := range []string{"pjp", "pjp_principles"} {
		sql := buildActivityReportGeotagSQL(schema, false)
		for _, check := range []string{
			schema + ".outlet_visit_list",
			schema + ".permanent_journey_plans",
			`p.cust_id IN ?`,
			`v.date::date BETWEEN ?::date AND ?::date`,
			`location_status = 1`,
			`location_status = 0`,
			`geotag_match_count`,
			`geotag_unmatch_count`,
			`ORDER BY`,
			`p.salesman_name`,
		} {
			if !strings.Contains(sql, check) {
				t.Fatalf("expected geotag SQL for %s to contain %q\nSQL:\n%s", schema, check, sql)
			}
		}
		if strings.Contains(sql, `p.salesman_id = ?`) {
			t.Fatalf("expected no emp filter in default geotag SQL for %s", schema)
		}
	}

	sqlWithEmp := buildActivityReportGeotagSQL("pjp", true)
	if !strings.Contains(sqlWithEmp, `AND p.salesman_id = ?`) {
		t.Fatalf("expected emp filter in geotag SQL with emp_id")
	}
}

func TestActivityReportGeotagSQLUsesYearDateRange(t *testing.T) {
	t.Parallel()

	db, recorded := newReportRepoDryRunDB(t)
	repo := NewReportRepo(db)

	if _, err := repo.ActivityReportGeotag("C26002", []string{"C260020001"}, 2026, nil); err != nil {
		t.Fatalf("ActivityReportGeotag returned error: %v", err)
	}

	query := latestRecordedQuery(t, recorded)
	expectedFrom := time.Date(2026, time.January, 1, 0, 0, 0, 0, time.UTC)
	expectedTo := time.Date(2026, time.December, 31, 0, 0, 0, 0, time.UTC)
	if !strings.Contains(query.SQL, `pjp.outlet_visit_list`) {
		t.Fatalf("expected distributor pjp schema in SQL:\n%s", query.SQL)
	}
	if len(query.Vars) < 3 {
		t.Fatalf("expected at least 3 vars, got %d: %#v", len(query.Vars), query.Vars)
	}
	if len(query.Vars) != 3 {
		t.Fatalf("expected 3 vars without emp_id filter, got %d: %#v", len(query.Vars), query.Vars)
	}
	if query.Vars[1] != expectedFrom {
		t.Fatalf("expected date_start %v, got %v", expectedFrom, query.Vars[1])
	}
	if query.Vars[2] != expectedTo {
		t.Fatalf("expected date_end %v, got %v", expectedTo, query.Vars[2])
	}
}

func TestMergeActivityReportGeotagRowsAggregatesBySalesman(t *testing.T) {
	t.Parallel()

	merged := mergeActivityReportGeotagRows([]model.ActivityReportGeotagRow{
		{SalesmanCode: 1, SalesmanName: "Budi", TotalVisit: 10, GeotagMatchCount: 2, GeotagUnmatchCount: 8},
		{SalesmanCode: 1, SalesmanName: "Budi", TotalVisit: 5, GeotagMatchCount: 1, GeotagUnmatchCount: 4},
	})
	if len(merged) != 1 {
		t.Fatalf("expected 1 merged row, got %d", len(merged))
	}
	if merged[0].TotalVisit != 15 || merged[0].GeotagMatchCount != 3 || merged[0].GeotagUnmatchCount != 12 {
		t.Fatalf("unexpected merged totals: %+v", merged[0])
	}
	if merged[0].GeotagMatchPct != 20 {
		t.Fatalf("expected match pct 20, got %v", merged[0].GeotagMatchPct)
	}
}
