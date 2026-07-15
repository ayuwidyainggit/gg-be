package controller

import (
	"encoding/json"
	"master/entity"
	"master/service"
	"net/http/httptest"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/gofiber/fiber/v2"
)

func productReportTestApp(t *testing.T, svc *productServiceStub) *fiber.App {
	t.Helper()
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller failed")
	}
	t.Chdir(filepath.Dir(filepath.Dir(file)))

	app := fiber.New()
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("requestid", "req-test")
		c.Locals("cust_id", "C26002")
		c.Locals("parent_cust_id", "C26002")
		c.Locals("distributor_id", int64(0))
		c.Locals("user_id", int64(0))
		return c.Next()
	})
	ctrl := NewProductController(svc, nil)
	ctrl.Route(app)
	return app
}

type productServiceStub struct {
	service.ProductService
	capturedFilter entity.ProductReportQueryFilter
	reportCalled   bool
	detailCalled   bool
}

func (s *productServiceStub) Detail(_ entity.DetailProductParams) (entity.ProductDetailResponse, error) {
	s.detailCalled = true
	return entity.ProductDetailResponse{}, nil
}

func (s *productServiceStub) ReportList(filter entity.ProductReportQueryFilter) ([]entity.ProductReportResponse, int, int, error) {
	s.capturedFilter = filter
	s.reportCalled = true
	return []entity.ProductReportResponse{}, 0, 0, nil
}

func TestProductReport_MissingCustID_Returns400(t *testing.T) {
	svc := &productServiceStub{}
	ctrl := &ProductController{ProductService: svc}
	app := fiber.New()
	app.Post("/v1/products/report", func(c *fiber.Ctx) error {
		c.Locals("requestid", "req-test")
		return ctrl.Report(c)
	})

	req := httptest.NewRequest("POST", "/v1/products/report", nil)
	res, err := app.Test(req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if res.StatusCode != fiber.StatusBadRequest {
		t.Fatalf("expected 400, got %d", res.StatusCode)
	}
	if svc.reportCalled {
		t.Fatal("service should not be called when cust_id is missing")
	}
	var resp map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&resp); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if resp["errors"] == nil {
		t.Fatal("expected errors in response")
	}
}

