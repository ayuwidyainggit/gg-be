package live_monitoring

import (
	"context"
	"scyllax-pjp/data/request"
	"scyllax-pjp/data/response"
	"scyllax-pjp/model"
	"testing"

	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type principalRepoStub struct {
	childCustIDs               []string
	childCustIDsErr            error
	principalEmployeeIDs       []int
	principalEmployeeIDsErr    error
	principalMonitoringRows    []model.LiveMonitoringPrincipalRow
	principalMonitoringErr     error
	principalExtraCallRows     []model.LiveMonitoringPrincipalRow
	principalExtraCallErr      error
	attendanceMap              map[int]model.AttendanceRow
	attendanceErr              error
	currentCoordinateMap       map[int]model.CurrentCoordinateRow
	currentCoordinateErr       error
	receivedChildCustID        string
	receivedScopeCustIDs       []string
	receivedScopedEmpIDs       []int
	receivedMonitoringCustIDs  []string
	receivedMonitoringEmpIDs   []int
	receivedMonitoringLimit    int
	receivedMonitoringOffset   int
	receivedMonitoringStatuses []string
	receivedMonitoringDate     string
	receivedMonitoringRegionID int
	receivedMonitoringAreaID   int
	receivedMonitoringDistID   int
	receivedExtraCallCustIDs   []string
	receivedExtraCallEmpIDs    []int
	receivedExtraCallStatuses  []string
	receivedExtraCallDate      string
	receivedExtraCallRegionID  int
	receivedExtraCallAreaID    int
	receivedExtraCallDistID    int
}

func (r *principalRepoStub) GetPrincipalEmployeeIDs(_ context.Context, _ *gorm.DB, custIDs []string, date string, regionID, areaID, distributorID int, empIDs []int, statuses []string) ([]int, error) {
	r.receivedScopeCustIDs = append([]string(nil), custIDs...)
	r.receivedScopedEmpIDs = append([]int(nil), empIDs...)
	r.receivedMonitoringStatuses = append([]string(nil), statuses...)
	r.receivedMonitoringDate = date
	r.receivedMonitoringRegionID = regionID
	r.receivedMonitoringAreaID = areaID
	r.receivedMonitoringDistID = distributorID
	return r.principalEmployeeIDs, r.principalEmployeeIDsErr
}

func (r *principalRepoStub) GetPrincipalMonitoring(_ context.Context, _ *gorm.DB, custIDs []string, date string, regionID, areaID, distributorID int, empIDs []int, statuses []string, limit, offset int) ([]model.LiveMonitoringPrincipalRow, error) {
	r.receivedMonitoringCustIDs = append([]string(nil), custIDs...)
	r.receivedMonitoringEmpIDs = append([]int(nil), empIDs...)
	r.receivedMonitoringStatuses = append([]string(nil), statuses...)
	r.receivedMonitoringDate = date
	r.receivedMonitoringRegionID = regionID
	r.receivedMonitoringAreaID = areaID
	r.receivedMonitoringDistID = distributorID
	r.receivedMonitoringLimit = limit
	r.receivedMonitoringOffset = offset
	return r.principalMonitoringRows, r.principalMonitoringErr
}

func (r *principalRepoStub) GetPrincipalExtraCalls(_ context.Context, _ *gorm.DB, custIDs []string, date string, regionID, areaID, distributorID int, empIDs []int, statuses []string) ([]model.LiveMonitoringPrincipalRow, error) {
	r.receivedExtraCallCustIDs = append([]string(nil), custIDs...)
	r.receivedExtraCallEmpIDs = append([]int(nil), empIDs...)
	r.receivedExtraCallStatuses = append([]string(nil), statuses...)
	r.receivedExtraCallDate = date
	r.receivedExtraCallRegionID = regionID
	r.receivedExtraCallAreaID = areaID
	r.receivedExtraCallDistID = distributorID
	return r.principalExtraCallRows, r.principalExtraCallErr
}

func (r *principalRepoStub) CountPrincipalMonitoring(context.Context, *gorm.DB, []string, string, int, int, int, []int, []string) (int64, error) {
	return 0, nil
}

func (r *principalRepoStub) GetDistributorEmployeeIDs(context.Context, *gorm.DB, []string, string, int, int, int, []int, []string) ([]int, error) {
	return nil, nil
}

func (r *principalRepoStub) GetDistributorMonitoring(context.Context, *gorm.DB, []string, string, int, int, int, []int, []string, int, int) ([]model.LiveMonitoringDistributorRow, error) {
	return nil, nil
}

func (r *principalRepoStub) GetDistributorLatestVisitCoordinates(context.Context, *gorm.DB, []string, string, []int) (map[string]model.LatestVisitCoordinateRow, error) {
	return nil, nil
}

func (r *principalRepoStub) GetDistributorEmployeeMeta(context.Context, *gorm.DB, []string, []int) (map[int]model.DistributorEmployeeMetaRow, error) {
	return nil, nil
}

func (r *principalRepoStub) GetDistributorRouteMeta(context.Context, *gorm.DB, []string, []int64) (map[string]model.DistributorRouteMetaRow, error) {
	return nil, nil
}

func (r *principalRepoStub) GetDistributorOutletMeta(context.Context, *gorm.DB, []string, []int) (map[string]model.DistributorOutletMetaRow, error) {
	return nil, nil
}

func (r *principalRepoStub) GetDistributorAttendance(context.Context, *gorm.DB, []string, string, int, int, int, []int, []string) (map[int]model.AttendanceRow, error) {
	if r.attendanceMap == nil {
		return map[int]model.AttendanceRow{}, r.attendanceErr
	}
	return r.attendanceMap, r.attendanceErr
}

func (r *principalRepoStub) GetDistributorCurrentCoordinates(context.Context, *gorm.DB, []string, string, int, int, int, []int, []string) (map[int]model.CurrentCoordinateRow, error) {
	if r.currentCoordinateMap == nil {
		return map[int]model.CurrentCoordinateRow{}, r.currentCoordinateErr
	}
	return r.currentCoordinateMap, r.currentCoordinateErr
}

func (r *principalRepoStub) CountDistributorMonitoring(context.Context, *gorm.DB, []string, string, int, int, int, []int, []string) (int64, error) {
	return 0, nil
}

func (r *principalRepoStub) GetVisitInformationPrincipal(context.Context, *gorm.DB, []string, string, int) (*model.VisitInformationRow, error) {
	return nil, nil
}

func (r *principalRepoStub) GetVisitInformationPrincipalFromHistory(context.Context, *gorm.DB, []string, string, int) (*model.VisitInformationRow, error) {
	return nil, nil
}

func (r *principalRepoStub) CountTotalVisitsPrincipal(context.Context, *gorm.DB, []string, string, int) (int64, error) {
	return 0, nil
}

func (r *principalRepoStub) GetVisitInformationDistributor(context.Context, *gorm.DB, string, int, int) (*model.VisitInformationRow, error) {
	return nil, nil
}

func (r *principalRepoStub) CountDistributorPlannedVisits(context.Context, *gorm.DB, string, int, int) (int64, error) {
	return 0, nil
}

func (r *principalRepoStub) CountDistributorExtraCalls(context.Context, *gorm.DB, string, int, int) (int64, error) {
	return 0, nil
}

func (r *principalRepoStub) CountDistributorOnGoingVisits(context.Context, *gorm.DB, string, int, int) (int64, error) {
	return 0, nil
}

func (r *principalRepoStub) CountDistributorVisitedVisits(context.Context, *gorm.DB, string, int, int) (int64, error) {
	return 0, nil
}

func (r *principalRepoStub) CountDistributorSkippedVisits(context.Context, *gorm.DB, string, int, int) (int64, error) {
	return 0, nil
}

func (r *principalRepoStub) GetSales(context.Context, *gorm.DB, []string, string, int) ([]model.SalesRow, error) {
	return nil, nil
}

func (r *principalRepoStub) GetReturns(context.Context, *gorm.DB, []string, string, int) ([]model.ReturnRow, error) {
	return nil, nil
}

func (r *principalRepoStub) GetCollections(context.Context, *gorm.DB, []string, string, int) ([]model.CollectionRow, error) {
	return nil, nil
}

func (r *principalRepoStub) GetExpenses(context.Context, *gorm.DB, string, int, string) ([]model.ExpenseRow, error) {
	return nil, nil
}

func (r *principalRepoStub) GetShipments(context.Context, *gorm.DB, []string, string, int) ([]model.ShipmentRow, error) {
	return nil, nil
}
func (r *principalRepoStub) GetSubmittedSurveyData(context.Context, *gorm.DB, []string, string, int) ([]model.SurveyDataRow, error) {
	return nil, nil
}
func (r *principalRepoStub) GetActivityTime(context.Context, *gorm.DB, string, int) (*string, error) {

	return nil, nil
}

func (r *principalRepoStub) GetDistributorInfo(context.Context, *gorm.DB, int) (*model.DistributorInfoRow, error) {
	return nil, nil
}

func (r *principalRepoStub) GetUserFullname(context.Context, *gorm.DB, string) (*string, error) {
	return nil, nil
}

func (r *principalRepoStub) GetChildCustIDs(context.Context, *gorm.DB, string) ([]string, error) {
	return r.childCustIDs, r.childCustIDsErr
}

func (r *principalRepoStub) GetSalesmanCustID(context.Context, *gorm.DB, int) (string, error) {
	return "", nil
}

func (r *principalRepoStub) GetEmployeeRole(context.Context, *gorm.DB, int, string) (string, error) {
	return "", nil
}

func (r *principalRepoStub) GetUpdateLocations(context.Context, *gorm.DB, int, string, string, string) ([]model.UpdateLocationRow, error) {
	return nil, nil
}

func TestTransformPrincipalRows_DoesNotDuplicateCrossJoinedDestinationsAfterRepoFix(t *testing.T) {
	arriveBMI260003 := int64(1779236763189)
	leaveBMI260003 := int64(1779237085474)
	arriveBMI260004 := int64(1779268286992)
	leaveBMI260004 := int64(1779268296910)

	rows := []model.LiveMonitoringPrincipalRow{
		{
			EmpID:           482,
			EmpCode:         "MS9990",
			EmpName:         "Sales 482",
			DistributorID:   22,
			AreaID:          7,
			RegionID:        3,
			PjpID:           101,
			PjpCode:         5001,
			ApprovalStatus:  "Approved",
			RouteCode:       7010,
			RouteName:       "Route 7010",
			DestinationID:   1,
			DestinationCode: "BMI260005",
			DestinationType: "OUTLET",
			DestinationName: "BMI260005",
			Longitude:       1,
			Latitude:        2,
		},
		{
			EmpID:           482,
			EmpCode:         "MS9990",
			EmpName:         "Sales 482",
			DistributorID:   22,
			AreaID:          7,
			RegionID:        3,
			PjpID:           101,
			PjpCode:         5001,
			ApprovalStatus:  "Approved",
			RouteCode:       7010,
			RouteName:       "Route 7010",
			DestinationID:   2,
			DestinationCode: "162612",
			DestinationType: "OUTLET",
			DestinationName: "162612",
			Longitude:       1,
			Latitude:        2,
		},
		{
			EmpID:           482,
			EmpCode:         "MS9990",
			EmpName:         "Sales 482",
			DistributorID:   22,
			AreaID:          7,
			RegionID:        3,
			PjpID:           101,
			PjpCode:         5001,
			ApprovalStatus:  "Approved",
			RouteCode:       7010,
			RouteName:       "Route 7010",
			DestinationID:   3,
			DestinationCode: "BMI260015",
			DestinationType: "OUTLET",
			DestinationName: "BMI260015",
			Longitude:       1,
			Latitude:        2,
		},
		{
			EmpID:           482,
			EmpCode:         "MS9990",
			EmpName:         "Sales 482",
			DistributorID:   22,
			AreaID:          7,
			RegionID:        3,
			PjpID:           101,
			PjpCode:         5001,
			ApprovalStatus:  "Approved",
			RouteCode:       7010,
			RouteName:       "Route 7010",
			DestinationID:   4,
			DestinationCode: "BMI260003",
			DestinationType: "OUTLET",
			DestinationName: "BMI260003",
			Longitude:       1,
			Latitude:        2,
			ArriveAt:        &arriveBMI260003,
			LeaveAt:         &leaveBMI260003,
			ArriveLongitude: -122.084000,
			ArriveLatitude:  37.421998,
		},
		{
			EmpID:           482,
			EmpCode:         "MS9990",
			EmpName:         "Sales 482",
			DistributorID:   22,
			AreaID:          7,
			RegionID:        3,
			PjpID:           101,
			PjpCode:         5001,
			ApprovalStatus:  "Approved",
			RouteCode:       7010,
			RouteName:       "Route 7010",
			DestinationID:   5,
			DestinationCode: "BMI260004",
			DestinationType: "OUTLET",
			DestinationName: "BMI260004",
			Longitude:       1,
			Latitude:        2,
			ArriveAt:        &arriveBMI260004,
			LeaveAt:         &leaveBMI260004,
		},
	}

	result := transformPrincipalRows(rows)
	if len(result) != 1 {
		t.Fatalf("employees = %d, want 1", len(result))
	}
	if len(result[0].PjpData) != 1 {
		t.Fatalf("pjp count = %d, want 1", len(result[0].PjpData))
	}
	if len(result[0].PjpData[0].RouteData) != 1 {
		t.Fatalf("route count = %d, want 1", len(result[0].PjpData[0].RouteData))
	}

	destinations := result[0].PjpData[0].RouteData[0].DestinationData
	if len(destinations) != 5 {
		t.Fatalf("destination count = %d, want 5", len(destinations))
	}

	seen := make(map[string]bool)
	for _, destination := range destinations {
		if seen[destination.DestinationCode] {
			t.Fatalf("duplicate destination_code found: %s", destination.DestinationCode)
		}
		seen[destination.DestinationCode] = true
	}

	bmi260003, ok := findPrincipalDestination(destinations, "BMI260003")
	if !ok {
		t.Fatalf("BMI260003 not found in destination_data: %#v", destinations)
	}
	if bmi260003.ArriveAt == nil || *bmi260003.ArriveAt != arriveBMI260003 {
		t.Fatalf("BMI260003 arrive_at = %#v, want %d", bmi260003.ArriveAt, arriveBMI260003)
	}
	if bmi260003.LeaveAt == nil || *bmi260003.LeaveAt != leaveBMI260003 {
		t.Fatalf("BMI260003 leave_at = %#v, want %d", bmi260003.LeaveAt, leaveBMI260003)
	}
	if bmi260003.ArriveLongitude != -122.084000 {
		t.Fatalf("BMI260003 arrive_longitude = %v, want -122.084000", bmi260003.ArriveLongitude)
	}
	if bmi260003.ArriveLatitude != 37.421998 {
		t.Fatalf("BMI260003 arrive_latitude = %v, want 37.421998", bmi260003.ArriveLatitude)
	}

	bmi260004, ok := findPrincipalDestination(destinations, "BMI260004")
	if !ok {
		t.Fatalf("BMI260004 not found in destination_data: %#v", destinations)
	}
	if bmi260004.ArriveAt == nil || *bmi260004.ArriveAt != arriveBMI260004 {
		t.Fatalf("BMI260004 arrive_at = %#v, want %d", bmi260004.ArriveAt, arriveBMI260004)
	}
}

func TestGetPrincipalMonitoring_PaginatesEmployeeScopeBeforeDetailRows(t *testing.T) {
	repo := &principalRepoStub{
		childCustIDs:         []string{"C26002", "C260020001"},
		principalEmployeeIDs: []int{481, 482, 483},
		principalMonitoringRows: []model.LiveMonitoringPrincipalRow{
			{EmpID: 482, EmpCode: "MS9990", EmpName: "Sales 482", DistributorID: 22, AreaID: 7, RegionID: 3, PjpID: 101, PjpCode: 5001, ApprovalStatus: "Approved", RouteCode: 7010, RouteName: "Route 7010", DestinationID: 1, DestinationCode: "BMI260005", DestinationType: "OUTLET", DestinationName: "BMI260005"},
			{EmpID: 482, EmpCode: "MS9990", EmpName: "Sales 482", DistributorID: 22, AreaID: 7, RegionID: 3, PjpID: 101, PjpCode: 5001, ApprovalStatus: "Approved", RouteCode: 7010, RouteName: "Route 7010", DestinationID: 2, DestinationCode: "162612", DestinationType: "OUTLET", DestinationName: "162612"},
			{EmpID: 482, EmpCode: "MS9990", EmpName: "Sales 482", DistributorID: 22, AreaID: 7, RegionID: 3, PjpID: 101, PjpCode: 5001, ApprovalStatus: "Approved", RouteCode: 7010, RouteName: "Route 7010", DestinationID: 3, DestinationCode: "BMI260015", DestinationType: "OUTLET", DestinationName: "BMI260015"},
			{EmpID: 482, EmpCode: "MS9990", EmpName: "Sales 482", DistributorID: 22, AreaID: 7, RegionID: 3, PjpID: 101, PjpCode: 5001, ApprovalStatus: "Approved", RouteCode: 7010, RouteName: "Route 7010", DestinationID: 4, DestinationCode: "BMI260003", DestinationType: "OUTLET", DestinationName: "BMI260003"},
			{EmpID: 482, EmpCode: "MS9990", EmpName: "Sales 482", DistributorID: 22, AreaID: 7, RegionID: 3, PjpID: 101, PjpCode: 5001, ApprovalStatus: "Approved", RouteCode: 7010, RouteName: "Route 7010", DestinationID: 5, DestinationCode: "BMI260004", DestinationType: "OUTLET", DestinationName: "BMI260004"},
		},
	}

	svc := &liveMonitoringService{repository: repo, validate: validator.New()}

	result, paging, err := svc.GetPrincipalMonitoring(context.Background(), request.LiveMonitoringRequest{
		Date:        1779278400,
		Status:      []string{"Approved", "Need Review"},
		Page:        2,
		Limit:       1,
		LegacyEmpID: "482",
	}, "C26002")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(repo.receivedScopeCustIDs) != 2 || repo.receivedScopeCustIDs[0] != "C26002" || repo.receivedScopeCustIDs[1] != "C260020001" {
		t.Fatalf("scope custIDs = %#v, want [C26002 C260020001]", repo.receivedScopeCustIDs)
	}

	if len(repo.receivedScopedEmpIDs) != 1 || repo.receivedScopedEmpIDs[0] != 482 {
		t.Fatalf("requested scoped empIDs = %#v, want [482]", repo.receivedScopedEmpIDs)
	}

	if len(repo.receivedMonitoringEmpIDs) != 1 || repo.receivedMonitoringEmpIDs[0] != 482 {
		t.Fatalf("monitoring empIDs = %#v, want [482]", repo.receivedMonitoringEmpIDs)
	}

	if repo.receivedMonitoringLimit != 0 || repo.receivedMonitoringOffset != 0 {
		t.Fatalf("monitoring limit/offset = (%d,%d), want (0,0)", repo.receivedMonitoringLimit, repo.receivedMonitoringOffset)
	}

	if paging.PageTotal != 3 || paging.PageCurrent != 2 || paging.PageLimit != 1 {
		t.Fatalf("paging = %#v, want total=3 current=2 limit=1", paging)
	}

	if len(result) != 1 {
		t.Fatalf("employees = %d, want 1", len(result))
	}

	destinations := result[0].PjpData[0].RouteData[0].DestinationData
	if len(destinations) != 5 {
		t.Fatalf("destination count = %d, want 5", len(destinations))
	}
}

func TestGetPrincipalMonitoring_FallsBackToCurrentCustIDWhenNoChildrenFound(t *testing.T) {
	repo := &principalRepoStub{
		childCustIDs:         []string{},
		principalEmployeeIDs: []int{482},
	}

	svc := &liveMonitoringService{repository: repo, validate: validator.New()}

	_, _, err := svc.GetPrincipalMonitoring(context.Background(), request.LiveMonitoringRequest{
		Date:   1779278400,
		Status: []string{"Approved"},
	}, "C26002")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(repo.receivedMonitoringCustIDs) != 1 || repo.receivedMonitoringCustIDs[0] != "C26002" {
		t.Fatalf("monitoring custIDs = %#v, want [C26002]", repo.receivedMonitoringCustIDs)
	}
}

func TestGetPrincipalMonitoring_PopulatesTopLevelAttendanceAndCurrentFields(t *testing.T) {
	attendanceAt := int64(1747720800)
	checkoutAt := int64(1747724400)

	repo := &principalRepoStub{
		childCustIDs:         []string{"C26002"},
		principalEmployeeIDs: []int{482},
		principalMonitoringRows: []model.LiveMonitoringPrincipalRow{
			{
				EmpID:           482,
				EmpCode:         "MS9990",
				EmpName:         "Sales 482",
				DistributorID:   22,
				AreaID:          7,
				RegionID:        3,
				PjpID:           101,
				PjpCode:         5001,
				ApprovalStatus:  "Approved",
				RouteCode:       7010,
				RouteName:       "Route 7010",
				DestinationID:   4,
				DestinationCode: "BMI260003",
				DestinationType: "OUTLET",
				DestinationName: "BMI260003",
			},
		},
		attendanceMap: map[int]model.AttendanceRow{
			482: {
				AttendanceID: ptrInt64(991),
				Timestamp:    &attendanceAt,
				Longitude:    106.8123,
				Latitude:     -6.2012,
				ClockOutID:   ptrInt64(992),
				ClockOutAt:   &checkoutAt,
				ClockOutLong: 106.8456,
				ClockOutLat:  -6.2456,
			},
		},
		currentCoordinateMap: map[int]model.CurrentCoordinateRow{
			482: {
				Longitude: 106.8456,
				Latitude:  -6.2456,
				Timestamp: &checkoutAt,
				Source:    "attendance_checkout",
			},
		},
	}

	svc := &liveMonitoringService{repository: repo, validate: validator.New()}

	result, _, err := svc.GetPrincipalMonitoring(context.Background(), request.LiveMonitoringRequest{
		Date:   1747674000,
		Status: []string{"Approved"},
	}, "C26002")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result) != 1 {
		t.Fatalf("employees = %d, want 1", len(result))
	}

	if result[0].AttendanceID == nil || *result[0].AttendanceID != 991 {
		t.Fatalf("attendance id = %#v, want 991", result[0].AttendanceID)
	}
	if result[0].AttendanceAt == nil || *result[0].AttendanceAt != attendanceAt {
		t.Fatalf("attendance at = %#v, want %d", result[0].AttendanceAt, attendanceAt)
	}
	if result[0].AttendanceLongitude != 106.8123 || result[0].AttendanceLatitude != -6.2012 {
		t.Fatalf("attendance coordinate = (%v,%v), want (106.8123,-6.2012)", result[0].AttendanceLongitude, result[0].AttendanceLatitude)
	}
	if result[0].ClockOut == nil || *result[0].ClockOut != 992 {
		t.Fatalf("clock out = %#v, want 992", result[0].ClockOut)
	}
	if result[0].ClockOutAt == nil || *result[0].ClockOutAt != checkoutAt {
		t.Fatalf("clock out at = %#v, want %d", result[0].ClockOutAt, checkoutAt)
	}
	if result[0].ClockOutLongitude != 106.8456 || result[0].ClockOutLatitude != -6.2456 {
		t.Fatalf("clock out coordinate = (%v,%v), want (106.8456,-6.2456)", result[0].ClockOutLongitude, result[0].ClockOutLatitude)
	}
	if result[0].CurrentCoordinateAt == nil || *result[0].CurrentCoordinateAt != checkoutAt {
		t.Fatalf("current coordinate at = %#v, want %d", result[0].CurrentCoordinateAt, checkoutAt)
	}
	if result[0].CurrentLongitude != 106.8456 || result[0].CurrentLatitude != -6.2456 {
		t.Fatalf("current coordinate = (%v,%v), want (106.8456,-6.2456)", result[0].CurrentLongitude, result[0].CurrentLatitude)
	}
	if result[0].CurrentCoordinateSource != "attendance_checkout" {
		t.Fatalf("current coordinate source = %q, want attendance_checkout", result[0].CurrentCoordinateSource)
	}
	if len(result[0].PjpData) != 1 || len(result[0].PjpData[0].RouteData) != 1 || len(result[0].PjpData[0].RouteData[0].DestinationData) != 1 {
		t.Fatalf("destination hierarchy regressed: %#v", result[0].PjpData)
	}
}

