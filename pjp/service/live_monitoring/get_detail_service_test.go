package live_monitoring

import (
	"context"
	"testing"

	"scyllax-pjp/data/request"
	"scyllax-pjp/model"

	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type detailRepoStub struct {
	visitDistributor          *model.VisitInformationRow
	visitErr                  error
	visitPrincipalHistory     *model.VisitInformationRow
	activityTime              *string
	activityErr               error
	distInfo                  *model.DistributorInfoRow
	distErr                   error
	plannedCount              int64
	extraCallCount            int64
	onGoingCount              int64
	visitedCount              int64
	skippedCount              int64
	expenses                  []model.ExpenseRow
	expenseErr                error
	collections               []model.CollectionRow
	surveys                   []model.SurveyDataRow
	receivedCustID            string
	receivedEmpID             int
	receivedDate              string
	receivedPrincipalEmpID    int
	receivedCollectionCustIDs []string
	receivedCollectionDate    string
	receivedCollectionEmpID   int
}

func (r *detailRepoStub) GetPrincipalEmployeeIDs(context.Context, *gorm.DB, []string, string, int, int, int, []int, []string) ([]int, error) {
	return nil, nil
}
func (r *detailRepoStub) GetPrincipalMonitoring(context.Context, *gorm.DB, []string, string, int, int, int, []int, []string, int, int) ([]model.LiveMonitoringPrincipalRow, error) {
	return nil, nil
}
func (r *detailRepoStub) GetPrincipalExtraCalls(context.Context, *gorm.DB, []string, string, int, int, int, []int, []string) ([]model.LiveMonitoringPrincipalRow, error) {
	return nil, nil
}
func (r *detailRepoStub) CountPrincipalMonitoring(context.Context, *gorm.DB, []string, string, int, int, int, []int, []string) (int64, error) {
	return 0, nil
}
func (r *detailRepoStub) GetDistributorEmployeeIDs(context.Context, *gorm.DB, []string, string, int, int, int, []int, []string) ([]int, error) {
	return nil, nil
}
func (r *detailRepoStub) GetDistributorMonitoring(context.Context, *gorm.DB, []string, string, int, int, int, []int, []string, int, int) ([]model.LiveMonitoringDistributorRow, error) {
	return nil, nil
}
func (r *detailRepoStub) GetDistributorLatestVisitCoordinates(context.Context, *gorm.DB, []string, string, []int) (map[string]model.LatestVisitCoordinateRow, error) {
	return map[string]model.LatestVisitCoordinateRow{}, nil
}
func (r *detailRepoStub) GetDistributorEmployeeMeta(context.Context, *gorm.DB, []string, []int) (map[int]model.DistributorEmployeeMetaRow, error) {
	return map[int]model.DistributorEmployeeMetaRow{}, nil
}
func (r *detailRepoStub) GetDistributorRouteMeta(context.Context, *gorm.DB, []string, []int64) (map[string]model.DistributorRouteMetaRow, error) {
	return map[string]model.DistributorRouteMetaRow{}, nil
}
func (r *detailRepoStub) GetDistributorOutletMeta(context.Context, *gorm.DB, []string, []int) (map[string]model.DistributorOutletMetaRow, error) {
	return map[string]model.DistributorOutletMetaRow{}, nil
}
func (r *detailRepoStub) GetDistributorAttendance(context.Context, *gorm.DB, []string, string, int, int, int, []int, []string) (map[int]model.AttendanceRow, error) {
	return nil, nil
}
func (r *detailRepoStub) GetDistributorCurrentCoordinates(context.Context, *gorm.DB, []string, string, int, int, int, []int, []string) (map[int]model.CurrentCoordinateRow, error) {
	return nil, nil
}
func (r *detailRepoStub) CountDistributorMonitoring(context.Context, *gorm.DB, []string, string, int, int, int, []int, []string) (int64, error) {
	return 0, nil
}
func (r *detailRepoStub) GetVisitInformationPrincipal(context.Context, *gorm.DB, []string, string, int) (*model.VisitInformationRow, error) {
	return nil, nil
}
func (r *detailRepoStub) GetVisitInformationPrincipalFromHistory(_ context.Context, _ *gorm.DB, _ []string, _ string, empID int) (*model.VisitInformationRow, error) {
	r.receivedPrincipalEmpID = empID
	return r.visitPrincipalHistory, nil
}
func (r *detailRepoStub) CountTotalVisitsPrincipal(context.Context, *gorm.DB, []string, string, int) (int64, error) {
	return 0, nil
}
func (r *detailRepoStub) GetVisitInformationDistributor(context.Context, *gorm.DB, string, int, int) (*model.VisitInformationRow, error) {
	return r.visitDistributor, r.visitErr
}
func (r *detailRepoStub) CountDistributorPlannedVisits(context.Context, *gorm.DB, string, int, int) (int64, error) {
	return r.plannedCount, nil
}
func (r *detailRepoStub) CountDistributorExtraCalls(context.Context, *gorm.DB, string, int, int) (int64, error) {
	return r.extraCallCount, nil
}
func (r *detailRepoStub) CountDistributorOnGoingVisits(context.Context, *gorm.DB, string, int, int) (int64, error) {
	return r.onGoingCount, nil
}
func (r *detailRepoStub) CountDistributorVisitedVisits(context.Context, *gorm.DB, string, int, int) (int64, error) {
	return r.visitedCount, nil
}
func (r *detailRepoStub) CountDistributorSkippedVisits(context.Context, *gorm.DB, string, int, int) (int64, error) {
	return r.skippedCount, nil
}
func (r *detailRepoStub) GetSales(context.Context, *gorm.DB, []string, string, int) ([]model.SalesRow, error) {
	return nil, nil
}
func (r *detailRepoStub) GetReturns(context.Context, *gorm.DB, []string, string, int) ([]model.ReturnRow, error) {
	return nil, nil
}
func (r *detailRepoStub) GetCollections(_ context.Context, _ *gorm.DB, custIDs []string, date string, empID int) ([]model.CollectionRow, error) {
	r.receivedCollectionCustIDs = append([]string(nil), custIDs...)
	r.receivedCollectionDate = date
	r.receivedCollectionEmpID = empID
	return r.collections, nil
}
func (r *detailRepoStub) GetExpenses(_ context.Context, _ *gorm.DB, custID string, empID int, date string) ([]model.ExpenseRow, error) {
	r.receivedCustID = custID
	r.receivedEmpID = empID
	r.receivedDate = date
	return r.expenses, r.expenseErr
}
func (r *detailRepoStub) GetShipments(context.Context, *gorm.DB, []string, string, int) ([]model.ShipmentRow, error) {
	return nil, nil
}
func (r *detailRepoStub) GetSubmittedSurveyData(context.Context, *gorm.DB, []string, string, int) ([]model.SurveyDataRow, error) {
	return r.surveys, nil
}
func (r *detailRepoStub) GetActivityTime(context.Context, *gorm.DB, string, int) (*string, error) {
	return r.activityTime, r.activityErr
}
func (r *detailRepoStub) GetDistributorInfo(context.Context, *gorm.DB, int) (*model.DistributorInfoRow, error) {
	return r.distInfo, r.distErr
}
func (r *detailRepoStub) GetUserFullname(context.Context, *gorm.DB, string) (*string, error) {
	return nil, nil
}
func (r *detailRepoStub) GetChildCustIDs(context.Context, *gorm.DB, string) ([]string, error) {
	return nil, nil
}
func (r *detailRepoStub) GetSalesmanCustID(context.Context, *gorm.DB, int) (string, error) {
	return "C220010001", nil
}

func (r *detailRepoStub) GetEmployeeRole(context.Context, *gorm.DB, int, string) (string, error) {
	return "", nil
}

func (r *detailRepoStub) GetUpdateLocations(context.Context, *gorm.DB, int, string, string, string) ([]model.UpdateLocationRow, error) {
	return nil, nil
}

func TestGetDistributorVisitInfo_UsesRepositoryExtraCallAggregation(t *testing.T) {
	activityTime := "2026-03-18T07:56:09.52461Z"
	repo := &detailRepoStub{
		visitDistributor: &model.VisitInformationRow{
			EmpID:   360,
			EmpCode: "2025120204",
			EmpName: "Yogie Setya",
		},
		activityTime:   &activityTime,
		plannedCount:   2,
		extraCallCount: 2,
		onGoingCount:   1,
		visitedCount:   2,
		skippedCount:   0,
		distInfo: &model.DistributorInfoRow{
			DistributorID:   67,
			DistributorCode: "3434",
			DistributorName: "Distributor iDetama",
		},
	}

	svc := &liveMonitoringService{repository: repo, validate: validator.New()}

	visitInfo, err := svc.getDistributorVisitInfo(context.Background(), "2026-03-18", 360, 67, "C220010001")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if visitInfo == nil {
		t.Fatal("visitInfo is nil")
	}

	if visitInfo.Planned != 2 {
		t.Fatalf("planned = %d, want 2", visitInfo.Planned)
	}
	if visitInfo.ExtraCall != 2 {
		t.Fatalf("extra_call = %d, want 2", visitInfo.ExtraCall)
	}
	if visitInfo.OnGoing != 1 {
		t.Fatalf("on_going = %d, want 1", visitInfo.OnGoing)
	}
	if visitInfo.Visited != 2 {
		t.Fatalf("visited = %d, want 2", visitInfo.Visited)
	}
	if visitInfo.Skipped != 0 {
		t.Fatalf("skipped = %d, want 0", visitInfo.Skipped)
	}
	if visitInfo.ReturnSummary.Count != 0 || visitInfo.ReturnSummary.Status != "none" {
		t.Fatalf("return_summary = %+v, want count=0 status=none", visitInfo.ReturnSummary)
	}
	if visitInfo.CollectionSummary.Count != 0 || visitInfo.CollectionSummary.Status != "none" {
		t.Fatalf("collection_summary = %+v, want count=0 status=none", visitInfo.CollectionSummary)
	}
}

func TestBuildVisitSummary_DefaultAndCompleted(t *testing.T) {
	defaultSummary := buildVisitSummary(0)
	if defaultSummary.Count != 0 || defaultSummary.Status != "none" {
		t.Fatalf("default summary = %+v, want count=0 status=none", defaultSummary)
	}

	completedSummary := buildVisitSummary(3)
	if completedSummary.Count != 3 || completedSummary.Status != "completed" {
		t.Fatalf("completed summary = %+v, want count=3 status=completed", completedSummary)
	}
}

func TestGetMonitoringDetail_FiltersExpenseByCollector(t *testing.T) {
	activityTime := "2026-03-18T07:56:09.52461Z"
	repo := &detailRepoStub{
		visitDistributor: &model.VisitInformationRow{
			EmpID:   358,
			EmpCode: "EMP-358",
			EmpName: "Firman Mulyawan",
		},
		activityTime: &activityTime,
		distInfo: &model.DistributorInfoRow{
			DistributorID:   67,
			DistributorCode: "3434",
			DistributorName: "Distributor iDetama",
		},
		expenses: []model.ExpenseRow{
			{ExpenseTypeID: 10, ExpenseTypeName: "Transport", Note: "Collector A", Amount: 15000},
			{ExpenseTypeID: 11, ExpenseTypeName: "Meal", Note: "Collector A second", Amount: 20000},
		},
	}

	svc := &liveMonitoringService{repository: repo, validate: validator.New()}
	distributorID := 67
	result, err := svc.GetMonitoringDetail(context.Background(), request.LiveMonitoringDetailRequest{
		EmpID:         358,
		DistributorID: &distributorID,
		Date:          "2026-03-18",
	}, "C220010001", 999)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result == nil {
		t.Fatal("result is nil")
	}

	if repo.receivedCustID != "C220010001" {
		t.Fatalf("expense cust_id = %s, want C220010001", repo.receivedCustID)
	}
	if repo.receivedEmpID != 358 {
		t.Fatalf("expense emp_id = %d, want 358", repo.receivedEmpID)
	}
	if repo.receivedDate != "2026-03-18" {
		t.Fatalf("expense date = %s, want 2026-03-18", repo.receivedDate)
	}
	if len(result.Expense) != 2 {
		t.Fatalf("len(expense) = %d, want 2", len(result.Expense))
	}
	if result.VisitInformation.EmpID != 358 {
		t.Fatalf("visit_information.emp_id = %d, want 358", result.VisitInformation.EmpID)
	}
	if result.VisitInformation.CompanyCode != "3434" {
		t.Fatalf("company_code = %s, want 3434", result.VisitInformation.CompanyCode)
	}
	if len(result.Sales) != 0 || len(result.Return) != 0 || len(result.Shipment) != 0 || len(result.Collection) != 0 {
		t.Fatalf("non-expense sections should remain empty, got sales=%d return=%d shipment=%d collection=%d", len(result.Sales), len(result.Return), len(result.Shipment), len(result.Collection))
	}
	if result.Expense[0].ExpenseTypeID != 10 || result.Expense[1].ExpenseTypeID != 11 {
		t.Fatalf("unexpected expense payload: %+v", result.Expense)
	}
}

func TestGetMonitoringDetail_DifferentCollectorsStayIsolated(t *testing.T) {
	repoCollectorA := &detailRepoStub{
		visitDistributor: &model.VisitInformationRow{EmpID: 358, EmpCode: "EMP-358", EmpName: "Collector A"},
		distInfo:         &model.DistributorInfoRow{DistributorID: 67, DistributorCode: "3434", DistributorName: "Distributor iDetama"},
		expenses: []model.ExpenseRow{
			{ExpenseTypeID: 10, ExpenseTypeName: "Transport", Note: "Only A", Amount: 10000},
		},
	}
	repoCollectorB := &detailRepoStub{
		visitDistributor: &model.VisitInformationRow{EmpID: 359, EmpCode: "EMP-359", EmpName: "Collector B"},
		distInfo:         &model.DistributorInfoRow{DistributorID: 67, DistributorCode: "3434", DistributorName: "Distributor iDetama"},
		expenses: []model.ExpenseRow{
			{ExpenseTypeID: 11, ExpenseTypeName: "Meal", Note: "Only B", Amount: 20000},
		},
	}

	svcA := &liveMonitoringService{repository: repoCollectorA, validate: validator.New()}
	svcB := &liveMonitoringService{repository: repoCollectorB, validate: validator.New()}
	distributorID := 67

	resultA, err := svcA.GetMonitoringDetail(context.Background(), request.LiveMonitoringDetailRequest{
		EmpID:         358,
		DistributorID: &distributorID,
		Date:          "2026-03-18",
	}, "C220010001", 111)
	if err != nil {
		t.Fatalf("collector A unexpected error: %v", err)
	}
	resultB, err := svcB.GetMonitoringDetail(context.Background(), request.LiveMonitoringDetailRequest{
		EmpID:         359,
		DistributorID: &distributorID,
		Date:          "2026-03-18",
	}, "C220010001", 111)
	if err != nil {
		t.Fatalf("collector B unexpected error: %v", err)
	}

	if len(resultA.Expense) != 1 || resultA.Expense[0].Note != "Only A" {
		t.Fatalf("collector A expense = %+v, want only collector A data", resultA.Expense)
	}
	if len(resultB.Expense) != 1 || resultB.Expense[0].Note != "Only B" {
		t.Fatalf("collector B expense = %+v, want only collector B data", resultB.Expense)
	}
	if repoCollectorA.receivedEmpID != 358 || repoCollectorB.receivedEmpID != 359 {
		t.Fatalf("received emp ids = %d and %d, want 358 and 359", repoCollectorA.receivedEmpID, repoCollectorB.receivedEmpID)
	}
	if repoCollectorA.receivedDate != "2026-03-18" || repoCollectorB.receivedDate != "2026-03-18" {
		t.Fatalf("received dates = %s and %s, want same request date", repoCollectorA.receivedDate, repoCollectorB.receivedDate)
	}
}

func TestGetMonitoringDetail_PrincipalIncludesExtraCallSummary(t *testing.T) {
	activityTime := "2026-05-22T08:00:00Z"
	repo := &detailRepoStub{
		visitPrincipalHistory: &model.VisitInformationRow{
			EmpID:     482,
			EmpCode:   "SLS482",
			EmpName:   "Salesman Test",
			Plan:      3,
			ExtraCall: 1,
			OnGoing:   0,
			Visited:   2,
			TotalSkip: 1,
			Matched:   2,
		},
		activityTime: &activityTime,
	}
	svc := &liveMonitoringService{repository: repo, validate: validator.New()}
	result, err := svc.GetMonitoringDetail(context.Background(), request.LiveMonitoringDetailRequest{
		EmpID: 482,
		Date:  "2026-05-22",
	}, "C26002", 999)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result == nil {
		t.Fatal("result is nil, expected non-nil for principal with data")
	}
	if result.VisitInformation.Planned != 3 {
		t.Fatalf("planned = %d, want 3", result.VisitInformation.Planned)
	}
	if result.VisitInformation.ExtraCall != 1 {
		t.Fatalf("extra_call = %d, want 1", result.VisitInformation.ExtraCall)
	}
	if result.VisitInformation.Visited != 2 {
		t.Fatalf("visited = %d, want 2", result.VisitInformation.Visited)
	}
	if result.VisitInformation.Skipped != 1 {
		t.Fatalf("skipped = %d, want 1", result.VisitInformation.Skipped)
	}
}

func TestGetMonitoringDetail_PrincipalNoDataReturnsNil(t *testing.T) {
	repo := &detailRepoStub{
		visitPrincipalHistory: nil,
	}
	svc := &liveMonitoringService{repository: repo, validate: validator.New()}
	result, err := svc.GetMonitoringDetail(context.Background(), request.LiveMonitoringDetailRequest{
		EmpID: 999,
		Date:  "2026-05-22",
	}, "C26002", 999)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != nil {
		t.Fatalf("expected nil result for no-data principal, got %+v", result)
	}
}

func TestGetMonitoringDetail_PrincipalPassesEmpIDToRepo(t *testing.T) {
	activityTime := "2026-05-22T08:00:00Z"
	repo := &detailRepoStub{
		visitPrincipalHistory: &model.VisitInformationRow{
			EmpID: 484, EmpCode: "SLS484", EmpName: "Test",
		},
		activityTime: &activityTime,
	}
	svc := &liveMonitoringService{repository: repo, validate: validator.New()}
	_, err := svc.GetMonitoringDetail(context.Background(), request.LiveMonitoringDetailRequest{
		EmpID: 484,
		Date:  "2026-05-22",
	}, "C26002", 999)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if repo.receivedPrincipalEmpID != 484 {
		t.Fatalf("repo received empID = %d, want 484", repo.receivedPrincipalEmpID)
	}
}

func TestGetMonitoringDetail_IncludesSurveyData(t *testing.T) {
	repo := &detailRepoStub{
		visitDistributor: &model.VisitInformationRow{EmpID: 421, EmpCode: "EMP-421", EmpName: "Collector"},
		distInfo:         &model.DistributorInfoRow{DistributorID: 67, DistributorCode: "3434", DistributorName: "Distributor"},
		surveys: []model.SurveyDataRow{
			{Submission: 2, SurveyTitle: "Survey A", OutletCode: "OUT-101", OutletName: "Outlet 101"},
		},
	}
	svc := &liveMonitoringService{repository: repo, validate: validator.New()}
	distributorID := 67
	result, err := svc.GetMonitoringDetail(context.Background(), request.LiveMonitoringDetailRequest{
		EmpID:         421,
		DistributorID: &distributorID,
		Date:          "2026-05-28",
	}, "C220010001", 999)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.SurveyData) != 1 {
		t.Fatalf("len(survey_data) = %d, want 1", len(result.SurveyData))
	}
	if result.SurveyData[0].Submission != 2 || result.SurveyData[0].SurveyTitle != "Survey A" {
		t.Fatalf("survey_data[0] = %+v, want mapped value", result.SurveyData[0])
	}
}

