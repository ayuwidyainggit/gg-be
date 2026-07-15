package controller

import (
	"encoding/json"
	"net/http/httptest"
	"testing"

	"master/entity"

	"github.com/gofiber/fiber/v2"
)

type businessUnitServiceStub struct {
	capturedFilter entity.BusinessUnitQueryFilter
}

func (s *businessUnitServiceStub) GetBusinessUnit(dataFilter entity.BusinessUnitQueryFilter) (interface{}, int, int, error) {
	s.capturedFilter = dataFilter
	return entity.BusinessUnitPrincipalResponse{}, 1, 1, nil
}

func TestBusinessUnitController_List_ParseCommaSeparatedArrayQuery(t *testing.T) {
	svc := &businessUnitServiceStub{}
	controller := &BusinessUnitController{Service: svc}

	app := fiber.New()
	app.Get("/v1/business-unit", func(c *fiber.Ctx) error {
		c.Locals("requestid", "req-1")
		c.Locals("cust_id", "C22001")
		c.Locals("parent_cust_id", "C22001")
		c.Locals("user_name", "princ@idetama.id")
		return controller.List(c)
	})

	req := httptest.NewRequest("GET", "/v1/business-unit?page=1&limit=10&q=dist&sort=area_id:asc&region_id[]=1,2,3&area_id[]=10,20&is_active[]=1", nil)
	res, err := app.Test(req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if res.StatusCode != fiber.StatusOK {
		t.Fatalf("expected status 200, got %d", res.StatusCode)
	}

	if len(svc.capturedFilter.RegionId) != 3 || svc.capturedFilter.RegionId[0] != 1 || svc.capturedFilter.RegionId[1] != 2 || svc.capturedFilter.RegionId[2] != 3 {
		t.Fatalf("expected region_id []int{1,2,3}, got %#v", svc.capturedFilter.RegionId)
	}

	if len(svc.capturedFilter.AreaId) != 2 || svc.capturedFilter.AreaId[0] != 10 || svc.capturedFilter.AreaId[1] != 20 {
		t.Fatalf("expected area_id []int{10,20}, got %#v", svc.capturedFilter.AreaId)
	}

	if svc.capturedFilter.Query != "dist" {
		t.Fatalf("expected q=dist, got %s", svc.capturedFilter.Query)
	}

	if svc.capturedFilter.Sort != "area_id:asc" {
		t.Fatalf("expected sort=area_id:asc, got %s", svc.capturedFilter.Sort)
	}

	if svc.capturedFilter.Page != 1 || svc.capturedFilter.Limit != 10 {
		t.Fatalf("expected page=1 limit=10, got page=%d limit=%d", svc.capturedFilter.Page, svc.capturedFilter.Limit)
	}

	var resp map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&resp); err != nil {
		t.Fatalf("failed decode response: %v", err)
	}
}

func TestBusinessUnitController_List_ParseRepeatedArrayQuery(t *testing.T) {
	svc := &businessUnitServiceStub{}
	controller := &BusinessUnitController{Service: svc}

	app := fiber.New()
	app.Get("/v1/business-unit", func(c *fiber.Ctx) error {
		c.Locals("requestid", "req-2")
		c.Locals("cust_id", "C22001")
		c.Locals("parent_cust_id", "C22001")
		c.Locals("user_name", "princ@idetama.id")
		return controller.List(c)
	})

	req := httptest.NewRequest("GET", "/v1/business-unit?region_id[]=1&region_id[]=2&area_id[]=10&is_active[]=1", nil)
	res, err := app.Test(req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if res.StatusCode != fiber.StatusOK {
		t.Fatalf("expected status 200, got %d", res.StatusCode)
	}

	if len(svc.capturedFilter.RegionId) != 2 || svc.capturedFilter.RegionId[0] != 1 || svc.capturedFilter.RegionId[1] != 2 {
		t.Fatalf("expected region_id []int{1,2}, got %#v", svc.capturedFilter.RegionId)
	}
}

func TestBusinessUnitController_List_ParseNonBracketCommaAndWhitespaceQuery(t *testing.T) {
	svc := &businessUnitServiceStub{}
	controller := &BusinessUnitController{Service: svc}

	app := fiber.New()
	app.Get("/v1/business-unit", func(c *fiber.Ctx) error {
		c.Locals("requestid", "req-3")
		c.Locals("cust_id", "C22001")
		c.Locals("parent_cust_id", "C22001")
		c.Locals("user_name", "princ@idetama.id")
		return controller.List(c)
	})

	req := httptest.NewRequest("GET", "/v1/business-unit?region_id=80,%2090&area_id=89,%2090&region_id=91&area_id=91", nil)
	res, err := app.Test(req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if res.StatusCode != fiber.StatusOK {
		t.Fatalf("expected status 200, got %d", res.StatusCode)
	}

	if got := svc.capturedFilter.RegionId; len(got) != 3 || got[0] != 80 || got[1] != 90 || got[2] != 91 {
		t.Fatalf("expected region_id []int{80,90,91}, got %#v", got)
	}
	if got := svc.capturedFilter.AreaId; len(got) != 3 || got[0] != 89 || got[1] != 90 || got[2] != 91 {
		t.Fatalf("expected area_id []int{89,90,91}, got %#v", got)
	}
}

func TestBusinessUnitController_List_InvalidNumericTokenReturns400(t *testing.T) {
	svc := &businessUnitServiceStub{}
	controller := &BusinessUnitController{Service: svc}

	app := fiber.New()
	app.Get("/v1/business-unit", func(c *fiber.Ctx) error {
		c.Locals("requestid", "req-4")
		c.Locals("cust_id", "C22001")
		c.Locals("parent_cust_id", "C22001")
		c.Locals("user_name", "princ@idetama.id")
		return controller.List(c)
	})

	req := httptest.NewRequest("GET", "/v1/business-unit?region_id=80,abc", nil)
	res, err := app.Test(req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if res.StatusCode != fiber.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", res.StatusCode)
	}
}
