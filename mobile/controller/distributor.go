package controller

import (
	"mobile/entity"
	"mobile/pkg/constant"
	"mobile/pkg/middleware"
	"mobile/pkg/responsebuild"
	"mobile/pkg/validation"
	"mobile/service"

	"github.com/gofiber/fiber/v2/log"

	"github.com/gofiber/fiber/v2"
)

type DistributorController struct {
	DistributorService service.DistributorService
	validator          *validation.Validate
}

func NewDistributorController(
	distributorService service.DistributorService,
	validator *validation.Validate,
) *DistributorController {
	return &DistributorController{
		DistributorService: distributorService,
		validator:          validator,
	}
}

func (controller *DistributorController) Route(app *fiber.App) {
	mobileDistributorRouteV1 := app.Group("/v1/distributor", middleware.JWTProtected())
	mobileDistributorRouteV1.Get("", controller.MobileDistributorList)
}

func (controller *DistributorController) MobileDistributorList(c *fiber.Ctx) error {
	var (
		headerAcceptLang string
		dataFilter       entity.MobileDistributorListQueryFilter
	)

	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	if err := c.QueryParser(&dataFilter); err != nil {
		log.Error("DistributorController, MobileDistributorList, query parser filter:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	if dataFilter.Page <= 0 {
		dataFilter.Page = 1
	}
	if dataFilter.Limit <= 0 {
		dataFilter.Limit = 5
	}
	if dataFilter.Sort == "" {
		dataFilter.Sort = "distributor_code:asc"
	}

	errs := controller.validator.ValidateStruct(dataFilter, headerAcceptLang)
	if errs != nil {
		log.Error("DistributorController, MobileDistributorList, ValidateStruct, errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	custId := c.Locals("cust_id").(string)

	principal, err := controller.DistributorService.GetPrincipalInfo(c.UserContext(), custId)
	if err != nil {
		log.Warnf("Error get principal info, err: %s", err.Error())
	}

	data, total, lastPage, err := controller.DistributorService.MobileDistributorList(c.UserContext(), dataFilter, custId)
	if err != nil {
		log.Error("DistributorController, MobileDistributorList, err:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusNotFound).JSON(responsePayload.GetRespPayload())
	}

	if len(data) == 0 {
		responsePayload.Setmsg(constant.STATUS_DB_NOT_FOUND)
		responsePayload.Setdata(nil)
		responsePayload.Setpaging(entity.Pagination{
			TotalRecord: 0,
			PageCurrent: dataFilter.Page,
			PageLimit:   dataFilter.Limit,
			PageTotal:   0,
		})
		resp := responsePayload.GetRespPayload()
		resp.CustID = principal.CustID
		resp.CustName = principal.CustName
		resp.DistributorID = principal.DistributorID

		return c.Status(fiber.StatusOK).JSON(resp)
	}

	responsePayload.Setmsg(constant.STATUS_OK)
	responsePayload.Setdata(data)
	responsePayload.Setpaging(entity.Pagination{
		TotalRecord: total,
		PageCurrent: dataFilter.Page,
		PageLimit:   dataFilter.Limit,
		PageTotal:   lastPage,
	})

	resp := responsePayload.GetRespPayload()
	resp.CustID = principal.CustID
	resp.CustName = principal.CustName
	resp.DistributorID = principal.DistributorID
	return c.Status(fiber.StatusOK).JSON(resp)
}
