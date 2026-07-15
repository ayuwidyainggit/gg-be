package model

import (
	"strconv"
	"time"

	"gorm.io/gorm"
)

type ApCndn struct {
	CustID     string         `gorm:"column:cust_id" json:"cust_id"`
	ApCndnNo   string         `gorm:"column:ap_cndn_no;primaryKey" json:"ap_cndn_no"`
	ApCndnDate *time.Time     `gorm:"column:ap_cndn_date" json:"ap_cndn_date"`
	TrCode     *string        `gorm:"column:tr_code" json:"tr_code"`
	DocNo      *string        `gorm:"column:doc_no" json:"doc_no"`
	CndnType   *string        `gorm:"column:cndn_type" json:"cndn_type"`
	CndnID     *int64         `gorm:"column:cndn_id" json:"cndn_id"`
	SupID      *int64         `gorm:"column:sup_id" json:"sup_id"`
	ApNo       *string        `gorm:"column:ap_no" json:"ap_no"`
	CndnValue  *float64       `gorm:"column:cndn_value" json:"cndn_value"`
	CndnUsed   *float64       `gorm:"column:cndn_used" json:"cndn_used"`
	Notes      *string        `gorm:"column:notes" json:"notes"`
	DataStatus *int64         `gorm:"column:data_status" json:"data_status"`
	CreatedBy  *int64         `gorm:"column:created_by" json:"created_by"`
	CreatedAt  time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedBy  *int64         `gorm:"column:updated_by" json:"updated_by"`
	UpdatedAt  *time.Time     `gorm:"column:updated_at" json:"updated_at"`
	IsDel      bool           `gorm:"column:is_del" json:"is_del"`
	DeletedBy  *int64         `gorm:"column:deleted_by" json:"deleted_by"`
	DeletedAt  gorm.DeletedAt `gorm:"column:deleted_at" json:"deleted_at"`
	IsPosted   *bool          `gorm:"column:is_posted" json:"is_posted"`
	PostedAt   *time.Time     `gorm:"column:posted_at" json:"posted_at"`
}

func (m *ApCndn) BeforeCreate(trx *gorm.DB) (err error) {
	now := time.Now()
	if m.IsPosted != nil {
		if *m.IsPosted {
			m.PostedAt = &now
		}
	}
	intTmpsStr := now.UnixNano() / int64(time.Millisecond)
	m.ApCndnNo = strconv.Itoa(int(intTmpsStr))
	m.CreatedAt = now
	m.UpdatedBy = m.CreatedBy
	return nil
}

func (ApCndn) TableName() string {
	return "acf.ap_cndn"
}

type ApCndnList struct {
	CustID        string         `gorm:"column:cust_id" json:"cust_id"`
	ApCndnNo      string         `gorm:"column:ap_cndn_no;primaryKey" json:"ap_cndn_no"`
	ApCndnDate    *time.Time     `gorm:"column:ap_cndn_date" json:"ap_cndn_date"`
	TrCode        *string        `gorm:"column:tr_code" json:"tr_code"`
	DocNo         *string        `gorm:"column:doc_no" json:"doc_no"`
	CndnType      *string        `gorm:"column:cndn_type" json:"cndn_type"`
	CndnID        *int64         `gorm:"column:cndn_id" json:"cndn_id"`
	CndnCode      *string        `gorm:"column:cndn_code" json:"cndn_code"`
	CndnName      *string        `gorm:"column:cndn_name" json:"cndn_name"`
	SupID         *int64         `gorm:"column:sup_id" json:"sup_id"`
	SupCode       *string        `gorm:"column:sup_code" json:"sup_code"`
	SupName       *string        `gorm:"column:sup_name" json:"sup_name"`
	ApNo          *string        `gorm:"column:ap_no" json:"ap_no"`
	CndnValue     *float64       `gorm:"column:cndn_value" json:"cndn_value"`
	CndnUsed      *float64       `gorm:"column:cndn_used" json:"cndn_used"`
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

func (ApCndnList) TableName() string {
	return "acf.ap_cndn"
}
