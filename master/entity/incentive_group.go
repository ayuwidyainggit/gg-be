package entity

import (
	"time"
)

type IncentiveGroupResponse struct {
	IncGrpID      int        `json:"inc_grp_id"`
	IncGrpCode    string     `json:"inc_grp_code"`
	IncGrpName    string     `json:"inc_grp_name"`
	IsActive      bool       `json:"is_active"`
	UpdatedBy     *int64     `json:"updated_by"`
	UpdatedAt     *time.Time `json:"updated_at"`
	UpdatedByName *string    `json:"updated_by_name"`
}

type IncentiveGroupLookupResponse struct {
	IncGrpID   int    `json:"inc_grp_id"`
	IncGrpCode string `json:"inc_grp_code"`
	IncGrpName string `json:"inc_grp_name"`
}

type CreateIncentiveGroupBody struct {
	CustId     string `json:"cust_id" validate:"required,max=10"`
	CreatedBy  int64  `json:"created_by" validate:"required"`
	IncGrpCode string `json:"inc_grp_code" validate:"required,max=3,numeric"`
	IncGrpName string `json:"inc_grp_name" validate:"required,max=40,alphanumericSpace"`
	IsActive   bool   `json:"is_active"`
}

type DetailIncentiveGroupParams struct {
	IncGrpID int `params:"inc_grp_id" validate:"required"`
}

type UpdateIncentiveGroupParams struct {
	IncGrpID int `params:"inc_grp_id" validate:"required"`
}

type DeleteIncentiveGroupParams struct {
	IncGrpID int `params:"inc_grp_id" validate:"required"`
}

type UpdateIncentiveGroupRequest struct {
	CustId     string `json:"cust_id" validate:"required,max=10"`
	UpdatedBy  int64  `json:"updated_by" validate:"required"`
	IncGrpCode string `json:"inc_grp_code,omitempty" validate:"required,max=3,numeric"`
	IncGrpName string `json:"inc_grp_name,omitempty" validate:"max=40,alphanumericSpace"`
	IsActive   *bool  `json:"is_active,omitempty"`
}
