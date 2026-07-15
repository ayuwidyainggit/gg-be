package model

import "time"

type OutletCr struct {
	CustId     string     `db:"cust_id" json:"cust_id"`
	OutletCrId int64      `db:"outlet_cr_id" json:"outlet_cr_id"`
	OutletId   int64      `db:"outlet_id" json:"outlet_id"`
	Source     int        `db:"source" json:"source"`
	Status     int        `db:"status" json:"status"`
	CreatedBy  *int64     `db:"created_by,omitempty" json:"created_by"`
	CreatedAt  time.Time  `db:"created_at" json:"created_at"`
	UpdatedBy  *int64     `db:"updated_by,omitempty" json:"updated_by"`
	UpdatedAt  *time.Time `db:"updated_at,omitempty" json:"updated_at"`
	ApprovalBy *int64     `db:"approval_by,omitempty" json:"approval_by"`
	ApprovalAt *time.Time `db:"approval_at,omitempty" json:"approval_at"`
}

type OutletCrDet struct {
	OutletCrDetId int64   `db:"outlet_cr_det_id" json:"outlet_cr_det_id"`
	OutletCrId    int64   `db:"outlet_cr_id" json:"outlet_cr_id"`
	FieldName     string  `db:"field_name" json:"field_name"`
	OldValue      *string `db:"old_value,omitempty" json:"old_value"`
	NewValue      *string `db:"new_value,omitempty" json:"new_value"`
}

type OutletCrList struct {
	OutletCrId  int64     `db:"outlet_cr_id" json:"outlet_cr_id"`
	OutletId    int64     `db:"outlet_id" json:"outlet_id"`
	OutletCode  string    `db:"outlet_code" json:"outlet_code"`
	OutletName  string    `db:"outlet_name" json:"outlet_name"`
	CurrentLong *string   `db:"current_long" json:"current_long"`
	CurrentLat  *string   `db:"current_lat" json:"current_lat"`
	NewLong     *string   `db:"new_long" json:"new_long"`
	NewLat      *string   `db:"new_lat" json:"new_lat"`
	Source      int       `db:"source" json:"source"`
	Status      int       `db:"status" json:"status"`
	StatusDesc  string    `db:"status_desc" json:"status_desc"`
	RequestBy   *string   `db:"request_by" json:"request_by"`
	RequestDate time.Time `db:"request_date" json:"request_date"`
}
