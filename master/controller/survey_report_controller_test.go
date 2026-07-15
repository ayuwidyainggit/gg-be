package controller

import (
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"io"
	"master/entity"
	"master/pkg/validation"
	"master/service"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
)

type surveyReportServiceStub struct {
	listFn   func(filter entity.SurveyReportQueryFilter) ([]entity.SurveyReportListResponse, int, int, error)
	detailFn func(surveyAnswerID int64, custID string) (entity.SurveyReportDetailResponse, error)
	exportFn func(filter entity.SurveyReportQueryFilter, createdBy string) (entity.SurveyReportExportResponse, error)
}

func (s *surveyReportServiceStub) List(filter entity.SurveyReportQueryFilter) ([]entity.SurveyReportListResponse, int, int, error) {
	if s.listFn != nil {
		return s.listFn(filter)
	}
	return []entity.SurveyReportListResponse{}, 0, 1, nil
}

func (s *surveyReportServiceStub) Detail(surveyAnswerID int64, custID string) (entity.SurveyReportDetailResponse, error) {
	if s.detailFn != nil {
		return s.detailFn(surveyAnswerID, custID)
	}
	return entity.SurveyReportDetailResponse{}, nil
}

func (s *surveyReportServiceStub) Export(filter entity.SurveyReportQueryFilter, createdBy string) (entity.SurveyReportExportResponse, error) {
	if s.exportFn != nil {
		return s.exportFn(filter, createdBy)
	}
	return entity.SurveyReportExportResponse{}, nil
}

func TestParseSurveyReportFilter_ShouldParseDateStringRange(t *testing.T) {
	app := fiber.New()
	var captured entity.SurveyReportQueryFilter

	app.Get("/v1/survey-report", func(c *fiber.Ctx) error {
		var err error
		captured, err = parseSurveyReportFilter(c)
		if err != nil {
			return err
		}
		return c.SendStatus(fiber.StatusOK)
	})

	req := httptest.NewRequest("GET", "/v1/survey-report?start_date=2026-04-24&end_date=2026-04-24", nil)
	res, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.StatusCode != fiber.StatusOK {
		t.Fatalf("expected status %d, got %d", fiber.StatusOK, res.StatusCode)
	}

	expectedDate := time.Date(2026, 4, 24, 0, 0, 0, 0, time.UTC)
	if captured.StartDate == nil || !captured.StartDate.Equal(expectedDate) {
		t.Fatalf("unexpected start_date: %v", captured.StartDate)
	}
	if captured.EndDate == nil || !captured.EndDate.Equal(expectedDate) {
		t.Fatalf("unexpected end_date: %v", captured.EndDate)
	}
}

func TestSurveyReportController_List_ShouldParseArrayFiltersAndCustID(t *testing.T) {
	app := fiber.New()
	v := validation.NewValidator()

	serviceStub := &surveyReportServiceStub{}
	controller := NewSurveyReportController(serviceStub, v)

	var captured entity.SurveyReportQueryFilter
	serviceStub.listFn = func(filter entity.SurveyReportQueryFilter) ([]entity.SurveyReportListResponse, int, int, error) {
		captured = filter
		return []entity.SurveyReportListResponse{}, 0, 1, nil
	}

	app.Get("/v1/survey-report", func(c *fiber.Ctx) error {
		c.Locals("requestid", "req-1")
		c.Locals("cust_id", "C260020001")
		return controller.List(c)
	})

	req := httptest.NewRequest("GET", "/v1/survey-report?survey_id=1,2&survey_title=Survey%20Outlet,Survey%20NOO&area_id[]=5&area_id[]=8&page=2&limit=7&q=test", nil)
	res, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if res.StatusCode != fiber.StatusOK {
		t.Fatalf("expected status %d, got %d", fiber.StatusOK, res.StatusCode)
	}

	if captured.CustID != "C260020001" {
		t.Fatalf("expected cust_id C260020001, got %s", captured.CustID)
	}
	if len(captured.SurveyID) != 2 || captured.SurveyID[0] != 1 || captured.SurveyID[1] != 2 {
		t.Fatalf("unexpected survey_id parsed: %+v", captured.SurveyID)
	}
	if len(captured.AreaID) != 2 || captured.AreaID[0] != 5 || captured.AreaID[1] != 8 {
		t.Fatalf("unexpected area_id parsed: %+v", captured.AreaID)
	}
	if len(captured.SurveyTitle) != 2 || captured.SurveyTitle[0] != "Survey Outlet" || captured.SurveyTitle[1] != "Survey NOO" {
		t.Fatalf("unexpected survey_title parsed: %+v", captured.SurveyTitle)
	}
	if captured.Page != 2 || captured.Limit != 7 || captured.Query != "test" {
		t.Fatalf("unexpected filter payload: %+v", captured)
	}
	if captured.Sort != "created_date:desc" {
		t.Fatalf("expected controller default sort created_date:desc, got %q", captured.Sort)
	}
}

