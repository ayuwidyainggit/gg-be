package entity

import (
	"time"
)

type PluGroupResponse struct {
	PluGrpId      int        `json:"plu_grp_id"`
	PluGrpCode    string     `json:"plu_grp_code"`
	PluGrpName    string     `json:"plu_grp_name"`
	IsActive      bool       `json:"is_active"`
	UpdatedBy     *int64     `json:"updated_by"`
	UpdatedAt     *time.Time `json:"updated_at"`
	UpdatedByName *string    `json:"updated_by_name"`
}

type PluGroupLookupResponse struct {
	PluGrpId   int    `json:"plu_grp_id"`
	PluGrpCode string `json:"plu_grp_code"`
	PluGrpName string `json:"plu_grp_name"`
}

type CreatePluGroupBody struct {
	CustId     string `json:"cust_id" validate:"required,max=10"`
	CreatedBy  int64  `json:"created_by" validate:"required"`
	PluGrpCode string `json:"plu_grp_code" validate:"required,max=10,alphanumericSpace"`
	PluGrpName string `json:"plu_grp_name" validate:"required,max=150"`
	IsActive   bool   `json:"is_active"`
}

type DetailPluGroupParams struct {
	PluGrpId int `params:"plu_grp_id" validate:"required"`
}

type UpdatePluGroupParams struct {
	PluGrpId int `params:"plu_grp_id" validate:"required"`
}

type DeletePluGroupParams struct {
	PluGrpId int `params:"plu_grp_id" validate:"required"`
}

type UpdatePluGroupRequest struct {
	CustId     string `json:"cust_id" validate:"required,max=10"`
	UpdatedBy  int64  `json:"updated_by" validate:"required"`
	PluGrpCode string `json:"plu_grp_code,omitempty" validate:"required,max=10,omitempty,alphanumericSpace"`
	PluGrpName string `json:"plu_grp_name,omitempty" validate:"max=150,omitempty"`
	IsActive   *bool  `json:"is_active,omitempty"`
}
