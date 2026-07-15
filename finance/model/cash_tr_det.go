package model

type CashTrDet struct {
	CashTrDetId *int64   `gorm:"column:cash_tr_det_id;primaryKey" json:"cash_tr_det_id"`
	CustID      string   `gorm:"column:cust_id" json:"cust_id"`
	CashTrNo    string   `gorm:"column:cash_tr_no" json:"cash_tr_no"`
	CoaId       *int64   `gorm:"column:coa_id" json:"coa_id"`
	Amount      *float64 `gorm:"column:amount" json:"amount"`
}

func (CashTrDet) TableName() string {
	return "acf.cash_tr_det"
}

type CashTrDetRead struct {
	CustID      string   `gorm:"column:cust_id" json:"cust_id"`
	CashTrNo    string   `gorm:"column:cash_tr_no" json:"cash_tr_no"`
	CashTrDetId *int64   `gorm:"column:cash_tr_det_id;primaryKey" json:"cash_tr_det_id"`
	CoaId       *int64   `gorm:"column:coa_id" json:"coa_id"`
	CoaCode     string   `gorm:"column:coa_code" json:"coa_code"`
	CoaName     string   `gorm:"column:coa_name" json:"coa_name"`
	Amount      *float64 `gorm:"column:amount" json:"amount"`
}

func (CashTrDetRead) TableName() string {
	return "acf.cash_tr_det"
}
