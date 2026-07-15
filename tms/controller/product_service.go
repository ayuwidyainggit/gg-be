package controller

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"net/http"
	"scyllax-tms/entity"
	"scyllax-tms/exception"
	"scyllax-tms/service"
	"scyllax-tms/utils"
	"time"
)

type ProductController struct {
	productService service.ProductService
}

func NewProductController(productService service.ProductService) *ProductController {
	return &ProductController{
		productService: productService,
	}
}

// Note 		        godoc
//
//	@Summary		Get all product.
//	@Description	Return the product.
//	@Param			shipment_no		query	string	false	"shipment_no"
//	@Param			outlet_id		query	string	false	"outlet_id"
//	@Param			outlet_code		query	string	false	"outlet_code"
//	@Param			cust_id			query	string	false	"cust_id"
//	@Param			product_id		query	string	false	"product_id"
//	@Param			product_name	query	string	false	"product_name"
//	@Param			driver_id		query	string	false	"driver_id"
//	@Param			status		    query	string	false	"status"
//	@Produce		application/json
//	@Tags			product
//	@Success		200	{object}	entity.JsonSuccess{data=[]entity.ProductResponse}	"Data"
//	@Failure		400	{object}	entity.JsonBadRequest{}								"Validation error"
//	@Failure		404	{object}	entity.JsonNotFound{}								"Data not found"
//	@Failure		500	{object}	entity.JsonInternalServerError{}					"Internal server error"
//	@Router			/mobile/products [get]
func (controller *ProductController) GetProduct(ctx *fiber.Ctx) error {
	c, cancel := context.WithTimeout(ctx.Context(), 30*time.Second)
	defer cancel()

	var dataFilter entity.ShipmentInvoicesQueryFilter

	if err := ctx.QueryParser(&dataFilter); err != nil {
		panic(exception.NewBadRequestError(err.Error()))
	}

	response := controller.productService.GetProduct(c, dataFilter)

	webResponse := entity.Response{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   response,
	}
	utils.ResponseInterceptor(c, &webResponse)
	return ctx.Status(fiber.StatusOK).JSON(webResponse)
}
