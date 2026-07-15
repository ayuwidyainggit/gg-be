package controller

import (
	"mobile/entity"
	"mobile/pkg/constant"
	"mobile/pkg/middleware"
	"mobile/pkg/responsebuild"
	"mobile/pkg/validation"
	"mobile/service"
	"sort"

	"github.com/gofiber/fiber/v2/log"

	"github.com/gofiber/fiber/v2"
)

type DiscountController struct {
	DiscountService service.DiscountService
	validator       *validation.Validate
}

func NewDiscountController(roService service.DiscountService, validator *validation.Validate) *DiscountController {
	return &DiscountController{
		DiscountService: roService,
		validator:       validator,
	}
}

func (controller *DiscountController) Route(app *fiber.App) {
	qParamId := ":discount_id"
	qParamIdDiscountGroup := ":disc_grp_id"
	roRouteV1 := app.Group("/v1/discounts", middleware.JWTProtected())
	roRouteV1.Post("", controller.Create)
	roRouteV1.Get("/statuses", controller.DiscountStatus)
	roRouteV1.Get("/publish/statuses", controller.PublishDiscountStatus)
	roRouteV1.Get("/"+qParamId, controller.Detail)
	roRouteV1.Get("/discounts-group/"+qParamIdDiscountGroup, controller.DetailGrp)
	roRouteV1.Get("", controller.List)
	roRouteV1.Patch("/"+qParamId, controller.Update)
	roRouteV1.Post("/publish", controller.Publish)
	roRouteV1.Delete("/"+qParamId, controller.Delete)
	roRouteV1.Post("/consult", controller.Consult)
}

