package live_monitoring

import (
	"context"
	"scyllax-pjp/data/response"
	"testing"

	"scyllax-pjp/data/request"
	"scyllax-pjp/model"

	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type distributorRepoStub struct {
	childCustIDs             []string
	childCustIDsErr          error
	employeeIDs              []int
	employeeIDsErr           error
	employeeMetaMap          map[int]model.DistributorEmployeeMetaRow
	routeMetaMap             map[string]model.DistributorRouteMetaRow
	outletMetaMap            map[string]model.DistributorOutletMetaRow
	receivedScopeCustIDs     []string
	receivedVisitEmpIDs      []int
	receivedMonitoringEmpIDs []int
	receivedScopedEmpIDs     []int
	receivedMetaEmpIDs       []int
	receivedMetaRouteCodes   []int64
	receivedMetaOutletIDs    []int
	receivedCountCustIDs     []string
	countResult              int64
}

func (r *distributorRepoStub) GetPrincipalEmployeeIDs(context.Context, *gorm.DB, []string, string, int, int, int, []int, []string) ([]int, error) {
	return nil, nil
}

func (r *distributorRepoStub) GetPrincipalMonitoring(context.Context, *gorm.DB, []string, string, int, int, int, []int, []string, int, int) ([]model.LiveMonitoringPrincipalRow, error) {
	return nil, nil
}

func (r *distributorRepoStub) GetPrincipalExtraCalls(context.Context, *gorm.DB, []string, string, int, int, int, []int, []string) ([]model.LiveMonitoringPrincipalRow, error) {
	return nil, nil
}

func (r *distributorRepoStub) CountPrincipalMonitoring(context.Context, *gorm.DB, []string, string, int, int, int, []int, []string) (int64, error) {
	return 0, nil
}

func (r *distributorRepoStub) GetDistributorEmployeeIDs(_ context.Context, _ *gorm.DB, custIDs []string, _ string, _, _, _ int, _ []int, _ []string) ([]int, error) {
	r.receivedScopeCustIDs = append([]string(nil), custIDs...)
	return r.employeeIDs, r.employeeIDsErr
}

func (r *distributorRepoStub) GetDistributorMonitoring(_ context.Context, _ *gorm.DB, _ []string, _ string, _, _, _ int, empIDs []int, _ []string, _, _ int) ([]model.LiveMonitoringDistributorRow, error) {
	r.receivedMonitoringEmpIDs = append([]int(nil), empIDs...)
	return nil, nil
}

func (r *distributorRepoStub) GetDistributorLatestVisitCoordinates(_ context.Context, _ *gorm.DB, _ []string, _ string, empIDs []int) (map[string]model.LatestVisitCoordinateRow, error) {
	r.receivedVisitEmpIDs = append([]int(nil), empIDs...)
	return map[string]model.LatestVisitCoordinateRow{}, nil
}

func (r *distributorRepoStub) GetDistributorEmployeeMeta(_ context.Context, _ *gorm.DB, _ []string, empIDs []int) (map[int]model.DistributorEmployeeMetaRow, error) {
	r.receivedMetaEmpIDs = append([]int(nil), empIDs...)
	if r.employeeMetaMap == nil {
		return map[int]model.DistributorEmployeeMetaRow{}, nil
	}
	return r.employeeMetaMap, nil
}

func (r *distributorRepoStub) GetDistributorRouteMeta(_ context.Context, _ *gorm.DB, _ []string, routeCodes []int64) (map[string]model.DistributorRouteMetaRow, error) {
	r.receivedMetaRouteCodes = append([]int64(nil), routeCodes...)
	if r.routeMetaMap == nil {
		return map[string]model.DistributorRouteMetaRow{}, nil
	}
	return r.routeMetaMap, nil
}

func (r *distributorRepoStub) GetDistributorOutletMeta(_ context.Context, _ *gorm.DB, _ []string, outletIDs []int) (map[string]model.DistributorOutletMetaRow, error) {
	r.receivedMetaOutletIDs = append([]int(nil), outletIDs...)
	if r.outletMetaMap == nil {
		return map[string]model.DistributorOutletMetaRow{}, nil
	}
	return r.outletMetaMap, nil
}

func (r *distributorRepoStub) GetDistributorAttendance(_ context.Context, _ *gorm.DB, _ []string, _ string, _, _, _ int, empIDs []int, _ []string) (map[int]model.AttendanceRow, error) {
	r.receivedScopedEmpIDs = append([]int(nil), empIDs...)
	return map[int]model.AttendanceRow{}, nil
}

func (r *distributorRepoStub) GetDistributorCurrentCoordinates(context.Context, *gorm.DB, []string, string, int, int, int, []int, []string) (map[int]model.CurrentCoordinateRow, error) {
	return map[int]model.CurrentCoordinateRow{}, nil
}

func (r *distributorRepoStub) CountDistributorMonitoring(_ context.Context, _ *gorm.DB, custIDs []string, _ string, _, _, _ int, _ []int, _ []string) (int64, error) {
	r.receivedCountCustIDs = append([]string(nil), custIDs...)
	return r.countResult, nil
}

func (r *distributorRepoStub) GetVisitInformationPrincipal(context.Context, *gorm.DB, []string, string, int) (*model.VisitInformationRow, error) {
	return nil, nil
}

func (r *distributorRepoStub) GetVisitInformationPrincipalFromHistory(context.Context, *gorm.DB, []string, string, int) (*model.VisitInformationRow, error) {
	return nil, nil
}

func (r *distributorRepoStub) CountTotalVisitsPrincipal(context.Context, *gorm.DB, []string, string, int) (int64, error) {
	return 0, nil
}

func (r *distributorRepoStub) GetVisitInformationDistributor(context.Context, *gorm.DB, string, int, int) (*model.VisitInformationRow, error) {
	return nil, nil
}

func (r *distributorRepoStub) CountDistributorPlannedVisits(context.Context, *gorm.DB, string, int, int) (int64, error) {
	return 0, nil
}

func (r *distributorRepoStub) CountDistributorExtraCalls(context.Context, *gorm.DB, string, int, int) (int64, error) {
	return 0, nil
}

func (r *distributorRepoStub) CountDistributorOnGoingVisits(context.Context, *gorm.DB, string, int, int) (int64, error) {
	return 0, nil
}

func (r *distributorRepoStub) CountDistributorVisitedVisits(context.Context, *gorm.DB, string, int, int) (int64, error) {
	return 0, nil
}

func (r *distributorRepoStub) CountDistributorSkippedVisits(context.Context, *gorm.DB, string, int, int) (int64, error) {
	return 0, nil
}

func (r *distributorRepoStub) GetSales(context.Context, *gorm.DB, []string, string, int) ([]model.SalesRow, error) {
	return nil, nil
}

func (r *distributorRepoStub) GetReturns(context.Context, *gorm.DB, []string, string, int) ([]model.ReturnRow, error) {
	return nil, nil
}

func (r *distributorRepoStub) GetCollections(context.Context, *gorm.DB, []string, string, int) ([]model.CollectionRow, error) {
	return nil, nil
}

func (r *distributorRepoStub) GetExpenses(context.Context, *gorm.DB, string, int, string) ([]model.ExpenseRow, error) {
	return nil, nil
}

func (r *distributorRepoStub) GetShipments(context.Context, *gorm.DB, []string, string, int) ([]model.ShipmentRow, error) {
	return nil, nil
}
func (r *distributorRepoStub) GetSubmittedSurveyData(context.Context, *gorm.DB, []string, string, int) ([]model.SurveyDataRow, error) {
	return nil, nil
}
func (r *distributorRepoStub) GetActivityTime(context.Context, *gorm.DB, string, int) (*string, error) {

	return nil, nil
}

func (r *distributorRepoStub) GetDistributorInfo(context.Context, *gorm.DB, int) (*model.DistributorInfoRow, error) {
	return nil, nil
}

func (r *distributorRepoStub) GetUserFullname(context.Context, *gorm.DB, string) (*string, error) {
	return nil, nil
}

func (r *distributorRepoStub) GetChildCustIDs(context.Context, *gorm.DB, string) ([]string, error) {
	return r.childCustIDs, r.childCustIDsErr
}

func (r *distributorRepoStub) GetSalesmanCustID(context.Context, *gorm.DB, int) (string, error) {
	return "", nil
}

func (r *distributorRepoStub) GetEmployeeRole(context.Context, *gorm.DB, int, string) (string, error) {
	return "", nil
}

func (r *distributorRepoStub) GetUpdateLocations(context.Context, *gorm.DB, int, string, string, string) ([]model.UpdateLocationRow, error) {
	return nil, nil
}

func TestGetDistributorMonitoring_UsesChildCustIDsForPrincipalScope(t *testing.T) {
	repo := &distributorRepoStub{
		childCustIDs: []string{"C22001", "C220010001"},
		employeeIDs:  []int{11, 12},
		countResult:  0,
	}

	svc := &liveMonitoringService{repository: repo, validate: validator.New()}

	_, _, err := svc.GetDistributorMonitoring(context.Background(), request.LiveMonitoringRequest{
		Date:   1776081600,
		Status: []string{"Approved"},
	}, "C22001")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(repo.receivedScopeCustIDs) != 2 || repo.receivedScopeCustIDs[0] != "C22001" || repo.receivedScopeCustIDs[1] != "C220010001" {
		t.Fatalf("scope custIDs = %#v, want [C22001 C220010001]", repo.receivedScopeCustIDs)
	}

	if len(repo.receivedScopedEmpIDs) != 2 || repo.receivedScopedEmpIDs[0] != 11 || repo.receivedScopedEmpIDs[1] != 12 {
		t.Fatalf("scoped empIDs = %#v, want [11 12]", repo.receivedScopedEmpIDs)
	}
}

func TestGetDistributorMonitoring_FallsBackToCurrentCustIDWhenNoChildrenFound(t *testing.T) {
	repo := &distributorRepoStub{
		childCustIDs: []string{},
		employeeIDs:  []int{33},
		countResult:  0,
	}

	svc := &liveMonitoringService{repository: repo, validate: validator.New()}

	_, _, err := svc.GetDistributorMonitoring(context.Background(), request.LiveMonitoringRequest{
		Date:   1776081600,
		Status: []string{"Approved"},
	}, "C220010001")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(repo.receivedScopeCustIDs) != 1 || repo.receivedScopeCustIDs[0] != "C220010001" {
		t.Fatalf("scope custIDs = %#v, want [C220010001]", repo.receivedScopeCustIDs)
	}
}

func TestGetDistributorMonitoring_PaginatesEmployeeScopeBeforeDetailQueries(t *testing.T) {
	repo := &distributorRepoStub{
		childCustIDs: []string{"C22001"},
		employeeIDs:  []int{11, 12, 13},
		employeeMetaMap: map[int]model.DistributorEmployeeMetaRow{
			12: {
				EmpID:         12,
				EmpCode:       "E12",
				EmpName:       "Emp 12",
				DistributorID: 67,
				AreaID:        82,
				RegionID:      67,
			},
		},
	}

	svc := &liveMonitoringService{repository: repo, validate: validator.New()}

	_, paging, err := svc.GetDistributorMonitoring(context.Background(), request.LiveMonitoringRequest{
		Date:   1776081600,
		Status: []string{"Approved"},
		Page:   2,
		Limit:  1,
	}, "C22001")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(repo.receivedMonitoringEmpIDs) != 1 || repo.receivedMonitoringEmpIDs[0] != 12 {
		t.Fatalf("monitoring empIDs = %#v, want [12]", repo.receivedMonitoringEmpIDs)
	}

	if len(repo.receivedScopedEmpIDs) != 1 || repo.receivedScopedEmpIDs[0] != 12 {
		t.Fatalf("attendance empIDs = %#v, want [12]", repo.receivedScopedEmpIDs)
	}

	if len(repo.receivedVisitEmpIDs) != 1 || repo.receivedVisitEmpIDs[0] != 12 {
		t.Fatalf("visit empIDs = %#v, want [12]", repo.receivedVisitEmpIDs)
	}

	if len(repo.receivedMetaEmpIDs) != 1 || repo.receivedMetaEmpIDs[0] != 12 {
		t.Fatalf("meta empIDs = %#v, want [12]", repo.receivedMetaEmpIDs)
	}

	if paging.PageTotal != 3 || paging.PageCurrent != 2 || paging.PageLimit != 1 {
		t.Fatalf("paging = %#v, want total=3 current=2 limit=1", paging)
	}
}

func TestEnrichDistributorRowsWithMetadata(t *testing.T) {
	rows := []model.LiveMonitoringDistributorRow{
		{
			CustID:    "C22001",
			EmpID:     12,
			RouteCode: 4135,
			OutletID:  1454,
		},
	}

	employeeMetaMap := map[int]model.DistributorEmployeeMetaRow{
		12: {
			EmpID:         12,
			EmpCode:       "E12",
			EmpName:       "Emp 12",
			DistributorID: 67,
			AreaID:        82,
			RegionID:      67,
		},
	}

	routeMetaMap := map[string]model.DistributorRouteMetaRow{
		buildDistributorRouteMetaKey("C22001", 4135): {
			CustID:    "C22001",
			RouteCode: 4135,
			RouteName: "Route 3",
		},
	}

	outletMetaMap := map[string]model.DistributorOutletMetaRow{
		buildDistributorOutletMetaKey("C22001", 1454): {
			CustID:     "C22001",
			OutletID:   1454,
			OutletCode: "1234567890",
			OutletName: "Toko Jan",
		},
	}

	enrichDistributorRowsWithMetadata(rows, employeeMetaMap, routeMetaMap, outletMetaMap)

	if rows[0].EmpCode != "E12" || rows[0].EmpName != "Emp 12" {
		t.Fatalf("employee metadata not enriched: %#v", rows[0])
	}

	if rows[0].RouteName != "Route 3" {
		t.Fatalf("route metadata not enriched: %#v", rows[0])
	}

	if rows[0].OutletCode != "1234567890" || rows[0].OutletName != "Toko Jan" {
		t.Fatalf("outlet metadata not enriched: %#v", rows[0])
	}
}

func TestEnrichDistributorMonitoringData_SelectsLatestValidCoordinateSource(t *testing.T) {
	attendanceAt := int64(1776132000)
	arriveAt := int64(1776132300)
	leaveAt := int64(1776132600)
	checkoutAt := int64(1776132942)

	result := []response.LiveMonitoringData{{
		EmpID: 360,
		PjpData: []response.LiveMonitoringPjpData{{
			RouteData: []response.LiveMonitoringRouteData{{
				DestinationData: []response.LiveMonitoringDestinationData{{
					ArriveAt: &arriveAt,
					LeaveAt:  &leaveAt,
				}},
			}},
		}},
	}}

	attendanceMap := map[int]model.AttendanceRow{
		360: {
			AttendanceID: ptrInt64(99),
			Timestamp:    &attendanceAt,
			Longitude:    106.7001,
			Latitude:     -6.2001,
			ClockOutID:   ptrInt64(100),
			ClockOutAt:   &checkoutAt,
			ClockOutLong: 106.7301,
			ClockOutLat:  -6.2301,
		},
	}

	currentCoordinateMap := map[int]model.CurrentCoordinateRow{
		360: {
			Longitude: 106.7301,
			Latitude:  -6.2301,
			Timestamp: &checkoutAt,
			Source:    "attendance_checkout",
		},
	}

	enrichDistributorMonitoringData("2026-04-14", result, attendanceMap, currentCoordinateMap)

	if result[0].CurrentCoordinateSource != "attendance_checkout" {
		t.Fatalf("current source = %s, want attendance_checkout", result[0].CurrentCoordinateSource)
	}

	if result[0].CurrentCoordinateAt == nil || *result[0].CurrentCoordinateAt != checkoutAt {
		t.Fatalf("current coordinate at = %#v, want %d", result[0].CurrentCoordinateAt, checkoutAt)
	}

	if result[0].CurrentLongitude != 106.7301 || result[0].CurrentLatitude != -6.2301 {
		t.Fatalf("current coordinate = (%v,%v), want (106.7301,-6.2301)", result[0].CurrentLongitude, result[0].CurrentLatitude)
	}

	if result[0].AttendanceID == nil || *result[0].AttendanceID != 99 {
		t.Fatalf("attendance id = %#v, want 99", result[0].AttendanceID)
	}

	if result[0].ClockOut == nil || *result[0].ClockOut != 100 {
		t.Fatalf("clock out = %#v, want 100", result[0].ClockOut)
	}

	if result[0].ClockOutAt == nil || *result[0].ClockOutAt != checkoutAt {
		t.Fatalf("clock out at = %#v, want %d", result[0].ClockOutAt, checkoutAt)
	}

	if result[0].ClockOutLongitude != 106.7301 || result[0].ClockOutLatitude != -6.2301 {
		t.Fatalf("clock out coordinate = (%v,%v), want (106.7301,-6.2301)", result[0].ClockOutLongitude, result[0].ClockOutLatitude)
	}
}

func TestEnrichDistributorMonitoringData_LeavesClockOutEmptyWhenCheckoutMissing(t *testing.T) {
	attendanceAt := int64(1776132000)

	result := []response.LiveMonitoringData{{EmpID: 360}}
	attendanceMap := map[int]model.AttendanceRow{
		360: {
			AttendanceID: ptrInt64(99),
			Timestamp:    &attendanceAt,
			Longitude:    106.7001,
			Latitude:     -6.2001,
		},
	}

	enrichDistributorMonitoringData("2026-04-14", result, attendanceMap, map[int]model.CurrentCoordinateRow{})

	if result[0].AttendanceID == nil || *result[0].AttendanceID != 99 {
		t.Fatalf("attendance id = %#v, want 99", result[0].AttendanceID)
	}

	if result[0].ClockOut != nil || result[0].ClockOutAt != nil {
		t.Fatalf("clock out fields should be empty, got %+v", result[0])
	}
}

func TestEnrichDistributorMonitoringData_FallsBackWhenCurrentCoordinateOutsideBusinessDay(t *testing.T) {
	attendanceAt := int64(1776132000)
	staleCheckoutAt := int64(1776045600)

	result := []response.LiveMonitoringData{{EmpID: 360}}
	attendanceMap := map[int]model.AttendanceRow{
		360: {
			AttendanceID: ptrInt64(99),
			Timestamp:    &attendanceAt,
			Longitude:    106.7001,
			Latitude:     -6.2001,
		},
	}
	currentCoordinateMap := map[int]model.CurrentCoordinateRow{
		360: {
			Longitude: 106.7301,
			Latitude:  -6.2301,
			Timestamp: &staleCheckoutAt,
			Source:    "attendance_checkout",
		},
	}

	enrichDistributorMonitoringData("2026-04-14", result, attendanceMap, currentCoordinateMap)

	if result[0].CurrentCoordinateAt != nil {
		t.Fatalf("current coordinate at = %#v, want nil", result[0].CurrentCoordinateAt)
	}

	if result[0].CurrentCoordinateSource != "" {
		t.Fatalf("current coordinate source = %q, want empty", result[0].CurrentCoordinateSource)
	}
}

func TestEnrichDistributorMonitoringData_ResetsWhenNoAttendanceOnRequestedDate(t *testing.T) {
	attendanceAt := int64(1776045600)
	checkoutAt := int64(1776132942)

	result := []response.LiveMonitoringData{{
		EmpID:               360,
		CurrentLongitude:    106.8,
		CurrentLatitude:     -6.2,
		CurrentCoordinateAt: &checkoutAt,
		AttendanceID:        ptrInt64(88),
		AttendanceLongitude: 106.7,
		AttendanceLatitude:  -6.1,
		AttendanceAt:        &attendanceAt,
		PjpData:             []response.LiveMonitoringPjpData{},
	}}

	attendanceMap := map[int]model.AttendanceRow{
		360: {
			AttendanceID: ptrInt64(88),
			Timestamp:    &attendanceAt,
			Longitude:    106.7,
			Latitude:     -6.1,
		},
	}
	currentCoordinateMap := map[int]model.CurrentCoordinateRow{
		360: {
			Longitude: 106.7301,
			Latitude:  -6.2301,
			Timestamp: &checkoutAt,
			Source:    "attendance_checkout",
		},
	}

	enrichDistributorMonitoringData("2026-04-14", result, attendanceMap, currentCoordinateMap)

	if result[0].AttendanceID != nil || result[0].CurrentCoordinateAt != nil {
		t.Fatalf("daily tracking state not reset: %#v", result[0])
	}
}

func TestEnrichDistributorRowsWithLatestVisits_AssignsFileURL(t *testing.T) {
	expectedFileURL := "https://files.example.com/arrive/distributor-123.jpg"
	rows := []model.LiveMonitoringDistributorRow{{
		CustID:       "C22001",
		SalesmanCode: "MS123",
		OutletCode:   "OUT001",
	}}

	latestVisitCoordinateMap := map[string]model.LatestVisitCoordinateRow{
		buildDistributorVisitCoordinateKey("C22001", "MS123", "OUT001"): {
			CustID:          "C22001",
			EmpCode:         "MS123",
			OutletCode:      "OUT001",
			ArriveLongitude: 106.81,
			ArriveLatitude:  -6.21,
			FileURL:         ptrString(expectedFileURL),
		},
	}

	enrichDistributorRowsWithLatestVisits(rows, latestVisitCoordinateMap)

	if rows[0].FileURL == nil || *rows[0].FileURL != expectedFileURL {
		t.Fatalf("file_url = %#v, want %q", rows[0].FileURL, expectedFileURL)
	}
	if rows[0].ArriveLongitude != 106.81 || rows[0].ArriveLatitude != -6.21 {
		t.Fatalf("arrive coordinate = (%v,%v), want (106.81,-6.21)", rows[0].ArriveLongitude, rows[0].ArriveLatitude)
	}
}

func TestTransformDistributorRows_AssignsLeaveLocation(t *testing.T) {
	expectedLeaveLong := "106.81234"
	expectedLeaveLat := "-6.21234"
	rows := []model.LiveMonitoringDistributorRow{{
		CustID:          "C22001",
		EmpID:           301,
		SalesmanCode:    "MS123",
		PjpID:           8001,
		ApprovalStatus:  "Approved",
		RouteCode:       9001,
		OutletID:        7001,
		OutletCode:      "OUT001",
		OutletName:      "Outlet 1",
		DestinationType: "Outlet",
		Longitude:       106.81,
		Latitude:        -6.21,
		LeaveLongitude:  &expectedLeaveLong,
		LeaveLatitude:   &expectedLeaveLat,
	}}

	result := transformDistributorRows(rows)
	if len(result) != 1 {
		t.Fatalf("result length = %d, want 1", len(result))
	}
	dest := result[0].PjpData[0].RouteData[0].DestinationData[0]
	if dest.LeaveLongitude == nil || *dest.LeaveLongitude != expectedLeaveLong {
		t.Fatalf("leave_longitude = %#v, want %q", dest.LeaveLongitude, expectedLeaveLong)
	}
	if dest.LeaveLatitude == nil || *dest.LeaveLatitude != expectedLeaveLat {
		t.Fatalf("leave_latitude = %#v, want %q", dest.LeaveLatitude, expectedLeaveLat)
	}
}

func TestTransformDistributorRows_NilLeaveLocation(t *testing.T) {
	rows := []model.LiveMonitoringDistributorRow{{
		CustID:          "C22001",
		EmpID:           301,
		SalesmanCode:    "MS123",
		PjpID:           8001,
		ApprovalStatus:  "Approved",
		RouteCode:       9001,
		OutletID:        7001,
		OutletCode:      "OUT001",
		OutletName:      "Outlet 1",
		DestinationType: "Outlet",
		Longitude:       106.81,
		Latitude:        -6.21,
	}}

	result := transformDistributorRows(rows)
	dest := result[0].PjpData[0].RouteData[0].DestinationData[0]
	if dest.LeaveLongitude != nil {
		t.Fatalf("leave_longitude = %#v, want nil", dest.LeaveLongitude)
	}
	if dest.LeaveLatitude != nil {
		t.Fatalf("leave_latitude = %#v, want nil", dest.LeaveLatitude)
	}
}

func ptrInt64(value int64) *int64 {
	return &value
}
