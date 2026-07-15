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

type SupplierController struct {
	SupplierService service.SupplierService
	validator       *validation.Validate
}

func NewSupplierController(supplierService service.SupplierService, validator *validation.Validate) *SupplierController {
	return &SupplierController{
		SupplierService: supplierService,
		validator:       validator,
	}
}

func (controller *SupplierController) Route(app *fiber.App) {
	qParamId := ":sup_id"
	suppliersRouteV1 := app.Group("/v1/suppliers", middleware.JWTProtected())
	suppliersRouteV1.Get("/"+qParamId, controller.Detail)
	suppliersRouteV1.Get("", controller.List)
	suppliersRouteV1.Post("", controller.Create)
	suppliersRouteV1.Patch("/"+qParamId, controller.Update)
	suppliersRouteV1.Delete("/"+qParamId, controller.Delete)
}

func (controller *SupplierController) Detail(c *fiber.Ctx) error {
	var params entity.DetailSupplierParams
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	if err := c.ParamsParser(&params); err != nil {
		log.Println("SupplierController, Detail, ParamsParser:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	errs := controller.validator.ValidateStruct(params, headerAcceptLang)
	if errs != nil {
		log.Println("SupplierController, Detail, ValidateStruct(params), errs:", errs)
		log.Println("BankController, Detail, ValidateStruct(params), errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	custId := c.Locals("parent_cust_id").(string)
	log.Println("SupplierController, Detail, params:", structs.StructToJson(params))
	log.Println("SupplierController, Detail, CustId:", custId)

	data, err := controller.SupplierService.Detail(params.SupplierId, custId)
	if err != nil {
		log.Println("SupplierController, Detail, FindOneBySupplierId, err:", err.Error())
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

func (controller *SupplierController) List(c *fiber.Ctx) error {
	var (
		err        error
		dataFilter entity.SupplierQueryFilter
		data       interface{}
		total      int
		lastPage   int
		sup []entity.SupplierResponse
	)
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	if err := c.QueryParser(&dataFilter); err != nil {
		log.Println("SupplierController, List, query parser filter:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	custId := c.Locals("parent_cust_id").(string)
	log.Println("SupplierController, List, CustId:", custId)

	switch dataFilter.Mode {
	case "lookup":
		lookupScope := entity.SupplierLookupScope{
			CustID:          c.Locals("cust_id").(string),
			ParentCustID:    custId,
			DistributorIDs:  supplierLookupDistributorIDsFromQuery(c),
			IncludeParentID: queryTruthy(c, "include_parent_id") || queryTruthy(c, "is_include_parent_id"),
		}
		data, total, lastPage, err = controller.SupplierService.LookupList(dataFilter, lookupScope)
		if err != nil {
			log.Println("SupplierController, List, data, err:", err.Error())
			return c.Status(fiber.StatusBadRequest).JSON(entity.ApiResponse{
				RequestId: c.Locals("requestid").(string),
				Message:   err.Error(),
			})
		}
	default:
		data, total, lastPage, err = controller.SupplierService.List(dataFilter, custId)
		if err != nil {
			log.Println("SupplierController, List, data, err:", err.Error())
			return c.Status(fiber.StatusBadRequest).JSON(entity.ApiResponse{
				RequestId: c.Locals("requestid").(string),
				Message:   err.Error(),
			})
		}
		err = structs.Automapper(sup, &data)
		if err != nil {
			log.Println("SubBrand1Controller, List, data, err:", err.Error())
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

func (controller *SupplierController) Create(c *fiber.Ctx) error {
	var request entity.CreateSupplierBody
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	if err := c.BodyParser(&request); err != nil {
		log.Println("SupplierController, Create, BodyParser:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	custId := c.Locals("cust_id").(string)
	userId := c.Locals("user_id").(int64)
	distributorId := c.Locals("distributor_id").(int64)
	// log.Println("SupplierController, Create, CustId:", custId)

	request.CustId = custId
	request.CreatedBy = userId
	if distributorId != 0 {
		request.DistributorId = &distributorId
	}

	errs := controller.validator.ValidateStruct(request, headerAcceptLang)
	if errs != nil {
		log.Println("SupplierController, Create, ValidateStruct, errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	_, err := controller.SupplierService.Store(request)
	if err != nil {
		log.Println("SupplierController, Create, Store, err:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	responsePayload.Setmsg(constant.SUCCESSFULLY_ADDED)
	return c.Status(fiber.StatusCreated).JSON(responsePayload.GetRespPayload())
}

func (controller *SupplierController) Update(c *fiber.Ctx) error {
	var (
		params  entity.UpdateSupplierParams
		request entity.UpdateSupplierRequest
	)
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	if err := c.ParamsParser(&params); err != nil {
		log.Println("SupplierController, Update, ParamsParser(params):", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	errs := controller.validator.ValidateStruct(params, headerAcceptLang)
	if errs != nil {
		log.Println("SupplierController, Update, ValidateStruct(params), errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	if err := c.BodyParser(&request); err != nil {
		log.Println("SupplierController, Update, BodyParser(request), err:", err.Error())
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	custId := c.Locals("parent_cust_id").(string)
	userId := c.Locals("user_id").(int64)
	// log.Println("SupplierController, Update, CustId:", custId)

	request.CustId = custId
	request.UpdatedBy = userId

	errs = controller.validator.ValidateStruct(request, headerAcceptLang)
	if errs != nil {
		log.Println("SupplierController, Update, ValidateStruct(request), errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	err := controller.SupplierService.Update(params.SupplierId, request)
	if err != nil {
		log.Println("SupplierController, Update, Service.Update, err:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	responsePayload.Setmsg(constant.SUCCESSFULLY_UPDATED)
	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}

func (controller *SupplierController) Delete(c *fiber.Ctx) error {
	var params entity.DeleteSupplierParams
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	if err := c.ParamsParser(&params); err != nil {
		log.Println("SupplierController, Delete, ParamsParser, err:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	errs := controller.validator.ValidateStruct(params, headerAcceptLang)
	if errs != nil {
		log.Println("SupplierController, Delete, ValidateStruct, errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())

	}

	custId := c.Locals("parent_cust_id").(string)
	userId := c.Locals("user_id").(int64)
	// log.Println("SupplierController, Delete, CustId:", custId)

	err := controller.SupplierService.Delete(custId, params.SupplierId, userId)
	if err != nil {
		log.Println("SupplierController, Delete, Service.Delete, err:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	responsePayload.Setmsg("Deleted Successfully")
	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}

// supplierLookupDistributorIDsFromQuery mendukung distributor_id=1,2,3 dan distributor_id=1&distributor_id=2
func supplierLookupDistributorIDsFromQuery(c *fiber.Ctx) []int {
	var out []int
	c.Context().QueryArgs().VisitAll(func(k, v []byte) {
		if string(k) != "distributor_id" {
			return
		}
		s := strings.TrimSpace(string(v))
		if s == "" {
			return
		}
		if strings.Contains(s, ",") {
			for _, p := range strings.Split(s, ",") {
				p = strings.TrimSpace(p)
				if p == "" {
					continue
				}
				n, err := strconv.Atoi(p)
				if err != nil {
					continue
				}
				out = append(out, n)
			}
			return
		}
		n, err := strconv.Atoi(s)
		if err != nil {
			return
		}
		out = append(out, n)
	})
	return out
}

func queryTruthy(c *fiber.Ctx, key string) bool {
	v := strings.ToLower(strings.TrimSpace(c.Query(key)))
	return v == "true" || v == "1" || v == "yes"
}
