package entity

type CreateOrderDetBody struct {
	SeqNo          int      `json:"seq_no"`
	ProId          int      `json:"pro_id"`
	ItemType       int      `json:"item_type"`
	Qty            *float64 `json:"qty"`
	Qty1           *float64 `json:"qty1"`
	Qty2           *float64 `json:"qty2"`
	Qty3           *float64 `json:"qty3"`
	Qty4           *float64 `json:"qty4"`
	Qty5           *float64 `json:"qty5"`
	Qty1Final      *float64 `json:"qty1_final"`
	Qty2Final      *float64 `json:"qty2_final"`
	Qty3Final      *float64 `json:"qty3_final"`
	Qty4Final      *float64 `json:"qty4_final"`
	Qty5Final      *float64 `json:"qty5_final"`
	Qty1Stok       *float64 `json:"qty1_stok"`
	Qty2Stok       *float64 `json:"qty2_stok"`
	Qty3Stok       *float64 `json:"qty3_stok"`
	PurchPrice1    *float64 `json:"purch_price1"`
	PurchPrice2    *float64 `json:"purch_price2"`
	PurchPrice3    *float64 `json:"purch_price3"`
	PurchPrice4    *float64 `json:"purch_price4"`
	PurchPrice5    *float64 `json:"purch_price5"`
	SellPrice1     *float64 `json:"sell_price1"`
	SellPrice2     *float64 `json:"sell_price2"`
	SellPrice3     *float64 `json:"sell_price3"`
	SellPrice4     *float64 `json:"sell_price4"`
	SellPrice5     *float64 `json:"sell_price5"`
	Amount         *float64 `json:"amount"`
	DiscValue      *float64 `json:"disc_value"`
	BatchNo        *string  `json:"batch_no"`
	ExpDate        *string  `json:"exp_date"`
	Vat            *float64 `gorm:"vat" json:"vat"`
	VatBg          *float64 `gorm:"vat_bg" json:"vat_bg"`
	VatLgSell      *float64 `gorm:"vat_lg_sell" json:"vat_lg_sell"`
	VatValue       *float64 `gorm:"vat_value" json:"vat_value"`
	VatBgValue     *float64 `gorm:"vat_bg_value" json:"vat_bg_value"`
	VatLgValue     *float64 `gorm:"vat_lg_value" json:"vat_lg_value"`
	VatLgSellValue *float64 `gorm:"vat_lg_sell_value" json:"vat_lg_sell_value"`
	UnitId1        *string  `gorm:"unit_id1" json:"unit_id1"`
	UnitId2        *string  `gorm:"unit_id2" json:"unit_id2"`
	UnitId3        *string  `gorm:"unit_id3" json:"unit_id3"`
	UnitId4        *string  `gorm:"unit_id4" json:"unit_id4"`
	UnitId5        *string  `gorm:"unit_id5" json:"unit_id5"`
	ConvUnit2      *int     `gorm:"conv_unit2" json:"conv_unit2"`
	ConvUnit3      *int     `gorm:"conv_unit3" json:"conv_unit3"`
	ConvUnit4      *int     `gorm:"conv_unit4" json:"conv_unit4"`
	ConvUnit5      *int     `gorm:"conv_unit5" json:"conv_unit5"`
	Notes          *string  `gorm:"notes" json:"notes"`
}

type OrderDetResponse struct {
	OrderDetId  int64    `json:"order_detail_id"`
	SeqNo       int      `json:"seq_no"`
	ProId       int      `json:"pro_id"`
	ProCode     string   `json:"pro_code"`
	ProName     string   `json:"pro_name"`
	OrderStatus string   `json:"order_status"`
	ItemType    int      `json:"item_type"`
	Qty         *float64 `json:"qty"`
	// QtyFinal       *float64 `json:"qty_final"`
	QtyPo          *float64 `json:"qty_po"`
	Qty1           *float64 `json:"qty1"`
	Qty2           *float64 `json:"qty2"`
	Qty3           *float64 `json:"qty3"`
	Qty4           *float64 `json:"qty4"`
	Qty5           *float64 `json:"qty5"`
	Qty1Final      *float64 `json:"qty1_final"`
	Qty2Final      *float64 `json:"qty2_final"`
	Qty3Final      *float64 `json:"qty3_final"`
	Qty4Final      *float64 `json:"qty4_final"`
	Qty5Final      *float64 `json:"qty5_final"`
	Qty1Stok       *float64 `json:"qty1_stok"`
	Qty2Stok       *float64 `json:"qty2_stok"`
	Qty3Stok       *float64 `json:"qty3_stok"`
	PurchPrice1    *float64 `json:"purch_price1"`
	PurchPrice2    *float64 `json:"purch_price2"`
	PurchPrice3    *float64 `json:"purch_price3"`
	PurchPrice4    *float64 `json:"purch_price4"`
	PurchPrice5    *float64 `json:"purch_price5"`
	SellPrice1     *float64 `json:"sell_price1"`
	SellPrice2     *float64 `json:"sell_price2"`
	SellPrice3     *float64 `json:"sell_price3"`
	SellPrice4     *float64 `json:"sell_price4"`
	SellPrice5     *float64 `json:"sell_price5"`
	Amount         *float64 `json:"amount"`
	DiscValue      *float64 `json:"disc_value"`
	BatchNo        *string  `json:"batch_no"`
	ExpDate        *string  `json:"exp_date"`
	Vat            *float64 `gorm:"vat" json:"vat"`
	VatBg          *float64 `gorm:"vat_bg" json:"vat_bg"`
	VatLgSell      *float64 `gorm:"vat_lg_sell" json:"vat_lg_sell"`
	VatValue       *float64 `gorm:"vat_value" json:"vat_value"`
	VatBgValue     *float64 `gorm:"vat_bg_value" json:"vat_bg_value"`
	VatLgValue     *float64 `gorm:"vat_lg_value" json:"vat_lg_value"`
	VatLgSellValue *float64 `gorm:"vat_lg_sell_value" json:"vat_lg_sell_value"`
	UnitId1        *string  `gorm:"unit_id1" json:"unit_id1"`
	UnitId2        *string  `gorm:"unit_id2" json:"unit_id2"`
	UnitId3        *string  `gorm:"unit_id3" json:"unit_id3"`
	UnitId4        *string  `gorm:"unit_id4" json:"unit_id4"`
	UnitId5        *string  `gorm:"unit_id5" json:"unit_id5"`
	ConvUnit2      *int     `gorm:"conv_unit2" json:"conv_unit2"`
	ConvUnit3      *int     `gorm:"conv_unit3" json:"conv_unit3"`
	MpConvUnit2    *int     `gorm:"mconv_unit2" json:"-"`
	MpConvUnit3    *int     `gorm:"mconv_unit3" json:"-"`
	ConvUnit4      *int     `gorm:"conv_unit4" json:"conv_unit4"`
	ConvUnit5      *int     `gorm:"conv_unit5" json:"conv_unit5"`
	Notes          *string  `gorm:"notes" json:"notes"`
}

