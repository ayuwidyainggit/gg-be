package controller

import (
	"mobile/entity"
	"mobile/pkg/constant"
	"mobile/pkg/middleware"
	"mobile/pkg/responsebuild"
	"mobile/pkg/structs"
	"mobile/pkg/validation"
	"mobile/service"

	"github.com/gofiber/fiber/v2/log"

	"github.com/gofiber/fiber/v2"
)

type EmployeeController struct {
	EmployeeService service.EmployeeService
	validator       *validation.Validate
}

func NewEmployeeController(empGroupService service.EmployeeService, validator *validation.Validate) EmployeeController {
	return EmployeeController{
		EmployeeService: empGroupService,
		validator:       validator,
	}
}

func (controller *EmployeeController) Route(app *fiber.App) {
	qParamId := ":emp_id"
	// employeeRouteV2 := app.Group("/create-multiple", middleware.JWTProtected())
	// employeeRouteV2.Post("", controller.CreateMultiple)

	employeeRouteV1 := app.Group("/v1/employees", middleware.JWTProtected())
	employeeRouteV1.Get("/"+qParamId, controller.Detail)
	employeeRouteV1.Get("", controller.List)
	employeeRouteV1.Post("", controller.Create)
	employeeRouteV1.Patch("/"+qParamId, controller.Update)
	employeeRouteV1.Delete("/"+qParamId, controller.Delete)
	employeeRouteV1.Post("/create-multiple", controller.CreateMultiple)

}

