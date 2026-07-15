package model

import (
	"fmt"
	"log"
	"strings"
	"time"

	"gorm.io/gorm"
)

type SampleIssue struct {
	CustID     string         `gorm:"column:cust_id" json:"cust_id"`
	SmpIssNo   string         `gorm:"column:smp_iss_no" json:"smp_iss_no"`
	SmpIssDate *time.Time     `gorm:"column:smp_iss_date" json:"smp_iss_date"`
	TrCode     *string        `gorm:"column:tr_code" json:"tr_code"`
	WhID       *int64         `gorm:"column:wh_id" json:"wh_id"`
	CndnID     *int64         `gorm:"column:cndn_id" json:"cndn_id"`
	OutletID   *int64         `gorm:"column:outlet_id" json:"outlet_id"`
	Notes      *string        `gorm:"column:notes" json:"notes"`
	SubTotal   *float64       `gorm:"column:sub_total" json:"sub_total"`
	Vat        *float64       `gorm:"column:vat" json:"vat"`
	VatValue   *float64       `gorm:"column:vat_value" json:"vat_value"`
	VatLg      *float64       `gorm:"column:vat_lg" json:"vat_lg"`
	VatLgValue *float64       `gorm:"column:vat_lg_value" json:"vat_lg_value"`
	Total      *float64       `gorm:"column:total" json:"total"`
	VatBg      *float64       `gorm:"column:vat_bg" json:"vat_bg"`
	VatBgValue *float64       `gorm:"column:vat_bg_value" json:"vat_bg_value"`
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

func (SampleIssue) TableName() string {
	return "inv.sample_issue"
}

type SmpIssNo struct {
	SmpIssNo string `gorm:"column:get_no_fn"`
}

func (m *SampleIssue) BeforeCreate(trx *gorm.DB) (err error) {
	var smpIssNo SmpIssNo
	trCode := *m.TrCode
	smpIssDateStr := m.SmpIssDate.Format("2006-01-02")
	smpIssDateSubtr := smpIssDateStr[2:4]

	queryStr := fmt.Sprintf(`SELECT 
	TRIM(to_char(COALESCE(MAX(TO_NUMBER(SUBSTR(smp_iss_no,6,5),'99999')),0)+1, '00000')) AS get_no_fn 
	FROM inv.sample_issue
	WHERE substr(smp_iss_no,4,2) = '%v' AND cust_id = '%v' AND tr_code = '%v'`, smpIssDateSubtr, strings.ToUpper(m.CustID), strings.ToUpper(trCode))
	err = trx.Raw(queryStr).Scan(&smpIssNo).Error
	if err != nil {
		return err
	}

	m.SmpIssNo = trCode + smpIssDateSubtr + smpIssNo.SmpIssNo
	log.Println("m.SmpIssNo:", m.SmpIssNo)

	m.CreatedAt = time.Now()
	m.UpdatedAt = time.Now()
	m.UpdatedBy = m.CreatedBy
	return nil
}

type SampleIssueList struct {
	CustID        string         `gorm:"column:cust_id" json:"cust_id"`
	SmpIssNo      string         `gorm:"column:smp_iss_no" json:"smp_iss_no"`
	SmpIssDate    *time.Time     `gorm:"column:smp_iss_date" json:"smp_iss_date"`
	TrCode        *string        `gorm:"column:tr_code" json:"tr_code"`
	WhID          *int64         `gorm:"column:wh_id" json:"wh_id"`
	WhCode        *string        `gorm:"column:wh_code" json:"wh_code"`
	WhName        *string        `gorm:"column:wh_name" json:"wh_name"`
	CndnID        *int64         `gorm:"column:cndn_id" json:"cndn_id"`
	CndnCode      *string        `gorm:"column:cndn_code" json:"cndn_code"`
	CndnName      *string        `gorm:"column:cndn_name" json:"cndn_name"`
	OutletID      *int64         `gorm:"column:outlet_id" json:"outlet_id"`
	Notes         *string        `gorm:"column:notes" json:"notes"`
	SubTotal      *float64       `gorm:"column:sub_total" json:"sub_total"`
	Vat           *float64       `gorm:"column:vat" json:"vat"`
	VatValue      *float64       `gorm:"column:vat_value" json:"vat_value"`
	VatLg         *float64       `gorm:"column:vat_lg" json:"vat_lg"`
	VatLgValue    *float64       `gorm:"column:vat_lg_value" json:"vat_lg_value"`
	Total         *float64       `gorm:"column:total" json:"total"`
	VatBg         *float64       `gorm:"column:vat_bg" json:"vat_bg"`
	VatBgValue    *float64       `gorm:"column:vat_bg_value" json:"vat_bg_value"`
	DataStatus    *int64         `gorm:"column:data_status" json:"data_status"`
	CreatedBy     int64          `gorm:"column:created_by" json:"created_by"`
	CreatedAt     time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedBy     *int64         `gorm:"column:updated_by" json:"updated_by"`
	UpdatedByName *string        `gorm:"column:updated_by_name" json:"updated_by_name"`
	UpdatedAt     *time.Time     `gorm:"column:updated_at" json:"updated_at"`
	IsDel         bool           `gorm:"column:is_del" json:"is_del"`
	DeletedBy     *int64         `gorm:"column:deleted_by" json:"deleted_by"`
	DeletedAt     gorm.DeletedAt `gorm:"column:deleted_at" json:"deleted_at"`
	IsClosed      bool           `gorm:"column:is_closed" json:"is_closed"`
	ClosedBy      *int64         `gorm:"column:closed_by" json:"closed_by"`
	ClosedAt      time.Time      `gorm:"column:closed_at" json:"closed_at"`
}

func (SampleIssueList) TableName() string {
	return "inv.sample_issue"
}
