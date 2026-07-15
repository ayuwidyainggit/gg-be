package service

import (
	"bytes"
	"context"
	"encoding/base64"
	"testing"
	"time"

	"sales/entity"
	"sales/model"
	"sales/repository"

	"github.com/xuri/excelize/v2"
)

type mockSoRepository struct {
	repository.SoRepository
	findDownloadDataPoFn         func(filter entity.SoDownloadQueryFilter) ([]model.SoDownloadPo, error)
	findDownloadDataSoFn         func(filter entity.SoDownloadQueryFilter) ([]model.SoDownloadSo, error)
	findDownloadDataFinalFn      func(filter entity.SoDownloadQueryFilter) ([]model.SoDownloadFinal, error)
	findDownloadQtySummaryDataFn func(filter entity.SoDownloadQueryFilter) ([]model.SoDownloadQtySummary, error)
}

func (m *mockSoRepository) FindDownloadDataPo(filter entity.SoDownloadQueryFilter) ([]model.SoDownloadPo, error) {
	return m.findDownloadDataPoFn(filter)
}

func (m *mockSoRepository) FindDownloadDataSo(filter entity.SoDownloadQueryFilter) ([]model.SoDownloadSo, error) {
	return m.findDownloadDataSoFn(filter)
}

func (m *mockSoRepository) FindDownloadDataFinal(filter entity.SoDownloadQueryFilter) ([]model.SoDownloadFinal, error) {
	return m.findDownloadDataFinalFn(filter)
}

func (m *mockSoRepository) FindDownloadQtySummary(filter entity.SoDownloadQueryFilter) ([]model.SoDownloadQtySummary, error) {
	return m.findDownloadQtySummaryDataFn(filter)
}

type mockReportRepository struct {
	repository.ReportRepository
	updatedReport *model.ReportList
	storedReport  *model.ReportList
}

func (m *mockReportRepository) UpdateReportByReportID(c context.Context, reportID string, data *model.ReportList) error {
	copied := *data
	m.updatedReport = &copied
	return nil
}

func (m *mockReportRepository) StoreReportList(c context.Context, data *model.ReportList) error {
	copied := *data
	m.storedReport = &copied
	return nil
}

func (m *mockReportRepository) CountDownloadSalesOrderInProgress(custID string) (int64, error) {
	return 0, nil
}

func (m *mockReportRepository) CountDownloadSalesOrderByDate(custID, exportDate string) int64 {
	return 1
}

func TestGenerateDownloadSalesOrderExcel_SalesmanHeader_AllSalesmenWhenNoFilter(t *testing.T) {
	soRepo := &mockSoRepository{
		findDownloadDataPoFn: func(filter entity.SoDownloadQueryFilter) ([]model.SoDownloadPo, error) {
			employeeName := "John Doe"
			salesmanCode := "EMP062"
			return []model.SoDownloadPo{{SoNo: "SO-1", SalesmanCode: &salesmanCode, EmployeeName: &employeeName}}, nil
		},
		findDownloadDataSoFn: func(filter entity.SoDownloadQueryFilter) ([]model.SoDownloadSo, error) {
			return []model.SoDownloadSo{}, nil
		},
		findDownloadDataFinalFn: func(filter entity.SoDownloadQueryFilter) ([]model.SoDownloadFinal, error) {
			return []model.SoDownloadFinal{}, nil
		},
		findDownloadQtySummaryDataFn: func(filter entity.SoDownloadQueryFilter) ([]model.SoDownloadQtySummary, error) {
			return []model.SoDownloadQtySummary{}, nil
		},
	}
	reportRepo := &mockReportRepository{}
	service := &soServiceImpl{SoRepository: soRepo, ReportRepository: reportRepo}

	service.generateDownloadSalesOrderExcel(entity.SoDownloadQueryFilter{ReportID: "r1", StartDate: 1704067200, EndDate: 1704153600, SalesmanId: []int64{}})

	if reportRepo.updatedReport == nil {
		t.Fatalf("expected report update")
	}
	verifySalesmanHeaderValue(t, reportRepo.updatedReport.FileBase64, "All Salesmen")
}

