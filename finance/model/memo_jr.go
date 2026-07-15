package model

import (
	"strconv"
	"time"

	"gorm.io/gorm"
)

type MemoJr struct {
	CustID     string         `gorm:"column:cust_id" json:"cust_id"`
	MjNo       string         `gorm:"column:mj_no" json:"mj_no"`
	MjDate     *time.Time     `gorm:"column:mj_date" json:"mj_date"`
	TrCode     *string        `gorm:"column:tr_code" json:"tr_code"`
	Notes      *string        `gorm:"column:notes" json:"notes"`
	DataStatus *int64         `gorm:"column:data_status" json:"data_status"`
	CreatedBy  *int64         `gorm:"column:created_by" json:"created_by"`
	CreatedAt  time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedBy  *int64         `gorm:"column:updated_by" json:"updated_by"`
	UpdatedAt  time.Time      `gorm:"column:updated_at" json:"updated_at"`
	IsDel      bool           `gorm:"column:is_del" json:"is_del"`
	DeletedBy  *int64         `gorm:"column:deleted_by" json:"deleted_by"`
	DeletedAt  gorm.DeletedAt `gorm:"column:deleted_at" json:"deleted_at"`
	IsPosted   *bool          `gorm:"column:is_posted" json:"is_posted"`
	PostedAt   *time.Time     `gorm:"column:posted_at" json:"posted_at"`
}

func (MemoJr) TableName() string {
	return "acf.memo_jr"
}
func (m *MemoJr) BeforeCreate(trx *gorm.DB) (err error) {
	intTmpsStr := time.Now().UnixNano() / int64(time.Millisecond)
	m.MjNo = strconv.Itoa(int(intTmpsStr))
	m.CreatedAt = time.Now()
	m.UpdatedAt = time.Now()
	m.UpdatedBy = m.CreatedBy
	return nil
}

type MemoJrList struct {
	CustID        string         `gorm:"column:cust_id" json:"cust_id"`
	MjNo          string         `gorm:"column:mj_no" json:"mj_no"`
	MjDate        *time.Time     `gorm:"column:mj_date" json:"mj_date"`
	TrCode        *string        `gorm:"column:tr_code" json:"tr_code"`
	Notes         *string        `gorm:"column:notes" json:"notes"`
	DataStatus    *int64         `gorm:"column:data_status" json:"data_status"`
	CreatedBy     *int64         `gorm:"column:created_by" json:"created_by"`
	CreatedAt     time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedBy     *int64         `gorm:"column:updated_by" json:"updated_by"`
	UpdatedAt     *time.Time     `gorm:"column:updated_at" json:"updated_at"`
	UpdatedByName *string        `gorm:"column:updated_by_name" json:"updated_by_name"`
	IsDel         bool           `gorm:"column:is_del" json:"is_del"`
	DeletedBy     *int64         `gorm:"column:deleted_by" json:"deleted_by"`
	DeletedAt     gorm.DeletedAt `gorm:"column:deleted_at" json:"deleted_at"`
	IsPosted      *bool          `gorm:"column:is_posted" json:"is_posted"`
	PostedAt      *time.Time     `gorm:"column:posted_at" json:"posted_at"`
}

func (MemoJrList) TableName() string {
	return "acf.memo_jr"
}
