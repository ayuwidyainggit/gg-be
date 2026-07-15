package controller

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"sales/entity"
	"sales/pkg/validation"
	"sales/service"

	"github.com/gofiber/fiber/v2"
)

type mockOrderServiceForCreateController struct {
	service.OrderService
	storeFn func(request entity.CreateOrderBody, validationData entity.ValidateResponse) (entity.CreateOrderResponse, error)
}

func (m *mockOrderServiceForCreateController) Store(request entity.CreateOrderBody, validationData entity.ValidateResponse) (entity.CreateOrderResponse, error) {
	if m.storeFn != nil {
		return m.storeFn(request, validationData)
	}
	return entity.CreateOrderResponse{}, nil
}

type mockValidateOrderServiceForCreateController struct {
	service.ValidateOrderService
	validateOrderFn             func(dataFilter entity.ValidateOrderBody) (entity.ValidateResponse, int64, int, error)
	validateOrderWithoutStockFn func(dataFilter entity.ValidateOrderBody) (entity.ValidateResponse, int64, int, error)
}

func (m *mockValidateOrderServiceForCreateController) ValidateOrder(dataFilter entity.ValidateOrderBody) (entity.ValidateResponse, int64, int, error) {
	if m.validateOrderFn != nil {
		return m.validateOrderFn(dataFilter)
	}
	return entity.ValidateResponse{}, 0, 0, nil
}

func (m *mockValidateOrderServiceForCreateController) ValidateOrderWithoutStock(dataFilter entity.ValidateOrderBody) (entity.ValidateResponse, int64, int, error) {
	if m.validateOrderWithoutStockFn != nil {
		return m.validateOrderWithoutStockFn(dataFilter)
	}
	return entity.ValidateResponse{}, 0, 0, nil
}

func (m *mockValidateOrderServiceForCreateController) ValidateOrderDetail(dataFilter entity.ValidateOrderDetailBody) (entity.ValidateDetailResponse, int64, int, error) {
	return entity.ValidateDetailResponse{}, 0, 0, nil
}

func TestOrderControllerCreate_SX2184TakingOrderBypassesStockValidation(t *testing.T) {
	validator := validation.NewValiditor()
	validator.RegisterCustomValidation()

	validateCalled := false
	nonStockValidateCalled := false
	orderType := "O"
	orderDate := "2026-06-08"
	qty1 := 10.0
	qty2 := 0.0
	qty3 := 0.0
	convUnit2 := 10
	convUnit3 := 5
	sellPrice1 := 100.0
	sellPrice2 := 0.0
	sellPrice3 := 0.0
	vat := 0.0
	whID := int64(301)

	var storedRequest entity.CreateOrderBody
	var storedValidation entity.ValidateResponse

	controller := NewOrderController(
		&mockOrderServiceForCreateController{storeFn: func(request entity.CreateOrderBody, validationData entity.ValidateResponse) (entity.CreateOrderResponse, error) {
			storedRequest = request
			storedValidation = validationData
			return entity.CreateOrderResponse{RoNo: "SO2606080001"}, nil
		}},
		&mockValidateOrderServiceForCreateController{
			validateOrderFn: func(dataFilter entity.ValidateOrderBody) (entity.ValidateResponse, int64, int, error) {
				validateCalled = true
				return entity.ValidateResponse{}, 0, 0, nil
			},
			validateOrderWithoutStockFn: func(dataFilter entity.ValidateOrderBody) (entity.ValidateResponse, int64, int, error) {
				nonStockValidateCalled = true
				return entity.ValidateResponse{
					Validate1Success: true,
					Validate1:        "Sufficient Stock",
					Validate2Success: false,
					Validate2:        "Over Limit (1.000)",
					Validate2value:   1000,
					Validate3Success: true,
					Validate3:        "Allowed",
					Validate4Success: true,
					Validate4:        "Allowed",
				}, 0, 0, nil
			},
		},
		validator,
	)

	app := fiber.New()
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("requestid", "req-order-o")
		c.Locals("cust_id", "C220010001")
		c.Locals("parent_cust_id", "PARENT1")
		c.Locals("user_id", int64(99))
		return c.Next()
	})
	app.Post("/orders", controller.Create)

	body, _ := json.Marshal(entity.CreateOrderBody{
		OrderType:  &orderType,
		RoDate:     &orderDate,
		SalesmanId: 11,
		WhId:       &whID,
		OutletID:   21,
		Details: entity.OrderDetWithGroup{Normal: []entity.CreateOrderDetBody{{
			ProId:           748,
			Qty1:            &qty1,
			Qty2:            &qty2,
			Qty3:            &qty3,
			ConvUnit2:       &convUnit2,
			ConvUnit3:       &convUnit3,
			SellPrice1:      &sellPrice1,
			SellPrice2:      &sellPrice2,
			SellPrice3:      &sellPrice3,
			PromoValue:      &vat,
			PromoValueFinal: &vat,
			DiscValue:       &vat,
			Vat:             &vat,
			VatValue:        &vat,
			Amount:          &vat,
			AmountFinal:     &vat,
		}}},
	})

	req := httptest.NewRequest("POST", "/orders", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app test request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != fiber.StatusCreated {
		t.Fatalf("expected status 201, got %d", resp.StatusCode)
	}
	if validateCalled {
		t.Fatal("expected ValidateOrder not to be called for order_type O")
	}
	if !nonStockValidateCalled {
		t.Fatal("expected ValidateOrderWithoutStock to be called for order_type O")
	}
	if storedRequest.OrderType == nil || *storedRequest.OrderType != orderType {
		t.Fatalf("expected store request order_type O, got %+v", storedRequest.OrderType)
	}
	if !storedValidation.Validate1Success {
		t.Fatalf("expected stock validation snapshot to stay successful, got %+v", storedValidation)
	}
	if storedValidation.Validate2Success {
		t.Fatalf("expected non-stock credit validation to be preserved, got %+v", storedValidation)
	}
	if storedValidation.Validate2 != "Over Limit (1.000)" || storedValidation.Validate2value != 1000 {
		t.Fatalf("expected non-stock validation details to propagate, got %+v", storedValidation)
	}
	if storedValidation.IsSuccessValidate {
		t.Fatalf("expected overall validation summary to reflect preserved non-stock failure, got %+v", storedValidation)
	}
}

