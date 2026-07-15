package service

import (
	"bytes"
	"context"
	"errors"
	"io"
	"strings"
	"testing"
	"time"

	"sales/entity"
	"sales/model"
	"sales/pkg/rabbitmq"
	"sales/repository"

	"github.com/xuri/excelize/v2"
)

// mockConfigEnv satisfies env.ConfigEnv for tests without loading a .env file.
type mockConfigEnv struct {
	values map[string]string
}

func (m *mockConfigEnv) Get(key string) string {
	return m.values[key]
}

type mockReportRepositoryForService struct {
	repository.ReportRepository
	storedReport                                 *model.ReportList
	updatedReport                                *model.ReportList
	getReportByReportIDFn                        func(reportID string) (model.ReportList, error)
	secondarySalesUnionFn                        func(filter entity.SecondarySalesReportQueryFilter) ([]model.SecondarySalesReportUnion, error)
	updateReportByReportIDFn                     func(c context.Context, reportID string, data *model.ReportList) error
	existsCustomerInParentScopeFn                func(custID string, parentCustID string) (bool, error)
	secondarySalesReportSumReportByMonthFn       func(custIDs []string, req entity.SecondarySalesReportDashboardSumPayload, year int) (model.SumReportByMonthModel, error)
	secondarySalesReportReturnSumReportByMonthFn func(custIDs []string, month int, year int) (model.SumReportReturnByMonthModel, error)
	secondarySalesReportGroupOutletFn            func(custIDs []string, month int, year int) ([]model.SecondarySalesReportGroup, error)
	secondarySalesReportGroupSalesmanFn          func(custIDs []string, month int, year int) ([]model.SecondarySalesReportGroup, error)
	secondarySalesReportProductCategoryFn        func(custIDs []string, month int, year int) ([]model.SecondarySalesReportGroup, error)
	secondarySalesReportProductFn                func(custIDs []string, month int, year int) ([]model.SecondarySalesReportGroup, error)
	secondarySalesReportTrendSalesFn             func(custIDs []string, year int) ([]model.TrendSalesSecondarySalesModel, error)
	salesmanActivityReportSumByMonthFn           func(custIDs []string, month int, year int) (model.SalesmanActivitySumByMonthModel, error)
	salesmanActivityReportTrendSalesFn           func(custIDs []string, year int) ([]model.ActivityReportTrendSalesModel, error)
	activityReportGeotagFn                       func(parentCustID string, custIDs []string, year int, empID *int) ([]model.ActivityReportGeotagRow, error)
	activitySalesmanReportGroupSalesmanFn        func(custIDs []string, month int, year int) ([]model.SecondarySalesReportGroup, error)
	activitySalesmanReturnReportGroupSalesmanFn  func(custIDs []string, month int, year int) ([]model.ReturnReportGroup, error)
}

func (m *mockReportRepositoryForService) StoreReportList(c context.Context, data *model.ReportList) error {
	copied := *data
	m.storedReport = &copied
	return nil
}

func (m *mockReportRepositoryForService) UpdateReportByReportID(c context.Context, reportID string, data *model.ReportList) error {
	if m.updateReportByReportIDFn != nil {
		return m.updateReportByReportIDFn(c, reportID, data)
	}
	copied := *data
	m.updatedReport = &copied
	return nil
}

func (m *mockReportRepositoryForService) GetReportByReportID(reportID string) (model.ReportList, error) {
	if m.getReportByReportIDFn != nil {
		return m.getReportByReportIDFn(reportID)
	}
	return model.ReportList{ReportID: reportID, ReportName: "SecondarySales-220426-001"}, nil
}

func (m *mockReportRepositoryForService) SecondarySalesUnion(filter entity.SecondarySalesReportQueryFilter) ([]model.SecondarySalesReportUnion, error) {
	if m.secondarySalesUnionFn != nil {
		return m.secondarySalesUnionFn(filter)
	}
	return nil, nil
}

func (m *mockReportRepositoryForService) CountSecondarySalesReportByDate(dataFilter entity.SecondarySalesReportQueryFilter) int64 {
	return 1
}

func (m *mockReportRepositoryForService) ExistsCustomerInParentScope(custID string, parentCustID string) (bool, error) {
	if m.existsCustomerInParentScopeFn != nil {
		return m.existsCustomerInParentScopeFn(custID, parentCustID)
	}
	return false, nil
}

func (m *mockReportRepositoryForService) SecondarySalesReportSumReportByMonth(custIDs []string, req entity.SecondarySalesReportDashboardSumPayload, year int) (model.SumReportByMonthModel, error) {
	if m.secondarySalesReportSumReportByMonthFn != nil {
		return m.secondarySalesReportSumReportByMonthFn(custIDs, req, year)
	}
	return model.SumReportByMonthModel{}, nil
}

func (m *mockReportRepositoryForService) SecondarySalesReportReturnSumReportByMonth(custIDs []string, month int, year int) (model.SumReportReturnByMonthModel, error) {
	if m.secondarySalesReportReturnSumReportByMonthFn != nil {
		return m.secondarySalesReportReturnSumReportByMonthFn(custIDs, month, year)
	}
	return model.SumReportReturnByMonthModel{}, nil
}

func (m *mockReportRepositoryForService) SecondarySalesReportGroupOutlet(custIDs []string, month int, year int) ([]model.SecondarySalesReportGroup, error) {
	if m.secondarySalesReportGroupOutletFn != nil {
		return m.secondarySalesReportGroupOutletFn(custIDs, month, year)
	}
	return nil, nil
}

func (m *mockReportRepositoryForService) SecondarySalesReportGroupSalesman(custIDs []string, month int, year int) ([]model.SecondarySalesReportGroup, error) {
	if m.secondarySalesReportGroupSalesmanFn != nil {
		return m.secondarySalesReportGroupSalesmanFn(custIDs, month, year)
	}
	return nil, nil
}

func (m *mockReportRepositoryForService) SecondarySalesReportProductCategory(custIDs []string, month int, year int) ([]model.SecondarySalesReportGroup, error) {
	if m.secondarySalesReportProductCategoryFn != nil {
		return m.secondarySalesReportProductCategoryFn(custIDs, month, year)
	}
	return nil, nil
}

func (m *mockReportRepositoryForService) SecondarySalesReportProduct(custIDs []string, month int, year int) ([]model.SecondarySalesReportGroup, error) {
	if m.secondarySalesReportProductFn != nil {
		return m.secondarySalesReportProductFn(custIDs, month, year)
	}
	return nil, nil
}

func (m *mockReportRepositoryForService) SecondarySalesReportTrendSales(custIDs []string, year int) ([]model.TrendSalesSecondarySalesModel, error) {
	if m.secondarySalesReportTrendSalesFn != nil {
		return m.secondarySalesReportTrendSalesFn(custIDs, year)
	}
	return nil, nil
}

func (m *mockReportRepositoryForService) SalesmanActivityReportSumByMonth(custIDs []string, month int, year int) (model.SalesmanActivitySumByMonthModel, error) {
	if m.salesmanActivityReportSumByMonthFn != nil {
		return m.salesmanActivityReportSumByMonthFn(custIDs, month, year)
	}
	return model.SalesmanActivitySumByMonthModel{}, nil
}

func (m *mockReportRepositoryForService) SalesmanActivityReportTrendSales(custIDs []string, year int) ([]model.ActivityReportTrendSalesModel, error) {
	if m.salesmanActivityReportTrendSalesFn != nil {
		return m.salesmanActivityReportTrendSalesFn(custIDs, year)
	}
	return nil, nil
}

func (m *mockReportRepositoryForService) ActivityReportGeotag(parentCustID string, custIDs []string, year int, empID *int) ([]model.ActivityReportGeotagRow, error) {
	if m.activityReportGeotagFn != nil {
		return m.activityReportGeotagFn(parentCustID, custIDs, year, empID)
	}
	return nil, nil
}

