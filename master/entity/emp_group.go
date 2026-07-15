package entity

import (
	"time"
)

type EmpGroupResponse struct {
	EmpGroupId    int        `json:"emp_grp_id"`
	EmpGroupCode  string     `json:"emp_grp_code"`
	EmpGroupName  string     `json:"emp_grp_name"`
	IsActive      bool       `json:"is_active"`
	UpdatedBy     *int64     `json:"updated_by,omitempty"`
	UpdatedAt     *time.Time `json:"updated_at,omitempty"`
	UpdatedByName string     `json:"updated_by_name"`
}

type EmpGroupLookupResponse struct {
	EmpGroupId   int    `json:"emp_grp_id"`
	EmpGroupCode string `json:"emp_grp_code"`
	EmpGroupName string `json:"emp_grp_name"`
}

type EmpTypeResponse struct {
	EmpTypeId   string `json:"emp_type_id"`
	EmpTypeName string `json:"emp_type_name"`
}

type CreateEmpGroupBody struct {
	CustId       string `json:"cust_id" validate:"required,max=10"`
	CreatedBy    int64  `json:"created_by" validate:"required"`
	EmpGroupCode string `json:"emp_grp_code" validate:"required,number,max=5"`
	EmpGroupName string `json:"emp_grp_name" validate:"required,alphanumericSpace,max=20"`
	IsActive     bool   `json:"is_active"`
}

type DetailEmpGroupParams struct {
	EmpGroupId int `params:"emp_grp_id" validate:"required"`
}

type DetailEmpTypeParams struct {
	EmpTypeId string `params:"emp_type_id" validate:"required"`
}

type UpdateEmpGroupParams struct {
	EmpGroupId int `params:"emp_grp_id" validate:"required"`
}

type DeleteEmpGroupParams struct {
	EmpGroupId int `params:"emp_grp_id" validate:"required"`
}

type UpdateEmpGroupRequest struct {
	CustId       string `json:"cust_id" validate:"required,max=10"`
	UpdatedBy    int64  `json:"updated_by" validate:"required"`
	EmpGroupCode string `json:"emp_grp_code,omitempty" validate:"required,number,max=5"`
	EmpGroupName string `json:"emp_grp_name,omitempty" validate:"max=20,alphanumericSpace,omitempty"`
	IsActive     *bool  `json:"is_active,omitempty"`
}
