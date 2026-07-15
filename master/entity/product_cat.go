package entity

import (
	"time"
)

type ProductCatResponse struct {
	PCatId    int        `json:"pcat_id"`
	PCatCode  string     `json:"pcat_code"`
	PCatName  string     `json:"pcat_name"`
	IsActive  bool       `json:"is_active"`
	UpdatedBy *int64     `json:"updated_by"`
	UpdatedAt *time.Time `json:"updated_at"`
}

type ProductCatListResponse struct {
	PCatId        int        `json:"pcat_id"`
	PCatCode      string     `json:"pcat_code"`
	PCatName      string     `json:"pcat_name"`
	IsActive      bool       `json:"is_active"`
	UpdatedBy     *int64     `json:"updated_by"`
	UpdatedAt     *time.Time `json:"updated_at"`
	UpdatedByName string     `json:"updated_by_name"`
}

type ProductCatLookupResponse struct {
	PCatId   int    `json:"pcat_id"`
	PCatCode string `json:"pcat_code"`
	PCatName string `json:"pcat_name"`
	IsActive bool   `json:"is_active"`
}

type CreateProductCatBody struct {
	CustId    string `json:"cust_id" validate:"required,max=10"`
	CreatedBy int64  `json:"created_by" validate:"required"`
	PCatCode  string `json:"pcat_code" validate:"required,max=5,alphanumericSpace"`
	PCatName  string `json:"pcat_name" validate:"required,max=40,alphanumericSpace"`
	IsActive  bool   `json:"is_active"`
}

type DetailProductCatParams struct {
	PCatId int `params:"pcat_id" validate:"required"`
}

type UpdateProductCatParams struct {
	PCatId int `params:"pcat_id" validate:"required"`
}

type DeleteProductCatParams struct {
	PCatId int `params:"pcat_id" validate:"required"`
}

type UpdateProductCatRequest struct {
	CustId    string `json:"cust_id" validate:"required,max=10"`
	UpdatedBy int64  `json:"updated_by" validate:"required"`
	PCatCode  string `json:"pcat_code,omitempty" validate:"max=5,omitempty,alphanumericSpace"`
	PCatName  string `json:"pcat_name,omitempty" validate:"max=40,omitempty,alphanumericSpace"`
	IsActive  *bool  `json:"is_active,omitempty"`
}
