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

type expenseEntryServiceMock struct {
	service.ExpenseEntryService
	capturedFilter entity.ExpenseEntryQueryFilter
	listFn         func(filter entity.ExpenseEntryQueryFilter) ([]entity.ExpenseEntryListResponse, int64, int, error)
}

func (m *expenseEntryServiceMock) List(filter entity.ExpenseEntryQueryFilter) ([]entity.ExpenseEntryListResponse, int64, int, error) {
	m.capturedFilter = filter
	if m.listFn != nil {
		return m.listFn(filter)
	}
	return []entity.ExpenseEntryListResponse{}, 0, 0, nil
}

func setupExpenseEntryListTestApp(svc service.ExpenseEntryService) *fiber.App {
	app := fiber.New()
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("requestid", "req-test-001")
		c.Locals("cust_id", "C001")
		c.Locals("parent_cust_id", "P001")
		c.Locals("user_id", int64(77))
		return c.Next()
	})

	controller := NewExpenseEntryController(svc, validation.NewValiditor())
	app.Get("/v1/expense", controller.List)
	return app
}

func TestExpenseEntryController_List_SuccessDefaultsAndCollectorParsing(t *testing.T) {
	svcMock := &expenseEntryServiceMock{
		listFn: func(filter entity.ExpenseEntryQueryFilter) ([]entity.ExpenseEntryListResponse, int64, int, error) {
			return []entity.ExpenseEntryListResponse{
				{
					ExpenseID:  10,
					DocumentNo: "EXP-001",
				},
			}, 1, 1, nil
		},
	}

	app := setupExpenseEntryListTestApp(svcMock)

	req := httptest.NewRequest("GET", "/v1/expense?collector_id=12,13&q=E20260424002", nil)
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
	if svcMock.capturedFilter.CustID != "C001" {
		t.Fatalf("expected cust_id C001, got %s", svcMock.capturedFilter.CustID)
	}
	if svcMock.capturedFilter.ParentCustID != "P001" {
		t.Fatalf("expected parent_cust_id P001, got %s", svcMock.capturedFilter.ParentCustID)
	}
	if svcMock.capturedFilter.UserID != 77 {
		t.Fatalf("expected user_id 77, got %d", svcMock.capturedFilter.UserID)
	}
	if len(svcMock.capturedFilter.CollectorIDs) != 2 || svcMock.capturedFilter.CollectorIDs[0] != 12 || svcMock.capturedFilter.CollectorIDs[1] != 13 {
		t.Fatalf("unexpected collector_ids parsed: %+v", svcMock.capturedFilter.CollectorIDs)
	}
	if svcMock.capturedFilter.Query != "E20260424002" {
		t.Fatalf("expected q to be parsed as E20260424002, got %q", svcMock.capturedFilter.Query)
	}

	var body map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		t.Fatalf("failed decode response: %v", err)
	}
	if body["message"] != "" {
		t.Fatalf("unexpected message: %v", body["message"])
	}
	if _, ok := body["paging"]; !ok {
		t.Fatalf("expected paging field in response")
	}
}

func TestExpenseEntryController_List_InvalidCollectorID(t *testing.T) {
	svcMock := &expenseEntryServiceMock{}
	app := setupExpenseEntryListTestApp(svcMock)

	req := httptest.NewRequest("GET", "/v1/expense?collector_id=abc", nil)
	res, err := app.Test(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != fiber.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", res.StatusCode)
	}
}

func TestExpenseEntryController_List_InvalidStartDateReturnsBadRequest(t *testing.T) {
	svcMock := &expenseEntryServiceMock{
		listFn: func(filter entity.ExpenseEntryQueryFilter) ([]entity.ExpenseEntryListResponse, int64, int, error) {
			return nil, 0, 0, service.ErrInvalidExpenseDateFilter
		},
	}
	app := setupExpenseEntryListTestApp(svcMock)

	req := httptest.NewRequest("GET", "/v1/expense?start_date=invalid-epoch", nil)
	res, err := app.Test(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != fiber.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", res.StatusCode)
	}
}

func TestExpenseEntryController_List_EmptyStateCompatible(t *testing.T) {
	svcMock := &expenseEntryServiceMock{}
	app := setupExpenseEntryListTestApp(svcMock)

	req := httptest.NewRequest("GET", "/v1/expense", nil)
	res, err := app.Test(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != fiber.StatusOK {
		t.Fatalf("expected status 200, got %d", res.StatusCode)
	}

	var body map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		t.Fatalf("failed decode response: %v", err)
	}
	if body["message"] != "No Data" {
		t.Fatalf("expected empty state message No Data, got %v", body["message"])
	}
	if data, ok := body["data"]; ok && data != nil {
		t.Fatalf("expected empty data to be nil, got %#v", data)
	}
}
