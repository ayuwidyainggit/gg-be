package attendance

import (
	"context"
	"scyllax-pjp/constant"
	"scyllax-pjp/data/request"
	"scyllax-pjp/data/response"
	"scyllax-pjp/repository/attendance"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// AttendanceService defines contract for attendance business logic
type AttendanceService interface {
	CheckAttendance(ctx context.Context, req request.AttendanceCheckRequest) response.AttendanceCheckResponse
}

type attendanceService struct {
	repo attendance.AttendanceRepository
	db   *gorm.DB
}

// NewAttendanceService creates a new attendance service instance
func NewAttendanceService(repo attendance.AttendanceRepository, db *gorm.DB) AttendanceService {
	return &attendanceService{
		repo: repo,
		db:   db,
	}
}

// CheckAttendance validates if check-in is available based on route plan (PJP) and warehouse stock
func (s *attendanceService) CheckAttendance(ctx context.Context, req request.AttendanceCheckRequest) response.AttendanceCheckResponse {
	requestID := uuid.New().String()
	date := time.Unix(req.Date, 0)

	res := response.AttendanceCheckResponse{
		RequestID: requestID,
	}

	// Principal Flow: Only validate plan
	if req.DistributorID == nil {
		return s.handlePrincipalValidation(ctx, req.EmpID, date, res)
	}

	// Distributor Flow: Validate based on salesman type
	return s.handleDistributorValidation(ctx, req.EmpID, date, res)
}

// handlePrincipalValidation processes validation for principal users
func (s *attendanceService) handlePrincipalValidation(ctx context.Context, empID int, date time.Time, res response.AttendanceCheckResponse) response.AttendanceCheckResponse {
	plan := s.repo.GetPrincipalPlanCount(ctx, s.db, empID, date)

	res.Data = &response.AttendanceCheckData{
		Plan: plan,
	}

	if plan > 0 {
		res.Success = true
		res.Message = constant.AttendanceMessage.CheckInAvailable
	} else {
		res.Success = false
		res.Message = constant.AttendanceMessage.CheckInUnavailable
		res.Description = constant.AttendanceDescription.NoPlanPrincipal
	}

	return res
}

// handleDistributorValidation processes validation for distributor users
func (s *attendanceService) handleDistributorValidation(ctx context.Context, empID int, date time.Time, res response.AttendanceCheckResponse) response.AttendanceCheckResponse {
	info, err := s.repo.GetSalesmanWithCanvas(ctx, s.db, empID)
	if err != nil {
		res.Success = false
		res.Message = constant.AttendanceDescription.ErrorFetchingSalesman
		return res
	}

	plan := s.repo.GetDistributorPlanCount(ctx, s.db, empID, date)

	res.Data = s.buildSalesmanData(info, plan)

	// Determine salesman type and validate accordingly
	if s.isTakingOrderOnly(info) {
		return s.validateTakingOrderOnly(res, plan)
	}

	// Canvas or TO+Canvas: Validate both plan and stock
	return s.validateCanvasOrTOCanvas(ctx, empID, info, res, plan)
}

// buildSalesmanData constructs the base salesman data for response
func (s *attendanceService) buildSalesmanData(info *attendance.SalesmanCanvasInfo, plan int) *response.AttendanceCheckData {
	data := &response.AttendanceCheckData{
		EmpID:         info.EmpID,
		EmpCode:       info.EmpCode,
		EmpName:       info.EmpName,
		OprType:       info.OprType,
		OprTypeCanvas: "",
		Plan:          plan,
	}

	if info.OprTypeCanvas != nil {
		data.OprTypeCanvas = *info.OprTypeCanvas
	}

	return data
}

// isTakingOrderOnly checks if salesman is Taking Order only (not Canvas)
func (s *attendanceService) isTakingOrderOnly(info *attendance.SalesmanCanvasInfo) bool {
	return info.OprType == constant.OprTypeTakingOrder &&
		(info.OprTypeCanvas == nil || !info.IsActiveCanvas)
}

// validateTakingOrderOnly handles validation for TO only salesman
func (s *attendanceService) validateTakingOrderOnly(res response.AttendanceCheckResponse, plan int) response.AttendanceCheckResponse {
	// According to docs, TO only response should include empty wh_id, wh_code, wh_name_canvas, stock
	res.Data.WhID = ""
	res.Data.WhCode = ""
	res.Data.WhNameCanvas = ""
	res.Data.Stock = ""

	if plan > 0 {
		res.Success = true
		res.Message = constant.AttendanceMessage.CheckInAvailable
	} else {
		res.Success = false
		res.Message = constant.AttendanceMessage.CheckInUnavailable
		res.Description = constant.AttendanceDescription.NoPlanDistributor
	}

	return res
}

// validateCanvasOrTOCanvas handles validation for Canvas or TO+Canvas salesman
func (s *attendanceService) validateCanvasOrTOCanvas(ctx context.Context, empID int, info *attendance.SalesmanCanvasInfo, res response.AttendanceCheckResponse, plan int) response.AttendanceCheckResponse {
	whID, stock, err := s.repo.GetWarehouseStock(ctx, s.db, empID)

	if err == nil {
		res.Data.WhID = whID
		res.Data.Stock = stock
		res.Data.WhCode = s.getStringOrEmpty(info.WhCode)
		res.Data.WhNameCanvas = s.getStringOrEmpty(info.WhName)
	} else {
		// Set default values if stock query fails
		res.Data.Stock = 0
		res.Data.WhID = 0
		res.Data.WhCode = ""
		res.Data.WhNameCanvas = ""
		stock = 0
	}

	// Apply validation rules based on plan and stock
	return s.applyPlanAndStockValidation(res, plan, stock)
}

// applyPlanAndStockValidation applies business rules for Canvas/TO+Canvas validation
func (s *attendanceService) applyPlanAndStockValidation(res response.AttendanceCheckResponse, plan, stock int) response.AttendanceCheckResponse {
	hasPlan := plan > 0
	hasStock := stock > 0

	switch {
	case hasPlan && hasStock:
		res.Success = true
		res.Message = constant.AttendanceMessage.CheckInAvailable
	case !hasPlan && hasStock:
		res.Success = false
		res.Message = constant.AttendanceMessage.CheckInUnavailable
		res.Description = constant.AttendanceDescription.NoPlanDistributor
	case hasPlan && !hasStock:
		res.Success = false
		res.Message = constant.AttendanceMessage.CheckInUnavailable
		res.Description = constant.AttendanceDescription.NoStock
	default:
		res.Success = false
		res.Message = constant.AttendanceMessage.CheckInUnavailable
		res.Description = constant.AttendanceDescription.NoPlanAndNoStock
	}

	return res
}

// getStringOrEmpty returns the string value or empty string if nil
func (s *attendanceService) getStringOrEmpty(str *string) string {
	if str != nil {
		return *str
	}
	return ""
}
