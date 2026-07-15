package controller

import (
	"encoding/json"
	"net/http/httptest"
	"testing"

	"master/entity"

	"github.com/gofiber/fiber/v2"
)

type regionServiceStub struct{ capturedFilter entity.RegionQueryFilter }
func (s *regionServiceStub) Detail(int, string) (entity.RegionResponse, error) { return entity.RegionResponse{}, nil }
func (s *regionServiceStub) LookupList(filter entity.RegionQueryFilter) ([]entity.RegionLookupResponse, int, int, error) {
	s.capturedFilter = filter
	return []entity.RegionLookupResponse{}, 0, 0, nil
}
func (s *regionServiceStub) List(filter entity.RegionQueryFilter) ([]entity.RegionResponse, int, int, error) {
	s.capturedFilter = filter
	return []entity.RegionResponse{}, 0, 0, nil
}
func (s *regionServiceStub) Store(entity.CreateRegionBody) (entity.RegionResponse, error) { return entity.RegionResponse{}, nil }
func (s *regionServiceStub) Update(int, entity.UpdateRegionRequest) error { return nil }
func (s *regionServiceStub) Delete(string, int, int64) error { return nil }

type areaServiceStub2 struct{ capturedFilter entity.AreaQueryFilter }
func (s *areaServiceStub2) Detail(int, string) (entity.AreaListResponse, error) { return entity.AreaListResponse{}, nil }
func (s *areaServiceStub2) LookupList(filter entity.AreaQueryFilter) ([]entity.AreaListResponse, int, int, error) {
	s.capturedFilter = filter
	return []entity.AreaListResponse{}, 0, 0, nil
}
func (s *areaServiceStub2) List(filter entity.AreaQueryFilter) ([]entity.AreaListResponse, int, int, error) {
	s.capturedFilter = filter
	return []entity.AreaListResponse{}, 0, 0, nil
}
func (s *areaServiceStub2) Store(entity.CreateAreaBody) (entity.AreaResponse, error) { return entity.AreaResponse{}, nil }
func (s *areaServiceStub2) Update(int, entity.UpdateAreaRequest) error { return nil }
func (s *areaServiceStub2) Delete(string, int, int64) error { return nil }

func TestRegionController_List_ParsesTokenContextAndArrayQuery(t *testing.T) {
	svc := &regionServiceStub{}
	controller := &RegionController{RegionService: svc}
	app := fiber.New()
	app.Get("/v1/regions", func(c *fiber.Ctx) error {
		c.Locals("requestid", "req-r")
		c.Locals("cust_id", "C22001")
		c.Locals("parent_cust_id", "C22001")
		c.Locals("employee_id", int64(77))
		c.Locals("distributor_id", int64(0))
		return controller.List(c)
	})

	req := httptest.NewRequest("GET", "/v1/regions?page=1&limit=10&region_id[]=1,2", nil)
	res, err := app.Test(req)
	if err != nil { t.Fatalf("expected no error, got %v", err) }
	if res.StatusCode != fiber.StatusOK { t.Fatalf("expected status 200, got %d", res.StatusCode) }
	if svc.capturedFilter.EmployeeId != 77 || svc.capturedFilter.DistributorId != 0 {
		t.Fatalf("expected employee/distributor from token, got %+v", svc.capturedFilter)
	}
	if len(svc.capturedFilter.RegionId) != 2 {
		t.Fatalf("expected region array parsed, got %+v", svc.capturedFilter.RegionId)
	}
	var resp map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&resp); err != nil { t.Fatalf("decode response: %v", err) }
}

func TestAreaController_List_ParsesTokenContextAndArrayQuery(t *testing.T) {
	svc := &areaServiceStub2{}
	controller := &AreaController{AreaService: svc}
	app := fiber.New()
	app.Get("/v1/areas", func(c *fiber.Ctx) error {
		c.Locals("requestid", "req-a")
		c.Locals("cust_id", "C22001")
		c.Locals("parent_cust_id", "C22001")
		c.Locals("employee_id", int64(88))
		c.Locals("distributor_id", int64(0))
		return controller.List(c)
	})

	req := httptest.NewRequest("GET", "/v1/areas?page=1&limit=10&region_id[]=1&area_id[]=10,20", nil)
	res, err := app.Test(req)
	if err != nil { t.Fatalf("expected no error, got %v", err) }
	if res.StatusCode != fiber.StatusOK { t.Fatalf("expected status 200, got %d", res.StatusCode) }
	if svc.capturedFilter.EmployeeId != 88 || svc.capturedFilter.DistributorId != 0 {
		t.Fatalf("expected employee/distributor from token, got %+v", svc.capturedFilter)
	}
	if len(svc.capturedFilter.AreaId) != 2 || len(svc.capturedFilter.RegionId) != 1 {
		t.Fatalf("expected area/region arrays parsed, got %+v", svc.capturedFilter)
	}
}