func TestGenerateDownloadSalesOrderExcel_SalesmanHeader_MultipleSalesmenWhenFilterHasMany(t *testing.T) {
	soRepo := &mockSoRepository{
		findDownloadDataPoFn: func(filter entity.SoDownloadQueryFilter) ([]model.SoDownloadPo, error) {
			employeeName := "John Doe"
			salesmanCode := "EMP062"
			return []model.SoDownloadPo{{SoNo: "SO-1", SalesmanCode: &salesmanCode, EmployeeName: &employeeName}}, nil
		},
		findDownloadDataSoFn: func(filter entity.SoDownloadQueryFilter) ([]model.SoDownloadSo, error) {
			return []model.SoDownloadSo{}, nil
		},
		findDownloadDataFinalFn: func(filter entity.SoDownloadQueryFilter) ([]model.SoDownloadFinal, error) {
			return []model.SoDownloadFinal{}, nil
		},
		findDownloadQtySummaryDataFn: func(filter entity.SoDownloadQueryFilter) ([]model.SoDownloadQtySummary, error) {
			return []model.SoDownloadQtySummary{}, nil
		},
	}
	reportRepo := &mockReportRepository{}
	service := &soServiceImpl{SoRepository: soRepo, ReportRepository: reportRepo}

	service.generateDownloadSalesOrderExcel(entity.SoDownloadQueryFilter{ReportID: "r2", StartDate: 1704067200, EndDate: 1704153600, SalesmanId: []int64{62, 204}})

	if reportRepo.updatedReport == nil {
		t.Fatalf("expected report update")
	}
	verifySalesmanHeaderValue(t, reportRepo.updatedReport.FileBase64, "Multiple Salesmen")
}

func TestGenerateDownloadSalesOrderExcel_SalesmanHeader_SingleSalesmanUsesDataRow(t *testing.T) {
	soRepo := &mockSoRepository{
		findDownloadDataPoFn: func(filter entity.SoDownloadQueryFilter) ([]model.SoDownloadPo, error) {
			employeeName := "John Doe"
			salesmanCode := "EMP062"
			return []model.SoDownloadPo{{SoNo: "SO-1", SalesmanCode: &salesmanCode, EmployeeName: &employeeName}}, nil
		},
		findDownloadDataSoFn: func(filter entity.SoDownloadQueryFilter) ([]model.SoDownloadSo, error) {
			return []model.SoDownloadSo{}, nil
		},
		findDownloadDataFinalFn: func(filter entity.SoDownloadQueryFilter) ([]model.SoDownloadFinal, error) {
			return []model.SoDownloadFinal{}, nil
		},
		findDownloadQtySummaryDataFn: func(filter entity.SoDownloadQueryFilter) ([]model.SoDownloadQtySummary, error) {
			return []model.SoDownloadQtySummary{}, nil
		},
	}
	reportRepo := &mockReportRepository{}
	service := &soServiceImpl{SoRepository: soRepo, ReportRepository: reportRepo}

	service.generateDownloadSalesOrderExcel(entity.SoDownloadQueryFilter{ReportID: "r3", StartDate: 1704067200, EndDate: 1704153600, SalesmanId: []int64{62}})

	if reportRepo.updatedReport == nil {
		t.Fatalf("expected report update")
	}
	verifySalesmanHeaderValue(t, reportRepo.updatedReport.FileBase64, "EMP062 - John Doe")
}

func TestGenerateDownloadSalesOrderExcel_DateRangeUsesUTCCalendarDate(t *testing.T) {
	soRepo := &mockSoRepository{
		findDownloadDataPoFn: func(filter entity.SoDownloadQueryFilter) ([]model.SoDownloadPo, error) {
			return []model.SoDownloadPo{}, nil
		},
		findDownloadDataSoFn: func(filter entity.SoDownloadQueryFilter) ([]model.SoDownloadSo, error) {
			return []model.SoDownloadSo{}, nil
		},
		findDownloadDataFinalFn: func(filter entity.SoDownloadQueryFilter) ([]model.SoDownloadFinal, error) {
			return []model.SoDownloadFinal{}, nil
		},
		findDownloadQtySummaryDataFn: func(filter entity.SoDownloadQueryFilter) ([]model.SoDownloadQtySummary, error) {
			return []model.SoDownloadQtySummary{}, nil
		},
	}
	reportRepo := &mockReportRepository{}
	service := &soServiceImpl{SoRepository: soRepo, ReportRepository: reportRepo}

	service.generateDownloadSalesOrderExcel(entity.SoDownloadQueryFilter{
		ReportID:  "range-utc",
		StartDate: time.Date(2026, time.April, 15, 0, 0, 0, 0, time.UTC).Unix(),
		EndDate:   time.Date(2026, time.April, 15, 23, 59, 59, 0, time.UTC).Unix(),
	})

	if reportRepo.updatedReport == nil {
		t.Fatalf("expected report update")
	}

	verifyDateRangeHeaderValue(t, reportRepo.updatedReport.FileBase64, "15 April 2026 - 15 April 2026")
}

