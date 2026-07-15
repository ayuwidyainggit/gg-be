package controller

import (
	"finance/entity"
	"finance/pkg/constant"
	"finance/pkg/responsebuild"
	"finance/pkg/validation"
	"finance/service"

	"github.com/gofiber/fiber/v2/log"

	"github.com/gofiber/fiber/v2"
)

type FilesController struct {
	FilesService service.FilesService
	validator    *validation.Validate
}

func NewFilesController(
	FilesService service.FilesService,
	validator *validation.Validate,
) *FilesController {
	return &FilesController{
		FilesService: FilesService,
		validator:    validator,
	}
}

func (controller *FilesController) Route(app *fiber.App) {
	FilesRouteV1 := app.Group("/v1/files")
	FilesRouteV1.Post("/uploads", controller.Uploads)

}
func (controller *FilesController) Uploads(c *fiber.Ctx) error {
	var (
		request entity.UploadRequest
	)
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	err := c.BodyParser(&request)
	if err != nil {
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	file, err := c.FormFile("file")
	if err != nil {
		responsePayload.Setmsg("upload file on form with key 'file'")
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	errs := controller.validator.ValidateStruct(request, headerAcceptLang)
	if errs != nil {
		log.Error("FileController, Create, ValidateStruct, errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	request.File = file

	data, err := controller.FilesService.Upload(request)
	if err != nil {
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}
	responsePayload.Data = data
	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}
