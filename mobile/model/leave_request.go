package model

import "time"

type LeaveRequest struct {
	LeaveID    int64      `gorm:"column:leave_id;primaryKey;autoIncrement"`
	CustID     string     `gorm:"column:cust_id"`
	EmpID      int64      `gorm:"column:emp_id"`
	StartDate  time.Time  `gorm:"column:start_date;type:date"`
	EndDate    time.Time  `gorm:"column:end_date;type:date"`
	Reason     string     `gorm:"column:reason"`
	FileURL    string     `gorm:"column:file_url"`
	FileName   string     `gorm:"column:file_name"`
	Approval   string     `gorm:"column:approval"`
	CreatedBy  int64      `gorm:"column:created_by"`
	CreatedAt  time.Time  `gorm:"column:created_at"`
	ApprovedBy *int64     `gorm:"column:approved_by"`
	ApprovedAt *time.Time `gorm:"column:approved_at"`
	CanceledBy *int64     `gorm:"column:canceled_by"`
	CanceledAt *time.Time `gorm:"column:canceled_at"`
}

func (LeaveRequest) TableName() string {
	return "mobile.leave_request"
}

type LeaveRequestRead struct {
	CustID     string     `gorm:"column:cust_id"`
	EmpID      int64      `gorm:"column:emp_id"`
	StartDate  time.Time  `gorm:"column:start_date"`
	EndDate    time.Time  `gorm:"column:end_date"`
	Reason     string     `gorm:"column:reason"`
	FileURL    string     `gorm:"column:file_url"`
	FileName   string     `gorm:"column:file_name"`
	Approval   string     `gorm:"column:approval"`
	Duration   string     `gorm:"column:duration"`
	CreatedBy  *string    `gorm:"column:created_by"`
	CreatedAt  time.Time  `gorm:"column:created_at"`
	ApprovedBy *string    `gorm:"column:approved_by"`
	ApprovedAt *time.Time `gorm:"column:approved_at"`
	CanceledBy *string    `gorm:"column:canceled_by"`
	CanceledAt *time.Time `gorm:"column:canceled_at"`
}