func TestDownload_StoresUTCDateRangeWithoutTimezoneShift(t *testing.T) {
	soRepo := &mockSoRepository{
		findDownloadDataPoFn: func(filter entity.SoDownloadQueryFilter) ([]model.SoDownloadPo, error) {
			return []model.SoDownloadPo{}, nil
		},
		findDownloadDataSoFn: func(filter entity.SoDownloadQueryFilter) ([]model.SoDownloadSo, error) {
			return []model.SoDownloadSo{}, nil
		},
		findDownloadDataFinalFn: func(filter entity.SoDownloadQueryFilter) ([]model.SoDownloadFinal, error) {
			return []model.SoDownloadFinal{}, nil
		},
		findDownloadQtySummaryDataFn: func(filter entity.SoDownloadQueryFilter) ([]model.SoDownloadQtySummary, error) {
			return []model.SoDownloadQtySummary{}, nil
		},
	}
	reportRepo := &mockReportRepository{}
	service := &soServiceImpl{SoRepository: soRepo, ReportRepository: reportRepo}

	response, err := service.Download(entity.SoDownloadQueryFilter{
		CustId:     "1100000001",
		ExportBy:   "QA User",
		StartDate:  time.Date(2026, time.April, 15, 0, 0, 0, 0, time.UTC).Unix(),
		EndDate:    time.Date(2026, time.April, 15, 23, 59, 59, 0, time.UTC).Unix(),
		SalesmanId: []int64{360},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if response.StartDate != "2026-04-15" {
		t.Fatalf("unexpected response start date: got=%q", response.StartDate)
	}
	if response.EndDate != "2026-04-15" {
		t.Fatalf("unexpected response end date: got=%q", response.EndDate)
	}
	if reportRepo.storedReport == nil {
		t.Fatalf("expected stored report")
	}
	if got := reportRepo.storedReport.StartDate.Format("2006-01-02"); got != "2026-04-15" {
		t.Fatalf("unexpected stored start date: got=%q", got)
	}
	if got := reportRepo.storedReport.EndDate.Format("2006-01-02"); got != "2026-04-15" {
		t.Fatalf("unexpected stored end date: got=%q", got)
	}
}

func TestFormatDownloadAmount_UsesIndonesianThousandsSeparator(t *testing.T) {
	tests := []struct {
		name     string
		input    *float64
		expected string
	}{
		{name: "nil", input: nil, expected: "0"},
		{name: "zero", input: ptrFloat64(0), expected: "0"},
		{name: "eleven million", input: ptrFloat64(11000000), expected: "11.000.000"},
		{name: "twelve million", input: ptrFloat64(12000000), expected: "12.000.000"},
		{name: "one point one million", input: ptrFloat64(1100000), expected: "1.100.000"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := formatDownloadAmount(tt.input); got != tt.expected {
				t.Fatalf("unexpected formatted amount: got=%q expected=%q", got, tt.expected)
			}
		})
	}
}

func ptrFloat64(v float64) *float64 { return &v }

func TestMapPoToEntity_SX1879MapsFinancialFieldsCorrectly(t *testing.T) {
	service := &soServiceImpl{}

	qty1 := 1.0
	sellPrice1 := 12000000.0
	discValueFinal := 0.0
	vatValueFinal := 1100000.0
	vatPct := 11.0

	rows := service.mapPoToEntity([]model.SoDownloadPo{{
		QtyPo1:         &qty1,
		SellPricePo1:   &sellPrice1,
		DiscValueFinal: &discValueFinal,
		VatValueFinal:  &vatValueFinal,
		Vat:            &vatPct,
	}})

	if len(rows) != 1 {
		t.Fatalf("expected one row, got %d", len(rows))
	}

	row := rows[0]
	if row.Discount == nil || *row.Discount != 0 {
		t.Fatalf("unexpected discount value: %+v", row.Discount)
	}
	if row.NetSales == nil || *row.NetSales != 12000000 {
		t.Fatalf("unexpected net sales value: %+v", row.NetSales)
	}
	if row.Vat == nil || *row.Vat != 1100000 {
		t.Fatalf("unexpected vat value: %+v", row.Vat)
	}
	if row.Gross == nil || *row.Gross != 12000000 {
		t.Fatalf("unexpected gross value: %+v", row.Gross)
	}
}

func TestMapSoToEntity_SX1879MapsFinancialFieldsCorrectly(t *testing.T) {
	service := &soServiceImpl{}

	qty1 := 1.0
	sellPrice1 := 12000000.0
	discValueFinal := 0.0
	vatValueFinal := 1100000.0
	vatPct := 11.0

	rows := service.mapSoToEntity([]model.SoDownloadSo{{
		Qty1:           &qty1,
		SellPrice1:     &sellPrice1,
		DiscValueFinal: &discValueFinal,
		VatValueFinal:  &vatValueFinal,
		Vat:            &vatPct,
	}})

	if len(rows) != 1 {
		t.Fatalf("expected one row, got %d", len(rows))
	}

	row := rows[0]
	if row.Discount == nil || *row.Discount != 0 {
		t.Fatalf("unexpected discount value: %+v", row.Discount)
	}
	if row.NetSales == nil || *row.NetSales != 12000000 {
		t.Fatalf("unexpected net sales value: %+v", row.NetSales)
	}
	if row.Vat == nil || *row.Vat != 1100000 {
		t.Fatalf("unexpected vat value: %+v", row.Vat)
	}
	if row.Gross == nil || *row.Gross != 12000000 {
		t.Fatalf("unexpected gross value: %+v", row.Gross)
	}
}

func TestMapFinalToEntity_SX1879MapsFinancialFieldsCorrectly(t *testing.T) {
	service := &soServiceImpl{}

	qty1Final := 1.0
	sellPriceFinal1 := 12000000.0
	discValueFinal := 0.0
	vatValueFinal := 1100000.0
	vatPct := 11.0

	rows := service.mapFinalToEntity([]model.SoDownloadFinal{{
		Qty1Final:       &qty1Final,
		SellPriceFinal1: &sellPriceFinal1,
		DiscValueFinal:  &discValueFinal,
		VatValueFinal:   &vatValueFinal,
		Vat:             &vatPct,
	}})

	if len(rows) != 1 {
		t.Fatalf("expected one row, got %d", len(rows))
	}

	row := rows[0]
	if row.Discount == nil || *row.Discount != 0 {
		t.Fatalf("unexpected discount value: %+v", row.Discount)
	}
	if row.NetSales == nil || *row.NetSales != 12000000 {
		t.Fatalf("unexpected net sales value: %+v", row.NetSales)
	}
	if row.Vat == nil || *row.Vat != 1100000 {
		t.Fatalf("unexpected vat value: %+v", row.Vat)
	}
	if row.Gross == nil || *row.Gross != 12000000 {
		t.Fatalf("unexpected gross value: %+v", row.Gross)
	}
}

func TestMapDownloadEntities_KeepsUnitPriceOrder(t *testing.T) {
	service := &soServiceImpl{}

	systemSmallest := 100.0
	systemMiddle := 200.0
	systemLargest := 300.0
	finalSmallest := 110.0
	finalMiddle := 220.0
	finalLargest := 330.0

	poRows := service.mapPoToEntity([]model.SoDownloadPo{{
		SellPriceSystem1: &systemSmallest,
		SellPriceSystem2: &systemMiddle,
		SellPriceSystem3: &systemLargest,
		SellPricePo1:     &finalSmallest,
		SellPricePo2:     &finalMiddle,
		SellPricePo3:     &finalLargest,
	}})
	if len(poRows) != 1 {
		t.Fatalf("expected one PO row, got %d", len(poRows))
	}
	assertDownloadPriceOrder(t, poRows[0].LargestSellingPrice, poRows[0].MiddleSellingPrice, poRows[0].SmallestSellingPrice, systemLargest, systemMiddle, systemSmallest)
	assertDownloadPriceOrder(t, poRows[0].FinalLargestSellingPrice, poRows[0].FinalMiddleSellingPrice, poRows[0].FinalSmallestSellingPrice, finalLargest, finalMiddle, finalSmallest)

	soRows := service.mapSoToEntity([]model.SoDownloadSo{{
		SellPriceSystem1: &systemSmallest,
		SellPriceSystem2: &systemMiddle,
		SellPriceSystem3: &systemLargest,
		SellPrice1:       &finalSmallest,
		SellPrice2:       &finalMiddle,
		SellPrice3:       &finalLargest,
	}})
	if len(soRows) != 1 {
		t.Fatalf("expected one SO row, got %d", len(soRows))
	}
	assertDownloadPriceOrder(t, soRows[0].LargestSellingPrice, soRows[0].MiddleSellingPrice, soRows[0].SmallestSellingPrice, systemLargest, systemMiddle, systemSmallest)
	assertDownloadPriceOrder(t, soRows[0].FinalLargestSellingPrice, soRows[0].FinalMiddleSellingPrice, soRows[0].FinalSmallestSellingPrice, finalLargest, finalMiddle, finalSmallest)

	finalRows := service.mapFinalToEntity([]model.SoDownloadFinal{{
		SellPriceSystem1: &systemSmallest,
		SellPriceSystem2: &systemMiddle,
		SellPriceSystem3: &systemLargest,
		SellPriceFinal1:  &finalSmallest,
		SellPriceFinal2:  &finalMiddle,
		SellPriceFinal3:  &finalLargest,
	}})
	if len(finalRows) != 1 {
		t.Fatalf("expected one Final row, got %d", len(finalRows))
	}
	assertDownloadPriceOrder(t, finalRows[0].LargestSellingPrice, finalRows[0].MiddleSellingPrice, finalRows[0].SmallestSellingPrice, systemLargest, systemMiddle, systemSmallest)
	assertDownloadPriceOrder(t, finalRows[0].FinalLargestSellingPrice, finalRows[0].FinalMiddleSellingPrice, finalRows[0].FinalSmallestSellingPrice, finalLargest, finalMiddle, finalSmallest)
}

func assertDownloadPriceOrder(t *testing.T, largest, middle, smallest *float64, expectedLargest, expectedMiddle, expectedSmallest float64) {
	t.Helper()

	if largest == nil || *largest != expectedLargest {
		t.Fatalf("unexpected largest price: got=%+v expected=%v", largest, expectedLargest)
	}
	if middle == nil || *middle != expectedMiddle {
		t.Fatalf("unexpected middle price: got=%+v expected=%v", middle, expectedMiddle)
	}
	if smallest == nil || *smallest != expectedSmallest {
		t.Fatalf("unexpected smallest price: got=%+v expected=%v", smallest, expectedSmallest)
	}
}

func TestCreateSalesOrderSheet_FormatsAmountColumnsWithIndonesianSeparator(t *testing.T) {
	service := &soServiceImpl{}
	f := excelize.NewFile()
	defer f.Close()

	service.createSalesOrderSheet(f, []entity.SoDownloadSoRow{{
		SoNo:                      "SO-001",
		LargestSellingPrice:       ptrFloat64(11000000),
		MiddleSellingPrice:        ptrFloat64(12000000),
		SmallestSellingPrice:      ptrFloat64(1100000),
		FinalLargestSellingPrice:  ptrFloat64(11000000),
		FinalMiddleSellingPrice:   ptrFloat64(12000000),
		FinalSmallestSellingPrice: ptrFloat64(1100000),
		GrossSales:                ptrFloat64(12000000),
		Promotion:                 ptrFloat64(0),
		Discount:                  ptrFloat64(0),
		NetSales:                  ptrFloat64(12000000),
		Vat:                       ptrFloat64(1100000),
		Gross:                     ptrFloat64(12000000),
	}}, "20 April 2026 - 28 April 2026", "All Salesmen")

	checks := map[string]string{
		"Q4":  "11.000.000",
		"R4":  "12.000.000",
		"S4":  "1.100.000",
		"T4":  "11.000.000",
		"U4":  "12.000.000",
		"V4":  "1.100.000",
		"Z4":  "12.000.000",
		"AA4": "0",
		"AB4": "0",
		"AC4": "12.000.000",
		"AD4": "1.100.000",
		"AE4": "12.000.000",
	}

	for cell, expected := range checks {
		value, err := f.GetCellValue("Sales Order", cell)
		if err != nil {
			t.Fatalf("failed to read %s: %v", cell, err)
		}
		if value != expected {
			t.Fatalf("unexpected %s value: got=%q expected=%q", cell, value, expected)
		}
	}
}

func TestCreatePurchaseOrderSheet_DefaultsNullableAmountColumnsToZero(t *testing.T) {
	service := &soServiceImpl{}
	f := excelize.NewFile()
	defer f.Close()

	service.createPurchaseOrderSheet(f, []entity.SoDownloadPoRow{{SoNo: "SO-001"}}, "20 April 2026 - 28 April 2026", "All Salesmen")

	for _, cell := range []string{"Q4", "R4", "S4", "T4", "U4", "V4", "Z4", "AA4", "AB4", "AC4", "AD4", "AE4"} {
		value, err := f.GetCellValue("Purchase Order", cell)
		if err != nil {
			t.Fatalf("failed to read %s: %v", cell, err)
		}
		if value != "0" {
			t.Fatalf("expected %s to default to 0, got %q", cell, value)
		}
	}
}

func TestMapPoToEntity_MapsOrderAndInvoiceFields(t *testing.T) {
	service := &soServiceImpl{}

	orderNo := "ORD-001"
	poNo := "PO-001"
	invoiceNo := "INV-001"
	roDate := mustParseDate(t, "2026-03-10")
	invoiceDate := mustParseDate(t, "2026-03-11")

	rows := service.mapPoToEntity([]model.SoDownloadPo{{
		OrderNo:     &orderNo,
		PoNo:        &poNo,
		SoNo:        "SO-001",
		RoDate:      &roDate,
		InvoiceDate: &invoiceDate,
		InvoiceNo:   &invoiceNo,
	}})

	if len(rows) != 1 {
		t.Fatalf("expected one row, got %d", len(rows))
	}

	row := rows[0]
	if row.OrderNo != orderNo {
		t.Fatalf("unexpected order no: got=%q expected=%q", row.OrderNo, orderNo)
	}
	if row.PoNo != poNo {
		t.Fatalf("unexpected po no: got=%q expected=%q", row.PoNo, poNo)
	}
	if row.InvoiceNo != invoiceNo {
		t.Fatalf("unexpected invoice no: got=%q expected=%q", row.InvoiceNo, invoiceNo)
	}
	if row.OrderDate != "2026-03-10" {
		t.Fatalf("unexpected order date: got=%q", row.OrderDate)
	}
	if row.InvoiceDate != "2026-03-11" {
		t.Fatalf("unexpected invoice date: got=%q", row.InvoiceDate)
	}
}

func TestMapSoToEntity_MapsOrderAndInvoiceFields(t *testing.T) {
	service := &soServiceImpl{}

	orderNo := "ORD-002"
	poNo := "PO-002"
	invoiceNo := "INV-002"
	roDate := mustParseDate(t, "2026-03-12")
	invoiceDate := mustParseDate(t, "2026-03-13")

	rows := service.mapSoToEntity([]model.SoDownloadSo{{
		OrderNo:     &orderNo,
		PoNo:        &poNo,
		SoNo:        "SO-002",
		RoDate:      &roDate,
		InvoiceDate: &invoiceDate,
		InvoiceNo:   &invoiceNo,
	}})

	if len(rows) != 1 {
		t.Fatalf("expected one row, got %d", len(rows))
	}

	row := rows[0]
	if row.OrderNo != orderNo {
		t.Fatalf("unexpected order no: got=%q expected=%q", row.OrderNo, orderNo)
	}
	if row.PoNo != poNo {
		t.Fatalf("unexpected po no: got=%q expected=%q", row.PoNo, poNo)
	}
	if row.InvoiceNo != invoiceNo {
		t.Fatalf("unexpected invoice no: got=%q expected=%q", row.InvoiceNo, invoiceNo)
	}
	if row.OrderDate != "2026-03-12" {
		t.Fatalf("unexpected order date: got=%q", row.OrderDate)
	}
	if row.InvoiceDate != "2026-03-13" {
		t.Fatalf("unexpected invoice date: got=%q", row.InvoiceDate)
	}
}

func TestMapFinalToEntity_MapsOrderAndInvoiceFields(t *testing.T) {
	service := &soServiceImpl{}

	orderNo := "ORD-003"
	poNo := "PO-003"
	invoiceNo := "INV-003"
	roDate := mustParseDate(t, "2026-03-14")
	invoiceDate := mustParseDate(t, "2026-03-15")

	rows := service.mapFinalToEntity([]model.SoDownloadFinal{{
		OrderNo:     &orderNo,
		PoNo:        &poNo,
		SoNo:        "SO-003",
		RoDate:      &roDate,
		InvoiceDate: &invoiceDate,
		InvoiceNo:   &invoiceNo,
	}})

	if len(rows) != 1 {
		t.Fatalf("expected one row, got %d", len(rows))
	}

	row := rows[0]
	if row.OrderNo != orderNo {
		t.Fatalf("unexpected order no: got=%q expected=%q", row.OrderNo, orderNo)
	}
	if row.PoNo != poNo {
		t.Fatalf("unexpected po no: got=%q expected=%q", row.PoNo, poNo)
	}
	if row.InvoiceNo != invoiceNo {
		t.Fatalf("unexpected invoice no: got=%q expected=%q", row.InvoiceNo, invoiceNo)
	}
	if row.OrderDate != "2026-03-14" {
		t.Fatalf("unexpected order date: got=%q", row.OrderDate)
	}
	if row.InvoiceDate != "2026-03-15" {
		t.Fatalf("unexpected invoice date: got=%q", row.InvoiceDate)
	}
}

func TestMapQtySummaryToEntity_MapsOrderAndInvoiceFields(t *testing.T) {
	service := &soServiceImpl{}

	orderNo := "ORD-004"
	poNo := "PO-004"
	invoiceNo := "INV-004"
	roDate := mustParseDate(t, "2026-03-16")
	invoiceDate := mustParseDate(t, "2026-03-17")

	rows := service.mapQtySummaryToEntity([]model.SoDownloadQtySummary{{
		OrderNo:     &orderNo,
		PoNo:        &poNo,
		SoNo:        "SO-004",
		RoDate:      &roDate,
		InvoiceDate: &invoiceDate,
		InvoiceNo:   &invoiceNo,
	}})

	if len(rows) != 1 {
		t.Fatalf("expected one row, got %d", len(rows))
	}

	row := rows[0]
	if row.OrderNo != orderNo {
		t.Fatalf("unexpected order no: got=%q expected=%q", row.OrderNo, orderNo)
	}
	if row.PoNo != poNo {
		t.Fatalf("unexpected po no: got=%q expected=%q", row.PoNo, poNo)
	}
	if row.InvoiceNo != invoiceNo {
		t.Fatalf("unexpected invoice no: got=%q expected=%q", row.InvoiceNo, invoiceNo)
	}
	if row.OrderDate != "2026-03-16" {
		t.Fatalf("unexpected order date: got=%q", row.OrderDate)
	}
	if row.InvoiceDate != "2026-03-17" {
		t.Fatalf("unexpected invoice date: got=%q", row.InvoiceDate)
	}
}

func TestResolveDownloadPONumber_UsesPoNoWhenAvailableAndFallsBackToOrderNo(t *testing.T) {
	poNo := "PO-100"
	orderNo := "ORD-100"
	emptyPoNo := ""

	tests := []struct {
		name     string
		poNo     *string
		orderNo  *string
		expected string
	}{
		{
			name:     "uses po number",
			poNo:     &poNo,
			orderNo:  &orderNo,
			expected: poNo,
		},
		{
			name:     "falls back when po number empty",
			poNo:     &emptyPoNo,
			orderNo:  &orderNo,
			expected: orderNo,
		},
		{
			name:     "falls back when po number nil",
			poNo:     nil,
			orderNo:  &orderNo,
			expected: orderNo,
		},
		{
			name:     "returns empty when both nil",
			poNo:     nil,
			orderNo:  nil,
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := resolveDownloadPONumber(tt.poNo, tt.orderNo)
			if actual != tt.expected {
				t.Fatalf("unexpected po number: got=%q expected=%q", actual, tt.expected)
			}
		})
	}
}