func TestGetMonitoringDetail_NoSurveyReturnsEmptyList(t *testing.T) {
	repo := &detailRepoStub{
		visitDistributor: &model.VisitInformationRow{EmpID: 421, EmpCode: "EMP-421", EmpName: "Collector"},
		distInfo:         &model.DistributorInfoRow{DistributorID: 67, DistributorCode: "3434", DistributorName: "Distributor"},
	}
	svc := &liveMonitoringService{repository: repo, validate: validator.New()}
	distributorID := 67
	result, err := svc.GetMonitoringDetail(context.Background(), request.LiveMonitoringDetailRequest{
		EmpID:         421,
		DistributorID: &distributorID,
		Date:          "2026-05-28",
	}, "C220010001", 999)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.SurveyData == nil || len(result.SurveyData) != 0 {
		t.Fatalf("survey_data = %+v, want empty non-nil list", result.SurveyData)
	}
}

func TestGetMonitoringDetail_IncludesCollectionPaid(t *testing.T) {
	repo := &detailRepoStub{
		visitDistributor: &model.VisitInformationRow{EmpID: 421, EmpCode: "EMP-421", EmpName: "Collector"},
		distInfo:         &model.DistributorInfoRow{DistributorID: 67, DistributorCode: "3434", DistributorName: "Distributor"},
		collections: []model.CollectionRow{
			{OutletID: 101, OutletCode: "OUT-101", OutletName: "Outlet 101", CollectionTotal: 1500000},
		},
	}
	svc := &liveMonitoringService{repository: repo, validate: validator.New()}
	distributorID := 67
	result, err := svc.GetMonitoringDetail(context.Background(), request.LiveMonitoringDetailRequest{
		EmpID:         421,
		DistributorID: &distributorID,
		Date:          "2026-05-28",
	}, "C220010001", 999)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Collection) != 1 {
		t.Fatalf("len(collection) = %d, want 1", len(result.Collection))
	}
	if result.Collection[0].CollectionTotal == nil || *result.Collection[0].CollectionTotal != 1500000 {
		t.Fatalf("collection_total = %+v, want 1500000", result.Collection[0].CollectionTotal)
	}
	if result.VisitInformation.CollectionSummary.Count != 1 || result.VisitInformation.CollectionSummary.Status != "completed" {
		t.Fatalf("collection_summary = %+v, want count=1 status=completed", result.VisitInformation.CollectionSummary)
	}
}

