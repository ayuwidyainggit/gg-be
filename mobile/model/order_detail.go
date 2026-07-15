package model

import "time"

type OrderDetail struct {
	CustId        string     `gorm:"cust_id" json:"cust_id"`
	RoNo          string     `gorm:"ro_no" json:"ro_no"`
	SeqNo         int        `gorm:"seq_no" json:"seq_no"`
	OrderDetailID *int       `gorm:"column:order_detail_id;primaryKey" json:"order_detail_id"`
	ProId         int        `gorm:"pro_id" json:"pro_id"`
	ItemType      int        `gorm:"item_type" json:"item_type"`
	Qty           float64    `gorm:"qty" json:"qty"`
	QtyFinal      float64    `gorm:"qty_final" json:"qty_final"`
	QtyPo         float64    `gorm:"qty_po" json:"qty_po"`
	Qty1          *float64   `gorm:"qty1" json:"qty1"`
	Qty2          *float64   `gorm:"qty2" json:"qty2"`
	Qty3          *float64   `gorm:"qty3" json:"qty3"`
	Qty4          *float64   `gorm:"qty4" json:"qty4"`
	Qty5          *float64   `gorm:"qty5" json:"qty5"`
	Qty1Final     *float64   `gorm:"qty1_final" json:"qty1_final"`
	Qty2Final     *float64   `gorm:"qty2_final" json:"qty2_final"`
	Qty3Final     *float64   `gorm:"qty3_final" json:"qty3_final"`
	Qty4Final     *float64   `gorm:"qty4_final" json:"qty4_final"`
	Qty5Final     *float64   `gorm:"qty5_final" json:"qty5_final"`
	Qty1Stok      *float64   `gorm:"qty1_stok" json:"qty1_stok"`
	Qty2Stok      *float64   `gorm:"qty2_stok" json:"qty2_stok"`
	Qty3Stok      *float64   `gorm:"qty3_stok" json:"qty3_stok"`
	PurchPrice1   *float64   `gorm:"purch_price1" json:"purch_price1"`
	PurchPrice2   *float64   `gorm:"purch_price2" json:"purch_price2"`
	PurchPrice3   *float64   `gorm:"purch_price3" json:"purch_price3"`
	PurchPrice4   *float64   `gorm:"purch_price4" json:"purch_price4"`
	PurchPrice5   *float64   `gorm:"purch_price5" json:"purch_price5"`
	SellPrice1    *float64   `gorm:"sell_price1" json:"sell_price1"`
	SellPrice2    *float64   `gorm:"sell_price2" json:"sell_price2"`
	SellPrice3    *float64   `gorm:"sell_price3" json:"sell_price3"`
	SellPrice4    *float64   `gorm:"sell_price4" json:"sell_price4"`
	SellPrice5    *float64   `gorm:"sell_price5" json:"sell_price5"`
	Amount        *float64   `gorm:"amount" json:"amount"`
	AmountFinal   *float64   `gorm:"amount_final" json:"amount_final"`
	DiscValue     *float64   `gorm:"disc_value" json:"disc_value"`
	BatchNo       *string    `gorm:"batch_no" json:"batch_no"`
	ExpDate       *time.Time `gorm:"exp_date" json:"exp_date"`
	Vat           *float64   `gorm:"vat" json:"vat"`
	VatBg         *float64   `gorm:"vat_bg" json:"vat_bg"`
	VatLgSell     *float64   `gorm:"vat_lg_sell" json:"vat_lg_sell"`
	VatValue      *float64   `gorm:"vat_value" json:"vat_value"`
	VatBgValue    *float64   `gorm:"vat_bg_value" json:"vat_bg_value"`
	// VatLgValue     *float64   `gorm:"vat_lg_value" json:"vat_lg_value"`
	VatLgSellValue *float64 `gorm:"vat_lg_sell_value" json:"vat_lg_sell_value"`
	VatValueFinal  *float64 `gorm:"vat_value_final" json:"vat_value_final"`

	SellPriceSystem1 *float64 `gorm:"sell_price_system1" json:"sell_price_system1"`
	SellPriceSystem2 *float64 `gorm:"sell_price_system2" json:"sell_price_system2"`
	SellPriceSystem3 *float64 `gorm:"sell_price_system3" json:"sell_price_system3"`
	SellPriceFinal1  *float64 `gorm:"sell_price_final1" json:"sell_price_final1"`
	SellPriceFinal2  *float64 `gorm:"sell_price_final2" json:"sell_price_final2"`
	SellPriceFinal3  *float64 `gorm:"sell_price_final3" json:"sell_price_final3"`
	QtyPo1           *float64 `gorm:"qty_po1" json:"qty_po1"`
	QtyPo2           *float64 `gorm:"qty_po2" json:"qty_po2"`
	QtyPo3           *float64 `gorm:"qty_po3" json:"qty_po3"`
	SellPricePo1     *float64 `gorm:"sell_price_po1" json:"sell_price_po1"`
	SellPricePo2     *float64 `gorm:"sell_price_po2" json:"sell_price_po2"`
	SellPricePo3     *float64 `gorm:"sell_price_po3" json:"sell_price_po3"`
	DiscPo           *float64 `gorm:"disc_po" json:"disc_po"`
	DiscValuePo      *float64 `gorm:"disc_value_po" json:"disc_value_po"`
	DiscValueFinal   *float64 `gorm:"disc_value_final" json:"disc_value_final"`
	VatPo            *float64 `gorm:"vat_po" json:"vat_po"`
	VatValuePo       *float64 `gorm:"vat_value_po" json:"vat_value_po"`

	UnitId1        *string  `gorm:"unit_id1" json:"unit_id1"`
	UnitId2        *string  `gorm:"unit_id2" json:"unit_id2"`
	UnitId3        *string  `gorm:"unit_id3" json:"unit_id3"`
	UnitId4        *string  `gorm:"unit_id4" json:"unit_id4"`
	UnitId5        *string  `gorm:"unit_id5" json:"unit_id5"`
	ConvUnit2      *float32 `gorm:"conv_unit2" json:"conv_unit2"`
	ConvUnit3      *float32 `gorm:"conv_unit3" json:"conv_unit3"`
	ConvUnit4      *float32 `gorm:"conv_unit4" json:"conv_unit4"`
	ConvUnit5      *float32 `gorm:"conv_unit5" json:"conv_unit5"`
	Notes          *string  `gorm:"notes" json:"notes"`
	OriginalQtyPo1 *float64 `gorm:"original_qty_po1" json:"original_qty_po1"`
	OriginalQtyPo2 *float64 `gorm:"original_qty_po2" json:"original_qty_po2"`
	OriginalQtyPo3 *float64 `gorm:"original_qty_po3" json:"original_qty_po3"`
}

