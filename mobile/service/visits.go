package service

import (
	"context"
	"errors"
	"fmt"
	"mobile/adapter"
	"mobile/entity"
	"mobile/model"
	"mobile/pkg/config/env"
	"mobile/pkg/constant"
	"mobile/pkg/str"
	"mobile/pkg/structs"
	"mobile/repository"
	"path/filepath"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2/log"
)

type VisitsService interface {
	Visits(entity.VisitsListRequest) (entity.VisitListDetailResponse, error)
	Summaries(ctx context.Context, params entity.SummariesRequest) (*entity.SummariesResponse, error)
	List(entity.VisitQueryFilter) ([]entity.VisitListResponse, error)
	Start(req entity.StartRequest) (err error)
	Skip(req entity.SkipRequest) (err error)
	SkipReasons(entity.SkipReasonsQueryFilter) ([]entity.SkipReasonsResponse, error)
	Arrive(req entity.ArriveRequest) (resp entity.ArriveResponse, err error)
	Hold(req entity.HoldRequest) (err error)
	Resume(req entity.ResumeRequest) (err error)
	Leave(req entity.LeaveRequest) (err error)
	End(req entity.EndRequest) (err error)
}

type VisitsServiceImpl struct {
	Config env.ConfigEnv
	// MCustomerRepository repository.MCustomerRepository,
	Transaction              repository.Dbtransaction
	MEmployeeRepository      repository.MEmployeeRepository
	VisitRepo                repository.VisitsRepository
	ObsAdapter               adapter.ObsAdapter
	InvoicesRepo             repository.InvoicesRepository
	PjpDistributorRepository repository.PjpDistributorRepository
	PjpPrincipleRepository   repository.PjpPrincipalRepository
}

func NewVisitsService(
	config env.ConfigEnv,
	transaction repository.Dbtransaction,
	mEmployee repository.MEmployeeRepository,
	visitRepo repository.VisitsRepository,
	obsAdapter adapter.ObsAdapter,
	invoicesRepo repository.InvoicesRepository,
	pjpDistributorRepository repository.PjpDistributorRepository,
	pjpPrincipalRepository repository.PjpPrincipalRepository,
) *VisitsServiceImpl {
	return &VisitsServiceImpl{
		Config:                   config,
		Transaction:              transaction,
		MEmployeeRepository:      mEmployee,
		VisitRepo:                visitRepo,
		ObsAdapter:               obsAdapter,
		InvoicesRepo:             invoicesRepo,
		PjpDistributorRepository: pjpDistributorRepository,
		PjpPrincipleRepository:   pjpPrincipalRepository,
	}
}

func (service *VisitsServiceImpl) Visits(request entity.VisitsListRequest) (resp entity.VisitListDetailResponse, err error) {
	pjp := model.PermanentJourneyPlan{}
	if request.IsDistributor {
		pjp, err = service.PjpDistributorRepository.GetPJPInfo(context.Background(), request.EmpID)
		if err != nil {
			log.Error("PjpService, GetSalesmanDetail, PjpDistributorRepository.GetPJPInfo, err:", err.Error())
			return resp, err
		}
	} else {
		pjp, err = service.PjpPrincipleRepository.GetPJPInfo(context.Background(), request.EmpID)
		if err != nil {
			log.Error("PjpService, GetSalesmanDetail, PjpDistributorRepository.GetPJPInfo, err:", err.Error())
			return resp, err
		}
	}

	request.PJPID = pjp.ID
	visits, err := service.VisitRepo.GetVisitListByCustID(request)
	if err != nil {
		return resp, err
	}

	if len(visits) == 0 {
		return resp, errors.New(constant.STATUS_DB_NOT_FOUND)
	}

	// Get the latest visit (first one after ordering by arrive_at DESC)
	latestVisit := visits[0]

	// Format checkin_at from arrive_at (unix timestamp in milliseconds)
	var checkinAt string
	if latestVisit.ArriveAt != nil {
		// arrive_at is stored in milliseconds, convert to seconds
		arriveAtSeconds := *latestVisit.ArriveAt / 1000
		checkinTime := time.Unix(arriveAtSeconds, 0).UTC()
		checkinAt = checkinTime.Format(time.RFC3339Nano)
	}

	// Get outlet info
	outletCode := ""
	if latestVisit.OutletCode != nil {
		outletCode = *latestVisit.OutletCode
	}

	outletName := ""
	if latestVisit.OutletName != nil {
		outletName = *latestVisit.OutletName
	}

	// Get coordinates
	latitude := ""
	if latestVisit.Latitude != nil {
		latitude = *latestVisit.Latitude
	}

	longitude := ""
	if latestVisit.Longitude != nil {
		longitude = *latestVisit.Longitude
	}

	// Build file info
	var fileInfo *entity.FileInfo
	if latestVisit.PhotoPath != nil && *latestVisit.PhotoPath != "" {
		fileInfo = service.buildFileInfo(*latestVisit.PhotoPath, latestVisit.Folder)
	}

	resp = entity.VisitListDetailResponse{
		CheckinAt:  checkinAt,
		OutletCode: outletCode,
		OutletName: outletName,
		Longitude:  longitude,
		Latitude:   latitude,
		File:       fileInfo,
	}

	return resp, nil
}

