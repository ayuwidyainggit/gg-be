package controller

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"master/entity"
	"master/pkg/validation"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
)

type mPriceServiceControllerStub struct {
	detailFn   func(entity.DetailMPriceParams) (entity.MPriceResponse, error)
	listFn     func(entity.MPriceQueryFilter, string) ([]entity.MPriceResponse, int, int, error)
	storeFn    func(entity.CreateMPriceBody) (entity.MPriceResponse, error)
	updateFn   func(string, entity.UpdateMPriceRequest) error
	publishFn  func(entity.PublishMPriceParams) error
	cancelFn   func(entity.CancelMPriceParams) error
	deleteFn   func(string, string, int64) error
	templateFn func(string, string, string, int64) (*bytes.Buffer, string, string, error)
	exportFn   func(entity.MPriceQueryFilter, string, string) (*bytes.Buffer, string, string, error)
	importFn   func(entity.MPriceImportRequest, string, string, int64, int64, string) (entity.MPriceImportResponse, error)

	capturedDetailParams   entity.DetailMPriceParams
	capturedListFilter     entity.MPriceQueryFilter
	capturedListCustID     string
	capturedStoreRequest   entity.CreateMPriceBody
	capturedUpdateID       string
	capturedUpdateReq      entity.UpdateMPriceRequest
	capturedPublish        entity.PublishMPriceParams
	capturedCancel         entity.CancelMPriceParams
	capturedDeleteCustID   string
	capturedDeleteID       string
	capturedDeleteUserID   int64
	capturedTemplateFormat string
	capturedExportFilter   entity.MPriceQueryFilter
	capturedImportReq      entity.MPriceImportRequest
	capturedImportCustID   string
}

func (s *mPriceServiceControllerStub) Detail(params entity.DetailMPriceParams) (entity.MPriceResponse, error) {
	s.capturedDetailParams = params
	if s.detailFn != nil {
		return s.detailFn(params)
	}
	return entity.MPriceResponse{PriceID: params.PriceID, Status: 1, StatusDesc: "Scheduled"}, nil
}

func (s *mPriceServiceControllerStub) List(filter entity.MPriceQueryFilter, custID string) ([]entity.MPriceResponse, int, int, error) {
	s.capturedListFilter = filter
	s.capturedListCustID = custID
	if s.listFn != nil {
		return s.listFn(filter, custID)
	}
	return []entity.MPriceResponse{{PriceID: "price-1", Status: 1}}, 1, 1, nil
}

func (s *mPriceServiceControllerStub) Store(request entity.CreateMPriceBody) (entity.MPriceResponse, error) {
	s.capturedStoreRequest = request
	if s.storeFn != nil {
		return s.storeFn(request)
	}
	return entity.MPriceResponse{PriceID: "price-1", Status: 1}, nil
}

func (s *mPriceServiceControllerStub) Update(priceID string, request entity.UpdateMPriceRequest) error {
	s.capturedUpdateID = priceID
	s.capturedUpdateReq = request
	if s.updateFn != nil {
		return s.updateFn(priceID, request)
	}
	return nil
}

func (s *mPriceServiceControllerStub) Publish(params entity.PublishMPriceParams) error {
	s.capturedPublish = params
	if s.publishFn != nil {
		return s.publishFn(params)
	}
	return nil
}

func (s *mPriceServiceControllerStub) PublishByRMQ(entity.PublishByRmqMPriceReq) error {
	return nil
}

func (s *mPriceServiceControllerStub) Cancel(params entity.CancelMPriceParams) error {
	s.capturedCancel = params
	if s.cancelFn != nil {
		return s.cancelFn(params)
	}
	return nil
}

func (s *mPriceServiceControllerStub) Delete(custID string, priceID string, userID int64) error {
	s.capturedDeleteCustID = custID
	s.capturedDeleteID = priceID
	s.capturedDeleteUserID = userID
	if s.deleteFn != nil {
		return s.deleteFn(custID, priceID, userID)
	}
	return nil
}

func (s *mPriceServiceControllerStub) Template(format string, custID string, parentCustID string, distributorID int64) (*bytes.Buffer, string, string, error) {
	s.capturedTemplateFormat = format
	if s.templateFn != nil {
		return s.templateFn(format, custID, parentCustID, distributorID)
	}
	return bytes.NewBufferString("template-bytes"), "text/csv", "manage-price-template.csv", nil
}

