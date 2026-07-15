package model

import "time"

type OrderApprovalRequestDetail struct {
	OrderApprovalRequestID int64      `gorm:"column:order_approval_request_id" json:"order_approval_request_id"`
	EmpID                  int64      `gorm:"column:emp_id" json:"emp_id"`
	Status                 *int       `gorm:"column:status" json:"status"`
	ActDate                *time.Time `gorm:"column:act_date" json:"act_date"`
	Seq                    int        `gorm:"column:seq" json:"seq"`
	Level                  int        `gorm:"column:level" json:"level"`
}

func (OrderApprovalRequestDetail) TableName() string {
	return "sls.order_approval_requests_details"
}

type OrderApprovalRequestDetailRead struct {
	OrderApprovalRequestID int64      `gorm:"column:order_approval_request_id" json:"order_approval_request_id"`
	EmployeeId             int        `gorm:"column:emp_id" json:"emp_id"`
	EmployeeCode           string     `gorm:"column:emp_code" json:"emp_code"`
	EmployeeName           string     `gorm:"column:emp_name" json:"emp_name"`
	ImageURL               string     `gorm:"column:image_url" json:"image_url"`
	Status                 *int       `gorm:"column:status" json:"status"`
	ActDate                *time.Time `gorm:"column:act_date" json:"act_date"`
	Seq                    int        `gorm:"column:seq" json:"seq"`
	Level                  int        `gorm:"column:level" json:"level"`
}

func (OrderApprovalRequestDetailRead) TableName() string {
	return "sls.order_approval_requests_details"
}
