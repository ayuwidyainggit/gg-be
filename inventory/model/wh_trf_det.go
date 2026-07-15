package model

type WhTrfDet struct {
	WhTrfDetId *int   `gorm:"column:wh_trf_det_id;primaryKey" json:"wh_trf_det_id"`
	CustID     string `gorm:"column:cust_id" json:"cust_id"`
	WhTrfNo    string `gorm:"column:wh_trf_no" json:"wh_trf_no"`
	SeqNo      int    `gorm:"column:seq_no" json:"seq_no"`
	ProID      int    `gorm:"column:pro_id" json:"pro_id"`
	Qty        int    `gorm:"column:qty" json:"qty"`
}

func (WhTrfDet) TableName() string {
	return "inv.wh_trf_det"
}

type WhTrfDetRead struct {
	WhTrfDetId  *int    `gorm:"column:wh_trf_det_id;primaryKey" json:"wh_trf_det_id"`
	CustID      string  `gorm:"column:cust_id" json:"cust_id"`
	WhTrfNo     string  `gorm:"column:wh_trf_no" json:"stock_trf_no"`
	SeqNo       int     `gorm:"column:seq_no" json:"seq_no"`
	ProID       int64   `gorm:"column:pro_id" json:"pro_id"`
	ProCode     string  `gorm:"column:pro_code" json:"pro_code"`
	ProName     string  `gorm:"column:pro_name" json:"pro_name"`
	Qty         float64 `gorm:"column:qty" json:"qty"`
	SellPrice1  float64 `gorm:"column:sell_price1" json:"sell_price1"`
	SellPrice2  float64 `gorm:"column:sell_price2" json:"sell_price2"`
	SellPrice3  float64 `gorm:"column:sell_price3" json:"sell_price3"`
	UnitId1     string  `gorm:"column:unit_id1" json:"unit_id1"`
	UnitId2     string  `gorm:"column:unit_id2" json:"unit_id2"`
	UnitId3     string  `gorm:"column:unit_id3" json:"unit_id3"`
	ConvUnit2   float64 `gorm:"column:conv_unit2" json:"conv_unit2"`
	ConvUnit3   float64 `gorm:"column:conv_unit3" json:"conv_unit3"`
	Vat         float64 `gorm:"column:vat" json:"vat"`
	VatBg       float64 `gorm:"column:vat_bg" json:"vat_bg"`
	VatLgPurch  float64 `gorm:"column:vat_lg_purch" json:"vat_lg_purch"`
	PurchPrice1 float64 `json:"purch_price1" gorm:"column:purch_price1"`
	PurchPrice2 float64 `json:"purch_price2" gorm:"column:purch_price2"`
	PurchPrice3 float64 `json:"purch_price3" gorm:"column:purch_price3"`
}

func (WhTrfDetRead) TableName() string {
	return "inv.wh_trf_det"
}