func TestSurveyReportController_Detail_ShouldReturnNotFound(t *testing.T) {
	app := fiber.New()
	v := validation.NewValidator()

	serviceStub := &surveyReportServiceStub{
		detailFn: func(surveyAnswerID int64, custID string) (entity.SurveyReportDetailResponse, error) {
			return entity.SurveyReportDetailResponse{}, sql.ErrNoRows
		},
	}
	controller := NewSurveyReportController(serviceStub, v)

	app.Get("/v1/survey-report/:survey_answer_id", func(c *fiber.Ctx) error {
		c.Locals("requestid", "req-2")
		c.Locals("cust_id", "C260020001")
		return controller.Detail(c)
	})

	req := httptest.NewRequest("GET", "/v1/survey-report/12", nil)
	res, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if res.StatusCode != fiber.StatusNotFound {
		t.Fatalf("expected status %d, got %d", fiber.StatusNotFound, res.StatusCode)
	}
}

func TestSurveyReportController_List_ShouldReturnDefaultPagingMetadata(t *testing.T) {
	app := fiber.New()
	v := validation.NewValidator()

	serviceStub := &surveyReportServiceStub{}
	controller := NewSurveyReportController(serviceStub, v)

	var captured entity.SurveyReportQueryFilter
	serviceStub.listFn = func(filter entity.SurveyReportQueryFilter) ([]entity.SurveyReportListResponse, int, int, error) {
		captured = filter
		return []entity.SurveyReportListResponse{}, 39, 8, nil
	}

	app.Get("/v1/survey-report", func(c *fiber.Ctx) error {
		c.Locals("requestid", "req-default")
		c.Locals("cust_id", "C260040001")
		return controller.List(c)
	})

	req := httptest.NewRequest("GET", "/v1/survey-report", nil)
	res, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if res.StatusCode != fiber.StatusOK {
		t.Fatalf("expected status %d, got %d", fiber.StatusOK, res.StatusCode)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("failed to read body: %v", err)
	}

	var payload struct {
		Paging entity.Pagination `json:"paging"`
	}
	if err := json.Unmarshal(body, &payload); err != nil {
		t.Fatalf("failed to decode payload: %v", err)
	}

	if captured.Page != 1 || captured.Limit != 5 || captured.Sort != "created_date:desc" {
		t.Fatalf("expected normalized defaults in controller, got %+v", captured)
	}
	if payload.Paging.PageCurrent != 1 || payload.Paging.PageLimit != 5 {
		t.Fatalf("expected paging defaults page_current=1 page_limit=5, got %+v", payload.Paging)
	}
	if payload.Paging.TotalRecord != 39 || payload.Paging.PageTotal != 8 {
		t.Fatalf("unexpected paging totals: %+v", payload.Paging)
	}
}

