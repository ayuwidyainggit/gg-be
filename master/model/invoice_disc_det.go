package model

import (
	"time"
)

type InvoiceDiscDet struct {
	CustId    string  `db:"cust_id" structs:"cust_id"`
	InvDiscId int     `db:"inv_disc_id" structs:"inv_disc_id"`
	RowNo     int     `db:"row_no" structs:"row_no"`
	MinValue  float64 `db:"min_value" structs:"min_value"`
	MaxValue  float64 `db:"max_value" structs:"max_value"`
	DiscPerc  float64 `db:"disc_perc" structs:"disc_perc"`
}

type InvoiceDiscDetUpdate struct {
	InvDiscCode *string    `json:"inv_disc_code,omitempty" sql:"inv_disc_code"`
	InvDiscName *string    `json:"inv_disc_name,omitempty" sql:"inv_disc_name"`
	IsActive    *bool      `json:"is_active" sql:"is_active"`
	UpdatedAt   *time.Time `json:"updated_at" sql:"updated_at"`
	UpdatedBy   *int64     `json:"updated_by" sql:"updated_by"`
}
