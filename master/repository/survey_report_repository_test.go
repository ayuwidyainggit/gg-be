package repository

import (
	"strings"
	"testing"
	"time"

	"master/entity"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
)

func newSurveyReportRepositoryTest(t *testing.T) *surveyReportRepositoryImpl {
	t.Helper()

	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	t.Cleanup(func() {
		_ = db.Close()
	})

	return &surveyReportRepositoryImpl{DB: sqlx.NewDb(db, "pgx")}
}

func TestBuildSurveyReportBaseQuery_UsesSurveyOwnerCustID(t *testing.T) {
	repo := newSurveyReportRepositoryTest(t)

	query, args, err := repo.buildListBaseQuery(entity.SurveyReportQueryFilter{CustID: "C26004"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(query, "WHERE s.cust_id = $1 AND sa.is_del = false") {
		t.Fatalf("expected query to filter by survey owner cust_id, got: %s", query)
	}
	if strings.Contains(query, "WHERE sa.cust_id = $1") {
		t.Fatalf("expected query not to filter ownership by answer cust_id, got: %s", query)
	}
	if len(args) != 1 || args[0] != "C26004" {
		t.Fatalf("unexpected args: %+v", args)
	}
}

func TestBuildSurveyReportBaseQuery_UsesSurveyIDWithSurveyOwnerScope(t *testing.T) {
	repo := newSurveyReportRepositoryTest(t)

	query, args, err := repo.buildListBaseQuery(entity.SurveyReportQueryFilter{
		CustID:   "C26004",
		SurveyID: []int64{62},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(query, "WHERE s.cust_id = $1 AND sa.is_del = false") {
		t.Fatalf("expected query to keep survey owner scope, got: %s", query)
	}
	if !strings.Contains(query, "AND s.survey_id IN ($2)") {
		t.Fatalf("expected query to filter by survey id on m_survey alias, got: %s", query)
	}
	if strings.Contains(query, "sa.survey_id IN") {
		t.Fatalf("expected query not to filter survey id from answer alias, got: %s", query)
	}
	if len(args) != 2 || args[0] != "C26004" || args[1] != int64(62) {
		t.Fatalf("unexpected args: %+v", args)
	}
}

func TestBuildSurveyReportBaseQuery_ShiftsMultiDigitPlaceholdersWithoutCorruption(t *testing.T) {
	repo := newSurveyReportRepositoryTest(t)
	startDate := time.Date(2026, 4, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2026, 4, 30, 0, 0, 0, 0, time.UTC)

	query, args, err := repo.buildListBaseQuery(entity.SurveyReportQueryFilter{
		CustID:    "C26004",
		StartDate: &startDate,
		EndDate:   &endDate,
		SurveyID:  []int64{68, 104, 63, 89, 74, 101, 62, 84, 77, 100, 71, 90},
		AreaID:    []int64{95, 92},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	for _, expectedPlaceholder := range []string{"$4", "$5", "$6", "$7", "$8", "$9", "$10", "$11", "$12", "$13", "$14", "$15", "$16", "$17"} {
		if !strings.Contains(query, expectedPlaceholder) {
			t.Fatalf("expected query to contain %s, got: %s", expectedPlaceholder, query)
		}
	}
	for _, corruptedPlaceholder := range []string{"$43", "$44", "$45"} {
		if strings.Contains(query, corruptedPlaceholder) {
			t.Fatalf("expected no corrupted placeholder %s, got: %s", corruptedPlaceholder, query)
		}
	}
	if len(args) != 17 {
		t.Fatalf("expected 17 args, got %d: %+v", len(args), args)
	}
}

func TestBuildSurveyReportBaseQuery_UsesInclusiveEndDateRange(t *testing.T) {
	repo := newSurveyReportRepositoryTest(t)
	startDate := time.Date(2026, 4, 24, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2026, 4, 24, 0, 0, 0, 0, time.UTC)

	query, args, err := repo.buildListBaseQuery(entity.SurveyReportQueryFilter{
		CustID:    "C26004",
		StartDate: &startDate,
		EndDate:   &endDate,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(query, "sa.answer_date >= $2::date") {
		t.Fatalf("expected query to use answer_date start boundary, got: %s", query)
	}
	if !strings.Contains(query, "sa.answer_date < ($3::date + INTERVAL '1 day')") {
		t.Fatalf("expected query to use exclusive next-day end boundary, got: %s", query)
	}
	if len(args) != 3 || args[1] != startDate || args[2] != endDate {
		t.Fatalf("unexpected date args: %+v", args)
	}
}

func TestFindExportRows_AttachmentQueryKeepsTenantScopeAndOrder(t *testing.T) {
	repo := newSurveyReportRepositoryTest(t)

	query, args, err := repo.buildExportRowsQuery(entity.SurveyReportQueryFilter{CustID: "C26004"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(args) != 1 || args[0] != "C26004" {
		t.Fatalf("unexpected args: %+v", args)
	}
	if !strings.Contains(query, "WHERE s.cust_id = $1 AND sa.is_del = false") {
		t.Fatalf("expected query to keep tenant scope, got: %s", query)
	}
	if strings.Count(query, "saf.cust_id = sa.cust_id") != 3 {
		t.Fatalf("expected all attachment subqueries to keep tenant scope, got: %s", query)
	}
	if strings.Count(query, "ORDER BY saf.survey_answer_files") != 3 {
		t.Fatalf("expected all attachment subqueries to keep stable order, got: %s", query)
	}
	if strings.Count(query, "COALESCE(NULLIF(saf.file_key, ''), NULLIF(saf.file_name, ''), '')") != 3 {
		t.Fatalf("expected all attachment subqueries to fallback from file_key to file_name, got: %s", query)
	}
}
