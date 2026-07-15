package live_monitoring

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"scyllax-pjp/data/request"
	"scyllax-pjp/data/response"

	"github.com/gin-gonic/gin"
)

type controllerServiceStub struct {
	distributorData   []response.LiveMonitoringData
	distributorPaging response.LiveMonitoringPaging
	distributorErr    error
	receivedRequest   request.LiveMonitoringRequest
	receivedCustID    string
	updateData        response.UpdateLocationsResponse
	updateErr         error
}

func (s *controllerServiceStub) GetPrincipalMonitoring(context.Context, request.LiveMonitoringRequest, string) ([]response.LiveMonitoringData, response.LiveMonitoringPaging, error) {
	return nil, response.LiveMonitoringPaging{}, nil
}

func (s *controllerServiceStub) GetDistributorMonitoring(_ context.Context, req request.LiveMonitoringRequest, custID string) ([]response.LiveMonitoringData, response.LiveMonitoringPaging, error) {
	s.receivedRequest = req
	s.receivedCustID = custID
	return s.distributorData, s.distributorPaging, s.distributorErr
}

func (s *controllerServiceStub) GetMonitoringDetail(context.Context, request.LiveMonitoringDetailRequest, string, int64) (*response.LiveMonitoringDetailData, error) {
	return nil, nil
}
func (s *controllerServiceStub) GetUpdateLocations(context.Context, request.UpdateLocationsRequest, string) (response.UpdateLocationsResponse, error) {
	return s.updateData, s.updateErr
}

func TestGetDistributorMonitoring_ReturnsMonitoringPayloadWithoutStaleData(t *testing.T) {
	gin.SetMode(gin.TestMode)

	service := &controllerServiceStub{
		distributorData: []response.LiveMonitoringData{
			{
				EmpID:               358,
				AttendanceID:        nil,
				AttendanceLongitude: 0,
				AttendanceLatitude:  0,
				AttendanceAt:        nil,
				CurrentLongitude:    0,
				CurrentLatitude:     0,
				CurrentCoordinateAt: nil,
				PjpData: []response.LiveMonitoringPjpData{
					{
						PjpID: 9001,
						RouteData: []response.LiveMonitoringRouteData{
							{
								RouteCode: "1201",
								DestinationData: []response.LiveMonitoringDestinationData{
									{
										DestinationID:   501,
										DestinationCode: "OUT-501",
										DestinationName: "Outlet A",
										ArriveAt:        nil,
										LeaveAt:         nil,
										ArriveLongitude: 0,
										ArriveLatitude:  0,
										Start:           nil,
										Finish:          nil,
									},
								},
							},
						},
					},
				},
			},
		},
		distributorPaging: response.LiveMonitoringPaging{TotalRecord: 1, PageCurrent: 1, PageLimit: 10, PageTotal: 1},
	}

	controller := NewLiveMonitoringController(service)
	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)

	query := url.Values{}
	query.Add("date", "1738256400")
	query.Add("status[]", "Approved")
	query.Add("emp_id", "358")

	ctx.Request = httptest.NewRequest(http.MethodGet, "/v1/live-monitoring-distributor?"+query.Encode(), nil)
	ctx.Set("currentCustomerId", "C220010001")

	controller.GetDistributorMonitoring(ctx)

	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusOK)
	}

	if service.receivedCustID != "C220010001" {
		t.Fatalf("received cust_id = %s, want C220010001", service.receivedCustID)
	}
	if service.receivedRequest.Date != 1738256400 {
		t.Fatalf("received date = %d, want 1738256400", service.receivedRequest.Date)
	}

	var payload struct {
		Message string                        `json:"message"`
		Data    []response.LiveMonitoringData `json:"data"`
		Paging  response.LiveMonitoringPaging `json:"paging"`
	}
	if err := json.Unmarshal(recorder.Body.Bytes(), &payload); err != nil {
		t.Fatalf("json.Unmarshal() error = %v", err)
	}

	if payload.Message != "Success" {
		t.Fatalf("message = %q, want Success", payload.Message)
	}
	if len(payload.Data) != 1 {
		t.Fatalf("len(data) = %d, want 1", len(payload.Data))
	}

	monitoring := payload.Data[0]
	if monitoring.AttendanceID != nil || monitoring.AttendanceAt != nil {
		t.Fatalf("attendance should be empty for requested day without attendance, got %+v", monitoring)
	}
	if monitoring.CurrentLongitude != 0 || monitoring.CurrentLatitude != 0 || monitoring.CurrentCoordinateAt != nil {
		t.Fatalf("current coordinate should be reset, got %+v", monitoring)
	}

	destination := monitoring.PjpData[0].RouteData[0].DestinationData[0]
	if destination.ArriveAt != nil || destination.LeaveAt != nil || destination.Start != nil || destination.Finish != nil {
		t.Fatalf("destination tracking should be empty, got %+v", destination)
	}
	if destination.ArriveLongitude != 0 || destination.ArriveLatitude != 0 {
		t.Fatalf("destination arrival coordinate should be reset, got %+v", destination)
	}
	if payload.Paging.TotalRecord != 1 || payload.Paging.PageTotal != 1 {
		t.Fatalf("paging = %+v, want total_record=1 page_total=1", payload.Paging)
	}
}
