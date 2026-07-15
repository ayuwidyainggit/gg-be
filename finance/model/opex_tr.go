package model

import (
	"strconv"
	"time"

	"gorm.io/gorm"
)

type OpexTr struct {
	CustID     string         `gorm:"column:cust_id" json:"cust_id"`
	OpexTrNo   string         `gorm:"column:opex_tr_no;primaryKey" json:"opex_tr_no"`
	OpexTrDate *time.Time     `gorm:"column:opex_tr_date" json:"opex_tr_date"`
	TrCode     *string        `gorm:"column:tr_code" json:"tr_code"`
	Notes      *string        `gorm:"column:notes" json:"notes"`
	TotAmount  *float64       `gorm:"column:tot_amount" json:"tot_amount"`
	DataStatus *int64         `gorm:"column:data_status" json:"data_status"`
	CreatedBy  *int64         `gorm:"column:created_by" json:"created_by"`
	CreatedAt  time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedBy  *int64         `gorm:"column:updated_by" json:"updated_by"`
	UpdatedAt  time.Time      `gorm:"column:updated_at" json:"updated_at"`
	IsDel      bool           `gorm:"column:is_del" json:"is_del"`
	DeletedBy  *int64         `gorm:"column:deleted_by" json:"deleted_by"`
	DeletedAt  gorm.DeletedAt `gorm:"column:deleted_at" json:"deleted_at"`
}

func (OpexTr) TableName() string {
	return "acf.opex_tr"
}
func (m *OpexTr) BeforeCreate(trx *gorm.DB) (err error) {
	intTmpsStr := time.Now().UnixNano() / int64(time.Millisecond)
	m.OpexTrNo = strconv.Itoa(int(intTmpsStr))
	m.CreatedAt = time.Now()
	m.UpdatedAt = time.Now()
	m.UpdatedBy = m.CreatedBy
	return nil
}

type OpexTrList struct {
	CustID        string         `gorm:"column:cust_id" json:"cust_id"`
	OpexTrNo      string         `gorm:"column:opex_tr_no;primaryKey" json:"opex_tr_no"`
	OpexTrDate    *time.Time     `gorm:"column:opex_tr_date" json:"opex_tr_date"`
	TrCode        *string        `gorm:"column:tr_code" json:"tr_code"`
	Notes         *string        `gorm:"column:notes" json:"notes"`
	TotAmount     *float64       `gorm:"column:tot_amount" json:"tot_amount"`
	DataStatus    *int64         `gorm:"column:data_status" json:"data_status"`
	CreatedBy     *int64         `gorm:"column:created_by" json:"created_by"`
	CreatedAt     time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedBy     *int64         `gorm:"column:updated_by" json:"updated_by"`
	UpdatedAt     time.Time      `gorm:"column:updated_at" json:"updated_at"`
	UpdatedByName *string        `gorm:"column:updated_by_name" json:"updated_by_name"`
	IsDel         bool           `gorm:"column:is_del" json:"is_del"`
	DeletedBy     *int64         `gorm:"column:deleted_by" json:"deleted_by"`
	DeletedAt     gorm.DeletedAt `gorm:"column:deleted_at" json:"deleted_at"`
}

func (OpexTrList) TableName() string {
	return "acf.opex_tr"
}
