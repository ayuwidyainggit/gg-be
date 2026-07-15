package entity

import (
	"errors"
	"fmt"
	"time"
)

const (
	APPROVAL_TYPE_DIRECT     = 1
	APPROVAL_TYPE_MULTILEVEL = 2
)

type CreateHierarchyApprovalBody struct {
	CustID       string
	ParentCustId string
	UserID       int64
	SetupFor     []string                            `json:"setup_for" validate:"required"`
	ApprovalType int                                 `json:"approval_type" validate:"required,oneof=1 2"`
	Details      []CreateDetailHierarchyApprovalBody `json:"details" validate:"required,min=1,dive"`
}

type CreateDetailHierarchyApprovalBody struct {
	IsActive     *bool                                  `json:"is_active" validate:"required"`
	Level        int                                    `json:"level" validate:"required"`
	CompanyID    string                                 `json:"company_id" validate:"required"`
	MaxOverLimit *float64                               `json:"max_over_limit"`
	EmpIDs       []CreateEmpDetailHierarchyApprovalBody `json:"emp_ids" validate:"required,min=1,dive"`
}

type CreateEmpDetailHierarchyApprovalBody struct {
	Sequence int   `json:"seq" validate:"required"`
	EmpID    int64 `json:"emp_id" validate:"required"`
}

type TempEmployeeValidationMap map[int64]int

func (m TempEmployeeValidationMap) SetTempEmployeeValidationMap(id int64, value int) {
	m[id] = value
}
func (m TempEmployeeValidationMap) GetByID(id int64) (value *int, err error) {
	val, ok := m[id]

	// If the key exists
	if !ok {
		return value, errors.New(fmt.Sprintf("%v Not Found", id))
	}

	return &val, nil
}

type HierarcyApprovalQueryFilter struct {
	CustId       string
	ParentCustId string
	Company      string `query:"company"`
	ApprovalType *int   `query:"approval_type"`
	Page         int    `query:"page"`
	Limit        int    `query:"limit"`
}

type HierarchyApprovalListResp struct {
	HierarchyApprovalID int64     `json:"hierarchy_approval_id"`
	SetupFor            string    `json:"setup_for"`
	SetupForName        string    `json:"setup_for_name"`
	SetupForCode        string    `json:"setup_for_code"`
	ApprovalType        int       `json:"approval_type"`
	UpdatedBy           *int64    `json:"updated_by"`
	UpdatedByname       string    `json:"updated_by_name"`
	UpdatedAt           time.Time `json:"updated_at"`
}

type DetailHierarchyApprovalParams struct {
	HierarchyApprovalID int64 `params:"hierarchy_approval_id" validate:"required"`
}

type DeleteHierarchyApprovalParams struct {
	HierarchyApprovalID int64 `params:"hierarchy_approval_id" validate:"required"`
}

type UpdateHierarchyApprovalParams struct {
	HierarchyApprovalID int64 `params:"hierarchy_approval_id" validate:"required"`
}

type DetailRequestApprovalParams struct {
	RequestApprovalID int64 `params:"request_approval_id" validate:"required"`
}

type HierarchyApprovalReadResp struct {
	HierarchyApprovalID int64                         `json:"hierarchy_approval_id"`
	SetupFor            string                        `json:"setup_for"`
	SetupForName        string                        `json:"setup_for_name"`
	SetupForCode        string                        `json:"setup_for_code"`
	ApprovalType        int                           `json:"approval_type"`
	Details             []HierarchyApprovalDetailResp `json:"details"`
}

type HierarchyApprovalDetailResp struct {
	IsActive     bool                             `json:"is_active" validate:"required"`
	Level        int                              `json:"level" validate:"required"`
	CompanyID    string                           `json:"company_id" validate:"required"`
	CompanyName  string                           `json:"company_name" validate:"required"`
	MaxOverLimit *float64                         `json:"max_over_limit"`
	EmpIDs       []HierarchyApprovalDetailEmpResp `json:"emp_ids" validate:"required,min=1,dive"`
}
type HierarchyApprovalDetailEmpResp struct {
	Sequence int    `json:"seq"`
	EmpID    int64  `json:"emp_id" `
	EmpName  string `json:"emp_name" `
}

type UpdateHierarchyApprovalBody struct {
	CustID       string
	UserID       int64
	SetupFor     string                              `json:"setup_for" validate:"required"`
	ApprovalType int                                 `json:"approval_type" validate:"required,oneof=1 2"`
	Details      []UpdateDetailHierarchyApprovalBody `json:"details" validate:"required,min=1,dive"`
}

type UpdateDetailHierarchyApprovalBody struct {
	IsActive     *bool                                  `json:"is_active" validate:"required"`
	Level        int                                    `json:"level" validate:"required"`
	CompanyID    string                                 `json:"company_id" validate:"required"`
	MaxOverLimit *float64                               `json:"max_over_limit"`
	EmpIDs       []UpdateEmpDetailHierarchyApprovalBody `json:"emp_ids" validate:"required,min=1,dive"`
}

type UpdateEmpDetailHierarchyApprovalBody struct {
	Sequence int   `json:"seq" validate:"required"`
	EmpID    int64 `json:"emp_id" validate:"required"`
}

type EmployeeHierarchyQueryFilter struct {
	CustId string `query:"cust_id" validate:"required"`
	Query  string `query:"q"`
	Sort   string `query:"sort"`
}

type EmployeeResp struct {
	CustID       string `json:"cust_id"`
	EmployeeId   int    `json:"emp_id"`
	EmployeeCode string `json:"emp_code"`
	EmployeeName string `json:"emp_name"`
}

type RequestApprovalBody struct {
	CustID string
	UserID int64
	RoNo   string `json:"ro_no"`
}

type RequestApprovalDetail struct {
	EmployeeId       int        `json:"emp_id"`
	EmployeeCode     string     `json:"emp_code"`
	EmployeeName     string     `json:"emp_name"`
	EmployeeImageURL string     `json:"emp_image_url"`
	Sequence         int        `json:"sequence"`
	Status           *int       `json:"status"`
	ActDate          *time.Time `json:"act_date"`
}

type ApprovalDetail struct {
	RoNo       string            `json:"ro_no"`
	FinishedAt *time.Time        `json:"finished_at"`
	Approvals  []GroupedApproval `json:"approvals"`
}

type GroupedApproval struct {
	Level   int                     `json:"level"`
	Details []RequestApprovalDetail `json:"details"`
}
