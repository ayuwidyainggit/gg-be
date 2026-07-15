package controller

import (
	"encoding/json"
	"io"
	"net/http/httptest"
	"strings"
	"testing"

	"master/entity"
	"master/pkg/constant"
	"master/pkg/validation"

	"github.com/gofiber/fiber/v2"
)

type distributorServiceStub struct {
	capturedFilter entity.DistributorQueryFilter
	capturedDetail entity.DetailDistributorParams
	capturedUpdate entity.UpdateDistributorRequest
	updateErr      error
}

func (s *distributorServiceStub) Store(request entity.CreateDistributorBody) (entity.DistributorResponse, error) {
	return entity.DistributorResponse{}, nil
}

func (s *distributorServiceStub) List(dataFilter entity.DistributorQueryFilter, custId string) ([]entity.DistributorListRespone, int, int, error) {
	s.capturedFilter = dataFilter
	return []entity.DistributorListRespone{}, 0, 0, nil
}

func (s *distributorServiceStub) LookupList(dataFilter entity.DistributorQueryFilter, custId string) ([]entity.DistributorLookupResponse, int, int, error) {
	s.capturedFilter = dataFilter
	return []entity.DistributorLookupResponse{}, 0, 0, nil
}

func (s *distributorServiceStub) Detail(params entity.DetailDistributorParams) (entity.DistributorResponse, error) {
	s.capturedDetail = params
	return entity.DistributorResponse{}, nil
}

func (s *distributorServiceStub) Update(distributorId int64, request entity.UpdateDistributorRequest) error {
	s.capturedUpdate = request
	return s.updateErr
}

func (s *distributorServiceStub) Delete(custId string, distributorId int64, userId int64) error {
	return nil
}

func (s *distributorServiceStub) ListWithCustomer(dataFilter entity.DistributorQueryFilter, custId string) ([]entity.DistributorCustomerResp, int, int, error) {
	return []entity.DistributorCustomerResp{}, 0, 0, nil
}

func strPtrDistributorController(v string) *string { return &v }

func TestNormalizeUpdateDistributorRequest_EmptyBarcodePreservedForNullPatch(t *testing.T) {
	request := entity.UpdateDistributorRequest{
		Barcode: strPtrDistributorController(""),
	}

	normalizeUpdateDistributorRequest(&request)

	if request.Barcode != nil {
		t.Fatalf("expected empty barcode to become nil during normalization")
	}
}

func TestNormalizeUpdateDistributorRequest_NonEmptyBarcodePreserved(t *testing.T) {
	request := entity.UpdateDistributorRequest{
		Barcode: strPtrDistributorController("123456"),
	}

	normalizeUpdateDistributorRequest(&request)

	if request.Barcode == nil || *request.Barcode != "123456" {
		t.Fatalf("expected non-empty barcode to be preserved")
	}
}

func TestNormalizeUpdateDistributorRequest_EmptyZipCodePreservedForNullPatch(t *testing.T) {
	request := entity.UpdateDistributorRequest{
		ZipCode: strPtrDistributorController(""),
	}

	normalizeUpdateDistributorRequest(&request)

	if request.ZipCode != nil {
		t.Fatalf("expected empty zip_code to become nil during normalization")
	}
}

func TestPreserveNullableDistributorStringPatch_EmptyStringRestoresPointer(t *testing.T) {
	raw := map[string]json.RawMessage{
		"zip_code": json.RawMessage(`""`),
	}

	var zipCode *string
	preserveNullableDistributorStringPatch(raw, "zip_code", &zipCode)

	if zipCode == nil {
		t.Fatalf("expected empty zip_code to be restored for null patch")
	}
	if *zipCode != "" {
		t.Fatalf("expected restored zip_code to be empty string")
	}
}

func TestMarkNullableDistributorStringPresence_ReturnsTrueWhenKeyExists(t *testing.T) {
	raw := map[string]json.RawMessage{
		"barcode": json.RawMessage(`""`),
	}

	if !markNullableDistributorStringPresence(raw, "barcode") {
		t.Fatalf("expected barcode presence to be detected")
	}
	if markNullableDistributorStringPresence(raw, "zip_code") {
		t.Fatalf("expected zip_code presence to be false when omitted")
	}
}

