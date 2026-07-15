package model

import "time"

type MOutletCode struct {
	Id             string    `db:"id" json:"id"`
	CustId         string    `db:"cust_id" json:"cust_id"`
	SerialCode     string    `db:"serial_code" json:"serial_code"`
	YearCode       int       `db:"year_code" json:"year_code"`
	LastSequenceNo string    `db:"last_sequence_no" json:"last_sequence_no"`
	Status         string    `db:"status" json:"status"`
	CreatedAt      time.Time `db:"created_at" json:"created_at"`
	CreatedBy      string    `db:"created_by" json:"created_by"`
	UpdatedAt      time.Time `db:"updated_at" json:"updated_at"`
	UpdatedBy      string    `db:"updated_by" json:"updated_by"`
}
