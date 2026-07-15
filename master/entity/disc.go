package entity

import (
	"time"
)

type DiscResponse struct {
	DiscId        int64      `json:"disc_id"`
	DiscCode      string     `json:"disc_code"`
	DiscName      string     `json:"disc_name"`
	StartDate     *time.Time `json:"start_date"`
	EndDate       *time.Time `json:"end_date"`
	RangeType     int        `json:"range_type"`
	IsMultiple    bool       `json:"is_multiple"`
	PurchaseLimit float64    `json:"purchase_limit"`
	DiscType      int        `json:"disc_type"`
	DiscPerc      float64    `json:"disc_perc"`
	DiscValue     float64    `json:"disc_value"`
	IsActive      bool       `json:"is_active"`
	UpdatedBy     *int64     `json:"updated_by"`
	UpdatedAt     *time.Time `json:"updated_at"`
	Details       []DiscDet  `json:"details,omitempty"`
}
type DiscListResponse struct {
	DiscId        int64     `json:"disc_id"`
	DiscCode      string    `json:"disc_code"`
	DiscName      string    `json:"disc_name"`
	StartDate     *string   `json:"start_date"`
	EndDate       *string   `json:"end_date"`
	RangeType     int       `json:"range_type"`
	IsMultiple    bool      `json:"is_multiple"`
	PurchaseLimit float64   `json:"purchase_limit"`
	DiscType      int       `json:"disc_type"`
	DiscPerc      float64   `json:"disc_perc"`
	DiscValue     float64   `json:"disc_value"`
	IsActive      bool      `json:"is_active"`
	UpdatedBy     *int64    `json:"updated_by"`
	UpdatedAt     *string   `json:"updated_at"`
	UpdatedByName *string   `json:"updated_by_name"`
	Details       []DiscDet `json:"details,omitempty"`
}
type CreateDiscBody struct {
	CustId        string              `json:"cust_id" validate:"required,max=10"`
	CreatedBy     int64               `json:"created_by" validate:"required"`
	DiscCode      string              `json:"disc_code" validate:"required,max=50,alphanumericSpace"`
	DiscName      string              `json:"disc_name" validate:"max=100,omitempty"`
	StartDate     string              `json:"start_date"`
	EndDate       string              `json:"end_date"`
	RangeType     *int                `json:"range_type"`
	IsMultiple    bool                `json:"is_multiple"`
	PurchaseLimit float64             `json:"purchase_limit"`
	DiscType      int                 `json:"disc_type"`
	DiscPerc      float64             `json:"disc_perc"`
	DiscValue     float64             `json:"disc_value"`
	IsActive      bool                `json:"is_active"`
	Details       []CreateDiscDetBody `json:"details"`
}

type DetailDiscParams struct {
	DiscId int64 `params:"disc_id" validate:"required"`
}

type UpdateDiscParams struct {
	DiscId int64 `params:"disc_id" validate:"required"`
}

type DeleteDiscParams struct {
	DiscId int `params:"disc_id" validate:"required"`
}

type UpdateDiscRequest struct {
	CustId        string              `json:"cust_id" validate:"required,max=10"`
	UpdatedBy     int64               `json:"updated_by" validate:"required"`
	DiscCode      string              `json:"disc_code,omitempty" validate:"max=50,alphanumericSpace"`
	DiscName      string              `json:"disc_name,omitempty" validate:"max=100,omitempty"`
	StartDate     string              `json:"start_date,omitempty" validate:"omitempty"`
	EndDate       string              `json:"end_date,omitempty" validate:"omitempty"`
	RangeType     int                 `json:"range_type,omitempty" validate:"omitempty"`
	IsMultiple    *bool               `json:"is_multiple,omitempty" validate:"omitempty"`
	PurchaseLimit float64             `json:"purchase_limit,omitempty" validate:"omitempty"`
	DiscType      int                 `json:"disc_type,omitempty" validate:"omitempty"`
	DiscPerc      float64             `json:"disc_perc,omitempty" validate:"omitempty"`
	DiscValue     float64             `json:"disc_value,omitempty" validate:"omitempty"`
	IsActive      *bool               `json:"is_active,omitempty"`
	Details       []UpdateDiscDetBody `json:"details"`
}
