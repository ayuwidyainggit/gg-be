package entity

type VanUlDetCreateGroup struct {
	Normal []VanUlDetCreateBody `json:"normal"`
	Promo  []VanUlDetCreateBody `json:"promo"`
}

type VanUlDetReadGroup struct {
	Normal []VanUlDetReadResponse `json:"normal"`
	Promo  []VanUlDetReadResponse `json:"promo"`
}

type VanUlDetUpdateGroup struct {
	Normal []VanUlDetUpdateBody `json:"normal"`
	Promo  []VanUlDetUpdateBody `json:"promo"`
}

type VanUlDetCreateBody struct {
	ProID      int      `json:"pro_id"`
	ItemType   int      `json:"item_type"`
	Qty        *float64 `json:"qty"`
	QtyStr     *string  `json:"qty_str"`
	QtyBs      *float64 `json:"qty_bs"`
	QtyBsStr   *string  `json:"qty_bs_str"`
	QtyExp     *float64 `json:"qty_exp"`
	QtyExpStr  *string  `json:"qty_exp_str"`
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
	SeqNo      int      `json:"seq_no"`
	ConvUnit2  float64  `json:"conv_unit2"`
	ConvUnit3  float64  `json:"conv_unit3"`
	ConvUnit4  float64  `json:"conv_unit4"`
	ConvUnit5  float64  `json:"conv_unit5"`
}

type VanUlDetUpdateBody struct {
	VanUlDetID *int64   `json:"van_ul_det_id"`
	CustID     string   `json:"cust_id"`
	VanUlNo    string   `json:"van_ul_no"`
	ProID      *int     `json:"pro_id"`
	ItemType   *int     `json:"item_type"`
	Qty        *float64 `json:"qty"`
	QtyStr     *string  `json:"qty_str"`
	QtyBs      *float64 `json:"qty_bs"`
	QtyBsStr   *string  `json:"qty_bs_str"`
	QtyExp     *float64 `json:"qty_exp"`
	QtyExpStr  *string  `json:"qty_exp_str"`
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
	SeqNo      *int     `json:"seq_no"`
	ConvUnit2  float64  `json:"conv_unit2"`
	ConvUnit3  float64  `json:"conv_unit3"`
	ConvUnit4  float64  `json:"conv_unit4"`
	ConvUnit5  float64  `json:"conv_unit5"`
}

type VanUlDetReadResponse struct {
	VanUlDetID int64    `json:"van_ul_det_id"`
	ProID      int      `json:"pro_id"`
	ProCode    string   `json:"pro_code"`
	ProName    string   `json:"pro_name"`
	ItemType   int      `json:"item_type"`
	Qty        *float64 `json:"qty"`
	QtyStr     *string  `json:"qty_str"`
	QtyBs      *float64 `json:"qty_bs"`
	QtyBsStr   *string  `json:"qty_bs_str"`
	QtyExp     *float64 `json:"qty_exp"`
	QtyExpStr  *string  `json:"qty_exp_str"`
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
	SeqNo      *int     `json:"seq_no"`
	ConvUnit2  float64  `json:"conv_unit2"`
	ConvUnit3  float64  `json:"conv_unit3"`
	ConvUnit4  float64  `json:"conv_unit4"`
	ConvUnit5  float64  `json:"conv_unit5"`
}
