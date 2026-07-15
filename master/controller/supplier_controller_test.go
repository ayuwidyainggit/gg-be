package controller

import (
	"net/http/httptest"
	"strings"
	"testing"

	"master/entity"
	"master/pkg/validation"

	"github.com/gofiber/fiber/v2"
)

type supplierServiceStub struct {
	capturedCreate entity.CreateSupplierBody
}

func (s *supplierServiceStub) Detail(int, string) (entity.SupplierResponse, error) {
	return entity.SupplierResponse{}, nil
}

func (s *supplierServiceStub) List(entity.SupplierQueryFilter, string) ([]entity.SupplierResponse, int, int, error) {
	return []entity.SupplierResponse{}, 0, 0, nil
}

func (s *supplierServiceStub) LookupList(entity.SupplierQueryFilter, entity.SupplierLookupScope) ([]entity.SupplierLookupResponse, int, int, error) {
	return []entity.SupplierLookupResponse{}, 0, 0, nil
}

func (s *supplierServiceStub) Store(request entity.CreateSupplierBody) (entity.SupplierResponse, error) {
	s.capturedCreate = request
	return entity.SupplierResponse{}, nil
}

func (s *supplierServiceStub) Update(int, entity.UpdateSupplierRequest) error {
	return nil
}

func (s *supplierServiceStub) Delete(string, int, int64) error {
	return nil
}

func TestSupplierController_Create_UsesLoggedInCustIDAndDistributorIDFromJWT(t *testing.T) {
	svc := &supplierServiceStub{}
	controller := &SupplierController{SupplierService: svc, validator: validation.NewValidator()}

	app := fiber.New()
	app.Post("/v1/suppliers", func(c *fiber.Ctx) error {
		c.Locals("requestid", "req-supplier-create")
		c.Locals("cust_id", "C22001")
		c.Locals("parent_cust_id", "P22001")
		c.Locals("user_id", int64(88))
		c.Locals("distributor_id", int64(55))
		return controller.Create(c)
	})

	body := `{
		"cust_id":"IGNORED",
		"sup_code":"SUP001",
		"sup_name":"Supplier One",
		"address1":"Jl. Mawar 1",
		"fax_no":"021123456",
		"contact_name":"Budi",
		"zip_code":"123456",
		"email":"budi@example.com",
		"phone_no":"08123456789",
		"sup_type":"I",
		"wa_no":"08123456789",
		"credit_limit_type":"L",
		"is_active":true
	}`
	req := httptest.NewRequest("POST", "/v1/suppliers", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	res, err := app.Test(req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if res.StatusCode != fiber.StatusCreated {
		t.Fatalf("expected status 201, got %d", res.StatusCode)
	}

	if svc.capturedCreate.CustId != "C22001" {
		t.Fatalf("expected cust_id C22001, got %s", svc.capturedCreate.CustId)
	}

	if svc.capturedCreate.CreatedBy != 88 {
		t.Fatalf("expected created_by 88, got %d", svc.capturedCreate.CreatedBy)
	}

	if svc.capturedCreate.DistributorId == nil || *svc.capturedCreate.DistributorId != 55 {
		t.Fatalf("expected distributor_id 55, got %#v", svc.capturedCreate.DistributorId)
	}
}

func TestSupplierController_Create_LeavesDistributorIDNilWhenJWTDistributorIsZero(t *testing.T) {
	svc := &supplierServiceStub{}
	controller := &SupplierController{SupplierService: svc, validator: validation.NewValidator()}

	app := fiber.New()
	app.Post("/v1/suppliers", func(c *fiber.Ctx) error {
		c.Locals("requestid", "req-supplier-create-zero")
		c.Locals("cust_id", "C22001")
		c.Locals("parent_cust_id", "P22001")
		c.Locals("user_id", int64(88))
		c.Locals("distributor_id", int64(0))
		return controller.Create(c)
	})

	body := `{
		"sup_code":"SUP002",
		"sup_name":"Supplier Two",
		"address1":"Jl. Melati 2",
		"fax_no":"021123450",
		"contact_name":"Sari",
		"zip_code":"654321",
		"email":"sari@example.com",
		"phone_no":"08123456780",
		"sup_type":"I",
		"wa_no":"08123456780",
		"credit_limit_type":"L",
		"is_active":true
	}`
	req := httptest.NewRequest("POST", "/v1/suppliers", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	res, err := app.Test(req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if res.StatusCode != fiber.StatusCreated {
		t.Fatalf("expected status 201, got %d", res.StatusCode)
	}

	if svc.capturedCreate.DistributorId != nil {
		t.Fatalf("expected distributor_id nil, got %#v", svc.capturedCreate.DistributorId)
	}
}
