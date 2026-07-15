package controller

import (
	"bytes"
	"database/sql"
	"master/entity"
	"master/pkg/validation"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
)

type workingDayCalendarServiceControllerStub struct {
	listFilter       entity.WorkingDayCalendarQueryFilter
	listCustID       string
	listParentCustID string
	listData         []entity.WorkingDayCalendarListItem
	listTotal        int
	listLastPage     int
	listErr          error

	createRequest        entity.CreateWorkingDayCalendarBody
	createCustID         string
	createParentCustID   string
	createUserID         int64
	createData           entity.WorkingDayCalendarDetailResponse
	createErr            error
	createCalled         bool
	detailID             int64
	detailCustID         string
	detailParentCustID   string
	detailData           entity.WorkingDayCalendarDetailResponse
	detailErr            error
	calendarID           int64
	calendarFilter       entity.WorkingDayCalendarViewFilter
	calendarCustID       string
	calendarParentCustID string
	calendarData         entity.WorkingDayCalendarViewResponse
	calendarErr          error
	templateID           int64
	templateFormat       string
	templateCustID       string
	templateParentCustID string
	templateBuffer       *bytes.Buffer
	templateContentType  string
	templateFilename     string
	templateErr          error
	importID             int64
	importRequest        entity.WorkingDayCalendarImportHolidayRequest
	importCustID         string
	importParentCustID   string
	importUserID         int64
	importData           entity.WorkingDayCalendarImportHolidayResponse
	importErr            error
}

func (s *workingDayCalendarServiceControllerStub) List(filter entity.WorkingDayCalendarQueryFilter, custID, parentCustID string) ([]entity.WorkingDayCalendarListItem, int, int, error) {
	s.listFilter = filter
	s.listCustID = custID
	s.listParentCustID = parentCustID
	return s.listData, s.listTotal, s.listLastPage, s.listErr
}

func (s *workingDayCalendarServiceControllerStub) Create(request entity.CreateWorkingDayCalendarBody, custID, parentCustID string, userID int64) (entity.WorkingDayCalendarDetailResponse, error) {
	s.createCalled = true
	s.createRequest = request
	s.createCustID = custID
	s.createParentCustID = parentCustID
	s.createUserID = userID
	return s.createData, s.createErr
}

func (s *workingDayCalendarServiceControllerStub) Detail(id int64, custID, parentCustID string) (entity.WorkingDayCalendarDetailResponse, error) {
	s.detailID = id
	s.detailCustID = custID
	s.detailParentCustID = parentCustID
	return s.detailData, s.detailErr
}

func (s *workingDayCalendarServiceControllerStub) Calendar(id int64, filter entity.WorkingDayCalendarViewFilter, custID, parentCustID string) (entity.WorkingDayCalendarViewResponse, error) {
	s.calendarID = id
	s.calendarFilter = filter
	s.calendarCustID = custID
	s.calendarParentCustID = parentCustID
	return s.calendarData, s.calendarErr
}

func (s *workingDayCalendarServiceControllerStub) DownloadHolidayTemplate(id int64, format, custID, parentCustID string) (*bytes.Buffer, string, string, error) {
	s.templateID = id
	s.templateFormat = format
	s.templateCustID = custID
	s.templateParentCustID = parentCustID
	if s.templateBuffer == nil {
		s.templateBuffer = bytes.NewBufferString("date,notes\n")
	}
	if s.templateContentType == "" {
		s.templateContentType = "text/csv"
	}
	if s.templateFilename == "" {
		s.templateFilename = "working_day_calendar_holidays.csv"
	}
	return s.templateBuffer, s.templateContentType, s.templateFilename, s.templateErr
}

func (s *workingDayCalendarServiceControllerStub) ImportHolidays(id int64, request entity.WorkingDayCalendarImportHolidayRequest, custID, parentCustID string, userID int64) (entity.WorkingDayCalendarImportHolidayResponse, error) {
	s.importID = id
	s.importRequest = request
	s.importCustID = custID
	s.importParentCustID = parentCustID
	s.importUserID = userID
	return s.importData, s.importErr
}

