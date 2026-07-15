package model

import (
	"time"

	"github.com/lib/pq"
)

type MPrice struct {
	CustID         string        `db:"cust_id" json:"cust_id"`
	PriceID        string        `db:"price_id" json:"price_id"`
	Coverage       string        `db:"coverage" json:"coverage"`
	DistributorIDs pq.Int64Array `db:"distributor_ids" json:"distributor_ids"`
	EffectiveDate  *time.Time    `db:"effective_date" json:"effective_date"`
	ProID          int64         `db:"pro_id" json:"pro_id"`
	UnitID1        string        `db:"unit_id1" json:"unit_id1"`
	UnitID2        string        `db:"unit_id2" json:"unit_id2"`
	UnitID3        string        `db:"unit_id3" json:"unit_id3"`
	ProductName    string        `db:"pro_name" json:"pro_name"`
	ProductCode    string        `db:"pro_code" json:"pro_code"`
	UnitName1      *string       `db:"unit_name1"`
	UnitName2      *string       `db:"unit_name2"`
	UnitName3      *string       `db:"unit_name3"`
	ConvUnit2      int           `db:"conv_unit2" json:"conv_unit2"`
	ConvUnit3      int           `db:"conv_unit3" json:"conv_unit3"`
	PurchPrice1    float64       `db:"purch_price1" json:"purch_price1"`
	PurchPrice2    float64       `db:"purch_price2" json:"purch_price2"`
	PurchPrice3    float64       `db:"purch_price3" json:"purch_price3"`
	SellPrice1     float64       `db:"sell_price1" json:"sell_price1"`
	SellPrice2     float64       `db:"sell_price2" json:"sell_price2"`
	SellPrice3     float64       `db:"sell_price3" json:"sell_price3"`
	NewPurchPrice1 float64       `db:"new_purch_price1" json:"new_purch_price1"`
	NewPurchPrice2 float64       `db:"new_purch_price2" json:"new_purch_price2"`
	NewPurchPrice3 float64       `db:"new_purch_price3" json:"new_purch_price3"`
	NewSellPrice1  float64       `db:"new_sell_price1" json:"new_sell_price1"`
	NewSellPrice2  float64       `db:"new_sell_price2" json:"new_sell_price2"`
	NewSellPrice3  float64       `db:"new_sell_price3" json:"new_sell_price3"`
	Status         int           `db:"status" json:"status"`
	CreatedByID    *int64        `db:"created_by_id" json:"created_by_id"`
	CreatedBy      string        `db:"created_by" json:"created_by"`
	CreatedAt      time.Time     `db:"created_at" json:"created_at"`
	UpdatedByID    *int64        `db:"updated_by_id" json:"updated_by_id"`
	UpdatedBy      string        `db:"updated_by" json:"updated_by"`
	UpdatedAt      time.Time     `db:"updated_at" json:"updated_at"`
}

type MPriceUpdate struct {
	EffectiveDate  *string    `sql:"effective_date" json:"effective_date"`
	Coverage       *string    `sql:"coverage" json:"coverage"`
	DistributorIDs *[]int64   `sql:"distributor_ids" json:"distributor_ids"`
	ProID          *int64     `sql:"pro_id" json:"pro_id"`
	UnitID1        *string    `sql:"unit_id1" json:"unit_id1"`
	UnitID2        *string    `sql:"unit_id2" json:"unit_id2"`
	UnitID3        *string    `sql:"unit_id3" json:"unit_id3"`
	ConvUnit2      *int       `sql:"conv_unit2" json:"conv_unit2"`
	ConvUnit3      *int       `sql:"conv_unit3" json:"conv_unit3"`
	PurchPrice1    *float64   `sql:"purch_price1" json:"purch_price1"`
	PurchPrice2    *float64   `sql:"purch_price2" json:"purch_price2"`
	PurchPrice3    *float64   `sql:"purch_price3" json:"purch_price3"`
	SellPrice1     *float64   `sql:"sell_price1" json:"sell_price1"`
	SellPrice2     *float64   `sql:"sell_price2" json:"sell_price2"`
	SellPrice3     *float64   `sql:"sell_price3" json:"sell_price3"`
	NewPurchPrice1 *float64   `sql:"new_purch_price1" json:"new_purch_price1"`
	NewPurchPrice2 *float64   `sql:"new_purch_price2" json:"new_purch_price2"`
	NewPurchPrice3 *float64   `sql:"new_purch_price3" json:"new_purch_price3"`
	NewSellPrice1  *float64   `sql:"new_sell_price1" json:"new_sell_price1"`
	NewSellPrice2  *float64   `sql:"new_sell_price2" json:"new_sell_price2"`
	NewSellPrice3  *float64   `sql:"new_sell_price3" json:"new_sell_price3"`
	UpdatedAt      *time.Time `sql:"updated_at" json:"updated_at"`
	UpdatedByID    *int64     `sql:"updated_by_id" json:"updated_by_id"`
	UpdatedBy      *string    `sql:"updated_by" json:"updated_by"`
}

