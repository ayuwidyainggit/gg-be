package controller

import (
	"fmt"
	"mobile/entity"
	"mobile/pkg/constant"
	"mobile/pkg/middleware"
	"mobile/pkg/responsebuild"
	"mobile/pkg/validation"
	"mobile/service"

	"github.com/gofiber/fiber/v2/log"

	"github.com/gofiber/fiber/v2"
)

type ProductController struct {
	ProductService service.ProductService
	validator      *validation.Validate
}

func NewProductController(
	productService service.ProductService,
	validator *validation.Validate,
) *ProductController {
	return &ProductController{
		ProductService: productService,
		validator:      validator,
	}
}

func (controller *ProductController) Route(app *fiber.App) {
	qParamId := ":pro_id"
	productRouteV1 := app.Group("/v1/products", middleware.JWTProtected())
	productRouteV1.Get("/"+qParamId, controller.Detail)
	productRouteV1.Get("", controller.List)

}

func (controller *ProductController) List(c *fiber.Ctx) error {
	var (
		headerAcceptLang string
		dataFilter       entity.ProductsQueryFilter
	)
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	if err := c.QueryParser(&dataFilter); err != nil {
		log.Error("ProductController, List, query parser filter:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}
	errs := controller.validator.ValidateStruct(dataFilter, headerAcceptLang)
	if errs != nil {
		log.Error("ProductController, List, ValidateStruct(params), errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	custId := c.Locals("cust_id").(string)
	parentCustId := c.Locals("parent_cust_id").(string)
	EmpId := c.Locals("emp_id").(int64)
	IsActiveGudangUtama := c.Locals("is_active_gudang_utama").(bool)
	IsActiveGudangCanvas := c.Locals("is_active_gudang_canvas").(bool)
	fmt.Print(IsActiveGudangCanvas, "====", IsActiveGudangUtama)
	// fmt.Println(">>>", EmpId)

	// log.Println("custId:", custId)
	// log.Println("parentCustId:", parentCustId)

	data, total, lastPage, err := controller.ProductService.List(dataFilter, custId, parentCustId, EmpId, IsActiveGudangCanvas, IsActiveGudangUtama)
	if err != nil {
		log.Error("ProductController, List, data, err:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	responsePayload.Setmsg(constant.STATUS_OK)
	responsePayload.Setdata(data)
	responsePayload.Setpaging(entity.Pagination{
		TotalRecord: total,
		PageCurrent: dataFilter.Page,
		PageLimit:   dataFilter.Limit,
		PageTotal:   lastPage,
	})

	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}

func (controller *ProductController) Detail(c *fiber.Ctx) error {
	var (
		params entity.DetailProductParams
	)
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), constant.HEADER_ACCEPT_LANG)

	if err := c.ParamsParser(&params); err != nil {
		log.Error("ProductController, Detail, ParamsParser:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	errs := controller.validator.ValidateStruct(params, constant.HEADER_ACCEPT_LANG)
	if errs != nil {
		log.Error("ProductController, Detail, ValidateStruct(params), errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	params.CustID = c.Locals("cust_id").(string)
	params.ParentCustID = c.Locals("parent_cust_id").(string)
	// log.Println("ProductController, Detail, CustId:", custId)

	data, err := controller.ProductService.Detail(params)
	if err != nil {
		log.Error("ProductController, Detail, FindOneByProductId, err:", err.Error())
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