func (m *mockReportRepositoryForService) ActivitySalesmanReportGroupSalesman(custIDs []string, month int, year int) ([]model.SecondarySalesReportGroup, error) {
	if m.activitySalesmanReportGroupSalesmanFn != nil {
		return m.activitySalesmanReportGroupSalesmanFn(custIDs, month, year)
	}
	return nil, nil
}

func (m *mockReportRepositoryForService) ActivitySalesmanReturnReportGroupSalesman(custIDs []string, month int, year int) ([]model.ReturnReportGroup, error) {
	if m.activitySalesmanReturnReportGroupSalesmanFn != nil {
		return m.activitySalesmanReturnReportGroupSalesmanFn(custIDs, month, year)
	}
	return nil, nil
}

type mockObsAdapter struct {
	uploadFileCsvFn func(req *model.UploadCsv) (string, error)
}

func (m *mockObsAdapter) UploadFile(req *model.Upload) (string, error) {
	return "", nil
}

func (m *mockObsAdapter) UploadFileCsv(req *model.UploadCsv) (string, error) {
	if m.uploadFileCsvFn != nil {
		return m.uploadFileCsvFn(req)
	}
	return "", nil
}

func TestBuildReportObjectFileNameUsesUniqueReportID(t *testing.T) {
	t.Parallel()

	got := buildReportObjectFileName("SecondarySales-220426-001", "69e8b56dc7df9759991d7868")
	want := "reports/SecondarySales-220426-001.xlsx"

	if got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}

func TestSecondarySalesReportNameFormatUsesSequenceSuffix(t *testing.T) {
	t.Parallel()

	reportName := "SecondarySales-220426-" + "001"
	want := "SecondarySales-220426-001"

	if reportName != want {
		t.Fatalf("expected %q, got %q", want, reportName)
	}
}

func TestBuildReportObjectFileNameFallsBackToReportID(t *testing.T) {
	t.Parallel()

	got := buildReportObjectFileName("", "69e8b56dc7df9759991d7868")
	want := "reports/69e8b56dc7df9759991d7868.xlsx"

	if got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}

func TestSecondarySalesExcelRowUsesDocumentDateAndProductFields(t *testing.T) {
	t.Parallel()

	row := model.SecondarySalesReportUnion{
		DistributorCode: "DST-01",
		DistributorName: "Distributor A",
		TrxType:         "RETURN",
		DocumentNo:      "RET-001",
		DocumentDate:    time.Date(2026, time.April, 16, 0, 0, 0, 0, time.UTC),
		OutletCode:      "OUT-01",
		OutletName:      "Outlet A",
		EmpCode:         "EMP-01",
		EmpName:         "Jane",
		SupCode:         "SUP001",
		SupName:         "Classic Jersey Inc",
		ProCode:         "LPI-003",
		ProName:         "Jersey Arema Indonesia",
		SellPrice3:      300,
		SellPrice2:      200,
		SellPrice1:      100,
		UnitID3:         "CTN",
		UnitID2:         "PAC",
		UnitID1:         "PCS",
		ConvUnit3:       12,
		ConvUnit2:       6,
		Qty3:            -1,
		Qty2:            0,
		Qty1:            -5,
		GrossSales:      -1500,
		SpecialDiscount: 0,
		Discount:        50,
		NetSalesExcPPN:  -1300,
		PPN:             130,
		NetSalesIncPPN:  -1430,
	}

	rec := secondarySalesExcelRow(row)

	if got := rec[3]; got != "RET-001" {
		t.Fatalf("unexpected document number: got=%v", got)
	}
	if got := rec[4]; got != "16-04-2026" {
		t.Fatalf("unexpected document date: got=%v", got)
	}
	if got := rec[9]; got != "SUP001" {
		t.Fatalf("unexpected supplier code: got=%v", got)
	}
	if got := rec[11]; got != "LPI-003" {
		t.Fatalf("unexpected product code: got=%v", got)
	}
	if got := rec[12]; got != "Jersey Arema Indonesia" {
		t.Fatalf("unexpected product name: got=%v", got)
	}
}

func TestSubscribeSecondarySalesReportMarksFailedWhenUploadReturnsEmptyURL(t *testing.T) {
	repo := &mockReportRepositoryForService{
		secondarySalesUnionFn: func(filter entity.SecondarySalesReportQueryFilter) ([]model.SecondarySalesReportUnion, error) {
			return []model.SecondarySalesReportUnion{{DocumentNo: "INV-1", DocumentDate: time.Date(2026, time.April, 15, 0, 0, 0, 0, time.UTC)}}, nil
		},
	}
	obs := &mockObsAdapter{uploadFileCsvFn: func(req *model.UploadCsv) (string, error) { return "", nil }}
	service := &reportServiceImpl{ReportRepository: repo, ObsAdapter: obs}

	err := service.SubscribeSecondarySalesReport(entity.SecondarySalesReportQueryFilter{ReportID: "report-1"})
	if err == nil {
		t.Fatal("expected error when upload returns empty url")
	}
	if repo.updatedReport == nil {
		t.Fatal("expected failed report update")
	}
	if repo.updatedReport.FileStatus != entity.FILE_STATUS_FAILED {
		t.Fatalf("expected failed status, got %d", repo.updatedReport.FileStatus)
	}
	if repo.updatedReport.FileURL != "" {
		t.Fatalf("expected empty file url, got %q", repo.updatedReport.FileURL)
	}
}