func TestSurveyReportController_Export_ShouldReturnValidationError(t *testing.T) {
	app := fiber.New()
	v := validation.NewValidator()

	serviceStub := &surveyReportServiceStub{
		exportFn: func(filter entity.SurveyReportQueryFilter, createdBy string) (entity.SurveyReportExportResponse, error) {
			return entity.SurveyReportExportResponse{}, service.ErrSurveyReportInvalidDateRange
		},
	}
	controller := NewSurveyReportController(serviceStub, v)

	app.Get("/v1/survey-report/export", func(c *fiber.Ctx) error {
		c.Locals("requestid", "req-3")
		c.Locals("cust_id", "C260020001")
		c.Locals("user_fullname", "Ujang")
		return controller.Export(c)
	})

	req := httptest.NewRequest("GET", "/v1/survey-report/export?start_date=200&end_date=100", nil)
	res, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if res.StatusCode != fiber.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", fiber.StatusBadRequest, res.StatusCode)
	}

	var payload struct {
		Message string `json:"message"`
	}
	if err := json.NewDecoder(res.Body).Decode(&payload); err != nil {
		t.Fatalf("failed to decode payload: %v", err)
	}
	if payload.Message != service.ErrSurveyReportInvalidDateRange.Error() {
		t.Fatalf("expected message %q, got %q", service.ErrSurveyReportInvalidDateRange.Error(), payload.Message)
	}
}

func TestSurveyReportController_Export_ShouldReturnXLSXWhenDownloadTrue(t *testing.T) {
	app := fiber.New()
	v := validation.NewValidator()

	serviceStub := &surveyReportServiceStub{
		exportFn: func(filter entity.SurveyReportQueryFilter, createdBy string) (entity.SurveyReportExportResponse, error) {
			return entity.SurveyReportExportResponse{
				ReportName: "DownloadSurveyReport-270426-001",
				FileBase64: base64.StdEncoding.EncodeToString([]byte("PK-test-workbook")),
			}, nil
		},
	}
	controller := NewSurveyReportController(serviceStub, v)

	app.Get("/v1/survey-report/export", func(c *fiber.Ctx) error {
		c.Locals("requestid", "req-download")
		c.Locals("cust_id", "C260020001")
		c.Locals("user_fullname", "Ujang")
		return controller.Export(c)
	})

	req := httptest.NewRequest("GET", "/v1/survey-report/export?download=true", nil)
	res, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if res.StatusCode != fiber.StatusOK {
		t.Fatalf("expected status %d, got %d", fiber.StatusOK, res.StatusCode)
	}
	if res.Header.Get("Content-Type") != "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet" {
		t.Fatalf("unexpected content type: %s", res.Header.Get("Content-Type"))
	}
	if res.Header.Get("Content-Disposition") != `attachment; filename="DownloadSurveyReport-270426-001.xlsx"` {
		t.Fatalf("unexpected content disposition: %s", res.Header.Get("Content-Disposition"))
	}
}

func TestSurveyReportController_RouteOrder_ShouldNotTreatExportAsDetail(t *testing.T) {
	app := fiber.New()
	v := validation.NewValidator()

	exportCalled := false
	detailCalled := false
	serviceStub := &surveyReportServiceStub{
		exportFn: func(filter entity.SurveyReportQueryFilter, createdBy string) (entity.SurveyReportExportResponse, error) {
			exportCalled = true
			return entity.SurveyReportExportResponse{}, nil
		},
		detailFn: func(surveyAnswerID int64, custID string) (entity.SurveyReportDetailResponse, error) {
			detailCalled = true
			return entity.SurveyReportDetailResponse{}, nil
		},
	}
	controller := NewSurveyReportController(serviceStub, v)

	app.Get("/v1/survey-report/export", func(c *fiber.Ctx) error {
		c.Locals("requestid", "req-route")
		c.Locals("cust_id", "C260020001")
		return controller.Export(c)
	})
	app.Get("/v1/survey-report/:survey_answer_id", func(c *fiber.Ctx) error {
		c.Locals("requestid", "req-route")
		c.Locals("cust_id", "C260020001")
		return controller.Detail(c)
	})

	req := httptest.NewRequest("GET", "/v1/survey-report/export", nil)
	res, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.StatusCode != fiber.StatusOK {
		t.Fatalf("expected status %d, got %d", fiber.StatusOK, res.StatusCode)
	}
	if !exportCalled {
		t.Fatal("expected export handler to be called")
	}
	if detailCalled {
		t.Fatal("did not expect detail handler to be called for /export")
	}
}