func (service *VisitsServiceImpl) buildFileInfo(photoPath string, folder *string) *entity.FileInfo {
	// photoPath is full URL like "https://bucket.endpoint/folder/filename.ext"
	// Extract file key (relative path) from URL
	var fileKey string
	if folder != nil && *folder != "" {
		// If folder is provided, construct file_key from folder + filename
		fileName := filepath.Base(photoPath)
		fileKey = *folder + "/" + fileName
	} else {
		// Extract key from URL (remove base URL part)
		// photoPath format: "https://bucket.endpoint/folder/filename.ext"
		// We need: "folder/filename.ext"
		parts := strings.Split(photoPath, "/")
		if len(parts) > 3 {
			// Skip protocol, empty, bucket.endpoint, then join the rest
			fileKey = strings.Join(parts[3:], "/")
		} else {
			// Fallback: use full path
			fileKey = photoPath
		}
	}

	// Extract file name from path
	fileName := filepath.Base(photoPath)
	// Remove extension for file_name (as per example: "arrival_IMG_20250119_155922")
	ext := filepath.Ext(fileName)
	fileNameWithoutExt := strings.TrimSuffix(fileName, ext)

	// Detect media category from extension
	mediaCategory := "image"
	extLower := strings.ToLower(ext)
	if strings.Contains(extLower, "mp4") || strings.Contains(extLower, "mov") || strings.Contains(extLower, "avi") {
		mediaCategory = "video"
	}

	// Use photo_path as file_url (since files are public with ACL PublicRead)
	fileURL := photoPath

	return &entity.FileInfo{
		FileName:      fileNameWithoutExt,
		FileURL:       fileURL,
		FileKey:       fileKey,
		MediaCategory: mediaCategory,
		FileSize:      0, // Not stored in database, set to 0
	}
}

