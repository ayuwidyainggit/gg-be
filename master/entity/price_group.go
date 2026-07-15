package entity

import (
	"time"
)

type PriceGroupResponse struct {
	PriceGrpId    int        `json:"price_grp_id"`
	PriceGrpCode  string     `json:"price_grp_code"`
	PriceGrpName  string     `json:"price_grp_name"`
	IsActive      bool       `json:"is_active"`
	UpdatedBy     *int64     `json:"updated_by"`
	UpdatedAt     *time.Time `json:"updated_at"`
	UpdatedByName *string    `json:"updated_by_name"`
}

type PriceGroupLookupResponse struct {
	PriceGrpId   int    `json:"price_grp_id"`
	PriceGrpCode string `json:"price_grp_code"`
	PriceGrpName string `json:"price_grp_name"`
}

type CreatePriceGroupBody struct {
	CustId       string `json:"cust_id" validate:"required,max=10"`
	CreatedBy    int64  `json:"created_by" validate:"required"`
	PriceGrpCode string `json:"price_grp_code" validate:"required,max=10,alphanumericSpace"`
	PriceGrpName string `json:"price_grp_name" validate:"required,max=150"`
	IsActive     bool   `json:"is_active"`
}

type DetailPriceGroupParams struct {
	PriceGrpId int `params:"price_grp_id" validate:"required"`
}

type UpdatePriceGroupParams struct {
	PriceGrpId int `params:"price_grp_id" validate:"required"`
}

type DeletePriceGroupParams struct {
	PriceGrpId int `params:"price_grp_id" validate:"required"`
}

type UpdatePriceGroupRequest struct {
	CustId       string `json:"cust_id" validate:"required,max=10"`
	UpdatedBy    int64  `json:"updated_by" validate:"required"`
	PriceGrpCode string `json:"price_grp_code,omitempty" validate:"required,max=10,omitempty,alphanumericSpace"`
	PriceGrpName string `json:"price_grp_name,omitempty" validate:"max=150,omitempty"`
	IsActive     *bool  `json:"is_active,omitempty"`
}
