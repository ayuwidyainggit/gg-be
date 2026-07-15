package model

import (
	"time"

	"gorm.io/gorm"
)

type Tls struct {
	CustId     string          `gorm:"cust_id" json:"cust_id"`
	TlsId      int64           `gorm:"column:tls_id;primaryKey" json:"tls_id"`
	TlsDate    *time.Time      `gorm:"tls_date" json:"tls_date"`
	SalesmanID *int64          `gorm:"salesman_id" json:"salesman_id"`
	Notes      *string         `gorm:"notes" json:"notes"`
	DataStatus *int64          `gorm:"data_status" json:"data_status"`
	CreatedBy  *int64          `gorm:"created_by" json:"created_by"`
	CreatedAt  time.Time       `gorm:"created_at" json:"created_at"`
	UpdatedBy  *int64          `gorm:"updated_by" json:"updated_by"`
	UpdatedAt  time.Time       `gorm:"updated_at" json:"updated_at"`
	IsDel      bool            `gorm:"is_del" json:"is_del"`
	DeletedBy  *int64          `gorm:"deleted_by" json:"deleted_by"`
	DeletedAt  *gorm.DeletedAt `gorm:"deleted_at" json:"deleted_at"`
}

func (Tls) TableName() string {
	return "sls.tls"
}

func (m *Tls) BeforeCreate(trx *gorm.DB) (err error) {
	m.CreatedAt = time.Now()
	m.UpdatedAt = time.Now()
	m.UpdatedBy = m.CreatedBy
	return nil
}

type TlsList struct {
	CustId        string          `gorm:"cust_id" json:"cust_id"`
	TlsId         int64           `gorm:"column:tls_id;primaryKey" json:"tls_id"`
	TlsDate       *time.Time      `gorm:"tls_date" json:"tls_date"`
	SalesmanID    *int64          `gorm:"salesman_id" json:"salesman_id"`
	SalesmanCode  *string         `gorm:"column:salesman_code" json:"salesman_code"`
	SalesmanName  *string         `gorm:"column:salesman_name" json:"salesman_name"`
	Notes         *string         `gorm:"notes" json:"notes"`
	DataStatus    *int64          `gorm:"data_status" json:"data_status"`
	CreatedBy     *int64          `gorm:"created_by" json:"created_by"`
	CreatedAt     time.Time       `gorm:"created_at" json:"created_at"`
	UpdatedBy     *int64          `gorm:"updated_by" json:"updated_by"`
	UpdatedByName *string         `gorm:"column:updated_by_name" json:"updated_by_name"`
	UpdatedAt     time.Time       `gorm:"updated_at" json:"updated_at"`
	IsDel         bool            `gorm:"is_del" json:"is_del"`
	DeletedBy     *int64          `gorm:"deleted_by" json:"deleted_by"`
	DeletedAt     *gorm.DeletedAt `gorm:"deleted_at" json:"deleted_at"`
}

func (TlsList) TableName() string {
	return "sls.tls"
}