func (service *VisitsServiceImpl) Summaries(ctx context.Context, request entity.SummariesRequest) (response *entity.SummariesResponse, err error) {
	var (
		currentTime = time.Now().UTC()
		todayStart  = time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), 0, 0, 0, 0, time.UTC)
		todayEnd    = time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), 23, 59, 59, 0, time.UTC)
		dateStr     = currentTime.Format(time.DateOnly)
	)

	if request.Date != "" {
		t, err := time.Parse(time.DateOnly, request.Date)
		if err != nil {
			return response, err
		}
		todayStart = time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
		todayEnd = time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, 0, time.UTC)
		dateStr = t.Format(time.DateOnly)
	}

	invoices, err := service.InvoicesRepo.GetInvoiceListByDate(ctx, request.EmpID, dateStr)
	if err != nil {
		return response, err
	}

	expenses, err := service.InvoicesRepo.GetExpenseSummaryByDate(ctx, request.EmpID, dateStr)
	if err != nil {
		return response, err
	}

	// collections, paymentsCollection, err := service.InvoicesRepo.GetCollectionSummaryByDate(ctx, request.EmpID, dateStr)
	// if err != nil {
	// 	return response, err
	// }

	// Ensure we return empty arrays instead of null in JSON if no records
	if invoices == nil {
		invoices = []model.InvoiceListItem{}
	}

	if expenses == nil {
		expenses = []model.ExpenseListItem{}
	}

	// if collections == nil {
	// 	collections = []model.CollectionListItem{}
	// }

	// var (
	// 	paymentCollectionMap     = make(map[string]model.PaymentInvoiceList)
	// 	paymentCollectionInfoMap = make(map[string]model.InvoicePayment)
	// )

	// for _, payment := range paymentsCollection {
	// 	paymentCollectionInfoMap[payment.InvoiceNo] = payment
	// 	paymentCollectionMap[payment.InvoiceNo] = payment.PaymentInvoiceList
	// }

	visits, err := service.VisitRepo.GetLastVisitGroupOutlet(request.CustID, todayStart, todayEnd)
	if err != nil {
		return response, err
	}

	var (
		totalPaymentAmountInvoice = 0.0
		orderNos                  = make([]string, 0)
	)

	// get payment invoice
	for _, invoice := range invoices {
		orderNos = append(orderNos, invoice.OrderNo)
	}

	paymentsInvoices, err := service.InvoicesRepo.GetPaymentsByInvoiceNo(ctx, request.EmpID, orderNos)
	if err != nil {
		return response, err
	}

	paymentInvoiceMap := make(map[string]model.InvoicePayment)
	paymentInvoiceListMap := make(map[string][]model.PaymentInvoiceList)
	for _, ip := range paymentsInvoices {
		// Handle paymentInvoiceMap
		if existing, ok := paymentInvoiceMap[ip.InvoiceNo]; ok {
			existing.PaymentAmount += ip.PaymentAmount

			// update
			paymentInvoiceMap[ip.InvoiceNo] = existing

		} else {
			paymentInvoiceMap[ip.InvoiceNo] = ip
		}

		// Handle paymentInvoiceListMap (unchanged)
		if _, ok := paymentInvoiceListMap[ip.InvoiceNo]; !ok {
			paymentInvoiceListMap[ip.InvoiceNo] = append([]model.PaymentInvoiceList{ip.PaymentInvoiceList})
		} else {
			paymentInvoiceListMap[ip.InvoiceNo] = append(paymentInvoiceListMap[ip.InvoiceNo], ip.PaymentInvoiceList)
		}
	} //end payment invoice

	for i, invoice := range invoices {
		invoices[i].InvoiceAmount = paymentInvoiceMap[invoice.OrderNo].InvoiceAmount
		invoices[i].RemainingAmount = paymentInvoiceMap[invoice.OrderNo].RemainingAmount
		invoices[i].PaidAmount = paymentInvoiceMap[invoice.OrderNo].PaymentAmount
		invoices[i].TotalPayment = paymentInvoiceMap[invoice.OrderNo].PaymentAmount
		invoices[i].PaymentBalance = 0
		invoices[i].RemainingPayment = paymentInvoiceMap[invoice.OrderNo].RemainingAmount
		invoices[i].Payments = paymentInvoiceListMap[invoice.OrderNo]
	}

	for _, payment := range paymentsInvoices {
		totalPaymentAmountInvoice += payment.PaymentAmount
	}

	totalExpense := 0.0
	for _, expense := range expenses {
		totalExpense += expense.Amount
	}

	totalCollection := 0.0
	totalPaymentAmountCollection := 0.0
	// for _, collection := range collections {
	// 	totalCollection += collection.TotalAmount
	// }

	// for _, payment := range paymentsCollection {
	// 	totalPaymentAmountCollection += payment.PaymentAmount
	// }

	totalReceived := totalPaymentAmountInvoice + totalExpense + totalCollection

	response = &entity.SummariesResponse{
		Totals: entity.TotalSummaryDeposit{
			TotalPayment:    totalPaymentAmountInvoice + totalPaymentAmountCollection,
			TotalExpense:    totalExpense,
			TotalCollection: totalCollection,
			TotalReceived:   totalReceived,
		},
		CurrentTime: currentTime.Format(time.RFC3339),
		StartTime:   todayStart.Format(time.RFC3339),
		EndTime:     todayEnd.Format(time.RFC3339),
		InvoiceList: invoices,
		ExpenseList: expenses,
		// CollectionSummary: collections,
	}

	for _, visit := range visits {
		switch *visit.Type {
		case entity.TYPE_ARRIVE_ID:
			response.InProgressIncrement()
		case entity.TYPE_RESUME_ID:
			response.InProgressIncrement()
		case entity.TYPE_ON_HOLD_ID:
			response.OnHoldIncrement()
		case entity.TYPE_SKIP_ID:
			response.SkippedIncrement()
		case entity.TYPE_LEAVE_ID:
			response.FinishIncrement()
		}
	}

	return response, err
}

