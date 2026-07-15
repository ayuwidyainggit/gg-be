package model

import (
	"time"

	"gorm.io/gorm"
)

type MApDisc struct {
	CustID     string         `gorm:"column:cust_id" json:"cust_id"`
	ApDiscID   int            `gorm:"column:ap_disc_id;primaryKey" json:"ap_disc_id"`
	ProID      *int64         `gorm:"column:pro_id" json:"pro_id"`
	PurchPrice *float64       `gorm:"column:purch_price" json:"purch_price"`
	NettPrice  *float64       `gorm:"column:nett_price" json:"nett_price"`
	DiscP      *float64       `gorm:"column:disc_p" json:"disc_p"`
	CreatedBy  *int64         `gorm:"column:created_by" json:"created_by"`
	CreatedAt  time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedBy  *int64         `gorm:"column:updated_by" json:"updated_by"`
	UpdatedAt  time.Time      `gorm:"column:updated_at" json:"updated_at"`
	IsDel      bool           `gorm:"column:is_del" json:"is_del"`
	DeletedBy  *int64         `gorm:"column:deleted_by" json:"deleted_by"`
	DeletedAt  gorm.DeletedAt `gorm:"column:deleted_at" json:"deleted_at"`
}

func (m *MApDisc) BeforeCreate(trx *gorm.DB) (err error) {
	m.CreatedAt = time.Now()
	m.UpdatedAt = time.Now()
	m.UpdatedBy = m.CreatedBy
	return nil
}
func (MApDisc) TableName() string {
	return "acf.m_ap_disc"
}

type MApDiscList struct {
	CustID        string         `gorm:"column:cust_id" json:"cust_id"`
	ApDiscID      int            `gorm:"column:ap_disc_id;primaryKey" json:"ap_disc_id"`
	ProID         *int64         `gorm:"column:pro_id" json:"pro_id"`
	ProCode       *string        `gorm:"column:pro_code" json:"pro_code"`
	ProName       *string        `gorm:"column:pro_name" json:"pro_name"`
	PurchPrice    *float64       `gorm:"column:purch_price" json:"purch_price"`
	NettPrice     *float64       `gorm:"column:nett_price" json:"nett_price"`
	DiscP         *float64       `gorm:"column:disc_p" json:"disc_p"`
	CreatedBy     *int64         `gorm:"column:created_by" json:"created_by"`
	CreatedAt     time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedBy     *int64         `gorm:"column:updated_by" json:"updated_by"`
	UpdatedAt     time.Time      `gorm:"column:updated_at" json:"updated_at"`
	UpdatedByName *string        `gorm:"column:updated_by_name" json:"updated_by_name"`
	IsDel         bool           `gorm:"column:is_del" json:"is_del"`
	DeletedBy     *int64         `gorm:"column:deleted_by" json:"deleted_by"`
	DeletedAt     gorm.DeletedAt `gorm:"column:deleted_at" json:"deleted_at"`
}

func (MApDiscList) TableName() string {
	return "acf.m_ap_disc"
}
