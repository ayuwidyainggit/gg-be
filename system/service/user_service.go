package service

import (
	"context"
	"errors"
	"log"
	"system/entity"
	"system/model"
	"system/pkg/config/env"
	"system/pkg/jwthelper"
	"system/pkg/mail/smtp"
	"system/pkg/str"
	"system/pkg/structs"
	"system/repository"
	"strconv"

	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	Login(entity.LoginRequest) (entity.LoginResponse, error)
	ForgotPassword(entity.ForgotPasswordRequest) error
	ForgotPasswordValidate(entity.ForgotPasswordValidateRequest) error
	UpdatePassword(request entity.UpdatePasswordRequest) (err error)
	// FindOneByUsername(string) (entity.LoginResponse, error)
	Store(request entity.CreateUserBody) (err error)
	Update(userID int64, request entity.UpdateUserBody) (err error)
	Detail(userID int64, custID string) (response entity.UserResponse, err error)
	DetailCust(custID string) (response entity.CustomerResponse, err error)
	Delete(custId string, UserIdDel int, userId int64) (err error)
	List(dataFilter entity.GeneralQueryFilter, custId string) (data []entity.UserListResponse, total int64, lastPage int, err error)
	UserMenus(userId int64, custId string) (data []entity.WebMenuResp, err error)
	UserMenusAll(menuParam string) (data []entity.WebMenuResp, err error)
	UserMenusDesktop(userID int64) (data entity.DesktopMenuResp, err error)
}

type userServiceImpl struct {
	SmcMCustomerRepository repository.SmcMCustomerRepository
	UserRepository         repository.UserRepository
	MMenuRepository        repository.MMenuRepository
	Transaction            repository.Dbtransaction
	Cache                  repository.CacheRepository
	Config                 env.ConfigEnv
}

func NewUserService(config env.ConfigEnv, smcMCustomer repository.SmcMCustomerRepository, userRepository repository.UserRepository,
	mMenuRepository repository.MMenuRepository, transaction repository.Dbtransaction, cache repository.CacheRepository) *userServiceImpl {
	return &userServiceImpl{
		Config:                 config,
		SmcMCustomerRepository: smcMCustomer,
		UserRepository:         userRepository,
		MMenuRepository:        mMenuRepository,
		Transaction:            transaction,
		Cache:                  cache,
	}
}

func (service *userServiceImpl) Login(request entity.LoginRequest) (response entity.LoginResponse, err error) {
	user, err := service.UserRepository.FindOneByEmail(request.Email)
	if err != nil {
		return response, errors.New("email not found")
	}

	if !service.CheckPasswordHash(request.Password, *user.Userpass) {
		return response, errors.New("invalid password")
	}

	credentials := make([]string, 0)
	additionalParam := make(map[string]interface{}, 0)

	customer, err := service.SmcMCustomerRepository.FindOneByCustId(user.CustId)
	if err != nil {
		return response, errors.New("customer data not found")
	}

	userData := entity.UserData{}
	err = structs.Automapper(user, &userData)
	if err != nil {
		return response, errors.New("error automapper")
	}
	userData.ParentCustId = customer.ParentCustID
	userData.DistPriceGrpId = customer.DistPriceGrpId
	userData.DistributorID = customer.DistributorID

	response.Token, err = jwthelper.GenerateNewToken(userData, credentials, additionalParam)
	if err != nil {
		return response, err
	}

	err = structs.Automapper(userData, &response)
	if err != nil {
		return response, errors.New("error automapper")
	}

	return response, err
}

func (service *userServiceImpl) ForgotPassword(request entity.ForgotPasswordRequest) (err error) {
	user, err := service.UserRepository.FindOneByEmail(request.Email)
	if err != nil {
		return errors.New("email not found")
	}

	otpCode, _ := str.Generate(`[0-9]{4}`)
	// save to redis
	err = service.Cache.SaveOTP(context.Background(), *user.UserId, otpCode)
	if err != nil {
		return err
	}

	subject := "Forgot Password"
	template := "forgot_password_otp.html"
	// var envCfg env.ConfigEnv
	data := entity.EmailData{
		Email:   request.Email,
		OtpCode: otpCode,
	}
	log.Println("data:", structs.StructToJson(data))
	err = smtp.BaseSendEmailWithTemplate(service.Config, request.Email, subject, template, data)
	if err != nil {
		return errors.New("ERROR SEND RECEIPT EMAIL: " + err.Error())
	}

	return nil
}