func TestOrderControllerCreate_SX2184NilOrderTypeStillValidatesStock(t *testing.T) {
	validator := validation.NewValiditor()
	validator.RegisterCustomValidation()

	validateCalled := false
	orderDate := "2026-06-08"
	qty1 := 10.0
	qty2 := 0.0
	qty3 := 0.0
	convUnit2 := 10
	convUnit3 := 5
	sellPrice1 := 100.0
	sellPrice2 := 0.0
	sellPrice3 := 0.0
	vat := 0.0
	whID := int64(301)

	controller := NewOrderController(
		&mockOrderServiceForCreateController{storeFn: func(request entity.CreateOrderBody, validationData entity.ValidateResponse) (entity.CreateOrderResponse, error) {
			return entity.CreateOrderResponse{RoNo: "SO2606080000"}, nil
		}},
		&mockValidateOrderServiceForCreateController{validateOrderFn: func(dataFilter entity.ValidateOrderBody) (entity.ValidateResponse, int64, int, error) {
			validateCalled = true
			return entity.ValidateResponse{}, 0, 0, errors.New("Insufficient Stock")
		}},
		validator,
	)

	app := fiber.New()
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("requestid", "req-order-nil")
		c.Locals("cust_id", "C220010001")
		c.Locals("parent_cust_id", "PARENT1")
		c.Locals("user_id", int64(99))
		return c.Next()
	})
	app.Post("/orders", controller.Create)

	body, _ := json.Marshal(entity.CreateOrderBody{
		RoDate:     &orderDate,
		SalesmanId: 11,
		WhId:       &whID,
		OutletID:   21,
		Details: entity.OrderDetWithGroup{Normal: []entity.CreateOrderDetBody{{
			ProId:           748,
			Qty1:            &qty1,
			Qty2:            &qty2,
			Qty3:            &qty3,
			ConvUnit2:       &convUnit2,
			ConvUnit3:       &convUnit3,
			SellPrice1:      &sellPrice1,
			SellPrice2:      &sellPrice2,
			SellPrice3:      &sellPrice3,
			PromoValue:      &vat,
			PromoValueFinal: &vat,
			DiscValue:       &vat,
			Vat:             &vat,
			VatValue:        &vat,
			Amount:          &vat,
			AmountFinal:     &vat,
		}}},
	})

	req := httptest.NewRequest("POST", "/orders", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app test request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != fiber.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", resp.StatusCode)
	}
	if !validateCalled {
		t.Fatal("expected ValidateOrder to be called when order_type is nil")
	}
}

