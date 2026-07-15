package entity

type CreateWhTrfDetBody struct {
	SeqNo int     `json:"seq_no"`
	ProID int64   `json:"pro_id"`
	Qty1  float64 `json:"qty1"`
	Qty2  float64 `json:"qty2"`
	Qty3  float64 `json:"qty3"`
}

type WhTrfDetRespose struct {
	WhTrfDetId  *int    `json:"wh_trf_det_id"`
	WhTrfNo     string  `json:"stock_trf_no"`
	SeqNo       int     `json:"seq_no"`
	ProID       int64   `json:"pro_id"`
	ProCode     string  `json:"pro_code"`
	ProName     string  `json:"pro_name"`
	UnitId1     string  `json:"unit_id1"`
	UnitId2     string  `json:"unit_id2"`
	UnitId3     string  `json:"unit_id3"`
	SellPrice1  float64 `json:"sell_price1"`
	SellPrice2  float64 `json:"sell_price2"`
	SellPrice3  float64 `json:"sell_price3"`
	Vat         float64 `json:"vat"`
	VatLgPurch  float64 `json:"vat_lg_purch"`
	VatBg       float64 `json:"vat_bg"`
	Qty1        float64 `json:"qty1"`
	Qty2        float64 `json:"qty2"`
	Qty3        float64 `json:"qty3"`
	SubTotal    float64 `json:"sub_total"`
	Total       float64 `json:"total"`
	PurchPrice1 float64 `json:"purch_price1"`
	PurchPrice2 float64 `json:"purch_price2"`
	PurchPrice3 float64 `json:"purch_price3"`
}
type UpdateWhTrfDetBody struct {
	WhTrfDetId *int     `json:"wh_trf_det_id"`
	CustID     string   `json:"cust_id"`
	WhTrfNo    string   `json:"wh_trf_no"`
	SeqNo      *int     `json:"seq_no"`
	ProID      *int     `json:"pro_id"`
	Qty        *float64 `json:"qty"`
	QtyStr     *string  `json:"qty_str"`
	UnitPrice1 *float64 `json:"unit_price1"`
	UnitPrice2 *float64 `json:"unit_price2"`
	UnitPrice3 *float64 `json:"unit_price3"`
	UnitPrice4 *float64 `json:"unit_price4"`
	UnitPrice5 *float64 `json:"unit_price5"`
	UnitId1    *string  `json:"unit_id1"`
	UnitId2    *string  `json:"unit_id2"`
	UnitId3    *string  `json:"unit_id3"`
	UnitId4    *string  `json:"unit_id4"`
	UnitId5    *string  `json:"unit_id5"`
	EmbInc     *float64 `json:"emb_inc"`
	EmbExc     *float64 `json:"emb_exc"`
	BatchNo    *string  `json:"batch_no"`
	ExpDate    *string  `json:"exp_date"`
	ConvUnit2  float64  `json:"conv_unit2"`
	ConvUnit3  float64  `json:"conv_unit3"`
	ConvUnit4  float64  `json:"conv_unit4"`
	ConvUnit5  float64  `json:"conv_unit5"`
}