func TestFilterDownloadDataPoWithPONumber_AllowsOrderNoFallbackForBlankPoNumber(t *testing.T) {
	validPoNo := "PO-001"
	emptyPoNo := ""
	spacePoNo := "   "
	orderNo := "ORD-001"
	blankOrderNo := "   "

	rows := filterDownloadDataPoWithPONumber([]model.SoDownloadPo{
		{PoNo: nil, OrderNo: &orderNo, SoNo: "SO-NIL"},
		{PoNo: &emptyPoNo, OrderNo: &orderNo, SoNo: "SO-EMPTY"},
		{PoNo: &spacePoNo, OrderNo: &orderNo, SoNo: "SO-SPACE"},
		{PoNo: &validPoNo, OrderNo: &orderNo, SoNo: "SO-VALID"},
		{PoNo: nil, OrderNo: nil, SoNo: "SO-MISSING"},
		{PoNo: &spacePoNo, OrderNo: &blankOrderNo, SoNo: "SO-BLANK"},
	})

	if len(rows) != 4 {
		t.Fatalf("expected four valid PO rows, got %d", len(rows))
	}

	gotSoNos := make([]string, 0, len(rows))
	for _, row := range rows {
		gotSoNos = append(gotSoNos, row.SoNo)
	}
	expectedSoNos := []string{"SO-NIL", "SO-EMPTY", "SO-SPACE", "SO-VALID"}
	for i, expected := range expectedSoNos {
		if gotSoNos[i] != expected {
			t.Fatalf("unexpected included PO rows: got=%v expected=%v", gotSoNos, expectedSoNos)
		}
	}
}

