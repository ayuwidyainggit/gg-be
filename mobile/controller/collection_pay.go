package controller

import (
	"mobile/entity"
	"mobile/pkg/constant"
	"mobile/pkg/middleware"
	"mobile/pkg/responsebuild"
	"mobile/pkg/validation"
	"mobile/service"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

type CollectionPayController struct {
	CollectionPayService service.CollectionPayService
	validator            *validation.Validate
}

func NewCollectionPayController(collectionPayService service.CollectionPayService, validator *validation.Validate) *CollectionPayController {
	return &CollectionPayController{
		CollectionPayService: collectionPayService,
		validator:            validator,
	}
}

func (controller *CollectionPayController) Route(app *fiber.App) {
	collectionRouteV1 := app.Group("/v1/collection-pay", middleware.JWTProtected())
	collectionRouteV1.Post("", controller.Create)
	collectionRouteV1.Get("", controller.List)
	collectionRouteV1.Post("/nopayment", controller.NoPayment)
}

func (controller *CollectionPayController) Create(c *fiber.Ctx) error {
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), constant.HEADER_ACCEPT_LANG)

	var request entity.CreateCollectionPayRequest
	if err := c.BodyParser(&request); err != nil {
		log.Error("CollectionPayController, Create, BodyParser:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	errs := controller.validator.ValidateStruct(request, constant.HEADER_ACCEPT_LANG)
	if errs != nil {
		log.Error("CollectionPayController, Create, ValidateStruct, errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}
	custId := c.Locals("cust_id").(string)
	empId := c.Locals("emp_id").(int64)
	userId := c.Locals("user_id").(int64)
	request.CustID = custId
	request.EmpID = empId
	request.UserID = userId
	data, err := controller.CollectionPayService.Store(c.Context(), request)
	if err != nil {
		log.Error("CollectionPayController, Create, Store, err:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}
	responsePayload.Setdata(data)
	responsePayload.Setmsg("Created Successfully")
	return c.Status(fiber.StatusCreated).JSON(responsePayload.GetRespPayload())
}

func (controller *CollectionPayController) List(c *fiber.Ctx) error {
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), constant.HEADER_ACCEPT_LANG)

	var request entity.CollectionPayQueryFilter
	if err := c.QueryParser(&request); err != nil {
		log.Error("CollectionPayController, List, BodyParser:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	errs := controller.validator.ValidateStruct(request, constant.HEADER_ACCEPT_LANG)
	if errs != nil {
		log.Error("CollectionPayController, List, ValidateStruct, errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	data, total, lastPage, err := controller.CollectionPayService.List(c.Context(), request)
	if err != nil {
		log.Error("CollectionPayController, List, data, err:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}
	responsePayload.Setdata(data)
	responsePayload.Setpaging(entity.Pagination{
		TotalRecord: total,
		PageCurrent: request.Page,
		PageLimit:   request.Limit,
		PageTotal:   int(lastPage),
	})

	responsePayload.Setmsg(constant.STATUS_OK)
	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}

func (controller *CollectionPayController) NoPayment(c *fiber.Ctx) error {
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), constant.HEADER_ACCEPT_LANG)

	var request entity.CreateNoPaymentRequest
	if err := c.BodyParser(&request); err != nil {
		log.Error("CollectionPayController, NoPayment, BodyParser:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	errs := controller.validator.ValidateStruct(request, constant.HEADER_ACCEPT_LANG)
	if errs != nil {
		log.Error("CollectionPayController, NoPayment, ValidateStruct, errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}
	custId := c.Locals("cust_id").(string)
	empId := c.Locals("emp_id").(int64)
	userId := c.Locals("user_id").(int64)
	request.CustID = custId
	request.EmpID = empId
	request.UserID = userId
	err := controller.CollectionPayService.StoreNoPayment(c.Context(), request)
	if err != nil {
		log.Error("CollectionPayController, NoPayment, Store, err:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}
	responsePayload.Setmsg("Created Successfully")
	return c.Status(fiber.StatusCreated).JSON(responsePayload.GetRespPayload())
}
