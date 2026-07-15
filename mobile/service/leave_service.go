package service

import (
	"errors"
	"mime/multipart"
	"mobile/adapter"
	"mobile/entity"
	"mobile/model"
	"mobile/pkg/times"
	"mobile/repository"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"
)

const (
	leaveUploadFolder = "leaves"
	maxLeaveFileSize  = 5 * 1024 * 1024
)

var (
	ErrAttendanceRecorded    = errors.New("Attendance has been recorded")
	ErrLeaveRequestRecorded  = errors.New("Leave request has been recorded")
	ErrInvalidLeaveDateRange = errors.New("start_date must be before or equal to end_date")
	ErrLeaveFileTypeInvalid  = errors.New("file must be jpg or png")
	ErrLeaveFileSizeExceeded = errors.New("file must not exceed 5MB")
	ErrLeaveFileUploader     = errors.New("file uploader is not configured")
)

type LeaveService interface {
	CreateLeaveRequest(request entity.LeaveRequestCreate) error
	ListLeaveRequests(filter entity.LeaveRequestQuery) ([]entity.LeaveRequestItem, int64, int, error)
	LeaveCheck(custID string, empID int64) (*entity.LeaveCheckResponse, error)
}

type LeaveServiceImpl struct {
	LeaveRequestRepo repository.LeaveRequestRepository
	AttendanceRepo   repository.AttendanceRepository
	ObsAdapter       adapter.ObsAdapter
}

func NewLeaveService(
	leaveRequestRepo repository.LeaveRequestRepository,
	attendanceRepo repository.AttendanceRepository,
	obsAdapter adapter.ObsAdapter,
) *LeaveServiceImpl {
	return &LeaveServiceImpl{
		LeaveRequestRepo: leaveRequestRepo,
		AttendanceRepo:   attendanceRepo,
		ObsAdapter:       obsAdapter,
	}
}

func (service *LeaveServiceImpl) CreateLeaveRequest(request entity.LeaveRequestCreate) error {
	startDate, err := time.Parse(time.DateOnly, request.StartDate)
	if err != nil {
		return err
	}
	endDate, err := time.Parse(time.DateOnly, request.EndDate)
	if err != nil {
		return err
	}
	if startDate.After(endDate) {
		return ErrInvalidLeaveDateRange
	}

	var fileURL, fileName string
	if request.File != nil {
		if err := validateLeaveFile(request.File); err != nil {
			return err
		}

		uploadedURL, err := service.uploadLeaveFile(request.File)
		if err != nil {
			return err
		}
		fileURL = uploadedURL
		fileName = request.File.Filename
	}

	hasAttendance, err := service.AttendanceRepo.ExistsAttendanceInDateRange(
		request.CustID,
		request.EmpCode,
		startDate,
		endDate,
	)
	if err != nil {
		return err
	}
	if hasAttendance {
		return ErrAttendanceRecorded
	}

	hasOverlap, err := service.LeaveRequestRepo.ExistsOverlappingLeave(
		request.CustID,
		request.EmpID,
		startDate,
		endDate,
	)
	if err != nil {
		return err
	}
	if hasOverlap {
		return ErrLeaveRequestRecorded
	}

	leaveRequest := &model.LeaveRequest{
		CustID:    request.CustID,
		EmpID:     request.EmpID,
		StartDate: startDate,
		EndDate:   endDate,
		Reason:    request.Reason,
		FileURL:   fileURL,
		FileName:  fileName,
		Approval:  "approved",
		CreatedBy: request.UserID,
	}

	return service.LeaveRequestRepo.Create(leaveRequest)
}

func validateLeaveFile(file *multipart.FileHeader) error {
	if file.Size > maxLeaveFileSize {
		return ErrLeaveFileSizeExceeded
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" {
		return ErrLeaveFileTypeInvalid
	}

	uploadModel := &model.Upload{File: file}
	contentType := strings.ToLower(uploadModel.GetFileContentType())
	if contentType != "" &&
		contentType != "image/jpeg" &&
		contentType != "image/jpg" &&
		contentType != "image/png" {
		return ErrLeaveFileTypeInvalid
	}

	return nil
}

func (service *LeaveServiceImpl) uploadLeaveFile(file *multipart.FileHeader) (string, error) {
	if service.ObsAdapter == nil {
		return "", ErrLeaveFileUploader
	}

	uploadModel := &model.Upload{
		Folder: leaveUploadFolder,
		File:   file,
	}

	return service.ObsAdapter.UploadFile(uploadModel)
}

func (service *LeaveServiceImpl) ListLeaveRequests(filter entity.LeaveRequestQuery) ([]entity.LeaveRequestItem, int64, int, error) {
	rows, total, lastPage, err := service.LeaveRequestRepo.List(filter)
	if err != nil {
		return nil, 0, 0, err
	}

	items := make([]entity.LeaveRequestItem, 0, len(rows))
	for _, row := range rows {
		items = append(items, mapLeaveRequestItem(row))
	}

	return items, total, lastPage, nil
}

func mapLeaveRequestItem(row model.LeaveRequestRead) entity.LeaveRequestItem {
	item := entity.LeaveRequestItem{
		CustID:     row.CustID,
		EmpID:      strconv.FormatInt(row.EmpID, 10),
		StartDate:  row.StartDate.Format(time.DateOnly),
		EndDate:    row.EndDate.Format(time.DateOnly),
		Reason:     row.Reason,
		Approval:   row.Approval,
		Duration:   row.Duration,
		CreatedBy:  row.CreatedBy,
		CreatedAt:  row.CreatedAt.Format(time.RFC3339),
		ApprovedBy: row.ApprovedBy,
		CanceledBy: row.CanceledBy,
	}

	if row.FileURL != "" {
		fileURL := row.FileURL
		item.FileURL = &fileURL
	}
	if row.FileName != "" {
		item.FileName = &row.FileName
	}

	if row.ApprovedAt != nil {
		approvedAt := row.ApprovedAt.Format(time.RFC3339)
		item.ApprovedAt = &approvedAt
	}
	if row.CanceledAt != nil {
		canceledAt := row.CanceledAt.Format(time.RFC3339)
		item.CanceledAt = &canceledAt
	}

	return item
}

func (service *LeaveServiceImpl) LeaveCheck(custID string, empID int64) (*entity.LeaveCheckResponse, error) {
	var resp entity.LeaveCheckResponse
	timeNow, err := times.GetCurrentTime()
	if err != nil {
		return nil, err
	}
	parseFormatDate := timeNow.Format("2006-01-02")
	data, err := service.LeaveRequestRepo.FindByCustIDAndEmpID(custID, empID, parseFormatDate, parseFormatDate)
	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
		default:
			return nil, err
		}
	}
	if data == nil {
		return nil, nil
	}

	resp = entity.LeaveCheckResponse{
		LeaveID:   &data.LeaveID,
		StartDate: data.StartDate.Format(time.DateOnly),
		EndDate:   data.EndDate.Format(time.DateOnly),
		Reason:    data.Reason,
	}
	return &resp, nil
}
