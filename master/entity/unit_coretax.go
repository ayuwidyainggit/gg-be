package entity

import (
	"time"
)

type UnitCoreTaxResponse struct {
	UnitIdCoreTax   string     `json:"unit_id_coretax"`
	UnitNameCoreTax string     `json:"unit_name_coretax"`
	IsActive        bool       `json:"is_active"`
	UpdatedBy       *int64     `json:"updated_by"`
	UpdatedAt       *time.Time `json:"updated_at"`
}
type UnitCoreTaxListResponse struct {
	UnitIdCoreTax   string     `json:"unit_id_coretax"`
	UnitNameCoreTax string     `json:"unit_name_coretax"`
	IsActive        bool       `json:"is_active"`
	UpdatedBy       *int64     `json:"updated_by"`
	UpdatedAt       *time.Time `json:"updated_at"`
	UpdatedByName   string     `json:"updated_by_name"`
}

type CreateUnitCoreTaxBody struct {
	CustId          string `json:"cust_id" validate:"required,max=10"`
	CreatedBy       int64  `json:"created_by" validate:"required"`
	UnitIdCoreTax   string `json:"core_tax_code" validate:"required,max=15,alphanumericSpace"`
	UnitNameCoreTax string `json:"unit_name_coretax" validate:"required,max=15,alphanumericSpace"`
	IsActive        bool   `json:"is_active"`
}

type DetailUnitCoreTaxParams struct {
	UnitIdCoreTax string `params:"unit_id_coretax" validate:"required"`
}

type UpdateUnitCoreTaxParams struct {
	UnitIdCoreTax string `params:"unit_id_coretax" validate:"required"`
}

type DeleteUnitCoreTaxParams struct {
	UnitIdCoreTax string `params:"unit_id_coretax" validate:"required"`
}

type UpdateUnitCoreTaxRequest struct {
	CustId          string `json:"cust_id" validate:"required,max=10"`
	UpdatedBy       int64  `json:"updated_by" validate:"required"`
	UnitIdCoreTax   string `json:"core_tax_code,omitempty" validate:"required,max=5,omitempty,alphanumericSpace"`
	UnitNameCoreTax string `json:"unit_name_coretax,omitempty" validate:"max=15,omitempty,alphanumericSpace"`
	IsActive        *bool  `json:"is_active,omitempty"`
}
