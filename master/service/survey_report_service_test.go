package service

import (
	"bytes"
	"encoding/base64"
	"master/adapter"
	"master/entity"
	"master/model"
	"testing"
	"time"

	"github.com/xuri/excelize/v2"
)

func TestNormalizeSurveyReportFilter_ShouldApplyDefaults(t *testing.T) {
	filter := entity.SurveyReportQueryFilter{}

	if err := normalizeSurveyReportFilter(&filter); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if filter.Page != 1 {
		t.Fatalf("expected default page 1, got %d", filter.Page)
	}
	if filter.Limit != 5 {
		t.Fatalf("expected default limit 5, got %d", filter.Limit)
	}
	if filter.Sort != "created_date:desc" {
		t.Fatalf("expected default sort created_date:desc, got %s", filter.Sort)
	}
}

func TestNormalizeSurveyReportFilter_ShouldRejectInvalidDateRange(t *testing.T) {
	start := time.Date(2026, 4, 24, 0, 0, 0, 0, time.UTC)
	end := time.Date(2026, 4, 23, 0, 0, 0, 0, time.UTC)
	filter := entity.SurveyReportQueryFilter{StartDate: &start, EndDate: &end}

	err := normalizeSurveyReportFilter(&filter)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if err != ErrSurveyReportInvalidDateRange {
		t.Fatalf("expected ErrSurveyReportInvalidDateRange, got %v", err)
	}
}

func TestGenerateSurveyReportExcel_ShouldWriteAttachmentHyperlinks(t *testing.T) {
	now := time.Now().UTC()
	rows := []model.SurveyReportExportRow{
		{
			SurveyDate:      &now,
			SurveyTitle:     "Survey Outlet",
			AreaCode:        "AR001",
			AreaName:        "Semarang",
			DistributorCode: "DIST001",
			DistributorName: "PT Distributor Jaya",
			EmpCode:         "EMP001",
			EmpName:         "Budi",
			OutletCode:      "OUT001",
			OutletName:      "Outlet A",
			Question:        "Bagaimana kualitas produk kami?",
			Answer:          "Baik",
			Attachment1:     "survey/C26004/file a.jpg",
			Attachment3:     "https://bucket.example.com/survey/C26004/file-b.jpg",
		},
	}

	fileBase64, err := generateSurveyReportExcel(rows, &adapter.ObsAdapterImpl{FileBaseUrl: "https://bucket.example.com"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if fileBase64 == "" {
		t.Fatal("expected non-empty base64")
	}

	decoded, err := base64.StdEncoding.DecodeString(fileBase64)
	if err != nil {
		t.Fatalf("failed to decode base64: %v", err)
	}
	if len(decoded) == 0 {
		t.Fatal("expected decoded workbook bytes")
	}
	if string(decoded[:2]) != "PK" {
		t.Fatalf("expected xlsx zip signature PK, got %q", string(decoded[:2]))
	}

	f, err := excelize.OpenReader(bytes.NewReader(decoded))
	if err != nil {
		t.Fatalf("failed to open workbook: %v", err)
	}
	defer f.Close()

	if got, err := f.GetCellValue("Survey Report", "M2"); err != nil || got != "survey/C26004/file a.jpg" {
		t.Fatalf("unexpected M2 value: %q", got)
	}
	if ok, link, err := f.GetCellHyperLink("Survey Report", "M2"); err != nil || !ok || link != "https://bucket.example.com/survey/C26004/file%20a.jpg" {
		t.Fatalf("unexpected M2 hyperlink: %v / %q / %v", ok, link, err)
	}
	if got, err := f.GetCellValue("Survey Report", "N2"); err != nil || got != "" {
		t.Fatalf("expected empty N2, got %q", got)
	}
	if ok, _, err := f.GetCellHyperLink("Survey Report", "N2"); err == nil && ok {
		t.Fatal("expected no hyperlink on N2")
	}
	if got, err := f.GetCellValue("Survey Report", "O2"); err != nil || got != "file-b.jpg" {
		t.Fatalf("unexpected O2 value: %q", got)
	}
	if ok, link, err := f.GetCellHyperLink("Survey Report", "O2"); err != nil || !ok || link != "https://bucket.example.com/survey/C26004/file-b.jpg" {
		t.Fatalf("unexpected O2 hyperlink: %v / %q / %v", ok, link, err)
	}
}