func (s *mPriceServiceControllerStub) Export(filter entity.MPriceQueryFilter, custID string, parentCustID string) (*bytes.Buffer, string, string, error) {
	s.capturedExportFilter = filter
	if s.exportFn != nil {
		return s.exportFn(filter, custID, parentCustID)
	}
	return bytes.NewBufferString("export-bytes"), "text/csv", "manage-price-export.csv", nil
}

func (s *mPriceServiceControllerStub) Import(request entity.MPriceImportRequest, custID string, parentCustID string, userID int64, distributorID int64, userFullName string) (entity.MPriceImportResponse, error) {
	s.capturedImportReq = request
	s.capturedImportCustID = custID
	if s.importFn != nil {
		return s.importFn(request, custID, parentCustID, userID, distributorID, userFullName)
	}
	return entity.MPriceImportResponse{FileURL: request.FileURL, TotalRow: 1, SuccessRow: 1}, nil
}

func newMPriceControllerTestApp(stub *mPriceServiceControllerStub) (*fiber.App, MPriceController) {
	app := fiber.New()
	controller := NewMPriceController(stub, validation.NewValidator())
	return app, controller
}

func setMPriceControllerLocals(c *fiber.Ctx) {
	c.Locals("requestid", "req-price-1")
	c.Locals("cust_id", "C22001")
	c.Locals("parent_cust_id", "P22001")
	c.Locals("user_fullname", "Price User")
	c.Locals("user_id", int64(99))
	c.Locals("distributor_id", int64(77))
}

func registerMPriceRoute(app *fiber.App, method string, path string, handler func(*fiber.Ctx) error) {
	app.Add(method, path, func(c *fiber.Ctx) error {
		setMPriceControllerLocals(c)
		return handler(c)
	})
}

func validCreateMPriceBody() string {
	return `{
		"coverage": "D",
		"distributor_ids": [67, 68],
		"effective_date": "2026-05-12",
		"pro_id": 123,
		"unit_id1": "PCS",
		"unit_id2": "BOX",
		"unit_id3": "CTN",
		"conv_unit2": 10,
		"conv_unit3": 100,
		"purch_price1": 100,
		"purch_price2": 200,
		"purch_price3": 300,
		"sell_price1": 150,
		"sell_price2": 250,
		"sell_price3": 350,
		"new_purch_price1": 110,
		"new_purch_price2": 210,
		"new_purch_price3": 310,
		"new_sell_price1": 160,
		"new_sell_price2": 260,
		"new_sell_price3": 360
	}`
}

func TestMPriceController_List_ShouldParseFiltersAndPaging(t *testing.T) {
	stub := &mPriceServiceControllerStub{
		listFn: func(filter entity.MPriceQueryFilter, custID string) ([]entity.MPriceResponse, int, int, error) {
			return []entity.MPriceResponse{{PriceID: "price-1", Status: 1}}, 12, 3, nil
		},
	}
	app, controller := newMPriceControllerTestApp(stub)
	registerMPriceRoute(app, fiber.MethodGet, "/v1/prices", controller.List)

	req := httptest.NewRequest("GET", "/v1/prices?page=2&limit=5&q=sku&sort=price_id:desc&status[]=1,5&status[]=10&distributor_id[]=0,67&distributor_id[]=68&effective_date_start=2026-05-01&effective_date_end=2026-05-31", nil)
	res, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.StatusCode != fiber.StatusOK {
		t.Fatalf("expected status %d, got %d", fiber.StatusOK, res.StatusCode)
	}
	if stub.capturedListCustID != "C22001" {
		t.Fatalf("expected cust id C22001, got %s", stub.capturedListCustID)
	}
	if stub.capturedListFilter.Page != 2 || stub.capturedListFilter.Limit != 5 || stub.capturedListFilter.Query != "sku" {
		t.Fatalf("unexpected list filter: %+v", stub.capturedListFilter)
	}
	if len(stub.capturedListFilter.Status) != 3 || stub.capturedListFilter.Status[0] != 1 || stub.capturedListFilter.Status[1] != 5 || stub.capturedListFilter.Status[2] != 10 {
		t.Fatalf("unexpected status filter: %+v", stub.capturedListFilter.Status)
	}
	if len(stub.capturedListFilter.DistributorIDs) != 2 || stub.capturedListFilter.DistributorIDs[0] != 67 || stub.capturedListFilter.DistributorIDs[1] != 68 {
		t.Fatalf("unexpected distributor filter: %+v", stub.capturedListFilter.DistributorIDs)
	}

	var payload struct {
		Paging entity.Pagination `json:"paging"`
	}
	if err := json.NewDecoder(res.Body).Decode(&payload); err != nil {
		t.Fatalf("failed decode response: %v", err)
	}
	if payload.Paging.TotalRecord != 12 || payload.Paging.PageCurrent != 2 || payload.Paging.PageLimit != 5 || payload.Paging.PageTotal != 3 {
		t.Fatalf("unexpected paging payload: %+v", payload.Paging)
	}
}

