package controller

import (
	"encoding/json"
	"master/entity"
	"master/pkg/validation"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
)

type surveyTemplateServiceControllerStub struct {
	storeReq         entity.CreateSurveyTemplateBody
	updateReq        entity.UpdateSurveyTemplateBody
	storeCalled      bool
	updateCalled     bool
	storeErr         error
	updateErr        error
	deleteErr        error
	listData         []entity.SurveyTemplateListResponse
	listTotal        int
	listLastPage     int
	listErr          error
	detailData       entity.SurveyTemplateDetailResponse
	detailErr        error
	deleteCalled     bool
	deleteCustID     string
	deleteTemplateID int
	deleteUserID     int64
}

func (s *surveyTemplateServiceControllerStub) List(_ entity.SurveyTemplateQueryFilter, _ string) ([]entity.SurveyTemplateListResponse, int, int, error) {
	return s.listData, s.listTotal, s.listLastPage, s.listErr
}

func (s *surveyTemplateServiceControllerStub) Detail(_ int, _ string) (entity.SurveyTemplateDetailResponse, error) {
	return s.detailData, s.detailErr
}

func (s *surveyTemplateServiceControllerStub) Store(req entity.CreateSurveyTemplateBody) error {
	s.storeCalled = true
	s.storeReq = req
	return s.storeErr
}

func (s *surveyTemplateServiceControllerStub) Update(_ int, req entity.UpdateSurveyTemplateBody) error {
	s.updateCalled = true
	s.updateReq = req
	return s.updateErr
}

func (s *surveyTemplateServiceControllerStub) Delete(custId string, surveyTemplateId int, userId int64) error {
	s.deleteCalled = true
	s.deleteCustID = custId
	s.deleteTemplateID = surveyTemplateId
	s.deleteUserID = userId
	return s.deleteErr
}

func TestSurveyTemplateController_Create_ShouldPassInputTypeToService(t *testing.T) {
	app := fiber.New()
	v := validation.NewValidator()
	stub := &surveyTemplateServiceControllerStub{}
	controller := NewSurveyTemplateController(stub, v)

	app.Post("/v1/survey_template", func(c *fiber.Ctx) error {
		c.Locals("requestid", "req-123")
		c.Locals("cust_id", "C22001")
		c.Locals("user_id", int64(10))
		return controller.Create(c)
	})

	body := `{
		"template_title":"Survey Kepuasan Pelanggan",
		"question_total":1,
		"use_image":true,
		"is_active":true,
		"question":[
			{
				"question":"Bagaimana layanan kami?",
				"input_type":"dropdown",
				"answer_type":"Single",
				"q_option":[{"option":"Baik"},{"option":"Cukup"}]
			}
		]
	}`

	req := httptest.NewRequest("POST", "/v1/survey_template", strings.NewReader(body))
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
	if len(stub.storeReq.Question) != 1 {
		t.Fatalf("expected 1 question, got %d", len(stub.storeReq.Question))
	}
	if stub.storeReq.Question[0].InputType != "dropdown" {
		t.Fatalf("expected input_type dropdown, got %s", stub.storeReq.Question[0].InputType)
	}
}

