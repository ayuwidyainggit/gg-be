package entity

import "errors"

var ErrOutletCodeDuplicate = errors.New("outlet code duplicate")

type CreateOutletCodeBody struct {
	SerialCode     string `json:"serial_code" validate:"required"`
	YearCode       int    `json:"year_code" validate:"required"`
	LastSequenceNo string `json:"last_sequence_no" validate:"required"`
}

type UpdateOutletCodeBody struct {
	SerialCode string `json:"serial_code" validate:"required"`
}

type UpdateOutletCodeStatusBody struct {
	Status string `json:"status" validate:"required"`
	Id     string `json:"id" validate:"required"`
}

type OutletCodeListFilter struct {
	Page   int      `query:"page"`
	Limit  int      `query:"limit"`
	Sort   string   `query:"sort"`
	Q      string   `query:"q"`
	CustId string   `query:"cust_id"`
	Status []string `query:"status"`
}

type OutletCodeItem struct {
	Id             string `json:"id"`
	CustId         string `json:"cust_id"`
	SerialCode     string `json:"serial_code"`
	YearCode       int    `json:"year_code"`
	LastSequenceNo string `json:"last_sequence_no"`
	Status         string `json:"status"`
	CreatedAt      string `json:"created_at"`
	CreatedBy      string `json:"created_by"`
	UpdatedAt      string `json:"updated_at"`
	UpdatedBy      string `json:"updated_by"`
}

type SetupOutletCheckData struct {
	Id             string `json:"id"`
	CustId         string `json:"cust_id"`
	SerialCode     string `json:"serial_code"`
	YearCode       int    `json:"year_code"`
	LastSequenceNo string `json:"last_sequence_no"`
}
