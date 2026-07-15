package controller

import (
	"encoding/json"
	"master/entity"
	"master/pkg/validation"
	"master/service"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
)

type surveyServiceControllerStub struct{}

type surveyServiceInvalidSalesmanStub struct{}

type surveyServiceCaptureStub struct {
	request       entity.CreateSurveyBody
	updateRequest entity.UpdateSurveyBody
}

type surveyServiceDetailStub struct {
	response entity.SurveyDetailResponse
}

func (s *surveyServiceControllerStub) List(_ entity.SurveyQueryFilter, _ string) ([]entity.SurveyListResponse, int, int, error) {
	return nil, 0, 0, nil
}

func (s *surveyServiceControllerStub) Detail(_ int, _ string) (entity.SurveyDetailResponse, error) {
	return entity.SurveyDetailResponse{}, nil
}

func (s *surveyServiceControllerStub) Store(_ entity.CreateSurveyBody) error {
	return service.ErrSurveyTitleConflict
}

func (s *surveyServiceControllerStub) Update(_ int, _ entity.UpdateSurveyBody) error {
	return nil
}

func (s *surveyServiceControllerStub) Deactivate(_ int, _ entity.DeactivateSurveyBody) error {
	return nil
}

func (s *surveyServiceInvalidSalesmanStub) List(_ entity.SurveyQueryFilter, _ string) ([]entity.SurveyListResponse, int, int, error) {
	return nil, 0, 0, nil
}

func (s *surveyServiceInvalidSalesmanStub) Detail(_ int, _ string) (entity.SurveyDetailResponse, error) {
	return entity.SurveyDetailResponse{}, nil
}

func (s *surveyServiceInvalidSalesmanStub) Store(_ entity.CreateSurveyBody) error {
	return &service.SurveyInvalidSalesmenError{
		InvalidEmpID:    []int{458, 459},
		InvalidSalesman: []string{"Bagus Prima", "Erling Braut Caraka"},
	}
}

func (s *surveyServiceInvalidSalesmanStub) Update(_ int, _ entity.UpdateSurveyBody) error {
	return &service.SurveyInvalidSalesmenError{
		InvalidEmpID:    []int{458, 459},
		InvalidSalesman: []string{"Bagus Prima", "Erling Braut Caraka"},
	}
}

func (s *surveyServiceInvalidSalesmanStub) Deactivate(_ int, _ entity.DeactivateSurveyBody) error {
	return nil
}

func (s *surveyServiceCaptureStub) List(_ entity.SurveyQueryFilter, _ string) ([]entity.SurveyListResponse, int, int, error) {
	return nil, 0, 0, nil
}

func (s *surveyServiceCaptureStub) Detail(_ int, _ string) (entity.SurveyDetailResponse, error) {
	return entity.SurveyDetailResponse{}, nil
}

func (s *surveyServiceCaptureStub) Store(request entity.CreateSurveyBody) error {
	s.request = request
	return nil
}

func (s *surveyServiceCaptureStub) Update(_ int, request entity.UpdateSurveyBody) error {
	s.updateRequest = request
	return nil
}

func (s *surveyServiceCaptureStub) Deactivate(_ int, _ entity.DeactivateSurveyBody) error {
	return nil
}

func (s *surveyServiceDetailStub) List(_ entity.SurveyQueryFilter, _ string) ([]entity.SurveyListResponse, int, int, error) {
	return nil, 0, 0, nil
}

func (s *surveyServiceDetailStub) Detail(_ int, _ string) (entity.SurveyDetailResponse, error) {
	return s.response, nil
}

func (s *surveyServiceDetailStub) Store(_ entity.CreateSurveyBody) error {
	return nil
}

func (s *surveyServiceDetailStub) Update(_ int, _ entity.UpdateSurveyBody) error {
	return nil
}

func (s *surveyServiceDetailStub) Deactivate(_ int, _ entity.DeactivateSurveyBody) error {
	return nil
}

