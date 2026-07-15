package entity

import (
	"time"
)

type SpecialPriceGroupQueryFilter struct {
	CustId              string
	ParentCustId        string
	Page                int    `query:"page"`
	Limit               int    `query:"limit" validate:"required"`
	Query               string `query:"q"`
	Mode                string `query:"mode"`
	Sort                string `query:"sort"`
	IsActive            *int   `query:"is_active"`
	SpecialPriceGroupId int    `query:"sp_price_grp_id"`
}

type SpecialPriceGroupResponse struct {
	SpecialPriceGroupId   int        `json:"sp_price_grp_id"`
	SpecialPriceGroupCode string     `json:"sp_price_grp_code"`
	SpecialPriceGroupName string     `json:"sp_price_grp_name"`
	IsActive              bool       `json:"is_active"`
	UpdatedBy             *int64     `json:"updated_by"`
	UpdatedByName         string     `json:"updated_by_name"`
	UpdatedAt             *time.Time `json:"updated_at"`
}
type SpecialPriceGroupListResponse struct {
	SpecialPriceGroupId   int        `json:"sp_price_grp_id"`
	SpecialPriceGroupCode string     `json:"sp_price_grp_code"`
	SpecialPriceGroupName string     `json:"sp_price_grp_name"`
	IsActive              bool       `json:"is_active"`
	UpdatedBy             *int64     `json:"updated_by"`
	UpdatedAt             *time.Time `json:"updated_at"`
	UpdatedByName         string     `json:"updated_by_name"`
}

type SpecialPriceGroupLookupResponse struct {
	SpecialPriceGroupId   int    `json:"sp_price_grp_id"`
	SpecialPriceGroupCode string `json:"sp_price_grp_code"`
	SpecialPriceGroupName string `json:"sp_price_grp_name"`
}

type CreateSpecialPriceGroupBody struct {
	CustId                string `json:"cust_id" validate:"required,max=10"`
	CreatedBy             int64  `json:"created_by" validate:"required"`
	SpecialPriceGroupCode string `json:"sp_price_grp_code" validate:"required,max=6"`
	SpecialPriceGroupName string `json:"sp_price_grp_name" validate:"required,max=25,alphanumericSpace"`
	IsActive              bool   `json:"is_active"`
}

type UpdateSpecialPriceGroupRequest struct {
	CustId                string `json:"cust_id" validate:"required,max=10"`
	UpdatedBy             int64  `json:"updated_by" validate:"required"`
	SpecialPriceGroupCode string `json:"sp_price_grp_code" validate:"required,max=6"`
	SpecialPriceGroupName string `json:"sp_price_grp_name,omitempty" validate:"max=25,omitempty,alphanumericSpace"`
	IsActive              *bool  `json:"is_active,omitempty"`
}

type DetailSpecialPriceGroupParams struct {
	SpecialPriceGroupId int `params:"sp_price_grp_id" validate:"required"`
}

type UpdateSpecialPriceGroupParams struct {
	SpecialPriceGroupId int `params:"sp_price_grp_id" validate:"required"`
}

type DeleteSpecialPriceGroupParams struct {
	SpecialPriceGroupId int `params:"sp_price_grp_id" validate:"required"`
}
