package repository

import (
	"strings"
	"testing"

	"master/entity"
)

// helper to build a pointer to int
func intPtr(v int) *int { return &v }

// helper to build a pointer to string
func strPtr(v string) *string { return &v }

// ---------------------------------------------------------------------------
// SX2003-01 — role filter EXISTS clause
// ---------------------------------------------------------------------------

func TestBuildEmployeePJPQuery_AppliesSalesmanRoleExistsFilter(t *testing.T) {
	filter := entity.EmployeePJPQueryFilter{
		CustId: "C22001",
		Page:   1,
		Limit:  10,
	}

	countQuery, countArgs, selectQuery, selectArgs := buildEmployeePJPQuery(filter)

	for _, q := range []struct{ name, query string }{
		{"countQuery", countQuery},
		{"selectQuery", selectQuery},
	} {
		if !strings.Contains(q.query, "EXISTS") {
			t.Errorf("%s: expected EXISTS subquery for role filter, got: %s", q.name, q.query)
		}
		if !strings.Contains(q.query, "sys.m_user") {
			t.Errorf("%s: expected sys.m_user in role filter, got: %s", q.name, q.query)
		}
		if !strings.Contains(q.query, "sys.user_roles") {
			t.Errorf("%s: expected sys.user_roles in role filter, got: %s", q.name, q.query)
		}
		if !strings.Contains(q.query, "sys.m_role") {
			t.Errorf("%s: expected sys.m_role in role filter, got: %s", q.name, q.query)
		}
		if !strings.Contains(q.query, "mu.emp_id = me.emp_id") {
			t.Errorf("%s: expected tenant join mu.emp_id = me.emp_id, got: %s", q.name, q.query)
		}
		if !strings.Contains(q.query, "mu.cust_id = me.cust_id") {
			t.Errorf("%s: expected tenant join mu.cust_id = me.cust_id, got: %s", q.name, q.query)
		}
		if !strings.Contains(q.query, "LOWER(mr.role_name) = 'salesman'") {
			t.Errorf("%s: expected LOWER(mr.role_name) = 'salesman', got: %s", q.name, q.query)
		}
		// Must NOT have a raw JOIN sys.m_role in the main FROM clause
		mainFrom := strings.SplitN(q.query, "EXISTS", 2)[0]
		if strings.Contains(mainFrom, "JOIN sys.m_role") {
			t.Errorf("%s: expected no raw JOIN sys.m_role in main query (use EXISTS), got: %s", q.name, q.query)
		}
	}

	// args must be equal length for count and select
	if len(countArgs) != len(selectArgs) {
		t.Errorf("expected countArgs len == selectArgs len, got count=%d select=%d", len(countArgs), len(selectArgs))
	}
}

// ---------------------------------------------------------------------------
// SX2003-02 — distributor scope uses smc.m_customer + parent_cust_id
// ---------------------------------------------------------------------------

func TestBuildEmployeePJPQuery_DistributorScopeUsesCustomerMapping(t *testing.T) {
	distId := 67
	filter := entity.EmployeePJPQueryFilter{
		CustId:        "C260020001",
		ParentCustId:  "C26002",
		DistributorId: &distId,
		Page:          1,
		Limit:         10,
	}

	countQuery, countArgs, _, _ := buildEmployeePJPQuery(filter)

	if !strings.Contains(countQuery, "smc.m_customer") {
		t.Errorf("expected smc.m_customer in distributor scope, got: %s", countQuery)
	}
	if !strings.Contains(countQuery, "mc.distributor_id") {
		t.Errorf("expected mc.distributor_id in distributor scope, got: %s", countQuery)
	}
	if !strings.Contains(countQuery, "mc.parent_cust_id") {
		t.Errorf("expected mc.parent_cust_id constraint in distributor scope, got: %s", countQuery)
	}

	// distributor_id and parent_cust_id must be in args, not concatenated
	foundDistId := false
	foundParent := false
	for _, a := range countArgs {
		switch v := a.(type) {
		case int:
			if v == 67 {
				foundDistId = true
			}
		case string:
			if v == "C26002" {
				foundParent = true
			}
		}
	}
	if !foundDistId {
		t.Errorf("expected distributor_id=67 in query args, got args: %v", countArgs)
	}
	if !foundParent {
		t.Errorf("expected parent_cust_id='C26002' in query args, got args: %v", countArgs)
	}
}

