package controller

import (
	"mobile/entity"
	"mobile/pkg/constant"
	"mobile/pkg/middleware"
	"mobile/pkg/responsebuild"
	"mobile/pkg/validation"
	"mobile/service"
	"strconv"

	"github.com/gofiber/fiber/v2/log"

	"github.com/gofiber/fiber/v2"
)

type VisitsController struct {
	VisitsService service.VisitsService
	validator     *validation.Validate
}

func NewVisitsController(
	VisitsService service.VisitsService,
	validator *validation.Validate,
) *VisitsController {
	return &VisitsController{
		VisitsService: VisitsService,
		validator:     validator,
	}
}

func (controller *VisitsController) Route(app *fiber.App) {
	VisitsRouteV1 := app.Group("/v1/visits", middleware.JWTProtected())
	VisitsRouteV1.Get("/", controller.Visits)
	VisitsRouteV1.Get("/summary", controller.Summaries)
	VisitsRouteV1.Get("/list", controller.List)
	VisitsRouteV1.Post("/start", controller.Start)
	VisitsRouteV1.Post("/skip", controller.Skip)
	VisitsRouteV1.Get("/skip/reasons", controller.SkipReasons)
	VisitsRouteV1.Post("/Arrive", controller.Arrive)
	VisitsRouteV1.Post("/Hold", controller.Hold)
	VisitsRouteV1.Post("/Resume", controller.Resume)
	VisitsRouteV1.Post("/Leave", controller.Leave)
	VisitsRouteV1.Post("/End", controller.End)
}

func (controller *VisitsController) Visits(c *fiber.Ctx) error {
	var (
		headerAcceptLang string
		request          entity.VisitsListRequest
	)

	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	// Extract cust_id from JWT context
	request.CustID = c.Locals("cust_id").(string)

	empCode, _ := c.Locals("emp_code").(string)
	empID, _ := c.Locals("emp_id").(int64)
	isDistributor, _ := c.Locals("is_distributor").(bool)

	request.EmpID = empID
	request.EmpCode = &empCode
	request.IsDistributor = isDistributor

	// Parse optional outlet_id query parameter
	if request.OutletID == nil {
		if outletIDStr := c.Query("outlet_id"); outletIDStr != "" {
			outletID, err := strconv.ParseInt(outletIDStr, 10, 64)
			if err != nil {
				responsePayload.Setmsg("invalid outlet_id format")
				return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
			}
			request.OutletID = &outletID
		}
	}

	// No validation needed for GET request - only cust_id from JWT is required
	data, err := controller.VisitsService.Visits(request)
	if err != nil {
		log.Error("VisitsController, Visits, err:", err.Error())
		statusCode := fiber.StatusBadRequest
		errMsg := err.Error()
		if err.Error() == constant.STATUS_DB_NOT_FOUND {
			statusCode = fiber.StatusNotFound
			errMsg = "Not found"
		}
		responsePayload.Setmsg(errMsg)
		return c.Status(statusCode).JSON(responsePayload.GetRespPayload())
	}
	responsePayload.Setdata(data)
	responsePayload.Setmsg(constant.STATUS_OK)
	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}

func (controller *VisitsController) Summaries(c *fiber.Ctx) error {
	var (
		headerAcceptLang string
		filter           entity.SummariesRequest
		responsePayload  = responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)
	)

	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}

	if err := c.QueryParser(&filter); err != nil {
		return err
	}

	empID, _ := c.Locals("emp_id").(int64)
	custID, _ := c.Locals("cust_id").(string)

	filter.EmpID = empID
	filter.CustID = custID

	response, err := controller.VisitsService.Summaries(c.UserContext(), filter)
	if err != nil {
		return err
	}

	responsePayload.Setdata(response)
	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}

