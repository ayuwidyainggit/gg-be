package model

type WhAdjDet struct {
	WhAdjDetId   int    `gorm:"column:wh_adj_det_id;primaryKey" json:"wh_adj_det_id"`
	CustID       string `gorm:"column:cust_id" json:"cust_id"`
	AdjNo        string `gorm:"column:adj_no" json:"stock_adjustment_no"`
	SeqNo        int    `gorm:"column:seq_no" json:"seq_no"`
	ProID        int    `gorm:"column:pro_id" json:"pro_id"`
	Qty          int    `gorm:"column:qty" json:"qty"`
	WhAdjDetType int    `gorm:"column:wh_adj_det_type" json:"wh_adj_det_type"`
}

func (WhAdjDet) TableName() string {
	return "inv.wh_adj_det"
}

type WhAdjDetRead struct {
	WhAdjDetId   *int    `gorm:"column:wh_adj_det_id;primaryKey" json:"wh_adj_det_id"`
	CustID       string  `gorm:"column:cust_id" json:"cust_id"`
	AdjNo        string  `gorm:"column:adj_no" json:"adj_no"`
	SeqNo        int     `gorm:"column:seq_no" json:"seq_no"`
	ProID        int     `gorm:"column:pro_id" json:"pro_id"`
	ProCode      string  `gorm:"column:pro_code" json:"pro_code"`
	ProName      string  `gorm:"column:pro_name" json:"pro_name"`
	Qty          float64 `gorm:"column:qty" json:"qty"`
	UnitId1      *string `gorm:"column:unit_id1" json:"unit_id1"`
	UnitId2      *string `gorm:"column:unit_id2" json:"unit_id2"`
	UnitId3      *string `gorm:"column:unit_id3" json:"unit_id3"`
	UnitId4      *string `gorm:"column:unit_id4" json:"unit_id4"`
	UnitId5      *string `gorm:"column:unit_id5" json:"unit_id5"`
	WhAdjDetType int     `gorm:"column:wh_adj_det_type" json:"wh_adj_det_type"`
	ConvUnit2    float64 `gorm:"column:conv_unit2" json:"conv_unit2"`
	ConvUnit3    float64 `gorm:"column:conv_unit3" json:"conv_unit3"`
}

func (WhAdjDetRead) TableName() string {
	return "inv.wh_adj_det"
}
