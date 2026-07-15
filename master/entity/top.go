package entity

import (
	"time"
)

type TopResponse struct {
	Top           int        `json:"top"`
	IsActive      bool       `json:"is_active"`
	UpdatedBy     *int64     `json:"updated_by"`
	UpdatedByName *string    `json:"updated_by_name"`
	UpdatedAt     *time.Time `json:"updated_at"`
}

type CreateTopBody struct {
	CustId    string `json:"cust_id" validate:"required,max=10"`
	CreatedBy int64  `json:"created_by" validate:"required"`
	Top       int    `json:"top" validate:"required,gte=1"`
	IsActive  bool   `json:"is_active"`
}

type DetailTopParams struct {
	Top int `params:"top" validate:"required"`
}

type UpdateTopParams struct {
	Top int `params:"top" validate:"required"`
}

type DeleteTopParams struct {
	Top int `params:"top" validate:"required"`
}

type UpdateTopRequest struct {
	CustId    string `json:"cust_id" validate:"required,max=10"`
	UpdatedBy int64  `json:"updated_by" validate:"required"`
	Top       int    `json:"top,omitempty" validate:"required,gte=1"`
	IsActive  *bool  `json:"is_active,omitempty"`
}