// ---------------------------------------------------------------------------
// SX2003-03 — FilterCustId is parameterized, not concatenated
// ---------------------------------------------------------------------------

func TestBuildEmployeePJPQuery_FilterCustIDUsesArgsNotConcatenation(t *testing.T) {
	injectionValue := "C26002' OR '1'='1"
	filter := entity.EmployeePJPQueryFilter{
		CustId:       "C22001",
		FilterCustId: strPtr(injectionValue),
		Page:         1,
		Limit:        10,
	}

	countQuery, countArgs, selectQuery, _ := buildEmployeePJPQuery(filter)

	// The raw injection string must NOT appear literally in the query
	if strings.Contains(countQuery, injectionValue) {
		t.Errorf("countQuery contains raw injection string — not parameterized: %s", countQuery)
	}
	if strings.Contains(selectQuery, injectionValue) {
		t.Errorf("selectQuery contains raw injection string — not parameterized: %s", selectQuery)
	}

	// The value must appear in args
	found := false
	for _, a := range countArgs {
		if s, ok := a.(string); ok && s == injectionValue {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected injection-shaped cust_id in args, got args: %v", countArgs)
	}
}

// ---------------------------------------------------------------------------
// SX2003-04 — default scope (no distributor, no FilterCustId) uses JWT CustId
// ---------------------------------------------------------------------------

func TestBuildEmployeePJPQuery_DefaultScopeUsesJWTCustId(t *testing.T) {
	filter := entity.EmployeePJPQueryFilter{
		CustId: "C22001",
		Page:   1,
		Limit:  10,
	}

	countQuery, countArgs, _, _ := buildEmployeePJPQuery(filter)

	if !strings.Contains(countQuery, "me.cust_id = ?") {
		t.Errorf("expected parameterized me.cust_id = ? in default scope, got: %s", countQuery)
	}

	found := false
	for _, a := range countArgs {
		if s, ok := a.(string); ok && s == "C22001" {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected JWT cust_id 'C22001' in args, got: %v", countArgs)
	}
}

// ---------------------------------------------------------------------------
// SX2003-05 — search query is parameterized, not concatenated
// ---------------------------------------------------------------------------

func TestBuildEmployeePJPQuery_SearchQueryIsParameterized(t *testing.T) {
	filter := entity.EmployeePJPQueryFilter{
		CustId: "C22001",
		Query:  "john",
		Page:   1,
		Limit:  10,
	}

	countQuery, countArgs, _, _ := buildEmployeePJPQuery(filter)

	if strings.Contains(countQuery, "'%john%'") || strings.Contains(countQuery, "john") {
		t.Errorf("expected search value not concatenated in query, got: %s", countQuery)
	}
	if !strings.Contains(countQuery, "ILIKE ?") {
		t.Errorf("expected ILIKE ? placeholder in query, got: %s", countQuery)
	}

	found := false
	for _, a := range countArgs {
		if s, ok := a.(string); ok && s == "%john%" {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected '%%john%%' in args, got: %v", countArgs)
	}
}

// ---------------------------------------------------------------------------
// SX2003-06 — is_active filter
// ---------------------------------------------------------------------------

func TestBuildEmployeePJPQuery_IsActiveFilterApplied(t *testing.T) {
	filter := entity.EmployeePJPQueryFilter{
		CustId:   "C22001",
		IsActive: []int{1},
		Page:     1,
		Limit:    10,
	}

	countQuery, countArgs, _, _ := buildEmployeePJPQuery(filter)

	if !strings.Contains(countQuery, "me.is_active = ?") {
		t.Errorf("expected me.is_active = ? in query, got: %s", countQuery)
	}

	found := false
	for _, a := range countArgs {
		if b, ok := a.(bool); ok && b == true {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected is_active=true in args, got: %v", countArgs)
	}
}

// ---------------------------------------------------------------------------
// SX2003-07 — pagination and sort
// ---------------------------------------------------------------------------

func TestBuildEmployeePJPQuery_PaginationAndSort(t *testing.T) {
	filter := entity.EmployeePJPQueryFilter{
		CustId: "C22001",
		Page:   2,
		Limit:  20,
		Sort:   "emp_name:asc",
	}

	_, _, selectQuery, _ := buildEmployeePJPQuery(filter)

	if !strings.Contains(selectQuery, "LIMIT 20 OFFSET 20") {
		t.Errorf("expected LIMIT 20 OFFSET 20, got: %s", selectQuery)
	}
	if !strings.Contains(selectQuery, "me.emp_name asc") {
		t.Errorf("expected ORDER BY me.emp_name asc, got: %s", selectQuery)
	}
}

func TestBuildEmployeePJPQuery_DefaultSortAndPagination(t *testing.T) {
	filter := entity.EmployeePJPQueryFilter{
		CustId: "C22001",
	}

	_, _, selectQuery, _ := buildEmployeePJPQuery(filter)

	if !strings.Contains(selectQuery, "me.emp_id DESC") {
		t.Errorf("expected default ORDER BY me.emp_id DESC, got: %s", selectQuery)
	}
	if !strings.Contains(selectQuery, "LIMIT 9999 OFFSET 0") {
		t.Errorf("expected default LIMIT 9999 OFFSET 0, got: %s", selectQuery)
	}
}

// ---------------------------------------------------------------------------
// SX2003-sort-01 — injection-shaped sort input is rejected
// ---------------------------------------------------------------------------

func TestBuildEmployeePJPSortClause_RejectsInjectionInput(t *testing.T) {
	result := buildEmployeePJPSortClause("emp_id:asc; DROP TABLE mst.m_employee--")

	if strings.Contains(result, "DROP") {
		t.Errorf("expected injection payload to be rejected, got: %s", result)
	}
	if strings.Contains(result, ";") {
		t.Errorf("expected semicolon to be stripped from sort output, got: %s", result)
	}
	// direction "asc; DROP TABLE mst.m_employee--" is not asc/desc → whole token rejected → fallback
	if result != "me.emp_id DESC" {
		t.Errorf("expected fallback to me.emp_id DESC for injection input, got: %s", result)
	}
}

// ---------------------------------------------------------------------------
// SX2003-sort-02 — unknown column falls back to default
// ---------------------------------------------------------------------------

func TestBuildEmployeePJPSortClause_FallsBackToDefaultForUnknownColumn(t *testing.T) {
	result := buildEmployeePJPSortClause("unknown_col:asc")

	if result != "me.emp_id DESC" {
		t.Errorf("expected fallback to me.emp_id DESC for unknown column, got: %s", result)
	}
}

// ---------------------------------------------------------------------------
// SX2003-sort-03 — valid column is accepted and qualified
// ---------------------------------------------------------------------------

func TestBuildEmployeePJPSortClause_AcceptsValidColumn(t *testing.T) {
	result := buildEmployeePJPSortClause("emp_name:desc")

	if !strings.Contains(result, "me.emp_name desc") {
		t.Errorf("expected me.emp_name desc in sort output, got: %s", result)
	}
	if strings.Contains(result, "me.emp_id DESC") {
		t.Errorf("expected no fallback default when valid column given, got: %s", result)
	}
}

// ---------------------------------------------------------------------------
// SX2003-08 — no duplicate employee: EXISTS not raw JOIN in main FROM
// ---------------------------------------------------------------------------

func TestBuildEmployeePJPQuery_NoRawRoleJoinInMainQuery(t *testing.T) {
	filter := entity.EmployeePJPQueryFilter{
		CustId: "C22001",
		Page:   1,
		Limit:  10,
	}

	countQuery, _, selectQuery, _ := buildEmployeePJPQuery(filter)

	// Split at EXISTS to get the main FROM/WHERE portion before the subquery
	for _, q := range []struct{ name, query string }{
		{"countQuery", countQuery},
		{"selectQuery", selectQuery},
	} {
		mainPart := strings.SplitN(q.query, "EXISTS", 2)[0]
		if strings.Contains(mainPart, "JOIN sys.user_roles") {
			t.Errorf("%s: raw JOIN sys.user_roles in main query would cause duplicates, got: %s", q.name, mainPart)
		}
		if strings.Contains(mainPart, "JOIN sys.m_role") {
			t.Errorf("%s: raw JOIN sys.m_role in main query would cause duplicates, got: %s", q.name, mainPart)
		}
	}
}