func (service *VisitsServiceImpl) List(dataFilter entity.VisitQueryFilter) (data []entity.VisitListResponse, err error) {

	outlets, err := service.VisitRepo.FindAllByCustId(dataFilter)
	if err != nil {
		return data, err
	}

	data = make([]entity.VisitListResponse, 0)
	for i, row := range outlets {
		fmt.Println("ROW >>>", row.NoOrderId)
		var vResp entity.VisitListResponse
		structs.Automapper(row, &vResp)
		vResp.Sequence = i + 1
		if row.NoOrderId != 0 {
			vResp.Status = "skipped"
		} else {
			vResp.Status = "planned"
		}
		vResp.OutletImg = row.ImageURL
		data = append(data, vResp)
	}

	return data, err
}

func (service *VisitsServiceImpl) Start(req entity.StartRequest) (err error) {
	timeRequest := time.Unix(req.CurrentTime, 0).UTC()
	employee, err := service.MEmployeeRepository.FindOneByEmailCustID(req.Email, req.CustID)
	if err != nil {
		return err
	}

	visitModel := model.Visit{
		CustID:    req.CustID,
		EmpCode:   &employee.EmpCode,
		Type:      &entity.TYPE_START_ID,
		CreatedAt: timeRequest,
	}
	err = service.VisitRepo.Store(&visitModel)
	if err != nil {
		return err
	}
	return nil
}

func (service *VisitsServiceImpl) Skip(req entity.SkipRequest) (err error) {
	now := time.Now()
	from := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.UTC().Location())
	timeRequest := time.Unix(req.CurrentTime, 0).UTC()
	endRange := timeRequest

	employee, err := service.MEmployeeRepository.FindOneByEmailCustID(req.Email, req.CustID)
	if err != nil {
		return err
	}

	visitStart, err := service.VisitRepo.GetStartByEmployeeBetweenTime(employee.EmpCode, req.CustID, from, timeRequest)
	if err != nil {
		if err.Error() != constant.STATUS_DB_NOT_FOUND {
			return err
		}
	}

	if visitStart == nil {
		return errors.New(entity.ERROR_SKIP_MUST_BE_START)
	}

	if timeRequest.Before(visitStart.CreatedAt) {
		return errors.New(entity.ERROR_ARRIVE_DATE_MUST_GREAT_START)
	}

	lastVisit, err := service.VisitRepo.GetLastStatusByOutletCodeBetweenTime(req.OutletCode, from, endRange)
	if err != nil {
		if err.Error() != constant.STATUS_DB_NOT_FOUND {
			return err
		}
	}

	if lastVisit != nil {
		if *lastVisit.Type == entity.TYPE_SKIP_ID {
			return fmt.Errorf(entity.ERROR_SKIP_ALREADY_EXIST, *lastVisit.OutletCode)
		}
		return fmt.Errorf(entity.ERROR_ARRIVE_MUST_BE_FIRST_STATUS, entity.GetTypeNameByID(*lastVisit.Type))
	}

	// parse time format YYYY-mm-dd to Rfc3339
	upComingDate, err := str.DateStrToRfc3339String(req.UpcomingVisit)
	if err != nil {
		return err
	}
	req.UpcomingVisit = upComingDate

	visitModel := model.Visit{
		CustID:        req.CustID,
		EmpCode:       &employee.EmpCode,
		Type:          &entity.TYPE_SKIP_ID,
		CreatedAt:     timeRequest,
		Latitude:      &req.Latitude,
		Longitude:     &req.Longitude,
		OutletCode:    &req.OutletCode,
		IsInOutlet:    &req.IsInOutlet,
		Reason:        &req.Reason,
		UpComingVisit: &req.UpcomingVisit,
	}
	err = service.VisitRepo.Store(&visitModel)
	if err != nil {
		return err
	}
	return nil
}

