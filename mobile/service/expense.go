package service

import (
	"context"
	"errors"
	"fmt"
	"mime/multipart"
	"mobile/adapter"
	"mobile/entity"
	"mobile/model"
	"mobile/pkg/config/env"
	"mobile/pkg/constant"
	"mobile/repository"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

var (
	ErrFileLimitExceeded      = errors.New("file must not exceed 3")
	ErrExpenseTypeNotFound    = errors.New("expense type not found")
	ErrOutletNotFound         = errors.New("data outlet not found")
	ErrExpenseNotFound        = errors.New("expense not found")
	ErrDuplicateDeleteFileIDs = errors.New("duplicate delete_file_ids")
	ErrFileNotFoundOrMismatch = errors.New("file not found or does not belong to expense")
)

type ExpenseService interface {
	List(dataFilter entity.ExpenseQueryFilter, custId string, userId int64) (data entity.ExpenseListDataResponse, total int64, lastPage int, err error)
	Detail(expenseId int64, custId string, userId int64) (response entity.ExpenseDetailResponse, err error)
	ExpenseTypeLookup(dataFilter entity.ExpenseTypeQueryFilter) (data []entity.ExpenseTypeLookupResponse, total int64, lastPage int, err error)
	OutletLookupByPJP(dataFilter entity.OutletLookupQueryFilter, salesmanId int64, custId string) (data []entity.OutletLookupResponse, total int64, lastPage int, err error)
	Create(ctx context.Context, request entity.CreateExpenseBody, files []*multipart.FileHeader, userId int64) (response entity.ExpenseDetailResponse, err error)
	Update(expenseId int64, request entity.UpdateExpenseBody, files []*multipart.FileHeader, userId int64, custId string) (response entity.ExpenseDetailResponse, err error)
	Delete(expenseId int64, custId string, userId int64) error
}

func NewExpenseService(
	config env.ConfigEnv,
	expenseRepository repository.ExpenseRepository,
	transaction repository.Dbtransaction,
	obsAdapter adapter.ObsAdapter,
	mOuteletRepository repository.MOutletRepository,
) *expenseServiceImpl {
	return &expenseServiceImpl{
		Config:            config,
		ExpenseRepository: expenseRepository,
		Transaction:       transaction,
		ObsAdapter:        obsAdapter,
		MOutletRepository: mOuteletRepository,
	}
}

type expenseServiceImpl struct {
	Config            env.ConfigEnv
	ExpenseRepository repository.ExpenseRepository
	Transaction       repository.Dbtransaction
	ObsAdapter        adapter.ObsAdapter
	MOutletRepository repository.MOutletRepository
}

func (service *expenseServiceImpl) List(dataFilter entity.ExpenseQueryFilter, custId string, userId int64) (entity.ExpenseListDataResponse, int64, int, error) {
	var response entity.ExpenseListDataResponse
	ctx := context.Background()

	expenses, total, lastPage, err := service.ExpenseRepository.FindAllByCustId(ctx, dataFilter, custId)
	if err != nil {
		return response, 0, 0, err
	}

	isClockOut := false
	attendanceType, err := service.ExpenseRepository.FindAttendanceByUserAndDate(ctx, userId, custId, time.Now())
	if err == nil {
		if attendanceType == 2 {
			isClockOut = true
		}
	}

	response.IsClockOut = isClockOut
	response.ExpenseData = make([]entity.ExpenseListResponse, 0)

	for _, exp := range expenses {
		expenseResp := entity.ExpenseListResponse{
			ExpenseID:   fmt.Sprintf("%d", exp.ExpenseID),
			Date:        exp.Date.Format("02/01/2006"), // DD/MM/YYYY
			ExpenseName: "",
			Amount:      exp.Amount,
			Reason:      "",
		}

		if exp.DocNo != nil {
			expenseResp.DocNo = *exp.DocNo
		}

		if exp.ExpenseTypeName != nil {
			expenseResp.ExpenseName = *exp.ExpenseTypeName
		}
		if exp.Note != nil {
			expenseResp.Reason = *exp.Note
		}

		response.ExpenseData = append(response.ExpenseData, expenseResp)
	}

	return response, total, lastPage, nil
}

func (service *expenseServiceImpl) Detail(expenseId int64, custId string, userId int64) (entity.ExpenseDetailResponse, error) {
	var response entity.ExpenseDetailResponse
	ctx := context.Background()

	expense, details, files, err := service.ExpenseRepository.FindOneByExpenseId(ctx, expenseId, custId)
	if err != nil {
		return response, err
	}

	response.ExpenseID = fmt.Sprintf("%d", expense.ExpenseID)
	response.Date = expense.Date.Format("02/01/2006") // DD/MM/YYYY
	response.Amount = expense.Amount

	if expense.DocNo != nil {
		response.DocNo = *expense.DocNo
	}

	// Check attendance status for is_clock_out based on expense date
	isClockOut := false
	attendanceType, err := service.ExpenseRepository.FindAttendanceByUserAndDate(ctx, userId, custId, expense.Date)
	if err == nil {
		if attendanceType == 2 {
			isClockOut = true
		}
	}
	response.IsClockOut = &isClockOut

	if expense.ExpenseTypeName != nil {
		response.ExpenseName = *expense.ExpenseTypeName
	}
	if expense.Note != nil {
		response.Reason = *expense.Note
	}

	// Map visits (outlets)
	response.Visits = make([]entity.ExpenseVisitResponse, 0, len(details))
	for _, det := range details {
		visit := entity.ExpenseVisitResponse{
			OutletID:       fmt.Sprintf("%d", det.OutletID),
			OutletCode:     det.OutletCode,
			OutletName:     det.OutletName,
			OutletAddress1: det.OutletAddress1,
		}
		response.Visits = append(response.Visits, visit)
	}

	// Map files
	response.File = make([]entity.ExpenseFileResponse, 0, len(files))
	for _, file := range files {
		fileResp := entity.ExpenseFileResponse{
			FileID:        fmt.Sprintf("%d", file.ExpenseFileID),
			FileName:      file.FileName,
			FileKey:       file.FileKey,
			FileURL:       file.FileURL,
			FileSize:      file.FileSize,
			MediaCategory: "",
		}

		// Extract file type from extension
		ext := strings.ToUpper(strings.TrimPrefix(filepath.Ext(file.FileName), "."))
		if ext == "" {
			ext = "UNKNOWN"
		}
		fileResp.FileType = ext

		if file.MediaCategory != nil {
			fileResp.MediaCategory = *file.MediaCategory
		}

		response.File = append(response.File, fileResp)
	}

	return response, nil
}

func (service *expenseServiceImpl) ExpenseTypeLookup(dataFilter entity.ExpenseTypeQueryFilter) ([]entity.ExpenseTypeLookupResponse, int64, int, error) {
	ctx := context.Background()

	expenseTypes, total, lastPage, err := service.ExpenseRepository.FindAllExpenseTypeLookup(ctx, dataFilter)
	if err != nil {
		return nil, 0, 0, err
	}

	response := make([]entity.ExpenseTypeLookupResponse, 0, len(expenseTypes))
	for _, et := range expenseTypes {
		etResp := entity.ExpenseTypeLookupResponse{
			ExpenseTypeID:   et.ExpenseTypeID,
			ExpenseTypeCode: "",
			ExpenseTypeName: "",
		}

		if et.ExpenseTypeCode != nil {
			etResp.ExpenseTypeCode = *et.ExpenseTypeCode
		}
		if et.ExpenseTypeName != nil {
			etResp.ExpenseTypeName = *et.ExpenseTypeName
		}

		response = append(response, etResp)
	}

	return response, total, lastPage, nil
}

func (service *expenseServiceImpl) OutletLookupByPJP(dataFilter entity.OutletLookupQueryFilter, salesmanId int64, custId string) ([]entity.OutletLookupResponse, int64, int, error) {
	ctx := context.Background()

	if len(dataFilter.Statuses) == 0 {
		dataFilter.Statuses = []int{
			constant.NewOutletStatus,
			constant.DormantOutletStatus,
			constant.RegisteredOutletStatus,
			constant.ActiveOutletStatus,
		}
	}

	// outlets, total, lastPage, err := service.ExpenseRepository.FindAllOutletByPJP(ctx, salesmanId, custId, dataFilter)
	outlets, total, lastPage, err := service.MOutletRepository.GetAllOutletLookupByCusttomerId(ctx, dataFilter, custId)
	if err != nil {
		return nil, 0, 0, err
	}

	response := make([]entity.OutletLookupResponse, 0, len(outlets))
	for _, outlet := range outlets {
		outletResp := entity.OutletLookupResponse{
			OutletID:      fmt.Sprintf("%d", outlet.OutletID),
			OutletCode:    outlet.OutletCode,
			OutletName:    outlet.OutletName,
			Address:       outlet.Address,
			Latitude:      outlet.Latitude,
			Longitude:     outlet.Longitude,
			OutletStatus:  outlet.OutletStatus,
			DistributorID: outlet.DistributorID,
			RegionID:      outlet.RegionID,
			AreaID:        outlet.AreaID,
		}
		response = append(response, outletResp)
	}

	return response, total, lastPage, nil
}

func (service *expenseServiceImpl) Create(ctx context.Context, request entity.CreateExpenseBody, files []*multipart.FileHeader, userId int64) (entity.ExpenseDetailResponse, error) {
	var (
		response     entity.ExpenseDetailResponse
		err          error
		errs         []error
		sourceMobile = 2
		collectorID  = int(request.EmpID)
		now          = time.Now()
	)

	// Validate max 3 files
	if len(files) > 3 {
		errs = append(errs, ErrFileLimitExceeded)
	}

	// Validate expense type exists
	_, err = service.ExpenseRepository.FindExpenseTypeById(ctx, request.ExpenseTypeID)
	if err != nil {
		errs = append(errs, ErrExpenseTypeNotFound)
	}

	// Validate outlets exist (only if outlet_id is provided)
	if len(request.OutletID) > 0 {
		validOutletIds, err := service.ExpenseRepository.ValidateOutlets(ctx, request.OutletID, request.CustID)
		if err != nil {
			return response, errors.New("failed to validate outlets")
		}
		if len(validOutletIds) != len(request.OutletID) {
			errs = append(errs, ErrOutletNotFound)
		}
	}

	if len(errs) > 0 {
		return response, errors.Join(errs...)
	}

	// Transaction: insert expense + details + files
	expenseID := int64(0)
	err = service.Transaction.WithinTransaction(ctx, func(txCtx context.Context) error {
		// Generate DocNo
		docNo, err := service.generateDocNo(txCtx, request.CustID, now)
		if err != nil {
			return fmt.Errorf("failed to generate doc no: %w", err)
		}

		expenseModel := model.Expense{
			CustID:        request.CustID,
			ExpenseTypeID: request.ExpenseTypeID,
			DocNo:         &docNo,
			Date:          now,
			Amount:        request.Amount,
			Balance:       request.Amount,
			CreatedBy:     int(userId),
			IsDel:         false,
			Source:        &sourceMobile,
			CollectorID:   &collectorID,
		}

		// Set Note only if provided
		if request.Note != "" {
			expenseModel.Note = &request.Note
		}

		// Insert expense header
		if err := service.ExpenseRepository.Store(txCtx, &expenseModel); err != nil {
			return err
		}

		// Insert expense details (outlets)
		for _, outletId := range request.OutletID {
			detailModel := model.ExpenseDet{
				CustID:    request.CustID,
				ExpenseID: expenseModel.ExpenseID,
				OutletID:  outletId,
			}
			if err := service.ExpenseRepository.StoreDetail(txCtx, &detailModel); err != nil {
				return err
			}
		}

		// Upload files and insert expense_file
		for _, fileHeader := range files {
			fileModel, err := service.processExpenseFile(fileHeader, request.Folder, expenseModel.ExpenseID, request.CustID)
			if err != nil {
				return err
			}

			if err := service.ExpenseRepository.StoreFile(txCtx, fileModel); err != nil {
				return err
			}
		}

		expenseID = expenseModel.ExpenseID
		return nil
	})

	if err != nil {
		return response, err
	}

	// Return detail response
	return service.Detail(expenseID, request.CustID, userId)
}

func (service *expenseServiceImpl) Update(expenseId int64, request entity.UpdateExpenseBody, files []*multipart.FileHeader, userId int64, custId string) (entity.ExpenseDetailResponse, error) {
	var (
		response    entity.ExpenseDetailResponse
		ctx         = context.Background()
		collectorID = int(request.EmpID)
	)

	// Validate expense exists
	_, _, _, err := service.ExpenseRepository.FindOneByExpenseId(ctx, expenseId, custId)
	if err != nil {
		return response, ErrExpenseNotFound
	}

	// Get existing files
	existingFiles, err := service.ExpenseRepository.FindFilesByExpenseId(ctx, expenseId, custId)
	if err != nil {
		return response, err
	}

	var errs []error

	// Validate delete_file_ids if provided
	if len(request.DeleteFileIDs) > 0 {
		// Check for duplicates
		deleteFileIdsMap := make(map[int64]bool)
		for _, id := range request.DeleteFileIDs {
			if deleteFileIdsMap[id] {
				errs = append(errs, ErrDuplicateDeleteFileIDs)
				break // Stop checking duplicates to avoid spamming error
			}
			deleteFileIdsMap[id] = true
		}

		// Validate files exist and belong to expense
		filesToDelete, err := service.ExpenseRepository.FindFilesByExpenseFileIds(ctx, expenseId, request.DeleteFileIDs, custId)
		if err != nil {
			return response, errors.New("failed to validate delete_file_ids")
		}
		if len(filesToDelete) != len(request.DeleteFileIDs) {
			errs = append(errs, ErrFileNotFoundOrMismatch)
		}
	}

	// Validate max 3 files (existing - delete + new)
	filesAfterDelete := len(existingFiles) - len(request.DeleteFileIDs)
	if filesAfterDelete+len(files) > 3 {
		errs = append(errs, ErrFileLimitExceeded)
	}

	// Validate expense type exists
	_, err = service.ExpenseRepository.FindExpenseTypeById(ctx, request.ExpenseTypeID)
	if err != nil {
		errs = append(errs, ErrExpenseTypeNotFound)
	}

	// Validate outlets exist (only if outlet_id is provided)
	if len(request.OutletID) > 0 {
		validOutletIds, err := service.ExpenseRepository.ValidateOutlets(ctx, request.OutletID, custId)
		if err != nil {
			return response, errors.New("failed to validate outlets")
		}
		if len(validOutletIds) != len(request.OutletID) {
			errs = append(errs, ErrOutletNotFound)
		}
	}

	if len(errs) > 0 {
		return response, errors.Join(errs...)
	}

	// Update expense model
	now := time.Now()
	expenseModel := model.Expense{
		ExpenseTypeID: request.ExpenseTypeID,
		Amount:        request.Amount,
		Balance:       request.Amount,
		CollectorID:   &collectorID,
		UpdatedBy:     &userId,
		UpdatedAt:     &now,
	}

	// Set Note only if provided
	if request.Note != "" {
		expenseModel.Note = &request.Note
	}

	// Transaction: update expense + delete old details + insert new details + delete files + insert new files
	err = service.Transaction.WithinTransaction(ctx, func(txCtx context.Context) error {
		// Update expense header
		if err := service.ExpenseRepository.Update(txCtx, expenseId, custId, &expenseModel); err != nil {
			return err
		}

		// Delete old details
		if err := service.ExpenseRepository.DeleteDetailsByExpenseId(txCtx, expenseId, custId); err != nil {
			return err
		}

		// Insert new details
		for _, outletId := range request.OutletID {
			detailModel := model.ExpenseDet{
				CustID:    custId,
				ExpenseID: expenseId,
				OutletID:  outletId,
			}
			if err := service.ExpenseRepository.StoreDetail(txCtx, &detailModel); err != nil {
				return err
			}
		}

		// Delete files if provided
		if len(request.DeleteFileIDs) > 0 {
			if err := service.ExpenseRepository.DeleteFilesByExpenseFileIds(txCtx, request.DeleteFileIDs, custId); err != nil {
				return err
			}
		}

		// Upload new files and insert expense_file
		for _, fileHeader := range files {
			fileModel, err := service.processExpenseFile(fileHeader, request.Folder, expenseId, custId)
			if err != nil {
				return err
			}

			if err := service.ExpenseRepository.StoreFile(txCtx, fileModel); err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return response, err
	}

	// Return detail response
	return service.Detail(expenseId, custId, userId)
}

func (service *expenseServiceImpl) Delete(expenseId int64, custId string, userId int64) error {
	ctx := context.Background()

	// Validate expense exists
	_, _, _, err := service.ExpenseRepository.FindOneByExpenseId(ctx, expenseId, custId)
	if err != nil {
		return ErrExpenseNotFound
	}

	// Soft delete
	return service.ExpenseRepository.Delete(ctx, expenseId, custId, userId)
}

func (service *expenseServiceImpl) uploadExpenseFile(fileHeader *multipart.FileHeader, folder string) (string, error) {
	if fileHeader == nil {
		return "", errors.New("file is required")
	}
	if service.ObsAdapter == nil {
		return "", errors.New("file uploader is not configured")
	}

	uploadModel := &model.Upload{
		Folder: folder,
		File:   fileHeader,
	}

	return service.ObsAdapter.UploadFile(uploadModel)
}

// determineFileMetadata determines file_key and media_category from file header
func (service *expenseServiceImpl) determineFileMetadata(fileHeader *multipart.FileHeader, folder string) (fileKey string, mediaCategory string) {
	uploadModel := &model.Upload{
		Folder: folder,
		File:   fileHeader,
	}
	contentType := uploadModel.GetFileContentType()
	fileKey = "image"
	if strings.Contains(contentType, "video") {
		fileKey = "video"
	}

	// Extract media category from file extension
	mediaCategory = "Receipt"
	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
	if ext == ".mp4" || ext == ".mov" || ext == ".avi" {
		mediaCategory = "Video"
	}

	return fileKey, mediaCategory
}

// processExpenseFile uploads file and returns ExpenseFile model ready to be stored
func (service *expenseServiceImpl) processExpenseFile(fileHeader *multipart.FileHeader, folder string, expenseId int64, custId string) (*model.ExpenseFile, error) {
	// Upload file
	fileURL, err := service.uploadExpenseFile(fileHeader, folder)
	if err != nil {
		return nil, fmt.Errorf("failed to upload file: %w", err)
	}

	// Determine file metadata
	fileKey, mediaCategory := service.determineFileMetadata(fileHeader, folder)

	// Create expense_file model
	fileModel := &model.ExpenseFile{
		CustID:        custId,
		ExpenseID:     expenseId,
		FileName:      fileHeader.Filename,
		FileURL:       fileURL,
		FileKey:       fileKey,
		MediaCategory: &mediaCategory,
		FileSize:      fileHeader.Size,
	}

	return fileModel, nil
}

func (service *expenseServiceImpl) generateDocNo(ctx context.Context, custId string, date time.Time) (string, error) {
	// Format: EYYYYMMDDNNN
	prefix := fmt.Sprintf("E%s", date.Format("20060102"))

	lastDocNo, err := service.ExpenseRepository.GetLastDocNo(ctx, custId, date)
	if err != nil {
		return "", err
	}

	if lastDocNo == "" {
		return prefix + "001", nil
	}

	// Extract running number (last 3 digits)
	if len(lastDocNo) < 3 {
		return prefix + "001", nil
	}

	runningNumStr := lastDocNo[len(lastDocNo)-3:]
	runningNum, err := strconv.Atoi(runningNumStr)
	if err != nil {
		// If failed to parse, start from 001
		return prefix + "001", nil
	}

	// Increment
	newRunningNum := runningNum + 1
	return fmt.Sprintf("%s%03d", prefix, newRunningNum), nil
}
