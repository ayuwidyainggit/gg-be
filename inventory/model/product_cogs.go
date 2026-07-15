package model

import (
	"time"

	"gorm.io/gorm"
)

type ProductCogs struct {
	CustID    string     `gorm:"column:cust_id" json:"cust_id"`
	ProID     *int64     `gorm:"column:pro_id" json:"pro_id"`
	CogsDate  *time.Time `gorm:"column:cogs_date" json:"cogs_date"`
	Cogs      *float64   `gorm:"column:cogs" json:"cogs"`
	TrCode    string     `gorm:"column:tr_code" json:"tr_code"`
	TrNo      string     `gorm:"column:tr_no" json:"tr_no"`
}

func (ProductCogs) TableName() string {
	return "mst.m_product_cogs"
}

func (m *ProductCogs) BeforeCreate(trx *gorm.DB) (err error) {

	return nil
}

func (m *ProductCogs) BeforeUpdate(trx *gorm.DB) (err error) {
	return nil
}
