package model

type AccountPayable struct {
	SupId     int64  `gorm:"column:sup_id" json:"sup_id"`
	InvoiceNo string `gorm:"column:invoice_no" json:"invoice_no"`
}

func (AccountPayable) TableName() string {
	return "acf.account_payable"
}
