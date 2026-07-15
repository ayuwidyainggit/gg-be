package controller

import (
	"mobile/entity"
	"mobile/model"
	"mobile/pkg/constant"
	"mobile/pkg/middleware"
	"mobile/pkg/responsebuild"
	"mobile/pkg/validation"
	"mobile/service"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

type SurveyController struct {
	SurveyService service.SurveyService
	validator     *validation.Validate
}

func NewSurveyController(surveyService service.SurveyService, validator *validation.Validate) *SurveyController {
	return &SurveyController{
		SurveyService: surveyService,
		validator:     validator,
	}
}

func (controller *SurveyController) Route(app *fiber.App) {
	surveyRouteV1 := app.Group("/v1/survey-list", middleware.JWTProtected())
	surveyRouteV1.Get("", controller.List)

	salesmanSurveyV1 := app.Group("/v1/salesman-survey", middleware.JWTProtected())
	salesmanSurveyV1.Post("", controller.SubmitSurvey)
	salesmanSurveyV1.Get("/:survey_id", controller.GetSubmittedSurvey)

	surveyDetailRouteV1 := app.Group("/v1/survey-detail", middleware.JWTProtected())
	surveyDetailRouteV1.Get("/:id", controller.Detail)

	surveyAnswerRouteV1 := app.Group("/v1/survey-answer-list", middleware.JWTProtected())
	surveyAnswerRouteV1.Get("", controller.ListSurveyAnswer)
}

func (controller *SurveyController) List(c *fiber.Ctx) error {
	var dataFilter entity.SurveyQueryFilter
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), constant.HEADER_ACCEPT_LANG)

	if err := c.QueryParser(&dataFilter); err != nil {
		log.Error("SurveyController, List, query parser filter:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	errs := controller.validator.ValidateStruct(dataFilter, constant.HEADER_ACCEPT_LANG)
	if errs != nil {
		log.Error("SurveyController, List, ValidateStruct(params), errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	distributorID, _ := c.Locals("distributor_id").(int64)
	empID, _ := c.Locals("emp_id").(int64)
	dataFilter.DistributorID = distributorID
	dataFilter.EmpID = empID

	data, total, lastPage, err := controller.SurveyService.List(c.Context(), dataFilter)
	if err != nil {
		log.Error("SurveyController, List, data, err:", err.Error())
		statusCode := fiber.StatusBadRequest
		if err.Error() == "sql: no rows in result set" {
			statusCode = fiber.StatusNotFound
		}
		responsePayload.Setmsg(err.Error())
		return c.Status(statusCode).JSON(responsePayload.GetRespPayload())
	}

	if data == nil {
		data = make([]model.SurveyResponse, 0)
	}

	responsePayload.Setmsg("Success")
	responsePayload.Setdata(data)
	responsePayload.Setpaging(entity.Pagination{
		TotalRecord: total,
		PageCurrent: int(dataFilter.Page),
		PageLimit:   int(dataFilter.Limit),
		PageTotal:   int(lastPage),
	})

	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}

func (controller *SurveyController) Detail(c *fiber.Ctx) error {
	surveyId, err := c.ParamsInt("id")
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), constant.HEADER_ACCEPT_LANG)

	if err != nil {
		log.Error("SurveyController, Detail, id parser:", err.Error())
		responsePayload.Setmsg("Invalid survey ID")
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	data, err := controller.SurveyService.Detail(c.Context(), int64(surveyId))
	if err != nil {
		log.Error("SurveyController, Detail, data, err:", err.Error())
		statusCode := fiber.StatusBadRequest
		if err.Error() == "sql: no rows in result set" {
			statusCode = fiber.StatusNotFound
		}
		responsePayload.Setmsg(err.Error())
		return c.Status(statusCode).JSON(responsePayload.GetRespPayload())
	}

	responsePayload.Setmsg("Success")
	responsePayload.Setdata(data)
	responsePayload.Setpaging(entity.Pagination{}) // standard empty pagination for detail

	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}

func (controller *SurveyController) SubmitSurvey(c *fiber.Ctx) error {
	requestID := c.Locals("requestid").(string)
	responsePayload := responsebuild.BuildResponse(requestID, constant.HEADER_ACCEPT_LANG)

	var req entity.SubmitSurveyRequest
	if err := c.BodyParser(&req); err != nil {
		log.Error("SurveyController, Create, body parser:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	errs := controller.validator.ValidateStruct(req, constant.HEADER_ACCEPT_LANG)
	if errs != nil {
		log.Error("SurveyController, Create, ValidateStruct(params), errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	custID, _ := c.Locals("cust_id").(string)
	if custID == "" {
		responsePayload.Setmsg("customer ID must be filled")
		return c.Status(fiber.StatusUnauthorized).JSON(responsePayload.GetRespPayload())
	}

	userID := c.Locals("user_id").(int64)
	if userID == 0 {
		responsePayload.Setmsg("user ID must be filled")
		return c.Status(fiber.StatusUnauthorized).JSON(responsePayload.GetRespPayload())
	}

	empID := c.Locals("emp_id").(int64)
	if empID == 0 {
		responsePayload.Setmsg("emp ID must be filled")
		return c.Status(fiber.StatusUnauthorized).JSON(responsePayload.GetRespPayload())
	}

	distributorID, _ := c.Locals("distributor_id").(int64)

	req.CustID = custID
	req.UserID = userID
	req.EmpID = empID
	req.DistributorID = distributorID

	err := controller.SurveyService.SubmitSurvey(c.Context(), req)
	if err != nil {
		log.Error("SurveyController, Create, service:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(responsePayload.GetRespPayload())
	}

	responsePayload.Setmsg("Success")
	responsePayload.Setdata(req)
	responsePayload.Setpaging(entity.Pagination{})

	return c.Status(fiber.StatusCreated).JSON(responsePayload.GetRespPayload())
}

func (controller *SurveyController) GetSubmittedSurvey(c *fiber.Ctx) error {
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), constant.HEADER_ACCEPT_LANG)
	surveyID, err := c.ParamsInt("survey_id")
	if err != nil {
		log.Error("SurveyController, GetSubmittedSurvey, survey_id from params:")
		responsePayload.Setmsg("survey_id is required")
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	surveyAnswerIDStr := c.Query("survey_answer_id")
	if surveyAnswerIDStr == "" {
		log.Error("SurveyController, GetSubmittedSurvey, survey_answer_id from query:")
		responsePayload.Setmsg("survey_answer_id is required")
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	surveyAnswerID, err := strconv.ParseInt(surveyAnswerIDStr, 10, 64)
	if err != nil {
		log.Error("SurveyController, GetSubmittedSurvey, survey_answer_id parser:", err.Error())
		responsePayload.Setmsg("Invalid survey_answer_id")
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	data, err := controller.SurveyService.GetSubmittedSurvey(c.Context(), int64(surveyID), surveyAnswerID)
	if err != nil {
		log.Error("SurveyController, GetSubmittedSurvey, service:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(responsePayload.GetRespPayload())
	}

	responsePayload.Setmsg("Success")
	responsePayload.Setdata(data)
	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}

func (controller *SurveyController) ListSurveyAnswer(c *fiber.Ctx) error {
	var dataFilter entity.SurveyAnswerListFilter
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), constant.HEADER_ACCEPT_LANG)

	if err := c.QueryParser(&dataFilter); err != nil {
		log.Error("SurveyController, ListSurveyAnswer, query parser filter:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	errs := controller.validator.ValidateStruct(dataFilter, constant.HEADER_ACCEPT_LANG)
	if errs != nil {
		log.Error("SurveyController, ListSurveyAnswer, ValidateStruct(params), errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	// Default pagination values
	if dataFilter.Page <= 0 {
		dataFilter.Page = 1
	}
	if dataFilter.Limit <= 0 {
		dataFilter.Limit = 10
	}

	data, total, lastPage, err := controller.SurveyService.ListSurveyAnswer(c.Context(), dataFilter.SurveyID, dataFilter.OutletID, dataFilter.Page, dataFilter.Limit)
	if err != nil {
		log.Error("SurveyController, ListSurveyAnswer, data, err:", err.Error())
		statusCode := fiber.StatusBadRequest
		if err.Error() == "sql: no rows in result set" {
			statusCode = fiber.StatusNotFound
		}
		responsePayload.Setmsg(err.Error())
		return c.Status(statusCode).JSON(responsePayload.GetRespPayload())
	}

	if data == nil {
		data = make([]model.SurveyAnswerListItem, 0)
	}

	responsePayload.Setmsg("Success")
	responsePayload.Setdata(data)
	responsePayload.Setpaging(entity.Pagination{
		TotalRecord: total,
		PageCurrent: dataFilter.Page,
		PageLimit:   dataFilter.Limit,
		PageTotal:   int(lastPage),
	})

	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}
