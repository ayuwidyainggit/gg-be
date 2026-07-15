package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"regexp"
	"time"

	"mobile/entity"
	"mobile/model"
	"mobile/pkg/apperr"
	"mobile/pkg/config/env"
	"mobile/pkg/constant"
	"mobile/pkg/jwthelper"
	"mobile/pkg/mail/smtp"
	"mobile/pkg/passwordreset"
	"mobile/pkg/str"
	"mobile/pkg/structs"
	"mobile/repository"

	"github.com/rs/xid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserService interface {
	ForgotPassword(entity.ForgotPasswordRequest) (requestID string, err error)
	Login(entity.LoginRequest) (entity.LoginResponse, error)
	ForgotPasswordValidate(entity.ForgotPasswordValidateRequest) (resetToken string, err error)
	ForgotPasswordResend(entity.ForgotPasswordResendRequest) error
	ResetPassword(entity.ResetPasswordRequest) error
	UpdatePassword(request entity.UpdatePasswordRequest) error
	ChangePassword(request entity.ChangePasswordRequest) error
	GetProfile(email string) (response entity.ProfileResponse, err error)
	Update(userId int64, request entity.UpdateUserImageBody) error
	SendLocation(request entity.SendLocationRequest) error
}

type userServiceImpl struct {
	Config                  env.ConfigEnv
	MCustomerRepository     repository.MCustomerRepository
	MEmployeeRepository     repository.MEmployeeRepository
	MSalesmanRepository     repository.MSalesmanRepository
	UserRepository          repository.UserRepository
	UserLocationRepository  repository.UserLocationRepository
	PasswordResetRepository repository.PasswordResetRequestRepository
	Transaction             repository.Dbtransaction
}

func NewUserService(
	config env.ConfigEnv,
	mCustomer repository.MCustomerRepository,
	mEmployee repository.MEmployeeRepository,
	mSalesman repository.MSalesmanRepository,
	userRepository repository.UserRepository,
	userLocationRepository repository.UserLocationRepository,
	passwordResetRepository repository.PasswordResetRequestRepository,
	transaction repository.Dbtransaction,
) *userServiceImpl {
	return &userServiceImpl{
		Config:                  config,
		MCustomerRepository:     mCustomer,
		MEmployeeRepository:     mEmployee,
		MSalesmanRepository:     mSalesman,
		UserRepository:          userRepository,
		UserLocationRepository:  userLocationRepository,
		PasswordResetRepository: passwordResetRepository,
		Transaction:             transaction,
	}
}

func (service *userServiceImpl) ForgotPassword(request entity.ForgotPasswordRequest) (requestID string, err error) {
	ctx := context.Background()

	employee, err := service.MEmployeeRepository.FindOneByEmail(request.Email)
	if err != nil {
		return "", mapForgotUserLookupErr(err)
	}

	_, err = service.MSalesmanRepository.FindOneByEmpId(employee.CustId, employee.EmpId)
	if err != nil {
		return "", mapForgotUserLookupErr(err)
	}

	user, err := service.UserRepository.FindOneByEmailCustIdEmpId(request.Email, employee.CustId, employee.EmpId)
	if err != nil {
		return "", mapForgotUserLookupErr(err)
	}

	if err := service.PasswordResetRepository.ExpirePendingByUserID(ctx, *user.UserId); err != nil {
		return "", err
	}

	otpCode, err := str.Generate(`[0-9]{4}`)
	if err != nil {
		return "", err
	}

	rid := xid.New().String()
	now := time.Now().UTC()
	cooldown := now.Add(30 * time.Second)
	exp := now.Add(5 * time.Minute)

	row := &model.PasswordResetRequest{
		UserID:              *user.UserId,
		Email:               request.Email,
		OtpCode:             otpCode,
		OtpExpiredAt:        exp,
		OtpAttemptCount:     0,
		OtpMaxAttempt:       3,
		ResendCount:         0,
		ResendMax:           3,
		ResendCooldownUntil: &cooldown,
		RequestID:           rid,
		Status:              constant.PasswordResetStatusPendingOTP,
		CreatedAt:           now,
		UpdatedAt:           now,
	}
	if err := service.PasswordResetRepository.Create(ctx, row); err != nil {
		return "", err
	}

	service.sendForgotPasswordOTPEmailAsync(rid, request.Email, otpCode)

	return rid, nil
}

