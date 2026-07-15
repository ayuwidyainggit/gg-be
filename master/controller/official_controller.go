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

type OfficialController struct {
	OfficialService service.OfficialService
	validator       *validation.Validate
}

func NewOfficialController(officialService service.OfficialService, validator *validation.Validate) *OfficialController {
	return &OfficialController{
		OfficialService: officialService,
		validator:       validator,
	}
}

func (controller *OfficialController) Route(app *fiber.App) {
	qParamId := ":official_id"
	officialsRouteV1 := app.Group("/v1/officials", middleware.JWTProtected())
	officialsRouteV1.Get("/hierarchy", controller.OfficialHierarchy)
	officialsRouteV1.Get("/"+qParamId, controller.Detail)
	officialsRouteV1.Get("", controller.List)
	officialsRouteV1.Post("hierarchy", controller.CreateHierarchy)
}

func (controller *OfficialController) Detail(c *fiber.Ctx) error {
	var params entity.DetailOfficialParams
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	if err := c.ParamsParser(&params); err != nil {
		log.Println("OfficialController, Detail, ParamsParser:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	errs := controller.validator.ValidateStruct(params, headerAcceptLang)
	if errs != nil {
		log.Println("OfficialController, Detail, ValidateStruct(params), errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	custId := c.Locals("parent_cust_id").(string)

	data, err := controller.OfficialService.Detail(params.OfficialId, custId)
	if err != nil {
		log.Println("OfficialController, Detail, FindOneByOfficialId, err:", err.Error())
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

func (controller *OfficialController) List(c *fiber.Ctx) error {
	var (
		err            error
		dataFilter     entity.OfficialQueryFilter
		data           interface{}
		total          int
		lastPage       int
		officialList   []entity.OfficialListResponse
		officialLookup []entity.OfficialLookupResponse
	)
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	if err := c.QueryParser(&dataFilter); err != nil {
		log.Println("OfficialController, List, query parser filter:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	// custId := c.Locals("parent_cust_id").(string)
	custId := c.Locals("cust_id").(string)

	switch dataFilter.Mode {
	case "lookup":
		data, total, lastPage, err = controller.OfficialService.LookupList(dataFilter, custId)
		if err != nil {
			responsePayload.Setmsg(err.Error())
			return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
		}
		err = structs.Automapper(officialLookup, &data)
		if err != nil {
			log.Println("OfficialController, Lookup, Automapper data, err:", err.Error())
			responsePayload.Setmsg(err.Error())
			return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
		}
	default:
		data, total, lastPage, err = controller.OfficialService.List(dataFilter, custId)
		if err != nil {
			log.Println("OfficialController, List, data, err:", err.Error())
			responsePayload.Setmsg(err.Error())
			return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
		}
		err = structs.Automapper(officialList, &data)
		if err != nil {
			log.Println("OfficialController, List, Automapper data, err:", err.Error())
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

func (controller *OfficialController) OfficialHierarchy(c *fiber.Ctx) error {
	var (
		err        error
		dataFilter entity.OfficialQueryFilter
	)

	// custId := c.Locals("parent_cust_id").(string)
	custId := c.Locals("cust_id").(string)

	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string))
	if err := c.QueryParser(&dataFilter); err != nil {
		// log.Println("OfficialController, Lookup, Automapper data, err:", err.Error())
		log.Println("OfficialController, QueryParser, Lookup:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	data, err := controller.OfficialService.OfficialHierarchy(dataFilter, custId)
	if err != nil {
		log.Println("OfficialController, OfficialHierarchy, data, err:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}
	responsePayload.Setdata(data)
	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())

}

func (controller *OfficialController) CreateHierarchy(c *fiber.Ctx) error {
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	var request entity.CreateOfficialBodyHierarchy
	if err := c.BodyParser(&request); err != nil {
		log.Println("OfficialController, Create, BodyParser:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	// custId := c.Locals("parent_cust_id").(string)
	custId := c.Locals("cust_id").(string)
	userId := c.Locals("user_id").(int64)
	// log.Println("OfficialController, Create, CustId:", custId)

	request.CustId = custId
	request.CreatedBy = userId

	errs := controller.validator.ValidateStruct(request, headerAcceptLang)
	if errs != nil {
		log.Println("OfficialController, Create, ValidateStruct, errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	_, err := controller.OfficialService.StoreHierarchy(request)
	if err != nil {
		log.Println("OfficialController, Create, Store, err:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	responsePayload.Setmsg(constant.SUCCESSFULLY_ADDED)
	return c.Status(fiber.StatusCreated).JSON(responsePayload.GetRespPayload())
}
