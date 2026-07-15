package controller

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"master/entity"
	"master/pkg/constant"
	"master/pkg/middleware"
	"master/pkg/responsebuild"
	"master/pkg/validation"
	"master/service"

	"github.com/gofiber/fiber/v2"
)

type WorkingDayCalendarController struct {
	WorkingDayCalendarService service.WorkingDayCalendarService
	validator                 *validation.Validate
}

func NewWorkingDayCalendarController(workingDayCalendarService service.WorkingDayCalendarService, validator *validation.Validate) *WorkingDayCalendarController {
	return &WorkingDayCalendarController{
		WorkingDayCalendarService: workingDayCalendarService,
		validator:                 validator,
	}
}

func (controller *WorkingDayCalendarController) Route(app *fiber.App) {
	route := app.Group("/v1/working-day-calendars", middleware.JWTProtected())
	route.Get("", controller.List)
	route.Post("", controller.Create)
	route.Get("/holidays/template", controller.DownloadBlankHolidayTemplate)
	route.Get("/:working_day_calendar_id/calendar", controller.Calendar)
	route.Get("/:working_day_calendar_id/holidays/template", controller.DownloadHolidayTemplate)
	route.Put("/:working_day_calendar_id/holidays/import", controller.ImportHolidays)
	route.Get("/:working_day_calendar_id", controller.Detail)
}

func (controller *WorkingDayCalendarController) getAcceptLang(c *fiber.Ctx) string {
	headers := c.GetReqHeaders()
	if val, ok := headers[constant.HEADER_ACCEPT_LANG]; ok {
		if len(val) > 0 && val[0] != "" {
			return val[0]
		}
	}
	return ""
}

