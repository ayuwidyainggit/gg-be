package controller

import (
	"encoding/json"
	"finance/entity"
	"finance/pkg/validation"
	"finance/service"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
)

type depositServiceMock struct {
	service.DepositService
	capturedFilter entity.DepositNumberListQueryFilter
	listFn         func(filter entity.DepositNumberListQueryFilter) ([]entity.DepositNumberListItemResponse, int64, int, error)
}

func (m *depositServiceMock) ListDepositNumber(filter entity.DepositNumberListQueryFilter) ([]entity.DepositNumberListItemResponse, int64, int, error) {
	m.capturedFilter = filter
	if m.listFn != nil {
		return m.listFn(filter)
	}
	return []entity.DepositNumberListItemResponse{}, 0, 0, nil
}

func setupDepositNumberListTestApp(svc service.DepositService) *fiber.App {
	app := fiber.New()
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("requestid", "req-test-001")
		c.Locals("cust_id", "C001")
		return c.Next()
	})

	controller := NewDepositController(svc, validation.NewValiditor())
	app.Get("/v1/deposits", controller.ListDepositNumber)
	return app
}

func TestDepositController_ListDepositNumber_Success(t *testing.T) {
	svcMock := &depositServiceMock{
		listFn: func(filter entity.DepositNumberListQueryFilter) ([]entity.DepositNumberListItemResponse, int64, int, error) {
			return []entity.DepositNumberListItemResponse{
				{
					DepositNo:   "DEP-2026-001",
					CollectorID: 123,
					DepositDate: "2026-02-01T00:00:00Z",
				},
			}, 72, 4, nil
		},
	}

	app := setupDepositNumberListTestApp(svcMock)

	req := httptest.NewRequest("GET", "/v1/deposits?collector_id=123,124", nil)
	res, err := app.Test(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != fiber.StatusOK {
		t.Fatalf("expected status 200, got %d", res.StatusCode)
	}

	if svcMock.capturedFilter.Page != 1 {
		t.Fatalf("expected default page 1, got %d", svcMock.capturedFilter.Page)
	}
	if svcMock.capturedFilter.Limit != 20 {
		t.Fatalf("expected default limit 20, got %d", svcMock.capturedFilter.Limit)
	}
	if svcMock.capturedFilter.Sort != "created_date:desc" {
		t.Fatalf("expected default sort created_date:desc, got %s", svcMock.capturedFilter.Sort)
	}
	if svcMock.capturedFilter.CustId != "C001" {
		t.Fatalf("expected cust_id C001, got %s", svcMock.capturedFilter.CustId)
	}
	if len(svcMock.capturedFilter.CollectorIDs) != 2 || svcMock.capturedFilter.CollectorIDs[0] != 123 || svcMock.capturedFilter.CollectorIDs[1] != 124 {
		t.Fatalf("unexpected collector_ids parsed: %+v", svcMock.capturedFilter.CollectorIDs)
	}

	var body map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		t.Fatalf("failed decode response: %v", err)
	}

	if body["message"] != "Data berhasil ditampilkan" {
		t.Fatalf("unexpected message: %v", body["message"])
	}
	if _, ok := body["pagination"]; !ok {
		t.Fatalf("pagination is missing in response")
	}
}

func TestDepositController_ListDepositNumber_InvalidCollectorID(t *testing.T) {
	svcMock := &depositServiceMock{}
	app := setupDepositNumberListTestApp(svcMock)

	req := httptest.NewRequest("GET", "/v1/deposits?collector_id=abc", nil)
	res, err := app.Test(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != fiber.StatusUnprocessableEntity {
		t.Fatalf("expected status 422, got %d", res.StatusCode)
	}
}

func TestDepositController_ListDepositNumber_CollectorIDRequired(t *testing.T) {
	svcMock := &depositServiceMock{}
	app := setupDepositNumberListTestApp(svcMock)

	req := httptest.NewRequest("GET", "/v1/deposits", nil)
	res, err := app.Test(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != fiber.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", res.StatusCode)
	}
}
