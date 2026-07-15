package repository

import (
	"finance/entity"
	"strings"
	"testing"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func newArRepositoryDryRunDB(t *testing.T) *gorm.DB {
	t.Helper()

	dsn := "host=localhost user=postgres password=postgres dbname=postgres port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{DryRun: true})
	if err != nil {
		t.Fatalf("failed to init dry-run db: %v", err)
	}

	return db
}

func TestBuildCollectorLookupQueries_UsesDepositSourceAndFallbackFields(t *testing.T) {
	db := newArRepositoryDryRunDB(t)
	repo := &RepositoryArImpl{DB: db}
	filter := entity.GeneralQueryFilter{
		CustId: "C22001",
		Page:   2,
		Limit:  70,
		Query:  "451",
		Sort:   "emp_name:asc",
	}

	query, queryCount, limit, page := repo.buildCollectorLookupQueries(filter)
	if limit != 70 {
		t.Fatalf("expected limit 70, got %d", limit)
	}
	if page != 2 {
		t.Fatalf("expected page 2, got %d", page)
	}

	dataTx := query.Find(&[]map[string]interface{}{})
	dataSQL := dataTx.Statement.SQL.String()
	if !strings.Contains(dataSQL, "FROM acf.deposit d") {
		t.Fatalf("expected deposit source in collector lookup SQL, got %s", dataSQL)
	}
	if strings.Contains(dataSQL, "acf.collection") {
		t.Fatalf("expected collector lookup SQL not to depend on acf.collection, got %s", dataSQL)
	}
	if !strings.Contains(dataSQL, "LEFT JOIN mst.m_employee emp ON emp.emp_id = d.emp_id") {
		t.Fatalf("expected employee join using deposit emp_id, got %s", dataSQL)
	}
	if !strings.Contains(dataSQL, "d.deleted_at IS NULL") {
		t.Fatalf("expected soft delete filter in SQL, got %s", dataSQL)
	}
	if !strings.Contains(dataSQL, "d.emp_id IS NOT NULL") {
		t.Fatalf("expected non-null collector filter in SQL, got %s", dataSQL)
	}
	if !strings.Contains(dataSQL, "COALESCE(NULLIF(emp.emp_name, ''), COALESCE(NULLIF(emp.emp_code, ''), CAST(d.emp_id AS varchar))) AS emp_name") {
		t.Fatalf("expected fallback collector name expression in SQL, got %s", dataSQL)
	}
	if !strings.Contains(dataSQL, "CAST(d.emp_id AS varchar) ILIKE") {
		t.Fatalf("expected search to include deposit emp_id fallback, got %s", dataSQL)
	}

	var total int64
	countTx := queryCount.Count(&total)
	countSQL := countTx.Statement.SQL.String()
	if !strings.Contains(countSQL, "COUNT(DISTINCT") {
		t.Fatalf("expected distinct collector count SQL, got %s", countSQL)
	}
}

func TestBuildCollectorLookupSort_FallbackAndWhitelist(t *testing.T) {
	repo := &RepositoryArImpl{}

	tests := []struct {
		name string
		sort string
		want string
	}{
		{
			name: "default when empty",
			sort: "",
			want: "d.emp_id DESC",
		},
		{
			name: "valid sort by name",
			sort: "emp_name:asc",
			want: "emp.emp_name ASC",
		},
		{
			name: "invalid field fallback default",
			sort: "collection_no:desc",
			want: "d.emp_id DESC",
		},
		{
			name: "invalid direction fallback default",
			sort: "emp_id:drop",
			want: "d.emp_id DESC",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := repo.buildCollectorLookupSort(tt.sort)
			if got != tt.want {
				t.Fatalf("buildCollectorLookupSort() = %q, want %q", got, tt.want)
			}
		})
	}
}
