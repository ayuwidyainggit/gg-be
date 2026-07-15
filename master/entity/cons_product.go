package entity

import (
	"time"
)

type ConsProductResponse struct {
	CProId    int        `json:"c_pro_id"`
	CProCode  string     `json:"c_pro_code"`
	CProName  string     `json:"c_pro_name"`
	IsActive  bool       `json:"is_active"`
	UpdatedBy *int64     `json:"updated_by"`
	UpdatedAt *time.Time `json:"updated_at"`
}

type ConsProductListResponse struct {
	CProId        int        `json:"c_pro_id"`
	CProCode      string     `json:"c_pro_code"`
	CProName      string     `json:"c_pro_name"`
	IsActive      bool       `json:"is_active"`
	UpdatedBy     *int64     `json:"updated_by"`
	UpdatedAt     *time.Time `json:"updated_at"`
	UpdatedByName string     `json:"updated_by_name"`
}

type ConsProductLookupResponse struct {
	CProId   int    `json:"c_pro_id"`
	CProCode string `json:"c_pro_code"`
	CProName string `json:"c_pro_name"`
}

type CreateConsProductBody struct {
	CustId    string `json:"cust_id" validate:"required,max=10"`
	CreatedBy int64  `json:"created_by" validate:"required"`
	CProCode  string `json:"c_pro_code" validate:"required,max=5,alphanumericSpace"`
	CProName  string `json:"c_pro_name" validate:"required,max=25,alphanumericSpace"`
	IsActive  bool   `json:"is_active"`
}

type DetailConsProductParams struct {
	CProId int `params:"c_pro_id" validate:"required"`
}

type UpdateConsProductParams struct {
	CProId int `params:"c_pro_id" validate:"required"`
}

type DeleteConsProductParams struct {
	CProId int `params:"c_pro_id" validate:"required"`
}

type UpdateConsProductRequest struct {
	CustId    string `json:"cust_id" validate:"required,max=10"`
	UpdatedBy int64  `json:"updated_by" validate:"required"`
	CProCode  string `json:"c_pro_code,omitempty" validate:"max=5,alphanumericSpace"`
	CProName  string `json:"c_pro_name,omitempty" validate:"max=25,omitempty,alphanumericSpace"`
	IsActive  *bool  `json:"is_active,omitempty"`
}
