package model

import (
	"time"

	"gorm.io/gorm"
)

type MCoa struct {
	CustID    string         `gorm:"column:cust_id" json:"cust_id"`
	CoaID     *int64         `gorm:"column:coa_id;primaryKey" json:"coa_id"`
	CoaCode   *string        `gorm:"column:coa_code" json:"coa_code"`
	CoaName   *string        `gorm:"column:coa_name" json:"coa_name"`
	Level     *int64         `gorm:"column:level" json:"level"`
	CoaTypeID *int64         `gorm:"column:coa_type_id" json:"coa_type_id"`
	ParentID  *int64         `gorm:"column:parent_id" json:"parent_id"`
	CashType  *int64         `gorm:"column:cash_type" json:"cash_type"`
	DefBlc    *string        `gorm:"column:def_blc" json:"def_blc"`
	IsDetail  *bool          `gorm:"column:is_detail" json:"is_detail"`
	CreatedBy *int64         `gorm:"column:created_by" json:"created_by"`
	CreatedAt time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedBy *int64         `gorm:"column:updated_by" json:"updated_by"`
	UpdatedAt time.Time      `gorm:"column:updated_at" json:"updated_at"`
	IsDel     bool           `gorm:"column:is_del" json:"is_del"`
	DeletedBy *int64         `gorm:"column:deleted_by" json:"deleted_by"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at" json:"deleted_at"`
}

func (m *MCoa) BeforeCreate(trx *gorm.DB) (err error) {
	now := time.Now()

	// intTmpsStr := now.UnixNano() / int64(time.Millisecond)
	// m.ChqTrNo = strconv.Itoa(int(intTmpsStr))
	m.CreatedAt = now
	m.UpdatedAt = now
	m.UpdatedBy = m.CreatedBy
	return nil
}

func (MCoa) TableName() string {
	return "acf.m_coa"
}

type MCoaList struct {
	CustID        string         `gorm:"column:cust_id" json:"cust_id"`
	CoaID         int64          `gorm:"column:coa_id;primaryKey" json:"coa_id"`
	CoaCode       *string        `gorm:"column:coa_code" json:"coa_code"`
	CoaName       *string        `gorm:"column:coa_name" json:"coa_name"`
	Level         *int64         `gorm:"column:level" json:"level"`
	CoaTypeID     *int64         `gorm:"column:coa_type_id" json:"coa_type_id"`
	CoaTypeName   *string        `gorm:"column:coa_type_name" json:"coa_type_name"`
	ParentID      *int64         `gorm:"column:parent_id" json:"parent_id"`
	CashType      *int64         `gorm:"column:cash_type" json:"cash_type"`
	DefBlc        *string        `gorm:"column:def_blc" json:"def_blc"`
	IsDetail      bool           `gorm:"column:is_detail" json:"is_detail"`
	CreatedBy     *int64         `gorm:"column:created_by" json:"created_by"`
	CreatedAt     time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedBy     *int64         `gorm:"column:updated_by" json:"updated_by"`
	UpdatedAt     time.Time      `gorm:"column:updated_at" json:"updated_at"`
	UpdatedByName *string        `gorm:"column:updated_by_name" json:"updated_by_name"`
	IsDel         bool           `gorm:"column:is_del" json:"is_del"`
	DeletedBy     *int64         `gorm:"column:deleted_by" json:"deleted_by"`
	DeletedAt     gorm.DeletedAt `gorm:"column:deleted_at" json:"deleted_at"`
}

func (MCoaList) TableName() string {
	return "acf.m_coa"
}
