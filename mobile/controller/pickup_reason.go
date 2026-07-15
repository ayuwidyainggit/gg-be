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

type PickupReasonController struct {
	PickupReasonService service.PickupReasonService
	validator           *validation.Validate
}

func NewPickupReasonController(pickupReasonService service.PickupReasonService, validator *validation.Validate) *PickupReasonController {
	return &PickupReasonController{
		PickupReasonService: pickupReasonService,
		validator:           validator,
	}
}

func (controller *PickupReasonController) Route(app *fiber.App) {
	// qParamId := ":pickup_reason_id"
	pickupReasonsRouteV1 := app.Group("/v1/pickup-reasons", middleware.JWTProtected())
	pickupReasonsRouteV1.Get("", controller.List)
}

func (controller *PickupReasonController) List(c *fiber.Ctx) error {
	var (
		err                error
		dataFilter         entity.PickupReasonQueryFilter
		data               interface{}
		total              int64
		lastPage           int
		pickupReason       []entity.PickupReasonResponse
		pickupReasonLookup []entity.PickupReasonLookupResponse
	)

	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), constant.HEADER_ACCEPT_LANG)
	if err := c.QueryParser(&dataFilter); err != nil {
		log.Error("PickupReasonController, List, query parser filter:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	custId := c.Locals("parent_cust_id").(string)

	switch dataFilter.Mode {
	case "lookup":
		data, total, lastPage, err = controller.PickupReasonService.LookupList(dataFilter, custId)
		if err != nil {
			log.Error("PickupReasonController, List, data, err:", err.Error())
			responsePayload.Setmsg(err.Error())
			return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
		}
		err = structs.Automapper(pickupReasonLookup, &data)
		if err != nil {
			log.Error("PickupReasonController, List, data, err:", err.Error())
			responsePayload.Setmsg(err.Error())
			return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
		}
	default:
		data, total, lastPage, err = controller.PickupReasonService.List(dataFilter, custId)
		if err != nil {
			log.Error("PickupReasonController, List, data, err:", err.Error())
			responsePayload.Setmsg(err.Error())
			return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
		}
		err = structs.Automapper(pickupReason, &data)
		if err != nil {
			log.Error("PickupReasonController, List, data, err:", err.Error())
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
