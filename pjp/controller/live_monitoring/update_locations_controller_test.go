package live_monitoring

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"scyllax-pjp/data/response"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func TestGetUpdateLocations_ReturnsHappy200ResponseShape(t *testing.T) {
	gin.SetMode(gin.TestMode)
	service := &controllerServiceStub{updateData: response.UpdateLocationsResponse{Timeline: []response.TimelineItem{{Sequence: 1, Type: "gps", RecordedAt: "2026-07-08T09:15:00+07:00"}}}}
	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Request = httptest.NewRequest(http.MethodGet, "/v1/monitoring_locations/update-locations?emp_id=479&date=2026-07-08", nil)
	ctx.Set("currentCustomerId", "C220010001")

	NewLiveMonitoringController(service).GetUpdateLocations(ctx)

	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusOK)
	}
	if body := recorder.Body.String(); body == "" || !containsAll(body, `"message":"Success"`, `"data":{"timeline"`, `"request_id"`) {
		t.Fatalf("body = %s, want success response shape", body)
	}
}

func TestGetUpdateLocations_EmptyTimelineReturnsNoData200(t *testing.T) {
	gin.SetMode(gin.TestMode)
	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Request = httptest.NewRequest(http.MethodGet, "/v1/monitoring_locations/update-locations?emp_id=479", nil)
	ctx.Set("currentCustomerId", "C220010001")

	NewLiveMonitoringController(&controllerServiceStub{updateData: response.UpdateLocationsResponse{Timeline: []response.TimelineItem{}}}).GetUpdateLocations(ctx)

	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusOK)
	}
	if !containsAll(recorder.Body.String(), `"message":"No Data"`, `"timeline"`) {
		t.Fatalf("body = %s, want no-data response", recorder.Body.String())
	}
}

func TestGetUpdateLocations_RecordNotFoundReturns404(t *testing.T) {
	gin.SetMode(gin.TestMode)
	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Request = httptest.NewRequest(http.MethodGet, "/v1/monitoring_locations/update-locations?emp_id=479", nil)
	ctx.Set("currentCustomerId", "C220010001")

	NewLiveMonitoringController(&controllerServiceStub{updateErr: gorm.ErrRecordNotFound}).GetUpdateLocations(ctx)

	if recorder.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusNotFound)
	}
}

func TestGetUpdateLocations_MissingEmployeeIDReturns400(t *testing.T) {
	gin.SetMode(gin.TestMode)
	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Request = httptest.NewRequest(http.MethodGet, "/v1/monitoring_locations/update-locations?date=2026-07-08", nil)
	ctx.Set("currentCustomerId", "C220010001")

	NewLiveMonitoringController(&controllerServiceStub{}).GetUpdateLocations(ctx)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusBadRequest)
	}
	if body := recorder.Body.String(); !contains(body, `"request_id"`) {
		t.Fatalf("body = %s, want request_id", body)
	}
}

func containsAll(value string, parts ...string) bool {
	for _, part := range parts {
		if !contains(value, part) {
			return false
		}
	}
	return true
}

func contains(value, part string) bool {
	for i := 0; i <= len(value)-len(part); i++ {
		if value[i:i+len(part)] == part {
			return true
		}
	}
	return false
}