func TestMPriceController_Detail_ShouldPassRouteAndTenantParams(t *testing.T) {
	stub := &mPriceServiceControllerStub{}
	app, controller := newMPriceControllerTestApp(stub)
	registerMPriceRoute(app, fiber.MethodGet, "/v1/prices/:price_id", controller.Detail)

	res, err := app.Test(httptest.NewRequest("GET", "/v1/prices/price-123", nil), -1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.StatusCode != fiber.StatusOK {
		t.Fatalf("expected status %d, got %d", fiber.StatusOK, res.StatusCode)
	}
	if stub.capturedDetailParams.PriceID != "price-123" || stub.capturedDetailParams.CustID != "C22001" || stub.capturedDetailParams.ParentCustID != "P22001" {
		t.Fatalf("unexpected detail params: %+v", stub.capturedDetailParams)
	}
}

func TestMPriceController_Create_ShouldValidateAndEnrichRequest(t *testing.T) {
	stub := &mPriceServiceControllerStub{}
	app, controller := newMPriceControllerTestApp(stub)
	registerMPriceRoute(app, fiber.MethodPost, "/v1/prices", controller.Create)

	req := httptest.NewRequest("POST", "/v1/prices", strings.NewReader(validCreateMPriceBody()))
	req.Header.Set("Content-Type", "application/json")
	res, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.StatusCode != fiber.StatusCreated {
		t.Fatalf("expected status %d, got %d", fiber.StatusCreated, res.StatusCode)
	}
	if stub.capturedStoreRequest.CustID != "C22001" || stub.capturedStoreRequest.ParentCustID != "P22001" {
		t.Fatalf("expected tenant values to be injected, got %+v", stub.capturedStoreRequest)
	}
	if stub.capturedStoreRequest.CreatedBy != "Price User" || stub.capturedStoreRequest.CreatedByID == nil || *stub.capturedStoreRequest.CreatedByID != 99 {
		t.Fatalf("expected creator values to be injected, got %+v", stub.capturedStoreRequest)
	}
	if stub.capturedStoreRequest.DistributorID != 77 || len(stub.capturedStoreRequest.DistributorIDs) != 2 {
		t.Fatalf("unexpected distributor values: %+v", stub.capturedStoreRequest)
	}
}

func TestMPriceController_Create_ShouldReturnValidationError(t *testing.T) {
	stub := &mPriceServiceControllerStub{}
	app, controller := newMPriceControllerTestApp(stub)
	registerMPriceRoute(app, fiber.MethodPost, "/v1/prices", controller.Create)

	req := httptest.NewRequest("POST", "/v1/prices", strings.NewReader(`{"coverage":"X"}`))
	req.Header.Set("Content-Type", "application/json")
	res, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.StatusCode != fiber.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", fiber.StatusBadRequest, res.StatusCode)
	}
	if stub.capturedStoreRequest.CustID != "" {
		t.Fatalf("expected service not to be called on validation error")
	}
}

func TestMPriceController_Update_ShouldSupportPatchAndPutRoutes(t *testing.T) {
	for _, method := range []string{fiber.MethodPatch, fiber.MethodPut} {
		t.Run(method, func(t *testing.T) {
			stub := &mPriceServiceControllerStub{}
			app, controller := newMPriceControllerTestApp(stub)
			registerMPriceRoute(app, method, "/v1/prices/:price_id", controller.Update)

			req := httptest.NewRequest(method, "/v1/prices/price-456", strings.NewReader(validCreateMPriceBody()))
			req.Header.Set("Content-Type", "application/json")
			res, err := app.Test(req, -1)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if res.StatusCode != fiber.StatusOK {
				t.Fatalf("expected status %d, got %d", fiber.StatusOK, res.StatusCode)
			}
			if stub.capturedUpdateID != "price-456" {
				t.Fatalf("expected update id price-456, got %s", stub.capturedUpdateID)
			}
			if stub.capturedUpdateReq.CustID != "C22001" || stub.capturedUpdateReq.ParentCustID != "P22001" {
				t.Fatalf("expected tenant values to be injected, got %+v", stub.capturedUpdateReq)
			}
			if stub.capturedUpdateReq.UpdatedByID == nil || *stub.capturedUpdateReq.UpdatedByID != 99 || stub.capturedUpdateReq.DistributorID != 77 {
				t.Fatalf("expected updater values to be injected, got %+v", stub.capturedUpdateReq)
			}
		})
	}
}

