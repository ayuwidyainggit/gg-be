package controller

import (
	"log"
	"master/entity"
	"master/pkg/constant"
	"master/pkg/jwthelper"
	"master/pkg/middleware"
	"master/pkg/responsebuild"
	"master/pkg/validation"
	"master/service"

	"github.com/gofiber/fiber/v2"
)

type SubBrand1Controller struct {
	SubBrand1Service service.SubBrand1Service
	validator        *validation.Validate
}

func NewSubBrand1Controller(subBrand1Service service.SubBrand1Service, validator *validation.Validate) *SubBrand1Controller {
	return &SubBrand1Controller{
		SubBrand1Service: subBrand1Service,
		validator:        validator,
	}
}

func (controller *SubBrand1Controller) Route(app *fiber.App) {
	qParamId := ":sbrand1_id"
	subBrand1RouteV1 := app.Group("/v1/sub-brand1", middleware.JWTProtected())
	subBrand1RouteV1.Get("/"+qParamId, controller.Detail)
	subBrand1RouteV1.Get("", controller.List)
	subBrand1RouteV1.Post("", controller.Create)
	subBrand1RouteV1.Patch("/"+qParamId, controller.Update)
	subBrand1RouteV1.Delete("/"+qParamId, controller.Delete)

	// New endpoint for /v1/sub-brand
	subBrandRouteV1 := app.Group("/v1/sub-brand", middleware.JWTProtected())
	subBrandRouteV1.Get("", controller.SubBrandList)
}

