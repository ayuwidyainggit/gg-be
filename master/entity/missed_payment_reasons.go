package entity

import (
	"time"
)

type MissedPaymentReasonsQueryFilter struct {
	CustId                 string
	ParentCustId           string
	Page                   int    `query:"page"`
	Limit                  int    `query:"limit" validate:"required"`
	Query                  string `query:"q"`
	Mode                   string `query:"mode"`
	Sort                   string `query:"sort"`
	IsActive               *int   `query:"is_active"`
	MissedPaymentReasonsId int    `query:"missed_payment_reasons_id"`
}

type MissedPaymentReasonsResponse struct {
	MissedPaymentReasonsId   int        `json:"missed_payment_reasons_id"`
	MissedPaymentReasonsCode string     `json:"missed_payment_reasons_code"`
	MissedPaymentReasonsName string     `json:"missed_payment_reasons_name"`
	ImageUrl                 string     `json:"image_url"`
	IsActive                 bool       `json:"is_active"`
	UpdatedBy                *int64     `json:"updated_by"`
	UpdatedByName            string     `json:"updated_by_name"`
	UpdatedAt                *time.Time `json:"updated_at"`
}
type MissedPaymentReasonsListResponse struct {
	MissedPaymentReasonsId   int        `json:"missed_payment_reasons_id"`
	MissedPaymentReasonsCode string     `json:"missed_payment_reasons_code"`
	MissedPaymentReasonsName string     `json:"missed_payment_reasons_name"`
	ImageUrl                 string     `json:"image_url"`
	IsActive                 bool       `json:"is_active"`
	UpdatedBy                *int64     `json:"updated_by"`
	UpdatedAt                *time.Time `json:"updated_at"`
	UpdatedByName            string     `json:"updated_by_name"`
}

type MissedPaymentReasonsLookupResponse struct {
	MissedPaymentReasonsId   int    `json:"missed_payment_reasons_id"`
	MissedPaymentReasonsCode string `json:"missed_payment_reasons_code"`
	MissedPaymentReasonsName string `json:"missed_payment_reasons_name"`
}

type CreateMissedPaymentReasonsBody struct {
	CustId                   string `json:"cust_id" validate:"required,max=10"`
	CreatedBy                int64  `json:"created_by" validate:"required"`
	MissedPaymentReasonsCode string `json:"missed_payment_reasons_code" validate:"required,max=50,alphanumericSpace"`
	MissedPaymentReasonsName string `json:"missed_payment_reasons_name" validate:"required,max=50,alphanumericSpace"`
	ImageUrl                 string `json:"image_url"`
	IsActive                 bool   `json:"is_active"`
}

type UpdateMissedPaymentReasonsRequest struct {
	CustId                   string `json:"cust_id" validate:"required,max=10"`
	UpdatedBy                int64  `json:"updated_by" validate:"required"`
	MissedPaymentReasonsCode string `json:"missed_payment_reasons_code,omitempty" validate:"max=50,omitempty,alphanumericSpace"`
	MissedPaymentReasonsName string `json:"missed_payment_reasons_name,omitempty" validate:"max=50,omitempty,alphanumericSpace"`
	ImageUrl                 string `json:"image_url"`
	IsActive                 *bool  `json:"is_active,omitempty"`
}

type DetailMissedPaymentReasonsParams struct {
	MissedPaymentReasonsId int `params:"missed_payment_reasons_id" validate:"required"`
}

type UpdateMissedPaymentReasonsParams struct {
	MissedPaymentReasonsId int `params:"missed_payment_reasons_id" validate:"required"`
}

type DeleteMissedPaymentReasonsParams struct {
	MissedPaymentReasonsId int `params:"missed_payment_reasons_id" validate:"required"`
}
