package entity

import (
	"time"
)

type OutletGroupResponse struct {
	OtGrpId       int        `json:"ot_grp_id"`
	OtGrpCode     string     `json:"ot_grp_code"`
	OtGrpName     string     `json:"ot_grp_name"`
	IsActive      bool       `json:"is_active"`
	UpdatedBy     *int64     `json:"updated_by"`
	UpdatedAt     *time.Time `json:"updated_at"`
	UpdatedByName *string    `json:"updated_by_name"`
}

type OutletGroupLookupResponse struct {
	OtGrpId   int    `json:"ot_grp_id"`
	OtGrpCode string `json:"ot_grp_code"`
	OtGrpName string `json:"ot_grp_name"`
}

type CreateOutletGroupBody struct {
	CustId    string `json:"cust_id" validate:"required,max=10"`
	CreatedBy int64  `json:"created_by" validate:"required"`
	OtGrpCode string `json:"ot_grp_code" validate:"required,max=3,numeric"`
	OtGrpName string `json:"ot_grp_name" validate:"required,max=40,alphanumericSpace"`
	IsActive  bool   `json:"is_active"`
}

type DetailOutletGroupParams struct {
	OtGrpId int64 `params:"ot_grp_id" validate:"required"`
}

type UpdateOutletGroupParams struct {
	OtGrpId int `params:"ot_grp_id" validate:"required"`
}

type DeleteOutletGroupParams struct {
	OtGrpId int `params:"ot_grp_id" validate:"required"`
}

type UpdateOutletGroupRequest struct {
	CustId    string `json:"cust_id" validate:"required,max=10"`
	UpdatedBy int64  `json:"updated_by" validate:"required"`
	OtGrpCode string `json:"ot_grp_code,omitempty" validate:"max=3,omitempty,numeric"`
	OtGrpName string `json:"ot_grp_name,omitempty" validate:"max=40,omitempty,alphanumericSpace"`
	IsActive  *bool  `json:"is_active,omitempty"`
}
