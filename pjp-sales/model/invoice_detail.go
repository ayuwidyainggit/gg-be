package model

import "time"

type InvoiceDet struct {
	CustId      string     `gorm:"cust_id" json:"cust_id"`
	RoNo        string     `gorm:"ro_no" json:"ro_no"`
	SeqNo       int        `gorm:"seq_no" json:"seq_no"`
	RoDetID     *int       `gorm:"column:ro_det_id;primaryKey" json:"ro_det_id"`
	ProId       int        `gorm:"pro_id" json:"pro_id"`
	ItemType    int        `gorm:"item_type" json:"item_type"`
	Qty         *float64   `gorm:"qty" json:"qty"`
	Qty1        *float64   `gorm:"qty1" json:"qty1"`
	Qty2        *float64   `gorm:"qty2" json:"qty2"`
	Qty3        *float64   `gorm:"qty3" json:"qty3"`
	QtyStr      *string    `gorm:"qty_str" json:"qty_str"`
	PurchPrice1 *float64   `gorm:"purch_price1" json:"purch_price1"`
	PurchPrice2 *float64   `gorm:"purch_price2" json:"purch_price2"`
	PurchPrice3 *float64   `gorm:"purch_price3" json:"purch_price3"`
	SellPrice1  *float64   `gorm:"sell_price1" json:"sell_price1"`
	SellPrice2  *float64   `gorm:"sell_price2" json:"sell_price2"`
	SellPrice3  *float64   `gorm:"sell_price3" json:"sell_price3"`
	UnitId1     *float64   `gorm:"unit_id1" json:"unit_id1"`
	UnitId2     *float64   `gorm:"unit_id2" json:"unit_id2"`
	UnitId3     *float64   `gorm:"unit_id3" json:"unit_id3"`
	ConvUnit2   *float64   `gorm:"conv_unit2" json:"conv_unit2"`
	ConvUnit3   *float64   `gorm:"conv_unit3" json:"conv_unit3"`
	Amount      *float64   `gorm:"amount" json:"amount"`
	DiscValue   *float64   `gorm:"disc_value" json:"disc_value"`
	BatchNo     *string    `gorm:"batch_no" json:"batch_no"`
	ExpDate     *time.Time `gorm:"exp_date" json:"exp_date"`
}

func (InvoiceDet) TableName() string {
	return "sls.order_detail"
}

type InvoiceDetRead struct {
	CustId          string     `gorm:"cust_id" json:"cust_id"`
	RoNo            string     `gorm:"ro_no" json:"ro_no"`
	SeqNo           int        `gorm:"seq_no" json:"seq_no"`
	OrderDetID      *int       `gorm:"column:order_detail_id;primaryKey" json:"order_detail_id"`
	ProId           int        `gorm:"pro_id" json:"pro_id"`
	ProCode         string     `gorm:"column:pro_code" json:"pro_code"`
	ProName         string     `gorm:"column:pro_name" json:"pro_name"`
	ItemType        int        `gorm:"item_type" json:"item_type"`
	Qty1            float64    `gorm:"qty1" json:"qty1"`
	Qty2            float64    `gorm:"qty2" json:"qty2"`
	Qty3            float64    `gorm:"qty3" json:"qty3"`
	Qty4            float64    `gorm:"qty4" json:"qty4"`
	Qty5            float64    `gorm:"qty5" json:"qty5"`
	Qty1Final       float64    `gorm:"qty1_final" json:"qty1_final"`
	Qty2Final       float64    `gorm:"qty2_final" json:"qty2_final"`
	Qty3Final       float64    `gorm:"qty3_final" json:"qty3_final"`
	Qty4Final       float64    `gorm:"qty4_final" json:"qty4_final"`
	Qty5Final       float64    `gorm:"qty5_final" json:"qty5_final"`
	Volume1         float64    `gorm:"volume1" json:"volume1"`
	Volume2         float64    `gorm:"volume2" json:"volume2"`
	Volume3         float64    `gorm:"volume3" json:"volume3"`
	Weight1         float64    `gorm:"weight1" json:"weight1"`
	Weight2         float64    `gorm:"weight2" json:"weight2"`
	Weight3         float64    `gorm:"weight3" json:"weight3"`
	PurchPrice1     float64    `gorm:"purch_price1" json:"purch_price1"`
	PurchPrice2     float64    `gorm:"purch_price2" json:"purch_price2"`
	PurchPrice3     float64    `gorm:"purch_price3" json:"purch_price3"`
	SellPrice1      float64    `gorm:"sell_price1" json:"sell_price1"`
	SellPrice2      float64    `gorm:"sell_price2" json:"sell_price2"`
	SellPrice3      float64    `gorm:"sell_price3" json:"sell_price3"`
	SellPriceFinal1 float64    `gorm:"sell_price_final1" json:"sell_price_final1"`
	SellPriceFinal2 float64    `gorm:"sell_price_final2" json:"sell_price_final2"`
	SellPriceFinal3 float64    `gorm:"sell_price_final3" json:"sell_price_final3"`
	ConvUnit1       float64    `gorm:"conv_unit1" json:"conv_unit1"`
	ConvUnit2       float64    `gorm:"conv_unit2" json:"conv_unit2"`
	ConvUnit3       float64    `gorm:"conv_unit3" json:"conv_unit3"`
	UnitId1         string     `gorm:"unit_id1" json:"unit_id1"`
	UnitId2         string     `gorm:"unit_id2" json:"unit_id2"`
	UnitId3         string     `gorm:"unit_id3" json:"unit_id3"`
	Length          float64    `gorm:"length" json:"length"`
	Width           float64    `gorm:"width" json:"width"`
	Height          float64    `gorm:"height" json:"height"`
	Weight          float64    `gorm:"weight" json:"weight"`
	Volume          float64    `gorm:"volume" json:"volume"`
	Amount          float64    `gorm:"amount" json:"amount"`
	AmountFinal     float64    `gorm:"amount_final" json:"amount_final"`
	DiscValue       float64    `gorm:"disc_value" json:"disc_value"`
	DiscValueFinal  *float64   `gorm:"disc_value_final" json:"disc_value_final"`
	PromoFinal1     *float64   `gorm:"column:promo_final1" json:"promo_final1"`
	PromoFinal2     *float64   `gorm:"column:promo_final2" json:"promo_final2"`
	PromoFinal3     *float64   `gorm:"column:promo_final3" json:"promo_final3"`
	PromoFinal4     *float64   `gorm:"column:promo_final4" json:"promo_final4"`
	PromoFinal5     *float64   `gorm:"column:promo_final5" json:"promo_final5"`
	BatchNo         string     `gorm:"batch_no" json:"batch_no"`
	ExpDate         *time.Time `gorm:"exp_date" json:"exp_date"`
	PriceIncludePpn float64    `gorm:"price_include_ppn" json:"price_include_ppn"`
	PriceExcludePpn float64    `gorm:"price_exclude_ppn" json:"price_exclude_ppn"`
	NetValue        float64    `gorm:"net_value" json:"net_value"`
	Vat             float64    `gorm:"vat" json:"vat"`
	VatValue        float64    `gorm:"vat_value" json:"vat_value"`
	VatValueFinal   *float64   `gorm:"vat_value_final" json:"vat_value_final"`
}

func (InvoiceDetRead) TableName() string {
	return "sls.order_detail"
}