func (controller *SubBrand1Controller) Detail(c *fiber.Ctx) error {
	var params entity.DetailSubBrand1Params
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	if err := c.ParamsParser(&params); err != nil {
		log.Println("SubBrand1Controller, Detail, ParamsParser:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	errs := controller.validator.ValidateStruct(params, headerAcceptLang)
	if errs != nil {
		log.Println("SubBrand1Controller, Detail, ValidateStruct(params), errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	custId := c.Locals("parent_cust_id").(string)
	// log.Println("SubBrand1Controller, Detail, CustId:", claims.CustId)

	data, err := controller.SubBrand1Service.Detail(params.Sbrand1Id, custId)
	if err != nil {
		log.Println("SubBrand1Controller, Detail, FindOneBySubBrand1Id, err:", err.Error())
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

func (controller *SubBrand1Controller) List(c *fiber.Ctx) error {
	var (
		err        error
		dataFilter entity.SubBrand1QueryFilter
		data       interface{}
		total      int
		lastPage   int
	)
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	if err := c.QueryParser(&dataFilter); err != nil {
		log.Println("SubBrand1Controller, List, query parser filter:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	custId := c.Locals("parent_cust_id").(string)
	// log.Println("SubBrand1Controller, List, CustId:", claims.CustId)

	switch dataFilter.Mode {
	case "lookup":
		data, total, lastPage, err = controller.SubBrand1Service.LookupList(dataFilter, custId)
		if err != nil {
			log.Println("SubBrand1Controller, LookupList, data, err:", err.Error())
			responsePayload.Setmsg(err.Error())
			return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
		}
	case "material_group":
		data, total, lastPage, err = controller.SubBrand1Service.MatGroupList(dataFilter, custId)
		if err != nil {
			log.Println("SubBrand1Controller, MatGroup, data, err:", err.Error())
			responsePayload.Setmsg(err.Error())
			return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
		}
	default:
		data, total, lastPage, err = controller.SubBrand1Service.List(dataFilter, custId)
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

func (controller *SubBrand1Controller) Create(c *fiber.Ctx) error {
	var request entity.CreateSubBrand1Body
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	if err := c.BodyParser(&request); err != nil {
		log.Println("SubBrand1Controller, Create, BodyParser:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	claims, err := jwthelper.ExtractTokenMetadata(c)
	if err != nil {
		log.Println("SubBrand1Controller, Create, ExtractTokenMetadata, err:", err.Error())
		// Return status 500 and JWT parse error.
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(responsePayload.GetRespPayload())
	}
	// log.Println("SubBrand1Controller, Create, CustId:", claims.CustId)

	request.CustId = claims.CustId
	request.CreatedBy = int64(claims.UserId)

	errs := controller.validator.ValidateStruct(request, headerAcceptLang)
	if errs != nil {
		log.Println("SubBrand1Controller, Create, ValidateStruct, errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	_, err = controller.SubBrand1Service.Store(request)
	if err != nil {
		log.Println("SubBrand1Controller, Create, Store, err:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	responsePayload.Setmsg(constant.SUCCESSFULLY_ADDED)
	return c.Status(fiber.StatusCreated).JSON(responsePayload.GetRespPayload())
}

func (controller *SubBrand1Controller) Update(c *fiber.Ctx) error {
	var (
		params  entity.UpdateSubBrand1Params
		request entity.UpdateSubBrand1Request
	)
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	if err := c.ParamsParser(&params); err != nil {
		log.Println("SubBrand1Controller, Update, ParamsParser(params):", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	errs := controller.validator.ValidateStruct(params, headerAcceptLang)
	if errs != nil {
		log.Println("SubBrand1Controller, Update, ValidateStruct(params), errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	if err := c.BodyParser(&request); err != nil {
		log.Println("SubBrand1Controller, Update, BodyParser(request), err:", err.Error())
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	claims, err := jwthelper.ExtractTokenMetadata(c)
	if err != nil {
		log.Println("SubBrand1Controller, Update, ExtractTokenMetadata, err:", err.Error())
		// Return status 500 and JWT parse error.
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(responsePayload.GetRespPayload())
	}
	// log.Println("SubBrand1Controller, Update, CustId:", claims.CustId)

	request.CustId = claims.CustId
	request.UpdatedBy = claims.UserId

	errs = controller.validator.ValidateStruct(request, headerAcceptLang)
	if errs != nil {
		log.Println("SubBrand1Controller, Update, ValidateStruct(request), errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	err = controller.SubBrand1Service.Update(params.Sbrand1Id, request)
	if err != nil {
		log.Println("SubBrand1Controller, Update, Service.Update, err:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	responsePayload.Setmsg(constant.SUCCESSFULLY_UPDATED)
	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}

func (controller *SubBrand1Controller) Delete(c *fiber.Ctx) error {
	var params entity.DeleteSubBrand1Params
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	if err := c.ParamsParser(&params); err != nil {
		log.Println("SubBrand1Controller, Delete, ParamsParser, err:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	errs := controller.validator.ValidateStruct(params, headerAcceptLang)
	if errs != nil {
		log.Println("SubBrand1Controller, Delete, ValidateStruct, errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	claims, err := jwthelper.ExtractTokenMetadata(c)
	if err != nil {
		log.Println("SubBrand1Controller, Delete, ExtractTokenMetadata, err:", err.Error())
		// Return status 500 and JWT parse error.
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(responsePayload.GetRespPayload())
	}
	// log.Println("SubBrand1Controller, Delete, CustId:", claims.CustId)

	err = controller.SubBrand1Service.Delete(claims.CustId, params.Sbrand1Id, claims.UserId)
	if err != nil {
		log.Println("SubBrand1Controller, Delete, Service.Delete, err:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}
	responsePayload.Setmsg("Deleted Successfully")
	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())

}

func (controller *SubBrand1Controller) SubBrandList(c *fiber.Ctx) error {
	var (
		err        error
		dataFilter entity.SubBrandQueryFilter
		data       []entity.SubBrandResponse
		total      int
		lastPage   int
	)
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	if err := c.QueryParser(&dataFilter); err != nil {
		log.Println("SubBrand1Controller, SubBrandList, query parser filter:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	custId := c.Locals("parent_cust_id").(string)
	dataFilter.CustId = c.Locals("cust_id").(string)
	dataFilter.ParentCustId = custId

	// Set default values
	if dataFilter.Page == 0 {
		dataFilter.Page = 1
	}
	if dataFilter.Limit == 0 {
		dataFilter.Limit = 9999
	}
	if dataFilter.Sort == "" {
		dataFilter.Sort = "created_date:desc"
	}

	data, total, lastPage, err = controller.SubBrand1Service.SubBrandList(dataFilter, custId)
	if err != nil {
		log.Println("SubBrand1Controller, SubBrandList, data, err:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	if len(data) == 0 {
		responsePayload.Setmsg(constant.NO_DATA)
		responsePayload.Setdata(nil)
	} else {
		responsePayload.Setmsg(constant.SUCCESS_NO_DATA_DISPLAYED)
		responsePayload.Setdata(data)
	}

	responsePayload.Setpaging(entity.Pagination{
		TotalRecord: total,
		PageCurrent: dataFilter.Page,
		PageLimit:   dataFilter.Limit,
		PageTotal:   lastPage,
	})

	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}