type MPriceDetail struct {
	CustID         string        `db:"cust_id" json:"cust_id"`
	PriceID        string        `db:"price_id" json:"price_id"`
	Coverage       string        `db:"coverage" json:"coverage"`
	EffectiveDate  *time.Time    `db:"effective_date" json:"effective_date"`
	ProID          int64         `db:"pro_id" json:"pro_id"`
	ProCode        string        `db:"pro_code" json:"pro_code"`
	ProName        string        `db:"pro_name" json:"pro_name"`
	UnitID1        string        `db:"unit_id1" json:"unit_id1"`
	UnitID2        string        `db:"unit_id2" json:"unit_id2"`
	UnitID3        string        `db:"unit_id3" json:"unit_id3"`
	UnitName1      string        `db:"unit_name1"`
	UnitName2      string        `db:"unit_name2"`
	UnitName3      string        `db:"unit_name3"`
	ConvUnit2      int           `db:"conv_unit2" json:"conv_unit2"`
	ConvUnit3      int           `db:"conv_unit3" json:"conv_unit3"`
	PurchPrice1    float64       `db:"purch_price1" json:"purch_price1"`
	PurchPrice2    float64       `db:"purch_price2" json:"purch_price2"`
	PurchPrice3    float64       `db:"purch_price3" json:"purch_price3"`
	SellPrice1     float64       `db:"sell_price1" json:"sell_price1"`
	SellPrice2     float64       `db:"sell_price2" json:"sell_price2"`
	SellPrice3     float64       `db:"sell_price3" json:"sell_price3"`
	NewPurchPrice1 float64       `db:"new_purch_price1" json:"new_purch_price1,omitempty"`
	NewPurchPrice2 float64       `db:"new_purch_price2" json:"new_purch_price2,omitempty"`
	NewPurchPrice3 float64       `db:"new_purch_price3" json:"new_purch_price3,omitempty"`
	NewSellPrice1  float64       `db:"new_sell_price1" json:"new_sell_price1,omitempty"`
	NewSellPrice2  float64       `db:"new_sell_price2" json:"new_sell_price2,omitempty"`
	NewSellPrice3  float64       `db:"new_sell_price3" json:"new_sell_price3,omitempty"`
	Status         int           `db:"status" json:"status"`
	CreatedByID    *int64        `db:"created_by_id,omitempty" json:"created_by_id"`
	CreatedBy      string        `db:"created_by,omitempty" json:"created_by"`
	CreatedAt      *time.Time    `db:"created_at,omitempty" json:"created_at"`
	UpdatedByID    *int64        `db:"updated_by_id,omitempty" json:"updated_by_id"`
	UpdatedBy      string        `db:"updated_by,omitempty" json:"updated_by"`
	UpdatedAt      *time.Time    `db:"updated_at,omitempty" json:"updated_at"`
	DistributorIDs pq.Int64Array `db:"distributor_ids" json:"distributor_ids"`
}

type MPricePublish struct {
	Status      *int       `sql:"status" json:"status"`
	UpdatedAt   *time.Time `sql:"updated_at" json:"updated_at"`
	UpdatedByID *int64     `sql:"updated_by_id" json:"updated_by_id"`
	UpdatedBy   *string    `sql:"updated_by" json:"updated_by"`
}

type MPriceProductSnapshot struct {
	ProID         int64   `db:"pro_id"`
	ProCode       string  `db:"pro_code"`
	ProName       string  `db:"pro_name"`
	UnitID1       string  `db:"unit_id1"`
	UnitID2       string  `db:"unit_id2"`
	UnitID3       string  `db:"unit_id3"`
	ConvUnit2     int     `db:"conv_unit2"`
	ConvUnit3     int     `db:"conv_unit3"`
	PurchPrice1   float64 `db:"purch_price1"`
	PurchPrice2   float64 `db:"purch_price2"`
	PurchPrice3   float64 `db:"purch_price3"`
	SellPrice1    float64 `db:"sell_price1"`
	SellPrice2    float64 `db:"sell_price2"`
	SellPrice3    float64 `db:"sell_price3"`
	DistributorID *int64  `db:"distributor_id"`
	ParentProID   int64   `db:"parent_pro_id"`
}