// sendForgotPasswordOTPEmailAsync sends OTP email in the background so HTTP handlers are not blocked by SMTP.
func (service *userServiceImpl) sendForgotPasswordOTPEmailAsync(requestID, toEmail, otpCode string) {
	cfg := service.Config
	email := toEmail
	code := otpCode
	rid := requestID
	go func() {
		data := entity.EmailData{Email: email, OtpCode: code}
		err := smtp.BaseSendEmailWithTemplate(cfg, email, "Forgot Password", "forgot_password_otp.html", data)
		if err != nil {
			log.Printf("[ForgotPassword] email send FAILED request_id=%s to=%s err=%v", rid, email, err)
			return
		}
		log.Printf("[ForgotPassword] email sent OK request_id=%s to=%s", rid, email)
	}()
}

func (service *userServiceImpl) Login(request entity.LoginRequest) (response entity.LoginResponse, err error) {
	employee, err := service.MEmployeeRepository.FindOneByEmail(request.Email)
	if err != nil {
		return response, err
	}

	if err := structs.Automapper(employee, &response); err != nil {
		return response, errors.New("parse error employee data")
	}
	// if employee.DeviceID != nil {
	// 	if request.DeviceId != *employee.DeviceID {
	// 		return response, errors.New("Device ID not match")
	// 	}
	// } else {
	// 	return response, errors.New("Please Register Device ID first")
	// }

	// if employee.MacAddress != nil {
	// 	if request.MacAddress != *employee.MacAddress {
	// 		return response, errors.New("Mac Address not match")
	// 	}
	// } else {
	// 	return response, errors.New("Please Register mac address first")

	// }

	var salesman model.MSalesmanRead
	salesman, err = service.MSalesmanRepository.FindOneByEmpId(employee.CustId, employee.EmpId)
	if err != nil {
		return response, err
	}
	if err := structs.Automapper(salesman, &response); err != nil {
		return response, errors.New("parse error employee data")
	}

	customer, err := service.MCustomerRepository.FindOneByCustId(employee.CustId)
	if err != nil {
		return response, err
	}

	// Deprecated
	// Get distributor data
	// distributor := model.DistributorDetail{}
	// distributorData, err := service.MEmployeeRepository.FindOneByDistributor(customer.DistributorId)
	// if err != nil {
	// 	log.Printf("Login: FindOneByDistributor error for distributor_id %d: %v", customer.DistributorId, err)
	// 	// Set default empty values if distributor not found
	// 	distributor.DistributorID = 0
	// 	distributor.DistributorCode = ""
	// 	distributor.DistributorName = ""
	// 	distributor.Address = ""
	// } else {
	// 	distributor = distributorData
	// }

	// user, err := service.UserRepository.FindOneByEmailAndCustId(request.Email, employee.CustId)
	user, err := service.UserRepository.FindOneByEmailCustIdEmpId(request.Email, employee.CustId, employee.EmpId)
	if err != nil {
		return response, err
	}

	if !service.CheckPasswordHash(request.Password, *user.Userpass) {
		return response, errors.New("Invalid password")
	}

	credentials := make([]string, 0)
	additionalParam := make(map[string]interface{}, 0)

	userData := entity.UserData{}
	if err = structs.Automapper(user, &userData); err != nil {
		return response, errors.New("error automapper")
	}
	userData.EmpCode = &employee.EmpCode
	userData.EmpGrpId = employee.EmpGrpId
	userData.ParentCustId = customer.ParentCustID
	userData.DistPriceGrpId = customer.DistPriceGrpId
	userData.DistributorID = customer.DistributorId
	response.UserId = *user.UserId
	if user.RoleName != nil {
		userData.UserRole = *user.RoleName
	}
	userData.OprTypeOrderTaking = salesman.OprTypeOrderTaking
	userData.OprTypeCanvas = salesman.OprTypeCanvas
	userData.AllowInputPrice = salesman.AllowInputPrice
	userData.TaxOption = salesman.TaxOption
	userData.IsActiveGudangCanvas = salesman.IsActiveGudangCanvas
	userData.IsActiveGudangUtama = salesman.IsActiveGudangUtama
	userData.Username = *user.Username
	userData.UserFullname = *user.Fullname
	userData.IsAdmin = *user.IsAdmin

	response.AccessToken, err = jwthelper.GenerateNewToken(userData, credentials, additionalParam)
	if err != nil {
		return response, err
	}

	err = service.UserRepository.UpdateFcmToken(context.Background(), *user.UserId, request.FcmToken)
	if err != nil {
		return response, err
	}

	// endhance distributor
	response.DistributorName = customer.DistributorName
	response.DistributorCode = customer.DistributorCode
	response.DistributorAddress = customer.DistributorAddress

	err = structs.Automapper(userData, &response)
	if err != nil {
		return response, errors.New("error mapping user data")
	}

	return response, err
}

