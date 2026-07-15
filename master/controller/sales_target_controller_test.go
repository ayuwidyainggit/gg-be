package controller

import (
	"master/entity"
	"master/pkg/validation"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
)

type salesTargetServiceControllerStub struct {
	storeCalled bool
	storeReq    entity.CreateSalesTargetRequest
	storeErr    error
}

func (s *salesTargetServiceControllerStub) List(_ entity.SalesTargetQueryFilter, _ string) ([]entity.SalesTargetListResponse, int, int, error) {
	return nil, 0, 0, nil
}

func (s *salesTargetServiceControllerStub) Detail(_ int64, _ string) (entity.SalesTargetDetailResponse, error) {
	return entity.SalesTargetDetailResponse{}, nil
}

func (s *salesTargetServiceControllerStub) MonthlyDistributor(_ entity.SalesTargetMonthlyDistQuery) (entity.SalesTargetMonthlyDistResp, error) {
	return entity.SalesTargetMonthlyDistResp{}, nil
}

func (s *salesTargetServiceControllerStub) Store(request entity.CreateSalesTargetRequest) error {
	s.storeCalled = true
	s.storeReq = request
	return s.storeErr
}

func (s *salesTargetServiceControllerStub) Update(_ int64, _ entity.UpdateSalesTargetRequest) error {
	return nil
}

func TestSalesTargetController_Create_WithoutStatusReturnsValidationError(t *testing.T) {
	app := fiber.New()
	v := validation.NewValidator()
	stub := &salesTargetServiceControllerStub{}
	controller := NewSalesTargetController(stub, v)

	app.Post("/v1/sales-target", func(c *fiber.Ctx) error {
		c.Locals("requestid", "req-123")
		c.Locals("cust_id", "C22001")
		c.Locals("parent_cust_id", "P22001")
		c.Locals("user_id", int64(10))
		return controller.Create(c)
	})

	body := `{
		"sales_target_distributor_yearly_id": 1,
		"sales_target_distributor_monthly_id": 2,
		"month": 2,
		"year": 2025,
		"allocated_total": 100,
		"monthly_target": 100,
		"remaining": 0,
		"data": [
			{
				"salesman_id": 1,
				"sales_team_id": 10,
				"allocated": 100
			}
		]
	}`

	req := httptest.NewRequest("POST", "/v1/sales-target", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	res, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if res.StatusCode != fiber.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", fiber.StatusBadRequest, res.StatusCode)
	}

	if stub.storeCalled {
		t.Fatalf("expected store not to be called when status is omitted")
	}
}

func TestSalesTargetController_Create_AllowsZeroRemaining(t *testing.T) {
	app := fiber.New()
	v := validation.NewValidator()
	stub := &salesTargetServiceControllerStub{}
	controller := NewSalesTargetController(stub, v)

	app.Post("/v1/sales-target", func(c *fiber.Ctx) error {
		c.Locals("requestid", "req-456")
		c.Locals("cust_id", "C22001")
		c.Locals("parent_cust_id", "P22001")
		c.Locals("user_id", int64(10))
		return controller.Create(c)
	})

	body := `{
		"sales_target_distributor_yearly_id": 43,
		"sales_target_distributor_monthly_id": 772,
		"month": 7,
		"year": 2026,
		"allocated_total": 100,
		"monthly_target": 100,
		"remaining": 0,
		"status": 1,
		"data": [
			{
				"salesman_id": 395,
				"sales_team_id": 53,
				"allocated": 100
			}
		]
	}`

	req := httptest.NewRequest("POST", "/v1/sales-target", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	res, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if res.StatusCode != fiber.StatusCreated {
		t.Fatalf("expected status %d, got %d", fiber.StatusCreated, res.StatusCode)
	}

	if !stub.storeCalled {
		t.Fatalf("expected store to be called")
	}

	if stub.storeReq.Remaining == nil {
		t.Fatalf("expected remaining pointer to be populated")
	}

	if *stub.storeReq.Remaining != 0 {
		t.Fatalf("expected remaining 0, got %d", *stub.storeReq.Remaining)
	}
}