func TestGetPrincipalMonitoring_ResetsTopLevelAttendanceFieldsWhenNoAttendanceOnRequestedDate(t *testing.T) {
	staleAttendanceAt := int64(1747634400)
	staleCurrentAt := int64(1747638000)

	repo := &principalRepoStub{
		childCustIDs:         []string{"C26002"},
		principalEmployeeIDs: []int{482},
		principalMonitoringRows: []model.LiveMonitoringPrincipalRow{
			{
				EmpID:           482,
				EmpCode:         "MS9990",
				EmpName:         "Sales 482",
				DistributorID:   22,
				AreaID:          7,
				RegionID:        3,
				PjpID:           101,
				PjpCode:         5001,
				ApprovalStatus:  "Approved",
				RouteCode:       7010,
				RouteName:       "Route 7010",
				DestinationID:   4,
				DestinationCode: "BMI260003",
				DestinationType: "OUTLET",
				DestinationName: "BMI260003",
				ArriveAt:        &staleAttendanceAt,
				ArriveLongitude: 106.8,
				ArriveLatitude:  -6.2,
			},
		},
		attendanceMap: map[int]model.AttendanceRow{
			482: {
				AttendanceID: ptrInt64(991),
				Timestamp:    &staleAttendanceAt,
				Longitude:    106.8123,
				Latitude:     -6.2012,
			},
		},
		currentCoordinateMap: map[int]model.CurrentCoordinateRow{
			482: {
				Longitude: 106.8456,
				Latitude:  -6.2456,
				Timestamp: &staleCurrentAt,
				Source:    "attendance_checkout",
			},
		},
	}

	svc := &liveMonitoringService{repository: repo, validate: validator.New()}

	result, _, err := svc.GetPrincipalMonitoring(context.Background(), request.LiveMonitoringRequest{
		Date:   1747674000,
		Status: []string{"Approved"},
	}, "C26002")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result) != 1 {
		t.Fatalf("employees = %d, want 1", len(result))
	}

	if result[0].AttendanceID != nil || result[0].AttendanceAt != nil || result[0].CurrentCoordinateAt != nil {
		t.Fatalf("top-level daily tracking state not reset: %#v", result[0])
	}
	if len(result[0].PjpData) != 1 || len(result[0].PjpData[0].RouteData) != 1 || len(result[0].PjpData[0].RouteData[0].DestinationData) != 1 {
		t.Fatalf("destination hierarchy regressed: %#v", result[0].PjpData)
	}
	if result[0].PjpData[0].RouteData[0].DestinationData[0].ArriveAt == nil {
		t.Fatalf("destination-level arrive_at should remain intact for principal rows")
	}
}

