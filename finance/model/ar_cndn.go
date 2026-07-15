package model

import (
	"time"

	"gorm.io/gorm"
)

type ArCndn struct {
	CustId     string         `gorm:"column:cust_id" json:"cust_id"`
	ArCndnNo   string         `gorm:"column:ar_cndn_no" json:"ar_cndn_no"`
	ArCndnId   *int64         `gorm:"column:ar_cndn_id;primaryKey" json:"ar_cndn_id"`
	ArCndnDate *time.Time     `gorm:"column:ar_cndn_date" json:"ar_cndn_date"`
	TrCode     *string        `gorm:"column:tr_code" json:"tr_code"`
	DocNo      *string        `gorm:"column:doc_no" json:"doc_no"`
	CndnType   *string        `gorm:"column:cndn_type" json:"cndn_type"`
	CndnId     *int64         `gorm:"column:cndn_id" json:"cndn_id"`
	OutletId   *int64         `gorm:"column:outlet_id" json:"outlet_id"`
	CndnValue  *float64       `gorm:"column:cndn_value" json:"cndn_value"`
	CndnUsed   *float64       `gorm:"column:cndn_used" json:"cndn_used"`
	Notes      *string        `gorm:"column:notes" json:"notes"`
	DataStatus *int           `gorm:"column:data_status" json:"data_status"`
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

func (m *ArCndn) BeforeCreate(trx *gorm.DB) (err error) {
	now := time.Now()
	if m.IsPosted != nil {
		if *m.IsPosted {
			m.PostedAt = &now
		}
	}
	// intTmpsStr := now.UnixNano() / int64(time.Millisecond)
	// m.ChqTrNo = strconv.Itoa(int(intTmpsStr))
	m.CreatedAt = now
	m.UpdatedAt = now
	m.UpdatedBy = m.CreatedBy
	return nil
}

func (ArCndn) TableName() string {
	return "acf.ar_cndn"
}

type ArCndnList struct {
	CustId        string         `gorm:"column:cust_id" json:"cust_id"`
	ArCndnNo      string         `gorm:"column:ar_cndn_no" json:"ar_cndn_no"`
	ArCndnId      *int64         `gorm:"column:ar_cndn_id;primaryKey" json:"ar_cndn_id"`
	ArCndnDate    *time.Time     `gorm:"column:ar_cndn_date" json:"ar_cndn_date"`
	TrCode        *string        `gorm:"column:tr_code" json:"tr_code"`
	DocNo         *string        `gorm:"column:doc_no" json:"doc_no"`
	CndnType      *string        `gorm:"column:cndn_type" json:"cndn_type"`
	CndnId        *int64         `gorm:"column:cndn_id" json:"cndn_id"`
	CndnCode      *string        `gorm:"column:cndn_code" json:"cndn_code"`
	CndnName      *string        `gorm:"column:cndn_name" json:"cndn_name"`
	OutletId      *int64         `gorm:"column:outlet_id" json:"outlet_id"`
	OutletCode    *string        `gorm:"column:outlet_code" json:"outlet_code"`
	OutletName    *string        `gorm:"column:outlet_name" json:"outlet_name"`
	CndnValue     *float64       `gorm:"column:cndn_value" json:"cndn_value"`
	CndnUsed      *float64       `gorm:"column:cndn_used" json:"cndn_used"`
	Notes         *string        `gorm:"column:notes" json:"notes"`
	DataStatus    *int           `gorm:"column:data_status" json:"data_status"`
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

func (ArCndnList) TableName() string {
	return "acf.ar_cndn"
}