func (service *userServiceImpl) ForgotPasswordValidate(request entity.ForgotPasswordValidateRequest) (err error) {
	user, err := service.UserRepository.FindOneByEmail(request.Email)
	if err != nil {
		return errors.New("email not found")
	}

	getOtp, err := service.Cache.GetOTP(context.Background(), *user.UserId)
	if err != nil {
		return err
	}

	if getOtp == nil {
		return errors.New("OTP expired")
	}

	if request.OtpCode != *getOtp {
		return errors.New("Invalid OTP")
	} else {
		err := service.Cache.DeleteOTP(context.Background(), *user.UserId)
		if err != nil {
			return err
		}
	}

	return nil
}

func (service *userServiceImpl) Store(request entity.CreateUserBody) (err error) {
	c := context.Background()
	var userModel model.User
	user, err := service.UserRepository.FindOneByUserID(request.CreatedBy)
	if err != nil {
		return err
	}
	if user.IsAdmin != nil {
		if !*user.IsAdmin {
			return errors.New("edit user not allowed")
		}
	}

	err = structs.Automapper(request, &userModel)
	if err != nil {
		return err
	}
	err = service.Transaction.WithinTransaction(c, func(txCtx context.Context) error {
		err := service.UserRepository.Store(txCtx, &userModel)
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

func (service *userServiceImpl) UpdatePassword(request entity.UpdatePasswordRequest) (err error) {
	user, err := service.UserRepository.FindOneByEmail(request.Email)
	if err != nil {
		return errors.New("email not found")
	}

	err = service.UserRepository.UpdatePassword(context.Background(), *user.UserId, request.Password)
	if err != nil {
		return err
	}
	return nil
}

func (service *userServiceImpl) Update(userID int64, request entity.UpdateUserBody) (err error) {
	c := context.Background()
	var userModel model.User
	user, err := service.UserRepository.FindOneByUserID(request.UpdatedBy)
	if err != nil {
		return err
	}
	if user.IsAdmin != nil {
		if !*user.IsAdmin {
			return errors.New("edit user not allowed")
		}
	}

	err = structs.Automapper(request, &userModel)
	if err != nil {
		return err
	}
	err = service.Transaction.WithinTransaction(c, func(txCtx context.Context) error {
		err := service.UserRepository.Update(txCtx, userID, userModel)
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

func (service *userServiceImpl) Detail(userID int64, custID string) (response entity.UserResponse, err error) {
	user, err := service.UserRepository.FindDetail(userID, custID)
	if err != nil {
		return response, err
	}
	err = structs.Automapper(user, &response)
	if err != nil {
		return response, err
	}
	return response, nil
}

func (service *userServiceImpl) DetailCust(custID string) (response entity.CustomerResponse, err error) {
	user, err := service.UserRepository.FindDetailCust(custID)
	if err != nil {
		return response, err
	}
	err = structs.Automapper(user, &response)
	if err != nil {
		return response, err
	}
	return response, nil
}

func (service *userServiceImpl) Delete(custId string, UserIdDel int, userId int64) (err error) {
	c := context.Background()
	err = service.Transaction.WithinTransaction(c, func(txCtx context.Context) error {
		err = service.UserRepository.Delete(txCtx, custId, UserIdDel, userId)
		if err != nil {
			return err
		}
		return nil
	})

	return err
}
func (service *userServiceImpl) List(dataFilter entity.GeneralQueryFilter, custId string) (data []entity.UserListResponse, total int64, lastPage int, err error) {
	whAdjs, total, lastPage, err := service.UserRepository.FindAllByCustId(dataFilter, custId)
	if err != nil {
		return data, total, lastPage, err
	}
	if len(whAdjs) > 0 {
		for _, row := range whAdjs {
			var vResp entity.UserListResponse
			structs.Automapper(row, &vResp)

			data = append(data, vResp)
		}
	}
	return data, total, lastPage, err
}

func (service *userServiceImpl) UserMenus(userId int64, custId string) (data []entity.WebMenuResp, err error) {
	userMenus, err := service.MMenuRepository.FindAllMenu(userId, custId)
	if err != nil {
		return data, err
	}

	// log.Println("userMenus:", structs.StructToJson(userMenus))

	// We will use a map to simulate the database (simplified example!)
	WebMenuMap := entity.NewWebMenuMap()
	for _, r := range userMenus {
		webMenuResp := entity.WebMenuResp{}
		structs.Automapper(r, &webMenuResp)
		webMenuResp.IsHeader = false
		// if webMenuResp.Url != "" {
		// 	webMenuResp.TargetType = "iframe-tab"
		// }
		WebMenuMap.Db[r.ParentID] = append(WebMenuMap.Db[r.ParentID], webMenuResp)
	}

	for i := range userMenus {
		if userMenus[i].ParentID == "" {
			umParent := entity.WebMenuResp{}
			umParent.IsHeader = false
			structs.Automapper(userMenus[i], &umParent)
			parent := umParent // no parents
			parent.IsHeader = false
			WebMenuMap.SetChildrenRecursively(&parent)
			WebMenuMap.Append(parent)
		}
	}
	data = WebMenuMap.Resp
	// log.Println("data:", structs.StructToJson(data))

	return data, err
}

func (service *userServiceImpl) UserMenusAll(menuParam string) (data []entity.WebMenuResp, err error) {

	menuInt, err := strconv.Atoi(menuParam) 
	if err != nil {
		return data, err 
	}

	userMenus, err := service.MMenuRepository.FindAllMenuWithoutCustId(menuInt)
	if err != nil {
		return data, err
	}

	// log.Println("userMenus:", structs.StructToJson(userMenus))

	// We will use a map to simulate the database (simplified example!)
	WebMenuMap := entity.NewWebMenuMap()
	for _, r := range userMenus {
		webMenuResp := entity.WebMenuResp{}
		structs.Automapper(r, &webMenuResp)
		webMenuResp.IsHeader = false
		// if webMenuResp.Url != "" {
		// 	webMenuResp.TargetType = "iframe-tab"
		// }
		WebMenuMap.Db[r.ParentID] = append(WebMenuMap.Db[r.ParentID], webMenuResp)
	}

	for i := range userMenus {
		if userMenus[i].ParentID == "" {
			umParent := entity.WebMenuResp{}
			umParent.IsHeader = false
			structs.Automapper(userMenus[i], &umParent)
			parent := umParent // no parents
			parent.IsHeader = false
			WebMenuMap.SetChildrenRecursively(&parent)
			WebMenuMap.Append(parent)
		}
	}
	data = WebMenuMap.Resp
	// log.Println("data:", structs.StructToJson(data))

	return data, err
}

func (service *userServiceImpl) UserMenusDesktop(userID int64) (data entity.DesktopMenuResp, err error) {
	userMenus, err := service.MMenuRepository.FindAllMenuDesktop(userID)
	if err != nil {
		return data, err
	}

	// We will use a map to simulate the database (simplified example!)
	DesktopMenuMap := entity.NewDesktopMenuMap()
	for _, r := range userMenus {
		DesktopMenu := entity.DesktopMenu{}
		structs.Automapper(r, &DesktopMenu)
		DesktopMenuMap.Db[r.ParentID] = append(DesktopMenuMap.Db[r.ParentID], DesktopMenu)
	}
	for i := range userMenus {
		if userMenus[i].ParentID == "" {
			umParent := entity.DesktopMenu{}
			structs.Automapper(userMenus[i], &umParent)
			parent := umParent // no parents
			DesktopMenuMap.SetChildrenRecursively(&parent)
			DesktopMenuMap.Append(parent)
		}
	}

	hakAksesPkgs, err := service.MMenuRepository.FindAllHakAksesPkg(userID)
	if err != nil {
		return data, err
	}

	for _, hakAksesPkg := range hakAksesPkgs {
		pkg := entity.DesktopPackage{}
		structs.Automapper(hakAksesPkg, &pkg)
		data.Package = append(data.Package, pkg)
	}
	data.Menu = DesktopMenuMap.Resp

	return data, err
}

// CheckPasswordHash compare password with hash
func (service *userServiceImpl) CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
