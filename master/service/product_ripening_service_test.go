package service

import (
	"bytes"
	"encoding/csv"
	"master/model"
	"testing"
	"time"

	"github.com/xuri/excelize/v2"
)

func TestCreateProductRipeningExportCSVFormatsAuditTimesInJakarta(t *testing.T) {
	rows := []model.ProductRipeningPlanListRow{productRipeningExportTestRow()}

	buf, err := createProductRipeningExportCSV(rows)
	if err != nil {
		t.Fatalf("createProductRipeningExportCSV() error = %v", err)
	}

	records, err := csv.NewReader(bytes.NewReader(buf.Bytes())).ReadAll()
	if err != nil {
		t.Fatalf("ReadAll() error = %v", err)
	}
	if got, want := records[1][11], "2026-05-24T01:30:00+07:00"; got != want {
		t.Fatalf("Created At = %q, want %q", got, want)
	}
	if got, want := records[1][14], "2026-05-24T07:05:00+07:00"; got != want {
		t.Fatalf("Updated At = %q, want %q", got, want)
	}
}

func TestCreateProductRipeningExportWorkbookFormatsAuditTimesInJakarta(t *testing.T) {
	rows := []model.ProductRipeningPlanListRow{productRipeningExportTestRow()}

	buf, err := createProductRipeningExportWorkbook(rows)
	if err != nil {
		t.Fatalf("createProductRipeningExportWorkbook() error = %v", err)
	}

	f, err := excelize.OpenReader(bytes.NewReader(buf.Bytes()))
	if err != nil {
		t.Fatalf("OpenReader() error = %v", err)
	}
	defer f.Close()

	if got, want := cellValue(t, f, "L2"), "2026-05-24T01:30:00+07:00"; got != want {
		t.Fatalf("Created At = %q, want %q", got, want)
	}
	if got, want := cellValue(t, f, "O2"), "2026-05-24T07:05:00+07:00"; got != want {
		t.Fatalf("Updated At = %q, want %q", got, want)
	}
}

func productRipeningExportTestRow() model.ProductRipeningPlanListRow {
	updatedBy := int64(2)
	updatedByName := "updater"
	updatedAt := time.Date(2026, 5, 24, 0, 5, 0, 0, time.UTC)

	return model.ProductRipeningPlanListRow{
		ID:              1,
		CustID:          "CUST",
		DistributorID:   10,
		DistributorCode: "D001",
		DistributorName: "Distributor",
		PerYear:         2026,
		PerID:           5,
		WeekID:          21,
		WeekStart:       "2026-05-24",
		WeekEnd:         "2026-05-30",
		TotalProduct:    3,
		CreatedBy:       1,
		CreatedByName:   "creator",
		CreatedAt:       time.Date(2026, 5, 23, 18, 30, 0, 0, time.UTC),
		UpdatedBy:       &updatedBy,
		UpdatedByName:   &updatedByName,
		UpdatedAt:       &updatedAt,
	}
}

func cellValue(t *testing.T, f *excelize.File, axis string) string {
	t.Helper()
	value, err := f.GetCellValue("Product Ripening", axis)
	if err != nil {
		t.Fatalf("GetCellValue(%s) error = %v", axis, err)
	}
	return value
}
