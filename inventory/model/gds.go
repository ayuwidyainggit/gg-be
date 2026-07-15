package model

import (
	"fmt"
	"log"
	"strings"
	"time"

	"gorm.io/gorm"
)

type Gds struct {
	CustID     string         `gorm:"column:cust_id" json:"cust_id"`
	GdsNo      string         `gorm:"column:gds_no" json:"gds_no"`
	GdsDate    *time.Time     `gorm:"column:gds_date" json:"gds_date"`
	TrCode     *string        `gorm:"column:tr_code" json:"tr_code"`
	RefNo      *string        `gorm:"column:ref_no" json:"ref_no"`
	WhID       *int64         `gorm:"column:wh_id" json:"wh_id"`
	SupID      *int64         `gorm:"column:sup_id" json:"sup_id"`
	SubTotal   *float64       `gorm:"column:sub_total" json:"sub_total"`
	Vat        *float64       `gorm:"column:vat" json:"vat"`
	VatValue   *float64       `gorm:"column:vat_value" json:"vat_value"`
	VatLg      *float64       `gorm:"column:vat_lg" json:"vat_lg"`
	VatLgValue *float64       `gorm:"column:vat_lg_value" json:"vat_lg_value"`
	Total      *float64       `gorm:"column:total" json:"total"`
	VatBg      *float64       `gorm:"column:vat_bg" json:"vat_bg"`
	VatBgValue *float64       `gorm:"column:vat_bg_value" json:"vat_bg_value"`
	TotEmbInc  *float64       `gorm:"column:tot_emb_inc" json:"tot_emb_inc"`
	TotEmbExc  *float64       `gorm:"column:tot_emb_exc" json:"tot_emb_exc"`
	Notes      *string        `gorm:"column:notes" json:"notes"`
	DataStatus *int64         `gorm:"column:data_status" json:"data_status"`
	CreatedBy  *int64         `gorm:"column:created_by" json:"created_by"`
	CreatedAt  time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedBy  *int64         `gorm:"column:updated_by" json:"updated_by"`
	UpdatedAt  time.Time      `gorm:"column:updated_at" json:"updated_at"`
	IsDel      bool           `gorm:"column:is_del" json:"is_del"`
	DeletedBy  *int64         `gorm:"column:deleted_by" json:"deleted_by"`
	DeletedAt  gorm.DeletedAt `gorm:"column:deleted_at" json:"deleted_at"`
	IsClosed   bool           `gorm:"column:is_closed" json:"is_closed"`
	ClosedBy   *int64         `gorm:"column:closed_by" json:"closed_by"`
	ClosedAt   time.Time      `gorm:"column:closed_at" json:"closed_at"`
}

type GdsNo struct {
	GdsNo string `gorm:"column:get_no_fn"`
}

func (m *Gds) BeforeCreate(trx *gorm.DB) (err error) {
	// intTmpsStr := time.Now().UnixNano() / int64(time.Millisecond)
	// m.GdsNo = strconv.Itoa(int(intTmpsStr))
	// m.CreatedAt = time.Now()
	// m.UpdatedAt = time.Now()
	// m.UpdatedBy = m.CreatedBy

	var gdsNo GdsNo
	trCode := *m.TrCode
	gdsDateStr := m.GdsDate.Format("2006-01-02")
	gdsDateSubtr := gdsDateStr[2:4]

	queryStr := fmt.Sprintf(`SELECT 
	TRIM(to_char(COALESCE(MAX(TO_NUMBER(SUBSTR(gds_no,6,5),'99999')),0)+1, '00000')) AS get_no_fn 
	FROM inv.gds
	WHERE substr(gds_no,4,2) = '%v' AND cust_id = '%v' AND tr_code = '%v'`, gdsDateSubtr, strings.ToUpper(m.CustID), strings.ToUpper(trCode))
	err = trx.Raw(queryStr).Scan(&gdsNo).Error
	if err != nil {
		return err
	}

	m.GdsNo = trCode + gdsDateSubtr + gdsNo.GdsNo
	log.Println("m.GdsNo:", m.GdsNo)

	m.CreatedAt = time.Now()
	m.UpdatedAt = time.Now()
	m.UpdatedBy = m.CreatedBy

	return nil
}

func (Gds) TableName() string {
	return "inv.gds"
}

type GdsList struct {
	CustID        string         `gorm:"column:cust_id" json:"cust_id"`
	GdsNo         string         `gorm:"column:gds_no" json:"gds_no"`
	GdsDate       *time.Time     `gorm:"column:gds_date" json:"gds_date"`
	TrCode        *string        `gorm:"column:tr_code" json:"tr_code"`
	RefNo         *string        `gorm:"column:ref_no" json:"ref_no"`
	WhID          *int64         `gorm:"column:wh_id" json:"wh_id"`
	WhCode        *string        `gorm:"column:wh_code" json:"wh_code"`
	WhName        *string        `gorm:"column:wh_name" json:"wh_name"`
	SupID         *int64         `gorm:"column:sup_id" json:"sup_id"`
	SupCode       *string        `gorm:"column:sup_code" json:"sup_code"`
	SupName       *string        `gorm:"column:sup_name" json:"sup_name"`
	SubTotal      *float64       `gorm:"column:sub_total" json:"sub_total"`
	Vat           *float64       `gorm:"column:vat" json:"vat"`
	VatValue      *float64       `gorm:"column:vat_value" json:"vat_value"`
	VatLg         *float64       `gorm:"column:vat_lg" json:"vat_lg"`
	VatLgValue    *float64       `gorm:"column:vat_lg_value" json:"vat_lg_value"`
	Total         *float64       `gorm:"column:total" json:"total"`
	VatBg         *float64       `gorm:"column:vat_bg" json:"vat_bg"`
	VatBgValue    *float64       `gorm:"column:vat_bg_value" json:"vat_bg_value"`
	TotEmbInc     *float64       `gorm:"column:tot_emb_inc" json:"tot_emb_inc"`
	TotEmbExc     *float64       `gorm:"column:tot_emb_exc" json:"tot_emb_exc"`
	Notes         *string        `gorm:"column:notes" json:"notes"`
	DataStatus    *int64         `gorm:"column:data_status" json:"data_status"`
	CreatedBy     *int64         `gorm:"column:created_by" json:"created_by"`
	CreatedAt     time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedBy     *int64         `gorm:"column:updated_by" json:"updated_by"`
	UpdatedByName string         `gorm:"column:updated_by_name" json:"updated_by_name"`
	UpdatedAt     *time.Time     `gorm:"column:updated_at" json:"updated_at"`
	IsDel         bool           `gorm:"column:is_del" json:"is_del"`
	DeletedBy     *int64         `gorm:"column:deleted_by" json:"deleted_by"`
	DeletedAt     gorm.DeletedAt `gorm:"column:deleted_at" json:"deleted_at"`
	IsClosed      bool           `gorm:"column:is_closed" json:"is_closed"`
	ClosedBy      *int64         `gorm:"column:closed_by" json:"closed_by"`
	ClosedAt      time.Time      `gorm:"column:closed_at" json:"closed_at"`
}

func (GdsList) TableName() string {
	return "inv.gds"
}