func TestPreserveNullableDistributorIntPatch_NullRestoresZeroSentinel(t *testing.T) {
	raw := map[string]json.RawMessage{
		"ot_loc_id": json.RawMessage(`null`),
	}

	var otLocID *int
	preserveNullableDistributorIntPatch(raw, "ot_loc_id", &otLocID)

	if otLocID == nil {
		t.Fatalf("expected ot_loc_id sentinel to be restored for null patch")
	}
	if *otLocID != 0 {
		t.Fatalf("expected ot_loc_id sentinel to be zero, got %d", *otLocID)
	}
}

func TestNormalizeUpdateDistributorRequest_ContactEmailNilBecomesEmptyString(t *testing.T) {
	request := entity.UpdateDistributorRequest{
		Contacts: []entity.DistributorContactUpdate{
			{},
		},
	}

	normalizeUpdateDistributorRequest(&request)

	if len(request.Contacts) != 1 {
		t.Fatalf("expected exactly one contact")
	}
	if request.Contacts[0].Email == nil {
		t.Fatalf("expected nil email to become empty string pointer")
	}
	if *request.Contacts[0].Email != "" {
		t.Fatalf("expected normalized email value to be empty string")
	}
}

func TestDistributorController_List_AssignsPrincipalScopeFromJWT(t *testing.T) {
	svc := &distributorServiceStub{}
	controller := &DistributorController{DistributorService: svc}

	app := fiber.New()
	app.Get("/v1/distributors", func(c *fiber.Ctx) error {
		c.Locals("requestid", "req-1")
		c.Locals("cust_id", "P22001")
		c.Locals("parent_cust_id", "P22001")
		c.Locals("distributor_id", int64(0))
		return controller.List(c)
	})

	req := httptest.NewRequest("GET", "/v1/distributors?page=1&limit=10", nil)
	res, err := app.Test(req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if res.StatusCode != fiber.StatusOK {
		t.Fatalf("expected status 200, got %d", res.StatusCode)
	}

	if svc.capturedFilter.CustId != "P22001" {
		t.Fatalf("expected cust_id P22001, got %s", svc.capturedFilter.CustId)
	}

	if svc.capturedFilter.ParentCustId != "P22001" {
		t.Fatalf("expected parent_cust_id P22001, got %s", svc.capturedFilter.ParentCustId)
	}

	if svc.capturedFilter.JwtDistributorId != 0 {
		t.Fatalf("expected principal distributor_id 0, got %d", svc.capturedFilter.JwtDistributorId)
	}
}

func TestDistributorController_Update_ReturnsNotFoundWhenServiceReportsMissingDistributor(t *testing.T) {
	svc := &distributorServiceStub{updateErr: constant.ErrNoRowsAffected}
	controller := &DistributorController{DistributorService: svc, validator: validation.NewValidator()}

	app := fiber.New()
	app.Patch("/v1/distributors/:distributor_id", func(c *fiber.Ctx) error {
		c.Locals("requestid", "req-update-not-found")
		c.Locals("cust_id", "C22001")
		c.Locals("user_id", int64(99))
		return controller.Update(c)
	})

	body := `{"distributor_name":"PT Besi Makmur"}`
	req := httptest.NewRequest("PATCH", "/v1/distributors/102", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	res, err := app.Test(req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if res.StatusCode != fiber.StatusNotFound {
		t.Fatalf("expected status 404, got %d", res.StatusCode)
	}
}

func TestDistributorController_Update_AllowsEmptyZipCode(t *testing.T) {
	svc := &distributorServiceStub{}
	controller := &DistributorController{DistributorService: svc, validator: validation.NewValidator()}

	app := fiber.New()
	app.Patch("/v1/distributors/:distributor_id", func(c *fiber.Ctx) error {
		c.Locals("requestid", "req-update-zip")
		c.Locals("cust_id", "C22001")
		c.Locals("user_id", int64(99))
		return controller.Update(c)
	})

	body := `{"zip_code":"","distributor_name":"PT Besi Makmur"}`
	req := httptest.NewRequest("PATCH", "/v1/distributors/102", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	res, err := app.Test(req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if res.StatusCode != fiber.StatusOK {
		responseBody, _ := io.ReadAll(res.Body)
		t.Fatalf("expected status 200, got %d with body %s", res.StatusCode, string(responseBody))
	}

	if svc.capturedUpdate.ZipCode == nil {
		t.Fatalf("expected empty zip_code to stay present so repository can patch null")
	}
	if *svc.capturedUpdate.ZipCode != "" {
		t.Fatalf("expected empty zip_code to be preserved, got %q", *svc.capturedUpdate.ZipCode)
	}
	if !svc.capturedUpdate.ZipCodeProvided {
		t.Fatalf("expected zip_code presence flag to be set")
	}
	if svc.capturedUpdate.CustId != "C22001" {
		t.Fatalf("expected cust_id C22001, got %s", svc.capturedUpdate.CustId)
	}
	if svc.capturedUpdate.UpdatedBy == nil || *svc.capturedUpdate.UpdatedBy != 99 {
		t.Fatalf("expected updated_by 99, got %v", svc.capturedUpdate.UpdatedBy)
	}
}

func TestDistributorController_Update_AllowsEmptyBarcode(t *testing.T) {
	svc := &distributorServiceStub{}
	controller := &DistributorController{DistributorService: svc, validator: validation.NewValidator()}

	app := fiber.New()
	app.Patch("/v1/distributors/:distributor_id", func(c *fiber.Ctx) error {
		c.Locals("requestid", "req-update-barcode")
		c.Locals("cust_id", "C22001")
		c.Locals("user_id", int64(99))
		return controller.Update(c)
	})

	body := `{"barcode":"","distributor_name":"PT Besi Makmur"}`
	req := httptest.NewRequest("PATCH", "/v1/distributors/102", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	res, err := app.Test(req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if res.StatusCode != fiber.StatusOK {
		responseBody, _ := io.ReadAll(res.Body)
		t.Fatalf("expected status 200, got %d with body %s", res.StatusCode, string(responseBody))
	}

	if svc.capturedUpdate.Barcode == nil {
		t.Fatalf("expected empty barcode to stay present so repository can patch null")
	}
	if *svc.capturedUpdate.Barcode != "" {
		t.Fatalf("expected empty barcode to be preserved, got %q", *svc.capturedUpdate.Barcode)
	}
	if !svc.capturedUpdate.BarcodeProvided {
		t.Fatalf("expected barcode presence flag to be set")
	}
}

func TestDistributorController_Update_AllowsNullOptionalLocationAndContactFields(t *testing.T) {
	svc := &distributorServiceStub{}
	controller := &DistributorController{DistributorService: svc, validator: validation.NewValidator()}

	app := fiber.New()
	app.Patch("/v1/distributors/:distributor_id", func(c *fiber.Ctx) error {
		c.Locals("requestid", "req-update-nullables")
		c.Locals("cust_id", "C22001")
		c.Locals("user_id", int64(99))
		return controller.Update(c)
	})

	body := `{"province_id":null,"regency_id":null,"sub_district_id":null,"ward_id":null,"ot_loc_id":null,"phone":"","fax_number":"","distributor_name":"PT Besi Makmur"}`
	req := httptest.NewRequest("PATCH", "/v1/distributors/102", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	res, err := app.Test(req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if res.StatusCode != fiber.StatusOK {
		responseBody, _ := io.ReadAll(res.Body)
		t.Fatalf("expected status 200, got %d with body %s", res.StatusCode, string(responseBody))
	}

	if !svc.capturedUpdate.ProvinceIdProvided || svc.capturedUpdate.ProvinceId == nil || *svc.capturedUpdate.ProvinceId != "" {
		t.Fatalf("expected province_id null patch sentinel to be preserved")
	}
	if !svc.capturedUpdate.RegencyIdProvided || svc.capturedUpdate.RegencyId == nil || *svc.capturedUpdate.RegencyId != "" {
		t.Fatalf("expected regency_id null patch sentinel to be preserved")
	}
	if !svc.capturedUpdate.SubDistrictIdProvided || svc.capturedUpdate.SubDistrictId == nil || *svc.capturedUpdate.SubDistrictId != "" {
		t.Fatalf("expected sub_district_id null patch sentinel to be preserved")
	}
	if !svc.capturedUpdate.WardIdProvided || svc.capturedUpdate.WardId == nil || *svc.capturedUpdate.WardId != "" {
		t.Fatalf("expected ward_id null patch sentinel to be preserved")
	}
	if !svc.capturedUpdate.OtLocIdProvided || svc.capturedUpdate.OtLocId == nil || *svc.capturedUpdate.OtLocId != 0 {
		t.Fatalf("expected ot_loc_id null patch sentinel to be preserved")
	}
	if !svc.capturedUpdate.PhoneProvided || svc.capturedUpdate.Phone == nil || *svc.capturedUpdate.Phone != "" {
		t.Fatalf("expected phone empty patch sentinel to be preserved")
	}
	if !svc.capturedUpdate.FaxNumberProvided || svc.capturedUpdate.FaxNumber == nil || *svc.capturedUpdate.FaxNumber != "" {
		t.Fatalf("expected fax_number empty patch sentinel to be preserved")
	}
}

func TestDistributorController_List_WithoutAcceptLanguageHeader_DoesNotPanic(t *testing.T) {
	svc := &distributorServiceStub{}
	controller := &DistributorController{DistributorService: svc}

	app := fiber.New()
	app.Get("/v1/distributors", func(c *fiber.Ctx) error {
		c.Locals("requestid", "req-no-lang")
		c.Locals("cust_id", "P22001")
		c.Locals("parent_cust_id", "P22001")
		c.Locals("distributor_id", int64(0))
		return controller.List(c)
	})

	req := httptest.NewRequest("GET", "/v1/distributors?page=1&limit=10", nil)
	res, err := app.Test(req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("expected response body to be readable, got %v", err)
	}

	if res.StatusCode != fiber.StatusOK {
		t.Fatalf("expected status 200, got %d with body %s", res.StatusCode, string(body))
	}

	if strings.Contains(strings.ToLower(string(body)), "panic") {
		t.Fatalf("expected response body without panic indication, got %s", string(body))
	}

	if svc.capturedFilter.CustId != "P22001" {
		t.Fatalf("expected cust_id P22001, got %s", svc.capturedFilter.CustId)
	}

	if svc.capturedFilter.ParentCustId != "P22001" {
		t.Fatalf("expected parent_cust_id P22001, got %s", svc.capturedFilter.ParentCustId)
	}
}

func TestDistributorController_List_AssignsDistributorScopeFromJWT(t *testing.T) {
	svc := &distributorServiceStub{}
	controller := &DistributorController{DistributorService: svc}

	app := fiber.New()
	app.Get("/v1/distributors", func(c *fiber.Ctx) error {
		c.Locals("requestid", "req-2")
		c.Locals("cust_id", "C22001")
		c.Locals("parent_cust_id", "P22001")
		c.Locals("distributor_id", int64(55))
		return controller.List(c)
	})

	req := httptest.NewRequest("GET", "/v1/distributors?page=1&limit=10", nil)
	res, err := app.Test(req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if res.StatusCode != fiber.StatusOK {
		t.Fatalf("expected status 200, got %d", res.StatusCode)
	}

	if svc.capturedFilter.CustId != "C22001" {
		t.Fatalf("expected cust_id C22001, got %s", svc.capturedFilter.CustId)
	}

	if svc.capturedFilter.ParentCustId != "P22001" {
		t.Fatalf("expected parent_cust_id P22001, got %s", svc.capturedFilter.ParentCustId)
	}

	if svc.capturedFilter.JwtDistributorId != 55 {
		t.Fatalf("expected distributor_id 55, got %d", svc.capturedFilter.JwtDistributorId)
	}
}

func TestDistributorController_Detail_AssignsDistributorScopeFromJWT(t *testing.T) {
	svc := &distributorServiceStub{}
	controller := &DistributorController{DistributorService: svc, validator: validation.NewValidator()}

	app := fiber.New()
	app.Get("/v1/distributors/:distributor_id", func(c *fiber.Ctx) error {
		c.Locals("requestid", "req-detail")
		c.Locals("cust_id", "C22001")
		c.Locals("parent_cust_id", "P22001")
		c.Locals("distributor_id", int64(55))
		return controller.Detail(c)
	})

	req := httptest.NewRequest("GET", "/v1/distributors/103", nil)
	res, err := app.Test(req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if res.StatusCode != fiber.StatusOK {
		t.Fatalf("expected status 200, got %d", res.StatusCode)
	}

	if svc.capturedDetail.CustId != "C22001" {
		t.Fatalf("expected cust_id C22001, got %s", svc.capturedDetail.CustId)
	}

	if svc.capturedDetail.ParentCustId != "P22001" {
		t.Fatalf("expected parent_cust_id P22001, got %s", svc.capturedDetail.ParentCustId)
	}

	if svc.capturedDetail.JwtDistributorId != 55 {
		t.Fatalf("expected distributor_id 55, got %d", svc.capturedDetail.JwtDistributorId)
	}

	if svc.capturedDetail.DistributorId != 103 {
		t.Fatalf("expected distributor_id param 103, got %d", svc.capturedDetail.DistributorId)
	}
}
