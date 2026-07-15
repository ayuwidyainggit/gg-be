package model

type VanSoDet struct {
	VanSoDetID *int64   `gorm:"column:van_so_det_id;primaryKey" json:"van_so_det_id"`
	CustID     string   `gorm:"column:cust_id" json:"cust_id"`
	VanSoNo    string   `gorm:"column:van_so_no" json:"van_so_no"`
	ProID      int      `gorm:"column:pro_id" json:"pro_id"`
	QtyPhs     *float64 `gorm:"column:qty_phs" json:"qty_phs"`
	QtyPhsStr  *string  `gorm:"column:qty_phs_str" json:"qty_phs_str"`
	QtyBs      *float64 `gorm:"column:qty_bs" json:"qty_bs"`
	QtyBsStr   *string  `gorm:"column:qty_bs_str" json:"qty_bs_str"`
	QtyExp     *float64 `gorm:"column:qty_exp" json:"qty_exp"`
	QtyExpStr  *string  `gorm:"column:qty_exp_str" json:"qty_exp_str"`
	QtyVan     *float64 `gorm:"column:qty_van" json:"qty_van"`
	QtyVanStr  *string  `gorm:"column:qty_van_str" json:"qty_van_str"`
	QtyBlc     *float64 `gorm:"column:qty_blc" json:"qty_blc"`
	QtyBlcStr  *string  `gorm:"column:qty_blc_str" json:"qty_blc_str"`
	UnitId1    *string  `gorm:"column:unit_id1" json:"unit_id1"`
	UnitId2    *string  `gorm:"column:unit_id2" json:"unit_id2"`
	UnitId3    *string  `gorm:"column:unit_id3" json:"unit_id3"`
	UnitId4    *string  `gorm:"column:unit_id4" json:"unit_id4"`
	UnitId5    *string  `gorm:"column:unit_id5" json:"unit_id5"`
	ConvUnit2  *float64 `gorm:"column:conv_unit2" json:"conv_unit2"`
	ConvUnit3  *float64 `gorm:"column:conv_unit3" json:"conv_unit3"`
	ConvUnit4  *float64 `gorm:"column:conv_unit4" json:"conv_unit4"`
	ConvUnit5  *float64 `gorm:"column:conv_unit5" json:"conv_unit5"`
}

func (VanSoDet) TableName() string {
	return "inv.van_so_det"
}

type VanSoDetRead struct {
	VanSoDetID *int64   `gorm:"column:van_so_det_id;primaryKey" json:"van_so_det_id"`
	CustID     string   `gorm:"column:cust_id" json:"cust_id"`
	VanSoNo    string   `gorm:"column:van_so_no" json:"van_so_no"`
	ProID      int      `gorm:"column:pro_id" json:"pro_id"`
	ProCode    string   `gorm:"column:pro_code" json:"pro_code"`
	ProName    string   `gorm:"column:pro_name" json:"pro_name"`
	QtyPhs     *float64 `gorm:"column:qty_phs" json:"qty_phs"`
	QtyPhsStr  *string  `gorm:"column:qty_phs_str" json:"qty_phs_str"`
	QtyBs      *float64 `gorm:"column:qty_bs" json:"qty_bs"`
	QtyBsStr   *string  `gorm:"column:qty_bs_str" json:"qty_bs_str"`
	QtyExp     *float64 `gorm:"column:qty_exp" json:"qty_exp"`
	QtyExpStr  *string  `gorm:"column:qty_exp_str" json:"qty_exp_str"`
	QtyVan     *float64 `gorm:"column:qty_van" json:"qty_van"`
	QtyVanStr  *string  `gorm:"column:qty_van_str" json:"qty_van_str"`
	QtyBlc     *float64 `gorm:"column:qty_blc" json:"qty_blc"`
	QtyBlcStr  *string  `gorm:"column:qty_blc_str" json:"qty_blc_str"`
	UnitId1    *string  `gorm:"column:unit_id1" json:"unit_id1"`
	UnitId2    *string  `gorm:"column:unit_id2" json:"unit_id2"`
	UnitId3    *string  `gorm:"column:unit_id3" json:"unit_id3"`
	UnitId4    *string  `gorm:"column:unit_id4" json:"unit_id4"`
	UnitId5    *string  `gorm:"column:unit_id5" json:"unit_id5"`
	ConvUnit2  *float64 `gorm:"column:conv_unit2" json:"conv_unit2"`
	ConvUnit3  *float64 `gorm:"column:conv_unit3" json:"conv_unit3"`
	ConvUnit4  *float64 `gorm:"column:conv_unit4" json:"conv_unit4"`
	ConvUnit5  *float64 `gorm:"column:conv_unit5" json:"conv_unit5"`
}

func (VanSoDetRead) TableName() string {
	return "inv.van_so_det"
}
