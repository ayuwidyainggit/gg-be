package entity

import (
	"time"
)

type InvoiceDiscResponse struct {
	InvDiscId   int        `json:"inv_disc_id"`
	InvDiscCode string     `json:"inv_disc_code"`
	InvDiscName string     `json:"inv_disc_name"`
	IsActive    bool       `json:"is_active"`
	UpdatedBy   *int64     `json:"updated_by"`
	UpdatedAt   *time.Time `json:"updated_at"`
}
type InvoiceDiscListResponse struct {
	InvDiscId     int        `json:"inv_disc_id"`
	InvDiscCode   string     `json:"inv_disc_code"`
	InvDiscName   string     `json:"inv_disc_name"`
	IsActive      bool       `json:"is_active"`
	UpdatedBy     *int64     `json:"updated_by"`
	UpdatedAt     *time.Time `json:"updated_at"`
	UpdatedByName string     `json:"updated_by_name"`
}

type InvDiscDet struct {
	RowNo    int     `json:"row_no"`
	MinValue float64 `json:"min_value"`
	MaxValue float64 `json:"max_value"`
	DiscPerc float64 `json:"disc_perc"`
}

type InvoiceDiscDetailsResponse struct {
	InvDiscId   int          `json:"inv_disc_id"`
	InvDiscCode string       `json:"inv_disc_code"`
	InvDiscName string       `json:"inv_disc_name"`
	IsActive    bool         `json:"is_active"`
	UpdatedBy   *int64       `json:"updated_by"`
	UpdatedAt   *time.Time   `json:"updated_at"`
	Details     []InvDiscDet `json:"details"`
}

type InvDiscDetReq struct {
	RowNo    int     `json:"row_no" validate:"required"`
	MinValue float64 `json:"min_value" validate:"required"`
	MaxValue float64 `json:"max_value" validate:"required"`
	DiscPerc float64 `json:"disc_perc" validate:"required"`
}

type CreateInvoiceDiscBody struct {
	CustId      string           `json:"cust_id" validate:"required,max=10"`
	CreatedBy   int64            `json:"created_by" validate:"required"`
	InvDiscCode string           `json:"inv_disc_code" validate:"required,max=5,alphanumericSpace"`
	InvDiscName string           `json:"inv_disc_name" validate:"required,max=150"`
	IsActive    bool             `json:"is_active"`
	Details     []*InvDiscDetReq `json:"details" validate:"min=1,dive"`
}

type DetailInvoiceDiscParams struct {
	InvDiscId int `params:"inv_disc_id" validate:"required"`
}

type UpdateInvoiceDiscParams struct {
	InvDiscId int `params:"inv_disc_id" validate:"required"`
}

type DeleteInvoiceDiscParams struct {
	InvDiscId int `params:"inv_disc_id" validate:"required"`
}

type UpdateInvoiceDiscRequest struct {
	CustId      string           `json:"cust_id" validate:"required,max=10"`
	UpdatedBy   int64            `json:"updated_by" validate:"required"`
	InvDiscCode string           `json:"inv_disc_code,omitempty" validate:"max=5,alphanumericSpace"`
	InvDiscName string           `json:"inv_disc_name,omitempty" validate:"max=150"`
	IsActive    *bool            `json:"is_active,omitempty"`
	Details     []*InvDiscDetReq `json:"details" validate:"min=1,dive"`
}