func (controller *DiscountController) Create(c *fiber.Ctx) error {
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	var request entity.CreateDiscountBody
	if err := c.BodyParser(&request); err != nil {
		log.Error("DiscountController, Create, BodyParser:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	request.CustID = c.Locals("cust_id").(string)
	request.ParentCustID = c.Locals("parent_cust_id").(string)
	request.CreatedBy = c.Locals("user_fullname").(string)

	errs := controller.validator.ValidateStruct(request, headerAcceptLang)
	if errs != nil {
		log.Error("DiscountController, Create, ValidateStruct, errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	err := controller.DiscountService.Store(request)
	if err != nil {
		log.Error("DiscountController, Create, Store, err:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	responsePayload.Setmsg("Successfully added")
	return c.Status(fiber.StatusCreated).JSON(responsePayload.GetRespPayload())
}

func (controller *DiscountController) Detail(c *fiber.Ctx) error {
	var params entity.DetailDiscountParams
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	if err := c.ParamsParser(&params); err != nil {
		log.Error("DiscountController, Detail, ParamsParser:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	errs := controller.validator.ValidateStruct(params, headerAcceptLang)
	if errs != nil {
		log.Error("DiscountController, Detail, ValidateStruct(params), errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	params.CustID = c.Locals("cust_id").(string)
	params.ParentCustId = c.Locals("parent_cust_id").(string)

	data, err := controller.DiscountService.Detail(params)
	if err != nil {
		log.Error("DiscountController, Detail, err:", err.Error())
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

func (controller *DiscountController) DetailGrp(c *fiber.Ctx) error {
	var params entity.DetailDiscountGrpParams
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	if err := c.ParamsParser(&params); err != nil {
		log.Error("DiscountController, Detail, ParamsParser:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	errs := controller.validator.ValidateStruct(params, headerAcceptLang)
	if errs != nil {
		log.Error("DiscountController, Detail, ValidateStruct(params), errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	params.CustID = c.Locals("cust_id").(string)
	params.ParentCustId = c.Locals("parent_cust_id").(string)

	data, err := controller.DiscountService.DetailGrp(params.DiscountGrpID)
	if err != nil {
		log.Error("DiscountController, Detail, err:", err.Error())
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
	// var (
	// 	dataFilter entity.DetailDiscountGrpParams
	// 	data       []entity.DetailDiscountGrp
	// )

	// var headerAcceptLang string
	// if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
	// 	headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	// }
	// responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)
	// if err := c.QueryParser(&dataFilter); err != nil {
	// 	log.Error("DiscountController, List, query parser filter:", err.Error())
	// 	responsePayload.Setmsg(err.Error())
	// 	return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	// }
	// errs := controller.validator.ValidateStruct(dataFilter, headerAcceptLang)
	// if errs != nil {
	// 	log.Error("DiscountController, Update, ValidateStruct(params), errs:", errs)
	// 	responsePayload.Setmsg(fiber.ErrBadRequest.Message)
	// 	responsePayload.Seterrors(errs)
	// 	return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	// }

	// // dataFilter.CustId = c.Locals("cust_id").(string)
	// dataFilter.ParentCustId = c.Locals("parent_cust_id").(string)
	// // log.Println("BankController, List, CustId:", custId)

	// data, total, lastPage, err := controller.DiscountService.DetailGrp(dataFilter)
	// if err != nil {
	// 	log.Error("DiscountController, List, data, err:", err.Error())
	// 	responsePayload.Setmsg(err.Error())
	// 	return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	// }

	// responsePayload.Setdata(data)
	// responsePayload.Setpaging(entity.Pagination{
	// 	TotalRecord: total,
	// 	PageCurrent: dataFilter.Page,
	// 	PageLimit:   dataFilter.Limit,
	// 	PageTotal:   lastPage,
	// })

	// return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}

func (controller *DiscountController) List(c *fiber.Ctx) error {
	var (
		dataFilter entity.DiscountQueryFilter
		data       []entity.Discount
	)

	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)
	if err := c.QueryParser(&dataFilter); err != nil {
		log.Error("DiscountController, List, query parser filter:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}
	errs := controller.validator.ValidateStruct(dataFilter, headerAcceptLang)
	if errs != nil {
		log.Error("DiscountController, Update, ValidateStruct(params), errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	dataFilter.CustId = c.Locals("cust_id").(string)
	dataFilter.ParentCustId = c.Locals("parent_cust_id").(string)
	// log.Println("BankController, List, CustId:", custId)

	data, total, lastPage, err := controller.DiscountService.List(dataFilter)
	if err != nil {
		log.Error("DiscountController, List, data, err:", err.Error())
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

func (controller *DiscountController) Update(c *fiber.Ctx) error {
	var (
		params  entity.UpdateDiscountParams
		request entity.UpdateDiscountBody
	)
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	if err := c.ParamsParser(&params); err != nil {
		log.Error("DiscountController, Update, ParamsParser(params):", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	errs := controller.validator.ValidateStruct(params, headerAcceptLang)
	if errs != nil {
		log.Error("DiscountController, Update, ValidateStruct(params), errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	if err := c.BodyParser(&request); err != nil {
		log.Error("DiscountController, Update, BodyParser(request), err:", err.Error())
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	request.CustID = c.Locals("cust_id").(string)
	request.ParentCustID = c.Locals("parent_cust_id").(string)
	request.UpdatedBy = c.Locals("user_fullname").(string)

	errs = controller.validator.ValidateStruct(request, headerAcceptLang)
	if errs != nil {
		log.Error("DiscountController, Update, ValidateStruct(request), errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	err := controller.DiscountService.Update(params.DiscountID, request)
	if err != nil {
		log.Error("DiscountController, Update, Service.Update, err:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}
	responsePayload.Setmsg("Updated Successfully")
	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}

func (controller *DiscountController) Delete(c *fiber.Ctx) error {
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)
	var params entity.DetailDiscountParams
	if err := c.ParamsParser(&params); err != nil {
		log.Error("DiscountController, Delete, ParamsParser(params):", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}
	errs := controller.validator.ValidateStruct(params, headerAcceptLang)
	if errs != nil {
		log.Error("DiscountController, Delete, ValidateStruct, errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	deletedBy := c.Locals("user_fullname").(string)

	params.CustID = c.Locals("cust_id").(string)
	params.ParentCustId = c.Locals("parent_cust_id").(string)

	err := controller.DiscountService.Delete(params, deletedBy)
	if err != nil {
		log.Error("DiscountController, Delete, Service.Delete, err:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	responsePayload.Setmsg("Deleted Successfully")
	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}

func (controller *DiscountController) DiscountStatus(c *fiber.Ctx) error {
	discountStatuses := make([]entity.DiscountStatus, 0)

	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	for index, element := range entity.DiscountStatusDesc {
		discountStatus := entity.DiscountStatus{
			DiscountStatusID:   index,
			DiscountStatusDesc: element,
		}
		discountStatuses = append(discountStatuses, discountStatus)
	}

	discountStatusesSorted := make(entity.DiscountStatusDescSlice, 0)
	for _, row := range discountStatuses {
		discountStatusesSorted = append(discountStatusesSorted, row)
	}
	sort.Sort(discountStatusesSorted)
	responsePayload.Setdata(discountStatusesSorted)
	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}

func (controller *DiscountController) PublishDiscountStatus(c *fiber.Ctx) error {
	publishStatuses := make([]entity.PublishStatus, 0)

	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	for index, element := range entity.PublishStatusDesc {
		publishStatus := entity.PublishStatus{
			PublishStatusID:   index,
			PublishStatusDesc: element,
		}
		publishStatuses = append(publishStatuses, publishStatus)
	}

	publishStatusesSorted := make(entity.PublishStatusDescSlice, 0)
	for _, row := range publishStatuses {
		publishStatusesSorted = append(publishStatusesSorted, row)
	}
	sort.Sort(publishStatusesSorted)
	responsePayload.Setdata(publishStatusesSorted)
	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}

func (controller *DiscountController) Publish(c *fiber.Ctx) error {
	var (
		request entity.PublishDiscountBody
	)
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	if err := c.BodyParser(&request); err != nil {
		log.Error("DiscountController, BulkUpdateStatus, BodyParser(request), err:", err.Error())
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	request.CustID = c.Locals("cust_id").(string)
	request.ParentCustID = c.Locals("parent_cust_id").(string)
	request.UpdatedBy = c.Locals("user_fullname").(string)

	errs := controller.validator.ValidateStruct(request, headerAcceptLang)
	if errs != nil {
		log.Error("DiscountController, PublishDiscount, ValidateStruct(request), errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	err := controller.DiscountService.PublishDiscount(request)
	if err != nil {
		log.Error("DiscountController, PublishDiscount, Service.Update, err:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	discount := entity.Discount{
		PublishStatusID: 2,
	}
	publishStatusDesc := discount.GetPublishStatusDesc()
	responsePayload.Setmsg(publishStatusDesc + " Successfully")
	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}

func (controller *DiscountController) Consult(c *fiber.Ctx) error {
	var (
		request entity.ConsultDiscountBody
	)
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	if err := c.BodyParser(&request); err != nil {
		log.Error("DiscountController, Consult, BodyParser(request), err:", err.Error())
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	request.CustID = c.Locals("cust_id").(string)
	request.ParentCustID = c.Locals("parent_cust_id").(string)

	if errs := controller.validator.ValidateStruct(request, headerAcceptLang); errs != nil {
		log.Error("DiscountController, Consult, ValidateStruct(request), errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	responses, err := controller.DiscountService.ConsultDiscount(request)
	if err != nil {
		log.Error("DiscountController, ConsultDiscount, Service.ConsultDiscount, err:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	responsePayload.Setdata(responses)
	responsePayload.Setmsg("Discount Consulted Successfully")
	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}
