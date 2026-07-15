package controller

import (
	"database/sql"
	"encoding/json"
	"errors"
	"master/entity"
	"master/pkg/constant"
	"master/pkg/middleware"
	"master/pkg/responsebuild"
	"master/pkg/validation"
	"master/service"

	"github.com/gofiber/fiber/v2"
)

type SurveyController struct {
	SurveyService service.SurveyService
	validator     *validation.Validate
}

func NewSurveyController(surveyService service.SurveyService, validator *validation.Validate) SurveyController {
	return SurveyController{
		SurveyService: surveyService,
		validator:     validator,
	}
}

func (controller *SurveyController) Route(app *fiber.App) {
	surveyRouteV1 := app.Group("/v1/survey", middleware.JWTProtected())
	surveyRouteV1.Get("", controller.List)
	surveyRouteV1.Get("/:survey_id", controller.Detail)
	surveyRouteV1.Post("", controller.Create)
	surveyRouteV1.Put("/:survey_id", controller.Update)
	surveyRouteV1.Patch("/:survey_id", controller.Deactivate)
}

func (controller *SurveyController) List(c *fiber.Ctx) error {
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	var filter entity.SurveyQueryFilter
	if err := c.QueryParser(&filter); err != nil {
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.Limit < 1 {
		filter.Limit = 10
	}
	if filter.Limit > 9999 {
		filter.Limit = 9999
	}

	custId := c.Locals("cust_id").(string)

	data, total, lastPage, err := controller.SurveyService.List(filter, custId)
	if err != nil {
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	if len(data) == 0 {
		responsePayload.Setmsg("No Data")
		responsePayload.Setdata(nil)
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

func (controller *SurveyController) Detail(c *fiber.Ctx) error {
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	var params entity.SurveyParams
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

	custId := c.Locals("cust_id").(string)

	data, err := controller.SurveyService.Detail(params.SurveyId, custId)
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

func (controller *SurveyController) Create(c *fiber.Ctx) error {
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	var request entity.CreateSurveyBody
	// Use json.Unmarshal directly to support FlexibleIntArray custom unmarshaling
	if err := json.Unmarshal(c.Body(), &request); err != nil {
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	custId := c.Locals("cust_id").(string)
	parentCustId := c.Locals("parent_cust_id").(string)
	userId := c.Locals("user_id").(int64)

	request.CustId = custId
	request.ParentCustId = parentCustId
	request.CreatedBy = userId

	errs := controller.validator.ValidateStruct(request, headerAcceptLang)
	if errs != nil {
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	err := controller.SurveyService.Store(request)
	if err != nil {
		if invalidErrors, ok := service.BuildInvalidSalesmanErrors(err); ok {
			responsePayload.Setmsg(err.Error())
			responsePayload.Seterrors(invalidErrors)
			return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
		}
		if errors.Is(err, service.ErrSurveyTitleConflict) ||
			errors.Is(err, service.ErrSurveyAreaDistributorRequired) ||
			errors.Is(err, service.ErrSurveyAreaDistributorMismatch) ||
			errors.Is(err, service.ErrSurveySalesmanNotFound) ||
			errors.Is(err, service.ErrSurveyInvalidDateFormat) {
			responsePayload.Setmsg(err.Error())
			statusCode := fiber.StatusBadRequest
			if errors.Is(err, service.ErrSurveyTitleConflict) {
				statusCode = fiber.StatusConflict
			}
			return c.Status(statusCode).JSON(responsePayload.GetRespPayload())
		}
		responsePayload.Setmsg("Failed to create survey data")
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	responsePayload.Setmsg("Survey has been successfully created")
	return c.Status(fiber.StatusCreated).JSON(responsePayload.GetRespPayload())
}

func (controller *SurveyController) Update(c *fiber.Ctx) error {
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	var params entity.SurveyParams
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

	var request entity.UpdateSurveyBody
	// Use json.Unmarshal directly to support FlexibleIntArray custom unmarshaling
	if err := json.Unmarshal(c.Body(), &request); err != nil {
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	custId := c.Locals("cust_id").(string)
	parentCustId := c.Locals("parent_cust_id").(string)
	userId := c.Locals("user_id").(int64)

	request.CustId = custId
	request.ParentCustId = parentCustId
	request.UpdatedBy = userId

	errs = controller.validator.ValidateStruct(request, headerAcceptLang)
	if errs != nil {
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	err := controller.SurveyService.Update(params.SurveyId, request)
	if err != nil {
		if invalidErrors, ok := service.BuildInvalidSalesmanErrors(err); ok {
			responsePayload.Setmsg(err.Error())
			responsePayload.Seterrors(invalidErrors)
			return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
		}
		if errors.Is(err, service.ErrSurveyTitleConflict) ||
			errors.Is(err, service.ErrSurveyAreaDistributorRequired) ||
			errors.Is(err, service.ErrSurveyAreaDistributorMismatch) ||
			errors.Is(err, service.ErrSurveySalesmanNotFound) ||
			errors.Is(err, service.ErrSurveyInvalidDateFormat) {
			responsePayload.Setmsg(err.Error())
			statusCode := fiber.StatusBadRequest
			if errors.Is(err, service.ErrSurveyTitleConflict) {
				statusCode = fiber.StatusConflict
			}
			return c.Status(statusCode).JSON(responsePayload.GetRespPayload())
		}
		responsePayload.Setmsg("Failed to update survey data")
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	responsePayload.Setmsg("Survey has been successfully updated")
	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}

func (controller *SurveyController) Deactivate(c *fiber.Ctx) error {
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	var params entity.SurveyParams
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

	var request entity.DeactivateSurveyBody
	if err := c.BodyParser(&request); err != nil {
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	custId := c.Locals("cust_id").(string)
	userId := c.Locals("user_id").(int64)

	request.CustId = custId
	request.UpdatedBy = userId

	err := controller.SurveyService.Deactivate(params.SurveyId, request)
	if err != nil {
		responsePayload.Setmsg("Survey not found")
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	msg := "Survey successfully deactivated"
	if request.IsActive {
		msg = "Survey successfully activated"
	}
	responsePayload.Setmsg(msg)
	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}
