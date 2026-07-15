package controller

import (
	"database/sql"
	"errors"
	"master/entity"
	"master/pkg/constant"
	"master/pkg/middleware"
	"master/pkg/responsebuild"
	"master/pkg/str"
	"master/pkg/validation"
	"master/service"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

type SurveyReportController struct {
	SurveyReportService service.SurveyReportService
	validator           *validation.Validate
}

const (
	defaultSurveyReportPage  = 1
	defaultSurveyReportLimit = 5
	defaultSurveyReportSort  = "created_date:desc"
)

func NewSurveyReportController(surveyReportService service.SurveyReportService, validator *validation.Validate) SurveyReportController {
	return SurveyReportController{
		SurveyReportService: surveyReportService,
		validator:           validator,
	}
}

func (controller *SurveyReportController) Route(app *fiber.App) {
	routeV1 := app.Group("/v1/survey-report", middleware.JWTProtected())
	routeV1.Get("", controller.List)
	routeV1.Get("/export", controller.Export)
	routeV1.Get("/:survey_answer_id", controller.Detail)
}

func (controller *SurveyReportController) List(c *fiber.Ctx) error {
	headerAcceptLang := getAcceptLanguage(c)
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	filter, err := parseSurveyReportFilter(c)
	if err != nil {
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}
	applySurveyReportPagingDefaults(&filter)
	filter.CustID = c.Locals("cust_id").(string)

	data, total, lastPage, err := controller.SurveyReportService.List(filter)
	if err != nil {
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	if len(data) == 0 {
		responsePayload.Setmsg("No Data")
		responsePayload.Setdata([]entity.SurveyReportListResponse{})
	} else {
		responsePayload.Setmsg("Success")
		responsePayload.Setdata(data)
	}
	responsePayload.Setpaging(entity.Pagination{
		TotalRecord: total,
		PageCurrent: filter.Page,
		PageLimit:   filter.Limit,
		PageTotal:   lastPage,
	})

	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}

func (controller *SurveyReportController) Detail(c *fiber.Ctx) error {
	headerAcceptLang := getAcceptLanguage(c)
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	params := entity.SurveyReportParams{}
	if err := c.ParamsParser(&params); err != nil {
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	errs := controller.validator.ValidateStruct(params, headerAcceptLang)
	if errs != nil {
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	data, err := controller.SurveyReportService.Detail(params.SurveyAnswerID, c.Locals("cust_id").(string))
	if err != nil {
		statusCode := fiber.StatusBadRequest
		errMsg := err.Error()
		if errors.Is(err, sql.ErrNoRows) {
			statusCode = fiber.StatusNotFound
			errMsg = "record not found"
		}
		responsePayload.Setmsg(errMsg)
		return c.Status(statusCode).JSON(responsePayload.GetRespPayload())
	}

	responsePayload.Setmsg("Success")
	responsePayload.Setdata(data)
	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}

func (controller *SurveyReportController) Export(c *fiber.Ctx) error {
	headerAcceptLang := getAcceptLanguage(c)
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	filter, err := parseSurveyReportFilter(c)
	if err != nil {
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}
	applySurveyReportPagingDefaults(&filter)
	filter.CustID = c.Locals("cust_id").(string)

	createdBy := resolveSurveyReportCreatedBy(c)

	data, err := controller.SurveyReportService.Export(filter, createdBy)
	if err != nil {
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	if strings.EqualFold(c.Query("download"), "true") {
		fileBytes, decodeErr := data.DecodeFileBase64()
		if decodeErr != nil {
			responsePayload.Setmsg(decodeErr.Error())
			return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
		}
		filename := data.ReportName
		if filename == "" {
			filename = "survey-report"
		}
		c.Set(fiber.HeaderContentType, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
		c.Set(fiber.HeaderContentDisposition, `attachment; filename="`+filename+`.xlsx"`)
		return c.Status(fiber.StatusOK).Send(fileBytes)
	}

	responsePayload.Setmsg("Success")
	responsePayload.Setdata(data)
	responsePayload.Setpaging(entity.Pagination{
		TotalRecord: 1,
		PageCurrent: 0,
		PageLimit:   0,
		PageTotal:   1,
	})
	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}

func parseSurveyReportFilter(c *fiber.Ctx) (entity.SurveyReportQueryFilter, error) {
	var filter entity.SurveyReportQueryFilter
	if err := c.QueryParser(&filter); err != nil {
		return filter, err
	}

	filter.SurveyID = parseInt64QueryArray(c, "survey_id")
	filter.SurveyTitle = parseStringQueryArray(c, "survey_title")
	filter.AreaID = parseInt64QueryArray(c, "area_id")
	startDate, err := parseSurveyReportDate(c.Query("start_date"))
	if err != nil {
		return filter, err
	}
	endDate, err := parseSurveyReportDate(c.Query("end_date"))
	if err != nil {
		return filter, err
	}
	filter.StartDate = startDate
	filter.EndDate = endDate

	return filter, nil
}

func parseSurveyReportDate(rawDate string) (*time.Time, error) {
	rawDate = strings.TrimSpace(rawDate)
	if rawDate == "" {
		return nil, nil
	}

	rfc3339Date, err := str.DateStrToRfc3339String(rawDate)
	if err == nil {
		parsedDate, err := time.Parse(time.RFC3339, rfc3339Date)
		if err != nil {
			return nil, err
		}
		return &parsedDate, nil
	}

	unixTimestamp, parseErr := strconv.ParseInt(rawDate, 10, 64)
	if parseErr != nil {
		return nil, err
	}

	parsedUnixDate := time.Unix(unixTimestamp, 0).UTC()
	parsedDate := time.Date(parsedUnixDate.Year(), parsedUnixDate.Month(), parsedUnixDate.Day(), 0, 0, 0, 0, time.UTC)
	return &parsedDate, nil
}

func applySurveyReportPagingDefaults(filter *entity.SurveyReportQueryFilter) {
	if filter.Page < 1 {
		filter.Page = defaultSurveyReportPage
	}
	if filter.Limit < 1 {
		filter.Limit = defaultSurveyReportLimit
	}
	if filter.Sort == "" {
		filter.Sort = defaultSurveyReportSort
	}
}

func resolveSurveyReportCreatedBy(c *fiber.Ctx) string {
	if value, ok := c.Locals("user_fullname").(string); ok && value != "" {
		return value
	}

	if value, ok := c.Locals("user_name").(string); ok && value != "" {
		return value
	}

	return "system"
}

func parseInt64QueryArray(c *fiber.Ctx, key string) []int64 {
	values := c.Context().QueryArgs().PeekMulti(key)
	valuesBracket := c.Context().QueryArgs().PeekMulti(key + "[]")
	rawValues := make([]string, 0, len(values)+len(valuesBracket)+1)

	for _, value := range values {
		rawValues = append(rawValues, string(value))
	}
	for _, value := range valuesBracket {
		rawValues = append(rawValues, string(value))
	}
	if len(rawValues) == 0 {
		fallback := c.Query(key)
		if fallback != "" {
			rawValues = append(rawValues, fallback)
		}
	}

	results := make([]int64, 0)
	for _, raw := range rawValues {
		parts := strings.Split(raw, ",")
		for _, part := range parts {
			trimmed := strings.TrimSpace(part)
			if trimmed == "" {
				continue
			}
			parsed, err := strconv.ParseInt(trimmed, 10, 64)
			if err != nil {
				continue
			}
			results = append(results, parsed)
		}
	}

	return results
}

func parseStringQueryArray(c *fiber.Ctx, key string) []string {
	values := c.Context().QueryArgs().PeekMulti(key)
	valuesBracket := c.Context().QueryArgs().PeekMulti(key + "[]")
	rawValues := make([]string, 0, len(values)+len(valuesBracket)+1)

	for _, value := range values {
		rawValues = append(rawValues, string(value))
	}
	for _, value := range valuesBracket {
		rawValues = append(rawValues, string(value))
	}
	if len(rawValues) == 0 {
		fallback := c.Query(key)
		if fallback != "" {
			rawValues = append(rawValues, fallback)
		}
	}

	results := make([]string, 0)
	seen := make(map[string]struct{})
	for _, raw := range rawValues {
		parts := strings.Split(raw, ",")
		for _, part := range parts {
			trimmed := strings.TrimSpace(part)
			if trimmed == "" {
				continue
			}
			if _, ok := seen[trimmed]; ok {
				continue
			}
			seen[trimmed] = struct{}{}
			results = append(results, trimmed)
		}
	}

	return results
}

func getAcceptLanguage(c *fiber.Ctx) string {
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		return c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	return ""
}
