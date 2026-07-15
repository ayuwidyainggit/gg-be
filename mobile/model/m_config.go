package model

import (
	"time"

	"gorm.io/gorm"
)

type MConfig struct {
	CustID      string    `gorm:"column:cust_id" json:"cust_id"`
	ConfigID    string    `gorm:"column:config_id;primaryKey" json:"config_id"`
	ConfigValue *string   `gorm:"column:config_value" json:"config_value"`
	DataType    *string   `gorm:"column:data_type" json:"data_type"`
	ConfigDesc  *string   `gorm:"column:config_desc" json:"config_desc"`
	Module      *string   `gorm:"column:module" json:"module"`
	CreatedDate time.Time `gorm:"column:created_date" json:"created_date"`
}

func (MConfig) TableName() string {
	return "sys.m_config"
}
func (m *MConfig) BeforeCreate(trx *gorm.DB) (err error) {

	m.CreatedDate = time.Now()
	return nil
}
