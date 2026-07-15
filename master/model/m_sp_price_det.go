package model

import "time"

type MSpPriceDet struct {
	CustID        string     `db:"cust_id" json:"cust_id"`
	SpPriceDetID  string     `db:"sp_price_det_id" json:"sp_price_det_id"`
	RefType       int        `db:"ref_type" json:"ref_type"`
	SpPriceID     string     `db:"sp_price_id" json:"sp_price_id"`
	RefID         int64      `db:"ref_id" json:"ref_id"`
	SellPrice1    float64    `db:"sell_price1" json:"sell_price1"`
	SellPrice2    float64    `db:"sell_price2" json:"sell_price2"`
	SellPrice3    float64    `db:"sell_price3" json:"sell_price3"`
	NewSellPrice1 float64    `db:"new_sell_price1" json:"new_sell_price1"`
	NewSellPrice2 float64    `db:"new_sell_price2" json:"new_sell_price2"`
	NewSellPrice3 float64    `db:"new_sell_price3" json:"new_sell_price3"`
	CreatedBy     string     `db:"created_by" json:"created_by"`
	CreatedAt     *time.Time `db:"created_at" json:"created_at"`
	UpdatedBy     string     `db:"updated_by" json:"updated_by"`
	UpdatedAt     *time.Time `db:"updated_at" json:"updated_at"`
}

type MSpPriceDetView struct {
	CustID        string     `db:"cust_id" json:"cust_id"`
	SpPriceDetID  string     `db:"sp_price_det_id" json:"sp_price_det_id"`
	RefType       int        `db:"ref_type" json:"ref_type"`
	SpPriceID     string     `db:"sp_price_id" json:"sp_price_id"`
	RefID         *int64     `db:"ref_id" json:"ref_id"`
	RefCode       *string    `db:"ref_code" json:"ref_code"`
	RefName       *string    `db:"ref_name" json:"ref_name"`
	SellPrice1    *float64   `db:"sell_price1" json:"sell_price1"`
	SellPrice2    *float64   `db:"sell_price2" json:"sell_price2"`
	SellPrice3    *float64   `db:"sell_price3" json:"sell_price3"`
	NewSellPrice1 *float64   `db:"new_sell_price1" json:"new_sell_price1"`
	NewSellPrice2 *float64   `db:"new_sell_price2" json:"new_sell_price2"`
	NewSellPrice3 *float64   `db:"new_sell_price3" json:"new_sell_price3"`
	CreatedBy     string     `db:"created_by" json:"created_by"`
	CreatedAt     *time.Time `db:"created_at" json:"created_at"`
	UpdatedBy     string     `db:"updated_by" json:"updated_by"`
	UpdatedAt     *time.Time `db:"updated_at" json:"updated_at"`
}

type MSpPriceUpdateDet struct {
	RefType       *int64     `db:"ref_type" json:"ref_type"`
	RefID         *int64     `db:"ref_id" json:"ref_id"`
	UnitIdS       *string    `db:"unit_id_s" json:"unit_id_s"`
	UnitIdM       *string    `db:"unit_id_m" json:"unit_id_m"`
	UnitIdL       *string    `db:"unit_id_l" json:"unit_id_l"`
	SellPriceS    *float64   `db:"sell_price_s" json:"sell_price_s"`
	SellPriceM    *float64   `db:"sell_price_m" json:"sell_price_m"`
	SellPriceL    *float64   `db:"sell_price_l" json:"sell_price_l"`
	NewSellPriceS *float64   `db:"new_sell_price_s" json:"new_sell_price_s"`
	NewSellPriceM *float64   `db:"new_sell_price_m" json:"new_sell_price_m"`
	NewSellPriceL *float64   `db:"new_sell_price_l" json:"new_sell_price_l"`
	CreatedBy     *string    `db:"created_by" json:"created_by"`
	CreatedAt     *time.Time `db:"created_at" json:"created_at"`
	UpdatedBy     *string    `db:"updated_by" json:"updated_by"`
	UpdatedAt     *time.Time `db:"updated_at" json:"updated_at"`
}

type MSpPriceDetPublish struct {
	CustID        string     `db:"cust_id" json:"cust_id"`
	SpPriceDetID  string     `db:"sp_price_det_id" json:"sp_price_det_id"`
	RefType       int        `db:"ref_type" json:"ref_type"`
	SpPriceID     string     `db:"sp_price_id" json:"sp_price_id"`
	RefID         *int64     `db:"ref_id" json:"ref_id"`
	RefCode       *string    `db:"ref_code" json:"ref_code"`
	RefName       *string    `db:"ref_name" json:"ref_name"`
	SellPrice1    *float64   `db:"sell_price1" json:"sell_price1"`
	SellPrice2    *float64   `db:"sell_price2" json:"sell_price2"`
	SellPrice3    *float64   `db:"sell_price3" json:"sell_price3"`
	NewSellPrice1 *float64   `db:"new_sell_price1" json:"new_sell_price1"`
	NewSellPrice2 *float64   `db:"new_sell_price2" json:"new_sell_price2"`
	NewSellPrice3 *float64   `db:"new_sell_price3" json:"new_sell_price3"`
	CreatedBy     string     `db:"created_by" json:"created_by"`
	CreatedAt     *time.Time `db:"created_at" json:"created_at"`
	UpdatedBy     string     `db:"updated_by" json:"updated_by"`
	UpdatedAt     *time.Time `db:"updated_at" json:"updated_at"`
	ProID         int64      `db:"pro_id" json:"pro_id"`
	StartDate     *time.Time `db:"start_date" json:"start_date"`
	EndDate       *time.Time `db:"end_date" json:"end_date"`
}