func TestWorkingDayCalendarControllerListDefaultsPagination(t *testing.T) {
	stub := &workingDayCalendarServiceControllerStub{
		listData: []entity.WorkingDayCalendarListItem{{
			WorkingDayCalendarID: 1,
			Title:                "Cycle A",
		}},
		listTotal:    1,
		listLastPage: 1,
	}
	app := newWorkingDayCalendarTestApp(stub)

	req := httptest.NewRequest("GET", "/v1/working-day-calendars", nil)
	res, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.StatusCode != fiber.StatusOK {
		t.Fatalf("expected status %d, got %d", fiber.StatusOK, res.StatusCode)
	}
	if stub.listFilter.Page != 1 || stub.listFilter.Limit != 10 {
		t.Fatalf("pagination defaults = page %d limit %d, want 1/10", stub.listFilter.Page, stub.listFilter.Limit)
	}
	if stub.listCustID != "D001" || stub.listParentCustID != "P001" {
		t.Fatalf("list locals = cust %q parent %q", stub.listCustID, stub.listParentCustID)
	}
}

func TestWorkingDayCalendarControllerCreatePassesRequestAndLocals(t *testing.T) {
	stub := &workingDayCalendarServiceControllerStub{
		createData: entity.WorkingDayCalendarDetailResponse{
			WorkingDayCalendarID: 11,
			Title:                "Cycle A",
		},
	}
	app := newWorkingDayCalendarTestApp(stub)

	body := `{"title":"Cycle A","start_date":"2026-01-01","number_of_weeks":4,"default_holidays":[0,6]}`
	req := httptest.NewRequest("POST", "/v1/working-day-calendars", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	res, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.StatusCode != fiber.StatusCreated {
		t.Fatalf("expected status %d, got %d", fiber.StatusCreated, res.StatusCode)
	}
	if !stub.createCalled {
		t.Fatalf("expected create service to be called")
	}
	if stub.createRequest.Title != "Cycle A" || stub.createRequest.NumberOfWeeks != 4 || len(stub.createRequest.DefaultHolidays) != 2 {
		t.Fatalf("create request = %+v", stub.createRequest)
	}
	if stub.createCustID != "D001" || stub.createParentCustID != "P001" || stub.createUserID != 77 {
		t.Fatalf("create locals = cust %q parent %q user %d", stub.createCustID, stub.createParentCustID, stub.createUserID)
	}
}

func TestWorkingDayCalendarControllerCreateValidationErrorDoesNotCallService(t *testing.T) {
	stub := &workingDayCalendarServiceControllerStub{}
	app := newWorkingDayCalendarTestApp(stub)

	body := `{"title":"","start_date":"2026-01-01","number_of_weeks":4}`
	req := httptest.NewRequest("POST", "/v1/working-day-calendars", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	res, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.StatusCode != fiber.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", fiber.StatusBadRequest, res.StatusCode)
	}
	if stub.createCalled {
		t.Fatalf("expected create service not to be called")
	}
}

func TestWorkingDayCalendarControllerDetailAndCalendar(t *testing.T) {
	stub := &workingDayCalendarServiceControllerStub{
		detailData: entity.WorkingDayCalendarDetailResponse{
			WorkingDayCalendarID: 12,
			Title:                "Cycle B",
		},
		calendarData: entity.WorkingDayCalendarViewResponse{
			WorkingDayCalendarID: 12,
			Title:                "Cycle B",
			View:                 "month",
		},
	}
	app := newWorkingDayCalendarTestApp(stub)

	res, err := app.Test(httptest.NewRequest("GET", "/v1/working-day-calendars/12", nil), -1)
	if err != nil {
		t.Fatalf("unexpected detail error: %v", err)
	}
	if res.StatusCode != fiber.StatusOK {
		t.Fatalf("expected detail status %d, got %d", fiber.StatusOK, res.StatusCode)
	}
	if stub.detailID != 12 || stub.detailCustID != "D001" || stub.detailParentCustID != "P001" {
		t.Fatalf("detail call = id %d cust %q parent %q", stub.detailID, stub.detailCustID, stub.detailParentCustID)
	}

	res, err = app.Test(httptest.NewRequest("GET", "/v1/working-day-calendars/12/calendar?view=year&year=2026", nil), -1)
	if err != nil {
		t.Fatalf("unexpected calendar error: %v", err)
	}
	if res.StatusCode != fiber.StatusOK {
		t.Fatalf("expected calendar status %d, got %d", fiber.StatusOK, res.StatusCode)
	}
	if stub.calendarID != 12 || stub.calendarFilter.View != "year" || stub.calendarFilter.Year != 2026 {
		t.Fatalf("calendar call = id %d filter %+v", stub.calendarID, stub.calendarFilter)
	}
}

func TestWorkingDayCalendarControllerTemplateAndImport(t *testing.T) {
	stub := &workingDayCalendarServiceControllerStub{
		importData: entity.WorkingDayCalendarImportHolidayResponse{
			FileURL:    "https://example.com/holidays.csv",
			TotalRow:   1,
			SuccessRow: 1,
		},
	}
	app := newWorkingDayCalendarTestApp(stub)

	res, err := app.Test(httptest.NewRequest("GET", "/v1/working-day-calendars/holidays/template?format=csv", nil), -1)
	if err != nil {
		t.Fatalf("unexpected blank template error: %v", err)
	}
	if res.StatusCode != fiber.StatusOK {
		t.Fatalf("expected blank template status %d, got %d", fiber.StatusOK, res.StatusCode)
	}
	if stub.templateID != 0 || stub.templateFormat != "csv" {
		t.Fatalf("blank template call id=%d format=%q", stub.templateID, stub.templateFormat)
	}

	res, err = app.Test(httptest.NewRequest("GET", "/v1/working-day-calendars/12/holidays/template?format=csv", nil), -1)
	if err != nil {
		t.Fatalf("unexpected template error: %v", err)
	}
	if res.StatusCode != fiber.StatusOK {
		t.Fatalf("expected template status %d, got %d", fiber.StatusOK, res.StatusCode)
	}
	if stub.templateID != 12 || stub.templateFormat != "csv" || stub.templateCustID != "D001" || stub.templateParentCustID != "P001" {
		t.Fatalf("template call id=%d format=%q cust=%q parent=%q", stub.templateID, stub.templateFormat, stub.templateCustID, stub.templateParentCustID)
	}

	body := `{"file_url":"https://example.com/holidays.csv"}`
	req := httptest.NewRequest("PUT", "/v1/working-day-calendars/12/holidays/import", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	res, err = app.Test(req, -1)
	if err != nil {
		t.Fatalf("unexpected import error: %v", err)
	}
	if res.StatusCode != fiber.StatusOK {
		t.Fatalf("expected import status %d, got %d", fiber.StatusOK, res.StatusCode)
	}
	if stub.importID != 12 || stub.importRequest.FileURL != "https://example.com/holidays.csv" || stub.importUserID != 77 {
		t.Fatalf("import call id=%d req=%+v user=%d", stub.importID, stub.importRequest, stub.importUserID)
	}
}

func TestWorkingDayCalendarControllerDetailNotFound(t *testing.T) {
	stub := &workingDayCalendarServiceControllerStub{detailErr: sql.ErrNoRows}
	app := newWorkingDayCalendarTestApp(stub)

	res, err := app.Test(httptest.NewRequest("GET", "/v1/working-day-calendars/99", nil), -1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.StatusCode != fiber.StatusNotFound {
		t.Fatalf("expected status %d, got %d", fiber.StatusNotFound, res.StatusCode)
	}
}

func newWorkingDayCalendarTestApp(stub *workingDayCalendarServiceControllerStub) *fiber.App {
	app := fiber.New()
	controller := NewWorkingDayCalendarController(stub, validation.NewValidator())
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("requestid", "req-123")
		c.Locals("cust_id", "D001")
		c.Locals("parent_cust_id", "P001")
		c.Locals("user_id", int64(77))
		return c.Next()
	})
	app.Get("/v1/working-day-calendars", controller.List)
	app.Post("/v1/working-day-calendars", controller.Create)
	app.Get("/v1/working-day-calendars/holidays/template", controller.DownloadBlankHolidayTemplate)
	app.Get("/v1/working-day-calendars/:working_day_calendar_id/calendar", controller.Calendar)
	app.Get("/v1/working-day-calendars/:working_day_calendar_id/holidays/template", controller.DownloadHolidayTemplate)
	app.Put("/v1/working-day-calendars/:working_day_calendar_id/holidays/import", controller.ImportHolidays)
	app.Get("/v1/working-day-calendars/:working_day_calendar_id", controller.Detail)
	return app
}
