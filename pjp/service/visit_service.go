package service

import (
	"context"
	"errors"
	"fmt"
	"math"
	"net/url"
	"path/filepath"
	"scyllax-pjp/constant"
	"scyllax-pjp/data/entity"
	"scyllax-pjp/data/request"
	"scyllax-pjp/helper"
	"scyllax-pjp/model"
	"scyllax-pjp/repository"
	arrivalreport "scyllax-pjp/repository/arrival_report"
	"scyllax-pjp/repository/outlet_visit_principle"
	"scyllax-pjp/repository/pjp"
	"strings"

	"scyllax-pjp/data/response"
	"time"

	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type VisitService interface {
	FinishVisit(ctx context.Context, request request.FinishVisitRequest)
	ArriveVisit(ctx context.Context, request request.ArriveVisitRequest, custId string, userId int64)
	SkipVisit(ctx context.Context, request request.SkipVisitRequest)
	ResumeVisit(ctx context.Context, request request.ResumeVisitRequest)
	LeaveVisit(ctx context.Context, request request.LeaveVisitRequest)
	UnloadVisit(ctx context.Context, request request.OnholdVisitRequest)
	OutletVisit(ctx context.Context, request request.OutletVisitRequest)
	GetAllOutletBySalesCode(ctx context.Context, salesCode, custId, date, routeCode string) (responses []response.OutletResponse)
	SummaryVisit(ctx context.Context, dataFilter entity.SummaryQueryFilter) (res response.SummaryResponse)
	TravelList(ctx context.Context, outletVisitId int, customerID string) (response response.TodoListResponse)
	SummaryVisitStatus(ctx context.Context, dataFilter entity.SummaryQueryFilter) (res response.VisitStatusResponse)
	GetVisitOutletList(ctx context.Context, dataFilter entity.SummaryQueryFilter) (res []response.OutletVisitListResponse)
	GetSalesmanReport(ctx context.Context, dataFilter entity.SalesmanReportQueryFilter) entity.SalesmanReport
}

type VisitServiceImpl struct {
	RouteOutletRepo              repository.RouteOutletRepository
	OutletVisitListRepo          repository.OutletVisitRepo
	OutletVisitListPrincipleRepo outlet_visit_principle.OutletVisitPrincipleRepository
	PjpRepo                      pjp.PjpRepository
	OutletCrRepo                 repository.OutletCrRepository
	ArrivalReportRepo            arrivalreport.ArrivalReportRepository
	CustomerRepo                 repository.CustomerRepository

	Validate *validator.Validate
	Db       *gorm.DB
}

func NewVisitServiceImpl(RouteOutletRepo repository.RouteOutletRepository, OutletVisitListRepo repository.OutletVisitRepo, OutletVisitListPrincipleRepo outlet_visit_principle.OutletVisitPrincipleRepository, PjpRepo pjp.PjpRepository, OutletCrRepo repository.OutletCrRepository, arrivalReportRepo arrivalreport.ArrivalReportRepository, customerRepo repository.CustomerRepository, validate *validator.Validate, db *gorm.DB) VisitService {
	return &VisitServiceImpl{
		RouteOutletRepo:              RouteOutletRepo,
		OutletVisitListRepo:          OutletVisitListRepo,
		OutletVisitListPrincipleRepo: OutletVisitListPrincipleRepo,
		PjpRepo:                      PjpRepo,
		OutletCrRepo:                 OutletCrRepo,
		ArrivalReportRepo:            arrivalReportRepo,
		CustomerRepo:                 customerRepo,
		Validate:                     validate,
		Db:                           db,
	}
}

func (service *VisitServiceImpl) FinishVisit(ctx context.Context, request request.FinishVisitRequest) {
	err := service.Validate.Struct(request)
	helper.ErrorPanic(err)

	service.OutletVisitListRepo.UpdateOutletVisitListMultiRowColumnAt(ctx, request.SalesmanCode, request.CustID, "finish", request.CurrentTime, request.Date)
}

// ArriveVisit processes arrive visit request with file upload and optional location update
// If is_update_location is true, saves location change request to outlet_cr for approval
// Location is not updated directly to m_outlet until approved
// Also saves arrival data to arrival_report table if geotaging data is provided
func (service *VisitServiceImpl) ArriveVisit(ctx context.Context, req request.ArriveVisitRequest, custId string, userId int64) {
	if err := service.Validate.Struct(req); err != nil {
		helper.ErrorPanic(err)
	}

	// Update arrive_at timestamp in outlet_visit_list
	var isPrinciple bool
	isPrinciple, errFind := service.CustomerRepo.CheckIsPrinciple(ctx, req.CustID)
	if errFind != nil {
		switch {
		case errors.Is(errFind, gorm.ErrRecordNotFound):
		default:
			helper.ErrorPanic(errFind)
		}
	}

	var (
		// Extract file metadata from file_url
		fileName      = extractFileNameFromURL(req.FileUrl)
		fileType      = extractFileTypeFromURL(req.FileUrl)
		mediaCategory = determineMediaCategory(fileType)

		// Location status
		locationStatus   = 0
		isUpdateLocation = false
	)

	if isPrinciple {
		service.OutletVisitListPrincipleRepo.UpdateOutletVisitListColumnAt(ctx, service.Db, "arrive_at", req.CurrentTime, req.Date, req.Id)
		if req.LocationStatus != nil {
			switch *req.LocationStatus {
			case response.LocationStatusInRadius:
				locationStatus = 1
			default:
				locationStatus = 0
			}
		}

		// Prepare file information for database update
		isUpdateLocation = req.IsUpdateLocation != nil && *req.IsUpdateLocation
		fileInfo := model.OutletVisitListPrinciple{
			IsUpdateLocation: isUpdateLocation,
			FileName:         fileName,
			FileType:         fileType,
			MediaCategory:    mediaCategory,
			FileUrl:          req.FileUrl,
			FileSize:         nil, // Not available from URL
			FileBase64:       "",  // Not available from URL
			PhotoPath:        req.FileUrl,
			Folder:           "", // Not provided in request
			DistanceMeter:    req.DistanceMeter,
			AllowedRadius:    req.AllowedRadius,
			LocationStatus:   &locationStatus,
		}

		// Set latitude and longitude if provided
		if req.Latitude != nil {
			fileInfo.Latitude = fmt.Sprintf("%f", *req.Latitude)
		}
		if req.Longitude != nil {
			fileInfo.Longitude = fmt.Sprintf("%f", *req.Longitude)
		}

		if err := service.OutletVisitListPrincipleRepo.UpdateOutletVisitListWithFile(ctx, service.Db, req.Id, req.Date, fileInfo); err != nil {
			helper.ErrorPanic(err)
		}
	} else {
		service.OutletVisitListRepo.UpdateOutletVisitListColumnAt(ctx, req.SalesmanCode, req.CustID, "arrive_at", req.CurrentTime, req.Date, req.Id)
		if req.LocationStatus != nil {
			switch *req.LocationStatus {
			case response.LocationStatusInRadius:
				locationStatus = 1
			default:
				locationStatus = 0
			}
		}

		// Prepare file information for database update
		isUpdateLocation = req.IsUpdateLocation != nil && *req.IsUpdateLocation
		fileInfo := model.OutletVisitList{
			IsUpdateLocation: isUpdateLocation,
			FileName:         fileName,
			FileType:         fileType,
			MediaCategory:    mediaCategory,
			FileUrl:          req.FileUrl,
			FileSize:         nil, // Not available from URL
			FileBase64:       "",  // Not available from URL
			PhotoPath:        req.FileUrl,
			Folder:           "", // Not provided in request
			DistanceMeter:    req.DistanceMeter,
			AllowedRadius:    req.AllowedRadius,
			LocationStatus:   &locationStatus,
		}

		// Set latitude and longitude if provided
		if req.Latitude != nil {
			fileInfo.Latitude = fmt.Sprintf("%f", *req.Latitude)
		}
		if req.Longitude != nil {
			fileInfo.Longitude = fmt.Sprintf("%f", *req.Longitude)
		}

		// Update outlet_visit_list with file information
		if err := service.OutletVisitListRepo.UpdateOutletVisitListWithFile(ctx, req.Id, req.Date, fileInfo); err != nil {
			helper.ErrorPanic(err)
		}
	}

	// Sync to mobile.visits table
	if req.OutletID != nil && req.Latitude != nil && req.Longitude != nil {
		service.OutletVisitListRepo.InsertOrUpdateMobileVisit(
			ctx,
			req.CustID,
			req.SalesmanCode,
			int(*req.OutletID),
			*req.Latitude,
			*req.Longitude,
			req.FileUrl,
			*req.CurrentTime,
			"Arrive",
		)
	}

	// If is_update_location is true, save to outlet_cr for approval
	// According to docs: location changes require approval and old data remains until approved
	if isUpdateLocation && req.OutletID != nil && req.Latitude != nil && req.Longitude != nil {
		latStr := fmt.Sprintf("%f", *req.Latitude)
		longStr := fmt.Sprintf("%f", *req.Longitude)
		if err := service.createOutletChangeRequest(ctx, custId, userId, *req.OutletID, latStr, longStr); err != nil {
			helper.ErrorPanic(err)
		}
	}

	// Save to arrival_report table if geotaging data is provided
	// Dreprecated, pjp.arrival_report is no longer used
	// if req.Activity != nil && req.LocationStatus != nil && req.OutletID != nil {
	// 	if err := service.saveArrivalReport(ctx, req, custId, userId); err != nil {
	// 		log.Printf("Error saving arrival report: %v", err)
	// 		// Don't panic, just log the error - arrival report is supplementary data
	// 	}
	// }
}

// saveArrivalReport saves arrival data to pjp.arrival_report table
func (service *VisitServiceImpl) saveArrivalReport(ctx context.Context, req request.ArriveVisitRequest, custId string, userId int64) error {
	now := time.Now()
	createdBy := fmt.Sprintf("%d", userId)

	var latStr, longStr *string
	if req.Latitude != nil {
		s := fmt.Sprintf("%f", *req.Latitude)
		latStr = &s
	}
	if req.Longitude != nil {
		s := fmt.Sprintf("%f", *req.Longitude)
		longStr = &s
	}

	arrivalReport := model.ArrivalReport{
		CustId:           custId,
		OutletId:         *req.OutletID,
		UserId:           userId,
		Activity:         *req.Activity,
		ArrivalLongitude: longStr,
		ArrivalLatitude:  latStr,
		OutletLongitude:  req.OutletLongitude,
		OutletLatitude:   req.OutletLatitude,
		DistanceMeter:    req.DistanceMeter,
		AllowedRadius:    req.AllowedRadius,
		LocationStatus:   req.LocationStatus,
		CreatedBy:        &createdBy,
		CreatedAt:        &now,
	}

	return service.ArrivalReportRepo.Create(service.Db, &arrivalReport)
}

// createOutletChangeRequest creates outlet change request record for approval
// Saves old and new location values to outlet_cr and outlet_cr_det tables
func (service *VisitServiceImpl) createOutletChangeRequest(ctx context.Context, custId string, userId int64, outletId int64, newLatitude, newLongitude string) error {
	// Get old values from mst.m_outlet
	oldLatitude, oldLongitude, err := service.OutletCrRepo.GetOutletLocation(ctx, service.Db, custId, outletId)
	if err != nil {
		return err
	}

	// Begin transaction
	tx := service.Db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	// Create outlet_cr record
	var createdBy *int64
	if userId > 0 {
		createdBy = &userId
	}
	outletCr := model.OutletCr{
		CustId:    custId,
		OutletId:  outletId,
		Source:    constant.OutletChangeRequestSourceMobile,
		Status:    constant.OutletChangeRequestStatusPending,
		CreatedBy: createdBy,
		CreatedAt: time.Now(),
	}

	outletCrId, err := service.OutletCrRepo.CreateOutletCr(ctx, tx, outletCr)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Create outlet_cr_det records for latitude and longitude
	details := []model.OutletCrDet{
		{
			OutletCrId: outletCrId,
			FieldName:  "latitude",
			OldValue:   &oldLatitude,
			NewValue:   &newLatitude,
		},
		{
			OutletCrId: outletCrId,
			FieldName:  "longitude",
			OldValue:   &oldLongitude,
			NewValue:   &newLongitude,
		},
	}

	for _, det := range details {
		if err := service.OutletCrRepo.CreateOutletCrDet(ctx, tx, det); err != nil {
			tx.Rollback()
			return err
		}
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

// extractFileNameFromURL extracts file name from URL using path.Base
// Returns the base name of the file from the URL path
func extractFileNameFromURL(fileURL string) string {
	parsedURL, err := url.Parse(fileURL)
	if err != nil {
		// If parsing fails, try to extract from string directly
		return filepath.Base(fileURL)
	}
	return filepath.Base(parsedURL.Path)
}

// extractFileTypeFromURL extracts file extension from URL and removes the dot
// Returns lowercase file extension without the dot (e.g., "jpg", "png", "mp4")
func extractFileTypeFromURL(fileURL string) string {
	parsedURL, err := url.Parse(fileURL)
	if err != nil {
		// If parsing fails, try to extract from string directly
		ext := filepath.Ext(fileURL)
		if ext != "" {
			return strings.ToLower(strings.TrimPrefix(ext, "."))
		}
		return ""
	}
	ext := filepath.Ext(parsedURL.Path)
	if ext != "" {
		return strings.ToLower(strings.TrimPrefix(ext, "."))
	}
	return ""
}

// determineMediaCategory determines media category based on file extension
// Returns "image" for jpg, jpeg, png or "video" for mp4
// Defaults to "image" if extension is not recognized
func determineMediaCategory(fileType string) string {
	fileTypeLower := strings.ToLower(fileType)
	switch fileTypeLower {
	case "jpg", "jpeg", "png":
		return "image"
	case "mp4":
		return "video"
	default:
		return "image" // Default to image if not recognized
	}
}

func (service *VisitServiceImpl) SkipVisit(ctx context.Context, request request.SkipVisitRequest) {
	err := service.Validate.Struct(request)
	helper.ErrorPanic(err)

	isPrinciple, errCheck := service.CustomerRepo.CheckIsPrinciple(ctx, request.CustID)
	if errCheck != nil {
		switch {
		case errors.Is(errCheck, gorm.ErrRecordNotFound):
		default:
			helper.ErrorPanic(errCheck)
		}
	}

	var (
		fileName         string
		fileType         string
		mediaCategory    string
		isUpdateLocation bool
	)
	if request.FileUrl != "" {
		fileName = extractFileNameFromURL(request.FileUrl)
		fileType = extractFileTypeFromURL(request.FileUrl)
		mediaCategory = determineMediaCategory(fileType)
		isUpdateLocation = request.IsUpdateLocation != nil && *request.IsUpdateLocation
	}

	if isPrinciple {
		var fileInfo model.OutletVisitListPrinciple
		if request.FileUrl != "" {
			fileInfo = model.OutletVisitListPrinciple{
				FileName:         fileName,
				FileType:         fileType,
				MediaCategory:    mediaCategory,
				FileUrl:          request.FileUrl,
				PhotoPath:        request.FileUrl,
				IsUpdateLocation: isUpdateLocation,
				SkipInOutlet:     request.InOutlet,
			}
			if request.Latitude != nil {
				fileInfo.Latitude = fmt.Sprintf("%f", *request.Latitude)
			}
			if request.Longitude != nil {
				fileInfo.Longitude = fmt.Sprintf("%f", *request.Longitude)
			}
		}

		service.OutletVisitListPrincipleRepo.UpdateOutletVisitListSkipColumnAt(ctx, service.Db, "skip_at", request.CurrentTime, request.Date, request.Id, request.SkipReason, request.InOutlet, fileInfo)
		return
	}

	var fileInfo model.OutletVisitList
	if request.FileUrl != "" {
		fileInfo = model.OutletVisitList{
			FileName:         fileName,
			FileType:         fileType,
			MediaCategory:    mediaCategory,
			FileUrl:          request.FileUrl,
			PhotoPath:        request.FileUrl,
			IsUpdateLocation: isUpdateLocation,
			SkipInOutlet:     request.InOutlet,
		}
		if request.Latitude != nil {
			fileInfo.Latitude = fmt.Sprintf("%f", *request.Latitude)
		}
		if request.Longitude != nil {
			fileInfo.Longitude = fmt.Sprintf("%f", *request.Longitude)
		}
	}
	service.OutletVisitListRepo.UpdateOutletVisitListSkipColumnAt(ctx, "skip_at", request.CurrentTime, request.Date, request.Id, request.SkipReason, request.InOutlet, fileInfo)
}

func (service *VisitServiceImpl) ResumeVisit(ctx context.Context, request request.ResumeVisitRequest) {
	err := service.Validate.Struct(request)
	helper.ErrorPanic(err)

	isPrinciple, errCheck := service.CustomerRepo.CheckIsPrinciple(ctx, request.CustID)
	if errCheck != nil {
		switch {
		case errors.Is(errCheck, gorm.ErrRecordNotFound):
		default:
			helper.ErrorPanic(errCheck)
		}
	}
	if isPrinciple {
		service.OutletVisitListPrincipleRepo.UpdateOutletVisitListColumnAt(ctx, service.Db, "resume_at", request.CurrentTime, request.Date, request.Id)
		// add on_hold to be null as resume is updated
		service.OutletVisitListPrincipleRepo.UpdateOutletVisitListColumnAt(ctx, service.Db, "on_hold", nil, request.Date, request.Id)
	} else {
		service.OutletVisitListRepo.UpdateOutletVisitListColumnAt(ctx, request.SalesmanCode, request.CustID, "resume_at", request.CurrentTime, request.Date, request.Id)
		// add on_hold to be null as resume is updated
		service.OutletVisitListRepo.UpdateOutletVisitListColumnAt(ctx, request.SalesmanCode, request.CustID, "on_hold", nil, request.Date, request.Id)
	}
}

func (service *VisitServiceImpl) LeaveVisit(ctx context.Context, request request.LeaveVisitRequest) {
	err := service.Validate.Struct(request)
	helper.ErrorPanic(err)
	fmt.Println("LeaveVisit request:", *request.LeaveLatitude, *request.LeaveLongitude)

	isPrinciple, errCheck := service.CustomerRepo.CheckIsPrinciple(ctx, request.CustID)
	if errCheck != nil {
		switch {
		case errors.Is(errCheck, gorm.ErrRecordNotFound):
		default:
			helper.ErrorPanic(errCheck)
		}
	}

	parseTime, errParse := time.Parse("2006-01-02", request.Date)
	if errParse != nil {
		helper.ErrorPanic(errParse)
	}

	if isPrinciple {
		datas := model.OutletVisitListPrinciple{
			LeaveAt: request.CurrentTime,
			Date:    parseTime,
		}
		if request.LeaveLatitude != nil && *request.LeaveLatitude != "" {
			datas.LeaveLatitude = *request.LeaveLatitude
		}
		if request.LeaveLongitude != nil && *request.LeaveLongitude != "" {
			datas.LeaveLongitude = *request.LeaveLongitude
		}
		err := service.OutletVisitListPrincipleRepo.UpdateOutletVisitListByID(ctx, service.Db, request.Id, datas)
		if err != nil {
			helper.ErrorPanic(err)
		}
	} else {
		datas := model.OutletVisitList{
			LeaveAt: request.CurrentTime,
			Date:    parseTime,
		}
		if request.LeaveLatitude != nil && *request.LeaveLatitude != "" {
			datas.LeaveLatitude = *request.LeaveLatitude
		}
		if request.LeaveLongitude != nil && *request.LeaveLongitude != "" {
			datas.LeaveLongitude = *request.LeaveLongitude
		}
		err := service.OutletVisitListRepo.UpdateOutletVisitListByID(ctx, service.Db, request.Id, datas)
		if err != nil {
			helper.ErrorPanic(err)
		}
	}
}

func (service *VisitServiceImpl) UnloadVisit(ctx context.Context, request request.OnholdVisitRequest) {
	err := service.Validate.Struct(request)
	helper.ErrorPanic(err)

	isPrinciple, errCheck := service.CustomerRepo.CheckIsPrinciple(ctx, request.CustID)
	if errCheck != nil {
		switch {
		case errors.Is(errCheck, gorm.ErrRecordNotFound):
		default:
			helper.ErrorPanic(errCheck)
		}
	}
	if isPrinciple {
		service.OutletVisitListPrincipleRepo.UpdateOutletVisitListColumnAt(ctx, service.Db, "on_hold", request.CurrentTime, request.Date, request.Id)
		service.OutletVisitListPrincipleRepo.UpdateOutletVisitListColumnAt(ctx, service.Db, "resume_at", nil, request.Date, request.Id)
	} else {
		service.OutletVisitListRepo.UpdateOutletVisitListColumnAt(ctx, request.SalesmanCode, request.CustID, "on_hold", request.CurrentTime, request.Date, request.Id)
		// add resume to be null as on_hold is updated
		service.OutletVisitListRepo.UpdateOutletVisitListColumnAt(ctx, request.SalesmanCode, request.CustID, "resume_at", nil, request.Date, request.Id)
	}
}

func (service *VisitServiceImpl) OutletVisit(ctx context.Context, request request.OutletVisitRequest) {
	err := service.Validate.Struct(request)
	helper.ErrorPanic(err)

	service.OutletVisitListRepo.CreateOutletVisit(ctx, request.EmpID, request.CustID, request.Date)
}

func (service *VisitServiceImpl) GetAllOutletBySalesCode(ctx context.Context, salesCode, custId, date, routeCode string) (responses []response.OutletResponse) {
	result := service.RouteOutletRepo.GetAllOutletBySalesCode(ctx, salesCode, custId, date, routeCode)

	for _, row := range result {
		var res response.OutletResponse
		helper.Automapper(row, &res)
		res.Status = row.Status
		responses = append(responses, res)
	}

	return responses
}

func (service *VisitServiceImpl) SummaryVisit(ctx context.Context, dataFilter entity.SummaryQueryFilter) (res response.SummaryResponse) {
	if dataFilter.CustID == "" {
		dataFilter.CustID, _ = service.CustomerRepo.GetCustIdByEmpCode(ctx, dataFilter.SalesmanCode)
	}
	isPrinciple, errCheck := service.CustomerRepo.CheckIsPrinciple(ctx, dataFilter.CustID)
	if errCheck != nil {
		switch {
		case errors.Is(errCheck, gorm.ErrRecordNotFound):
		default:
			helper.ErrorPanic(errCheck)
		}
	}
	dataFilter.IsPrinciple = isPrinciple

	result, err := service.OutletVisitListRepo.GetSummary(ctx, dataFilter)
	if err != nil {
		helper.ErrorPanic(err)
	}

	if result.Start != nil {
		res.StartTime = *result.Start
	}

	if result.Finish != nil {
		res.EndTime = *result.Finish
	}

	return res
}

func (service *VisitServiceImpl) SummaryVisitStatus(ctx context.Context, dataFilter entity.SummaryQueryFilter) (res response.VisitStatusResponse) {
	if dataFilter.CustID == "" {
		dataFilter.CustID, _ = service.CustomerRepo.GetCustIdByEmpCode(ctx, dataFilter.SalesmanCode)
	}
	isPrinciple, errCheck := service.CustomerRepo.CheckIsPrinciple(ctx, dataFilter.CustID)
	if errCheck != nil {
		switch {
		case errors.Is(errCheck, gorm.ErrRecordNotFound):
		default:
			helper.ErrorPanic(errCheck)
		}
	}
	dataFilter.IsPrinciple = isPrinciple
	result, err := service.OutletVisitListRepo.GetSummaryStatus(ctx, dataFilter)
	if err != nil {
		helper.ErrorPanic(err)
	}

	res = response.VisitStatusResponse{
		Planned:    result.Planned,
		Finished:   result.Finished,
		Skipped:    result.Skipped,
		OnProgress: result.OnProgress,
		OnHold:     result.OnHold,
		ExtraCall:  result.ExtraCall,
	}

	return res
}

func (service *VisitServiceImpl) TravelList(ctx context.Context, outleVisitId int, customerID string) response.TodoListResponse {
	isPrinciple, errCheck := service.CustomerRepo.CheckIsPrinciple(ctx, customerID)
	if errCheck != nil {
		switch {
		case errors.Is(errCheck, gorm.ErrRecordNotFound):
		default:
			helper.ErrorPanic(errCheck)
		}
	}
	var err error
	var data response.TodoListResponse
	if isPrinciple {
		data, err = service.OutletVisitListPrincipleRepo.FindById(ctx, service.Db, outleVisitId)
		if err != nil {
			helper.ErrorPanic(err)
		}
	} else {
		data, err = service.OutletVisitListRepo.FindById(ctx, outleVisitId)
		if err != nil {
			helper.ErrorPanic(err)
		}
	}

	return data

}

func (service *VisitServiceImpl) GetVisitOutletList(ctx context.Context, dataFilter entity.SummaryQueryFilter) (resp []response.OutletVisitListResponse) {
	isPrinciple, errCheck := service.CustomerRepo.CheckIsPrinciple(ctx, dataFilter.CustID)
	if errCheck != nil {
		switch {
		case errors.Is(errCheck, gorm.ErrRecordNotFound):
		default:
			helper.ErrorPanic(errCheck)
		}
	}
	dataFilter.IsPrinciple = isPrinciple

	data, err := service.OutletVisitListRepo.FindAll(ctx, dataFilter)
	if err != nil {
		helper.ErrorPanic(err)
	}

	if len(data) == 0 {
		return resp
	}

	for _, d := range data {
		currentDate := time.Now()
		newDate := currentDate.AddDate(0, 0, d.Top)
		formattedDate := newDate.Format("2006-01-02")
		resp = append(resp, response.OutletVisitListResponse{
			ID:              d.ID,
			Year:            d.Year,
			Week:            d.Week,
			Date:            d.Date,
			Day:             d.Day,
			RouteCode:       d.RouteCode,
			OutletID:        d.OutletID,
			OutletCode:      d.OutletCode,
			OutletName:      d.OutletName,
			OutletAddress:   d.OutletAddress,
			OutletLongitude: d.OutletLongitude,
			OutletLatitude:  d.OutletLatitude,
			DueDate:         formattedDate,
			PjpID:           d.PjpID,
			PjpCode:         d.PjpCode,
			Start:           d.Start,
			Finish:          d.Finish,
			SkipAt:          d.SkipAt,
			LeaveAt:         d.LeaveAt,
			ArriveAt:        d.ArriveAt,
			OnHold:          d.OnHold,
			ResumeAt:        d.ResumeAt,
			CreatedAt:       d.CreatedAt,
			UpdatedAt:       d.UpdatedAt,
			Status:          d.Status,
			IsPlanned:       d.IsPlanned,
			DestinationType: &d.DestinationType,
			SkipReason:      d.SkipReason,
			InOutlet:        d.SkipInOutlet,
		})
	}

	//helper.Automapper(data, &resp)
	return resp
}

func (service *VisitServiceImpl) GetSalesmanReport(ctx context.Context, dataFilter entity.SalesmanReportQueryFilter) entity.SalesmanReport {
	var report entity.SalesmanReport

	// Get visit data
	custID, errCheck := service.CustomerRepo.GetCustIdByEmpId(ctx, dataFilter.SalesmanId)
	if errCheck != nil {
		helper.ErrorPanic(errCheck)
	}

	isPrinciple, errCheck := service.CustomerRepo.CheckIsPrinciple(ctx, custID)
	if errCheck != nil {
		switch {
		case errors.Is(errCheck, gorm.ErrRecordNotFound):
		default:
			helper.ErrorPanic(errCheck)
		}
	}
	skipReasonCount := make(map[string]int)
	if isPrinciple {
		visits, err := service.OutletVisitListPrincipleRepo.GetVisitsByDateAndSalesman(ctx, service.Db, dataFilter)
		if err != nil {
			helper.ErrorPanic(err)
		}

		report.TotalOutlets = len(visits)

		for _, visit := range visits {
			var isPlanned bool
			destinationCode := fmt.Sprintf("%d", visit.RouteCode)
			existsInAdditional, err := service.OutletVisitListPrincipleRepo.CheckRouteOutletAdditional(ctx, service.Db, destinationCode, visit.OutletID)
			if err != nil {
				helper.ErrorPanic(err)
			}
			isPlanned = !existsInAdditional

			if visit.LeaveAt != nil {
				report.TotalVisit++
				if isPlanned {
					report.Visit.Planned++
				} else {
					report.Visit.NotPlanned++
				}
			} else {
				report.TotalNotVisit++
				if isPlanned {
					report.NotVisit.Planned++
				} else {
					report.NotVisit.NotPlanned++
				}
			}

			if visit.SkipAt != nil && visit.SkipReason != nil && *visit.SkipReason != "" {
				skipReasonCount[*visit.SkipReason]++
			}
		}
	} else {
		visits, err := service.OutletVisitListRepo.GetVisitsByDateAndSalesman(ctx, dataFilter)
		if err != nil {
			helper.ErrorPanic(err)
		}
		report.TotalOutlets = len(visits)

		for _, visit := range visits {
			var isPlanned bool
			existsInAdditional, err := service.OutletVisitListRepo.CheckRouteOutletAdditional(ctx, visit.RouteCode, visit.OutletID)
			if err != nil {
				helper.ErrorPanic(err)
			}
			isPlanned = !existsInAdditional

			if visit.LeaveAt != nil {
				report.TotalVisit++
				if isPlanned {
					report.Visit.Planned++
				} else {
					report.Visit.NotPlanned++
				}
			} else {
				report.TotalNotVisit++
				if isPlanned {
					report.NotVisit.Planned++
				} else {
					report.NotVisit.NotPlanned++
				}
			}

			if visit.SkipAt != nil && visit.SkipReason != nil && *visit.SkipReason != "" {
				skipReasonCount[*visit.SkipReason]++
			}
		}
	}

	totalNotVisited := report.TotalNotVisit
	for reason, count := range skipReasonCount {
		percentage := int(math.Round((float64(count) / float64(totalNotVisited)) * 100))
		report.NotVisitReasons = append(report.NotVisitReasons, entity.NotVisitReason{
			Reason:     reason,
			Count:      count,
			Percentage: percentage,
		})
	}

	return report
}
