package service

import (
	"bytes"
	"context"
	"encoding/base64"
	"finance/entity"
	"finance/model"
	"finance/repository"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/xuri/excelize/v2"
)

type paymentDepositReportRepositoryMock struct {
	repository.PaymentDepositReportRepository
	findAllFn      func(dataFilter entity.PaymentDepositReportQueryFilter, custId string) ([]model.PaymentDepositReportRow, int64, int, error)
	summaryFn      func(dataFilter entity.PaymentDepositReportQueryFilter, custId string) (model.PaymentDepositReportSummaryRow, error)
	findNoLimitFn  func(dataFilter entity.PaymentDepositReportQueryFilter, custId string) ([]model.PaymentDepositReportRow, error)
	findDownloadFn func(dataFilter entity.PaymentDepositReportQueryFilter, custId, parentCustId string) ([]model.PaymentDepositReportDownloadRow, error)
	findRecapFn    func(dataFilter entity.PaymentDepositReportQueryFilter, custId, parentCustId string) ([]model.PaymentDepositReportRecapRow, error)
	runningNumFn   func(custId string, date time.Time) (int, error)
	updateFn       func(c context.Context, reportID string, status int, fileBase64 string) error
}

func (m *paymentDepositReportRepositoryMock) FindAllPaymentDeposit(dataFilter entity.PaymentDepositReportQueryFilter, custId string) ([]model.PaymentDepositReportRow, int64, int, error) {
	if m.findAllFn != nil {
		return m.findAllFn(dataFilter, custId)
	}
	return nil, 0, 0, nil
}

func (m *paymentDepositReportRepositoryMock) FindPaymentDepositSummary(dataFilter entity.PaymentDepositReportQueryFilter, custId string) (model.PaymentDepositReportSummaryRow, error) {
	if m.summaryFn != nil {
		return m.summaryFn(dataFilter, custId)
	}
	return model.PaymentDepositReportSummaryRow{}, nil
}

func (m *paymentDepositReportRepositoryMock) FindAllPaymentDepositNoLimit(dataFilter entity.PaymentDepositReportQueryFilter, custId string) ([]model.PaymentDepositReportRow, error) {
	if m.findNoLimitFn != nil {
		return m.findNoLimitFn(dataFilter, custId)
	}
	return nil, nil
}

func (m *paymentDepositReportRepositoryMock) FindAllPaymentDepositDownload(dataFilter entity.PaymentDepositReportQueryFilter, custId, parentCustId string) ([]model.PaymentDepositReportDownloadRow, error) {
	if m.findDownloadFn != nil {
		return m.findDownloadFn(dataFilter, custId, parentCustId)
	}
	return nil, nil
}

func (m *paymentDepositReportRepositoryMock) FindPaymentDepositRecapRows(dataFilter entity.PaymentDepositReportQueryFilter, custId, parentCustId string) ([]model.PaymentDepositReportRecapRow, error) {
	if m.findRecapFn != nil {
		return m.findRecapFn(dataFilter, custId, parentCustId)
	}
	return nil, nil
}

func (m *paymentDepositReportRepositoryMock) InsertReportList(c context.Context, report model.ReportList) error {
	return nil
}

func (m *paymentDepositReportRepositoryMock) UpdateReportList(c context.Context, reportID string, status int, fileBase64 string) error {
	if m.updateFn != nil {
		return m.updateFn(c, reportID, status, fileBase64)
	}
	return nil
}

func (m *paymentDepositReportRepositoryMock) GetReportRunningNumber(custId string, date time.Time) (int, error) {
	if m.runningNumFn != nil {
		return m.runningNumFn(custId, date)
	}
	return 0, nil
}