func (service *VisitsServiceImpl) SkipReasons(dataFilter entity.SkipReasonsQueryFilter) (data []entity.SkipReasonsResponse, err error) {

	skipReasons, err := service.VisitRepo.FindAllSkipReasonByCustId(dataFilter)
	if err != nil {
		return data, err
	}

	data = make([]entity.SkipReasonsResponse, 0)
	for _, row := range skipReasons {
		var vResp entity.SkipReasonsResponse
		structs.Automapper(row, &vResp)
		vResp.Reason = row.SkipReasonName
		data = append(data, vResp)
	}

	return data, err
}

func (service *VisitsServiceImpl) Arrive(req entity.ArriveRequest) (resp entity.ArriveResponse, err error) {
	timeRequest := time.Unix(req.CurrentTime, 0).UTC()
	from := time.Date(timeRequest.Year(), timeRequest.Month(), timeRequest.Day(), 0, 0, 0, 0, timeRequest.Location())
	endRange := timeRequest

	employee, err := service.MEmployeeRepository.FindOneByEmailCustID(req.Email, req.CustID)
	if err != nil {
		return resp, err
	}

	visitStart, err := service.VisitRepo.GetStartByEmployeeBetweenTime(employee.EmpCode, req.CustID, from, timeRequest)
	if err != nil {
		if err.Error() != constant.STATUS_DB_NOT_FOUND {
			return resp, err
		}
	}

	if visitStart == nil {
		return resp, errors.New(entity.ERROR_ARRIVE_MUST_BE_START)
	}

	if timeRequest.Before(visitStart.CreatedAt) {
		return resp, errors.New(entity.ERROR_ARRIVE_DATE_MUST_GREAT_START)
	}

	lastVisit, err := service.VisitRepo.GetLastStatusByOutletCodeBetweenTime(req.OutletCode, from, endRange)
	if err != nil {
		if err.Error() != constant.STATUS_DB_NOT_FOUND {
			return resp, err
		}
	}

	if lastVisit != nil {
		if lastVisit.Type != nil && *lastVisit.Type == entity.TYPE_ARRIVE_ID {
			return resp, errors.New(entity.ERROR_ARRIVE_ALREADY_EXIST)
		}
		if lastVisit.Type != nil {
			return resp, fmt.Errorf(entity.ERROR_ARRIVE_MUST_BE_FIRST_STATUS, entity.GetTypeNameByID(*lastVisit.Type))
		}
		return resp, errors.New(entity.ERROR_ARRIVE_MUST_BE_FIRST_STATUS)
	}

	outlet, err := service.VisitRepo.FindOutletByCode(req.CustID, req.OutletCode)
	if err != nil {
		return resp, err
	}

	photoURL, err := service.uploadArriveFile(&req)
	if err != nil {
		return resp, err
	}

	visitModel := model.Visit{
		CustID:     req.CustID,
		EmpCode:    &employee.EmpCode,
		Type:       &entity.TYPE_ARRIVE_ID,
		CreatedAt:  timeRequest,
		Latitude:   &req.Latitude,
		Longitude:  &req.Longitude,
		OutletCode: &req.OutletCode,
	}
	err = service.VisitRepo.Store(&visitModel)
	if err != nil {
		return resp, err
	}

	visitList := service.buildOutletVisitListModel(req, outlet.OutletId, employee.EmpCode, photoURL, timeRequest)
	if err = service.VisitRepo.StoreOutletVisitList(&visitList); err != nil {
		return resp, err
	}

	if req.IsUpdateLocation {
		if err = service.VisitRepo.UpdateOutletLocation(outlet.OutletId, req.Latitude, req.Longitude); err != nil {
			return resp, err
		}
	}

	resp.PhotoURL = photoURL
	resp.IsUpdateLocation = req.IsUpdateLocation
	resp.OutletVisitListID = visitList.ID

	return resp, nil
}

