package service

import (
	"errors"
	"mobile/entity"
	"mobile/model"
	"mobile/pkg/config/env"
	"mobile/pkg/constant"
	"mobile/pkg/times"
	"mobile/repository"
	"time"
)

type AttendanceService interface {
	AttendanceRequest(request entity.AttendanceRequest) (response entity.AttendanceResponse, err error)
	AttendanceGet(request entity.AttendanceGetRequest) (response entity.AttendanceResponse, err error)
	AttendanceCheck(request entity.AttendanceCheckRequest) (response entity.AttendanceCheckResponse, err error)
}
type AttendanceServiceImpl struct {
	Config              env.ConfigEnv
	Transaction         repository.Dbtransaction
	AttendanceRepo      repository.AttendanceRepository
	MEmployeeRepository repository.MEmployeeRepository
}

func NewAttendanceService(
	config env.ConfigEnv,
	transaction repository.Dbtransaction,
	attendanceRepo repository.AttendanceRepository,
	mEmployee repository.MEmployeeRepository,
) *AttendanceServiceImpl {
	return &AttendanceServiceImpl{
		Config:              config,
		Transaction:         transaction,
		AttendanceRepo:      attendanceRepo,
		MEmployeeRepository: mEmployee,
	}
}
func (service *AttendanceServiceImpl) AttendanceRequest(request entity.AttendanceRequest) (response entity.AttendanceResponse, err error) {
	now, err := times.GetCurrentTime()
	if err != nil {
		return response, err
	}

	var (
		isAlreadyCheckin bool
	)

	employee, err := service.MEmployeeRepository.FindOneByEmailCustID(request.Email, request.CustID)
	if err != nil {
		return response, err
	}

	from := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	to := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, now.Location())

	attendancesExist, err := service.AttendanceRepo.FindBetween(request.CustID, employee.EmpCode, from, to)
	if err != nil {
		return response, err
	}

	for _, attendanceExist := range attendancesExist {
		createdAt := attendanceExist.CreatedAt

		if *attendanceExist.Type == entity.TYPE_CHECKIN_ID {
			isAlreadyCheckin = true
			response.CheckIn = &createdAt // set response checkin time
		}

		if *request.Type == entity.TYPE_CHECKIN {
			if *attendanceExist.Type == entity.TYPE_CHECKIN_ID {
				return response, errors.New("employee already checkin")
			}
		}
		if *request.Type == entity.TYPE_CHECKOUT {
			if *attendanceExist.Type == entity.TYPE_CHECKOUT_ID {
				return response, errors.New("employee already checkout")
			}
		}
	}

	if *request.Type == entity.TYPE_CHECKOUT {
		if !isAlreadyCheckin {
			return response, errors.New("employee must be checkin before checkout")
		}
	}

	attType := entity.ConvTypeStringToCode(*request.Type)
	attendenceModel := model.Attendance{
		CustID:    request.CustID,
		EmpCode:   &employee.EmpCode,
		Latitude:  &request.Latitude,
		Longitude: &request.Longitude,
		Type:      &attType,
		CreatedAt: now,
	}
	if *request.Type == entity.TYPE_CHECKIN && request.LeaveID != nil {
		attendenceModel.LeaveID = request.LeaveID
	}
	err = service.AttendanceRepo.Store(&attendenceModel)
	if err != nil {
		return response, err
	}
	if *request.Type == entity.TYPE_CHECKIN {
		response.CheckIn = &now // set response checkin time
	}
	if *request.Type == entity.TYPE_CHECKOUT {
		response.CheckOut = &now // set response checkout time
	}
	return response, err
}
func (service *AttendanceServiceImpl) AttendanceGet(request entity.AttendanceGetRequest) (response entity.AttendanceResponse, err error) {
	now, err := times.GetCurrentTime()
	if err != nil {
		return response, err
	}
	response.CurrentTime = now

	from := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	to := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, now.Location())

	employee, err := service.MEmployeeRepository.FindOneByEmailCustID(request.Email, request.CustID)
	if err != nil {
		return response, err
	}
	attendancesExist, err := service.AttendanceRepo.FindBetween(request.CustID, employee.EmpCode, from, to)
	if err != nil {
		return response, err
	}

	for _, attendanceExist := range attendancesExist {
		createdAt := attendanceExist.CreatedAt
		if *attendanceExist.Type == entity.TYPE_CHECKIN_ID {
			response.CheckIn = &createdAt // set response checkin time
		}
		if *attendanceExist.Type == entity.TYPE_CHECKOUT_ID {
			response.CheckOut = &createdAt // set response checkout time
		}
	}
	return response, err
}