func (controller *WorkingDayCalendarController) List(c *fiber.Ctx) error {
	var filter entity.WorkingDayCalendarQueryFilter
	headerAcceptLang := controller.getAcceptLang(c)
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	if err := c.QueryParser(&filter); err != nil {
		log.Println("WorkingDayCalendarController, List, QueryParser:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}
	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.Limit < 1 {
		filter.Limit = 10
	}

	errs := controller.validator.ValidateStruct(filter, headerAcceptLang)
	if errs != nil {
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	custID := c.Locals("cust_id").(string)
	parentCustID := c.Locals("parent_cust_id").(string)
	data, total, lastPage, err := controller.WorkingDayCalendarService.List(filter, custID, parentCustID)
	if err != nil {
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	if total == 0 {
		responsePayload.Setmsg(constant.NO_DATA)
	} else {
		responsePayload.Setmsg(constant.SUCCESS_NO_DATA_DISPLAYED)
	}
	responsePayload.Setdata(data)
	responsePayload.Setpaging(entity.Pagination{
		TotalRecord: total,
		PageCurrent: filter.Page,
		PageLimit:   filter.Limit,
		PageTotal:   lastPage,
	})
	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}

func (controller *WorkingDayCalendarController) Create(c *fiber.Ctx) error {
	var request entity.CreateWorkingDayCalendarBody
	headerAcceptLang := controller.getAcceptLang(c)
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	if err := c.BodyParser(&request); err != nil {
		log.Println("WorkingDayCalendarController, Create, BodyParser:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	errs := controller.validator.ValidateStruct(request, headerAcceptLang)
	if errs != nil {
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	custID := c.Locals("cust_id").(string)
	parentCustID := c.Locals("parent_cust_id").(string)
	userID := c.Locals("user_id").(int64)
	data, err := controller.WorkingDayCalendarService.Create(request, custID, parentCustID, userID)
	if err != nil {
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	responsePayload.Setmsg(constant.SUCCESSFULLY_ADDED)
	responsePayload.Setdata(data)
	return c.Status(fiber.StatusCreated).JSON(responsePayload.GetRespPayload())
}

func (controller *WorkingDayCalendarController) Detail(c *fiber.Ctx) error {
	var params entity.WorkingDayCalendarIDParams
	headerAcceptLang := controller.getAcceptLang(c)
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	if err := c.ParamsParser(&params); err != nil {
		log.Println("WorkingDayCalendarController, Detail, ParamsParser:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	errs := controller.validator.ValidateStruct(params, headerAcceptLang)
	if errs != nil {
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	custID := c.Locals("cust_id").(string)
	parentCustID := c.Locals("parent_cust_id").(string)
	data, err := controller.WorkingDayCalendarService.Detail(params.WorkingDayCalendarID, custID, parentCustID)
	if err != nil {
		return controller.handleReadError(c, responsePayload, err)
	}

	responsePayload.Setmsg(constant.SUCCESS_NO_DATA_DISPLAYED)
	responsePayload.Setdata(data)
	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}

func (controller *WorkingDayCalendarController) Calendar(c *fiber.Ctx) error {
	var params entity.WorkingDayCalendarIDParams
	var filter entity.WorkingDayCalendarViewFilter
	headerAcceptLang := controller.getAcceptLang(c)
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	if err := c.ParamsParser(&params); err != nil {
		log.Println("WorkingDayCalendarController, Calendar, ParamsParser:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}
	if err := c.QueryParser(&filter); err != nil {
		log.Println("WorkingDayCalendarController, Calendar, QueryParser:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	errs := controller.validator.ValidateStruct(params, headerAcceptLang)
	if errs != nil {
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	custID := c.Locals("cust_id").(string)
	parentCustID := c.Locals("parent_cust_id").(string)
	data, err := controller.WorkingDayCalendarService.Calendar(params.WorkingDayCalendarID, filter, custID, parentCustID)
	if err != nil {
		return controller.handleReadError(c, responsePayload, err)
	}

	responsePayload.Setmsg(constant.SUCCESS_NO_DATA_DISPLAYED)
	responsePayload.Setdata(data)
	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}

func (controller *WorkingDayCalendarController) DownloadHolidayTemplate(c *fiber.Ctx) error {
	var params entity.WorkingDayCalendarIDParams
	headerAcceptLang := controller.getAcceptLang(c)
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	if err := c.ParamsParser(&params); err != nil {
		log.Println("WorkingDayCalendarController, DownloadHolidayTemplate, ParamsParser:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}
	if errs := controller.validator.ValidateStruct(params, headerAcceptLang); errs != nil {
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	format := c.Query("format", "xlsx")
	switch format {
	case "csv", "xls", "xlsx":
	default:
		responsePayload.Setmsg("Unsupported format. Use csv, xls, or xlsx")
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	custID := c.Locals("cust_id").(string)
	parentCustID := c.Locals("parent_cust_id").(string)
	buffer, contentType, filename, err := controller.WorkingDayCalendarService.DownloadHolidayTemplate(params.WorkingDayCalendarID, format, custID, parentCustID)
	if err != nil {
		return controller.handleReadError(c, responsePayload, err)
	}

	c.Set("Content-Type", contentType)
	c.Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filename))
	return c.Send(buffer.Bytes())
}

func (controller *WorkingDayCalendarController) DownloadBlankHolidayTemplate(c *fiber.Ctx) error {
	headerAcceptLang := controller.getAcceptLang(c)
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	format := c.Query("format", "xlsx")
	switch format {
	case "csv", "xls", "xlsx":
	default:
		responsePayload.Setmsg("Unsupported format. Use csv, xls, or xlsx")
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	custID := c.Locals("cust_id").(string)
	parentCustID := c.Locals("parent_cust_id").(string)
	buffer, contentType, filename, err := controller.WorkingDayCalendarService.DownloadHolidayTemplate(0, format, custID, parentCustID)
	if err != nil {
		return controller.handleReadError(c, responsePayload, err)
	}

	c.Set("Content-Type", contentType)
	c.Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filename))
	return c.Send(buffer.Bytes())
}

func (controller *WorkingDayCalendarController) ImportHolidays(c *fiber.Ctx) error {
	var params entity.WorkingDayCalendarIDParams
	var request entity.WorkingDayCalendarImportHolidayRequest
	headerAcceptLang := controller.getAcceptLang(c)
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	if err := c.ParamsParser(&params); err != nil {
		log.Println("WorkingDayCalendarController, ImportHolidays, ParamsParser:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}
	if err := c.BodyParser(&request); err != nil {
		log.Println("WorkingDayCalendarController, ImportHolidays, BodyParser:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	if errs := controller.validator.ValidateStruct(params, headerAcceptLang); errs != nil {
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}
	if errs := controller.validator.ValidateStruct(request, headerAcceptLang); errs != nil {
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	custID := c.Locals("cust_id").(string)
	parentCustID := c.Locals("parent_cust_id").(string)
	userID := c.Locals("user_id").(int64)
	data, err := controller.WorkingDayCalendarService.ImportHolidays(params.WorkingDayCalendarID, request, custID, parentCustID, userID)
	if err != nil {
		responsePayload.Setmsg(err.Error())
		responsePayload.Setdata(data)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	responsePayload.Setmsg("File imported successfully")
	responsePayload.Setdata(data)
	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}

func (controller *WorkingDayCalendarController) handleReadError(c *fiber.Ctx, responsePayload *responsebuild.DataRespReq, err error) error {
	statusCode := fiber.StatusBadRequest
	errMsg := err.Error()
	if errors.Is(err, sql.ErrNoRows) || err.Error() == "sql: no rows in result set" {
		statusCode = fiber.StatusNotFound
		errMsg = constant.RECORD_NOT_FOUND
	}

	responsePayload.Setmsg(errMsg)
	return c.Status(statusCode).JSON(responsePayload.GetRespPayload())
}