func (service *VisitsServiceImpl) uploadArriveFile(req *entity.ArriveRequest) (string, error) {
	if req.File == nil {
		return "", errors.New("file is required")
	}
	if service.ObsAdapter == nil {
		return "", errors.New("file uploader is not configured")
	}

	uploadModel := &model.Upload{
		Folder: req.Folder,
		File:   req.File,
	}

	return service.ObsAdapter.UploadFile(uploadModel)
}

func (service *VisitsServiceImpl) buildOutletVisitListModel(req entity.ArriveRequest, outletID int64, empCode string, photoURL string, timeRequest time.Time) model.OutletVisitList {
	visitDate := time.Date(timeRequest.Year(), timeRequest.Month(), timeRequest.Day(), 0, 0, 0, 0, time.UTC)
	nowUTC := time.Now().UTC()
	arriveUnix := timeRequest.Unix()
	outletIDCopy := outletID
	outletCode := req.OutletCode
	lat := req.Latitude
	lng := req.Longitude
	photo := photoURL
	folder := req.Folder
	isUpdate := req.IsUpdateLocation

	return model.OutletVisitList{
		OutletID:         &outletIDCopy,
		OutletCode:       &outletCode,
		Date:             visitDate,
		ArriveAt:         &arriveUnix,
		Latitude:         &lat,
		Longitude:        &lng,
		PhotoPath:        &photo,
		Folder:           &folder,
		IsUpdateLocation: &isUpdate,
		CreatedAt:        nowUTC,
	}
}

func (service *VisitsServiceImpl) Hold(req entity.HoldRequest) (err error) {
	now := time.Now()
	from := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.UTC().Location())
	to := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, now.UTC().Location())
	timeRequest := time.Unix(req.CurrentTime, 0).UTC()

	lastVisit, err := service.VisitRepo.GetLastStatusByOutletCodeBetweenTime(req.OutletCode, from, to)
	if err != nil {
		if err.Error() != constant.STATUS_DB_NOT_FOUND {
			return err
		}
	}

	if lastVisit == nil {
		return errors.New(entity.ERROR_ON_HOLD_MUST_BE_ARRIVE)
	}

	if lastVisit.Type != nil {
		if *lastVisit.Type != entity.TYPE_ARRIVE_ID {
			return errors.New(entity.ERROR_ON_HOLD_LAST_MUST_BE_ARRIVE)
		}
	}

	if timeRequest.Before(lastVisit.CreatedAt) {
		return errors.New(entity.ERROR_ON_HOLD_DATE_MUST_GREAT_ARRIVE)
	}

	employee, err := service.MEmployeeRepository.FindOneByEmailCustID(req.Email, req.CustID)
	if err != nil {
		return err
	}

	visitModel := model.Visit{
		CustID:     req.CustID,
		EmpCode:    &employee.EmpCode,
		Type:       &entity.TYPE_ON_HOLD_ID,
		CreatedAt:  timeRequest,
		OutletCode: &req.OutletCode,
	}

	err = service.VisitRepo.Store(&visitModel)
	if err != nil {
		return err
	}
	return nil
}

