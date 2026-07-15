package controller

import (
	"system/entity"
	"system/pkg/constant"
	"system/pkg/middleware"
	"system/pkg/responsebuild"
	"system/pkg/validation"
	"system/service"

	"github.com/gofiber/fiber/v2/log"

	"github.com/gofiber/fiber/v2"
)

type UserController struct {
	UserService service.UserService
	validator   *validation.Validate
}

func NewUserController(userService service.UserService, validator *validation.Validate) *UserController {
	return &UserController{
		UserService: userService,
		validator:   validator,
	}
}

func (controller *UserController) Route(app *fiber.App) {
	qParamId := ":user_id"
	qParamMenuID := ":menu_param"
	app.Post("v1/users/login", controller.Login)
	app.Post("v1/users/forgot-password", controller.ForgotPassword)
	app.Post("v1/users/forgot-password-validate", controller.ForgotPasswordValidate)
	app.Patch("v1/users/password", controller.UpdatePassword)
	userMenuRouteV1 := app.Group("/v1/users-menu")
	userMenuRouteV1.Get("/menus/"+qParamMenuID, controller.UserMenusAll)
	userRouteV1 := app.Group("/v1/users", middleware.JWTProtected())
	custRouteV1 := app.Group("/v1/cust", middleware.JWTProtected())
	custRouteV1.Get("", controller.DetailCust)
	userRouteV1.Post("", controller.Create)
	userRouteV1.Patch("/"+qParamId, controller.Update)
	userRouteV1.Get("/menus/"+qParamMenuID, controller.UserMenus)
	userRouteV1.Get("/"+qParamId, controller.Detail)
	userRouteV1.Delete("/"+qParamId, controller.Delete)
	userRouteV1.Get("", controller.List)

}

