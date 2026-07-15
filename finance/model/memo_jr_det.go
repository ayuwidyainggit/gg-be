package model

type MemoJrDet struct {
	MemoJrDetID *int64   `gorm:"column:memo_jr_det_id;primaryKey" json:"memo_jr_det_id"`
	CustID      string   `gorm:"column:cust_id" json:"cust_id"`
	MjNo        string   `gorm:"column:mj_no" json:"mj_no"`
	CcID        int      `gorm:"column:cc_id" json:"cc_id"`
	CoaID       int      `gorm:"column:coa_id" json:"coa_id"`
	Debit       *float64 `gorm:"column:debit" json:"debit"`
	Credit      *float64 `gorm:"column:credit" json:"credit"`
	Notes       *string  `gorm:"column:notes" json:"notes"`
}

func (MemoJrDet) TableName() string {
	return "acf.memo_jr_det"
}

type MemoJrDetRead struct {
	MemoJrDetID *int64   `gorm:"column:memo_jr_det_id;primaryKey" json:"memo_jr_det_id"`
	CustID      string   `gorm:"column:cust_id" json:"cust_id"`
	MjNo        string   `gorm:"column:mj_no" json:"mj_no"`
	CcID        int      `gorm:"column:cc_id" json:"cc_id"`
	CoaID       int      `gorm:"column:coa_id" json:"coa_id"`
	CoaCode     string   `gorm:"column:coa_code" json:"coa_code"`
	CoaName     string   `gorm:"column:coa_name" json:"coa_name"`
	Debit       *float64 `gorm:"column:debit" json:"debit"`
	Credit      *float64 `gorm:"column:credit" json:"credit"`
	Notes       *string  `gorm:"column:notes" json:"notes"`
}

func (MemoJrDetRead) TableName() string {
	return "acf.memo_jr_det"
}
