package entity

import (
	"time"
)

type ConvGroupResponse struct {
	ConvGrpId     int        `json:"conv_grp_id"`
	ConvGrpCode   string     `json:"conv_grp_code"`
	ConvGrpName   string     `json:"conv_grp_name"`
	IsActive      bool       `json:"is_active"`
	UpdatedBy     *int64     `json:"updated_by"`
	UpdatedAt     *time.Time `json:"updated_at"`
	UpdatedByName *string    `json:"updated_by_name"`
}

type ConvGroupLookupResponse struct {
	ConvGrpId   int    `json:"conv_grp_id"`
	ConvGrpCode string `json:"conv_grp_code"`
	ConvGrpName string `json:"conv_grp_name"`
}

type CreateConvGroupBody struct {
	CustId      string `json:"cust_id" validate:"required,max=10"`
	CreatedBy   int64  `json:"created_by" validate:"required"`
	ConvGrpCode string `json:"conv_grp_code" validate:"required,max=10,alphanumericSpace"`
	ConvGrpName string `json:"conv_grp_name" validate:"required,max=150"`
	IsActive    bool   `json:"is_active"`
}

type DetailConvGroupParams struct {
	ConvGrpId int `params:"conv_grp_id" validate:"required"`
}

type UpdateConvGroupParams struct {
	ConvGrpId int `params:"conv_grp_id" validate:"required"`
}

type DeleteConvGroupParams struct {
	ConvGrpId int `params:"conv_grp_id" validate:"required"`
}

type UpdateConvGroupRequest struct {
	CustId      string `json:"cust_id" validate:"required,max=10"`
	UpdatedBy   int64  `json:"updated_by" validate:"required"`
	ConvGrpCode string `json:"conv_grp_code,omitempty" validate:"required,max=10,alphanumericSpace"`
	ConvGrpName string `json:"conv_grp_name,omitempty" validate:"max=150,omitempty"`
	IsActive    *bool  `json:"is_active,omitempty"`
}
