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

	"github.com/gofiber/fiber/v2"
)

type ProductDistController struct {
	ProductDistService service.ProductDistService
	validator          *validation.Validate
}

func NewProductDistController(ProductDistService service.ProductDistService, validator *validation.Validate) *ProductDistController {
	return &ProductDistController{
		ProductDistService: ProductDistService,
		validator:          validator,
	}
}

func (controller *ProductDistController) Route(app *fiber.App) {
	qParamId := ":pro_id"
	productsDistRouteV1 := app.Group("/v1/distributor-products", middleware.JWTProtected())
	productsDistRouteV1.Get("/"+qParamId, controller.Detail)
	productsDistRouteV1.Get("", controller.List)
	productsDistRouteV1.Post("", controller.Create)
	productsDistRouteV1.Patch("/"+qParamId, controller.Update)
	productsDistRouteV1.Delete("/"+qParamId, controller.Delete)
}

func (controller *ProductDistController) Detail(c *fiber.Ctx) error {
	var params entity.DetailProductDistParams
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	if err := c.ParamsParser(&params); err != nil {
		log.Println("ProductDistController, Detail, ParamsParser:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	errs := controller.validator.ValidateStruct(params, headerAcceptLang)
	if errs != nil {
		log.Println("ProductDistController, Detail, ValidateStruct(params), errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	custId := c.Locals("cust_id").(string)
	parentCustId := c.Locals("parent_cust_id").(string)

	data, err := controller.ProductDistService.Detail(params.ProductId, custId, parentCustId)
	if err != nil {
		log.Println("ProductDistController, Detail, FindOneByProductId, err:", err.Error())
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

func (controller *ProductDistController) List(c *fiber.Ctx) error {
	var (
		err           error
		dataFilter    entity.ProductDistQueryFilter
		data          interface{}
		total         int
		lastPage      int
		products      []entity.ProductDistResponse
		productLookup []entity.ProductDistLookupResponse
		productSearch []entity.ProductDistSearchResponse
	)
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	if err := c.QueryParser(&dataFilter); err != nil {
		log.Println("ProductDistController, List, query parser filter:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	custId := c.Locals("cust_id").(string)
	parentCustId := c.Locals("parent_cust_id").(string)
	distPriceGrpId := c.Locals("dist_price_grp_id").(int)
	// log.Println("distPriceGrpId", distPriceGrpId)
	switch dataFilter.Mode {
	case "search":
		data, total, lastPage, err = controller.ProductDistService.SearchList(dataFilter, custId, parentCustId)
		if err != nil {
			log.Println("ProductDistController, List, data, err:", err.Error())
			responsePayload.Setmsg(err.Error())
			return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
		}
		err = structs.Automapper(productSearch, &data)
		if err != nil {
			log.Println("ProductDistController, List, Automapper data, err:", err.Error())
			responsePayload.Setmsg(err.Error())
			return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
		}

	case "lookup":
		dataFilter.DistPriceGroupId = distPriceGrpId
		data, total, lastPage, err = controller.ProductDistService.LookupList(dataFilter, custId, parentCustId)
		if err != nil {
			log.Println("ProductDistController, List, data, err:", err.Error())
			responsePayload.Setmsg(err.Error())
			return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
		}
		err = structs.Automapper(productLookup, &data)
		if err != nil {
			log.Println("ProductDistController, List, Automapper data, err:", err.Error())
			responsePayload.Setmsg(err.Error())
			return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
		}

	default:
		data, total, lastPage, err = controller.ProductDistService.List(dataFilter, custId, parentCustId)
		if err != nil {
			log.Println("ProductDistController, List, data, err:", err.Error())
			responsePayload.Setmsg(err.Error())
			return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
		}
		err = structs.Automapper(products, &data)
		if err != nil {
			log.Println("ProductDistController, List, Automapper data, err:", err.Error())
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

func (controller *ProductDistController) Create(c *fiber.Ctx) error {
	var request entity.CreateProductBody
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	if err := c.BodyParser(&request); err != nil {
		log.Println("ProductDistController, Create, BodyParser:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	custId := c.Locals("cust_id").(string)
	userId := c.Locals("user_id").(int64)
	// log.Println("ProductDistController, Create, CustId:", custId)

	request.CustId = custId
	request.CreatedBy = userId

	errs := controller.validator.ValidateStruct(request, headerAcceptLang)
	if errs != nil {
		log.Println("ProductDistController, Create, ValidateStruct, errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	_, err := controller.ProductDistService.Store(request)
	if err != nil {
		log.Println("ProductDistController, Create, Store, err:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())

	}

	responsePayload.Setmsg("Created Successfully")
	return c.Status(fiber.StatusCreated).JSON(responsePayload.GetRespPayload())
}

func (controller *ProductDistController) Update(c *fiber.Ctx) error {
	var (
		params  entity.UpdateProductDistParams
		request entity.UpdateProductDistRequest
	)
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	if err := c.ParamsParser(&params); err != nil {
		log.Println("ProductDistController, Update, ParamsParser(params):", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	errs := controller.validator.ValidateStruct(params, headerAcceptLang)
	if errs != nil {
		log.Println("ProductDistController, Update, ValidateStruct(params), errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	if err := c.BodyParser(&request); err != nil {
		log.Println("ProductDistController, Update, BodyParser(request), err:", err.Error())
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	custId := c.Locals("cust_id").(string)
	userId := c.Locals("user_id").(int64)
	// log.Println("ProductDistController, Update, CustId:", custId)

	request.CustId = custId
	request.UpdatedBy = userId

	errs = controller.validator.ValidateStruct(request, headerAcceptLang)
	if errs != nil {
		log.Println("ProductDistController, Update, ValidateStruct(request), errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	err := controller.ProductDistService.Update(params.ProductId, request)
	if err != nil {
		log.Println("ProductDistController, Update, Service.Update, err:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	responsePayload.Setmsg("Updated Successfully")
	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}

func (controller *ProductDistController) Delete(c *fiber.Ctx) error {
	var params entity.DeleteProductDistParams
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	if err := c.ParamsParser(&params); err != nil {
		log.Println("ProductDistController, Delete, ParamsParser, err:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	errs := controller.validator.ValidateStruct(params, headerAcceptLang)
	if errs != nil {
		log.Println("ProductDistController, Delete, ValidateStruct, errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	custId := c.Locals("cust_id").(string)
	userId := c.Locals("user_id").(int64)
	// log.Println("ProductDistController, Delete, CustId:", custId)

	err := controller.ProductDistService.Delete(custId, params.ProductId, userId)
	if err != nil {
		log.Println("ProductDistController, Delete, Service.Delete, err:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())

	}
	responsePayload.Setmsg("Deleted Successfully")
	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}
