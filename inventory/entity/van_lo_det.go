package entity

type VanLoDetCreateBody struct {
	SeqNo      int      `json:"seq_no"`
	ProID      int      `json:"pro_id"`
	ItemType   int      `json:"item_type"`
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

type VanLoDetUpdateBody struct {
	VanLoDetID *int64   `json:"van_lo_det_id"`
	SeqNo      int      `json:"seq_no"`
	ProID      int      `json:"pro_id"`
	ItemType   int      `json:"item_type"`
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

type VanLoDetReadResponse struct {
	VanLoDetID *int64   `json:"van_lo_det_id"`
	SeqNo      int      `json:"seq_no"`
	ProID      int      `json:"pro_id"`
	ItemType   int      `json:"item_type"`
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

type VanLoDetCreateGroup struct {
	Normal []VanLoDetCreateBody `json:"normal"`
	Promo  []VanLoDetCreateBody `json:"promo"`
}

type VanLoDetReadGroup struct {
	Normal []VanLoDetReadResponse `json:"normal"`
	Promo  []VanLoDetReadResponse `json:"promo"`
}

type VanLoDetUpdateGroup struct {
	Normal []VanLoDetUpdateBody `json:"normal"`
	Promo  []VanLoDetUpdateBody `json:"promo"`
}