func TestOrderControllerCreate_SX2184EmptyOrderTypeStillValidatesStock(t *testing.T) {
	validator := validation.NewValiditor()
	validator.RegisterCustomValidation()

	validateCalled := false
	orderType := ""
	orderDate := "2026-06-08"
	qty1 := 10.0
	qty2 := 0.0
	qty3 := 0.0
	convUnit2 := 10
	convUnit3 := 5
	sellPrice1 := 100.0
	sellPrice2 := 0.0
	sellPrice3 := 0.0
	vat := 0.0
	whID := int64(301)

	controller := NewOrderController(
		&mockOrderServiceForCreateController{storeFn: func(request entity.CreateOrderBody, validationData entity.ValidateResponse) (entity.CreateOrderResponse, error) {
			return entity.CreateOrderResponse{RoNo: "SO2606080001"}, nil
		}},
		&mockValidateOrderServiceForCreateController{validateOrderFn: func(dataFilter entity.ValidateOrderBody) (entity.ValidateResponse, int64, int, error) {
			validateCalled = true
			return entity.ValidateResponse{}, 0, 0, errors.New("Insufficient Stock")
		}},
		validator,
	)

	app := fiber.New()
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("requestid", "req-order-empty")
		c.Locals("cust_id", "C220010001")
		c.Locals("parent_cust_id", "PARENT1")
		c.Locals("user_id", int64(99))
		return c.Next()
	})
	app.Post("/orders", controller.Create)

	body, _ := json.Marshal(entity.CreateOrderBody{
		OrderType:  &orderType,
		RoDate:     &orderDate,
		SalesmanId: 11,
		WhId:       &whID,
		OutletID:   21,
		Details: entity.OrderDetWithGroup{Normal: []entity.CreateOrderDetBody{{
			ProId:           748,
			Qty1:            &qty1,
			Qty2:            &qty2,
			Qty3:            &qty3,
			ConvUnit2:       &convUnit2,
			ConvUnit3:       &convUnit3,
			SellPrice1:      &sellPrice1,
			SellPrice2:      &sellPrice2,
			SellPrice3:      &sellPrice3,
			PromoValue:      &vat,
			PromoValueFinal: &vat,
			DiscValue:       &vat,
			Vat:             &vat,
			VatValue:        &vat,
			Amount:          &vat,
			AmountFinal:     &vat,
		}}},
	})

	req := httptest.NewRequest("POST", "/orders", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app test request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != fiber.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", resp.StatusCode)
	}
	if !validateCalled {
		t.Fatal("expected ValidateOrder to be called when order_type is empty")
	}
}