func TestSurveyTemplateController_Update_ShouldPassInputTypeToService(t *testing.T) {
	app := fiber.New()
	v := validation.NewValidator()
	stub := &surveyTemplateServiceControllerStub{}
	controller := NewSurveyTemplateController(stub, v)

	app.Put("/v1/survey_template/:survey_template_id", func(c *fiber.Ctx) error {
		c.Locals("requestid", "req-123")
		c.Locals("cust_id", "C22001")
		c.Locals("user_id", int64(10))
		return controller.Update(c)
	})

	body := `{
		"template_title":"Survey Kepuasan Pelanggan Update",
		"question_total":1,
		"use_image":false,
		"is_active":true,
		"question":[
			{
				"question_template_id":23,
				"question":"Bagaimana layanan kami sekarang?",
				"input_type":"radiobutton",
				"answer_type":"Single",
				"q_option":[{"option":"Baik"}]
			}
		]
	}`

	req := httptest.NewRequest("PUT", "/v1/survey_template/12", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	res, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if res.StatusCode != fiber.StatusOK {
		t.Fatalf("expected status %d, got %d", fiber.StatusOK, res.StatusCode)
	}
	if !stub.updateCalled {
		t.Fatalf("expected update to be called")
	}
	if len(stub.updateReq.Question) != 1 {
		t.Fatalf("expected 1 question, got %d", len(stub.updateReq.Question))
	}
	if stub.updateReq.Question[0].InputType != "radiobutton" {
		t.Fatalf("expected input_type radiobutton, got %s", stub.updateReq.Question[0].InputType)
	}
}

func TestSurveyTemplateController_Create_ShouldReturnBadRequest_WhenInputTypeInvalid(t *testing.T) {
	app := fiber.New()
	v := validation.NewValidator()
	stub := &surveyTemplateServiceControllerStub{}
	controller := NewSurveyTemplateController(stub, v)

	app.Post("/v1/survey_template", func(c *fiber.Ctx) error {
		c.Locals("requestid", "req-123")
		c.Locals("cust_id", "C22001")
		c.Locals("user_id", int64(10))
		return controller.Create(c)
	})

	body := `{
		"template_title":"Survey Kepuasan Pelanggan",
		"question_total":1,
		"use_image":true,
		"is_active":true,
		"question":[
			{
				"question":"Bagaimana layanan kami?",
				"input_type":"slider",
				"answer_type":"Single",
				"q_option":[]
			}
		]
	}`

	req := httptest.NewRequest("POST", "/v1/survey_template", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	res, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if res.StatusCode != fiber.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", fiber.StatusBadRequest, res.StatusCode)
	}
	if stub.storeCalled {
		t.Fatalf("expected store not to be called on invalid validation")
	}

	var payload struct {
		Message string `json:"message"`
	}
	if err := json.NewDecoder(res.Body).Decode(&payload); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}
	if payload.Message == "" {
		t.Fatalf("expected non-empty message")
	}
}

func TestSurveyTemplateController_Detail_ShouldIncludeInputTypeInResponseData(t *testing.T) {
	app := fiber.New()
	v := validation.NewValidator()
	stub := &surveyTemplateServiceControllerStub{
		detailData: entity.SurveyTemplateDetailResponse{
			SurveyTemplateId: 12,
			TemplateCode:     "TMP001",
			TemplateTitle:    "Template",
			QuestionTotal:    1,
			UseImage:         true,
			IsActive:         true,
			CreatedAt:        func() *time.Time { n := time.Now().UTC(); return &n }(),
			QuestionTemplate: []entity.QuestionTemplateResponse{
				{
					QuestionTemplateId: 23,
					SurveyTemplateId:   12,
					Question:           "Q1",
					InputType:          "toggle",
					AnswerType:         "Single",
					MQOptionTemplate:   []entity.QOptionTemplateResponse{},
				},
			},
		},
	}
	controller := NewSurveyTemplateController(stub, v)

	app.Get("/v1/survey_template/:survey_template_id", func(c *fiber.Ctx) error {
		c.Locals("requestid", "req-123")
		c.Locals("cust_id", "C22001")
		c.Locals("user_id", int64(10))
		return controller.Detail(c)
	})

	req := httptest.NewRequest("GET", "/v1/survey_template/12", nil)
	res, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.StatusCode != fiber.StatusOK {
		t.Fatalf("expected status %d, got %d", fiber.StatusOK, res.StatusCode)
	}

	var payload struct {
		Data struct {
			QuestionTemplate []struct {
				InputType string `json:"input_type"`
			} `json:"question_template"`
		} `json:"data"`
	}
	if err := json.NewDecoder(res.Body).Decode(&payload); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}
	if len(payload.Data.QuestionTemplate) != 1 {
		t.Fatalf("expected 1 question template, got %d", len(payload.Data.QuestionTemplate))
	}
	if payload.Data.QuestionTemplate[0].InputType != "toggle" {
		t.Fatalf("expected input_type toggle, got %s", payload.Data.QuestionTemplate[0].InputType)
	}
}
