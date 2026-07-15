package model

import "time"

// ReportList model for report.list table
type ReportList struct {
	ReportID   string     `gorm:"column:report_id;primaryKey" json:"report_id"`
	CustID     string     `gorm:"column:cust_id;not null" json:"cust_id"`
	ReportName string     `gorm:"column:report_name;not null" json:"report_name"`
	StartDate  time.Time  `gorm:"column:start_date;not null" json:"start_date"`
	EndDate    time.Time  `gorm:"column:end_date;not null" json:"end_date"`
	FileStatus int        `gorm:"column:file_status;not null;default:0" json:"file_status"`
	FileURL    *string    `gorm:"column:file_url" json:"file_url"`
	FileBase64 *string    `gorm:"column:file_base64" json:"file_base64"`
	CreatedBy  *string    `gorm:"column:created_by" json:"created_by"`
	CreatedAt  time.Time  `gorm:"column:created_at;not null" json:"created_at"`
	UpdatedAt  *time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func (ReportList) TableName() string {
	return "report.list"
}
