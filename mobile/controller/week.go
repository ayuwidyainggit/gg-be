package controller

import (
	"mobile/entity"
	"mobile/pkg/constant"
	"mobile/pkg/middleware"
	"mobile/pkg/responsebuild"
	"mobile/pkg/validation"
	"mobile/service"

	"github.com/gofiber/fiber/v2/log"

	"github.com/gofiber/fiber/v2"
)

type WeekController struct {
	WeekService service.WeekService
	validator   *validation.Validate
}

func NewWeekController(
	weekService service.WeekService,
	validator *validation.Validate,
) *WeekController {
	return &WeekController{
		WeekService: weekService,
		validator:   validator,
	}
}

func (controller *WeekController) Route(app *fiber.App) {
	weekMobileRouteV1 := app.Group("/v1/week", middleware.JWTProtected())
	weekMobileRouteV1.Get("", controller.List)
}

func (controller *WeekController) List(c *fiber.Ctx) error {
	var (
		headerAcceptLang string
		dataFilter       entity.WeekListQueryFilter
	)

	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	if err := c.QueryParser(&dataFilter); err != nil {
		log.Error("WeekController, List, query parser filter:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	if dataFilter.Page <= 0 {
		dataFilter.Page = 1
	}
	if dataFilter.Limit <= 0 {
		dataFilter.Limit = 5
	}
	if dataFilter.Sort == "" {
		dataFilter.Sort = "week_start:desc"
	}

	dataFilter.CustID = c.Locals("cust_id").(string)

	// Detect if user is distributor
	if distributorID, ok := c.Locals("distributor_id").(int64); ok && distributorID > 0 {
		dataFilter.IsDistributor = true
	}

	if empID, ok := c.Locals("emp_id").(int64); ok {
		dataFilter.EmpID = empID
	}

	errs := controller.validator.ValidateStruct(dataFilter, headerAcceptLang)
	if errs != nil {
		log.Error("WeekController, List, ValidateStruct, errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	data, total, lastPage, err := controller.WeekService.List(dataFilter)
	if err != nil {
		log.Error("WeekController, List, err:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusNotFound).JSON(responsePayload.GetRespPayload())
	}

	if len(data) == 0 {
		responsePayload.Setmsg(constant.STATUS_DB_NOT_FOUND)
		responsePayload.Setdata(nil)
		responsePayload.Setpaging(entity.Pagination{
			TotalRecord: 0,
			PageCurrent: dataFilter.Page,
			PageLimit:   dataFilter.Limit,
			PageTotal:   0,
		})
		return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
	}

	responsePayload.Setmsg(constant.STATUS_OK)
	responsePayload.Setdata(data)
	responsePayload.Setpaging(entity.Pagination{
		TotalRecord: total,
		PageCurrent: dataFilter.Page,
		PageLimit:   dataFilter.Limit,
		PageTotal:   lastPage,
	})

	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}