func TestPaymentDepositReportService_ListReportMapsCollectorFields(t *testing.T) {
	collectorID := 381
	collectorCode := "COL-381"
	collectorName := "Collector 381"
	date := time.Date(2026, 2, 2, 0, 0, 0, 0, time.UTC)
	repoMock := &paymentDepositReportRepositoryMock{
		findAllFn: func(dataFilter entity.PaymentDepositReportQueryFilter, custId string) ([]model.PaymentDepositReportRow, int64, int, error) {
			return []model.PaymentDepositReportRow{{
				DepositDate:   date,
				DepositType:   "AR",
				DepositNo:     "DP001",
				CollectorID:   &collectorID,
				CollectorCode: &collectorCode,
				CollectorName: &collectorName,
				CashAmount:    10,
				TotalPayment:  10,
			}}, 1, 1, nil
		},
		summaryFn: func(dataFilter entity.PaymentDepositReportQueryFilter, custId string) (model.PaymentDepositReportSummaryRow, error) {
			return model.PaymentDepositReportSummaryRow{TotalCash: 10}, nil
		},
	}

	svc := &paymentDepositReportServiceImpl{Repo: repoMock}
	resp, err := svc.ListReport(entity.PaymentDepositReportQueryFilter{CustId: "C001", Page: 1, Limit: 20, StartDate: "2026-02-02", EndDate: "2026-02-02", DepositType: []string{"AR"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resp.Items) != 1 {
		t.Fatalf("expected 1 item, got %d", len(resp.Items))
	}
	if resp.Items[0].CollectorID == nil || *resp.Items[0].CollectorID != collectorID {
		t.Fatalf("expected collector_id %d, got %#v", collectorID, resp.Items[0].CollectorID)
	}
	if resp.Items[0].CollectorCode == nil || *resp.Items[0].CollectorCode != collectorCode {
		t.Fatalf("expected collector_code %s, got %#v", collectorCode, resp.Items[0].CollectorCode)
	}
	if resp.Items[0].CollectorName == nil || *resp.Items[0].CollectorName != collectorName {
		t.Fatalf("expected collector_name %s, got %#v", collectorName, resp.Items[0].CollectorName)
	}
}

func TestPaymentDepositReportService_ListReportKeepsAPNullCollector(t *testing.T) {
	repoMock := &paymentDepositReportRepositoryMock{
		findAllFn: func(dataFilter entity.PaymentDepositReportQueryFilter, custId string) ([]model.PaymentDepositReportRow, int64, int, error) {
			return []model.PaymentDepositReportRow{{
				DepositType:  "AP",
				DepositNo:    "AP001",
				CashAmount:   0,
				TotalPayment: 0,
			}}, 1, 1, nil
		},
		summaryFn: func(dataFilter entity.PaymentDepositReportQueryFilter, custId string) (model.PaymentDepositReportSummaryRow, error) {
			return model.PaymentDepositReportSummaryRow{}, nil
		},
	}

	svc := &paymentDepositReportServiceImpl{Repo: repoMock}
	resp, err := svc.ListReport(entity.PaymentDepositReportQueryFilter{CustId: "C001", Page: 1, Limit: 20, StartDate: "2026-02-02", EndDate: "2026-02-02", DepositType: []string{"AP"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resp.Items) != 1 {
		t.Fatalf("expected 1 item, got %d", len(resp.Items))
	}
	if resp.Items[0].CollectorID != nil || resp.Items[0].CollectorCode != nil || resp.Items[0].CollectorName != nil {
		t.Fatalf("expected AP collector fields to remain nil, got %#v", resp.Items[0])
	}
}

func TestPaymentDepositReportService_ListReportSummaryByDepositType(t *testing.T) {
	repoMock := &paymentDepositReportRepositoryMock{
		findAllFn: func(dataFilter entity.PaymentDepositReportQueryFilter, custId string) ([]model.PaymentDepositReportRow, int64, int, error) {
			return []model.PaymentDepositReportRow{{DepositType: "AR", DepositNo: "DP001"}}, 1, 1, nil
		},
		summaryFn: func(dataFilter entity.PaymentDepositReportQueryFilter, custId string) (model.PaymentDepositReportSummaryRow, error) {
			return model.PaymentDepositReportSummaryRow{}, nil
		},
		findRecapFn: func(dataFilter entity.PaymentDepositReportQueryFilter, custId, parentCustId string) ([]model.PaymentDepositReportRecapRow, error) {
			return []model.PaymentDepositReportRecapRow{
				{DepositType: "Account Receivable", Cash: 100, Expense: -25},
				{DepositType: "Account Payable", Cash: 50, Expense: 99},
			}, nil
		},
	}

	svc := &paymentDepositReportServiceImpl{Repo: repoMock}
	resp, err := svc.ListReport(entity.PaymentDepositReportQueryFilter{CustId: "C001", ParentCustId: "P001", Page: 1, Limit: 20, StartDate: "2026-02-02", EndDate: "2026-02-02", DepositType: []string{"AR", "AP"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resp.SummaryByDepositType) != 2 {
		t.Fatalf("expected 2 summary buckets, got %d", len(resp.SummaryByDepositType))
	}
	if resp.SummaryByDepositType[0].DepositTypeLabel != "Account Receivable" || resp.SummaryByDepositType[0].TotalExpense != -25 {
		t.Fatalf("unexpected AR bucket: %#v", resp.SummaryByDepositType[0])
	}
	if resp.SummaryByDepositType[1].DepositTypeLabel != "Account Payable" || resp.SummaryByDepositType[1].TotalExpense != 0 {
		t.Fatalf("unexpected AP bucket: %#v", resp.SummaryByDepositType[1])
	}
}

func TestPaymentDepositReportService_GenerateExcelRecapLayout(t *testing.T) {
	svc := &paymentDepositReportServiceImpl{}
	rows := []model.PaymentDepositReportDownloadRow{
		{DepositDate: time.Date(2026, 5, 5, 0, 0, 0, 0, time.UTC), DepositType: "Account Receivable", DepositNo: "AR1", Cash: 10},
		{DepositDate: time.Date(2026, 5, 6, 0, 0, 0, 0, time.UTC), DepositType: "Account Receivable", DepositNo: "AR1", Expense: -16},
		{DepositDate: time.Date(2026, 5, 7, 0, 0, 0, 0, time.UTC), DepositType: "Account Payable", DepositNo: "AP1", Cash: 20},
	}

	encoded, err := svc.generateExcel(rows, paymentDepositExportMetadata{
		StartDate:      time.Date(2026, 5, 5, 0, 0, 0, 0, time.UTC),
		EndDate:        time.Date(2026, 5, 7, 0, 0, 0, 0, time.UTC),
		CollectorLabel: "All",
	})
	if err != nil {
		t.Fatalf("generateExcel error: %v", err)
	}
	binary, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		t.Fatalf("decode base64 error: %v", err)
	}
	f, err := excelize.OpenReader(strings.NewReader(string(binary)))
	if err != nil {
		f, err = excelize.OpenReader(bytes.NewReader(binary))
	}
	if err != nil {
		t.Fatalf("open workbook error: %v", err)
	}
	sheet := "Payment Deposit Report"
	assertCellEquals(t, f, sheet, "B11", "Account Receivable")
	assertCellEquals(t, f, sheet, "E11", "Account Payable")
	assertCellEquals(t, f, sheet, "A12", "Total Cash")
	assertCellEquals(t, f, sheet, "D12", "Total Cash")
	assertCellEquals(t, f, sheet, "A19", "Total Expense")
	assertCellFloat(t, f, sheet, "B19", -16)
	assertCellEquals(t, f, sheet, "D19", "")
	assertCellEquals(t, f, sheet, "E19", "")
}

func TestPaymentDepositReportService_GenerateExcelExpenseNameMobileFormat(t *testing.T) {
	svc := &paymentDepositReportServiceImpl{}
	expenseName := "000 - Uang Parkir"
	rows := []model.PaymentDepositReportDownloadRow{
		{
			DepositDate: time.Date(2026, 6, 11, 0, 0, 0, 0, time.UTC),
			DepositType: "Account Receivable",
			DepositNo:   "E20260611001",
			Expense:     -5000,
			ExpenseName: &expenseName,
		},
	}

	encoded, err := svc.generateExcel(rows, paymentDepositExportMetadata{
		StartDate:      time.Date(2026, 6, 11, 0, 0, 0, 0, time.UTC),
		EndDate:        time.Date(2026, 6, 11, 0, 0, 0, 0, time.UTC),
		CollectorLabel: "All",
	})
	if err != nil {
		t.Fatalf("generateExcel error: %v", err)
	}
	binary, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		t.Fatalf("decode base64 error: %v", err)
	}
	f, err := excelize.OpenReader(bytes.NewReader(binary))
	if err != nil {
		t.Fatalf("open workbook error: %v", err)
	}
	assertCellEquals(t, f, "Payment Deposit Report", "Q6", "000 - Uang Parkir")
}

func TestPaymentDepositReportService_GenerateExcelUsesExplicitMetadata(t *testing.T) {
	svc := &paymentDepositReportServiceImpl{}
	collectorA := "Collector A"
	collectorB := "Collector B"
	rows := []model.PaymentDepositReportDownloadRow{
		{DepositDate: time.Date(2026, 5, 7, 0, 0, 0, 0, time.UTC), DepositType: "Account Receivable", DepositNo: "AR1", Collector: &collectorA, Cash: 10},
		{DepositDate: time.Date(2026, 5, 5, 0, 0, 0, 0, time.UTC), DepositType: "Account Payable", DepositNo: "AP1", Collector: &collectorB, Cash: 20},
	}

	encoded, err := svc.generateExcel(rows, paymentDepositExportMetadata{
		StartDate:      time.Date(2026, 5, 1, 0, 0, 0, 0, time.UTC),
		EndDate:        time.Date(2026, 5, 31, 0, 0, 0, 0, time.UTC),
		CollectorLabel: "All",
	})
	if err != nil {
		t.Fatalf("generateExcel error: %v", err)
	}
	binary, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		t.Fatalf("decode base64 error: %v", err)
	}
	f, err := excelize.OpenReader(bytes.NewReader(binary))
	if err != nil {
		t.Fatalf("open workbook error: %v", err)
	}
	sheet := "Payment Deposit Report"
	assertCellEquals(t, f, sheet, "B2", "01-05-2026 - 31-05-2026")
	assertCellEquals(t, f, sheet, "B3", "All")
}

func TestResolvePaymentDepositCollectorLabel(t *testing.T) {
	collectorA := "Collector A"
	collectorB := "Collector B"
	cases := []struct {
		name   string
		filter entity.PaymentDepositReportQueryFilter
		rows   []model.PaymentDepositReportDownloadRow
		want   string
	}{
		{
			name:   "no collector filter",
			filter: entity.PaymentDepositReportQueryFilter{},
			rows:   []model.PaymentDepositReportDownloadRow{{Collector: &collectorA}},
			want:   "All",
		},
		{
			name:   "multiple collectors in filter",
			filter: entity.PaymentDepositReportQueryFilter{EmpID: []string{"1", "2"}},
			rows:   []model.PaymentDepositReportDownloadRow{{Collector: &collectorA}},
			want:   "All",
		},
		{
			name:   "single collector filter and single name in rows",
			filter: entity.PaymentDepositReportQueryFilter{EmpID: []string{"1"}},
			rows:   []model.PaymentDepositReportDownloadRow{{Collector: &collectorA}, {Collector: &collectorA}},
			want:   "Collector A",
		},
		{
			name:   "single collector filter but mixed names in rows",
			filter: entity.PaymentDepositReportQueryFilter{EmpID: []string{"1"}},
			rows:   []model.PaymentDepositReportDownloadRow{{Collector: &collectorA}, {Collector: &collectorB}},
			want:   "All",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := resolvePaymentDepositCollectorLabel(tc.filter, tc.rows)
			if got != tc.want {
				t.Fatalf("collector label = %q, want %q", got, tc.want)
			}
		})
	}
}

func TestPaymentDepositReportService_DownloadReportPersistsDatesAndReadyState(t *testing.T) {
	var gotReportID string
	var gotStatus int
	var gotBase64 string
	updated := make(chan struct{}, 1)
	repoMock := &paymentDepositReportRepositoryMock{
		runningNumFn: func(custId string, date time.Time) (int, error) { return 2, nil },
		findDownloadFn: func(dataFilter entity.PaymentDepositReportQueryFilter, custId, parentCustId string) ([]model.PaymentDepositReportDownloadRow, error) {
			return []model.PaymentDepositReportDownloadRow{{
				DepositDate: time.Date(2026, 2, 2, 0, 0, 0, 0, time.UTC),
				DepositType: "AR",
				DepositNo:   "DP001",
				Cash:        10,
			}}, nil
		},
		updateFn: func(c context.Context, reportID string, status int, fileBase64 string) error {
			gotReportID = reportID
			gotStatus = status
			gotBase64 = fileBase64
			updated <- struct{}{}
			return nil
		},
	}
	svc := &paymentDepositReportServiceImpl{Repo: repoMock, Transaction: noopTransaction{}}
	resp, err := svc.DownloadReport(entity.PaymentDepositReportQueryFilter{CustId: "C001", ParentCustId: "P001", StartDate: "2026-02-02", EndDate: "2026-02-03", DepositType: []string{"AR"}}, "tester")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StartDate != "2026-02-02" || resp.EndDate != "2026-02-03" {
		t.Fatalf("expected persisted filter dates, got %s - %s", resp.StartDate, resp.EndDate)
	}
	if resp.FileStatus != 0 || resp.FileStatusName != entity.PaymentDepositReportStatusNameProcessing {
		t.Fatalf("expected processing response, got %#v", resp)
	}
	select {
	case <-updated:
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for async update")
	}
	if gotStatus != entity.PaymentDepositReportStatusReady {
		t.Fatalf("expected ready status update, got %d", gotStatus)
	}
	if gotReportID == "" || gotBase64 == "" {
		t.Fatalf("expected report id and base64 to be set, got id=%q base64=%q", gotReportID, gotBase64)
	}
	if !strings.HasPrefix(resp.ReportName, entity.PaymentDepositReportDownloadPrefix) {
		t.Fatalf("expected report name prefix %q, got %q", entity.PaymentDepositReportDownloadPrefix, resp.ReportName)
	}
}

type noopTransaction struct{}

func (noopTransaction) WithinTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	return fn(ctx)
}

func assertCellEquals(t *testing.T, f *excelize.File, sheet, cell, want string) {
	t.Helper()
	got, err := f.GetCellValue(sheet, cell)
	if err != nil {
		t.Fatalf("get cell %s error: %v", cell, err)
	}
	if got != want {
		t.Fatalf("cell %s = %q, want %q", cell, got, want)
	}
}

func assertCellFloat(t *testing.T, f *excelize.File, sheet, cell string, want float64) {
	t.Helper()
	got, err := f.GetCellValue(sheet, cell)
	if err != nil {
		t.Fatalf("get cell %s error: %v", cell, err)
	}
	parsed, err := strconv.ParseFloat(got, 64)
	if err != nil {
		t.Fatalf("parse cell %s value %q error: %v", cell, got, err)
	}
	if parsed != want {
		t.Fatalf("cell %s = %v, want %v", cell, parsed, want)
	}
}
