package model

import (
	"time"
)

type MTransactionPrice struct {
	CustID             string     `db:"cust_id" json:"cust_id"`
	TransactionPriceID string     `db:"transaction_price_id" json:"transaction_price_id"`
	ProID              int64      `db:"pro_id" json:"pro_id"`
	PurchPrice1        float64    `db:"purch_price1" json:"purch_price1"`
	PurchPrice2        float64    `db:"purch_price2" json:"purch_price2"`
	PurchPrice3        float64    `db:"purch_price3" json:"purch_price3"`
	SellPrice1         float64    `db:"sell_price1" json:"sell_price1"`
	SellPrice2         float64    `db:"sell_price2" json:"sell_price2"`
	SellPrice3         float64    `db:"sell_price3" json:"sell_price3"`
	Source             int        `db:"source" json:"source"`
	CreatedBy          string     `db:"created_by" json:"created_by"`
	CreatedAt          time.Time  `db:"created_at" json:"created_at"`
	StartDate          *time.Time `db:"start_date" json:"start_date"`
	EndDate            *time.Time `db:"end_date" json:"end_date"`
	Coverage           string     `db:"coverage" json:"coverage"`
	DistributorID      int64      `db:"distributor_id" json:"distributor_id"`
	PriceGroupReff     int64      `db:"price_group_reff" json:"price_group_reff"`
	ReferenceID        string     `db:"reference_id" json:"reference_id"`
	OutletID           int64      `db:"outlet_id" json:"outlet_id"`
}
