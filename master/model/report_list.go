package model

import "time"

type ReportList struct {
	ReportID   string     `db:"report_id" json:"report_id"`
	CustID     string     `db:"cust_id" json:"cust_id"`
	ReportName string     `db:"report_name" json:"report_name"`
	StartDate  *time.Time `db:"start_date" json:"start_date"`
	EndDate    *time.Time `db:"end_date" json:"end_date"`
	FileStatus int        `db:"file_status" json:"file_status"`
	FileURL    *string    `db:"file_url" json:"file_url"`
	FileBase64 *string    `db:"file_base64" json:"file_base64"`
	CreatedBy  *string    `db:"created_by" json:"created_by"`
	CreatedAt  *time.Time `db:"created_at" json:"created_at"`
	UpdatedAt  *time.Time `db:"updated_at" json:"updated_at"`
}
