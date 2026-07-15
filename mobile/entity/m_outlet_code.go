package entity

import "time"

type MOutletCodeResponse struct {
	Id             string    `json:"id"`
	CustId         string    `json:"cust_id"`
	SerialCode     string    `json:"serial_code"`
	YearCode       int       `json:"year_code"`
	LastSequenceNo string    `json:"last_sequence_no"`
	Status         string    `json:"status"`
	CreatedAt      time.Time `json:"created_at"`
	CreatedBy      string    `json:"created_by"`
	UpdatedAt      time.Time `json:"updated_at"`
	UpdatedBy      string    `json:"updated_by"`
}
