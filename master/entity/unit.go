package entity

import (
	"time"
)

type UnitResponse struct {
	UnitId          string     `json:"unit_id"`
	UnitName        string     `json:"unit_name"`
	UnitIdCoreTax   *string    `json:"unit_id_coretax"`
	UnitNameCoreTax *string    `json:"unit_name_coretax"`
	IsActive        bool       `json:"is_active"`
	UpdatedBy       *int64     `json:"updated_by"`
	UpdatedAt       *time.Time `json:"updated_at"`
}
type UnitListResponse struct {
	UnitId          string     `json:"unit_id"`
	UnitName        string     `json:"unit_name"`
	UnitIdCoreTax   *string    `json:"unit_id_coretax"`
	UnitNameCoreTax *string    `json:"unit_name_coretax"`
	IsActive        bool       `json:"is_active"`
	UpdatedBy       *int64     `json:"updated_by"`
	UpdatedAt       *time.Time `json:"updated_at"`
	UpdatedByName   string     `json:"updated_by_name"`
}

type CreateUnitBody struct {
	CustId        string `json:"cust_id" validate:"required,max=10"`
	CreatedBy     int64  `json:"created_by" validate:"required"`
	UnitId        string `json:"unit_id" validate:"required,max=5,alphanumericSpace"`
	UnitName      string `json:"unit_name" validate:"required,max=15,alphanumericSpace"`
	UnitIdCoreTax string `json:"unit_id_coretax" validate:"required,max=15"`
	IsActive      bool   `json:"is_active"`
}

type DetailUnitParams struct {
	UnitId string `params:"unit_id" validate:"required"`
}

type UpdateUnitParams struct {
	UnitId string `params:"unit_id" validate:"required"`
}

type DeleteUnitParams struct {
	UnitId string `params:"unit_id" validate:"required"`
}

type UpdateUnitRequest struct {
	CustId        string `json:"cust_id" validate:"required,max=10"`
	UpdatedBy     int64  `json:"updated_by" validate:"required"`
	UnitId        string `json:"unit_id,omitempty" validate:"required,max=5,omitempty,alphanumericSpace"`
	UnitName      string `json:"unit_name,omitempty" validate:"max=15,omitempty,alphanumericSpace"`
	UnitIdCoreTax string `json:"unit_id_coretax,omitempty" validate:"max=15"`
	IsActive      *bool  `json:"is_active,omitempty"`
}
