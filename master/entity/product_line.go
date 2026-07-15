package entity

import (
	"time"
)

type ProductLineQueryFilter struct {
	Page     int    `query:"page"`
	Limit    int    `query:"limit" validate:"required"`
	Query    string `query:"q"`
	Mode     string `query:"mode"`
	Sort     string `query:"sort"`
	PlIds    []int  `query:"pl_ids"`
	IsActive *int   `query:"is_active"`
}

type ProductLineResponse struct {
	PLId          int        `json:"pl_id"`
	PLCode        string     `json:"pl_code"`
	PLName        string     `json:"pl_name"`
	EffCall       int        `json:"eff_call"`
	MinItem       int        `json:"min_item"`
	IsActive      bool       `json:"is_active"`
	UpdatedBy     *int64     `json:"updated_by,omitempty"`
	UpdatedAt     *time.Time `json:"updated_at"`
	UpdatedByName string     `json:"updated_by_name"`
}

type ProductLineListResponse struct {
	PLId          int        `json:"pl_id"`
	PLCode        string     `json:"pl_code"`
	PLName        string     `json:"pl_name"`
	EffCall       int        `json:"eff_call"`
	MinItem       int        `json:"min_item"`
	IsActive      bool       `json:"is_active"`
	UpdatedBy     *int64     `json:"updated_by,omitempty"`
	UpdatedAt     *time.Time `json:"updated_at"`
	UpdatedByName string     `json:"updated_by_name"`
}

type ProductLineLookupResponse struct {
	PlId     int    `json:"pl_id"`
	PlCode   string `json:"pl_code"`
	PlName   string `json:"pl_name"`
	EffCall  int    `json:"eff_call"`
	MinItem  int    `json:"min_item"`
	IsActive bool   `json:"is_active"`
}

type CreateProductLineBody struct {
	CustId    string `json:"cust_id" validate:"required,max=10"`
	CreatedBy int64  `json:"created_by" validate:"required"`
	PlCode    string `json:"pl_code" validate:"required,max=5,alphanumericSpace"`
	PlName    string `json:"pl_name" validate:"required,max=40,alphanumericSpace"`
	EffCall   int    `json:"eff_call" validate:"min=0"`
	MinItem   int    `json:"min_item" validate:"min=0"`
	IsActive  bool   `json:"is_active"`
}

type DetailProductLineParams struct {
	PlId int `params:"pl_id" validate:"required"`
}

type UpdateProductLineParams struct {
	PlId int `params:"pl_id" validate:"required"`
}

type DeleteProductLineParams struct {
	PlId int `params:"pl_id" validate:"required"`
}

type UpdateProductLineRequest struct {
	CustId    string  `json:"cust_id" validate:"required,max=10"`
	UpdatedBy int64   `json:"updated_by" validate:"required"`
	PlCode    string  `json:"pl_code,omitempty" validate:"max=5,omitempty,alphanumericSpace"`
	PlName    string  `json:"pl_name,omitempty" validate:"max=40,omitempty,alphanumericSpace"`
	EffCall   float32 `json:"eff_call" validate:"min=0"`
	MinItem   float32 `json:"min_item" validate:"min=0"`
	IsActive  *bool   `json:"is_active,omitempty"`
}
