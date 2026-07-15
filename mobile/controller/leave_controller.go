package controller

import (
	"mobile/entity"
	"mobile/pkg/constant"
	"mobile/pkg/middleware"
	"mobile/pkg/responsebuild"
	"mobile/pkg/validation"
	"mobile/service"

	"github.com/gofiber/fiber/v2"
)

type LeaveController struct {
	leaveService service.LeaveService
	validator    *validation.Validate
}

func NewLeaveController(
	leaveService service.LeaveService,
	validator *validation.Validate,
) *LeaveController {
	return &LeaveController{
		leaveService: leaveService,
		validator:    validator,
	}
}

func (controller *LeaveController) Route(app *fiber.App) {
	leaveRouteV1 := app.Group("/v1/leave-request", middleware.JWTProtected())
	leaveRouteV1.Post("", controller.Create)
	leaveRouteV1.Get("", controller.List)
	leaveCheckV1 := app.Group("/v1/leave-check", middleware.JWTProtected())
	leaveCheckV1.Get("", controller.LeaveCheck)
}

func (controller *LeaveController) Create(c *fiber.Ctx) error {
	var (
		headerAcceptLang string
		request          entity.LeaveRequestCreate
	)

	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	request.StartDate = c.FormValue("start_date")
	request.EndDate = c.FormValue("end_date")
	request.Reason = c.FormValue("reason")

	if file, err := c.FormFile("file"); err == nil {
		request.File = file
	}

	if errs := controller.validator.ValidateStruct(request, headerAcceptLang); errs != nil {
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	request.CustID = c.Locals("cust_id").(string)
	request.EmpID = c.Locals("emp_id").(int64)
	request.EmpCode = c.Locals("emp_code").(string)
	request.UserID = c.Locals("user_id").(int64)

	if err := controller.leaveService.CreateLeaveRequest(request); err != nil {
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	responsePayload.Setmsg("Leave request successfully submitted")
	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}

func (controller *LeaveController) List(c *fiber.Ctx) error {
	var (
		headerAcceptLang string
		dataFilter       entity.LeaveRequestQuery
	)

	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	if err := c.QueryParser(&dataFilter); err != nil {
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	if errs := controller.validator.ValidateStruct(dataFilter, headerAcceptLang); errs != nil {
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	dataFilter.CustID = c.Locals("cust_id").(string)
	dataFilter.EmpID = c.Locals("emp_id").(int64)
	dataFilter.Limit, dataFilter.Page = controller.setPaginationDefaults(dataFilter.Limit, dataFilter.Page, 10)

	data, total, lastPage, err := controller.leaveService.ListLeaveRequests(dataFilter)
	if err != nil {
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	responsePayload.Setdata(data)
	responsePayload.Setpaging(entity.Pagination{
		TotalRecord: total,
		PageCurrent: dataFilter.Page,
		PageLimit:   dataFilter.Limit,
		PageTotal:   lastPage,
	})
	responsePayload.Setmsg(constant.STATUS_OK)

	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}

func (controller *LeaveController) setPaginationDefaults(limit, page int, defaultLimit int) (int, int) {
	if limit <= 0 || limit > 9999 {
		limit = defaultLimit
	}
	if page <= 0 {
		page = 1
	}
	return limit, page
}

func (controller *LeaveController) LeaveCheck(c *fiber.Ctx) error {
	var (
		headerAcceptLang string
	)

	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}

	custID := c.Locals("cust_id").(string)
	empID := c.Locals("emp_id").(int64)

	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	data, err := controller.leaveService.LeaveCheck(custID, empID)
	if err != nil {
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(responsePayload.GetRespPayload())
	}
	statusData := "Check-in Available"
	if data != nil {
		statusData = "Check-in Not Available"
	}

	responsePayload.Setdata(data)
	responsePayload.Setmsg(statusData)

	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}