func TestCreateQtySummarySheet_DefaultsPurchaseOrderQtyToZero(t *testing.T) {
	service := &soServiceImpl{}
	f := excelize.NewFile()
	defer f.Close()

	service.createQtySummarySheet(f, []entity.SoDownloadQtySummaryRow{{SoNo: "SO-001"}}, "20 April 2026 - 28 April 2026", "All Salesmen")

	for _, cell := range []string{"Q4", "R4", "S4"} {
		value, err := f.GetCellValue("QTY Summary", cell)
		if err != nil {
			t.Fatalf("failed to read %s: %v", cell, err)
		}
		if value != "0" {
			t.Fatalf("expected %s to default to 0, got %q", cell, value)
		}
	}
}

func TestCreateSalesOrderSheet_UsesEmployeeCodeValue(t *testing.T) {
	service := &soServiceImpl{}
	f := excelize.NewFile()
	defer f.Close()

	salesmanCode := "EMP001"
	service.createSalesOrderSheet(f, []entity.SoDownloadSoRow{{SoNo: "SO-001", SalesmanCode: &salesmanCode}}, "20 April 2026 - 28 April 2026", "All Salesmen")

	value, err := f.GetCellValue("Sales Order", "H4")
	if err != nil {
		t.Fatalf("failed to read Sales Order H4: %v", err)
	}
	if value != salesmanCode {
		t.Fatalf("expected salesman code %q, got %q", salesmanCode, value)
	}
}