func (OrderDetail) TableName() string {
	return "sls.order_detail"
}

type OrderDetailRead struct {
	CustId        string     `gorm:"cust_id" json:"cust_id"`
	RoNo          string     `gorm:"ro_no" json:"ro_no"`
	SeqNo         int        `gorm:"seq_no" json:"seq_no"`
	OrderDetailID *int       `gorm:"column:order_detail_id;primaryKey" json:"order_detail_id"`
	ProId         int        `gorm:"pro_id" json:"pro_id"`
	ProCode       string     `gorm:"column:pro_code" json:"pro_code"`
	ProName       string     `gorm:"column:pro_name" json:"pro_name"`
	ItemType      int        `gorm:"item_type" json:"item_type"`
	Qty1          *float64   `gorm:"qty1" json:"qty1"`
	Qty2          *float64   `gorm:"qty2" json:"qty2"`
	Qty3          *float64   `gorm:"qty3" json:"qty3"`
	Qty4          *float64   `gorm:"qty4" json:"qty4"`
	Qty5          *float64   `gorm:"qty5" json:"qty5"`
	Qty1Final     *float64   `gorm:"qty1_final" json:"qty1_final"`
	Qty2Final     *float64   `gorm:"qty2_final" json:"qty2_final"`
	Qty3Final     *float64   `gorm:"qty3_final" json:"qty3_final"`
	Qty4Final     *float64   `gorm:"qty4_final" json:"qty4_final"`
	Qty5Final     *float64   `gorm:"qty5_final" json:"qty5_final"`
	Qty1Stok      *float64   `gorm:"qty1_stok" json:"qty1_stok"`
	Qty2Stok      *float64   `gorm:"qty2_stok" json:"qty2_stok"`
	Qty3Stok      *float64   `gorm:"qty3_stok" json:"qty3_stok"`
	PurchPrice1   *float64   `gorm:"purch_price1" json:"purch_price1"`
	PurchPrice2   *float64   `gorm:"purch_price2" json:"purch_price2"`
	PurchPrice3   *float64   `gorm:"purch_price3" json:"purch_price3"`
	PurchPrice4   *float64   `gorm:"purch_price4" json:"purch_price4"`
	PurchPrice5   *float64   `gorm:"purch_price5" json:"purch_price5"`
	SellPrice1    *float64   `gorm:"sell_price1" json:"sell_price1"`
	SellPrice2    *float64   `gorm:"sell_price2" json:"sell_price2"`
	SellPrice3    *float64   `gorm:"sell_price3" json:"sell_price3"`
	SellPrice4    *float64   `gorm:"sell_price4" json:"sell_price4"`
	SellPrice5    *float64   `gorm:"sell_price5" json:"sell_price5"`
	Amount        *float64   `gorm:"amount" json:"amount"`
	DiscValue     *float64   `gorm:"disc_value" json:"disc_value"`
	BatchNo       *string    `gorm:"batch_no" json:"batch_no"`
	ExpDate       *time.Time `gorm:"exp_date" json:"exp_date"`
	Vat           *float64   `gorm:"vat" json:"vat"`
	VatBg         *float64   `gorm:"vat_bg" json:"vat_bg"`
	VatLgSell     *float64   `gorm:"vat_lg_sell" json:"vat_lg_sell"`
	VatValue      *float64   `gorm:"vat_value" json:"vat_value"`
	VatBgValue    *float64   `gorm:"vat_bg_value" json:"vat_bg_value"`
	// VatLgValue     *float64   `gorm:"vat_lg_value" json:"vat_lg_value"`
	VatLgSellValue *float64 `gorm:"vat_lg_sell_value" json:"vat_lg_sell_value"`

	SellPriceSystem1 *float64 `gorm:"sell_price_system1" json:"sell_price_system1"`
	SellPriceSystem2 *float64 `gorm:"sell_price_system2" json:"sell_price_system2"`
	SellPriceSystem3 *float64 `gorm:"sell_price_system3" json:"sell_price_system3"`
	QtyPo1           *float64 `gorm:"qty_po1" json:"qty_po1"`
	QtyPo2           *float64 `gorm:"qty_po2" json:"qty_po2"`
	QtyPo3           *float64 `gorm:"qty_po3" json:"qty_po3"`
	SellPricePo1     *float64 `gorm:"sell_price_po1" json:"sell_price_po1"`
	SellPricePo2     *float64 `gorm:"sell_price_po2" json:"sell_price_po2"`
	SellPricePo3     *float64 `gorm:"sell_price_po3" json:"sell_price_po3"`
	DiscPo           *float64 `gorm:"disc_po" json:"disc_po"`
	DiscValuePo      *float64 `gorm:"disc_value_po" json:"disc_value_po"`
	VatPo            *float64 `gorm:"vat_po" json:"vat_po"`
	VatValuePo       *float64 `gorm:"vat_value_po" json:"vat_value_po"`

	UnitId1   *string `gorm:"unit_id1" json:"unit_id1"`
	UnitId2   *string `gorm:"unit_id2" json:"unit_id2"`
	UnitId3   *string `gorm:"unit_id3" json:"unit_id3"`
	UnitId4   *string `gorm:"unit_id4" json:"unit_id4"`
	UnitId5   *string `gorm:"unit_id5" json:"unit_id5"`
	ConvUnit2 *int    `gorm:"conv_unit2" json:"conv_unit2"`
	ConvUnit3 *int    `gorm:"conv_unit3" json:"conv_unit3"`
	ConvUnit4 *int    `gorm:"conv_unit4" json:"conv_unit4"`
	ConvUnit5 *int    `gorm:"conv_unit5" json:"conv_unit5"`
	Notes     *string `gorm:"notes" json:"notes"`
}

func (OrderDetailRead) TableName() string {
	return "sls.order_detail"
}