func (controller *VisitsController) List(c *fiber.Ctx) error {
	var (
		headerAcceptLang string
		dataFilter       entity.VisitQueryFilter
	)

	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	if err := c.QueryParser(&dataFilter); err != nil {
		log.Error("VisitController, List, query parser filter:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}
	errs := controller.validator.ValidateStruct(dataFilter, headerAcceptLang)
	if errs != nil {
		log.Error("VisitController, List, ValidateStruct(params), errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	dataFilter.CustId = c.Locals("cust_id").(string)
	dataFilter.ParentCustId = c.Locals("parent_cust_id").(string)

	data, err := controller.VisitsService.List(dataFilter)
	if err != nil {
		log.Error("VisitController, List, data, err:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	responsePayload.Setmsg(constant.STATUS_OK)
	responsePayload.Setdata(data)
	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}

func (controller *VisitsController) Start(c *fiber.Ctx) error {
	var (
		headerAcceptLang string
		request          entity.StartRequest
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
	if err = controller.VisitsService.Start(request); err != nil {
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}
	// log.Println("response:", response)
	responsePayload.Setmsg(constant.STATUS_OK)
	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}

func (controller *VisitsController) Skip(c *fiber.Ctx) error {
	var (
		headerAcceptLang string
		request          entity.SkipRequest
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
	if err = controller.VisitsService.Skip(request); err != nil {
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}
	// log.Println("response:", response)
	responsePayload.Setmsg(constant.STATUS_OK)
	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}

func (controller *VisitsController) SkipReasons(c *fiber.Ctx) error {
	var (
		headerAcceptLang string
		dataFilter       entity.SkipReasonsQueryFilter
	)

	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	if err := c.QueryParser(&dataFilter); err != nil {
		log.Error("VisitController, List, query parser filter:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}
	errs := controller.validator.ValidateStruct(dataFilter, headerAcceptLang)
	if errs != nil {
		log.Error("VisitController, List, ValidateStruct(params), errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	dataFilter.CustId = c.Locals("cust_id").(string)
	dataFilter.ParentCustId = c.Locals("parent_cust_id").(string)

	data, err := controller.VisitsService.SkipReasons(dataFilter)
	if err != nil {
		log.Error("VisitController, List, data, err:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	responsePayload.Setmsg(constant.STATUS_OK)
	responsePayload.Setdata(data)
	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}

func (controller *VisitsController) Arrive(c *fiber.Ctx) error {
	var (
		headerAcceptLang string
		request          entity.ArriveRequest
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

	file, err := c.FormFile("file")
	if err != nil {
		responsePayload.Setmsg("file is required")
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}
	request.File = file

	currentTimeValue := c.FormValue("current_time")
	if currentTimeValue == "" {
		responsePayload.Setmsg("current_time is required")
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}
	currentTimeParsed, err := strconv.ParseInt(currentTimeValue, 10, 64)
	if err != nil {
		responsePayload.Setmsg("invalid current_time value")
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}
	request.CurrentTime = currentTimeParsed

	if outletCode := c.FormValue("outlet_code"); outletCode != "" {
		request.OutletCode = outletCode
	}
	if latitude := c.FormValue("latitude"); latitude != "" {
		request.Latitude = latitude
	}
	if longitude := c.FormValue("longitude"); longitude != "" {
		request.Longitude = longitude
	}
	if folder := c.FormValue("folder"); folder != "" {
		request.Folder = folder
	}

	isUpdateLocationRaw := c.FormValue("is_update_location")
	if isUpdateLocationRaw == "" {
		responsePayload.Setmsg("is_update_location is required")
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}
	isUpdateLocation, err := strconv.ParseBool(isUpdateLocationRaw)
	if err != nil {
		responsePayload.Setmsg("invalid is_update_location value")
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}
	request.IsUpdateLocation = isUpdateLocation

	errs := controller.validator.ValidateStruct(request, headerAcceptLang)
	if errs != nil {
		// log.Println("errs, UserController-Login:", structs.StructToJson(errs))
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	request.Email = c.Locals("email").(string)
	request.CustID = c.Locals("cust_id").(string)
	responseData, err := controller.VisitsService.Arrive(request)
	if err != nil {
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}
	// log.Println("response:", response)
	responsePayload.Setmsg(constant.STATUS_OK)
	responsePayload.Setdata(responseData)
	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}

func (controller *VisitsController) Hold(c *fiber.Ctx) error {
	var (
		headerAcceptLang string
		request          entity.HoldRequest
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
	if err = controller.VisitsService.Hold(request); err != nil {
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}
	// log.Println("response:", response)
	responsePayload.Setmsg(constant.STATUS_OK)
	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}

func (controller *VisitsController) Resume(c *fiber.Ctx) error {
	var (
		headerAcceptLang string
		request          entity.ResumeRequest
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
	if err = controller.VisitsService.Resume(request); err != nil {
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}
	// log.Println("response:", response)
	responsePayload.Setmsg(constant.STATUS_OK)
	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}

func (controller *VisitsController) Leave(c *fiber.Ctx) error {
	var (
		headerAcceptLang string
		request          entity.LeaveRequest
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
	if err = controller.VisitsService.Leave(request); err != nil {
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}
	// log.Println("response:", response)
	responsePayload.Setmsg(constant.STATUS_OK)
	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}

func (controller *VisitsController) End(c *fiber.Ctx) error {
	var (
		headerAcceptLang string
		request          entity.EndRequest
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
	if err = controller.VisitsService.End(request); err != nil {
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}
	// log.Println("response:", response)
	responsePayload.Setmsg(constant.STATUS_OK)
	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}