func TestSubscribeSecondarySalesReportUploadsWorkbookWithExpectedValues(t *testing.T) {
	repo := &mockReportRepositoryForService{
		secondarySalesUnionFn: func(filter entity.SecondarySalesReportQueryFilter) ([]model.SecondarySalesReportUnion, error) {
			return []model.SecondarySalesReportUnion{{
				DistributorCode: "DST-01",
				DistributorName: "Distributor A",
				TrxType:         "ORDER",
				DocumentNo:      "INV-001",
				DocumentDate:    time.Date(2026, time.April, 15, 0, 0, 0, 0, time.UTC),
				OutletCode:      "OUT-01",
				OutletName:      "Outlet A",
				EmpCode:         "EMP-01",
				EmpName:         "Jane",
				SupCode:         "SUP001",
				SupName:         "Classic Jersey Inc",
				ProCode:         "LPI-003",
				ProName:         "Jersey Arema Indonesia",
				SellPrice3:      300,
				SellPrice2:      200,
				SellPrice1:      100,
				UnitID3:         "CTN",
				UnitID2:         "PAC",
				UnitID1:         "PCS",
				ConvUnit3:       12,
				ConvUnit2:       6,
				Qty3:            1,
				Qty2:            0,
				Qty1:            5,
				GrossSales:      1500,
				SpecialDiscount: 25,
				Discount:        50,
				NetSalesExcPPN:  1300,
				PPN:             130,
				NetSalesIncPPN:  1430,
			}}, nil
		},
	}

	var uploaded []byte
	obs := &mockObsAdapter{uploadFileCsvFn: func(req *model.UploadCsv) (string, error) {
		data, err := io.ReadAll(req.FileData)
		if err != nil {
			return "", err
		}
		uploaded = data
		return "https://files.example/reports/SecondarySales-220426-001.xlsx", nil
	}}
	service := &reportServiceImpl{ReportRepository: repo, ObsAdapter: obs}

	if err := service.SubscribeSecondarySalesReport(entity.SecondarySalesReportQueryFilter{ReportID: "report-2"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if repo.updatedReport == nil {
		t.Fatal("expected ready report update")
	}
	if repo.updatedReport.FileStatus != entity.FILE_STATUS_READY {
		t.Fatalf("expected ready status, got %d", repo.updatedReport.FileStatus)
	}
	if repo.updatedReport.FileURL == "" {
		t.Fatal("expected file url to be set")
	}

	f, err := excelize.OpenReader(bytes.NewReader(uploaded))
	if err != nil {
		t.Fatalf("failed to open uploaded workbook: %v", err)
	}
	defer func() {
		_ = f.Close()
	}()

	if got, err := f.GetCellValue("Report", "D2"); err != nil || got != "INV-001" {
		t.Fatalf("unexpected document number cell: got=%q err=%v", got, err)
	}
	if got, err := f.GetCellValue("Report", "E2"); err != nil || got != "15-04-2026" {
		t.Fatalf("unexpected document date cell: got=%q err=%v", got, err)
	}
	if got, err := f.GetCellValue("Report", "J2"); err != nil || got != "SUP001" {
		t.Fatalf("unexpected supplier code cell: got=%q err=%v", got, err)
	}
	if got, err := f.GetCellValue("Report", "L2"); err != nil || got != "LPI-003" {
		t.Fatalf("unexpected product code cell: got=%q err=%v", got, err)
	}
	if got, err := f.GetCellValue("Report", "M2"); err != nil || got != "Jersey Arema Indonesia" {
		t.Fatalf("unexpected product name cell: got=%q err=%v", got, err)
	}
	if got, err := f.GetCellValue("Report", "V2"); err != nil || got != "1" {
		t.Fatalf("unexpected qty3 cell: got=%q err=%v", got, err)
	}
}

func TestMarkReportFailedKeepsOriginalErrorWhenStatusUpdateFails(t *testing.T) {
	repo := &mockReportRepositoryForService{
		secondarySalesUnionFn: func(filter entity.SecondarySalesReportQueryFilter) ([]model.SecondarySalesReportUnion, error) {
			return nil, errors.New("query failed")
		},
		updateReportByReportIDFn: func(c context.Context, reportID string, data *model.ReportList) error {
			return errors.New("update failed")
		},
	}
	service := &reportServiceImpl{ReportRepository: repo, ObsAdapter: &mockObsAdapter{}}

	err := service.SubscribeSecondarySalesReport(entity.SecondarySalesReportQueryFilter{ReportID: "report-3"})
	if err == nil || err.Error() != "query failed" {
		t.Fatalf("expected original query error, got %v", err)
	}
}

func TestPublishSecondarySalesReportMarksFailedWhenRabbitMQPublishFails(t *testing.T) {
	t.Cleanup(func() {
		publishReportMessage = rabbitmq.PublishMessage
	})
	publishReportMessage = func(rmq *rabbitmq.RmqConfig) error {
		return errors.New("publish failed")
	}

	from := time.Date(2026, time.April, 15, 0, 0, 0, 0, time.UTC).Unix()
	to := time.Date(2026, time.April, 16, 0, 0, 0, 0, time.UTC).Unix()
	repo := &mockReportRepositoryForService{}
	service := &reportServiceImpl{
		Config:           &mockConfigEnv{values: map[string]string{"REPORT_DELAY_SECONDS": "0"}},
		ReportRepository: repo,
	}

	_, err := service.PublishSecondarySalesReport(entity.SecondarySalesReportQueryFilter{
		CustID:   "CUST-1",
		ExportBy: "tester",
		From:     &from,
		To:       &to,
	})
	if err == nil || err.Error() != "publish failed" {
		t.Fatalf("expected publish failed error, got %v", err)
	}
	if repo.storedReport == nil {
		t.Fatal("expected processing report to be stored before publish")
	}
	if repo.storedReport.FileStatus != entity.FILE_STATUS_PROCESSING {
		t.Fatalf("expected stored processing status, got %d", repo.storedReport.FileStatus)
	}
	if repo.updatedReport == nil {
		t.Fatal("expected failed report update after publish error")
	}
	if repo.updatedReport.ReportID != repo.storedReport.ReportID {
		t.Fatalf("expected failed update for report %q, got %q", repo.storedReport.ReportID, repo.updatedReport.ReportID)
	}
	if repo.updatedReport.FileStatus != entity.FILE_STATUS_FAILED {
		t.Fatalf("expected failed status, got %d", repo.updatedReport.FileStatus)
	}
	if repo.updatedReport.FileURL != "" {
		t.Fatalf("expected empty file url on failed publish, got %q", repo.updatedReport.FileURL)
	}
}

func TestSecondarySalesReportSumReportByMonthUsesChildCustAndExplicitYear(t *testing.T) {
	t.Parallel()

	var gotCustID string
	var gotMonth int
	var gotYear int
	repo := &mockReportRepositoryForService{
		existsCustomerInParentScopeFn: func(custID string, parentCustID string) (bool, error) {
			if custID != "CHILD1" {
				t.Fatalf("unexpected cust id: %s", custID)
			}
			if parentCustID != "PARENT1" {
				t.Fatalf("unexpected parent cust id: %s", parentCustID)
			}
			return true, nil
		},
		secondarySalesReportSumReportByMonthFn: func(custIDs []string, req entity.SecondarySalesReportDashboardSumPayload, year int) (model.SumReportByMonthModel, error) {
			if len(custIDs) != 1 {
				t.Fatalf("expected single cust id, got %#v", custIDs)
			}
			gotCustID = custIDs[0]
			gotMonth = req.Month
			gotYear = year
			return model.SumReportByMonthModel{}, nil
		},
		secondarySalesReportReturnSumReportByMonthFn: func(custIDs []string, month int, year int) (model.SumReportReturnByMonthModel, error) {
			return model.SumReportReturnByMonthModel{}, nil
		},
	}
	service := &reportServiceImpl{ReportRepository: repo}
	year := 2026

	_, err := service.SecondarySalesReportSumReportByMonth("PARENT1", "PARENT1", entity.SecondarySalesReportDashboardSumPayload{
		Month:  5,
		Year:   &year,
		CustID: "CHILD1",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotCustID != "CHILD1" || gotMonth != 5 || gotYear != 2026 {
		t.Fatalf("unexpected repository args: cust=%s month=%d year=%d", gotCustID, gotMonth, gotYear)
	}
}

func TestSecondarySalesReportSumReportByMonthRejectsUnauthorizedCustID(t *testing.T) {
	t.Parallel()

	called := false
	repo := &mockReportRepositoryForService{
		secondarySalesReportSumReportByMonthFn: func(custIDs []string, req entity.SecondarySalesReportDashboardSumPayload, year int) (model.SumReportByMonthModel, error) {
			called = true
			return model.SumReportByMonthModel{}, nil
		},
	}
	service := &reportServiceImpl{ReportRepository: repo}
	year := 2026

	_, err := service.SecondarySalesReportSumReportByMonth("DIST1", "PARENT1", entity.SecondarySalesReportDashboardSumPayload{
		Month:  5,
		Year:   &year,
		CustID: "SIBLING1",
	})
	if !errors.Is(err, ErrUnauthorizedCustID) {
		t.Fatalf("expected ErrUnauthorizedCustID, got %v", err)
	}
	if called {
		t.Fatal("expected report repository not called")
	}
}

func TestSecondarySalesReportSumReportByMonthFallsBackToCurrentYearAndAuthCust(t *testing.T) {
	t.Parallel()

	currentYear := time.Now().Year()
	var gotCustID string
	var gotYear int
	repo := &mockReportRepositoryForService{
		existsCustomerInParentScopeFn: func(custID string, parentCustID string) (bool, error) {
			t.Fatal("scope lookup should not be called for empty cust_id")
			return false, nil
		},
		secondarySalesReportSumReportByMonthFn: func(custIDs []string, req entity.SecondarySalesReportDashboardSumPayload, year int) (model.SumReportByMonthModel, error) {
			if len(custIDs) != 1 {
				t.Fatalf("expected single cust id, got %#v", custIDs)
			}
			gotCustID = custIDs[0]
			gotYear = year
			return model.SumReportByMonthModel{}, nil
		},
		secondarySalesReportReturnSumReportByMonthFn: func(custIDs []string, month int, year int) (model.SumReportReturnByMonthModel, error) {
			return model.SumReportReturnByMonthModel{}, nil
		},
	}
	service := &reportServiceImpl{ReportRepository: repo}

	_, err := service.SecondarySalesReportSumReportByMonth("AUTH1", "PARENT1", entity.SecondarySalesReportDashboardSumPayload{Month: 5})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotCustID != "AUTH1" {
		t.Fatalf("expected auth cust fallback, got %s", gotCustID)
	}
	if gotYear != currentYear {
		t.Fatalf("expected current year %d, got %d", currentYear, gotYear)
	}
}

func TestSecondarySalesReportSumReportByMonthUsesRepositoryReturnRateWhenOrderQtyZero(t *testing.T) {
	t.Parallel()

	repo := &mockReportRepositoryForService{
		secondarySalesReportSumReportByMonthFn: func(custIDs []string, req entity.SecondarySalesReportDashboardSumPayload, year int) (model.SumReportByMonthModel, error) {
			return model.SumReportByMonthModel{Qty: 0, QtyReturn: 7, ReturnRate: 12.34}, nil
		},
		secondarySalesReportReturnSumReportByMonthFn: func(custIDs []string, month int, year int) (model.SumReportReturnByMonthModel, error) {
			return model.SumReportReturnByMonthModel{Qty: 7}, nil
		},
	}
	service := &reportServiceImpl{ReportRepository: repo}
	year := 2026

	data, err := service.SecondarySalesReportSumReportByMonth("AUTH1", "PARENT1", entity.SecondarySalesReportDashboardSumPayload{Month: 5, Year: &year})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if data.ReturnRate != 12.34 {
		t.Fatalf("expected repository return rate 12.34, got %v", data.ReturnRate)
	}
}

func TestSecondarySalesReportSumReportByMonthPassesOptionalFiltersToRepository(t *testing.T) {
	t.Parallel()

	from := int64(1778086800)
	to := int64(1778173199)
	var gotReq entity.SecondarySalesReportDashboardSumPayload
	repo := &mockReportRepositoryForService{
		secondarySalesReportSumReportByMonthFn: func(custIDs []string, req entity.SecondarySalesReportDashboardSumPayload, year int) (model.SumReportByMonthModel, error) {
			gotReq = req
			if len(custIDs) != 1 || custIDs[0] != "AUTH1" {
				t.Fatalf("unexpected custIDs %#v", custIDs)
			}
			if year != 2026 {
				t.Fatalf("unexpected year %d", year)
			}
			return model.SumReportByMonthModel{}, nil
		},
		secondarySalesReportReturnSumReportByMonthFn: func(custIDs []string, month int, year int) (model.SumReportReturnByMonthModel, error) {
			return model.SumReportReturnByMonthModel{}, nil
		},
	}
	service := &reportServiceImpl{ReportRepository: repo}
	year := 2026

	_, err := service.SecondarySalesReportSumReportByMonth("AUTH1", "PARENT1", entity.SecondarySalesReportDashboardSumPayload{
		Month:       5,
		Year:        &year,
		From:        &from,
		To:          &to,
		OutletIDs:   []int64{301},
		SalesmanIDs: []int64{101},
		ProIDs:      []int64{501},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotReq.From == nil || *gotReq.From != from {
		t.Fatalf("expected from filter %d, got %#v", from, gotReq.From)
	}
	if gotReq.To == nil || *gotReq.To != to {
		t.Fatalf("expected to filter %d, got %#v", to, gotReq.To)
	}
	if gotReq.Month != 5 {
		t.Fatalf("expected month 5, got %d", gotReq.Month)
	}
	if gotReq.Year == nil || *gotReq.Year != 2026 {
		t.Fatalf("expected year 2026, got %#v", gotReq.Year)
	}
	if len(gotReq.OutletIDs) != 1 || gotReq.OutletIDs[0] != 301 {
		t.Fatalf("unexpected outlet filters %#v", gotReq.OutletIDs)
	}
	if len(gotReq.SalesmanIDs) != 1 || gotReq.SalesmanIDs[0] != 101 {
		t.Fatalf("unexpected salesman filters %#v", gotReq.SalesmanIDs)
	}
	if len(gotReq.ProIDs) != 1 || gotReq.ProIDs[0] != 501 {
		t.Fatalf("unexpected product filters %#v", gotReq.ProIDs)
	}
}

func TestSecondarySalesReportSumReportByMonthMapsPPNAndNetSalesExcPPN(t *testing.T) {
	t.Parallel()

	orderUpdatedAt := time.Date(2026, time.June, 20, 10, 0, 0, 0, time.UTC)
	returnUpdatedAt := time.Date(2026, time.June, 18, 9, 0, 0, 0, time.UTC)
	repo := &mockReportRepositoryForService{
		secondarySalesReportSumReportByMonthFn: func(custIDs []string, req entity.SecondarySalesReportDashboardSumPayload, year int) (model.SumReportByMonthModel, error) {
			return model.SumReportByMonthModel{
				TotalGrossSales:    1_000,
				TotalDiscountPromo: 150,
				TotalPPN:           85,
				NetSalesExcPPN:     850,
				NetSales:           935,
				Qty:                10,
				QtyReturn:          2,
				ReturnRate:         22.22,
				NetSalesReturn:     100,
				LastUpdate:         &orderUpdatedAt,
			}, nil
		},
		secondarySalesReportReturnSumReportByMonthFn: func(custIDs []string, month int, year int) (model.SumReportReturnByMonthModel, error) {
			return model.SumReportReturnByMonthModel{LastUpdate: &returnUpdatedAt}, nil
		},
	}
	service := &reportServiceImpl{ReportRepository: repo}
	year := 2026

	data, err := service.SecondarySalesReportSumReportByMonth("AUTH1", "PARENT1", entity.SecondarySalesReportDashboardSumPayload{Month: 6, Year: &year})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if data.TotalPPN != 85 {
		t.Fatalf("expected total_ppn 85, got %v", data.TotalPPN)
	}
	if data.NetSalesExcPPN != 850 {
		t.Fatalf("expected net_sales_exc_ppn 850, got %v", data.NetSalesExcPPN)
	}
	if data.NetSales != 935 {
		t.Fatalf("expected net_sales 935, got %v", data.NetSales)
	}
	if data.QtyReturn != 2 || data.NetSalesReturn != 100 || data.ReturnRate != 22.22 {
		t.Fatalf("unexpected return metrics: qty_return=%d net_sales_return=%v return_rate=%v", data.QtyReturn, data.NetSalesReturn, data.ReturnRate)
	}
	if data.LastUpdate == nil || !data.LastUpdate.Equal(orderUpdatedAt) {
		t.Fatalf("expected summary last_update to remain order timestamp, got %#v", data.LastUpdate)
	}
}

func TestSecondarySalesReportSumReportByMonthPropagatesSubtractFromRepository(t *testing.T) {
	t.Parallel()

	repo := &mockReportRepositoryForService{
		secondarySalesReportSumReportByMonthFn: func(custIDs []string, req entity.SecondarySalesReportDashboardSumPayload, year int) (model.SumReportByMonthModel, error) {
			return model.SumReportByMonthModel{
				Qty:                134,
				TotalDiscountPromo: 1_238_740,
			}, nil
		},
		secondarySalesReportReturnSumReportByMonthFn: func(custIDs []string, month int, year int) (model.SumReportReturnByMonthModel, error) {
			return model.SumReportReturnByMonthModel{}, nil
		},
	}
	service := &reportServiceImpl{ReportRepository: repo}
	year := 2026

	data, err := service.SecondarySalesReportSumReportByMonth("AUTH1", "PARENT1", entity.SecondarySalesReportDashboardSumPayload{Month: 6, Year: &year})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if data.Qty != 134 {
		t.Fatalf("expected qty 134, got %d", data.Qty)
	}
	if data.TotalDiscountPromo != 1_238_740 {
		t.Fatalf("expected total_discount_promo 1238740, got %v", data.TotalDiscountPromo)
	}
}

func TestSecondarySalesReportGroupSalesUsesFallbackYearForAllBranches(t *testing.T) {
	t.Parallel()

	currentYear := time.Now().Year()
	cases := []struct {
		name    string
		groupBy string
		setup   func(repo *mockReportRepositoryForService, t *testing.T)
	}{
		{
			name:    "outlet",
			groupBy: entity.SECONDARY_SALES_GROUP_OUTLET,
			setup: func(repo *mockReportRepositoryForService, t *testing.T) {
				repo.secondarySalesReportGroupOutletFn = func(custIDs []string, month int, year int) ([]model.SecondarySalesReportGroup, error) {
					if len(custIDs) != 1 || custIDs[0] != "AUTH1" || month != 5 || year != currentYear {
						t.Fatalf("unexpected args: cust=%#v month=%d year=%d", custIDs, month, year)
					}
					return []model.SecondarySalesReportGroup{{ID: 1, Code: "OUT-1", Name: "Outlet", NetSales: 10}}, nil
				}
			},
		},
		{
			name:    "salesman",
			groupBy: entity.SECONDARY_SALES_GROUP_SALESMAN,
			setup: func(repo *mockReportRepositoryForService, t *testing.T) {
				repo.secondarySalesReportGroupSalesmanFn = func(custIDs []string, month int, year int) ([]model.SecondarySalesReportGroup, error) {
					if len(custIDs) != 1 || custIDs[0] != "AUTH1" || month != 5 || year != currentYear {
						t.Fatalf("unexpected args: cust=%#v month=%d year=%d", custIDs, month, year)
					}
					return []model.SecondarySalesReportGroup{{ID: 2, Code: "SLS-2", Name: "Salesman", NetSales: 20}}, nil
				}
			},
		},
		{
			name:    "product_category",
			groupBy: entity.SECONDARY_SALES_GROUP_PRODUCT_CATEGORY,
			setup: func(repo *mockReportRepositoryForService, t *testing.T) {
				repo.secondarySalesReportProductCategoryFn = func(custIDs []string, month int, year int) ([]model.SecondarySalesReportGroup, error) {
					if len(custIDs) != 1 || custIDs[0] != "AUTH1" || month != 5 || year != currentYear {
						t.Fatalf("unexpected args: cust=%#v month=%d year=%d", custIDs, month, year)
					}
					return []model.SecondarySalesReportGroup{{ID: 3, Code: "CAT-3", Name: "Category", NetSales: 30}}, nil
				}
			},
		},
		{
			name:    "default product",
			groupBy: "",
			setup: func(repo *mockReportRepositoryForService, t *testing.T) {
				repo.secondarySalesReportProductFn = func(custIDs []string, month int, year int) ([]model.SecondarySalesReportGroup, error) {
					if len(custIDs) != 1 || custIDs[0] != "AUTH1" || month != 5 || year != currentYear {
						t.Fatalf("unexpected args: cust=%#v month=%d year=%d", custIDs, month, year)
					}
					return []model.SecondarySalesReportGroup{{ID: 4, Code: "PRO-4", Name: "Product", NetSales: 40}}, nil
				}
			},
		},
		{
			name:    "unknown group falls back to product",
			groupBy: "unknown_group",
			setup: func(repo *mockReportRepositoryForService, t *testing.T) {
				repo.secondarySalesReportProductFn = func(custIDs []string, month int, year int) ([]model.SecondarySalesReportGroup, error) {
					if len(custIDs) != 1 || custIDs[0] != "AUTH1" || month != 5 || year != currentYear {
						t.Fatalf("unexpected args: cust=%#v month=%d year=%d", custIDs, month, year)
					}
					return []model.SecondarySalesReportGroup{{ID: 5, Code: "PRO-5", Name: "Product Fallback", NetSales: 50}}, nil
				}
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			repo := &mockReportRepositoryForService{}
			tc.setup(repo, t)
			service := &reportServiceImpl{ReportRepository: repo}

			data, err := service.SecondarySalesReportGroupSales("AUTH1", "PARENT1", entity.SecondarySalesReportDashboardGroupPayload{Month: 5, GroupBy: tc.groupBy})
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if len(data) != 1 {
				t.Fatalf("expected single row, got %d", len(data))
			}
			if data[0].Code == "" {
				t.Fatalf("expected group code to be mapped, got %#v", data[0])
			}
		})
	}
}

func TestResolveSecondaryDashboardCustIDsAllowsPrincipalMultiChildren(t *testing.T) {
	repo := &mockReportRepositoryForService{
		existsCustomerInParentScopeFn: func(custID string, parentCustID string) (bool, error) {
			return custID == "CHILD1" || custID == "CHILD2", nil
		},
	}
	service := &reportServiceImpl{ReportRepository: repo}

	custIDs, err := service.resolveSecondaryDashboardCustIDs("PARENT1", "PARENT1", []string{"CHILD1", "CHILD2"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(custIDs) != 2 || custIDs[0] != "CHILD1" || custIDs[1] != "CHILD2" {
		t.Fatalf("unexpected cust ids: %#v", custIDs)
	}
}

func TestResolveSecondaryDashboardCustIDsRejectsDistributorSiblingMulti(t *testing.T) {
	service := &reportServiceImpl{ReportRepository: &mockReportRepositoryForService{}}
	_, err := service.resolveSecondaryDashboardCustIDs("DIST1", "PARENT1", []string{"DIST1", "SIBLING1"})
	if !errors.Is(err, ErrUnauthorizedCustID) {
		t.Fatalf("expected ErrUnauthorizedCustID, got %v", err)
	}
}

func TestApplyActivityReportListCustIDsAllowsPrincipalMultiBusinessUnit(t *testing.T) {
	repo := &mockReportRepositoryForService{
		existsCustomerInParentScopeFn: func(custID string, parentCustID string) (bool, error) {
			return custID == "C26002" || custID == "C26003", nil
		},
	}
	service := &reportServiceImpl{ReportRepository: repo}

	filter := entity.ActivityReportQueryFilterList{
		AuthCustID:   "C26002",
		CustID:       "C26002",
		ParentCustID: "C26002",
		CustIDs:      []string{"C26002", "C26003"},
	}
	if err := service.applyActivityReportListCustIDs(&filter); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(filter.CustIDs) != 2 || filter.CustIDs[0] != "C26002" || filter.CustIDs[1] != "C26003" {
		t.Fatalf("unexpected resolved cust ids: %#v", filter.CustIDs)
	}
	if filter.CustID != "C26002" {
		t.Fatalf("expected auth cust id preserved for multi select, got %q", filter.CustID)
	}
}

func TestSalesmanActivityReportSumReportByMonthMapsNetSalesAndYear(t *testing.T) {
	lastUpdate := time.Date(2026, 6, 12, 0, 1, 0, 0, time.UTC)
	repo := &mockReportRepositoryForService{
		existsCustomerInParentScopeFn: func(custID, parentCustID string) (bool, error) {
			return custID == "C260020001" && parentCustID == "C26002", nil
		},
		salesmanActivityReportSumByMonthFn: func(custIDs []string, month int, year int) (model.SalesmanActivitySumByMonthModel, error) {
			if len(custIDs) != 1 || custIDs[0] != "C260020001" || month != 6 || year != 2026 {
				t.Fatalf("unexpected repo args: custIDs=%v month=%d year=%d", custIDs, month, year)
			}
			return model.SalesmanActivitySumByMonthModel{
				TotalSales:    1250000.5,
				TotalReturn:   50000,
				TotalSalesman: 5,
				LastUpdate:    &lastUpdate,
			}, nil
		},
	}
	service := &reportServiceImpl{ReportRepository: repo}
	year := 2026

	data, err := service.SalesmanActivityReportSumReportByMonth("C26002", "C26002", entity.SalesmanActivityReportDashboardSumPayload{
		Month:   6,
		Year:    &year,
		CustID:  "C260020001",
		CustIDs: []string{"C260020001"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if data.TotalSales != 1250000.5 || data.TotalReturn != 50000 || data.SalesmanTotal != 5 {
		t.Fatalf("unexpected response: %#v", data)
	}
	if data.LastUpdate == nil || !data.LastUpdate.Equal(lastUpdate) {
		t.Fatalf("unexpected last update: %#v", data.LastUpdate)
	}
}

func TestApplyActivityReportListCustIDsRejectsDistributorMultiBusinessUnit(t *testing.T) {
	service := &reportServiceImpl{ReportRepository: &mockReportRepositoryForService{}}
	filter := entity.ActivityReportQueryFilterList{
		AuthCustID:   "C260020001",
		CustID:       "C260020001",
		ParentCustID: "C26002",
		CustIDs:      []string{"C260020001", "C260030001"},
	}
	err := service.applyActivityReportListCustIDs(&filter)
	if !errors.Is(err, ErrUnauthorizedCustID) {
		t.Fatalf("expected ErrUnauthorizedCustID, got %v", err)
	}
}

func TestPublishSecondarySalesReportUsesEffectiveCustButStoresAuthOwner(t *testing.T) {
	t.Cleanup(func() {
		publishReportMessage = rabbitmq.PublishMessage
	})

	from := time.Date(2026, time.April, 15, 0, 0, 0, 0, time.UTC).Unix()
	to := time.Date(2026, time.April, 16, 0, 0, 0, 0, time.UTC).Unix()
	repo := &mockReportRepositoryForService{
		existsCustomerInParentScopeFn: func(custID string, parentCustID string) (bool, error) {
			if custID != "CHILD1" || parentCustID != "PARENT1" {
				t.Fatalf("unexpected scope args cust=%s parent=%s", custID, parentCustID)
			}
			return true, nil
		},
	}
	var published rabbitmq.RmqConfig
	publishReportMessage = func(rmq *rabbitmq.RmqConfig) error {
		published = *rmq
		return nil
	}
	service := &reportServiceImpl{
		Config:           &mockConfigEnv{values: map[string]string{"REPORT_DELAY_SECONDS": "0"}},
		ReportRepository: repo,
	}

	_, err := service.PublishSecondarySalesReport(entity.SecondarySalesReportQueryFilter{
		CustID:          "PARENT1",
		ParentCustID:    "PARENT1",
		RequestedCustID: "CHILD1",
		ExportBy:        "tester",
		From:            &from,
		To:              &to,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if repo.storedReport == nil {
		t.Fatal("expected report list stored")
	}
	if repo.storedReport.CustID != "PARENT1" {
		t.Fatalf("expected stored auth cust_id, got %s", repo.storedReport.CustID)
	}
	if !bytes.Contains([]byte(published.Message), []byte(`"_cust_id":"CHILD1"`)) {
		t.Fatalf("expected published effective cust payload, got %s", published.Message)
	}
	if !bytes.Contains([]byte(published.Message), []byte(`"cust_ids":["CHILD1"]`)) {
		t.Fatalf("expected cust_ids payload, got %s", published.Message)
	}
	if !bytes.Contains([]byte(published.Message), []byte(`"cust_ids":["CHILD1"]`)) {
		t.Fatalf("expected requested cust payload, got %s", published.Message)
	}
}

func TestPublishSecondarySalesReportRejectsUnauthorizedDistributorSibling(t *testing.T) {
	t.Cleanup(func() {
		publishReportMessage = rabbitmq.PublishMessage
	})
	publishCalled := false
	publishReportMessage = func(rmq *rabbitmq.RmqConfig) error {
		publishCalled = true
		return nil
	}

	from := time.Date(2026, time.April, 15, 0, 0, 0, 0, time.UTC).Unix()
	to := time.Date(2026, time.April, 16, 0, 0, 0, 0, time.UTC).Unix()
	repo := &mockReportRepositoryForService{}
	service := &reportServiceImpl{
		Config:           &mockConfigEnv{values: map[string]string{"REPORT_DELAY_SECONDS": "0"}},
		ReportRepository: repo,
	}

	_, err := service.PublishSecondarySalesReport(entity.SecondarySalesReportQueryFilter{
		CustID:          "DIST1",
		ParentCustID:    "PARENT1",
		RequestedCustID: "SIBLING1",
		ExportBy:        "tester",
		From:            &from,
		To:              &to,
	})
	if !errors.Is(err, ErrUnauthorizedCustID) {
		t.Fatalf("expected ErrUnauthorizedCustID, got %v", err)
	}
	if repo.storedReport != nil {
		t.Fatal("expected no stored report on unauthorized cust")
	}
	if publishCalled {
		t.Fatal("expected no publish on unauthorized cust")
	}
}

func TestPublishSecondarySalesReportFallsBackToAuthCustWithoutScopeLookup(t *testing.T) {
	t.Cleanup(func() {
		publishReportMessage = rabbitmq.PublishMessage
	})

	from := time.Date(2026, time.April, 15, 0, 0, 0, 0, time.UTC).Unix()
	to := time.Date(2026, time.April, 16, 0, 0, 0, 0, time.UTC).Unix()
	repo := &mockReportRepositoryForService{
		existsCustomerInParentScopeFn: func(custID string, parentCustID string) (bool, error) {
			t.Fatal("scope lookup should not be called for empty requested cust")
			return false, nil
		},
	}
	var published rabbitmq.RmqConfig
	publishReportMessage = func(rmq *rabbitmq.RmqConfig) error {
		published = *rmq
		return nil
	}
	service := &reportServiceImpl{
		Config:           &mockConfigEnv{values: map[string]string{"REPORT_DELAY_SECONDS": "0"}},
		ReportRepository: repo,
	}

	_, err := service.PublishSecondarySalesReport(entity.SecondarySalesReportQueryFilter{
		CustID:       "AUTH1",
		ParentCustID: "PARENT1",
		ExportBy:     "tester",
		From:         &from,
		To:           &to,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !bytes.Contains([]byte(published.Message), []byte(`"_cust_id":"AUTH1"`)) {
		t.Fatalf("expected auth cust payload fallback, got %s", published.Message)
	}
}

func TestSubscribeSecondarySalesReportUsesRequestedCustAsEffectiveFallback(t *testing.T) {
	var gotFilter entity.SecondarySalesReportQueryFilter
	repo := &mockReportRepositoryForService{
		secondarySalesUnionFn: func(filter entity.SecondarySalesReportQueryFilter) ([]model.SecondarySalesReportUnion, error) {
			gotFilter = filter
			return nil, nil
		},
	}
	obs := &mockObsAdapter{uploadFileCsvFn: func(req *model.UploadCsv) (string, error) {
		return "https://files.example/reports/report.xlsx", nil
	}}
	service := &reportServiceImpl{ReportRepository: repo, ObsAdapter: obs}

	err := service.SubscribeSecondarySalesReport(entity.SecondarySalesReportQueryFilter{
		ReportID:         "report-4",
		CustID:           "AUTH1",
		RequestedCustID:  "CHILD1",
		RequestedCustIDs: entity.StringListOrScalar{"CHILD1"},
		ParentCustID:     "PARENT1",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotFilter.CustID != "CHILD1" {
		t.Fatalf("expected subscribe to use requested cust as effective cust, got %s", gotFilter.CustID)
	}
}

func TestSubscribeSecondarySalesReportSupportsMultiCustPayload(t *testing.T) {
	var gotFilter entity.SecondarySalesReportQueryFilter
	repo := &mockReportRepositoryForService{
		secondarySalesUnionFn: func(filter entity.SecondarySalesReportQueryFilter) ([]model.SecondarySalesReportUnion, error) {
			gotFilter = filter
			return nil, nil
		},
	}
	obs := &mockObsAdapter{uploadFileCsvFn: func(req *model.UploadCsv) (string, error) {
		return "https://files.example/reports/report.xlsx", nil
	}}
	service := &reportServiceImpl{ReportRepository: repo, ObsAdapter: obs}

	err := service.SubscribeSecondarySalesReport(entity.SecondarySalesReportQueryFilter{
		ReportID:     "report-4b",
		CustID:       "AUTH1",
		CustIDs:      []string{"CHILD1", "CHILD2"},
		ParentCustID: "PARENT1",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !bytes.Equal([]byte(strings.Join(gotFilter.CustIDs, ",")), []byte("CHILD1,CHILD2")) {
		t.Fatalf("expected subscribe to preserve multi cust ids, got %#v", gotFilter.CustIDs)
	}
}

func TestSecondarySalesReportTrendSalesUsesChildCustWhenAllowed(t *testing.T) {
	var gotCustID string
	var gotYear int
	repo := &mockReportRepositoryForService{
		existsCustomerInParentScopeFn: func(custID string, parentCustID string) (bool, error) {
			if custID != "CHILD1" || parentCustID != "PARENT1" {
				t.Fatalf("unexpected scope args cust=%s parent=%s", custID, parentCustID)
			}
			return true, nil
		},
		secondarySalesReportTrendSalesFn: func(custIDs []string, year int) ([]model.TrendSalesSecondarySalesModel, error) {
			if len(custIDs) != 1 {
				t.Fatalf("expected single cust id, got %#v", custIDs)
			}
			gotCustID = custIDs[0]
			gotYear = year
			return nil, nil
		},
	}

	service := &reportServiceImpl{ReportRepository: repo}

	_, err := service.SecondarySalesReportTrendSales("PARENT1", "PARENT1", 2026, []string{"CHILD1"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotCustID != "CHILD1" || gotYear != 2026 {
		t.Fatalf("unexpected repo args cust=%s year=%d", gotCustID, gotYear)
	}
}

func TestSecondarySalesReportTrendSalesRejectsUnauthorizedDistributorSibling(t *testing.T) {
	called := false
	repo := &mockReportRepositoryForService{
		secondarySalesReportTrendSalesFn: func(custIDs []string, year int) ([]model.TrendSalesSecondarySalesModel, error) {
			called = true
			return nil, nil
		},
	}
	service := &reportServiceImpl{ReportRepository: repo}

	_, err := service.SecondarySalesReportTrendSales("DIST1", "PARENT1", 2026, []string{"SIBLING1"})
	if !errors.Is(err, ErrUnauthorizedCustID) {
		t.Fatalf("expected ErrUnauthorizedCustID, got %v", err)
	}
	if called {
		t.Fatal("expected repo not called")
	}
}

func TestSecondarySalesReportTrendSalesFallsBackToAuthCust(t *testing.T) {
	var gotCustID string
	repo := &mockReportRepositoryForService{
		existsCustomerInParentScopeFn: func(custID string, parentCustID string) (bool, error) {
			t.Fatal("scope lookup should not be called for empty requested cust")
			return false, nil
		},
		secondarySalesReportTrendSalesFn: func(custIDs []string, year int) ([]model.TrendSalesSecondarySalesModel, error) {
			if len(custIDs) != 1 {
				t.Fatalf("expected single cust id, got %#v", custIDs)
			}
			gotCustID = custIDs[0]
			return nil, nil
		},
	}
	service := &reportServiceImpl{ReportRepository: repo}

	_, err := service.SecondarySalesReportTrendSales("AUTH1", "PARENT1", 2026, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotCustID != "AUTH1" {
		t.Fatalf("expected auth cust fallback, got %s", gotCustID)
	}
}

func TestSalesmanActivityReportTrendSalesUsesChildCustWhenAllowed(t *testing.T) {
	var gotCustID string
	var gotYear int
	repo := &mockReportRepositoryForService{
		existsCustomerInParentScopeFn: func(custID string, parentCustID string) (bool, error) {
			if custID != "CHILD1" || parentCustID != "PARENT1" {
				t.Fatalf("unexpected scope args cust=%s parent=%s", custID, parentCustID)
			}
			return true, nil
		},
		salesmanActivityReportTrendSalesFn: func(custIDs []string, year int) ([]model.ActivityReportTrendSalesModel, error) {
			if len(custIDs) != 1 {
				t.Fatalf("expected single cust id, got %#v", custIDs)
			}
			gotCustID = custIDs[0]
			gotYear = year
			return nil, nil
		},
	}

	service := &reportServiceImpl{ReportRepository: repo}

	_, err := service.SalesmanActivityReportTrendSales("PARENT1", "PARENT1", 2026, []string{"CHILD1"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotCustID != "CHILD1" || gotYear != 2026 {
		t.Fatalf("unexpected repo args cust=%s year=%d", gotCustID, gotYear)
	}
}

func TestSalesmanActivityReportTrendSalesRejectsUnauthorizedDistributorSibling(t *testing.T) {
	repo := &mockReportRepositoryForService{
		salesmanActivityReportTrendSalesFn: func(custIDs []string, year int) ([]model.ActivityReportTrendSalesModel, error) {
			t.Fatal("repository should not be called for unauthorized cust")
			return nil, nil
		},
	}
	service := &reportServiceImpl{ReportRepository: repo}

	_, err := service.SalesmanActivityReportTrendSales("DIST1", "PARENT1", 2026, []string{"SIBLING1"})
	if !errors.Is(err, ErrUnauthorizedCustID) {
		t.Fatalf("expected ErrUnauthorizedCustID, got %v", err)
	}
}

func TestSalesmanActivityReportTrendSalesFallsBackToAuthCust(t *testing.T) {
	var gotCustID string
	repo := &mockReportRepositoryForService{
		salesmanActivityReportTrendSalesFn: func(custIDs []string, year int) ([]model.ActivityReportTrendSalesModel, error) {
			if len(custIDs) != 1 {
				t.Fatalf("expected single cust id, got %#v", custIDs)
			}
			gotCustID = custIDs[0]
			return nil, nil
		},
	}
	service := &reportServiceImpl{ReportRepository: repo}

	_, err := service.SalesmanActivityReportTrendSales("AUTH1", "PARENT1", 2026, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotCustID != "AUTH1" {
		t.Fatalf("expected auth cust fallback, got %s", gotCustID)
	}
}

func TestSalesmanActivityReportGroupSalesUsesChildCustWhenAllowed(t *testing.T) {
	var gotCustIDs []string
	repo := &mockReportRepositoryForService{
		existsCustomerInParentScopeFn: func(custID string, parentCustID string) (bool, error) {
			if custID != "C260020001" || parentCustID != "C26002" {
				t.Fatalf("unexpected scope args cust=%s parent=%s", custID, parentCustID)
			}
			return true, nil
		},
		activitySalesmanReportGroupSalesmanFn: func(custIDs []string, month int, year int) ([]model.SecondarySalesReportGroup, error) {
			gotCustIDs = custIDs
			return []model.SecondarySalesReportGroup{{ID: 471, Code: "SLS003", Name: "Yabes Roni", NetSales: 29810000}}, nil
		},
	}
	service := &reportServiceImpl{ReportRepository: repo}

	data, err := service.SalesmanActivityReportGroupSales("C26002", "C26002", entity.SalesmanActivityReportDashboardGroupPayload{
		Month:        6,
		ActivityType: entity.ACTIVITY_SALESMAN_GROUP_SALES,
		CustIDs:      []string{"C260020001"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(gotCustIDs) != 1 || gotCustIDs[0] != "C260020001" {
		t.Fatalf("expected child cust id, got %#v", gotCustIDs)
	}
	if len(data) != 1 || data[0].ID != 471 {
		t.Fatalf("unexpected data: %#v", data)
	}
}

func TestSalesmanActivityReportGroupSalesRejectsUnauthorizedCust(t *testing.T) {
	repo := &mockReportRepositoryForService{
		activitySalesmanReportGroupSalesmanFn: func(custIDs []string, month int, year int) ([]model.SecondarySalesReportGroup, error) {
			t.Fatal("repository should not be called for unauthorized cust")
			return nil, nil
		},
	}
	service := &reportServiceImpl{ReportRepository: repo}

	_, err := service.SalesmanActivityReportGroupSales("DIST1", "PARENT1", entity.SalesmanActivityReportDashboardGroupPayload{
		Month:        6,
		ActivityType: entity.ACTIVITY_SALESMAN_GROUP_SALES,
		CustIDs:      []string{"SIBLING1"},
	})
	if !errors.Is(err, ErrUnauthorizedCustID) {
		t.Fatalf("expected ErrUnauthorizedCustID, got %v", err)
	}
}

func TestFormatPrincipalActivityReportColumns(t *testing.T) {
	tests := []struct {
		name         string
		row          model.SalesActivityReportRow
		wantBU       string
		wantDistCode string
		wantDistName string
	}{
		{
			name: "principal outlet visit",
			row: model.SalesActivityReportRow{
				BusinessUnitCode: "",
				DistributorCode:  "3100022",
				DistributorName:  "Dist A",
				OutletCode:       "OUT001",
			},
			wantBU: "-", wantDistCode: "-", wantDistName: "-",
		},
		{
			name: "principal distributor visit",
			row: model.SalesActivityReportRow{
				BusinessUnitCode: "",
				DistributorCode:  "3100022",
				DistributorName:  "Dist A",
				OutletCode:       "",
			},
			wantBU: "-", wantDistCode: "3100022", wantDistName: "Dist A",
		},
		{
			name: "distributor business unit row",
			row: model.SalesActivityReportRow{
				BusinessUnitCode: "3100063",
				DistributorCode:  "",
				DistributorName:  "",
				OutletCode:       "OUT002",
			},
			wantBU: "3100063", wantDistCode: "", wantDistName: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotBU, gotDistCode, gotDistName := formatPrincipalActivityReportColumns(tt.row)
			if gotBU != tt.wantBU || gotDistCode != tt.wantDistCode || gotDistName != tt.wantDistName {
				t.Fatalf("got (%q,%q,%q), want (%q,%q,%q)", gotBU, gotDistCode, gotDistName, tt.wantBU, tt.wantDistCode, tt.wantDistName)
			}
		})
	}
}

func TestMapActivityReportRowPrincipalUserDisplay(t *testing.T) {
	row := model.SalesActivityReportRow{
		BusinessUnitCode: "",
		BusinessUnitName: "Principal Co",
		DistributorCode:  "3100022",
		DistributorName:  "Dist A",
		OutletCode:       "OUT001",
		PJPCode:          "12",
		IsPlanned:        true,
	}

	mapped := mapActivityReportRow(row, true)
	if mapped.BusinessUnitCode != "-" {
		t.Fatalf("expected business unit code '-', got %q", mapped.BusinessUnitCode)
	}
	if mapped.DistributorCode != "-" || mapped.DistributorName != "-" {
		t.Fatalf("expected distributor '-', got %q / %q", mapped.DistributorCode, mapped.DistributorName)
	}

	nonPrincipal := mapActivityReportRow(row, false)
	if nonPrincipal.BusinessUnitCode != "" || nonPrincipal.DistributorCode != "3100022" {
		t.Fatalf("non-principal user should keep raw values")
	}
}

func TestMapActivityReportRowRemarks(t *testing.T) {
	t.Parallel()

	onLeave := mapActivityReportRow(model.SalesActivityReportRow{
		PJPCode:  "1",
		Remarks:  "On Leave",
	}, false)
	if onLeave.Remarks != "On Leave" {
		t.Fatalf("expected On Leave, got %q", onLeave.Remarks)
	}

	empty := mapActivityReportRow(model.SalesActivityReportRow{
		PJPCode: "1",
		Remarks: "",
	}, false)
	if empty.Remarks != "-" {
		t.Fatalf("expected default '-', got %q", empty.Remarks)
	}

	dash := mapActivityReportRow(model.SalesActivityReportRow{
		PJPCode: "1",
		Remarks: "-",
	}, false)
	if dash.Remarks != "-" {
		t.Fatalf("expected '-', got %q", dash.Remarks)
	}
}

func TestSalesmanActivityReportGeotagUsesChildCustWhenAllowed(t *testing.T) {
	t.Parallel()

	var gotCustIDs []string
	service := &reportServiceImpl{
		ReportRepository: &mockReportRepositoryForService{
			existsCustomerInParentScopeFn: func(custID string, parentCustID string) (bool, error) {
				return custID == "CHILD1" && parentCustID == "PARENT1", nil
			},
			activityReportGeotagFn: func(parentCustID string, custIDs []string, year int, empID *int) ([]model.ActivityReportGeotagRow, error) {
				gotCustIDs = custIDs
				return []model.ActivityReportGeotagRow{
					{SalesmanCode: 1, SalesmanName: "Budi", TotalVisit: 10, GeotagMatchCount: 2, GeotagUnmatchCount: 8, GeotagMatchPct: 20, GeotagUnmatchPct: 80},
				}, nil
			},
		},
	}

	data, err := service.SalesmanActivityReportGeotag("PARENT1", "PARENT1", entity.ActivityReportGeotagPayload{
		Year:    2026,
		CustIDs: []string{"CHILD1"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(gotCustIDs) != 1 || gotCustIDs[0] != "CHILD1" {
		t.Fatalf("expected CHILD1 cust filter, got %#v", gotCustIDs)
	}
	if data.TotalGeotagMatchPercentage != 20 || data.TotalGeotagUnmatchPercentage != 80 {
		t.Fatalf("unexpected totals: %+v", data)
	}
	if len(data.Details) != 1 || data.Details[0].SalesmanCode != "1" {
		t.Fatalf("unexpected details: %+v", data.Details)
	}
}

func TestSalesmanActivityReportGeotagRejectsUnauthorizedDistributorSibling(t *testing.T) {
	t.Parallel()

	service := &reportServiceImpl{
		ReportRepository: &mockReportRepositoryForService{
			activityReportGeotagFn: func(parentCustID string, custIDs []string, year int, empID *int) ([]model.ActivityReportGeotagRow, error) {
				t.Fatal("repository should not be called for unauthorized cust_id")
				return nil, nil
			},
		},
	}

	_, err := service.SalesmanActivityReportGeotag("DIST1", "PARENT1", entity.ActivityReportGeotagPayload{
		Year:    2026,
		CustIDs: []string{"SIBLING1"},
	})
	if !errors.Is(err, ErrUnauthorizedCustID) {
		t.Fatalf("expected ErrUnauthorizedCustID, got %v", err)
	}
}

func TestSalesmanActivityReportGeotagFallsBackToAuthCust(t *testing.T) {
	t.Parallel()

	var gotCustIDs []string
	service := &reportServiceImpl{
		ReportRepository: &mockReportRepositoryForService{
			activityReportGeotagFn: func(parentCustID string, custIDs []string, year int, empID *int) ([]model.ActivityReportGeotagRow, error) {
				gotCustIDs = custIDs
				return nil, nil
			},
		},
	}

	_, err := service.SalesmanActivityReportGeotag("AUTH1", "PARENT1", entity.ActivityReportGeotagPayload{Year: 2026})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(gotCustIDs) != 1 || gotCustIDs[0] != "AUTH1" {
		t.Fatalf("expected AUTH1 cust filter, got %#v", gotCustIDs)
	}
}