func (service *AttendanceServiceImpl) AttendanceCheck(request entity.AttendanceCheckRequest) (response entity.AttendanceCheckResponse, err error) {
	date := time.Unix(request.Date, 0)

	isDistributor := request.DistributorID != nil && *request.DistributorID > 0

	planCount, err := service.AttendanceRepo.CountPlanBySalesmanIDAndDate(request.EmpID, date, isDistributor)
	if err != nil {
		return response, err
	}

	if planCount == 0 {
		planCountFallback, errFallback := service.AttendanceRepo.CountPlanBySalesmanIDAndDate(request.EmpID, date, !isDistributor)
		if errFallback == nil && planCountFallback > 0 {
			planCount = planCountFallback
		}
	}

	// Initialize response data
	response.Data.Plan = int(planCount)

	var salesmanType string // "Taking Order", "Canvas", or "Taking Order + Canvas"
	var stock int64 = 0

	if isDistributor {
		salesmanInfo, err := service.AttendanceRepo.GetSalesmanInfoWithCanvas(request.EmpID)
		if err != nil {
			return response, err
		}

		// Determine salesman type based on opr_type and opr_type_canvas
		// opr_type: "O" if is_taking_order = true, "" if false
		// opr_type_canvas: "C" if is_active = true, "" if false
		isTakingOrder := salesmanInfo.OprType != nil && *salesmanInfo.OprType == "O"
		isCanvas := salesmanInfo.OprTypeCanvas != nil && *salesmanInfo.OprTypeCanvas == "C"

		if isTakingOrder && isCanvas {
			salesmanType = "Taking Order + Canvas"
		} else if isCanvas {
			salesmanType = "Canvas"
		} else if isTakingOrder {
			salesmanType = "Taking Order"
		} else {
			salesmanType = "Unknown"
		}

		response.Data.EmpID = &salesmanInfo.EmpID
		if salesmanInfo.EmpCode != nil {
			response.Data.EmpCode = salesmanInfo.EmpCode
		}
		if salesmanInfo.EmpName != nil {
			response.Data.EmpName = salesmanInfo.EmpName
		}

		// Set opr_type: "O" if is_taking_order = true, "" if false
		if salesmanInfo.OprType != nil {
			response.Data.OprType = salesmanInfo.OprType
		} else {
			emptyStr := ""
			response.Data.OprType = &emptyStr
		}

		// Set opr_type_canvas: "C" if is_active = true, "" if false
		if salesmanInfo.OprTypeCanvas != nil {
			response.Data.OprTypeCanvas = salesmanInfo.OprTypeCanvas
		} else {
			emptyStr := ""
			response.Data.OprTypeCanvas = &emptyStr
		}

		if salesmanInfo.WhID != nil {
			response.Data.WhID = salesmanInfo.WhID
		}
		if salesmanInfo.WhCode != nil {
			response.Data.WhCode = salesmanInfo.WhCode
		} else {
			// Empty string for Taking Order only
			emptyStr := ""
			response.Data.WhCode = &emptyStr
		}
		if salesmanInfo.WhNameCanvas != nil {
			response.Data.WhNameCanvas = salesmanInfo.WhNameCanvas
		} else {
			// Empty string for Taking Order only
			emptyStr := ""
			response.Data.WhNameCanvas = &emptyStr
		}

		// Query STOCK (only for distributor with Canvas type)
		if salesmanType == "Canvas" || salesmanType == "Taking Order + Canvas" {
			stock, err = service.AttendanceRepo.GetWarehouseStockByEmpIDAndCustID(request.EmpID, request.CustID)
			if err != nil {
				return response, err
			}
			stockInt := int(stock)
			response.Data.Stock = &stockInt
		} else {
			response.Data.Stock = nil
		}
	}

	// Validate based on user role and salesman type
	if isDistributor {
		// Distributor validation
		if salesmanType == "Taking Order" {
			// Distributor Taking Order: Validate PLAN only
			if planCount > 0 {
				response.Success = true
				response.Message = constant.CHECKIN_AVAILABLE
				response.Description = ""
			} else {
				response.Success = false
				response.Message = constant.CHECKIN_UNAVAILABLE
				response.Description = constant.CHECKIN_UNAVAILABLE_NO_PLAN
			}
		} else if salesmanType == "Canvas" || salesmanType == "Taking Order + Canvas" {
			// Distributor Canvas or Taking Order + Canvas: Validate PLAN AND STOCK
			if planCount > 0 && stock > 0 {
				response.Success = true
				response.Message = constant.CHECKIN_AVAILABLE
				response.Description = ""
			} else if planCount == 0 && stock > 0 {
				response.Success = false
				response.Message = constant.CHECKIN_UNAVAILABLE
				response.Description = constant.CHECKIN_UNAVAILABLE_NO_PLAN
			} else if planCount > 0 && stock == 0 {
				response.Success = false
				response.Message = constant.CHECKIN_UNAVAILABLE
				response.Description = constant.CHECKIN_UNAVAILABLE_NO_STOCK
			} else {
				response.Success = false
				response.Message = constant.CHECKIN_UNAVAILABLE
				response.Description = constant.CHECKIN_UNAVAILABLE_NO_PLAN_AND_STOCK
			}
		} else {
			// Unknown salesman type, default to PLAN validation
			if planCount > 0 {
				response.Success = true
				response.Message = constant.CHECKIN_AVAILABLE
				response.Description = ""
			} else {
				response.Success = false
				response.Message = constant.CHECKIN_UNAVAILABLE
				response.Description = constant.CHECKIN_UNAVAILABLE_NO_PLAN
			}
		}
	} else {
		// Principal validation: Validate PLAN only
		if planCount > 0 {
			response.Success = true
			response.Message = constant.CHECKIN_AVAILABLE
			response.Description = ""
		} else {
			response.Success = false
			response.Message = constant.CHECKIN_UNAVAILABLE
			response.Description = constant.CHECKIN_UNAVAILABLE_NO_PLAN
		}
	}

	return response, nil
}
