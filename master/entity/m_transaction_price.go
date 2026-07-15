package entity

import (
	"time"
)

type MTransactionPriceQueryFilter struct {
	CustID       string
	ParentCustID string
	Page         int    `query:"page"`
	Limit        int    `query:"limit"`
	Query        string `query:"q"`
	Mode         string `query:"mode"`
	Sort         string `query:"sort"`
}

type MTransactionPriceResp struct {
	CustID        string    `json:"cust_id,omitempty"`
	ParentCustID  string    `json:"parent_cust_id,omitempty"`
	Coverage      string    `json:"coverage"`
	EffectiveDate string    `json:"effective_date"`
	ProID         int64     `json:"pro_id"`
	ProName       string    `json:"pro_name"`
	ProCode       string    `json:"pro_code"`
	UnitID1       string    `json:"unit_id1"`
	UnitID2       string    `json:"unit_id2"`
	UnitID3       string    `json:"unit_id3"`
	ConvUnit2     int       `json:"conv_unit2"`
	ConvUnit3     int       `json:"conv_unit3"`
	PurchPrice1   float64   `json:"purch_price1"`
	PurchPrice2   float64   `json:"purch_price2"`
	PurchPrice3   float64   `json:"purch_price3"`
	SellPrice1    float64   `json:"sell_price1"`
	SellPrice2    float64   `json:"sell_price2"`
	SellPrice3    float64   `json:"sell_price3"`
	UpdatedBy     string    `json:"updated_by"`
	UpdatedAt     time.Time `json:"updated_at"`
	DistributorID int64     `json:"distributor_id"`
}