func TestProductReport_BlankCustID_Returns400(t *testing.T) {
	svc := &productServiceStub{}
	ctrl := &ProductController{ProductService: svc}
	app := fiber.New()
	app.Post("/v1/products/report", func(c *fiber.Ctx) error {
		c.Locals("requestid", "req-test")
		return ctrl.Report(c)
	})

	req := httptest.NewRequest("POST", "/v1/products/report?cust_id[]=", nil)
	res, err := app.Test(req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if res.StatusCode != fiber.StatusBadRequest {
		t.Fatalf("expected 400, got %d", res.StatusCode)
	}
	if svc.reportCalled {
		t.Fatal("service should not be called when cust_id is blank")
	}
}

func TestProductReport_InvalidSortBy_Returns400(t *testing.T) {
	svc := &productServiceStub{}
	ctrl := &ProductController{ProductService: svc}
	app := fiber.New()
	app.Post("/v1/products/report", func(c *fiber.Ctx) error {
		c.Locals("requestid", "req-test")
		return ctrl.Report(c)
	})

	req := httptest.NewRequest("POST", "/v1/products/report?cust_id[]=C26002&sort_by=invalid", nil)
	res, err := app.Test(req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if res.StatusCode != fiber.StatusBadRequest {
		t.Fatalf("expected 400, got %d", res.StatusCode)
	}
	if svc.reportCalled {
		t.Fatal("service should not be called when sort_by is invalid")
	}
}

func TestProductReport_InvalidSortOrder_Returns400(t *testing.T) {
	svc := &productServiceStub{}
	ctrl := &ProductController{ProductService: svc}
	app := fiber.New()
	app.Post("/v1/products/report", func(c *fiber.Ctx) error {
		c.Locals("requestid", "req-test")
		return ctrl.Report(c)
	})

	req := httptest.NewRequest("POST", "/v1/products/report?cust_id[]=C26002&sort_order=invalid", nil)
	res, err := app.Test(req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if res.StatusCode != fiber.StatusBadRequest {
		t.Fatalf("expected 400, got %d", res.StatusCode)
	}
	if svc.reportCalled {
		t.Fatal("service should not be called when sort_order is invalid")
	}
}

func TestProductReportRoute_GET_ReachesReport(t *testing.T) {
	svc := &productServiceStub{}
	app := productReportTestApp(t, svc)

	req := httptest.NewRequest("GET", "/v1/products/report?cust_id[]=C26002&page=1&limit=20", nil)
	req.Header.Set("Cust_id", "C26002")
	res, err := app.Test(req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if res.StatusCode != fiber.StatusOK {
		t.Fatalf("expected 200, got %d", res.StatusCode)
	}
	if !svc.reportCalled {
		t.Fatal("report service should be called")
	}
	if svc.detailCalled {
		t.Fatal("detail service should not be called")
	}
}

func TestProductReportRoute_GET_Detail_ReachesDetail(t *testing.T) {
	svc := &productServiceStub{}
	app := productReportTestApp(t, svc)

	req := httptest.NewRequest("GET", "/v1/products/12345", nil)
	req.Header.Set("Cust_id", "C26002")
	_, err := app.Test(req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !svc.detailCalled {
		t.Fatal("detail service should be called")
	}
	if svc.reportCalled {
		t.Fatal("report service should not be called")
	}
}

func TestProductReport_ValidRequest_CallsService(t *testing.T) {
	svc := &productServiceStub{}
	ctrl := &ProductController{ProductService: svc}
	app := fiber.New()
	app.Post("/v1/products/report", func(c *fiber.Ctx) error {
		c.Locals("requestid", "req-test")
		return ctrl.Report(c)
	})

	req := httptest.NewRequest("POST", "/v1/products/report?cust_id[]=C26002&cust_id[]=C260020001&q=ABC&page=1&limit=20&sort_by=pro_code&sort_order=desc", nil)
	res, err := app.Test(req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if res.StatusCode != fiber.StatusOK {
		t.Fatalf("expected 200, got %d", res.StatusCode)
	}
	if !svc.reportCalled {
		t.Fatal("service should be called for valid request")
	}

	// Verify captured filter
	if len(svc.capturedFilter.CustIDs) != 2 {
		t.Fatalf("expected 2 cust_ids, got %d: %v", len(svc.capturedFilter.CustIDs), svc.capturedFilter.CustIDs)
	}
	if svc.capturedFilter.CustIDs[0] != "C26002" || svc.capturedFilter.CustIDs[1] != "C260020001" {
		t.Fatalf("expected [C26002 C260020001], got %v", svc.capturedFilter.CustIDs)
	}
	if svc.capturedFilter.Query != "ABC" {
		t.Fatalf("expected q=ABC, got %s", svc.capturedFilter.Query)
	}
	if svc.capturedFilter.Page != 1 {
		t.Fatalf("expected page=1, got %d", svc.capturedFilter.Page)
	}
	if svc.capturedFilter.Limit != 20 {
		t.Fatalf("expected limit=20, got %d", svc.capturedFilter.Limit)
	}
	if svc.capturedFilter.SortBy != "pro_code" {
		t.Fatalf("expected sort_by=pro_code, got %s", svc.capturedFilter.SortBy)
	}
	if svc.capturedFilter.SortOrder != "desc" {
		t.Fatalf("expected sort_order=desc, got %s", svc.capturedFilter.SortOrder)
	}

	// Verify response envelope
	var resp map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&resp); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if resp["paging"] == nil {
		t.Fatal("expected paging in response")
	}
	paging := resp["paging"].(map[string]interface{})
	if paging["total_record"] == nil {
		t.Fatal("expected total_record in paging")
	}
	if paging["page_current"] == nil {
		t.Fatal("expected page_current in paging")
	}
	if paging["page_limit"] == nil {
		t.Fatal("expected page_limit in paging")
	}
	if paging["page_total"] == nil {
		t.Fatal("expected page_total in paging")
	}
}
