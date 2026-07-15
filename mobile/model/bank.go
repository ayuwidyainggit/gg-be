package model

type OutletBank struct {
	BankName     string `gorm:"column:bank_name;type:varchar(255);not null" json:"bank_name"`
	BankCode     string `gorm:"column:bank_code;type:varchar(255);not null" json:"bank_code"`
	AccountName  string `gorm:"column:account_name;type:varchar(255);not null" json:"account_name"`
	AccountNo    string `gorm:"column:account_no;type:varchar(255);not null" json:"account_no"`
	BankID       int64  `gorm:"column:bank_id;type:int64;not null" json:"bank_id"`
	OutletID     int64  `gorm:"column:outlet_id;type:int64;not null" json:"outlet_id"`
	OutletBankID int64  `gorm:"column:outlet_bank_id;type:int64;not null" json:"outlet_bank_id"`
}

func (OutletBank) TableName() string {
	return "mst.m_outlet_bank"
}
