package entity

import (
	"time"
)

type ProductCoreTaxResponse struct {
	CatCoreTax     string     `json:"cat_coretax"`
	ProCodeCoreTax string     `json:"pro_code_coretax"`
	ProNameCoreTax string     `json:"pro_name_coretax"`
	IsActive       bool       `json:"is_active"`
	UpdatedBy      *int64     `json:"updated_by"`
	UpdatedAt      *time.Time `json:"updated_at"`
}
type ProductCoreTaxListResponse struct {
	CatCoreTax     string     `json:"cat_coretax"`
	ProCodeCoreTax string     `json:"pro_code_coretax"`
	ProNameCoreTax string     `json:"pro_name_coretax"`
	IsActive       bool       `json:"is_active"`
	UpdatedBy      *int64     `json:"updated_by"`
	UpdatedAt      *time.Time `json:"updated_at"`
	UpdatedByName  string     `json:"updated_by_name"`
}

type CreateProductCoreTaxBody struct {
	CustId         string `json:"cust_id" validate:"required,max=10"`
	CreatedBy      int64  `json:"created_by" validate:"required"`
	CatCoreTax     string `json:"cat_coretax" validate:"required,max=15,alphanumericSpace"`
	ProCodeCoreTax string `json:"pro_code_coretax" validate:"required,max=15,alphanumericSpace"`
	ProNameCoreTax string `json:"pro_name_coretax" validate:"required,max=15,alphanumericSpace"`
	IsActive       bool   `json:"is_active"`
}

type DetailProductCoreTaxParams struct {
	ProCodeCoreTax string `params:"pro_code_coretax" validate:"required"`
}

type UpdateProductCoreTaxParams struct {
	ProCodeCoreTax string `params:"pro_code_coretax" validate:"required"`
}

type DeleteProductCoreTaxParams struct {
	ProCodeCoreTax string `params:"pro_code_coretax" validate:"required"`
}

type UpdateProductCoreTaxRequest struct {
	CustId         string `json:"cust_id" validate:"required,max=10"`
	UpdatedBy      int64  `json:"updated_by" validate:"required"`
	CatCoreTax     string `json:"cat_coretax,omitempty" validate:"required,max=15,alphanumericSpace"`
	ProCodeCoreTax string `json:"pro_code_coretax,omitempty" validate:"required,max=5,omitempty,alphanumericSpace"`
	ProNameCoreTax string `json:"pro_name_coretax,omitempty" validate:"max=15,omitempty,alphanumericSpace"`
	IsActive       *bool  `json:"is_active,omitempty"`
}
