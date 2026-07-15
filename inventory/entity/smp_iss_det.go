package entity

type CreateSampleIssueDetBody struct {
	SmpIssDetId *int     `json:"smp_iss_det_id"`
	CustID      string   `json:"cust_id"`
	SmpIssNo    string   `json:"smp_iss_no"`
	SeqNo       int      `json:"seq_no"`
	ProID       int      `json:"pro_id" validate:"required"`
	Qty         *float64 `json:"qty" validate:"required"`
	QtyStr      *string  `json:"qty_str" validate:"required,qtystr"`
	UnitPrice1  *float64 `json:"unit_price1"`
	UnitPrice2  *float64 `json:"unit_price2"`
	UnitPrice3  *float64 `json:"unit_price3"`
	UnitPrice4  *float64 `json:"unit_price4"`
	UnitPrice5  *float64 `json:"unit_price5"`
	UnitId1     *string  `json:"unit_id1"`
	UnitId2     *string  `json:"unit_id2"`
	UnitId3     *string  `json:"unit_id3"`
	UnitId4     *string  `json:"unit_id4"`
	UnitId5     *string  `json:"unit_id5"`
	EmbInc      *float64 `json:"emb_inc"`
	EmbExc      *float64 `json:"emb_exc"`
	BatchNo     *string  `json:"batch_no"`
	ExpDate     *string  `json:"exp_date"`
	ConvUnit2   float64  `json:"conv_unit2"`
	ConvUnit3   float64  `json:"conv_unit3"`
	ConvUnit4   float64  `json:"conv_unit4"`
	ConvUnit5   float64  `json:"conv_unit5"`
}

type SampleIssueDetResp struct {
	SmpIssDetId *int     `json:"smp_iss_det_id"`
	SeqNo       int      `json:"seq_no"`
	ProID       int      `json:"pro_id" validate:"required"`
	ProCode     string   `json:"pro_code"`
	ProName     string   `json:"pro_name"`
	Qty         *float64 `json:"qty" validate:"required"`
	QtyStr      *string  `json:"qty_str" validate:"required,qtystr"`
	UnitPrice1  *float64 `json:"unit_price1"`
	UnitPrice2  *float64 `json:"unit_price2"`
	UnitPrice3  *float64 `json:"unit_price3"`
	UnitPrice4  *float64 `json:"unit_price4"`
	UnitPrice5  *float64 `json:"unit_price5"`
	UnitId1     *string  `json:"unit_id1"`
	UnitId2     *string  `json:"unit_id2"`
	UnitId3     *string  `json:"unit_id3"`
	UnitId4     *string  `json:"unit_id4"`
	UnitId5     *string  `json:"unit_id5"`
	EmbInc      *float64 `json:"emb_inc"`
	EmbExc      *float64 `json:"emb_exc"`
	BatchNo     *string  `json:"batch_no"`
	ExpDate     *string  `json:"exp_date"`
	ConvUnit2   float64  `json:"conv_unit2"`
	ConvUnit3   float64  `json:"conv_unit3"`
	ConvUnit4   float64  `json:"conv_unit4"`
	ConvUnit5   float64  `json:"conv_unit5"`
}

type UpdateSampleIssueDetBody struct {
	SmpIssDetId *int     `json:"smp_iss_det_id"`
	CustID      string   `json:"cust_id"`
	SmpIssNo    string   `json:"smp_iss_no"`
	SeqNo       int      `json:"seq_no"`
	ProID       *int     `json:"pro_id"`
	Qty         *float64 `json:"qty"`
	QtyStr      *string  `json:"qty_str"`
	UnitPrice1  *float64 `json:"unit_price1"`
	UnitPrice2  *float64 `json:"unit_price2"`
	UnitPrice3  *float64 `json:"unit_price3"`
	UnitPrice4  *float64 `json:"unit_price4"`
	UnitPrice5  *float64 `json:"unit_price5"`
	UnitId1     *string  `json:"unit_id1"`
	UnitId2     *string  `json:"unit_id2"`
	UnitId3     *string  `json:"unit_id3"`
	UnitId4     *string  `json:"unit_id4"`
	UnitId5     *string  `json:"unit_id5"`
	EmbInc      *float64 `json:"emb_inc"`
	EmbExc      *float64 `json:"emb_exc"`
	BatchNo     *string  `json:"batch_no"`
	ExpDate     *string  `json:"exp_date"`
	ConvUnit2   float64  `json:"conv_unit2"`
	ConvUnit3   float64  `json:"conv_unit3"`
	ConvUnit4   float64  `json:"conv_unit4"`
	ConvUnit5   float64  `json:"conv_unit5"`
}
