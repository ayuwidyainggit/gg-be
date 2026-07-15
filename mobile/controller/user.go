package controller

import (
	"encoding/json"
	"errors"
	"mobile/entity"
	"mobile/pkg/apperr"
	"mobile/pkg/constant"
	"mobile/pkg/middleware"
	"mobile/pkg/responsebuild"
	"mobile/pkg/validation"
	"mobile/service"

	"github.com/gofiber/fiber/v2/log"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type UserController struct {
	UserService service.UserService
	validator   *validation.Validate
}

func NewUserController(
	userService service.UserService,
	validator *validation.Validate,
) *UserController {
	return &UserController{
		UserService: userService,
		validator:   validator,
	}
}

func (controller *UserController) Route(app *fiber.App) {
	qParamId := ":user_id"
	app.Post("v1/users/register", controller.ForgotPassword)
	app.Post("v1/users/forgot-password", controller.ForgotPassword)
	app.Post("v1/users/forgot-password-validate", controller.ForgotPasswordValidate)
	app.Post("v1/users/forgot-password-resend", controller.ForgotPasswordResend)
	app.Post("v1/users/reset-password", controller.ResetPassword)
	app.Patch("v1/users/password", controller.UpdatePassword)
	app.Patch("v1/users/change-password", controller.ChangePassword)

	app.Post("v1/users/login", controller.Login)
	userRouteV1 := app.Group("/v1/users", middleware.JWTProtected())
	userRouteV1.Patch("/"+qParamId, controller.Update)
	userRouteV1.Get("/profile", controller.Profile)

	app.Post("v1/send_location", middleware.JWTProtected(), controller.SendLocation)
}

func (controller *UserController) Register(c *fiber.Ctx) error {
	var (
		headerAcceptLang string
		request          entity.RegisterRequest
	)

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
		// log.Println("errs, UserController-Login:", structs.StructToJson(errs))
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	// find employee with type salesman, update m_user set password
	// if err = controller.UserService.ForgotPassword(request); err != nil {
	// 	isErrRecordNotFound := errors.Is(err, gorm.ErrRecordNotFound)
	// 	if isErrRecordNotFound {
	// 		responsePayload.Setmsg("Email not registered")
	// 	}
	// 	return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	// }

	// log.Println("response:", response)
	responsePayload.Setmsg(constant.STATUS_OK)
	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}

func (controller *UserController) ForgotPassword(c *fiber.Ctx) error {
	var (
		headerAcceptLang string
		request          entity.ForgotPasswordRequest
	)

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
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	requestID, err := controller.UserService.ForgotPassword(request)
	if err != nil {
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	responsePayload.Setmsg("OTP sent")
	responsePayload.Setdata(map[string]string{"request_id": requestID})
	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}

func (controller *UserController) ForgotPasswordValidate(c *fiber.Ctx) error {
	var (
		headerAcceptLang string
		request          entity.ForgotPasswordValidateRequest
	)

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
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	resetToken, err := controller.UserService.ForgotPasswordValidate(request)
	if err != nil {
		ae, ok := err.(*apperr.AppError)
		message := err.Error()
		if ok {
			message = ae.Msg
		}
		responsePayload.Setmsg(message)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	responsePayload.Setmsg("OTP valid")
	responsePayload.Setdata(map[string]string{"reset_token": resetToken})
	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}

func (controller *UserController) ForgotPasswordResend(c *fiber.Ctx) error {
	var (
		headerAcceptLang string
		request          entity.ForgotPasswordResendRequest
	)

	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	if err := c.BodyParser(&request); err != nil {
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}
	if errs := controller.validator.ValidateStruct(request, headerAcceptLang); errs != nil {
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	if err := controller.UserService.ForgotPasswordResend(request); err != nil {
		ae, ok := err.(*apperr.AppError)
		message := err.Error()
		if ok {
			message = ae.Msg
		}
		responsePayload.Setmsg(message)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	responsePayload.Setmsg("OTP resent")
	responsePayload.Setdata(json.RawMessage("null"))
	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}

func (controller *UserController) ResetPassword(c *fiber.Ctx) error {
	var (
		headerAcceptLang string
		request          entity.ResetPasswordRequest
	)

	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	if err := c.BodyParser(&request); err != nil {
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}
	if errs := controller.validator.ValidateStruct(request, headerAcceptLang); errs != nil {
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	if err := controller.UserService.ResetPassword(request); err != nil {
		ae, ok := err.(*apperr.AppError)
		message := err.Error()
		if ok {
			message = ae.Msg
		}
		responsePayload.Setmsg(message)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	responsePayload.Setmsg("Password successfully updated")
	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}

func (controller *UserController) Login(c *fiber.Ctx) error {
	var (
		headerAcceptLang string
		request          entity.LoginRequest
	)

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
		// log.Println("errs, UserController-Login:", structs.StructToJson(errs))
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	data, err := controller.UserService.Login(request)
	if err != nil {
		responsePayload.Setmsg("Login failed")
		isErrRecordNotFound := errors.Is(err, gorm.ErrRecordNotFound)
		if isErrRecordNotFound {
			errs = append(errs, map[string]interface{}{
				"key":     "email",
				"message": "User not found",
			})
			responsePayload.Seterrors(errs)
		}
		if err.Error() == "Invalid password" {
			errs = append(errs, map[string]interface{}{
				"key":     "password",
				"message": err.Error(),
			})
			responsePayload.Seterrors(errs)
		}
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}
	responsePayload.Setdata(data)
	// log.Println("response:", response)
	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}

func (controller *UserController) UpdatePassword(c *fiber.Ctx) error {
	var (
		headerAcceptLang string
		request          entity.UpdatePasswordRequest
	)

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
		// log.Println("errs, UserController-Login:", structs.StructToJson(errs))
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	if request.Password != request.PasswordConfirm {
		responsePayload.Setmsg("Password does not match")
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}
	err = controller.UserService.UpdatePassword(request)
	if err != nil {
		responsePayload.Setmsg("update failed")
		isErrRecordNotFound := errors.Is(err, gorm.ErrRecordNotFound)
		if isErrRecordNotFound {
			errs = append(errs, map[string]interface{}{
				"key":     "email",
				"message": "User not found",
			})
			responsePayload.Seterrors(errs)
		} else {
			responsePayload.Setmsg(err.Error())
		}

		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}
	// log.Println("response:", response)
	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}

func (controller *UserController) ChangePassword(c *fiber.Ctx) error {
	var (
		headerAcceptLang string
		request          entity.ChangePasswordRequest
	)

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
		// log.Println("errs, UserController-Login:", structs.StructToJson(errs))
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	if request.NewPassword != request.PasswordConfirm {
		responsePayload.Setmsg("Password does not match")
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	err = controller.UserService.ChangePassword(request)
	if err != nil {
		responsePayload.Setmsg("update failed")
		isErrRecordNotFound := errors.Is(err, gorm.ErrRecordNotFound)
		if isErrRecordNotFound {
			errs = append(errs, map[string]interface{}{
				"key":     "email",
				"message": "User not found",
			})
			responsePayload.Seterrors(errs)
		} else {
			responsePayload.Setmsg(err.Error())
		}

		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}
	// log.Println("response:", response)
	responsePayload.Setmsg("Updated Password Successfully")
	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}

func (controller *UserController) Profile(c *fiber.Ctx) error {
	var (
		headerAcceptLang string
	)
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	email := c.Locals("email").(string)

	data, err := controller.UserService.GetProfile(email)
	if err != nil {
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}
	responsePayload.Setdata(data)
	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}

func (controller *UserController) Update(c *fiber.Ctx) error {
	var (
		params           entity.UpdateUserBodyParam
		request          entity.UpdateUserImageBody
		headerAcceptLang string
	)
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	if err := c.ParamsParser(&params); err != nil {
		log.Error("UserController, Update, ParamsParser(params):", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	if err := c.BodyParser(&request); err != nil {
		log.Error("UserController, Update, BodyParser(request), err:", err.Error())
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	errs := controller.validator.ValidateStruct(request, headerAcceptLang)
	if errs != nil {
		log.Error("UserController, Update, ValidateStruct(params), errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}
	custId := c.Locals("cust_id").(string)
	// userId := c.Locals("user_id").(int64)
	// log.Println("UserController, Update, CustId:", custId)

	request.CustId = custId
	// request.UpdatedBy = userId

	errs = controller.validator.ValidateStruct(request, headerAcceptLang)

	if errs != nil {
		log.Error("UserController, Update, ValidateStruct(request), errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	err := controller.UserService.Update(params.UserID, request)
	if err != nil {
		log.Error("UserController, Update, Service.Update, err:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	responsePayload.Setmsg(constant.SUCCESSFULLY_UPDATED)
	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())

}

func (controller *UserController) SendLocation(c *fiber.Ctx) error {
	var (
		headerAcceptLang string
		request          entity.SendLocationRequest
	)

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
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}
	request.CustID = c.Locals("cust_id").(string)
	request.EmpID = c.Locals("emp_id").(int64)

	err = controller.UserService.SendLocation(request)
	if err != nil {
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}
	responsePayload.Setmsg(constant.SUCCESSFULLY_SAVED)
	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}