func TestOrderControllerCreate_SX2184OrderTypeCStillValidatesStock(t *testing.T) {
	validator := validation.NewValiditor()
	validator.RegisterCustomValidation()

	validateCalled := false
	orderType := "C"
	orderDate := "2026-06-08"
	qty1 := 10.0
	qty2 := 0.0
	qty3 := 0.0
	convUnit2 := 10
	convUnit3 := 5
	sellPrice1 := 100.0
	sellPrice2 := 0.0
	sellPrice3 := 0.0
	vat := 0.0
	whID := int64(301)

	controller := NewOrderController(
		&mockOrderServiceForCreateController{storeFn: func(request entity.CreateOrderBody, validationData entity.ValidateResponse) (entity.CreateOrderResponse, error) {
			return entity.CreateOrderResponse{RoNo: "SO2606080002"}, nil
		}},
		&mockValidateOrderServiceForCreateController{validateOrderFn: func(dataFilter entity.ValidateOrderBody) (entity.ValidateResponse, int64, int, error) {
			validateCalled = true
			return entity.ValidateResponse{}, 0, 0, errors.New("Insufficient Stock")
		}},
		validator,
	)

	app := fiber.New()
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("requestid", "req-order-c")
		c.Locals("cust_id", "C220010001")
		c.Locals("parent_cust_id", "PARENT1")
		c.Locals("user_id", int64(99))
		return c.Next()
	})
	app.Post("/orders", controller.Create)

	body, _ := json.Marshal(entity.CreateOrderBody{
		OrderType:  &orderType,
		RoDate:     &orderDate,
		SalesmanId: 11,
		WhId:       &whID,
		OutletID:   21,
		Details: entity.OrderDetWithGroup{Normal: []entity.CreateOrderDetBody{{
			ProId:           748,
			Qty1:            &qty1,
			Qty2:            &qty2,
			Qty3:            &qty3,
			ConvUnit2:       &convUnit2,
			ConvUnit3:       &convUnit3,
			SellPrice1:      &sellPrice1,
			SellPrice2:      &sellPrice2,
			SellPrice3:      &sellPrice3,
			PromoValue:      &vat,
			PromoValueFinal: &vat,
			DiscValue:       &vat,
			Vat:             &vat,
			VatValue:        &vat,
			Amount:          &vat,
			AmountFinal:     &vat,
		}}},
	})

	req := httptest.NewRequest("POST", "/orders", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app test request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != fiber.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", resp.StatusCode)
	}
	if !validateCalled {
		t.Fatal("expected ValidateOrder to be called for order_type C")
	}
}

func TestOrderControllerCreate_SX2184SalesOrderStillValidatesStock(t *testing.T) {
	validator := validation.NewValiditor()
	validator.RegisterCustomValidation()

	validateCalled := false
	orderType := "SO"
	orderDate := "2026-06-08"
	qty1 := 10.0
	qty2 := 0.0
	qty3 := 0.0
	convUnit2 := 10
	convUnit3 := 5
	sellPrice1 := 100.0
	sellPrice2 := 0.0
	sellPrice3 := 0.0
	vat := 0.0
	whID := int64(301)

	controller := NewOrderController(
		&mockOrderServiceForCreateController{storeFn: func(request entity.CreateOrderBody, validationData entity.ValidateResponse) (entity.CreateOrderResponse, error) {
			return entity.CreateOrderResponse{RoNo: "SO2606080002"}, nil
		}},
		&mockValidateOrderServiceForCreateController{validateOrderFn: func(dataFilter entity.ValidateOrderBody) (entity.ValidateResponse, int64, int, error) {
			validateCalled = true
			return entity.ValidateResponse{}, 0, 0, errors.New("Insufficient Stock")
		}},
		validator,
	)

	app := fiber.New()
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("requestid", "req-order-so")
		c.Locals("cust_id", "C220010001")
		c.Locals("parent_cust_id", "PARENT1")
		c.Locals("user_id", int64(99))
		return c.Next()
	})
	app.Post("/orders", controller.Create)

	body, _ := json.Marshal(entity.CreateOrderBody{
		OrderType:  &orderType,
		RoDate:     &orderDate,
		SalesmanId: 11,
		WhId:       &whID,
		OutletID:   21,
		Details: entity.OrderDetWithGroup{Normal: []entity.CreateOrderDetBody{{
			ProId:           748,
			Qty1:            &qty1,
			Qty2:            &qty2,
			Qty3:            &qty3,
			ConvUnit2:       &convUnit2,
			ConvUnit3:       &convUnit3,
			SellPrice1:      &sellPrice1,
			SellPrice2:      &sellPrice2,
			SellPrice3:      &sellPrice3,
			PromoValue:      &vat,
			PromoValueFinal: &vat,
			DiscValue:       &vat,
			Vat:             &vat,
			VatValue:        &vat,
			Amount:          &vat,
			AmountFinal:     &vat,
		}}},
	})

	req := httptest.NewRequest("POST", "/orders", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app test request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != fiber.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", resp.StatusCode)
	}
	if !validateCalled {
		t.Fatal("expected ValidateOrder to be called for order_type SO")
	}
}

type mockOrderImportService struct {
	service.OrderService
	exportFn   func(format string) (*bytes.Buffer, string, string, error)
	importFn   func(custId, parentCustId string, userId int64, file io.Reader, filename string) (entity.OrderImportResult, []entity.OrderImportError, error)
	validateFn func(custId, parentCustId string, userId int64, file io.Reader, filename string) (entity.OrderImportSummary, error)
}