func TestSurveyController_Create_ShouldReturnConflict409_WhenDuplicateTitleOverlap(t *testing.T) {
	app := fiber.New()
	v := validation.NewValidator()
	controller := NewSurveyController(&surveyServiceControllerStub{}, v)

	app.Post("/v1/survey", func(c *fiber.Ctx) error {
		c.Locals("requestid", "req-123")
		c.Locals("cust_id", "C1001")
		c.Locals("parent_cust_id", "C1001")
		c.Locals("user_id", int64(10))
		return controller.Create(c)
	})

	body := `{
		"survey_title":"Sales Visit",
		"efective_date_start":"2026-01-10",
		"efective_date_end":"2026-01-20",
		"answer_frequency":"One Time",
		"response_type":"Mandatory",
		"survey_template_id":[1]
	}`

	req := httptest.NewRequest("POST", "/v1/survey", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	res, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if res.StatusCode != fiber.StatusConflict {
		t.Fatalf("expected status %d, got %d", fiber.StatusConflict, res.StatusCode)
	}

	var payload struct {
		Message string `json:"message"`
	}
	if err := json.NewDecoder(res.Body).Decode(&payload); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}

	if payload.Message != service.ErrSurveyTitleConflict.Error() {
		t.Fatalf("expected message %q, got %q", service.ErrSurveyTitleConflict.Error(), payload.Message)
	}
}

func TestSurveyController_Create_ShouldParseDistributorAndEmpArrays(t *testing.T) {
	app := fiber.New()
	v := validation.NewValidator()
	serviceStub := &surveyServiceCaptureStub{}
	controller := NewSurveyController(serviceStub, v)

	app.Post("/v1/survey", func(c *fiber.Ctx) error {
		c.Locals("requestid", "req-123")
		c.Locals("cust_id", "C260020001")
		c.Locals("parent_cust_id", "C26002")
		c.Locals("user_id", int64(10))
		return controller.Create(c)
	})

	body := `{
		"survey_title":"survey salesman",
		"efective_date_start":"2026-04-01",
		"efective_date_end":"2026-04-18",
		"answer_frequency":"One Time",
		"response_type":"Mandatory",
		"target_type":"Specific",
		"distributor_id":[102],
		"area_id":[88],
		"outlet_id":[],
		"survey_template_id":40,
		"emp_id":[421,422]
	}`

	req := httptest.NewRequest("POST", "/v1/survey", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	res, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if res.StatusCode != fiber.StatusCreated {
		t.Fatalf("expected status %d, got %d", fiber.StatusCreated, res.StatusCode)
	}

	if len(serviceStub.request.DistributorId) != 1 || serviceStub.request.DistributorId[0] != 102 {
		t.Fatalf("expected distributor_id [102], got %v", serviceStub.request.DistributorId)
	}
	if len(serviceStub.request.EmpId) != 2 || serviceStub.request.EmpId[0] != 421 || serviceStub.request.EmpId[1] != 422 {
		t.Fatalf("expected emp_id [421 422], got %v", serviceStub.request.EmpId)
	}
	if serviceStub.request.CustId != "C260020001" {
		t.Fatalf("expected cust_id C260020001, got %s", serviceStub.request.CustId)
	}
	if serviceStub.request.ParentCustId != "C26002" {
		t.Fatalf("expected parent_cust_id C26002, got %s", serviceStub.request.ParentCustId)
	}
}

func TestSurveyController_Create_ShouldReturnDetailedInvalidSalesmanErrors(t *testing.T) {
	app := fiber.New()
	v := validation.NewValidator()
	controller := NewSurveyController(&surveyServiceInvalidSalesmanStub{}, v)

	app.Post("/v1/survey", func(c *fiber.Ctx) error {
		c.Locals("requestid", "req-123")
		c.Locals("cust_id", "C26002")
		c.Locals("parent_cust_id", "C26002")
		c.Locals("user_id", int64(140))
		return controller.Create(c)
	})

	body := `{
		"survey_title":"Testing May 12",
		"efective_date_start":"2026-05-12",
		"efective_date_end":"2026-05-13",
		"answer_frequency":"One Time",
		"response_type":"Optional",
		"target_type":"Specific",
		"distributor_id":[0,102,103,119],
		"area_id":[91,88],
		"outlet_id":[],
		"survey_template_id":53,
		"emp_id":[450,435,415,421,458,459,466]
	}`

	req := httptest.NewRequest("POST", "/v1/survey", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	res, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if res.StatusCode != fiber.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", fiber.StatusBadRequest, res.StatusCode)
	}

	var payload struct {
		Message string         `json:"message"`
		Errors  map[string]any `json:"errors"`
	}
	if err := json.NewDecoder(res.Body).Decode(&payload); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}

	if payload.Message != service.ErrSurveySalesmanNotFound.Error() {
		t.Fatalf("expected message %q, got %q", service.ErrSurveySalesmanNotFound.Error(), payload.Message)
	}

	invalidEmpIDs, ok := payload.Errors["invalid_emp_id"].([]interface{})
	if !ok || len(invalidEmpIDs) != 2 {
		t.Fatalf("expected invalid_emp_id array with 2 entries, got %+v", payload.Errors["invalid_emp_id"])
	}
	if invalidEmpIDs[0].(float64) != 458 || invalidEmpIDs[1].(float64) != 459 {
		t.Fatalf("unexpected invalid_emp_id payload: %+v", invalidEmpIDs)
	}

	invalidSalesmen, ok := payload.Errors["invalid_salesman"].([]interface{})
	if !ok || len(invalidSalesmen) != 2 {
		t.Fatalf("expected invalid_salesman array with 2 entries, got %+v", payload.Errors["invalid_salesman"])
	}
	if invalidSalesmen[0].(string) != "Bagus Prima" || invalidSalesmen[1].(string) != "Erling Braut Caraka" {
		t.Fatalf("unexpected invalid_salesman payload: %+v", invalidSalesmen)
	}
}

