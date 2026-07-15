package entity

import "mime/multipart"

type LeaveRequestCreate struct {
	StartDate string `form:"start_date" validate:"required,len=10,datetime=2006-01-02"`
	EndDate   string `form:"end_date" validate:"required,len=10,datetime=2006-01-02"`
	Reason    string `form:"reason" validate:"required"`
	File      *multipart.FileHeader
	CustID    string
	EmpID     int64
	EmpCode   string
	UserID    int64
}

type LeaveRequestQuery struct {
	CustID      string
	EmpID       int64
	FilterStart string `query:"filter_start" validate:"omitempty,len=10,datetime=2006-01-02"`
	FilterEnd   string `query:"filter_end" validate:"omitempty,len=10,datetime=2006-01-02"`
	Page        int    `query:"page"`
	Limit       int    `query:"limit"`
}

type LeaveRequestItem struct {
	CustID     string  `json:"cust_id"`
	EmpID      string  `json:"emp_id"`
	StartDate  string  `json:"start_date"`
	EndDate    string  `json:"end_date"`
	Reason     string  `json:"reason"`
	FileURL    *string `json:"file_url"`
	FileName   *string `json:"file_name"`
	Approval   string  `json:"approval"`
	Duration   string  `json:"duration"`
	CreatedBy  *string `json:"created_by"`
	CreatedAt  string  `json:"created_at"`
	ApprovedBy *string `json:"approved_by"`
	ApprovedAt *string `json:"approved_at"`
	CanceledBy *string `json:"canceled_by"`
	CanceledAt *string `json:"canceled_at"`
}

type LeaveCheckResponse struct {
	LeaveID   *int64 `json:"leave_id"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
	Reason    string `json:"reason"`
}