func (m *mockOrderImportService) ExportTemplate(format string) (*bytes.Buffer, string, string, error) {
	if m.exportFn != nil {
		return m.exportFn(format)
	}
	return bytes.NewBufferString(""), "application/octet-stream", "x", nil
}

func (m *mockOrderImportService) ImportOrders(custId, parentCustId string, userId int64, file io.Reader, filename string) (entity.OrderImportResult, []entity.OrderImportError, error) {
	if m.importFn != nil {
		return m.importFn(custId, parentCustId, userId, file, filename)
	}
	return entity.OrderImportResult{}, nil, nil
}

func (m *mockOrderImportService) ValidateImport(custId, parentCustId string, userId int64, file io.Reader, filename string) (entity.OrderImportSummary, error) {
	if m.validateFn != nil {
		return m.validateFn(custId, parentCustId, userId, file, filename)
	}
	return entity.OrderImportSummary{}, nil
}

func TestOrderControllerExportTemplate_SX2435(t *testing.T) {
	validator := validation.NewValiditor()
	validator.RegisterCustomValidation()

	controller := NewOrderController(
		&mockOrderImportService{exportFn: func(format string) (*bytes.Buffer, string, string, error) {
			if format != "xlsx" {
				t.Fatalf("expected default format xlsx, got %s", format)
			}
			return bytes.NewBufferString("excel"), "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", "order_import_template.xlsx", nil
		}},
		&mockValidateOrderServiceForCreateController{},
		validator,
	)

	app := fiber.New()
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("requestid", "req-export")
		c.Locals("cust_id", "C220010001")
		c.Locals("parent_cust_id", "PARENT1")
		c.Locals("user_id", int64(99))
		return c.Next()
	})
	app.Get("/v1/orders/export-template", controller.ExportTemplate)

	req := httptest.NewRequest("GET", "/v1/orders/export-template", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app test request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != fiber.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.StatusCode)
	}
	if got := resp.Header.Get("Content-Type"); got != "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet" {
		t.Fatalf("unexpected content type: %s", got)
	}
	if got := resp.Header.Get("Content-Disposition"); got != `attachment; filename="order_import_template.xlsx"` {
		t.Fatalf("unexpected content disposition: %s", got)
	}
}

func TestOrderControllerImport_SX2451RejectsWhenAnyRowInvalid(t *testing.T) {
	validator := validation.NewValiditor()
	validator.RegisterCustomValidation()

	controller := NewOrderController(
		&mockOrderImportService{importFn: func(custId, parentCustId string, userId int64, file io.Reader, filename string) (entity.OrderImportResult, []entity.OrderImportError, error) {
			return entity.OrderImportResult{StartDate: "2026-07-09", EndDate: "2026-07-09", NumberOfInvoice: 1, NumberOfOutlet: 1, Amount: 1000}, nil, &entity.ImportFailedError{FailedReasons: []string{"row 2: salesman_code not found"}}
		}},
		&mockValidateOrderServiceForCreateController{},
		validator,
	)

	app := fiber.New()
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("requestid", "req-import")
		c.Locals("cust_id", "C220010001")
		c.Locals("parent_cust_id", "PARENT1")
		c.Locals("user_id", int64(99))
		return c.Next()
	})
	app.Post("/v1/orders/import", controller.Import)

	body := &bytes.Buffer{}
	mw := multipart.NewWriter(body)
	_, _ = mw.CreateFormFile("file", "template.xlsx")
	_ = mw.Close()

	req := httptest.NewRequest("POST", "/v1/orders/import", body)
	req.Header.Set("Content-Type", mw.FormDataContentType())
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app test request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != fiber.StatusUnprocessableEntity {
		t.Fatalf("expected status 422 on validation errors, got %d", resp.StatusCode)
	}
}