func TestSurveyController_Update_ShouldReturnDetailedInvalidSalesmanErrors(t *testing.T) {
	app := fiber.New()
	v := validation.NewValidator()
	controller := NewSurveyController(&surveyServiceInvalidSalesmanStub{}, v)

	app.Put("/v1/survey/:survey_id", func(c *fiber.Ctx) error {
		c.Locals("requestid", "req-123")
		c.Locals("cust_id", "C26002")
		c.Locals("parent_cust_id", "C26002")
		c.Locals("user_id", int64(140))
		return controller.Update(c)
	})

	body := `{
		"survey_title":"Testing May 12",
		"efective_date_start":"2026-05-12",
		"efective_date_end":"2026-05-13",
		"answer_frequency":"One Time",
		"response_type":"Optional",
		"target_type":"Specific",
		"distributor_id":[0,102,103,119],
		"area_id":[91,88],
		"outlet_id":[],
		"survey_template_id":53,
		"emp_id":[450,435,415,421,458,459,466]
	}`

	req := httptest.NewRequest("PUT", "/v1/survey/123", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	res, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if res.StatusCode != fiber.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", fiber.StatusBadRequest, res.StatusCode)
	}

	var payload struct {
		Message string         `json:"message"`
		Errors  map[string]any `json:"errors"`
	}
	if err := json.NewDecoder(res.Body).Decode(&payload); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}

	if payload.Message != service.ErrSurveySalesmanNotFound.Error() {
		t.Fatalf("expected message %q, got %q", service.ErrSurveySalesmanNotFound.Error(), payload.Message)
	}

	invalidEmpIDs, ok := payload.Errors["invalid_emp_id"].([]interface{})
	if !ok || len(invalidEmpIDs) != 2 {
		t.Fatalf("expected invalid_emp_id array with 2 entries, got %+v", payload.Errors["invalid_emp_id"])
	}
	if invalidEmpIDs[0].(float64) != 458 || invalidEmpIDs[1].(float64) != 459 {
		t.Fatalf("unexpected invalid_emp_id payload: %+v", invalidEmpIDs)
	}

	invalidSalesmen, ok := payload.Errors["invalid_salesman"].([]interface{})
	if !ok || len(invalidSalesmen) != 2 {
		t.Fatalf("expected invalid_salesman array with 2 entries, got %+v", payload.Errors["invalid_salesman"])
	}
	if invalidSalesmen[0].(string) != "Bagus Prima" || invalidSalesmen[1].(string) != "Erling Braut Caraka" {
		t.Fatalf("unexpected invalid_salesman payload: %+v", invalidSalesmen)
	}
}