type UpdateOrderDetBody struct {
	SeqNo          int      `json:"seq_no"`
	OrderDetId     *int64   `json:"order_detail_id"`
	ProId          int      `json:"pro_id"`
	ItemType       int      `json:"item_type"`
	Qty            *float64 `json:"qty"`
	Qty1           *float64 `json:"qty1"`
	Qty2           *float64 `json:"qty2"`
	Qty3           *float64 `json:"qty3"`
	Qty4           *float64 `json:"qty4"`
	Qty5           *float64 `json:"qty5"`
	PurchPrice1    *float64 `json:"purch_price1"`
	PurchPrice2    *float64 `json:"purch_price2"`
	PurchPrice3    *float64 `json:"purch_price3"`
	PurchPrice4    *float64 `json:"purch_price4"`
	PurchPrice5    *float64 `json:"purch_price5"`
	SellPrice1     *float64 `json:"sell_price1"`
	SellPrice2     *float64 `json:"sell_price2"`
	SellPrice3     *float64 `json:"sell_price3"`
	SellPrice4     *float64 `json:"sell_price4"`
	SellPrice5     *float64 `json:"sell_price5"`
	Amount         *float64 `json:"amount"`
	DiscValue      *float64 `json:"disc_value"`
	BatchNo        *string  `json:"batch_no"`
	ExpDate        *string  `json:"exp_date"`
	Vat            *float64 `gorm:"vat" json:"vat"`
	VatBg          *float64 `gorm:"vat_bg" json:"vat_bg"`
	VatLgSell      *float64 `gorm:"vat_lg_sell" json:"vat_lg_sell"`
	VatValue       *float64 `gorm:"vat_value" json:"vat_value"`
	VatBgValue     *float64 `gorm:"vat_bg_value" json:"vat_bg_value"`
	VatLgValue     *float64 `gorm:"vat_lg_value" json:"vat_lg_value"`
	VatLgSellValue *float64 `gorm:"vat_lg_sell_value" json:"vat_lg_sell_value"`
	UnitId1        *string  `gorm:"unit_id1" json:"unit_id1"`
	UnitId2        *string  `gorm:"unit_id2" json:"unit_id2"`
	UnitId3        *string  `gorm:"unit_id3" json:"unit_id3"`
	UnitId4        *string  `gorm:"unit_id4" json:"unit_id4"`
	UnitId5        *string  `gorm:"unit_id5" json:"unit_id5"`
	ConvUnit2      *int     `gorm:"conv_unit2" json:"conv_unit2"`
	ConvUnit3      *int     `gorm:"conv_unit3" json:"conv_unit3"`
	ConvUnit4      *int     `gorm:"conv_unit4" json:"conv_unit4"`
	ConvUnit5      *int     `gorm:"conv_unit5" json:"conv_unit5"`
	Notes          *string  `gorm:"notes" json:"notes"`
}

func (req CreateOrderDetBody) GetSafeQTY(qty int) int {
	getSafeQty := func(qty *float64) int {
		if qty == nil {
			return 0
		}
		return int(*qty)
	}

	switch qty {
	case 1:
		return getSafeQty(req.Qty1)
	case 2:
		return getSafeQty(req.Qty2)
	case 3:
		return getSafeQty(req.Qty3)
	default:
		return 0
	}
}

type OrderDetWithGroup struct {
	Normal []CreateOrderDetBody `json:"normal"`
	Promo  []CreateOrderDetBody `json:"promo"`
}
type OrderDetReadWithGroup struct {
	Normal []OrderDetResponse `json:"normal"`
	Promo  []OrderDetResponse `json:"promo"`
}
type UpdateOrderDetWithGroup struct {
	Normal []UpdateOrderDetBody `json:"normal"`
	Promo  []UpdateOrderDetBody `json:"promo"`
}