func TestOrderControllerExportTemplate_SX2435RejectsXls(t *testing.T) {
	validator := validation.NewValiditor()
	validator.RegisterCustomValidation()

	controller := NewOrderController(
		&mockOrderImportService{exportFn: func(format string) (*bytes.Buffer, string, string, error) {
			// SX-2470: xls is a legacy alias for xlsx; the underlying
			// service must be called with the canonical "xlsx" value.
			if format != "xlsx" {
				t.Fatalf("expected service to receive xlsx, got %s", format)
			}
			return bytes.NewBufferString("excel"), "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", "order_import_template.xlsx", nil
		}},
		&mockValidateOrderServiceForCreateController{},
		validator,
	)

	app := fiber.New()
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("requestid", "req-export-xls")
		c.Locals("cust_id", "C220010001")
		c.Locals("parent_cust_id", "PARENT1")
		c.Locals("user_id", int64(99))
		return c.Next()
	})
	app.Get("/v1/orders/export-template", controller.ExportTemplate)

	req := httptest.NewRequest("GET", "/v1/orders/export-template?format=xls", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app test request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != fiber.StatusOK {
		t.Fatalf("expected status 200 for xls alias, got %d", resp.StatusCode)
	}
}

func TestOrderControllerImportFromUrl_SX2475ValidateOnly(t *testing.T) {
	// SX-2475: FE uploads file then POSTs {url, validate:"False"}.
	// BE must download, run validation, and return summary.
	// No DB calls: we mock the service so no real xlsx is required.
	fakeServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("not-a-real-xlsx"))
	}))
	defer fakeServer.Close()

	validateCalled := false
	importCalled := false
	validator := validation.NewValiditor()
	validator.RegisterCustomValidation()

	controller := NewOrderController(
		&mockOrderImportService{
			validateFn: func(custId, parentCustId string, userId int64, file io.Reader, filename string) (entity.OrderImportSummary, error) {
				validateCalled = true
				return entity.OrderImportSummary{StartDate: "2026-07-09", EndDate: "2026-07-09", NumberOfInvoice: 2, NumberOfOutlet: 1, Amount: 6069000, FailedReasons: []string{"row 2: qty1 quantity must be > 0"}}, nil
			},
			importFn: func(custId, parentCustId string, userId int64, file io.Reader, filename string) (entity.OrderImportResult, []entity.OrderImportError, error) {
				importCalled = true
				return entity.OrderImportResult{}, nil, nil
			},
		},
		&mockValidateOrderServiceForCreateController{},
		validator,
	)

	app := fiber.New()
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("requestid", "req-import-url")
		c.Locals("cust_id", "C220010001")
		c.Locals("parent_cust_id", "PARENT1")
		c.Locals("user_id", int64(99))
		return c.Next()
	})
	app.Post("/v1/orders/export-template/import", controller.ImportFromUrl)

	body, _ := json.Marshal(entity.OrderImportFromURLRequest{URL: fakeServer.URL + "/file.xlsx", Validate: "False"})
	req := httptest.NewRequest("POST", "/v1/orders/export-template/import", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("app test request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != fiber.StatusUnprocessableEntity {
		t.Fatalf("expected status 422, got %d", resp.StatusCode)
	}
	if !validateCalled {
		t.Fatalf("expected ValidateImport to be called")
	}
	if importCalled {
		t.Fatalf("expected ImportOrders NOT to be called when validate=False")
	}
}

func TestOrderControllerImportFromUrl_SX2475MissingURL(t *testing.T) {
	validator := validation.NewValiditor()
	validator.RegisterCustomValidation()

	controller := NewOrderController(
		&mockOrderImportService{},
		&mockValidateOrderServiceForCreateController{},
		validator,
	)

	app := fiber.New()
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("requestid", "req-import-url-empty")
		c.Locals("cust_id", "C220010001")
		c.Locals("parent_cust_id", "PARENT1")
		c.Locals("user_id", int64(99))
		return c.Next()
	})
	app.Post("/v1/orders/export-template/import", controller.ImportFromUrl)

	body, _ := json.Marshal(entity.OrderImportFromURLRequest{URL: "   ", Validate: "False"})
	req := httptest.NewRequest("POST", "/v1/orders/export-template/import", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("app test request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != fiber.StatusBadRequest {
		t.Fatalf("expected status 400 for missing url, got %d", resp.StatusCode)
	}
	respBody, _ := io.ReadAll(resp.Body)
	if !strings.Contains(string(respBody), "url is required") {
		t.Fatalf("expected response to mention url is required, got: %s", string(respBody))
	}
}