func (controller *EmployeeController) List(c *fiber.Ctx) error {
	var (
		err             error
		dataFilter      entity.EmployeeQueryFilter
		data            interface{}
		total           int64
		lastPage        int
		employees       []entity.EmployeeResponse
		employeesLookup []entity.EmployeeLookupResponse
	)
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), constant.HEADER_ACCEPT_LANG)
	if err := c.QueryParser(&dataFilter); err != nil {
		log.Error("EmployeeController, List, query parser filter:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	dataFilter.CustId = c.Locals("cust_id").(string)
	dataFilter.ParentCustId = c.Locals("parent_cust_id").(string)

	switch dataFilter.Mode {
	case "lookup":
		data, total, lastPage, err = controller.EmployeeService.LookupList(dataFilter)
		if err != nil {
			log.Error("EmployeeController, Lookup List, data, err:", err.Error())
			responsePayload.Setmsg(err.Error())
			return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
		}

		err = structs.Automapper(employeesLookup, &data)
		if err != nil {
			log.Error("EmployeeController, Automapper Lookup List, data, err:", err.Error())
			responsePayload.Setmsg(err.Error())
			return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
		}
	default:
		data, total, lastPage, err = controller.EmployeeService.List(dataFilter)
		if err != nil {
			log.Error("EmployeeController, List, data, err:", err.Error())
			responsePayload.Setmsg(err.Error())
			return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
		}

		err = structs.Automapper(employees, &data)
		if err != nil {
			log.Error("EmployeeController, Automapper List, data, err:", err.Error())
			responsePayload.Setmsg(err.Error())
			return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
		}
	}

	responsePayload.Setdata(data)
	responsePayload.Setpaging(entity.Pagination{
		TotalRecord: total,
		PageCurrent: dataFilter.Page,
		PageLimit:   dataFilter.Limit,
		PageTotal:   lastPage,
	})
	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}

func (controller *EmployeeController) Detail(c *fiber.Ctx) error {
	var params entity.DetailEmployeeParams
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), constant.HEADER_ACCEPT_LANG)

	if err := c.ParamsParser(&params); err != nil {
		log.Error("EmployeeController, Detail, ParamsParser:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	errs := controller.validator.ValidateStruct(params, constant.HEADER_ACCEPT_LANG)
	if errs != nil {
		log.Error("EmployeeController, Detail, ValidateStruct(params), errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	params.CustId = c.Locals("cust_id").(string)
	params.ParentCustId = c.Locals("parent_cust_id").(string)
	// log.Println("EmployeeController, Detail, CustId:", custId)

	data, err := controller.EmployeeService.Detail(params)
	if err != nil {
		log.Error("EmployeeController, Detail, FindOneByEmployeeId, err:", err.Error())
		statusCode := fiber.StatusBadRequest
		errMsg := err.Error()
		if err.Error() == "sql: no rows in result set" {
			statusCode = fiber.StatusNotFound
			errMsg = "Not found"
		}
		responsePayload.Setmsg(errMsg)
		return c.Status(statusCode).JSON(responsePayload.GetRespPayload())
	}

	responsePayload.Setdata(data)
	return c.JSON(responsePayload.GetRespPayload())
}

func (controller *EmployeeController) Create(c *fiber.Ctx) error {
	var request entity.CreateEmployeeBody
	if err := c.BodyParser(&request); err != nil {
		log.Error("EmployeeController, Create, BodyParser:", err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(entity.ApiResponse{
			RequestId: c.Locals("requestid").(string),
			Message:   err.Error(),
		})
	}

	request.CustId = c.Locals("cust_id").(string)
	request.ParentCustId = c.Locals("parent_cust_id").(string)
	request.CreatedBy = c.Locals("user_id").(int64)
	errs := controller.validator.ValidateStruct(request, constant.HEADER_ACCEPT_LANG)
	if errs != nil {
		log.Error("EmployeeController, Create, ValidateStruct, errs:", errs)
		return c.Status(fiber.StatusBadRequest).JSON(entity.ApiResponse{
			RequestId: c.Locals("requestid").(string),
			Message:   fiber.ErrBadRequest.Message,
			Errors:    errs,
		})
	}

	_, err := controller.EmployeeService.Store(request)
	if err != nil {
		log.Error("EmployeeController, Create, Store, err:", err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(entity.ApiResponse{
			RequestId: c.Locals("requestid").(string),
			Message:   err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(entity.ApiResponse{
		RequestId: c.Locals("requestid").(string),
		Message:   constant.SUCCESSFULLY_ADDED,
	})
}

func (controller *EmployeeController) Update(c *fiber.Ctx) error {
	var (
		params  entity.UpdateEmployeeParams
		request entity.UpdateEmployeeRequest
	)

	if err := c.ParamsParser(&params); err != nil {
		log.Error("EmployeeController, Update, ParamsParser(params):", err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(entity.ApiResponse{
			RequestId: c.Locals("requestid").(string),
			Message:   err.Error(),
		})
	}

	request.CustId = c.Locals("cust_id").(string)
	request.ParentCustId = c.Locals("parent_cust_id").(string)
	request.UpdatedBy = c.Locals("user_id").(int64)

	errs := controller.validator.ValidateStruct(params, constant.HEADER_ACCEPT_LANG)
	if errs != nil {
		log.Error("EmployeeController, Update, ValidateStruct(params), errs:", errs)
		return c.Status(fiber.StatusBadRequest).JSON(entity.ApiResponse{
			RequestId: c.Locals("requestid").(string),
			Message:   fiber.ErrBadRequest.Message,
			Errors:    errs,
		})
	}

	if err := c.BodyParser(&request); err != nil {
		log.Error("EmployeeController, Update, BodyParser(request), err:", err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(entity.ApiResponse{
			RequestId: c.Locals("requestid").(string),
			Message:   err.Error(),
		})
	}

	errs = controller.validator.ValidateStruct(request, constant.HEADER_ACCEPT_LANG)
	if errs != nil {
		// log.Println("EmployeeController, Update, ValidateStruct(request), errs:", errs)
		return c.Status(fiber.StatusBadRequest).JSON(entity.ApiResponse{
			RequestId: c.Locals("requestid").(string),
			Message:   fiber.ErrBadRequest.Message,
			Errors:    errs,
		})
	}

	err := controller.EmployeeService.Update(params.EmployeeId, request)
	if err != nil {
		log.Error("EmployeeController, Update, Service.Update, err:", err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(entity.ApiResponse{
			RequestId: c.Locals("requestid").(string),
			Message:   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(entity.ApiResponse{
		RequestId: c.Locals("requestid").(string),
		Message:   constant.SUCCESSFULLY_UPDATED,
	})
}

func (controller *EmployeeController) Delete(c *fiber.Ctx) error {
	var params entity.DeleteEmployeeParams
	if err := c.ParamsParser(&params); err != nil {
		log.Error("EmployeeController, Delete, ParamsParser, err:", err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(entity.ApiResponse{
			RequestId: c.Locals("requestid").(string),
			Message:   err.Error(),
		})
	}

	errs := controller.validator.ValidateStruct(params, constant.HEADER_ACCEPT_LANG)
	if errs != nil {
		log.Error("EmployeeController, Delete, ValidateStruct, errs:", errs)
		return c.Status(fiber.StatusBadRequest).JSON(entity.ApiResponse{
			RequestId: c.Locals("requestid").(string),
			Message:   fiber.ErrBadRequest.Message,
			Errors:    errs,
		})
	}

	custId := c.Locals("cust_id").(string)
	userId := c.Locals("user_id").(int64)
	// log.Println("EmployeeController, Delete, CustId:", custId)

	err := controller.EmployeeService.Delete(custId, params.EmployeeId, userId)
	if err != nil {
		log.Error("EmployeeController, Delete, Service.Delete, err:", err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(entity.ApiResponse{
			RequestId: c.Locals("requestid").(string),
			Message:   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(entity.ApiResponse{
		RequestId: c.Locals("requestid").(string),
		Message:   "Deleted Successfully",
	})
}

func (controller *EmployeeController) CreateMultiple(c *fiber.Ctx) error {
	var request entity.CreateMultipleEmployeeBody
	if err := c.BodyParser(&request); err != nil {
		log.Error("EmployeeController, Create, BodyParser:", err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(entity.ApiResponse{
			RequestId: c.Locals("requestid").(string),
			Message:   err.Error(),
		})
	}

	custId := c.Locals("cust_id").(string)
	userId := c.Locals("user_id").(int64)
	// log.Println("EmployeeController, Create, CustId:", custId)

	// request.CustId = custId
	// request.CreatedBy = userId
	errs := controller.validator.ValidateStruct(request, constant.HEADER_ACCEPT_LANG)
	if errs != nil {
		log.Error("EmployeeController, Create, ValidateStruct, errs:", errs)
		return c.Status(fiber.StatusBadRequest).JSON(entity.ApiResponse{
			RequestId: c.Locals("requestid").(string),
			Message:   fiber.ErrBadRequest.Message,
			Errors:    errs,
		})
	}

	_, err := controller.EmployeeService.StoreMultiple(request, custId, userId)
	if err != nil {
		log.Error("EmployeeController, Create, Store, err:", err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(entity.ApiResponse{
			RequestId: c.Locals("requestid").(string),
			Message:   err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(entity.ApiResponse{
		RequestId: c.Locals("requestid").(string),
		Message:   "Created Successfully",
	})
}