func (service *userServiceImpl) ForgotPasswordValidate(request entity.ForgotPasswordValidateRequest) (resetToken string, err error) {
	ctx := context.Background()

	row, err := service.PasswordResetRepository.FindByRequestID(ctx, request.RequestID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", apperr.New(constant.ErrCodeOTPInvalidOrExpired, "Invalid OTP")
		}
		return "", err
	}

	if row.Status != constant.PasswordResetStatusPendingOTP {
		return "", apperr.New(constant.ErrCodeOTPInvalidOrExpired, "Invalid OTP")
	}
	if time.Now().UTC().After(row.OtpExpiredAt) {
		_ = service.PasswordResetRepository.Update(ctx, row.ID, map[string]interface{}{
			"status": constant.PasswordResetStatusExpired,
		})
		return "", apperr.New(constant.ErrCodeOTPInvalidOrExpired, "Invalid OTP")
	}
	if row.OtpAttemptCount >= row.OtpMaxAttempt {
		return "", apperr.New(constant.ErrCodeOTPMaxAttempts, "Invalid OTP")
	}

	if row.OtpCode != request.OtpCode {
		next := row.OtpAttemptCount + 1
		updates := map[string]interface{}{"otp_attempt_count": next}
		if next >= row.OtpMaxAttempt {
			updates["status"] = constant.PasswordResetStatusLocked
		}
		if uerr := service.PasswordResetRepository.Update(ctx, row.ID, updates); uerr != nil {
			return "", uerr
		}
		return "", apperr.New(constant.ErrCodeOTPInvalidOrExpired, "Invalid OTP")
	}

	plainToken, err := randomHexToken(32)
	if err != nil {
		return "", err
	}
	h := passwordreset.ResetTokenHMACHex(service.resetTokenSecret(), plainToken)
	exp := time.Now().UTC().Add(60 * time.Minute)
	if err := service.PasswordResetRepository.Update(ctx, row.ID, map[string]interface{}{
		"reset_token":            h,
		"reset_token_expired_at": exp,
		"status":                 constant.PasswordResetStatusOTPValidated,
	}); err != nil {
		return "", err
	}

	return plainToken, nil
}

func (service *userServiceImpl) ForgotPasswordResend(request entity.ForgotPasswordResendRequest) error {
	ctx := context.Background()

	row, err := service.PasswordResetRepository.FindByRequestID(ctx, request.RequestID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperr.New(constant.ErrCodeOTPInvalidOrExpired, "Invalid OTP")
		}
		return err
	}
	if row.Status != constant.PasswordResetStatusPendingOTP {
		return apperr.New(constant.ErrCodeOTPInvalidOrExpired, "Invalid OTP")
	}
	if row.ResendCount >= row.ResendMax {
		return apperr.New(constant.ErrCodeOTPMaxResend, "OTP resent limit reached")
	}
	if row.ResendCooldownUntil != nil && time.Now().UTC().Before(*row.ResendCooldownUntil) {
		return apperr.New(constant.ErrCodeOTPCooldown, "Please wait before resending OTP")
	}

	otpCode, err := str.Generate(`[0-9]{4}`)
	if err != nil {
		return err
	}
	now := time.Now().UTC()
	cooldown := now.Add(30 * time.Second)
	exp := now.Add(5 * time.Minute)

	if err := service.PasswordResetRepository.Update(ctx, row.ID, map[string]interface{}{
		"otp_code":               otpCode,
		"otp_expired_at":         exp,
		"otp_attempt_count":      0,
		"resend_count":           row.ResendCount + 1,
		"resend_cooldown_until":  cooldown,
		"reset_token":            nil,
		"reset_token_expired_at": nil,
	}); err != nil {
		return err
	}

	service.sendForgotPasswordOTPEmailAsync(row.RequestID, row.Email, otpCode)
	return nil
}