func verifySalesmanHeaderValue(t *testing.T, fileBase64 string, expected string) {
	t.Helper()

	decoded, err := base64.StdEncoding.DecodeString(fileBase64)
	if err != nil {
		t.Fatalf("failed to decode base64: %v", err)
	}

	excelFile, err := excelize.OpenReader(bytes.NewReader(decoded))
	if err != nil {
		t.Fatalf("failed to open excel from bytes: %v", err)
	}
	defer excelFile.Close()

	for _, sheetName := range []string{"Purchase Order", "Sales Order", "Final Order", "QTY Summary"} {
		salesmanLabel, err := excelFile.GetCellValue(sheetName, "A2")
		if err != nil {
			t.Fatalf("failed to read %s A2: %v", sheetName, err)
		}
		if salesmanLabel != "Salesman" {
			t.Fatalf("expected Salesman label in %s A2, got %q", sheetName, salesmanLabel)
		}

		salesmanValue, err := excelFile.GetCellValue(sheetName, "B2")
		if err != nil {
			t.Fatalf("failed to read %s B2: %v", sheetName, err)
		}
		if salesmanValue != expected {
			t.Fatalf("unexpected salesman header in %s B2: got=%q expected=%q", sheetName, salesmanValue, expected)
		}
	}
}

func verifyDateRangeHeaderValue(t *testing.T, fileBase64 string, expected string) {
	t.Helper()

	decoded, err := base64.StdEncoding.DecodeString(fileBase64)
	if err != nil {
		t.Fatalf("failed to decode base64: %v", err)
	}

	excelFile, err := excelize.OpenReader(bytes.NewReader(decoded))
	if err != nil {
		t.Fatalf("failed to open excel from bytes: %v", err)
	}
	defer excelFile.Close()

	for _, sheetName := range []string{"Purchase Order", "Sales Order", "Final Order", "QTY Summary"} {
		dateLabel, err := excelFile.GetCellValue(sheetName, "A1")
		if err != nil {
			t.Fatalf("failed to read %s A1: %v", sheetName, err)
		}
		if dateLabel != "Order date" {
			t.Fatalf("expected Order date label in %s A1, got %q", sheetName, dateLabel)
		}

		dateValue, err := excelFile.GetCellValue(sheetName, "B1")
		if err != nil {
			t.Fatalf("failed to read %s B1: %v", sheetName, err)
		}
		if dateValue != expected {
			t.Fatalf("unexpected date header in %s B1: got=%q expected=%q", sheetName, dateValue, expected)
		}
	}
}

func mustParseDate(t *testing.T, value string) time.Time {
	t.Helper()

	parsed, err := time.Parse("2006-01-02", value)
	if err != nil {
		t.Fatalf("failed to parse date %q: %v", value, err)
	}

	return parsed
}