func TestMPriceController_Delete_ShouldPassCustIDPriceIDAndUserID(t *testing.T) {
	stub := &mPriceServiceControllerStub{}
	app, controller := newMPriceControllerTestApp(stub)
	registerMPriceRoute(app, fiber.MethodDelete, "/v1/prices/:price_id", controller.Delete)

	res, err := app.Test(httptest.NewRequest("DELETE", "/v1/prices/price-789", nil), -1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.StatusCode != fiber.StatusOK {
		t.Fatalf("expected status %d, got %d", fiber.StatusOK, res.StatusCode)
	}
	if stub.capturedDeleteCustID != "C22001" || stub.capturedDeleteID != "price-789" || stub.capturedDeleteUserID != 99 {
		t.Fatalf("unexpected delete params: cust=%s price=%s user=%d", stub.capturedDeleteCustID, stub.capturedDeleteID, stub.capturedDeleteUserID)
	}
}

func TestMPriceController_GetStatuses_ShouldReturnSortedStatuses(t *testing.T) {
	stub := &mPriceServiceControllerStub{}
	app, controller := newMPriceControllerTestApp(stub)
	registerMPriceRoute(app, fiber.MethodGet, "/v1/prices/statuses", controller.GetStatuses)

	res, err := app.Test(httptest.NewRequest("GET", "/v1/prices/statuses", nil), -1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.StatusCode != fiber.StatusOK {
		t.Fatalf("expected status %d, got %d", fiber.StatusOK, res.StatusCode)
	}

	var payload struct {
		Data []entity.MPriceStatus `json:"data"`
	}
	if err := json.NewDecoder(res.Body).Decode(&payload); err != nil {
		t.Fatalf("failed decode response: %v", err)
	}
	if len(payload.Data) != 3 || payload.Data[0].StatusID != 1 || payload.Data[1].StatusID != 5 || payload.Data[2].StatusID != 10 {
		t.Fatalf("unexpected statuses: %+v", payload.Data)
	}
}

func TestMPriceController_Publish_ShouldSupportLegacyAndCurrentRoutes(t *testing.T) {
	for _, path := range []string{"/v1/prices/publish/:price_id", "/v1/prices/:price_id/publish"} {
		t.Run(path, func(t *testing.T) {
			stub := &mPriceServiceControllerStub{}
			app, controller := newMPriceControllerTestApp(stub)
			registerMPriceRoute(app, fiber.MethodPatch, path, controller.Publish)

			realPath := strings.Replace(path, ":price_id", "price-901", 1)
			res, err := app.Test(httptest.NewRequest("PATCH", realPath, nil), -1)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if res.StatusCode != fiber.StatusOK {
				t.Fatalf("expected status %d, got %d", fiber.StatusOK, res.StatusCode)
			}
			if stub.capturedPublish.PriceID != "price-901" || stub.capturedPublish.CustID != "C22001" || stub.capturedPublish.ParentCustID != "P22001" {
				t.Fatalf("unexpected publish params: %+v", stub.capturedPublish)
			}
			if stub.capturedPublish.UpdatedByID == nil || *stub.capturedPublish.UpdatedByID != 99 || stub.capturedPublish.DistributorID != 77 {
				t.Fatalf("expected publish audit values, got %+v", stub.capturedPublish)
			}
		})
	}
}

func TestMPriceController_Cancel_ShouldSupportLegacyAndCurrentRoutes(t *testing.T) {
	for _, path := range []string{"/v1/prices/cancel/:price_id", "/v1/prices/:price_id/cancel"} {
		t.Run(path, func(t *testing.T) {
			stub := &mPriceServiceControllerStub{}
			app, controller := newMPriceControllerTestApp(stub)
			registerMPriceRoute(app, fiber.MethodPatch, path, controller.Cancel)

			realPath := strings.Replace(path, ":price_id", "price-902", 1)
			res, err := app.Test(httptest.NewRequest("PATCH", realPath, nil), -1)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if res.StatusCode != fiber.StatusOK {
				t.Fatalf("expected status %d, got %d", fiber.StatusOK, res.StatusCode)
			}
			if stub.capturedCancel.PriceID != "price-902" || stub.capturedCancel.CustID != "C22001" || stub.capturedCancel.ParentCustID != "P22001" {
				t.Fatalf("unexpected cancel params: %+v", stub.capturedCancel)
			}
			if stub.capturedCancel.UpdatedByID == nil || *stub.capturedCancel.UpdatedByID != 99 || stub.capturedCancel.UpdatedBy != "Price User" {
				t.Fatalf("expected cancel audit values, got %+v", stub.capturedCancel)
			}
		})
	}
}

func TestMPriceController_Template_ShouldReturnGeneratedFile(t *testing.T) {
	stub := &mPriceServiceControllerStub{}
	app, controller := newMPriceControllerTestApp(stub)
	registerMPriceRoute(app, fiber.MethodGet, "/v1/prices/template", controller.Template)

	res, err := app.Test(httptest.NewRequest("GET", "/v1/prices/template?format=csv", nil), -1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.StatusCode != fiber.StatusOK {
		t.Fatalf("expected status %d, got %d", fiber.StatusOK, res.StatusCode)
	}
	if stub.capturedTemplateFormat != "csv" {
		t.Fatalf("expected format csv, got %s", stub.capturedTemplateFormat)
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("failed read response: %v", err)
	}
	if string(body) != "template-bytes" || res.Header.Get("Content-Type") != "text/csv" {
		t.Fatalf("unexpected template response: content_type=%s body=%s", res.Header.Get("Content-Type"), string(body))
	}
}

func TestMPriceController_Export_ShouldParseFiltersAndReturnGeneratedFile(t *testing.T) {
	stub := &mPriceServiceControllerStub{}
	app, controller := newMPriceControllerTestApp(stub)
	registerMPriceRoute(app, fiber.MethodGet, "/v1/prices/export", controller.Export)

	res, err := app.Test(httptest.NewRequest("GET", "/v1/prices/export?file_type=csv&status=1,10&distributor_id=67,68", nil), -1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.StatusCode != fiber.StatusOK {
		t.Fatalf("expected status %d, got %d", fiber.StatusOK, res.StatusCode)
	}
	if stub.capturedExportFilter.FileType != "csv" || len(stub.capturedExportFilter.Status) != 2 || len(stub.capturedExportFilter.DistributorIDs) != 2 {
		t.Fatalf("unexpected export filter: %+v", stub.capturedExportFilter)
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("failed read response: %v", err)
	}
	if string(body) != "export-bytes" || res.Header.Get("Content-Type") != "text/csv" {
		t.Fatalf("unexpected export response: content_type=%s body=%s", res.Header.Get("Content-Type"), string(body))
	}
}

func TestMPriceController_Import_ShouldValidateAndPassContext(t *testing.T) {
	stub := &mPriceServiceControllerStub{}
	app, controller := newMPriceControllerTestApp(stub)
	registerMPriceRoute(app, fiber.MethodPost, "/v1/prices/import", controller.Import)

	req := httptest.NewRequest("POST", "/v1/prices/import", strings.NewReader(`{"file_url":"https://example.test/price.csv"}`))
	req.Header.Set("Content-Type", "application/json")
	res, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.StatusCode != fiber.StatusOK {
		t.Fatalf("expected status %d, got %d", fiber.StatusOK, res.StatusCode)
	}
	if stub.capturedImportReq.FileURL != "https://example.test/price.csv" || stub.capturedImportCustID != "C22001" {
		t.Fatalf("unexpected import request: %+v cust=%s", stub.capturedImportReq, stub.capturedImportCustID)
	}
}

func TestMPriceController_Import_ShouldReturnServiceErrorWithResultData(t *testing.T) {
	stub := &mPriceServiceControllerStub{
		importFn: func(request entity.MPriceImportRequest, custID string, parentCustID string, userID int64, distributorID int64, userFullName string) (entity.MPriceImportResponse, error) {
			return entity.MPriceImportResponse{FileURL: request.FileURL, FailedRow: 1}, errors.New("invalid price import")
		},
	}
	app, controller := newMPriceControllerTestApp(stub)
	registerMPriceRoute(app, fiber.MethodPost, "/v1/prices/import", controller.Import)

	req := httptest.NewRequest("POST", "/v1/prices/import", strings.NewReader(`{"file_url":"https://example.test/bad.csv"}`))
	req.Header.Set("Content-Type", "application/json")
	res, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.StatusCode != fiber.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", fiber.StatusBadRequest, res.StatusCode)
	}
}
