package entity

import "time"

// ImportHistoryRow is a row of import.import_history for listing
type ImportHistoryRow struct {
	HistoryId      int64     `db:"history_id" json:"history_id"`
	FileName       string    `db:"file_name" json:"file_name"`
	UploadedBy     int64     `db:"uploaded_by" json:"-"`
	UploadedByName *string   `db:"uploaded_by_name" json:"uploaded_by"`
	UploadDate     time.Time `db:"upload_date" json:"upload_date"`
	Successful     int       `db:"successful_data" json:"successful_data"`
	Failed         int       `db:"failed_data" json:"failed_data"`
	Total          int       `db:"total_data" json:"total_data"`
	Status         string    `db:"status" json:"status"`
	StatusReupload bool      `db:"status_reupload" json:"status_reupload"`
	UploadType     *string   `db:"upload_type" json:"upload_type"`
	ErrorMessages  []string  `db:"-" json:"error_messages,omitempty"`
	ErrorStatuses  []string  `db:"-" json:"error_statuses,omitempty"`
	LogError       *string   `db:"-" json:"log_error,omitempty"`
}
