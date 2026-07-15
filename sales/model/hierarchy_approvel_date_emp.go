package model

import (
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type HierarchyApprovalDetEmp struct {
	HierarchyApprovalDetailEmpID *int64 `gorm:"column:hierarchy_approval_detail_emp_id;default:nextval('sls.hierarchy_approval_detail_emp_seq'::regclass);not null" json:"hierarchy_approval_detail_emp_id"`
	HierarchyApprovalDetailID    int64  `gorm:"column:hierarchy_approval_detail_id;not null" json:"hierarchy_approval_detail_id"`
	EmpID                        int64  `gorm:"emp_id" json:"emp_id"`
	Seq                          int    `gorm:"seq" json:"seq"`
}

func (HierarchyApprovalDetEmp) TableName() string {
	return "sls.hierarchy_approvals_details_emp"
}

type HierarcyApprovalEmployee struct {
	CustId       string         `gorm:"column:cust_id" json:"cust_id"`
	EmployeeId   int            `gorm:"column:emp_id" json:"emp_id"`
	EmployeeCode string         `gorm:"column:emp_code" json:"emp_code"`
	EmployeeName string         `gorm:"column:emp_name" json:"emp_name"`
	DeletedAt    gorm.DeletedAt `gorm:"column:deleted_at" json:"deleted_at"`
}

func (HierarcyApprovalEmployee) TableName() string {
	return "mst.m_employee"
}

type HierarchyApprovalDetEmpRead struct {
	HierarchyApprovalDetailEmpID int64  `gorm:"column:hierarchy_approval_detail_emp_id;default:nextval('sls.hierarchy_approval_detail_emp_seq'::regclass);not null" json:"hierarchy_approval_detail_emp_id"`
	EmpID                        int64  `gorm:"column:emp_id" json:"emp_id"`
	EmpName                      string `gorm:"column:emp_name" json:"emp_name"`
	Seq                          int    `gorm:"column:seq" json:"seq"`
}

func (HierarchyApprovalDetEmpRead) TableName() string {
	return "sls.hierarchy_approvals_details_emp"
}

type HierarchyApprovalDetEmpReadMap map[int]*HierarchyApprovalDetEmpRead

func (m HierarchyApprovalDetEmpReadMap) Set(seq int, value HierarchyApprovalDetEmpRead) {
	m[seq] = &value
}
func (m HierarchyApprovalDetEmpReadMap) GetBySequence(seq int) (value *HierarchyApprovalDetEmpRead, err error) {
	val, ok := m[seq]

	// If the key exists
	if !ok {
		return value, errors.New(fmt.Sprintf("%v Not Found", seq))
	}

	return val, nil
}

type HierarchyApprovalDetEmpUpdate struct {
	EmpID int64 `gorm:"emp_id" json:"emp_id"`
	Seq   int   `gorm:"seq" json:"seq"`
}

func (HierarchyApprovalDetEmpUpdate) TableName() string {
	return "sls.hierarchy_approvals_details_emp"
}
