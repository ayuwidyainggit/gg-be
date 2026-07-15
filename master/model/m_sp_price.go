package model

import "time"

type MSpPrice struct {
	CustID        string     `db:"cust_id" json:"cust_id"`
	SpPriceID     string     `db:"sp_price_id" json:"sp_price_id"`
	StartDate     *time.Time `db:"start_date" json:"start_date"`
	EndDate       *time.Time `db:"end_date" json:"end_date"`
	PriceGrpID    int64      `db:"price_grp_id" json:"price_grp_id"`
	ProID         int64      `db:"pro_id" json:"pro_id"`
	UnitId1       string     `db:"unit_id1" json:"unit_id1"`
	UnitId2       string     `db:"unit_id2" json:"unit_id2"`
	UnitId3       string     `db:"unit_id3" json:"unit_id3"`
	SellPrice1    float64    `db:"sell_price1" json:"sell_price1"`
	SellPrice2    float64    `db:"sell_price2" json:"sell_price2"`
	SellPrice3    float64    `db:"sell_price3" json:"sell_price3"`
	NewSellPrice1 float64    `db:"new_sell_price1" json:"new_sell_price1"`
	NewSellPrice2 float64    `db:"new_sell_price2" json:"new_sell_price2"`
	NewSellPrice3 float64    `db:"new_sell_price3" json:"new_sell_price3"`
	ConvUnit2     float32    `db:"conv_unit2" json:"conv_unit2"`
	ConvUnit3     float32    `db:"conv_unit3" json:"conv_unit3"`
	Status        int        `db:"status" json:"status"`
	CreatedBy     string     `db:"created_by" json:"created_by"`
	CreatedAt     *time.Time `db:"created_at" json:"created_at"`
	UpdatedBy     string     `db:"updated_by" json:"updated_by"`
	UpdatedAt     *time.Time `db:"updated_at" json:"updated_at"`
}

type MSpPriceDetail struct {
	CustID        string     `db:"cust_id" json:"cust_id"`
	SpPriceID     string     `db:"sp_price_id" json:"sp_price_id"`
	StartDate     *time.Time `db:"start_date" json:"start_date"`
	EndDate       *time.Time `db:"end_date" json:"end_date"`
	PriceGrpID    int64      `db:"price_grp_id" json:"price_grp_id"`
	PriceGrpCode  string     `db:"price_grp_code" json:"price_grp_code"`
	PriceGrpName  string     `db:"price_grp_name" json:"price_grp_name"`
	ProID         int64      `db:"pro_id" json:"pro_id"`
	ProCode       string     `db:"pro_code" json:"pro_code"`
	ProName       string     `db:"pro_name" json:"pro_name"`
	UnitId1       string     `db:"unit_id1" json:"unit_id1"`
	UnitId2       string     `db:"unit_id2" json:"unit_id2"`
	UnitId3       string     `db:"unit_id3" json:"unit_id3"`
	UnitName1     string     `db:"unit_name1" json:"unit_name1"`
	UnitName2     string     `db:"unit_name2" json:"unit_name2"`
	UnitName3     string     `db:"unit_name3" json:"unit_name3"`
	SellPrice1    float64    `db:"sell_price1" json:"sell_price1"`
	SellPrice2    float64    `db:"sell_price2" json:"sell_price2"`
	SellPrice3    float64    `db:"sell_price3" json:"sell_price3"`
	NewSellPrice1 float64    `db:"new_sell_price1" json:"new_sell_price1"`
	NewSellPrice2 float64    `db:"new_sell_price2" json:"new_sell_price2"`
	NewSellPrice3 float64    `db:"new_sell_price3" json:"new_sell_price3"`
	ConvUnit2     float32    `db:"conv_unit2" json:"conv_unit2"`
	ConvUnit3     float32    `db:"conv_unit3" json:"conv_unit3"`
	Status        int        `db:"status" json:"status"`
	CreatedBy     string     `db:"created_by" json:"created_by"`
	CreatedAt     *time.Time `db:"created_at" json:"created_at"`
	UpdatedBy     string     `db:"updated_by" json:"updated_by"`
	UpdatedAt     *time.Time `db:"updated_at" json:"updated_at"`
}