func (service *VisitsServiceImpl) Resume(req entity.ResumeRequest) (err error) {
	now := time.Now()
	from := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.UTC().Location())
	to := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, now.UTC().Location())
	timeRequest := time.Unix(req.CurrentTime, 0).UTC()

	lastVisit, err := service.VisitRepo.GetLastStatusByOutletCodeBetweenTime(req.OutletCode, from, to)
	if err != nil {
		if err.Error() != constant.STATUS_DB_NOT_FOUND {
			return err
		}
	}

	if lastVisit == nil {
		return errors.New(entity.ERROR_RESUME_MUST_BE_ON_HOLD)
	}

	if lastVisit.Type != nil {
		if *lastVisit.Type != entity.TYPE_ON_HOLD_ID {
			return errors.New(entity.ERROR_RESUME_LAST_MUST_BE_ON_HOLD)
		}
	}

	if timeRequest.Before(lastVisit.CreatedAt) {
		return errors.New(entity.ERROR_RESUME_DATE_MUST_GREAT_ON_HOLD)
	}

	employee, err := service.MEmployeeRepository.FindOneByEmailCustID(req.Email, req.CustID)
	if err != nil {
		return err
	}

	visitModel := model.Visit{
		CustID:     req.CustID,
		EmpCode:    &employee.EmpCode,
		Type:       &entity.TYPE_RESUME_ID,
		CreatedAt:  timeRequest,
		OutletCode: &req.OutletCode,
	}

	err = service.VisitRepo.Store(&visitModel)
	if err != nil {
		return err
	}
	return nil
}

func (service *VisitsServiceImpl) Leave(req entity.LeaveRequest) (err error) {
	now := time.Now()
	from := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.UTC().Location())
	to := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, now.UTC().Location())
	timeRequest := time.Unix(req.CurrentTime, 0).UTC()

	lastVisit, err := service.VisitRepo.GetLastStatusByOutletCodeBetweenTime(req.OutletCode, from, to)
	if err != nil {
		if err.Error() != constant.STATUS_DB_NOT_FOUND {
			return err
		}
	}

	if lastVisit == nil {
		return errors.New(entity.ERROR_LEAVE_MUST_BE_ARRIVE_OR_RESUME)
	}
	fmt.Println(*lastVisit.Type)
	if lastVisit.Type != nil {
		if *lastVisit.Type != entity.TYPE_ARRIVE_ID && *lastVisit.Type != entity.TYPE_RESUME_ID {
			return errors.New(entity.ERROR_LEAVE_LAST_MUST_BE_ARRIVE_OR_RESUME)
		}
	}

	if timeRequest.Before(lastVisit.CreatedAt) {
		fmt.Println(timeRequest)
		fmt.Println(lastVisit.CreatedAt)
		return errors.New(entity.ERROR_LEAVE_DATE_MUST_GREAT_ARRIVE_OR_RESUME)
	}

	employee, err := service.MEmployeeRepository.FindOneByEmailCustID(req.Email, req.CustID)
	if err != nil {
		return err
	}

	visitModel := model.Visit{
		CustID:     req.CustID,
		EmpCode:    &employee.EmpCode,
		Type:       &entity.TYPE_LEAVE_ID,
		CreatedAt:  timeRequest,
		OutletCode: &req.OutletCode,
	}

	err = service.VisitRepo.Store(&visitModel)
	if err != nil {
		return err
	}
	return nil
}

func (service *VisitsServiceImpl) End(req entity.EndRequest) (err error) {
	timeRequest := time.Unix(req.CurrentTime, 0).UTC()
	employee, err := service.MEmployeeRepository.FindOneByEmailCustID(req.Email, req.CustID)
	if err != nil {
		return err
	}

	visitModel := model.Visit{
		CustID:    req.CustID,
		EmpCode:   &employee.EmpCode,
		Type:      &entity.TYPE_END_ID,
		CreatedAt: timeRequest,
	}
	err = service.VisitRepo.Store(&visitModel)
	if err != nil {
		return err
	}
	return nil
}
