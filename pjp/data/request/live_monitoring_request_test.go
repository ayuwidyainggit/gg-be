package request

import (
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestLiveMonitoringRequest_GetEmpIDs(t *testing.T) {
	tests := []struct {
		name        string
		empIDs      []int
		legacyEmpID string
		expected    []int
	}{
		{"empty input", nil, "", nil},
		{"single array value", []int{210}, "", []int{210}},
		{"multiple array values", []int{210, 358}, "", []int{210, 358}},
		{"legacy comma separated values", nil, "210,358", []int{210, 358}},
		{"legacy values with spaces and invalid tokens", nil, "210, 358, abc", []int{210, 358}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := LiveMonitoringRequest{
				EmpIDs:      tt.empIDs,
				LegacyEmpID: tt.legacyEmpID,
			}
			got := req.GetEmpIDs()
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("GetEmpIDs() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestLiveMonitoringRequest_ShouldBindQuery_DistributorContract(t *testing.T) {
	gin.SetMode(gin.TestMode)

	query := url.Values{}
	query.Add("date", "1738195200")
	query.Add("emp_id[]", "358")
	query.Add("emp_id[]", "360")
	query.Add("status[]", "Approved")
	query.Add("status[]", "Need Review")
	query.Add("page", "2")
	query.Add("limit", "25")

	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	request := httptest.NewRequest("GET", "/v1/live-monitoring-distributor?"+query.Encode(), nil)
	ctx.Request = request

	var req LiveMonitoringRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		t.Fatalf("ShouldBindQuery() error = %v", err)
	}

	if req.Date != 1738195200 {
		t.Fatalf("Date = %d, want %d", req.Date, int64(1738195200))
	}

	expectedEmpIDs := []int{358, 360}
	if !reflect.DeepEqual(req.EmpIDs, expectedEmpIDs) {
		t.Fatalf("EmpIDs = %v, want %v", req.EmpIDs, expectedEmpIDs)
	}

	expectedStatus := []string{"Approved", "Need Review"}
	if !reflect.DeepEqual(req.Status, expectedStatus) {
		t.Fatalf("Status = %v, want %v", req.Status, expectedStatus)
	}

	if req.Page != 2 {
		t.Fatalf("Page = %d, want %d", req.Page, 2)
	}

	if req.Limit != 25 {
		t.Fatalf("Limit = %d, want %d", req.Limit, 25)
	}

	if req.RegionID != 0 || req.AreaID != 0 || req.DistributorID != 0 {
		t.Fatalf("region/area/distributor filters should remain optional, got region=%d area=%d distributor=%d", req.RegionID, req.AreaID, req.DistributorID)
	}
}

func TestLiveMonitoringRequest_ShouldBindQuery_DistributorContract_WithoutOptionalFilters(t *testing.T) {
	gin.SetMode(gin.TestMode)

	query := url.Values{}
	query.Add("date", "1738195200")
	query.Add("status[]", "Approved")

	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	request := httptest.NewRequest("GET", "/v1/live-monitoring-distributor?"+query.Encode(), nil)
	ctx.Request = request

	var req LiveMonitoringRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		t.Fatalf("ShouldBindQuery() error = %v", err)
	}

	if req.Page != 0 {
		t.Fatalf("Page = %d, want %d before service defaulting", req.Page, 0)
	}

	if req.Limit != 0 {
		t.Fatalf("Limit = %d, want %d before service defaulting", req.Limit, 0)
	}

	if req.EmpIDs != nil {
		t.Fatalf("EmpIDs = %v, want nil when query is omitted", req.EmpIDs)
	}
}