type MSpPriceUpdate struct {
	StartDate     *string  `db:"start_date" json:"start_date"`
	EndDate       *string  `db:"end_date" json:"end_date"`
	PriceGrpID    *int64   `db:"price_grp_id" json:"price_grp_id"`
	ProID         *int64   `db:"pro_id" json:"pro_id"`
	UnitId1       *string  `db:"unit_id1" json:"unit_id1"`
	UnitId2       *string  `db:"unit_id2" json:"unit_id2"`
	UnitId3       *string  `db:"unit_id3" json:"unit_id3"`
	SellPrice1    *float64 `db:"sell_price1" json:"sell_price1"`
	SellPrice2    *float64 `db:"sell_price2" json:"sell_price2"`
	SellPrice3    *float64 `db:"sell_price3" json:"sell_price3"`
	NewSellPrice1 *float64 `db:"new_sell_price1" json:"new_sell_price1"`
	NewSellPrice2 *float64 `db:"new_sell_price2" json:"new_sell_price2"`
	NewSellPrice3 *float64 `db:"new_sell_price3" json:"new_sell_price3"`
	ConvUnit2     *float32 `db:"conv_unit2" json:"conv_unit2"`
	ConvUnit3     *float32 `db:"conv_unit3" json:"conv_unit3"`
	Status        *int     `db:"status" json:"status"`
	UpdatedBy     *string  `db:"updated_by" json:"updated_by"`
}

type MSpPriceWithDetails struct {
	CustID        string     `db:"cust_id" json:"cust_id"`
	SpPriceID     string     `db:"sp_price_id" json:"sp_price_id"`
	StartDate     *time.Time `db:"start_date" json:"start_date"`
	EndDate       *time.Time `db:"end_date" json:"end_date"`
	PriceGrpID    int64      `db:"price_grp_id" json:"price_grp_id"`
	PriceGrpCode  string     `db:"price_grp_code" json:"price_grp_code"`
	PriceGrpName  string     `db:"price_grp_name" json:"price_grp_name"`
	ProID         int64      `db:"pro_id" json:"pro_id"`
	ProCode       string     `db:"pro_code" json:"pro_code"`
	ProName       string     `db:"pro_name" json:"pro_name"`
	UnitId1       string     `db:"unit_id1" json:"unit_id1"`
	UnitId2       string     `db:"unit_id2" json:"unit_id2"`
	UnitId3       string     `db:"unit_id3" json:"unit_id3"`
	UnitName1     string     `db:"unit_name1" json:"unit_name1"`
	UnitName2     string     `db:"unit_name2" json:"unit_name2"`
	UnitName3     string     `db:"unit_name3" json:"unit_name3"`
	PurchPrice1   float64    `db:"purch_price1" json:"purch_price1"`
	PurchPrice2   float64    `db:"purch_price2" json:"purch_price2"`
	PurchPrice3   float64    `db:"purch_price3" json:"purch_price3"`
	SellPrice1    float64    `db:"sell_price1" json:"sell_price1"`
	SellPrice2    float64    `db:"sell_price2" json:"sell_price2"`
	SellPrice3    float64    `db:"sell_price3" json:"sell_price3"`
	NewSellPrice1 float64    `db:"new_sell_price1" json:"new_sell_price1"`
	NewSellPrice2 float64    `db:"new_sell_price2" json:"new_sell_price2"`
	NewSellPrice3 float64    `db:"new_sell_price3" json:"new_sell_price3"`
	ConvUnit2     float32    `db:"conv_unit2" json:"conv_unit2"`
	ConvUnit3     float32    `db:"conv_unit3" json:"conv_unit3"`
	Status        int        `db:"status" json:"status"`
	CreatedBy     string     `db:"created_by" json:"created_by"`
	CreatedAt     *time.Time `db:"created_at" json:"created_at"`
	UpdatedBy     string     `db:"updated_by" json:"updated_by"`
	UpdatedAt     *time.Time `db:"updated_at" json:"updated_at"`
}

type MSpPricePublish struct {
	Status    *int       `sql:"status" json:"status"`
	UpdatedAt *time.Time `sql:"updated_at" json:"updated_at"`
	UpdatedBy *string    `sql:"updated_by" json:"updated_by"`
}
