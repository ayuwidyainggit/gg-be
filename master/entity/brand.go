package entity

import (
	"time"
)

type BrandResponse struct {
	BrandId   int        `json:"brand_id"`
	BrandCode string     `json:"brand_code"`
	BrandName string     `json:"brand_name"`
	PlId      int        `json:"pl_id"`
	PlCode    string     `json:"pl_code"`
	PlName    string     `json:"pl_name"`
	EffCall   float32    `json:"eff_call"`
	MinItem   float32    `json:"min_item"`
	IsActive  bool       `json:"is_active"`
	UpdatedBy *int64     `json:"updated_by"`
	UpdatedAt *time.Time `json:"updated_at"`
}
type BrandListResponse struct {
	BrandId       int        `json:"brand_id"`
	BrandCode     string     `json:"brand_code"`
	BrandName     string     `json:"brand_name"`
	PlId          int        `json:"pl_id"`
	PlCode        string     `json:"pl_code"`
	PlName        string     `json:"pl_name"`
	EffCall       float32    `json:"eff_call"`
	MinItem       float32    `json:"min_item"`
	IsActive      bool       `json:"is_active"`
	UpdatedBy     *int64     `json:"updated_by"`
	UpdatedAt     *time.Time `json:"updated_at"`
	UpdatedByName string     `json:"updated_by_name"`
}

type BrandLookupResponse struct {
	BrandId   int     `json:"brand_id"`
	BrandCode string  `json:"brand_code"`
	BrandName string  `json:"brand_name"`
	PlId      int     `json:"pl_id"`
	PlCode    string  `json:"pl_code"`
	PlName    string  `json:"pl_name"`
	EffCall   float32 `json:"eff_call"`
	MinItem   float32 `json:"min_item"`
}
type BrandLookupListResponse struct {
	BrandId   int     `json:"brand_id"`
	BrandCode string  `json:"brand_code"`
	BrandName string  `json:"brand_name"`
	PlId      int     `json:"pl_id"`
	PlCode    string  `json:"pl_code"`
	PlName    string  `json:"pl_name"`
	EffCall   float32 `json:"eff_call"`
	MinItem   float32 `json:"min_item"`
}

type CreateBrandBody struct {
	CustId    string  `json:"cust_id" validate:"required,max=10"`
	BrandCode string  `json:"brand_code" validate:"required,max=5,alphanumericSpace"`
	BrandName string  `json:"brand_name" validate:"required,max=40,alphanumericSpace"`
	PlId      int     `json:"pl_id" validate:"required"`
	EffCall   float32 `json:"eff_call" validate:"min=0"`
	MinItem   float32 `json:"min_item" validate:"min=0"`
	IsActive  bool    `json:"is_active"`
	CreatedBy int64   `json:"created_by" validate:"required"`
}

type DetailBrandParams struct {
	BrandId int `params:"brand_id" validate:"required"`
}

type UpdateBrandParams struct {
	BrandId int `params:"brand_id" validate:"required"`
}

type DeleteBrandParams struct {
	BrandId int `params:"brand_id" validate:"required"`
}

type UpdateBrandRequest struct {
	CustId    string  `json:"cust_id" validate:"required,max=10"`
	BrandCode string  `json:"brand_code,omitempty" validate:"max=5,alphanumericSpace"`
	BrandName string  `json:"brand_name,omitempty" validate:"max=40,omitempty,alphanumericSpace"`
	PlId      int     `json:"pl_id" validate:"required"`
	EffCall   float32 `json:"eff_call" validate:"min=0"`
	MinItem   float32 `json:"min_item" validate:"min=0"`
	IsActive  *bool   `json:"is_active,omitempty"`
	UpdatedBy int64   `json:"updated_by" validate:"required"`
}

type BrandQueryFilter struct {
	Page     int    `query:"page"`
	Limit    int    `query:"limit" validate:"required"`
	Query    string `query:"q"`
	Mode     string `query:"mode"`
	Sort     string `query:"sort"`
	PlId     int    `query:"pl_id"`
	PlIds    []int  `query:"pl_ids"`
	BrandId  int    `query:"brand_id"`
	IsActive *int   `query:"is_active"`
	BrandIds []int  `query:"brand_ids"`
}
