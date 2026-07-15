package model

type OpexTrDet struct {
	CustID      string   `gorm:"column:cust_id" json:"cust_id"`
	OpexTrNo    string   `gorm:"column:opex_tr_no" json:"opex_tr_no"`
	OpexTrDetID *int64   `gorm:"column:opex_tr_det_id;primaryKey" json:"opex_tr_det_id"`
	CcID        *int64   `gorm:"column:cc_id" json:"cc_id"`
	OpexID      *int64   `gorm:"column:opex_id" json:"opex_id"`
	TrDesc      *string  `gorm:"column:tr_desc" json:"tr_desc"`
	Amount      *float64 `gorm:"column:amount" json:"amount"`
}

func (OpexTrDet) TableName() string {
	return "acf.opex_tr_det"
}

type OpexTrDetRead struct {
	CustID      string   `gorm:"column:cust_id" json:"cust_id"`
	OpexTrNo    string   `gorm:"column:opex_tr_no" json:"opex_tr_no"`
	OpexTrDetID *int64   `gorm:"column:opex_tr_det_id;primaryKey" json:"opex_tr_det_id"`
	CcID        *int64   `gorm:"column:cc_id" json:"cc_id"`
	OpexID      *int64   `gorm:"column:opex_id" json:"opex_id"`
	OpexCode    string   `gorm:"column:opex_code" json:"opex_code"`
	OpexName    string   `gorm:"column:opex_name" json:"opex_name"`
	TrDesc      *string  `gorm:"column:tr_desc" json:"tr_desc"`
	Amount      *float64 `gorm:"column:amount" json:"amount"`
}

func (OpexTrDetRead) TableName() string {
	return "acf.opex_tr_det"
}