func TestGetPrincipalMonitoring_PopulatesExtraCallData(t *testing.T) {
	repo := &principalRepoStub{
		childCustIDs:         []string{"C26002", "C260020001"},
		principalEmployeeIDs: []int{482},
		principalMonitoringRows: []model.LiveMonitoringPrincipalRow{
			{
				EmpID: 482, EmpCode: "MS9990", EmpName: "Sales 482",
				DistributorID: 22, AreaID: 7, RegionID: 3,
				PjpID: 101, PjpCode: 5001, ApprovalStatus: "Approved",
				RouteCode: 7010, RouteName: "Route 7010",
				DestinationID: 1, DestinationCode: "BMI260005", DestinationType: "outlet",
				DestinationName: "Toko Akbar",
				Longitude:       106.8, Latitude: -6.2,
				IsExtraCall: false,
			},
		},
		principalExtraCallRows: []model.LiveMonitoringPrincipalRow{
			{
				EmpID: 482, EmpCode: "MS9990", EmpName: "Sales 482",
				DistributorID: 22, AreaID: 7, RegionID: 3,
				PjpID: 101, PjpCode: 5001, ApprovalStatus: "Approved",
				RouteCode: 7010, RouteName: "Route 7010",
				DestinationID: 99, DestinationCode: "BMI260099", DestinationType: "outlet",
				DestinationName: "Toko Principal Extra",
				Longitude:       106.9, Latitude: -6.3,
				IsExtraCall: true,
			},
		},
	}

	svc := &liveMonitoringService{repository: repo, validate: validator.New()}

	result, _, err := svc.GetPrincipalMonitoring(context.Background(), request.LiveMonitoringRequest{
		Date:        1779364800,
		Status:      []string{"Approved", "Need Review"},
		LegacyEmpID: "482",
	}, "C26002")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result) != 1 {
		t.Fatalf("employees = %d, want 1", len(result))
	}
	if len(result[0].PjpData) != 1 {
		t.Fatalf("pjp count = %d, want 1", len(result[0].PjpData))
	}

	pjp := result[0].PjpData[0]

	// route_data must contain only the non-extra-call destination
	if len(pjp.RouteData) != 1 {
		t.Fatalf("route_data count = %d, want 1", len(pjp.RouteData))
	}
	if len(pjp.RouteData[0].DestinationData) != 1 {
		t.Fatalf("route_data[0].destination_data count = %d, want 1", len(pjp.RouteData[0].DestinationData))
	}
	if pjp.RouteData[0].DestinationData[0].DestinationCode != "BMI260005" {
		t.Fatalf("route_data destination_code = %q, want BMI260005", pjp.RouteData[0].DestinationData[0].DestinationCode)
	}

	// extra_call_data must contain the extra call destination
	if len(pjp.ExtraCallData) != 1 {
		t.Fatalf("extra_call_data count = %d, want 1", len(pjp.ExtraCallData))
	}
	if len(pjp.ExtraCallData[0].DestinationData) != 1 {
		t.Fatalf("extra_call_data[0].destination_data count = %d, want 1", len(pjp.ExtraCallData[0].DestinationData))
	}
	extra := pjp.ExtraCallData[0].DestinationData[0]
	if extra.DestinationCode != "BMI260099" {
		t.Fatalf("extra_call destination_code = %q, want BMI260099", extra.DestinationCode)
	}
	if extra.DestinationName != "Toko Principal Extra" {
		t.Fatalf("extra_call destination_name = %q, want Toko Principal Extra", extra.DestinationName)
	}
	if extra.Longitude != 106.9 || extra.Latitude != -6.3 {
		t.Fatalf("extra_call coordinate = (%v,%v), want (106.9,-6.3)", extra.Longitude, extra.Latitude)
	}
}

