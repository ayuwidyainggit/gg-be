package model

import (
	"time"

	"gorm.io/gorm"
)

type MOpex struct {
	CustId    string         `gorm:"cust_id" json:"cust_id"`
	OpexId    *int           `gorm:"column:opex_id;primaryKey" json:"opex_id"`
	OpexCode  string         `gorm:"opex_code" json:"opex_code"`
	OpexName  string         `gorm:"opex_name" json:"opex_name"`
	CoaId     *int64         `gorm:"coa_id" json:"coa_id"`
	IsActive  *bool          `gorm:"is_active" json:"is_active"`
	CreatedBy *int64         `gorm:"created_by" json:"created_by"`
	CreatedAt time.Time      `gorm:"created_at" json:"created_at"`
	UpdatedBy *int64         `gorm:"updated_by" json:"updated_by"`
	UpdatedAt time.Time      `gorm:"updated_at" json:"updated_at"`
	IsDel     bool           `gorm:"is_del" json:"is_del"`
	DeletedBy *int64         `gorm:"deleted_by" json:"deleted_by"`
	DeletedAt gorm.DeletedAt `grom:"deleted_at" json:"deleted_at"`
}

func (m *MOpex) BeforeCreate(trx *gorm.DB) (err error) {
	m.CreatedAt = time.Now()
	m.UpdatedAt = time.Now()
	m.UpdatedBy = m.CreatedBy
	return nil
}
func (MOpex) TableName() string {
	return "acf.m_opex"
}

type MOpexList struct {
	CustId        string         `gorm:"cust_id" json:"cust_id"`
	OpexId        *int           `gorm:"column:opex_id;primaryKey" json:"opex_id"`
	OpexCode      *string        `gorm:"opex_code" json:"opex_code"`
	OpexName      *string        `gorm:"opex_name" json:"opex_name"`
	CoaId         *int64         `gorm:"coa_id" json:"coa_id"`
	CoaCode       *string        `gorm:"coa_code" json:"coa_code"`
	CoaName       *string        `gorm:"coa_name" json:"coa_name"`
	IsActive      *bool          `gorm:"is_active" json:"is_active"`
	CreatedBy     *int64         `gorm:"created_by" json:"created_by"`
	CreatedAt     time.Time      `gorm:"created_at" json:"created_at"`
	UpdatedBy     *int64         `gorm:"updated_by" json:"updated_by"`
	UpdatedAt     time.Time      `gorm:"updated_at" json:"updated_at"`
	UpdatedByName *string        `gorm:"column:updated_by_name" json:"updated_by_name"`
	IsDel         bool           `gorm:"is_del" json:"is_del"`
	DeletedBy     *int64         `gorm:"deleted_by" json:"deleted_by"`
	DeletedAt     gorm.DeletedAt `grom:"deleted_at" json:"deleted_at"`
}

func (MOpexList) TableName() string {
	return "acf.m_opex"
}
