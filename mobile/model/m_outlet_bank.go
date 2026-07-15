package model

import "time"

// MBank represents mst.m_bank table
type MBank struct {
	CustID    string     `gorm:"column:cust_id;type:varchar(10);primaryKey" json:"cust_id"`
	BankID    int64      `gorm:"column:bank_id;primaryKey;autoIncrement" json:"bank_id"`
	BankCode  *string    `gorm:"column:bank_code;type:varchar(10)" json:"bank_code"`
	BankName  *string    `gorm:"column:bank_name;type:varchar(150)" json:"bank_name"`
	IsActive  bool       `gorm:"column:is_active;default:true" json:"is_active"`
	CreatedBy *int64     `gorm:"column:created_by;type:int8" json:"created_by"`
	CreatedAt *time.Time `gorm:"column:created_at;type:timestamptz(6)" json:"created_at"`
	UpdatedBy *int64     `gorm:"column:updated_by;type:int8" json:"updated_by"`
	UpdatedAt *time.Time `gorm:"column:updated_at;type:timestamptz(6)" json:"updated_at"`
	IsDel     bool       `gorm:"column:is_del;default:false" json:"is_del"`
	DeletedBy *int64     `gorm:"column:deleted_by;type:int8" json:"deleted_by"`
	DeletedAt *time.Time `gorm:"column:deleted_at;type:timestamptz(6)" json:"deleted_at"`
}

func (MBank) TableName() string {
	return "mst.m_bank"
}

// MOutletBank represents mst.m_outlet_bank table
type MOutletBank struct {
	CustID       string  `gorm:"column:cust_id;type:varchar(10);primaryKey" json:"cust_id"`
	OutletID     int64   `gorm:"column:outlet_id;type:int4;primaryKey" json:"outlet_id"`
	BankID       int64   `gorm:"column:bank_id;type:int4" json:"bank_id"`
	AccountNo    *string `gorm:"column:account_no;type:varchar(50)" json:"account_no"`
	AccountName  *string `gorm:"column:account_name;type:varchar(150)" json:"account_name"`
	OutletBankID int64   `gorm:"column:outlet_bank_id;primaryKey;autoIncrement" json:"outlet_bank_id"`
}

func (MOutletBank) TableName() string {
	return "mst.m_outlet_bank"
}

// OutletBankInfo is a query result combining m_outlet_bank and m_bank data
type OutletBankInfo struct {
	OutletBankID int64   `gorm:"column:outlet_bank_id" json:"outlet_bank_id"`
	OutletID     int64   `gorm:"column:outlet_id" json:"outlet_id"`
	BankID       int64   `gorm:"column:bank_id" json:"bank_id"`
	BankCode     *string `gorm:"column:bank_code" json:"bank_code"`
	BankName     *string `gorm:"column:bank_name" json:"bank_name"`
	AccountNo    *string `gorm:"column:account_no" json:"account_no"`
	AccountName  *string `gorm:"column:account_name" json:"account_name"`
}
