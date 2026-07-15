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

type SubBrand2Controller struct {
	SubBrand2Service service.SubBrand2Service
	validator        *validation.Validate
}

func NewSubBrand2Controller(consProductService service.SubBrand2Service, validator *validation.Validate) SubBrand2Controller {
	return SubBrand2Controller{
		SubBrand2Service: consProductService,
		validator:        validator,
	}
}

func (controller *SubBrand2Controller) Route(app *fiber.App) {
	qParamId := ":sbrand2_id"
	sbrand2RouteV1 := app.Group("/v1/sub-brand2", middleware.JWTProtected())
	sbrand2RouteV1.Get("/"+qParamId, controller.Detail)
	sbrand2RouteV1.Get("/", controller.List)
	sbrand2RouteV1.Post("", controller.Create)
	sbrand2RouteV1.Patch("/"+qParamId, controller.Update)
	sbrand2RouteV1.Delete("/"+qParamId, controller.Delete)
}

func (controller *SubBrand2Controller) Detail(c *fiber.Ctx) error {
	var params entity.DetailSubBrand2Params
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	if err := c.ParamsParser(&params); err != nil {
		log.Println("SubBrand2Controller, Detail, ParamsParser:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	errs := controller.validator.ValidateStruct(params, headerAcceptLang)
	if errs != nil {
		log.Println("SubBrand2Controller, Detail, ValidateStruct(params), errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	custId := c.Locals("parent_cust_id").(string)
	data, err := controller.SubBrand2Service.Detail(params.SBrand2Id, custId)
	if err != nil {
		log.Println("SubBrand2Controller, Detail, FindOneBySubBrand2Id, err:", err.Error())
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

func (controller *SubBrand2Controller) List(c *fiber.Ctx) error {
	var dataFilter entity.GeneralQueryFilter
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	if err := c.QueryParser(&dataFilter); err != nil {
		log.Println("SubBrand2Controller, List, query parser filter:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	custId := c.Locals("parent_cust_id").(string)
	data, total, lastPage, err := controller.SubBrand2Service.List(dataFilter, custId)
	if err != nil {
		log.Println("SubBrand2Controller, List, data, err:", err.Error())
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

	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())

}

func (controller *SubBrand2Controller) Create(c *fiber.Ctx) error {
	var request entity.CreateSubBrand2Body
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	if err := c.BodyParser(&request); err != nil {
		log.Println("SubBrand2Controller, Create, BodyParser:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	claims, err := jwthelper.ExtractTokenMetadata(c)
	if err != nil {
		log.Println("SubBrand2Controller, Create, ExtractTokenMetadata, err:", err.Error())
		// Return status 500 and JWT parse error.
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(responsePayload.GetRespPayload())
	}
	// log.Println("SubBrand2Controller, Create, CustId:", claims.CustId)

	request.CustId = claims.CustId
	request.CreatedBy = int64(claims.UserId)

	errs := controller.validator.ValidateStruct(request, headerAcceptLang)
	if errs != nil {
		log.Println("SubBrand2Controller, Create, ValidateStruct, errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	_, err = controller.SubBrand2Service.Store(request)
	if err != nil {
		log.Println("SubBrand2Controller, Create, Store, err:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	responsePayload.Setmsg(constant.SUCCESSFULLY_ADDED)
	return c.Status(fiber.StatusCreated).JSON(responsePayload.GetRespPayload())
}

func (controller *SubBrand2Controller) Update(c *fiber.Ctx) error {
	var (
		params  entity.UpdateSubBrand2Params
		request entity.UpdateSubBrand2Request
	)
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	if err := c.ParamsParser(&params); err != nil {
		log.Println("SubBrand2Controller, Update, ParamsParser(params):", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	errs := controller.validator.ValidateStruct(params, headerAcceptLang)
	if errs != nil {
		log.Println("SubBrand2Controller, Update, ValidateStruct(params), errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	if err := c.BodyParser(&request); err != nil {
		log.Println("SubBrand2Controller, Update, BodyParser(request), err:", err.Error())
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	claims, err := jwthelper.ExtractTokenMetadata(c)
	if err != nil {
		log.Println("SubBrand2Controller, Update, ExtractTokenMetadata, err:", err.Error())
		// Return status 500 and JWT parse error.
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(responsePayload.GetRespPayload())
	}
	// log.Println("SubBrand2Controller, Update, CustId:", claims.CustId)

	request.CustId = claims.CustId
	request.UpdatedBy = claims.UserId

	errs = controller.validator.ValidateStruct(request, headerAcceptLang)
	if errs != nil {
		log.Println("SubBrand2Controller, Update, ValidateStruct(request), errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	err = controller.SubBrand2Service.Update(params.SBrand2Id, request)
	if err != nil {
		log.Println("SubBrand2Controller, Update, Service.Update, err:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	responsePayload.Setmsg(constant.SUCCESSFULLY_UPDATED)
	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}

func (controller *SubBrand2Controller) Delete(c *fiber.Ctx) error {
	var params entity.DeleteSubBrand2Params
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	if err := c.ParamsParser(&params); err != nil {
		log.Println("SubBrand2Controller, Delete, ParamsParser, err:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	errs := controller.validator.ValidateStruct(params, headerAcceptLang)
	if errs != nil {
		log.Println("SubBrand2Controller, Delete, ValidateStruct, errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	claims, err := jwthelper.ExtractTokenMetadata(c)
	if err != nil {
		log.Println("SubBrand2Controller, Delete, ExtractTokenMetadata, err:", err.Error())
		// Return status 500 and JWT parse error.
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(responsePayload.GetRespPayload())
	}
	// log.Println("SubBrand2Controller, Delete, CustId:", claims.CustId)

	err = controller.SubBrand2Service.Delete(claims.CustId, params.SBrand2Id, claims.UserId)
	if err != nil {
		log.Println("SubBrand2Controller, Delete, Service.Delete, err:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	responsePayload.Setmsg("Deleted Successfully")
	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}
