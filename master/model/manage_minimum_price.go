package model

import (
	"time"
)

type ManageMinimumPrice struct {
	CustId                   string    `db:"cust_id" json:"cust_id"`
	ManageMinimumPriceId     *int      `db:"manage_minimum_price_id" json:"manage_minimum_price_id"`
	BasePrice                *int      `db:"base_price" json:"base_price"`
	LimitAction              *int      `db:"limit_action" json:"limit_action"`
	Threshold                *float64  `db:"threshold" json:"threshold"`
	StatusManageMinimumPrice *int      `db:"status_manage_minimum_price" json:"status_manage_minimum_price"`
	ProId                    *int      `db:"pro_id" json:"pro_id"`
	Price1                   *float64  `db:"price1" json:"price1"`
	Price2                   *float64  `db:"price2" json:"price2"`
	Price3                   *float64  `db:"price3" json:"price3"`
	Price4                   *float64  `db:"price4" json:"price4"`
	Price5                   *float64  `db:"price5" json:"price5"`
	PriceMinimum1            *float64  `db:"price1_minimum" json:"price1_minimum"`
	PriceMinimum2            *float64  `db:"price2_minimum" json:"price2_minimum"`
	PriceMinimum3            *float64  `db:"price3_minimum" json:"price3_minimum"`
	PriceMinimum4            *float64  `db:"price4_minimum" json:"price4_minimum"`
	PriceMinimum5            *float64  `db:"price5_minimum" json:"price5_minimum"`
	UnitId1                  *string   `db:"unit_id1" json:"unit_id1"`
	UnitId2                  *string   `db:"unit_id2" json:"unit_id2"`
	UnitId3                  *string   `db:"unit_id3" json:"unit_id3"`
	UnitId4                  *string   `db:"unit_id4" json:"unit_id4"`
	UnitId5                  *string   `db:"unit_id5" json:"unit_id5"`
	ConvUnit2                *int      `db:"conv_unit2" json:"conv_unit2"`
	ConvUnit3                *int      `db:"conv_unit3" json:"conv_unit3"`
	ConvUnit4                *int      `db:"conv_unit4" json:"conv_unit4"`
	ConvUnit5                *int      `db:"conv_unit5" json:"conv_unit5"`
	CreatedBy                *int64    `db:"created_by" json:"created_by"`
	CreatedAt                time.Time `db:"created_at" json:"created_at"`
	UpdatedBy                *int64    `db:"updated_by" json:"updated_by"`
	UpdatedAt                time.Time `db:"updated_at" json:"updated_at"`
	DeletedBy                *int64    `db:"deleted_by" json:"deleted_by"`
	DeletedAt                time.Time `db:"deleted_at" json:"deleted_at"`
}

func (ManageMinimumPrice) TableName() string {
	return "mst.manage_minimum_price"
}

type ManageMinimumPriceRead struct {
	ManageMinimumPrice
	// BasePriceName                string  `db:"base_price_name" json:"base_price_name"`
	// LimitActionName              string  `db:"limit_action_name" json:"limit_action_name"`
	StatusManageMinimumPriceName *string `db:"status_manage_minimum_price_name" json:"status_manage_minimum_price_name"`
	ProductName                  *string `db:"pro_name" json:"pro_name"`
	ProductCode                  *string `db:"pro_code" json:"pro_code"`
}

func (ManageMinimumPriceRead) TableName() string {
	return "mst.manage_minimum_price"
}

type ManageMinimumPriceUpdate struct {
	BasePrice     *int     `db:"base_price" json:"base_price"`
	LimitAction   *int     `db:"limit_action" json:"limit_action"`
	Threshold     *float64 `db:"threshold" json:"threshold"`
	PriceMinimum1 *float64 `sql:"price1_minimum" json:"price1_minimum"`
	PriceMinimum2 *float64 `sql:"price2_minimum" json:"price2_minimum"`
	PriceMinimum3 *float64 `sql:"price3_minimum" json:"price3_minimum"`
	PriceMinimum4 *float64 `sql:"price4_minimum" json:"price4_minimum"`
	PriceMinimum5 *float64 `sql:"price5_minimum" json:"price5_minimum"`
	UpdatedBy     *int64   `db:"updated_by" json:"updated_by"`
}

func (ManageMinimumPriceUpdate) TableName() string {
	return "mst.manage_minimum_price"
}