func (service *userServiceImpl) ResetPassword(request entity.ResetPasswordRequest) error {
	ctx := context.Background()

	if request.NewPassword != request.ConfirmPassword {
		return apperr.New(constant.ErrCodePasswordConfirmMismatch, "Password does not match")
	}
	if err := validatePasswordPolicy(request.NewPassword); err != nil {
		return err
	}

	h := passwordreset.ResetTokenHMACHex(service.resetTokenSecret(), request.ResetToken)
	row, err := service.PasswordResetRepository.FindByResetToken(ctx, h)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperr.New(constant.ErrCodeResetTokenInvalidOrExpired, "Invalid reset token")
		}
		return err
	}
	if row.Status != constant.PasswordResetStatusOTPValidated {
		return apperr.New(constant.ErrCodeResetTokenInvalidOrExpired, "Invalid reset token")
	}
	if row.ResetTokenExpiredAt == nil || time.Now().UTC().After(*row.ResetTokenExpiredAt) {
		return apperr.New(constant.ErrCodeResetTokenInvalidOrExpired, "Invalid reset token")
	}

	if err := service.UserRepository.UpdatePassword(ctx, row.UserID, request.NewPassword); err != nil {
		return err
	}
	return service.PasswordResetRepository.Update(ctx, row.ID, map[string]interface{}{
		"status":                 constant.PasswordResetStatusCompleted,
		"reset_token":            nil,
		"reset_token_expired_at": nil,
	})
}

func (service *userServiceImpl) UpdatePassword(request entity.UpdatePasswordRequest) error {
	return apperr.New(constant.ErrCodeDeprecatedPasswordPatch, "This endpoint is deprecated. Use POST /v1/users/reset-password")
}

func (service *userServiceImpl) ChangePassword(request entity.ChangePasswordRequest) error {
	employee, err := service.MEmployeeRepository.FindOneByEmail(request.Email)
	if err != nil {
		return err
	}

	// user, err := service.UserRepository.FindOneByEmailAndCustId(request.Email, employee.CustId)
	user, err := service.UserRepository.FindOneByEmailCustIdEmpId(request.Email, employee.CustId, employee.EmpId)
	if err != nil {
		return err
	}
	userUpdate := model.User{}
	if err = structs.Automapper(user, &userUpdate); err != nil {
		return errors.New("error automapper")
	}

	if !service.CheckPasswordHash(request.CurrentPassword, *user.Userpass) {
		return errors.New("invalid Old password")
	}
	if service.CheckPasswordHash(request.NewPassword, *user.Userpass) {
		return errors.New("new password must be different from the current password")
	}
	err = service.UserRepository.UpdatePassword(context.Background(), *user.UserId, request.NewPassword)
	if err != nil {
		return err
	}
	return nil
}

func mapForgotUserLookupErr(err error) error {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return apperr.New(constant.ErrCodeEmailNotRegistered, "Invalid email")
	}
	return err
}

func (service *userServiceImpl) resetTokenSecret() string {
	s := service.Config.Get("PASSWORD_RESET_HMAC_SECRET")
	if s != "" {
		return s
	}
	return service.Config.Get("JWT_SECRET_KEY")
}

