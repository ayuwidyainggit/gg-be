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

type AttendanceController struct {
	attendanceService service.AttendanceService
	validator         *validation.Validate
}

func NewAttendanceController(
	attendanceService service.AttendanceService,
	validator *validation.Validate,
) *AttendanceController {
	return &AttendanceController{
		attendanceService: attendanceService,
		validator:         validator,
	}
}

func (controller *AttendanceController) Route(app *fiber.App) {
	RouteV1 := app.Group("/v1/attendances", middleware.JWTProtected())
	RouteV1.Post("", controller.AttendanceRequest)
	RouteV1.Get("", controller.AttendanceGet)
	RouteV1.Get("/check", controller.AttendanceCheck)
}
func (controller *AttendanceController) AttendanceRequest(c *fiber.Ctx) error {
	var (
		headerAcceptLang string
		request          entity.AttendanceRequest
	)

	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	err := c.BodyParser(&request)
	if err != nil {
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())

	}
	errs := controller.validator.ValidateStruct(request, headerAcceptLang)
	if errs != nil {
		// log.Println("errs, UserController-Login:", structs.StructToJson(errs))
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	request.Email = c.Locals("email").(string)
	request.CustID = c.Locals("cust_id").(string)

	data, err := controller.attendanceService.AttendanceRequest(request)
	if err != nil {
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}
	responsePayload.Setdata(data)
	responsePayload.Setmsg(constant.STATUS_OK)
	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}
func (controller *AttendanceController) AttendanceGet(c *fiber.Ctx) error {
	var (
		headerAcceptLang string
		request          entity.AttendanceGetRequest
	)

	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	request.Email = c.Locals("email").(string)
	request.CustID = c.Locals("cust_id").(string)

	data, err := controller.attendanceService.AttendanceGet(request)
	if err != nil {
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}
	responsePayload.Setdata(data)
	responsePayload.Setmsg(constant.STATUS_OK)
	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}

func (controller *AttendanceController) AttendanceCheck(c *fiber.Ctx) error {
	var (
		headerAcceptLang string
		request          entity.AttendanceCheckRequest
	)

	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}

	request.CustID = c.Locals("cust_id").(string)

	err := c.QueryParser(&request)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success":     false,
			"message":     err.Error(),
			"description": "",
			"data":        nil,
			"request_id":  c.Locals("requestid").(string),
		})
	}

	errs := controller.validator.ValidateStruct(request, headerAcceptLang)
	if errs != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success":     false,
			"message":     fiber.ErrBadRequest.Message,
			"description": "",
			"errors":      errs,
			"request_id":  c.Locals("requestid").(string),
		})
	}

	data, err := controller.attendanceService.AttendanceCheck(request)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success":     false,
			"message":     err.Error(),
			"description": "",
			"data":        nil,
			"request_id":  c.Locals("requestid").(string),
		})
	}

	responseData := fiber.Map{
		"plan": data.Data.Plan,
	}

	if data.Data.EmpID != nil {
		responseData["emp_id"] = *data.Data.EmpID

		if data.Data.EmpCode != nil {
			responseData["emp_code"] = *data.Data.EmpCode
		}
		if data.Data.EmpName != nil {
			responseData["emp_name"] = *data.Data.EmpName
		}
		if data.Data.OprType != nil {
			responseData["opr_type"] = *data.Data.OprType
		} else {
			responseData["opr_type"] = ""
		}

		if data.Data.OprTypeCanvas != nil {
			responseData["opr_type_canvas"] = *data.Data.OprTypeCanvas
		} else {
			responseData["opr_type_canvas"] = ""
		}

		if data.Data.WhID != nil {
			responseData["wh_id"] = *data.Data.WhID
		} else {
			responseData["wh_id"] = nil
		}

		if data.Data.WhCode != nil {
			responseData["wh_code"] = *data.Data.WhCode
		} else {
			responseData["wh_code"] = ""
		}

		if data.Data.WhNameCanvas != nil {
			responseData["wh_name_canvas"] = *data.Data.WhNameCanvas
		} else {
			responseData["wh_name_canvas"] = ""
		}

		if data.Data.Stock != nil {
			responseData["stock"] = *data.Data.Stock
		} else {
			responseData["stock"] = 0
		}
	}

	response := fiber.Map{
		"success":     data.Success,
		"message":     data.Message,
		"description": data.Description,
		"data":        responseData,
		"request_id":  c.Locals("requestid").(string),
	}

	return c.Status(fiber.StatusOK).JSON(response)
}
