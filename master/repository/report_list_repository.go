package repository

import (
	"database/sql"
	"errors"
	"master/model"
	"time"

	"github.com/jmoiron/sqlx"
)

type ReportListRepository interface {
	CountInProgress(reportNamePrefix, custID string, processingStatus int) (int64, error)
	CountByPrefixAndDate(reportNamePrefix, custID string, startOfDay, endOfDay time.Time) (int64, error)
	Store(report model.ReportList) error
	UpdateFileResult(reportID string, fileStatus int, fileBase64 *string, fileURL *string) error
	FindByReportID(reportID, custID string) (model.ReportList, error)
}

type reportListRepositoryImpl struct {
	*sqlx.DB
}

func NewReportListRepository(db *sqlx.DB) ReportListRepository {
	return &reportListRepositoryImpl{DB: db}
}

func (r *reportListRepositoryImpl) CountInProgress(reportNamePrefix, custID string, processingStatus int) (int64, error) {
	var total int64
	query := `SELECT COUNT(*)
		FROM report.list
		WHERE cust_id = $1 AND report_name LIKE $2 AND file_status = $3`
	err := r.DB.Get(&total, query, custID, reportNamePrefix+"%", processingStatus)
	return total, err
}

func (r *reportListRepositoryImpl) CountByPrefixAndDate(reportNamePrefix, custID string, startOfDay, endOfDay time.Time) (int64, error) {
	var total int64
	query := `SELECT COUNT(*)
		FROM report.list
		WHERE cust_id = $1
			AND report_name LIKE $2
			AND created_at >= $3
			AND created_at < $4`
	err := r.DB.Get(&total, query, custID, reportNamePrefix+"%", startOfDay, endOfDay)
	return total, err
}

func (r *reportListRepositoryImpl) Store(report model.ReportList) error {
	query := `INSERT INTO report.list
		(report_id, cust_id, report_name, start_date, end_date, file_status, file_url, file_base64, created_by, created_at, updated_at)
		VALUES (:report_id, :cust_id, :report_name, :start_date, :end_date, :file_status, :file_url, :file_base64, :created_by, :created_at, :updated_at)`
	_, err := r.DB.NamedExec(query, report)
	return err
}

func (r *reportListRepositoryImpl) UpdateFileResult(reportID string, fileStatus int, fileBase64 *string, fileURL *string) error {
	now := time.Now().UTC()
	query := `UPDATE report.list
		SET file_status = $1, file_base64 = $2, file_url = $3, updated_at = $4
		WHERE report_id = $5`
	_, err := r.DB.Exec(query, fileStatus, fileBase64, fileURL, now, reportID)
	return err
}

func (r *reportListRepositoryImpl) FindByReportID(reportID, custID string) (model.ReportList, error) {
	var report model.ReportList
	query := `SELECT report_id, cust_id, report_name, start_date, end_date, file_status, file_url, file_base64, created_by, created_at, updated_at
		FROM report.list
		WHERE report_id = $1 AND cust_id = $2`
	err := r.DB.Get(&report, query, reportID, custID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.ReportList{}, err
		}
		return model.ReportList{}, err
	}
	return report, nil
}
