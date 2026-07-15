package entity

type CreateVanSoDetBody struct {
	VanSoNo   string   `json:"van_so_no"`
	ProID     int      `json:"pro_id"`
	QtyPhs    *float64 `json:"qty_phs"`
	QtyPhsStr *string  `json:"qty_phs_str"`
	QtyBs     *float64 `json:"qty_bs"`
	QtyBsStr  *string  `json:"qty_bs_str"`
	QtyExp    *float64 `json:"qty_exp"`
	QtyExpStr *string  `json:"qty_exp_str"`
	QtyVan    *float64 `json:"qty_van"`
	QtyVanStr *string  `json:"qty_van_str"`
	QtyBlc    *float64 `json:"qty_blc"`
	QtyBlcStr *string  `json:"qty_blc_str"`
	UnitId1   *string  `json:"unit_id1"`
	UnitId2   *string  `json:"unit_id2"`
	UnitId3   *string  `json:"unit_id3"`
	UnitId4   *string  `json:"unit_id4"`
	UnitId5   *string  `json:"unit_id5"`
	ConvUnit2 float64  `json:"conv_unit2"`
	ConvUnit3 float64  `json:"conv_unit3"`
	ConvUnit4 float64  `json:"conv_unit4"`
	ConvUnit5 float64  `json:"conv_unit5"`
}

type VanSoDetResponse struct {
	VanSoDetID *int64   `json:"van_so_det_id"`
	VanSoNo    string   `json:"van_so_no"`
	ProID      int      `json:"pro_id"`
	ProCode    string   `json:"pro_code"`
	ProName    string   `json:"pro_name"`
	QtyPhs     *float64 `json:"qty_phs"`
	QtyPhsStr  *string  `json:"qty_phs_str"`
	QtyBs      *float64 `json:"qty_bs"`
	QtyBsStr   *string  `json:"qty_bs_str"`
	QtyExp     *float64 `json:"qty_exp"`
	QtyExpStr  *string  `json:"qty_exp_str"`
	QtyVan     *float64 `json:"qty_van"`
	QtyVanStr  *string  `json:"qty_van_str"`
	QtyBlc     *float64 `json:"qty_blc"`
	QtyBlcStr  *string  `json:"qty_blc_str"`
	UnitId1    *string  `json:"unit_id1"`
	UnitId2    *string  `json:"unit_id2"`
	UnitId3    *string  `json:"unit_id3"`
	UnitId4    *string  `json:"unit_id4"`
	UnitId5    *string  `json:"unit_id5"`
	ConvUnit2  float64  `json:"conv_unit2"`
	ConvUnit3  float64  `json:"conv_unit3"`
	ConvUnit4  float64  `json:"conv_unit4"`
	ConvUnit5  float64  `json:"conv_unit5"`
}
type UpdateVanSoDetBody struct {
	VanSoDetID *int64   `json:"van_so_det_id"`
	ProID      *int     `json:"pro_id"`
	QtyPhs     *float64 `json:"qty_phs"`
	QtyPhsStr  *string  `json:"qty_phs_str"`
	QtyBs      *float64 `json:"qty_bs"`
	QtyBsStr   *string  `json:"qty_bs_str"`
	QtyExp     *float64 `json:"qty_exp"`
	QtyExpStr  *string  `json:"qty_exp_str"`
	QtyVan     *float64 `json:"qty_van"`
	QtyVanStr  *string  `json:"qty_van_str"`
	QtyBlc     *float64 `json:"qty_blc"`
	QtyBlcStr  *string  `json:"qty_blc_str"`
	UnitId1    *string  `json:"unit_id1"`
	UnitId2    *string  `json:"unit_id2"`
	UnitId3    *string  `json:"unit_id3"`
	UnitId4    *string  `json:"unit_id4"`
	UnitId5    *string  `json:"unit_id5"`
	ConvUnit2  float64  `json:"conv_unit2"`
	ConvUnit3  float64  `json:"conv_unit3"`
	ConvUnit4  float64  `json:"conv_unit4"`
	ConvUnit5  float64  `json:"conv_unit5"`
}
