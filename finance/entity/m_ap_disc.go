package entity

import "time"

type CreateMApDiscBody struct {
	CustID     string   `json:"cust_id"`
	ProID      *int64   `json:"pro_id"`
	PurchPrice *float64 `json:"purch_price"`
	NettPrice  *float64 `json:"nett_price"`
	DiscP      *float64 `json:"disc_p"`
	CreatedBy  *int64   `json:"created_by"`
}

type MApDiscResponse struct {
	ApDiscID      int       `json:"ap_disc_id"`
	ProID         *int64    `json:"pro_id"`
	ProCode       string    `json:"pro_code"`
	ProName       string    `json:"pro_name"`
	PurchPrice    *float64  `json:"purch_price"`
	NettPrice     *float64  `json:"nett_price"`
	DiscP         *float64  `json:"disc_p"`
	UpdatedAt     time.Time `json:"updated_at"`
	UpdatedByName *string   `json:"updated_by_name"`
}
type MApDiscListResponse struct {
	ApDiscID      int       `json:"ap_disc_id"`
	ProID         *int64    `json:"pro_id"`
	ProCode       string    `json:"pro_code"`
	ProName       string    `json:"pro_name"`
	PurchPrice    *float64  `json:"purch_price"`
	NettPrice     *float64  `json:"nett_price"`
	DiscP         *float64  `json:"disc_p"`
	UpdatedAt     time.Time `json:"updated_at"`
	UpdatedByName *string   `json:"updated_by_name"`
}
type UpdateMApDiscBody struct {
	ApDiscID   int      `json:"ap_disc_id"`
	CustID     string   `json:"cust_id"`
	ProID      *int64   `json:"pro_id"`
	PurchPrice *float64 `json:"purch_price"`
	NettPrice  *float64 `json:"nett_price"`
	DiscP      *float64 `json:"disc_p"`
	CreatedBy  *int64   `json:"created_by"`
	UpdatedBy  int64    `json:"updated_by"`
}
type DetailMApDiscParams struct {
	ApDiscID int64 `params:"ap_disc_id" validate:"required"`
}

type UpdateMApDiscParams struct {
	ApDiscID int64 `params:"ap_disc_id" validate:"required"`
}
