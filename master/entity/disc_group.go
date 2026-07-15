package entity

import (
	"time"
)

type DiscGroupResponse struct {
	DiscGrpId     int        `json:"disc_grp_id"`
	DiscGrpCode   string     `json:"disc_grp_code"`
	DiscGrpName   string     `json:"disc_grp_name"`
	IsActive      bool       `json:"is_active"`
	UpdatedBy     *int64     `json:"updated_by"`
	UpdatedAt     *time.Time `json:"updated_at"`
	UpdatedByName *string    `json:"updated_by_name"`
}

type DiscGroupLookupResponse struct {
	DiscGrpId   int    `json:"disc_grp_id"`
	DiscGrpCode string `json:"disc_grp_code"`
	DiscGrpName string `json:"disc_grp_name"`
}

type CreateDiscGroupBody struct {
	CustId      string `json:"cust_id" validate:"required,max=10"`
	CreatedBy   int64  `json:"created_by" validate:"required"`
	DiscGrpCode string `json:"disc_grp_code" validate:"required,max=5,alphanumericSpace"`
	DiscGrpName string `json:"disc_grp_name" validate:"required,max=40,alphanumericSpace"`
	IsActive    bool   `json:"is_active"`
}

type DetailDiscGroupParams struct {
	DiscGrpId int `params:"disc_grp_id" validate:"required"`
}

type UpdateDiscGroupParams struct {
	DiscGrpId int `params:"disc_grp_id" validate:"required"`
}

type DeleteDiscGroupParams struct {
	DiscGrpId int `params:"disc_grp_id" validate:"required"`
}

type UpdateDiscGroupRequest struct {
	CustId      string `json:"cust_id" validate:"required,max=10"`
	UpdatedBy   int64  `json:"updated_by" validate:"required"`
	DiscGrpCode string `json:"disc_grp_code,omitempty" validate:"required,max=5,alphanumericSpace"`
	DiscGrpName string `json:"disc_grp_name,omitempty" validate:"max=40,omitempty,alphanumericSpace"`
	IsActive    *bool  `json:"is_active,omitempty"`
}
