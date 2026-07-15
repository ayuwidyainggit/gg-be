package entity

import (
	"time"
)

type MProductResponse struct {
	MProId    int        `json:"pro_id"`
	MProCode  string     `json:"pro_code"`
	MProName  string     `json:"pro_name"`
	IsActive  bool       `json:"is_active"`
	UpdatedBy *int64     `json:"updated_by"`
	UpdatedAt *time.Time `json:"updated_at"`
}

type CreateMProductBody struct {
	CustId    string `json:"cust_id" validate:"required,max=10"`
	CreatedBy int64  `json:"created_by" validate:"required"`
	MProCode  string `json:"pro_code" validate:"required,max=5"`
	MProName  string `json:"pro_name" validate:"required,max=150"`
	IsActive  bool   `json:"is_active"`
}

type DetailMProductParams struct {
	MProId int `params:"pro_id" validate:"required"`
}

type UpdateMProductParams struct {
	MProId int `params:"pro_id" validate:"required"`
}

type DeleteMProductParams struct {
	MProId int `params:"pro_id" validate:"required"`
}

type UpdateMProductRequest struct {
	CustId    string `json:"cust_id" validate:"required,max=10"`
	UpdatedBy int64  `json:"updated_by" validate:"required"`
	MProCode  string `json:"pro_code,omitempty" validate:"max=5,omitempty"`
	MProName  string `json:"pro_name,omitempty" validate:"max=150,omitempty"`
	IsActive  *bool  `json:"is_active,omitempty"`
}