func (controller *UserController) Login(c *fiber.Ctx) error {
	var request entity.LoginRequest

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
	errs := controller.validator.ValidateStruct(request, headerAcceptLang)
	if errs != nil {
		log.Error(errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	data, err := controller.UserService.Login(request)
	if err != nil {
		responsePayload.Setmsg("Login failed")
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}
	responsePayload.Setdata(data)
	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}

func (controller *UserController) ForgotPassword(c *fiber.Ctx) error {
	var request entity.ForgotPasswordRequest

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
	errs := controller.validator.ValidateStruct(request, headerAcceptLang)
	if errs != nil {
		log.Error(errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	err = controller.UserService.ForgotPassword(request)
	if err != nil {
		responsePayload.Setmsg("forgot password failed")
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}
	// log.Println("response:", response)
	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}

func (controller *UserController) ForgotPasswordValidate(c *fiber.Ctx) error {
	var request entity.ForgotPasswordValidateRequest

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
	errs := controller.validator.ValidateStruct(request, headerAcceptLang)
	if errs != nil {
		log.Error(errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	err = controller.UserService.ForgotPasswordValidate(request)
	if err != nil {
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}
	// log.Println("response:", response)
	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}

func (controller *UserController) UpdatePassword(c *fiber.Ctx) error {
	var request entity.UpdatePasswordRequest

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
	errs := controller.validator.ValidateStruct(request, headerAcceptLang)
	if errs != nil {
		log.Error(errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	err = controller.UserService.UpdatePassword(request)
	if err != nil {
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}
	// log.Println("response:", response)
	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}

func (controller *UserController) Create(c *fiber.Ctx) error {

	var request entity.CreateUserBody
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	if err := c.BodyParser(&request); err != nil {
		log.Error(err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	custId := c.Locals("cust_id").(string)
	userId := c.Locals("user_id").(int64)

	request.CustId = custId
	request.CreatedBy = userId
	errs := controller.validator.ValidateStruct(request, headerAcceptLang)
	if errs != nil {
		log.Error(errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	err := controller.UserService.Store(request)
	if err != nil {
		log.Error(err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}
	responsePayload.Setmsg("Created Successfully")
	return c.Status(fiber.StatusCreated).JSON(responsePayload.GetRespPayload())
}

func (controller *UserController) Update(c *fiber.Ctx) error {
	var (
		params  entity.UpdateUserBodyParam
		request entity.UpdateUserBody
	)
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	if err := c.ParamsParser(&params); err != nil {
		log.Error(err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	errs := controller.validator.ValidateStruct(params, headerAcceptLang)
	if errs != nil {
		log.Error(errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(fiber.ErrBadRequest.Message)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	if err := c.BodyParser(&request); err != nil {
		log.Error(err.Error())
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	custId := c.Locals("cust_id").(string)
	userId := c.Locals("user_id").(int64)
	// log.Println("BankController, Update, CustId:", custId)
	request.CustId = custId
	request.UpdatedBy = userId

	errs = controller.validator.ValidateStruct(request, headerAcceptLang)
	if errs != nil {
		log.Error(errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	err := controller.UserService.Update(params.UserID, request)
	if err != nil {
		log.Error(err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}
	responsePayload.Setmsg("Updated Successfully")
	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}

func (controller *UserController) Detail(c *fiber.Ctx) error {
	var params entity.DetailUserBodyParam
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	if err := c.ParamsParser(&params); err != nil {
		log.Error(err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	errs := controller.validator.ValidateStruct(params, headerAcceptLang)
	if errs != nil {
		log.Error(errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	custId := c.Locals("cust_id").(string)
	data, err := controller.UserService.Detail(params.UserID, custId)
	if err != nil {
		log.Error(err.Error())
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

func (controller *UserController) DetailCust(c *fiber.Ctx) error {
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	custId := c.Locals("cust_id").(string)

	data, err := controller.UserService.DetailCust(custId)
	if err != nil {
		log.Error(err.Error())
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

func (controller *UserController) Delete(c *fiber.Ctx) error {
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)
	var params entity.DeleteUserBodyParams
	if err := c.ParamsParser(&params); err != nil {
		log.Error(err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}
	errs := controller.validator.ValidateStruct(params, headerAcceptLang)
	if errs != nil {
		log.Error(errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	custId := c.Locals("cust_id").(string)
	userId := c.Locals("user_id").(int64)
	// log.Println("VehicleController, Delete, CustId:", custId)

	err := controller.UserService.Delete(custId, int(params.UserID), userId)
	if err != nil {
		log.Error(err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	responsePayload.Setmsg("Deleted Successfully")
	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}

func (controller *UserController) List(c *fiber.Ctx) error {
	var dataFilter entity.GeneralQueryFilter
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)
	if err := c.QueryParser(&dataFilter); err != nil {
		log.Error(err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	custId := c.Locals("cust_id").(string)
	data, total, lastPage, err := controller.UserService.List(dataFilter, custId)
	if err != nil {
		log.Error(err.Error())
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

func (controller *UserController) UserMenus(c *fiber.Ctx) error {
	var (
		params entity.MenuBodyParam
	)
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string))
	if err := c.ParamsParser(&params); err != nil {
		log.Error(err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}
	userId := c.Locals("user_id").(int64)
	custId := c.Locals("cust_id").(string)
	if params.MenuParam == entity.MENUWEB {
		data, err := controller.UserService.UserMenus(userId, custId)
		if err != nil {
			log.Error(err.Error())
			responsePayload.Setmsg(err.Error())
			return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
		}
		responsePayload.Setdata(data)
		return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
	} else {
		data, err := controller.UserService.UserMenusDesktop(userId)
		if err != nil {
			log.Error(err.Error())
			responsePayload.Setmsg(err.Error())
			return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
		}
		responsePayload.Setdata(data)
		return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
	}

}

func (controller *UserController) UserMenusAll(c *fiber.Ctx) error {
	var (
		params entity.MenuBodyParam
	)
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string))
	if err := c.ParamsParser(&params); err != nil {
		log.Error(err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}
	// userId := c.Locals("user_id").(int64)
	// custId := c.Locals("cust_id").(string)
		data, err := controller.UserService.UserMenusAll(params.MenuParam)
		if err != nil {
			log.Error(err.Error())
			responsePayload.Setmsg(err.Error())
			return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
		}
		responsePayload.Setdata(data)
		return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())

}