func randomHexToken(nBytes int) (string, error) {
	b := make([]byte, nBytes)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

var (
	rePasswordHasLetter = regexp.MustCompile(`[A-Za-z]`)
	rePasswordHasDigit  = regexp.MustCompile(`[0-9]`)
)

func validatePasswordPolicy(pw string) error {
	if len(pw) < 8 {
		return apperr.New(constant.ErrCodePasswordPolicy, "Password must be at least 8 characters and include letters and numbers")
	}
	if !rePasswordHasLetter.MatchString(pw) || !rePasswordHasDigit.MatchString(pw) {
		return apperr.New(constant.ErrCodePasswordPolicy, "Password must be at least 8 characters and include letters and numbers")
	}
	return nil
}

// CheckPasswordHash compare password with hash
func (service *userServiceImpl) CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		fmt.Println(err.Error())
	}
	return err == nil
}
func (service *userServiceImpl) GetProfile(email string) (response entity.ProfileResponse, err error) {
	employee, err := service.MEmployeeRepository.FindOneByEmail(email)
	if err != nil {
		return response, err
	}

	var salesman model.MSalesmanRead
	salesman, err = service.MSalesmanRepository.FindOneByEmpId(employee.CustId, employee.EmpId)
	if err != nil {
		return response, err
	}

	customer, err := service.MCustomerRepository.FindOneByCustId(employee.CustId)
	if err != nil {
		return response, err
	}

	distributor, _ := service.MEmployeeRepository.FindOneByDistributor(customer.DistributorId)
	// if err != nil { // comment because for tracing
	// 	return response, err
	// }

	// user, err := service.UserRepository.FindOneByEmailAndCustId(request.Email, employee.CustId)
	user, err := service.UserRepository.FindOneByEmailCustIdEmpId(email, employee.CustId, employee.EmpId)
	if err != nil {
		return response, err
	}

	location, err := service.UserLocationRepository.FindFirst()
	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
		default:
			return response, err
		}
	}
	response.Duration = 5
	response.Distance = 25
	if location != nil {
		response.Duration = location.Duration
		response.Distance = location.Distance
	}
	response.UserRole = safeString(user.RoleName)
	response.UserImg = safeString(user.ImageUrl)
	response.SalesCode = employee.EmpCode
	response.SalesName = salesman.SalesName
	response.CustID = customer.CustId
	response.Custname = customer.CustName
	response.Email = safeString(employee.Email)
	response.SalesTeamName = salesman.SalesTeamName
	response.LangID = safeString(user.LangId)
	response.MobileNo = safeString(user.MobileNo)
	response.Whatsapp = safeString(user.Whatsapp)
	response.IsValidRoute = salesman.SmValidRoute
	response.MaxRadius = salesman.SmRadius
	response.IsCaptureOutlet = false
	if distributor.DistributorID != 0 {
		response.DistributorCode = distributor.DistributorCode
		response.DistributorName = distributor.DistributorName
	}

	return response, nil
}

func safeString(s *string) string {
	if s == nil {
		return "" // Default empty string if nil
	}
	return *s
}

func (service *userServiceImpl) Update(userId int64, request entity.UpdateUserImageBody) (err error) {
	c := context.Background()

	// employee, err := service.MEmployeeRepository.FindOneByEmail(*request.Email)
	// if err != nil {
	// 	return err
	// }

	// user, err := service.UserRepository.FindOneByEmailCustIdEmpId(*request.Email, employee.CustId, employee.EmpId)
	// if err != nil {
	// 	return err
	// }
	user, err := service.UserRepository.FindOneByUserID(userId)
	if err != nil {
		return err
	}

	log.Println(user)

	var Model model.User
	err = structs.Automapper(request, &Model)
	if err != nil {
		return err
	}

	Model.CustId = ""
	err = service.Transaction.WithinTransaction(c, func(txCtx context.Context) error {
		err = service.UserRepository.Update(txCtx, userId, Model)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func (service *userServiceImpl) SendLocation(request entity.SendLocationRequest) (err error) {
	c := context.Background()
	err = service.UserLocationRepository.CreateLocation(c, &model.UserLocation{
		CustId:    request.CustID,
		EmpID:     request.EmpID,
		Latitude:  request.Latitude,
		Longitude: request.Longitude,
	})
	if err != nil {
		return err
	}

	return nil
}
