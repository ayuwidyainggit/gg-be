package entity

import (
	"time"
)

type DistPriceGroupResponse struct {
	DistPriceGrpId   int        `json:"dist_price_grp_id"`
	DistPriceGrpCode string     `json:"dist_price_grp_code"`
	DistPriceGrpName string     `json:"dist_price_grp_name"`
	IsActive         bool       `json:"is_active"`
	UpdatedAt        *time.Time `json:"updated_at"`
	UpdatedByName    *string    `json:"updated_by_name"`
}

type CreateDistPriceGroupBody struct {
	CustId           string `json:"cust_id" validate:"required,max=10"`
	CreatedBy        int64  `json:"created_by" validate:"required"`
	DistPriceGrpCode string `json:"dist_price_grp_code" validate:"required,max=6,alphanumericSpace"`
	DistPriceGrpName string `json:"dist_price_grp_name" validate:"required,max=25"`
	IsActive         bool   `json:"is_active"`
}

type DetailDistPriceGroupParams struct {
	DistPriceGrpId int `params:"dist_price_grp_id" validate:"required" json:"dist_price_grp_id"`
}

type UpdateDistPriceGroupParams struct {
	DistPriceGroupId int `params:"dist_price_grp_id" validate:"required"`
}

type DeleteDistPriceGroupParams struct {
	DistPriceGroupId int `params:"dist_price_grp_id" validate:"required"`
}

type UpdateDistPriceGroupRequest struct {
	CustId           string `json:"cust_id" validate:"required,max=10"`
	UpdatedBy        int64  `json:"updated_by" validate:"required"`
	DistPriceGrpCode string `json:"dist_price_grp_code,omitempty" validate:"required,max=6,alphanumericSpace"`
	DistPriceGrpName string `json:"dist_price_grp_name,omitempty" validate:"max=25,omitempty"`
	IsActive         *bool  `json:"is_active,omitempty"`
}
