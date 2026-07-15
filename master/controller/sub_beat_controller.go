package controller

import (
	"log"
	"master/entity"
	"master/pkg/constant"
	"master/pkg/middleware"
	"master/pkg/responsebuild"
	"master/pkg/validation"
	"master/service"

	"github.com/gofiber/fiber/v2"
)

type SubBeatController struct {
	SubBeatService service.SubBeatService
	validator      *validation.Validate
}

func NewSubBeatController(subBeatService service.SubBeatService, validator *validation.Validate) SubBeatController {
	return SubBeatController{
		SubBeatService: subBeatService,
		validator:      validator,
	}
}

func (controller *SubBeatController) Route(app *fiber.App) {
	qParamId := ":sbeat_id"
	subBeatRouteV1 := app.Group("/v1/sub-beats", middleware.JWTProtected())
	subBeatRouteV1.Get("/"+qParamId, controller.Detail)
	subBeatRouteV1.Get("", controller.List)
	subBeatRouteV1.Post("", controller.Create)
	subBeatRouteV1.Patch("/"+qParamId, controller.Update)
	subBeatRouteV1.Delete("/"+qParamId, controller.Delete)
}

func (controller *SubBeatController) Detail(c *fiber.Ctx) error {
	var params entity.DetailSubBeatParams
	if err := c.ParamsParser(&params); err != nil {
		log.Println("SubBeatController, Detail, ParamsParser:", err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(entity.ApiResponse{
			RequestId: c.Locals("requestid").(string),
			Message:   err.Error(),
		})
	}
	var headerAcceptLang string
	errs := controller.validator.ValidateStruct(params, headerAcceptLang)
	if errs != nil {
		log.Println("SubBeatController, Detail, ValidateStruct(params), errs:", errs)
		return c.Status(fiber.StatusBadRequest).JSON(entity.ApiResponse{
			RequestId: c.Locals("requestid").(string),
			Message:   fiber.ErrBadRequest.Message,
			Errors:    errs,
		})
	}

	custId := c.Locals("cust_id").(string)

	data, err := controller.SubBeatService.Detail(params.SbeatId, custId)
	if err != nil {
		log.Println("SubBeatController, Detail, FindOneBySubBeatId, err:", err.Error())
		statusCode := fiber.StatusBadRequest
		errMsg := err.Error()
		if err.Error() == "sql: no rows in result set" {
			statusCode = fiber.StatusNotFound
			errMsg = "Not found"
		}
		return c.Status(statusCode).JSON(entity.ApiResponse{
			RequestId: c.Locals("requestid").(string),
			Message:   errMsg,
		})
	}

	return c.JSON(entity.ApiResponse{
		RequestId: c.Locals("requestid").(string),
		Data:      data,
	})
}

func (controller *SubBeatController) List(c *fiber.Ctx) error {
	var (
		err        error
		dataFilter entity.GeneralQueryFilter
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
		log.Println("SubBeatController, List, query parser filter:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	custId := c.Locals("cust_id").(string)

	switch dataFilter.Mode {
	case "lookup":
		data, total, lastPage, err = controller.SubBeatService.LookupList(dataFilter, custId)
		if err != nil {
			log.Println("SubBeatController, LookupList, data, err:", err.Error())
			return c.Status(fiber.StatusBadRequest).JSON(entity.ApiResponse{
				RequestId: c.Locals("requestid").(string),
				Message:   err.Error(),
			})
		}
	default:
		data, total, lastPage, err = controller.SubBeatService.List(dataFilter, custId)
		if err != nil {
			log.Println("SubBeatController, List, data, err:", err.Error())
			return c.Status(fiber.StatusBadRequest).JSON(entity.ApiResponse{
				RequestId: c.Locals("requestid").(string),
				Message:   err.Error(),
			})
		}
	}

	return c.Status(fiber.StatusOK).JSON(entity.ApiResponse{
		RequestId: c.Locals("requestid").(string),
		Data:      data,
		Paging: entity.Pagination{
			TotalRecord: total,
			PageCurrent: dataFilter.Page,
			PageLimit:   dataFilter.Limit,
			PageTotal:   lastPage,
		},
	})
}

func (controller *SubBeatController) Create(c *fiber.Ctx) error {
	var request entity.CreateSubBeatBody
	if err := c.BodyParser(&request); err != nil {
		log.Println("SubBeatController, Create, BodyParser:", err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(entity.ApiResponse{
			RequestId: c.Locals("requestid").(string),
			Message:   err.Error(),
		})
	}

	custId := c.Locals("cust_id").(string)
	userId := c.Locals("user_id").(int64)
	// log.Println("SubBeatController, Create, CustId:", custId)

	request.CustId = custId
	request.CreatedBy = userId
	var headerAcceptLang string
	errs := controller.validator.ValidateStruct(request, headerAcceptLang)
	if errs != nil {
		log.Println("SubBeatController, Create, ValidateStruct, errs:", errs)
		return c.Status(fiber.StatusBadRequest).JSON(entity.ApiResponse{
			RequestId: c.Locals("requestid").(string),
			Message:   fiber.ErrBadRequest.Message,
			Errors:    errs,
		})
	}

	_, err := controller.SubBeatService.Store(request)
	if err != nil {
		log.Println("SubBeatController, Create, Store, err:", err.Error())
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

func (controller *SubBeatController) Update(c *fiber.Ctx) error {
	var (
		params  entity.UpdateSubBeatParams
		request entity.UpdateSubBeatRequest
	)

	if err := c.ParamsParser(&params); err != nil {
		log.Println("SubBeatController, Update, ParamsParser(params):", err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(entity.ApiResponse{
			RequestId: c.Locals("requestid").(string),
			Message:   err.Error(),
		})
	}
	var headerAcceptLang string
	errs := controller.validator.ValidateStruct(params, headerAcceptLang)
	if errs != nil {
		log.Println("SubBeatController, Update, ValidateStruct(params), errs:", errs)
		return c.Status(fiber.StatusBadRequest).JSON(entity.ApiResponse{
			RequestId: c.Locals("requestid").(string),
			Message:   fiber.ErrBadRequest.Message,
			Errors:    errs,
		})
	}

	if err := c.BodyParser(&request); err != nil {
		log.Println("SubBeatController, Update, BodyParser(request), err:", err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(entity.ApiResponse{
			RequestId: c.Locals("requestid").(string),
			Message:   err.Error(),
		})
	}

	custId := c.Locals("cust_id").(string)
	userId := c.Locals("user_id").(int64)
	// log.Println("SubBeatController, Update, CustId:", custId)

	request.CustId = custId
	request.UpdatedBy = userId

	errs = controller.validator.ValidateStruct(request, headerAcceptLang)
	if errs != nil {
		log.Println("SubBeatController, Update, ValidateStruct(request), errs:", errs)
		return c.Status(fiber.StatusBadRequest).JSON(entity.ApiResponse{
			RequestId: c.Locals("requestid").(string),
			Message:   fiber.ErrBadRequest.Message,
			Errors:    errs,
		})
	}

	err := controller.SubBeatService.Update(params.SbeatId, request)
	if err != nil {
		log.Println("SubBeatController, Update, Service.Update, err:", err.Error())
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

func (controller *SubBeatController) Delete(c *fiber.Ctx) error {
	var params entity.DeleteSubBeatParams
	if err := c.ParamsParser(&params); err != nil {
		log.Println("SubBeatController, Delete, ParamsParser, err:", err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(entity.ApiResponse{
			RequestId: c.Locals("requestid").(string),
			Message:   err.Error(),
		})
	}
	var headerAcceptLang string
	errs := controller.validator.ValidateStruct(params, headerAcceptLang)
	if errs != nil {
		log.Println("SubBeatController, Delete, ValidateStruct, errs:", errs)
		return c.Status(fiber.StatusBadRequest).JSON(entity.ApiResponse{
			RequestId: c.Locals("requestid").(string),
			Message:   fiber.ErrBadRequest.Message,
			Errors:    errs,
		})
	}

	custId := c.Locals("cust_id").(string)
	userId := c.Locals("user_id").(int64)
	log.Println("SubBeatController, Delete, CustId:", custId)

	err := controller.SubBeatService.Delete(custId, params.SbeatId, userId)
	if err != nil {
		log.Println("SubBeatController, Delete, Service.Delete, err:", err.Error())
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