func findPrincipalDestination(destinations []response.LiveMonitoringDestinationData, code string) (response.LiveMonitoringDestinationData, bool) {
	for _, destination := range destinations {
		if destination.DestinationCode == code {
			return destination, true
		}
	}
	return response.LiveMonitoringDestinationData{}, false
}

func ptrString(value string) *string {
	return &value
}

func TestTransformPrincipalRows_PropagatesFileURL(t *testing.T) {
	expectedFileURL := "https://files.example.com/arrive/arrive-123.jpg"
	expectedLeaveLong := "106.81234"
	expectedLeaveLat := "-6.21234"
	extraLeaveLong := "107.10001"
	extraLeaveLat := "-6.10001"

	rows := []model.LiveMonitoringPrincipalRow{
		{
			EmpID:           482,
			EmpCode:         "MS9990",
			EmpName:         "Sales 482",
			DistributorID:   22,
			AreaID:          7,
			RegionID:        3,
			PjpID:           101,
			PjpCode:         5001,
			ApprovalStatus:  "Approved",
			RouteCode:       7010,
			RouteName:       "Route 7010",
			DestinationID:   1,
			DestinationCode: "BMI260001",
			DestinationType: "OUTLET",
			DestinationName: "BMI260001",
			Longitude:       1,
			Latitude:        2,
			LeaveLongitude:  ptrString(expectedLeaveLong),
			LeaveLatitude:   ptrString(expectedLeaveLat),
			FileURL:         ptrString(expectedFileURL),
		},
		{
			EmpID:           482,
			EmpCode:         "MS9990",
			EmpName:         "Sales 482",
			DistributorID:   22,
			AreaID:          7,
			RegionID:        3,
			PjpID:           101,
			PjpCode:         5001,
			ApprovalStatus:  "Approved",
			RouteCode:       7010,
			RouteName:       "Route 7010",
			DestinationID:   2,
			DestinationCode: "BMI260002",
			DestinationType: "OUTLET",
			DestinationName: "BMI260002",
			Longitude:       1,
			Latitude:        2,
			LeaveLongitude:  ptrString(extraLeaveLong),
			LeaveLatitude:   ptrString(extraLeaveLat),
			IsExtraCall:     true,
		},
	}

	result := transformPrincipalRows(rows)
	if len(result) != 1 {
		t.Fatalf("employees = %d, want 1", len(result))
	}
	if len(result[0].PjpData) != 1 {
		t.Fatalf("pjp count = %d, want 1", len(result[0].PjpData))
	}
	pjp := result[0].PjpData[0]

	if len(pjp.RouteData) != 1 {
		t.Fatalf("route_data count = %d, want 1", len(pjp.RouteData))
	}
	routeDest := pjp.RouteData[0].DestinationData
	withURL, ok := findPrincipalDestination(routeDest, "BMI260001")
	if !ok {
		t.Fatalf("BMI260001 not found in route_data destination_data: %#v", routeDest)
	}
	if withURL.FileURL == nil || *withURL.FileURL != expectedFileURL {
		t.Fatalf("BMI260001 file_url = %#v, want %q", withURL.FileURL, expectedFileURL)
	}

	if len(pjp.ExtraCallData) != 1 {
		t.Fatalf("extra_call_data count = %d, want 1", len(pjp.ExtraCallData))
	}
	extraDest := pjp.ExtraCallData[0].DestinationData
	withoutURL, ok := findPrincipalDestination(extraDest, "BMI260002")
	if !ok {
		t.Fatalf("BMI260002 not found in extra_call_data destination_data: %#v", extraDest)
	}
	if withoutURL.FileURL != nil {
		t.Fatalf("BMI260002 file_url = %#v, want nil", withoutURL.FileURL)
	}

	if withURL.LeaveLongitude == nil || *withURL.LeaveLongitude != expectedLeaveLong {
		t.Fatalf("BMI260001 leave_longitude = %#v, want %q", withURL.LeaveLongitude, expectedLeaveLong)
	}
	if withURL.LeaveLatitude == nil || *withURL.LeaveLatitude != expectedLeaveLat {
		t.Fatalf("BMI260001 leave_latitude = %#v, want %q", withURL.LeaveLatitude, expectedLeaveLat)
	}
	if withoutURL.LeaveLongitude == nil || *withoutURL.LeaveLongitude != extraLeaveLong {
		t.Fatalf("BMI260002 leave_longitude = %#v, want %q", withoutURL.LeaveLongitude, extraLeaveLong)
	}
	if withoutURL.LeaveLatitude == nil || *withoutURL.LeaveLatitude != extraLeaveLat {
		t.Fatalf("BMI260002 leave_latitude = %#v, want %q", withoutURL.LeaveLatitude, extraLeaveLat)
	}
}
