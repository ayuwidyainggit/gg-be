package entity

type InvoiceDetResponse struct {
	RoDetId         int64   `json:"order_detail_id"`
	SeqNo           int     `json:"seq_no"`
	ProId           int     `json:"pro_id"`
	ProCode         string  `json:"pro_code"`
	ProName         string  `json:"pro_name"`
	ItemType        int     `json:"item_type"`
	Qty1            float64 `json:"qty1"`
	Qty2            float64 `json:"qty2"`
	Qty3            float64 `json:"qty3"`
	Volume1         float64 `json:"volume1"`
	Volume2         float64 `json:"volume2"`
	Volume3         float64 `json:"volume3"`
	Weight1         float64 `json:"weight1"`
	Weight2         float64 `json:"weight2"`
	Weight3         float64 `json:"weight3"`
	Stock1          float64 `json:"stock1"`
	Stock2          float64 `json:"stock2"`
	Stock3          float64 `json:"stock3"`
	PurchPrice1     float64 `json:"purch_price1"`
	PurchPrice2     float64 `json:"purch_price2"`
	PurchPrice3     float64 `json:"purch_price3"`
	SellPrice1      float64 `json:"sell_price1"`
	SellPrice2      float64 `json:"sell_price2"`
	SellPrice3      float64 `json:"sell_price3"`
	ConvUnit1       float64 `json:"conv_unit1"`
	ConvUnit2       float64 `json:"conv_unit2"`
	ConvUnit3       float64 `json:"conv_unit3"`
	UnitId1         string  `json:"unit_id1"`
	UnitId2         string  `json:"unit_id2"`
	UnitId3         string  `json:"unit_id3"`
	Amount          float64 `json:"amount"`
	DiscValue       float64 `json:"disc_value"`
	BatchNo         string  `json:"batch_no"`
	ExpDate         string  `json:"exp_date"`
	PriceIncludePpn float64 `json:"price_include_ppn"`
	PriceExcludePpn float64 `json:"price_exclude_ppn"`
	NetValue        float64 `json:"net_value"`
	Vat             float64 `json:"vat"`
	VatValue        float64 `json:"vat_value"`
}

type UpdateInvoiceDetBody struct {
	SeqNo       int      `json:"seq_no"`
	RoDetId     *int64   `json:"ro_det_id"`
	ProId       int      `json:"pro_id"`
	ItemType    int      `json:"item_type"`
	Qty         *float64 `json:"qty"`
	QtyStr      *string  `json:"qty_str"`
	PurchPrice1 *float64 `json:"purch_price1"`
	PurchPrice2 *float64 `json:"purch_price2"`
	PurchPrice3 *float64 `json:"purch_price3"`
	SellPrice1  *float64 `json:"sell_price1"`
	SellPrice2  *float64 `json:"sell_price2"`
	SellPrice3  *float64 `json:"sell_price3"`
	Amount      *float64 `json:"amount"`
	DiscValue   *float64 `json:"disc_value"`
	BatchNo     *string  `json:"batch_no"`
	ExpDate     *string  `json:"exp_date"`
}

type InvoiceDetReadWithGroup struct {
	Normal []InvoiceDetResponse `json:"normal"`
	Promo  []InvoiceDetResponse `json:"promo"`
}
type UpdateInvoiceDetWithGroup struct {
	Normal []UpdateInvoiceDetBody `json:"normal"`
	Promo  []UpdateInvoiceDetBody `json:"promo"`
}
