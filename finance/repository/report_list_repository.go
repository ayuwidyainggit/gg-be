package repository

import (
	"context"
	"finance/model"
	"time"

	"gorm.io/gorm"
)

type ReportListRepository interface {
	CheckReportInProgress(ctx context.Context, prefix string, custID string) (bool, error)
	GetNextSequenceNumber(ctx context.Context, prefix string, custID string) (int, error)
	Create(ctx context.Context, report *model.ReportList) error
	UpdateFileReady(ctx context.Context, reportID string, fileBase64 string) error
	FindByID(ctx context.Context, reportID string, custID string) (*model.ReportList, error)
}

type ReportListRepositoryImpl struct {
	*gorm.DB
}

func NewReportListRepo(db *gorm.DB) *ReportListRepositoryImpl {
	return &ReportListRepositoryImpl{db}
}

func (r *ReportListRepositoryImpl) model(ctx context.Context) *gorm.DB {
	tx := extractTx(ctx)
	if tx != nil {
		return tx.WithContext(ctx)
	}
	return r.WithContext(ctx)
}

func (r *ReportListRepositoryImpl) CheckReportInProgress(ctx context.Context, prefix string, custID string) (bool, error) {
	var count int64
	err := r.model(ctx).
		Table("report.list").
		Where("report_name LIKE ? AND cust_id = ? AND file_status = 0", prefix+"%", custID).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *ReportListRepositoryImpl) GetNextSequenceNumber(ctx context.Context, prefix string, custID string) (int, error) {
	var count int64
	err := r.model(ctx).
		Table("report.list").
		Where("report_name LIKE ? AND cust_id = ?", prefix+"%", custID).
		Count(&count).Error
	if err != nil {
		return 0, err
	}
	return int(count) + 1, nil
}

func (r *ReportListRepositoryImpl) Create(ctx context.Context, report *model.ReportList) error {
	return r.model(ctx).Create(report).Error
}

func (r *ReportListRepositoryImpl) UpdateFileReady(ctx context.Context, reportID string, fileBase64 string) error {
	return r.model(ctx).Table("report.list").
		Where("report_id = ?", reportID).
		Updates(map[string]interface{}{
			"file_status": 1,
			"file_base64": fileBase64,
			"updated_at":  time.Now(),
		}).Error
}

func (r *ReportListRepositoryImpl) FindByID(ctx context.Context, reportID string, custID string) (*model.ReportList, error) {
	var report model.ReportList
	err := r.model(ctx).
		Table("report.list").
		Where("report_id = ? AND cust_id = ?", reportID, custID).
		First(&report).Error
	if err != nil {
		return nil, err
	}
	return &report, nil
}
