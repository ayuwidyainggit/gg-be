package entity

type BpprDet struct {
	ID         *int     `json:"bppr_det_id"`
	ProID      int64    `json:"pro_id" validate:"required"`
	ProCode    string   `json:"pro_code"`
	ProName    string   `json:"pro_name"`
	Qty        *float64 `json:"qty" validate:"required"`
	QtyStr     *string  `json:"qty_str" validate:"required,qtystr"`
	UnitPrice1 *float64 `json:"unit_price1" validate:"required"`
	UnitPrice2 *float64 `json:"unit_price2"`
	UnitPrice3 *float64 `json:"unit_price3"`
	UnitPrice4 *float64 `json:"unit_price4"`
	UnitPrice5 *float64 `json:"unit_price5"`
	UnitId1    *string  `json:"unit_id1" validate:"required"`
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

type BpprDetUpdateRequest struct {
	ID         *int     `json:"bppr_det_id"`
	CustID     string   `json:"cust_id"`
	ProID      int64    `json:"pro_id" validate:"required"`
	Qty        float64  `json:"qty" validate:"required"`
	QtyStr     string   `json:"qty_str" validate:"required,qtystr"`
	UnitPrice1 *float64 `json:"unit_price1" validate:"required"`
	UnitPrice2 *float64 `json:"unit_price2"`
	UnitPrice3 *float64 `json:"unit_price3"`
	UnitPrice4 *float64 `json:"unit_price4"`
	UnitPrice5 *float64 `json:"unit_price5"`
	UnitId1    *string  `json:"unit_id1" validate:"required"`
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
