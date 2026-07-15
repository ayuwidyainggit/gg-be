package model

import (
	"fmt"
	"log"
	"strings"
	"time"

	"gorm.io/gorm"
)

type ItemStCh struct {
	CustID      string         `gorm:"column:cust_id;primaryKey" json:"cust_id"`
	IscNo       string         `gorm:"column:isc_no;primaryKey" json:"isc_no"`
	IscDate     *time.Time     `gorm:"column:isc_date" json:"isc_date"`
	TrCode      *string        `gorm:"column:tr_code" json:"tr_code"`
	WhID        *int64         `gorm:"column:wh_id" json:"wh_id"`
	ItemCdnFrom *int64         `gorm:"column:item_cdn_from" json:"item_cdn_from"`
	ItemCdnTo   *int64         `gorm:"column:item_cdn_to" json:"item_cdn_to"`
	Notes       *string        `gorm:"column:notes" json:"notes"`
	DataStatus  *int64         `gorm:"column:data_status" json:"data_status"`
	CreatedBy   *int64         `gorm:"column:created_by" json:"created_by"`
	CreatedAt   time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedBy   *int64         `gorm:"column:updated_by" json:"updated_by"`
	UpdatedAt   time.Time      `gorm:"column:updated_at" json:"updated_at"`
	IsDel       bool           `gorm:"column:is_del" json:"is_del"`
	DeletedBy   *int64         `gorm:"column:deleted_by" json:"deleted_by"`
	DeletedAt   gorm.DeletedAt `gorm:"column:deleted_at" json:"deleted_at"`
	IsClosed    bool           `gorm:"column:is_closed" json:"is_closed"`
	ClosedBy    *int64         `gorm:"column:closed_by" json:"closed_by"`
	ClosedAt    time.Time      `gorm:"column:closed_at" json:"closed_at"`
}

func (ItemStCh) TableName() string {
	return "inv.item_st_ch"
}

type IscNo struct {
	IscNo string `gorm:"column:get_no_fn"`
}

func (m *ItemStCh) BeforeCreate(trx *gorm.DB) (err error) {
	// intTmpsStr := time.Now().UnixNano() / int64(time.Millisecond)
	// m.IscNo = strconv.Itoa(int(intTmpsStr))
	// m.CreatedAt = time.Now()
	// m.UpdatedAt = time.Now()
	// m.UpdatedBy = m.CreatedBy

	var iscNo IscNo
	trCode := *m.TrCode
	iscDateStr := m.IscDate.Format("2006-01-02")
	iscDateSubtr := iscDateStr[2:4]

	queryStr := fmt.Sprintf(`SELECT 
	TRIM(to_char(COALESCE(MAX(TO_NUMBER(SUBSTR(isc_no,6,5),'99999')),0)+1, '00000')) AS get_no_fn 
	FROM inv.item_st_ch
	WHERE substr(isc_no,4,2) = '%v' AND cust_id = '%v' AND tr_code = '%v'`, iscDateSubtr, strings.ToUpper(m.CustID), strings.ToUpper(trCode))
	err = trx.Raw(queryStr).Scan(&iscNo).Error
	if err != nil {
		return err
	}

	m.IscNo = trCode + iscDateSubtr + iscNo.IscNo
	log.Println("m.IscNo:", m.IscNo)

	m.CreatedAt = time.Now()
	m.UpdatedAt = time.Now()
	m.UpdatedBy = m.CreatedBy

	return nil
}

type ItemStChList struct {
	CustID        string         `gorm:"column:cust_id;primaryKey" json:"cust_id"`
	IscNo         string         `gorm:"column:isc_no;primaryKey" json:"isc_no"`
	IscDate       *time.Time     `gorm:"column:isc_date" json:"isc_date"`
	TrCode        *string        `gorm:"column:tr_code" json:"tr_code"`
	WhID          *int64         `gorm:"column:wh_id" json:"wh_id"`
	WhCode        *string        `gorm:"column:wh_code" json:"wh_code"`
	WhName        *string        `gorm:"column:wh_name" json:"wh_name"`
	ItemCdnFrom   *int64         `gorm:"column:item_cdn_from" json:"item_cdn_from"`
	ItemCdnTo     *int64         `gorm:"column:item_cdn_to" json:"item_cdn_to"`
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
	IsClosed      bool           `gorm:"column:is_closed" json:"is_closed"`
	ClosedBy      *int64         `gorm:"column:closed_by" json:"closed_by"`
	ClosedAt      time.Time      `gorm:"column:closed_at" json:"closed_at"`
}

func (ItemStChList) TableName() string {
	return "inv.item_st_ch"
}
