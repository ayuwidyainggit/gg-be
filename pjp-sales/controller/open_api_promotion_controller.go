package controller

import (
	"fmt"
	"strings"

	"sales/entity"
	"sales/pkg/constant"
	"sales/pkg/middleware"
	"sales/pkg/responsebuild"
	"sales/pkg/validation"
	"sales/service"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

type OpenAPIPromotionController struct {
	PromotionController *PromotionController
	OpenAPIService      service.OpenAPIService
	validator           *validation.Validate
}

func NewOpenAPIPromotionController(
	promotionController *PromotionController,
	openAPIService service.OpenAPIService,
	validator *validation.Validate,
) *OpenAPIPromotionController {
	return &OpenAPIPromotionController{
		PromotionController: promotionController,
		OpenAPIService:      openAPIService,
		validator:           validator,
	}
}

func (controller *OpenAPIPromotionController) Route(app *fiber.App) {
	group := app.Group(
		"/open-api/v1",
		middleware.OpenAPIProtected(controller.OpenAPIService),
	)
	group.Post("/promotions", controller.CreatePromotion)
}

func (controller *OpenAPIPromotionController) CreatePromotion(c *fiber.Ctx) error {
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	var request entity.CreatePromotionV2Body
	if err := c.BodyParser(&request); err != nil {
		log.Error("OpenAPIPromotionController, CreatePromotion, BodyParser:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	parentCustID := fmt.Sprint(c.Locals("parent_cust_id"))
	request.CustID = parentCustID
	request.ParentCustID = parentCustID

	if strings.TrimSpace(request.DistributorCustID) == "" {
		request.DistributorCustID = parentCustID
	}

	request.CreatedBy = fmt.Sprint(c.Locals("user_fullname"))
	request.UpdatedBy = request.CreatedBy

	if request.PromoStatus == "" {
		request.PromoStatus = entity.PromoStatusDraft
	}
	if request.Coverage == "" {
		request.Coverage = entity.CoverageNational
	}

	if request.PromoCreationType == entity.CreationTypeNew {
		request.ExistingPromoID = ""
	}

	return controller.PromotionController.processCreatePromotionV2(
		c, request, headerAcceptLang, constant.MsgOpenAPICreatePromotionSuccess,
	)
}
