package controller

import (
	"log"
	"master/entity"
	"master/pkg/constant"
	"master/pkg/middleware"
	"master/pkg/responsebuild"
	"master/pkg/structs"
	"master/pkg/validation"
	"master/service"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type WarehouseController struct {
	WarehouseService service.WarehouseService
	validator        *validation.Validate
}

func NewWarehouseController(warehouseService service.WarehouseService, validator *validation.Validate) *WarehouseController {
	return &WarehouseController{
		WarehouseService: warehouseService,
		validator:        validator,
	}
}

func (controller *WarehouseController) Route(app *fiber.App) {
	qParamId := ":wh_id"
	outletTypesRouteV1 := app.Group("/v1/warehouses", middleware.JWTProtected())
	outletTypesRouteV1.Get("/"+qParamId, controller.Detail)
	outletTypesRouteV1.Get("", controller.List)
	outletTypesRouteV1.Post("", controller.Create)
	outletTypesRouteV1.Patch("/"+qParamId, controller.Update)
	outletTypesRouteV1.Delete("/"+qParamId, controller.Delete)
}

func (controller *WarehouseController) Detail(c *fiber.Ctx) error {
	var params entity.DetailWarehouseParams
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	if err := c.ParamsParser(&params); err != nil {
		log.Println("WarehouseController, Detail, ParamsParser:", err.Error())
		responsePayload.Setmsg(err.Error())

		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())

	}

	errs := controller.validator.ValidateStruct(params, headerAcceptLang)
	if errs != nil {
		log.Println("WarehouseController, Detail, ValidateStruct(params), errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	custId := c.Locals("cust_id").(string)
	// log.Println("WarehouseController, Detail, CustId:", custId)

	data, err := controller.WarehouseService.Detail(params.WarehouseId, custId)
	if err != nil {
		log.Println("WarehouseController, Detail, FindOneByWarehouseId, err:", err.Error())
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

func (controller *WarehouseController) List(c *fiber.Ctx) error {
	var (
		err              error
		dataFilter       entity.WarehouseQueryFilter
		data             interface{}
		total            int
		lastPage         int
		warehouses       []entity.WarehouseResponse
		warehousesLookup []entity.WarehouseLookupResponse
	)
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	if err := c.QueryParser(&dataFilter); err != nil {
		log.Println("WarehouseController, List, query parser filter:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())

	}

	if dataFilter.Page < 1 {
		dataFilter.Page = 1
	}
	if dataFilter.Limit <= 0 {
		dataFilter.Limit = 10
	}
	dataFilter.DistributorIDs = warehouseQueryIntSliceFromRequest(c, "distributor_id")

	headerCustId := ""
	if len(c.GetReqHeaders()[constant.CUST_ID]) > 0 {
		headerCustId = c.GetReqHeaders()[constant.CUST_ID][0]
	}
	custId := headerCustId
	if custId == "" {
		custId = c.Locals("cust_id").(string)
	}
	// custId := c.Locals("cust_id").(string)

	switch dataFilter.Mode {
	case "lookup":
		data, total, lastPage, err = controller.WarehouseService.LookupList(dataFilter, custId)
		if err != nil {
			log.Println("WarehouseController, List, data, err:", err.Error())
			responsePayload.Setmsg(err.Error())
			return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
		}

		err = structs.Automapper(warehousesLookup, &data)
		if err != nil {
			log.Println("WarehouseController, List, data, err:", err.Error())
			responsePayload.Setmsg(err.Error())
			return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
		}
	default:
		data, total, lastPage, err = controller.WarehouseService.List(dataFilter, custId)
		if err != nil {
			log.Println("WarehouseController, List, data, err:", err.Error())
			return c.Status(fiber.StatusBadRequest).JSON(entity.ApiResponse{
				RequestId: c.Locals("requestid").(string),
				Message:   err.Error(),
			})
		}

		err = structs.Automapper(warehouses, &data)
		if err != nil {
			log.Println("WarehouseController, List, data, err:", err.Error())
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
		RequestID:   c.Locals("requestid").(string),
	})
	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}

func warehouseQueryIntSliceFromRequest(c *fiber.Ctx, key string) []int {
	var out []int
	c.Context().QueryArgs().VisitAll(func(k, v []byte) {
		if string(k) != key {
			return
		}
		s := strings.TrimSpace(string(v))
		if s == "" {
			return
		}
		if strings.Contains(s, ",") {
			out = append(out, warehouseParseCommaSeparatedInts(s)...)
			return
		}
		n, err := strconv.Atoi(s)
		if err != nil {
			return
		}
		if n <= 0 {
			return
		}
		out = append(out, n)
	})
	return out
}

func warehouseParseCommaSeparatedInts(raw string) []int {
	if raw == "" {
		return nil
	}
	parts := strings.Split(raw, ",")
	out := make([]int, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		v, err := strconv.Atoi(p)
		if err != nil {
			continue
		}
		if v <= 0 {
			continue
		}
		out = append(out, v)
	}
	if len(out) == 0 {
		return nil
	}
	return out
}

func (controller *WarehouseController) Create(c *fiber.Ctx) error {
	var request entity.CreateWarehouseBody
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	if err := c.BodyParser(&request); err != nil {
		log.Println("WarehouseController, Create, BodyParser:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	custId := c.Locals("cust_id").(string)
	userId := c.Locals("user_id").(int64)
	// log.Println("WarehouseController, Create, CustId:", custId)

	request.CustId = custId
	request.CreatedBy = userId

	errs := controller.validator.ValidateStruct(request, headerAcceptLang)
	if errs != nil {
		log.Println("WarehouseController, Create, ValidateStruct, errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())

	}

	_, err := controller.WarehouseService.Store(request)
	if err != nil {
		log.Println("WarehouseController, Create, Store, err:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	responsePayload.Setmsg(constant.SUCCESSFULLY_ADDED)
	return c.Status(fiber.StatusCreated).JSON(responsePayload.GetRespPayload())
}

func (controller *WarehouseController) Update(c *fiber.Ctx) error {
	var (
		params  entity.UpdateWarehouseParams
		request entity.UpdateWarehouseRequest
	)
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	if err := c.ParamsParser(&params); err != nil {
		log.Println("WarehouseController, Update, ParamsParser(params):", err.Error())
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(err)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	errs := controller.validator.ValidateStruct(params, headerAcceptLang)
	if errs != nil {
		log.Println("WarehouseController, Update, ValidateStruct(params), errs:", errs)
		return c.Status(fiber.StatusBadRequest).JSON(entity.ApiResponse{
			RequestId: c.Locals("requestid").(string),
			Message:   fiber.ErrBadRequest.Message,
			Errors:    errs,
		})
	}

	if err := c.BodyParser(&request); err != nil {
		log.Println("WarehouseController, Update, BodyParser(request), err:", err.Error())
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	custId := c.Locals("cust_id").(string)
	userId := c.Locals("user_id").(int64)
	// log.Println("WarehouseController, Update, CustId:", custId)

	request.CustId = custId
	request.UpdatedBy = userId

	errs = controller.validator.ValidateStruct(request, headerAcceptLang)
	if errs != nil {
		log.Println("WarehouseController, Update, ValidateStruct(request), errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	err := controller.WarehouseService.Update(params.WarehouseId, request)
	if err != nil {
		log.Println("WarehouseController, Update, Service.Update, err:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	responsePayload.Setmsg(constant.SUCCESSFULLY_UPDATED)
	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}

func (controller *WarehouseController) Delete(c *fiber.Ctx) error {
	var params entity.DeleteWarehouseParams
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	if err := c.ParamsParser(&params); err != nil {
		log.Println("WarehouseController, Delete, ParamsParser, err:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	errs := controller.validator.ValidateStruct(params, headerAcceptLang)
	if errs != nil {
		log.Println("WarehouseController, Delete, ValidateStruct, errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	custId := c.Locals("cust_id").(string)
	userId := c.Locals("user_id").(int64)
	// log.Println("WarehouseController, Delete, CustId:", custId)

	err := controller.WarehouseService.Delete(custId, params.WarehouseId, userId)
	if err != nil {
		log.Println("WarehouseController, Delete, Service.Delete, err:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}
	responsePayload.Setmsg("Deleted Successfully")
	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}
