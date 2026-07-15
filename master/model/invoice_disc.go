package model

import (
	"time"
)

type InvoiceDisc struct {
	CustId        string     `db:"cust_id"`
	InvDiscId     int        `db:"inv_disc_id"`
	InvDiscCode   string     `db:"inv_disc_code"`
	InvDiscName   string     `db:"inv_disc_name"`
	IsActive      bool       `db:"is_active"`
	IsDel         bool       `db:"is_del"`
	CreatedBy     *int64     `db:"created_by,omitempty"`
	CreatedAt     *time.Time `db:"created_at,omitempty"`
	UpdatedBy     *int64     `db:"updated_by,omitempty"`
	UpdatedAt     *time.Time `db:"updated_at,omitempty"`
	UpdatedByName *string    `json:"updated_by_name" db:"updated_by_name"`
	DeletedBy     *int64     `db:"deleted_by,omitempty"`
	DeletedAt     *time.Time `db:"deleted_at,omitempty"`
}

type InvoiceDiscUpdate struct {
	InvDiscCode *string    `json:"inv_disc_code,omitempty" sql:"inv_disc_code"`
	InvDiscName *string    `json:"inv_disc_name,omitempty" sql:"inv_disc_name"`
	IsActive    *bool      `json:"is_active" sql:"is_active"`
	UpdatedAt   *time.Time `json:"updated_at" sql:"updated_at"`
	UpdatedBy   *int64     `json:"updated_by" sql:"updated_by"`
}