func TestSurveyController_Create_ShouldAcceptNewAnswerFrequencyValues(t *testing.T) {
	values := []string{"Multiple Times, One Day", "Multiple Times, Different Day"}
	for _, value := range values {
		t.Run(value, func(t *testing.T) {
			app := fiber.New()
			v := validation.NewValidator()
			serviceStub := &surveyServiceCaptureStub{}
			controller := NewSurveyController(serviceStub, v)

			app.Post("/v1/survey", func(c *fiber.Ctx) error {
				c.Locals("requestid", "req-123")
				c.Locals("cust_id", "C260020001")
				c.Locals("parent_cust_id", "C26002")
				c.Locals("user_id", int64(10))
				return controller.Create(c)
			})

			body := `{"survey_title":"survey salesman","efective_date_start":"2026-04-01","efective_date_end":"2026-04-18","answer_frequency":"` + value + `","response_type":"Mandatory","survey_template_id":[40]}`
			req := httptest.NewRequest("POST", "/v1/survey", strings.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			res, err := app.Test(req, -1)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if res.StatusCode != fiber.StatusCreated {
				t.Fatalf("expected status %d, got %d", fiber.StatusCreated, res.StatusCode)
			}
			if serviceStub.request.AnswerFrequency != value {
				t.Fatalf("expected answer_frequency %q, got %q", value, serviceStub.request.AnswerFrequency)
			}
		})
	}
}

func TestSurveyController_Create_ShouldRejectLegacyMultipleAnswerFrequency(t *testing.T) {
	app := fiber.New()
	v := validation.NewValidator()
	serviceStub := &surveyServiceCaptureStub{}
	controller := NewSurveyController(serviceStub, v)

	app.Post("/v1/survey", func(c *fiber.Ctx) error {
		c.Locals("requestid", "req-123")
		c.Locals("cust_id", "C260020001")
		c.Locals("parent_cust_id", "C26002")
		c.Locals("user_id", int64(10))
		return controller.Create(c)
	})

	body := `{"survey_title":"survey salesman","efective_date_start":"2026-04-01","efective_date_end":"2026-04-18","answer_frequency":"Multiple","response_type":"Mandatory","survey_template_id":[40]}`
	req := httptest.NewRequest("POST", "/v1/survey", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	res, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.StatusCode != fiber.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", fiber.StatusBadRequest, res.StatusCode)
	}
}

func TestSurveyController_Update_ShouldAcceptNewAnswerFrequencyValues(t *testing.T) {
	values := []string{"Multiple Times, One Day", "Multiple Times, Different Day"}
	for _, value := range values {
		t.Run(value, func(t *testing.T) {
			app := fiber.New()
			v := validation.NewValidator()
			serviceStub := &surveyServiceCaptureStub{}
			controller := NewSurveyController(serviceStub, v)

			app.Put("/v1/survey/:survey_id", func(c *fiber.Ctx) error {
				c.Locals("requestid", "req-123")
				c.Locals("cust_id", "C260020001")
				c.Locals("parent_cust_id", "C26002")
				c.Locals("user_id", int64(10))
				return controller.Update(c)
			})

			body := `{"survey_title":"survey salesman","efective_date_start":"2026-04-01","efective_date_end":"2026-04-18","answer_frequency":"` + value + `","response_type":"Mandatory","survey_template_id":[40]}`
			req := httptest.NewRequest("PUT", "/v1/survey/99", strings.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			res, err := app.Test(req, -1)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if res.StatusCode != fiber.StatusOK {
				t.Fatalf("expected status %d, got %d", fiber.StatusOK, res.StatusCode)
			}
			if serviceStub.updateRequest.AnswerFrequency != value {
				t.Fatalf("expected answer_frequency %q, got %q", value, serviceStub.updateRequest.AnswerFrequency)
			}
		})
	}
}

func TestSurveyController_Update_ShouldRejectLegacyMultipleAnswerFrequency(t *testing.T) {
	app := fiber.New()
	v := validation.NewValidator()
	serviceStub := &surveyServiceCaptureStub{}
	controller := NewSurveyController(serviceStub, v)

	app.Put("/v1/survey/:survey_id", func(c *fiber.Ctx) error {
		c.Locals("requestid", "req-123")
		c.Locals("cust_id", "C260020001")
		c.Locals("parent_cust_id", "C26002")
		c.Locals("user_id", int64(10))
		return controller.Update(c)
	})

	body := `{"survey_title":"survey salesman","efective_date_start":"2026-04-01","efective_date_end":"2026-04-18","answer_frequency":"Multiple","response_type":"Mandatory","survey_template_id":[40]}`
	req := httptest.NewRequest("PUT", "/v1/survey/99", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	res, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.StatusCode != fiber.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", fiber.StatusBadRequest, res.StatusCode)
	}
}

func TestSurveyController_Create_ShouldParseLevelTargetAndTargetDistributor(t *testing.T) {
	app := fiber.New()
	v := validation.NewValidator()
	serviceStub := &surveyServiceCaptureStub{}
	controller := NewSurveyController(serviceStub, v)

	app.Post("/v1/survey", func(c *fiber.Ctx) error {
		c.Locals("requestid", "req-123")
		c.Locals("cust_id", "C22001")
		c.Locals("parent_cust_id", "C22001")
		c.Locals("user_id", int64(10))
		return controller.Create(c)
	})

	body := `{
		"survey_title":"survey distributor",
		"efective_date_start":"2026-07-01",
		"efective_date_end":"2026-07-30",
		"answer_frequency":"One Time",
		"response_type":"Mandatory",
		"level_target":"Distributor",
		"target_cust_id":"C22001",
		"distributor_id":[0],
		"area_id":[82],
		"target_distributor_id":[120, 0, 121],
		"outlet_id":[],
		"survey_template_id":53
	}`

	req := httptest.NewRequest("POST", "/v1/survey", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	res, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.StatusCode != fiber.StatusCreated {
		t.Fatalf("expected status %d, got %d", fiber.StatusCreated, res.StatusCode)
	}
	if serviceStub.request.LevelTarget != "Distributor" {
		t.Fatalf("expected level_target Distributor, got %q", serviceStub.request.LevelTarget)
	}
	if serviceStub.request.TargetCustId != "C22001" {
		t.Fatalf("expected target_cust_id C22001, got %q", serviceStub.request.TargetCustId)
	}
	if len(serviceStub.request.TargetDistributorId) != 3 || serviceStub.request.TargetDistributorId[2] != 121 {
		t.Fatalf("expected target_distributor_id [120 0 121], got %v", serviceStub.request.TargetDistributorId)
	}
}

func TestSurveyController_Update_ShouldParseLevelTargetAndTargetDistributor(t *testing.T) {
	app := fiber.New()
	v := validation.NewValidator()
	serviceStub := &surveyServiceCaptureStub{}
	controller := NewSurveyController(serviceStub, v)

	app.Put("/v1/survey/:survey_id", func(c *fiber.Ctx) error {
		c.Locals("requestid", "req-123")
		c.Locals("cust_id", "C22001")
		c.Locals("parent_cust_id", "C22001")
		c.Locals("user_id", int64(10))
		return controller.Update(c)
	})

	body := `{
		"survey_title":"survey distributor",
		"efective_date_start":"2026-07-01",
		"efective_date_end":"2026-07-30",
		"answer_frequency":"One Time",
		"response_type":"Mandatory",
		"level_target":"Outlet",
		"target_cust_id":"C22001",
		"distributor_id":[67],
		"area_id":[82],
		"target_distributor_id":[],
		"outlet_id":[3489],
		"survey_template_id":53
	}`

	req := httptest.NewRequest("PUT", "/v1/survey/107", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	res, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.StatusCode != fiber.StatusOK {
		t.Fatalf("expected status %d, got %d", fiber.StatusOK, res.StatusCode)
	}
	if serviceStub.updateRequest.LevelTarget != "Outlet" || serviceStub.updateRequest.TargetCustId != "C22001" {
		t.Fatalf("expected parsed level_target/target_cust_id, got %+v", serviceStub.updateRequest)
	}
}

func TestSurveyController_Create_ShouldRejectInvalidLevelTarget(t *testing.T) {
	app := fiber.New()
	v := validation.NewValidator()
	serviceStub := &surveyServiceCaptureStub{}
	controller := NewSurveyController(serviceStub, v)

	app.Post("/v1/survey", func(c *fiber.Ctx) error {
		c.Locals("requestid", "req-123")
		c.Locals("cust_id", "C22001")
		c.Locals("parent_cust_id", "C22001")
		c.Locals("user_id", int64(10))
		return controller.Create(c)
	})

	body := `{
		"survey_title":"survey invalid level",
		"efective_date_start":"2026-07-01",
		"efective_date_end":"2026-07-30",
		"answer_frequency":"One Time",
		"response_type":"Mandatory",
		"level_target":"Manager",
		"target_cust_id":"C22001",
		"distributor_id":[0],
		"area_id":[82],
		"survey_template_id":53
	}`

	req := httptest.NewRequest("POST", "/v1/survey", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	res, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.StatusCode != fiber.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", fiber.StatusBadRequest, res.StatusCode)
	}
}

func TestSurveyController_Detail_ShouldReturnSurveyDetailContractFields(t *testing.T) {
	app := fiber.New()
	v := validation.NewValidator()
	distributorId := 102
	salesTeamId := 30
	otClassId := 1
	otGrpId := 2
	otTypeId := 3
	serviceStub := &surveyServiceDetailStub{
		response: entity.SurveyDetailResponse{
			SurveyId:        80,
			DistributorId:   entity.FlexibleIntArray{distributorId},
			AreaId:          entity.FlexibleIntArray{88},
			DistributorCode: "DIST001",
			DistributorName: "Distributor One",
			BusinessUnits: []entity.SurveyBusinessUnit{
				{DistributorId: distributorId, AreaId: 88, BusinessUnitName: "Distributor One", Name: "Distributor One", Type: "distributor"},
			},
			Outlet: []entity.SurveyOutletResponse{
				{
					SurveyOutletId: 1,
					OutletId:       10,
					OutletCode:     "OUT001",
					OutletName:     "Outlet One",
					OtClassId:      &otClassId,
					OtClassName:    "Class A",
					OtGrpId:        &otGrpId,
					OtGrpName:      "Group B",
					OtTypeId:       &otTypeId,
					OtTypeName:     "Type C",
				},
			},
			Salesman: []entity.SurveySalesmanResponse{
				{
					MSurveySalesmanId: 2,
					SalesId:           20,
					SalesTeamId:       &salesTeamId,
					SalesTeamName:     "Team Alpha",
					SalesName:         "Budi Sales",
				},
			},
			TargetSurvey: &entity.SurveyTargetResponse{
				Area:     []entity.SurveyAreaResponse{},
				Outlet:   []entity.SurveyOutletResponse{},
				Salesman: []entity.SurveySalesmanResponse{},
			},
			Template: []entity.SurveyTemplateNested{},
		},
	}
	controller := NewSurveyController(serviceStub, v)

	app.Get("/v1/survey/:survey_id", func(c *fiber.Ctx) error {
		c.Locals("requestid", "req-123")
		c.Locals("cust_id", "C260020001")
		return controller.Detail(c)
	})

	req := httptest.NewRequest("GET", "/v1/survey/80", nil)
	res, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.StatusCode != fiber.StatusOK {
		t.Fatalf("expected status %d, got %d", fiber.StatusOK, res.StatusCode)
	}

	var payload struct {
		Data entity.SurveyDetailResponse `json:"data"`
	}
	if err := json.NewDecoder(res.Body).Decode(&payload); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}

	if len(payload.Data.DistributorId) != 1 || payload.Data.DistributorId[0] != 102 || payload.Data.DistributorCode != "DIST001" || payload.Data.DistributorName != "Distributor One" {
		t.Fatalf("unexpected distributor response: %+v", payload.Data)
	}
	if len(payload.Data.AreaId) != 1 || payload.Data.AreaId[0] != 88 {
		t.Fatalf("unexpected area_id response: %+v", payload.Data.AreaId)
	}
	if len(payload.Data.BusinessUnits) != 1 || payload.Data.BusinessUnits[0].DistributorId != 102 || payload.Data.BusinessUnits[0].AreaId != 88 {
		t.Fatalf("unexpected business_units response: %+v", payload.Data.BusinessUnits)
	}
	if payload.Data.BusinessUnits[0].BusinessUnitName != "Distributor One" || payload.Data.BusinessUnits[0].Name != "Distributor One" || payload.Data.BusinessUnits[0].Type != "distributor" {
		t.Fatalf("unexpected business unit labels: %+v", payload.Data.BusinessUnits[0])
	}

	serviceStub.response.BusinessUnits = []entity.SurveyBusinessUnit{{DistributorId: 0, AreaId: 88, BusinessUnitName: "Principal", Name: "Principal", Type: "principal"}}
	serviceStub.response.DistributorId = entity.FlexibleIntArray{0}
	res, err = app.Test(httptest.NewRequest("GET", "/v1/survey/80", nil), -1)
	if err != nil {
		t.Fatalf("unexpected error on principal detail: %v", err)
	}
	if res.StatusCode != fiber.StatusOK {
		t.Fatalf("expected status %d, got %d", fiber.StatusOK, res.StatusCode)
	}
	if err := json.NewDecoder(res.Body).Decode(&payload); err != nil {
		t.Fatalf("failed to decode principal response body: %v", err)
	}
	if len(payload.Data.DistributorId) != 1 || payload.Data.DistributorId[0] != 0 {
		t.Fatalf("unexpected principal distributor response: %+v", payload.Data.DistributorId)
	}
	if len(payload.Data.BusinessUnits) != 1 || payload.Data.BusinessUnits[0].DistributorId != 0 || payload.Data.BusinessUnits[0].Name != "Principal" || payload.Data.BusinessUnits[0].Type != "principal" {
		t.Fatalf("unexpected principal business_units response: %+v", payload.Data.BusinessUnits)
	}
	if len(payload.Data.Outlet) != 1 || payload.Data.Outlet[0].SurveyOutletId != 1 || payload.Data.Outlet[0].OutletCode != "OUT001" {
		t.Fatalf("unexpected outlet response: %+v", payload.Data.Outlet)
	}
	if len(payload.Data.Salesman) != 1 || payload.Data.Salesman[0].MSurveySalesmanId != 2 || payload.Data.Salesman[0].SalesName != "Budi Sales" {
		t.Fatalf("unexpected salesman response: %+v", payload.Data.Salesman)
	}
}