func TestGetMonitoringDetail_NoCollectionReturnsEmptyList(t *testing.T) {
	repo := &detailRepoStub{
		visitDistributor: &model.VisitInformationRow{EmpID: 421, EmpCode: "EMP-421", EmpName: "Collector"},
		distInfo:         &model.DistributorInfoRow{DistributorID: 67, DistributorCode: "3434", DistributorName: "Distributor"},
		collections:      []model.CollectionRow{},
	}
	svc := &liveMonitoringService{repository: repo, validate: validator.New()}
	distributorID := 67
	result, err := svc.GetMonitoringDetail(context.Background(), request.LiveMonitoringDetailRequest{
		EmpID:         421,
		DistributorID: &distributorID,
		Date:          "2026-05-28",
	}, "C220010001", 999)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Collection) != 0 {
		t.Fatalf("len(collection) = %d, want 0", len(result.Collection))
	}
	if result.VisitInformation.CollectionSummary.Count != 0 || result.VisitInformation.CollectionSummary.Status != "none" {
		t.Fatalf("collection_summary = %+v, want count=0 status=none", result.VisitInformation.CollectionSummary)
	}
}

func TestGetMonitoringDetail_CollectionUsesRequestDateAndEmpID(t *testing.T) {
	repo := &detailRepoStub{
		visitDistributor: &model.VisitInformationRow{EmpID: 421, EmpCode: "EMP-421", EmpName: "Collector"},
		distInfo:         &model.DistributorInfoRow{DistributorID: 67, DistributorCode: "3434", DistributorName: "Distributor"},
	}
	svc := &liveMonitoringService{repository: repo, validate: validator.New()}
	distributorID := 67
	_, err := svc.GetMonitoringDetail(context.Background(), request.LiveMonitoringDetailRequest{
		EmpID:         421,
		DistributorID: &distributorID,
		Date:          "2026-05-28",
	}, "C220010001", 999)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if repo.receivedCollectionDate != "2026-05-28" {
		t.Fatalf("collection date = %s, want 2026-05-28", repo.receivedCollectionDate)
	}
	if repo.receivedCollectionEmpID != 421 {
		t.Fatalf("collection emp_id = %d, want 421", repo.receivedCollectionEmpID)
	}
	if len(repo.receivedCollectionCustIDs) != 1 || repo.receivedCollectionCustIDs[0] != "C220010001" {
		t.Fatalf("collection cust_ids = %+v, want [C220010001]", repo.receivedCollectionCustIDs)
	}
}
