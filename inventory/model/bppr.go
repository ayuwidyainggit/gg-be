package model

import (
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"
)

type Bppr struct {
	CustID         *string        `gorm:"column:cust_id" json:"cust_id"`
	BpprNo         string         `gorm:"column:bppr_no" json:"bppr_no"`
	BpprDate       *time.Time     `gorm:"column:bppr_date" json:"bppr_date"`
	TrCode         *string        `gorm:"column:tr_code" json:"tr_code"`
	SupID          *int64         `gorm:"column:sup_id" json:"sup_id"`
	WhID           *int64         `gorm:"column:wh_id" json:"wh_id"`
	ItemCdn        *int64         `gorm:"column:item_cdn" json:"item_cdn"`
	ReturnReasonID *int64         `gorm:"column:return_reason_id" json:"return_reason_id"`
	ReturnNo       string         `gorm:"column:return_no" json:"return_no"`
	ReturnDate     *time.Time     `gorm:"column:return_date" json:"return_date,omitempty"`
	Notes          *string        `gorm:"column:notes" json:"notes"`
	TotEmbInc      *float64       `gorm:"column:tot_emb_inc" json:"tot_emb_inc"`
	TotEmbExc      *float64       `gorm:"column:tot_emb_exc" json:"tot_emb_exc"`
	SubTotal       *float64       `gorm:"column:sub_total" json:"sub_total"`
	Vat            *float64       `gorm:"column:vat" json:"vat"`
	VatValue       *float64       `gorm:"column:vat_value" json:"vat_value"`
	VatLg          *float64       `gorm:"column:vat_lg" json:"vat_lg"`
	VatLgValue     *float64       `gorm:"column:vat_lg_value" json:"vat_lg_value"`
	Total          *float64       `gorm:"column:total" json:"total"`
	VatBg          *float64       `gorm:"column:vat_bg" json:"vat_bg"`
	VatBgValue     *float64       `gorm:"column:vat_bg_value" json:"vat_bg_value"`
	DataStatus     *int64         `gorm:"column:data_status" json:"data_status"`
	CreatedBy      *int64         `gorm:"column:created_by" json:"created_by"`
	CreatedAt      time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedBy      *int64         `gorm:"column:updated_by" json:"updated_by"`
	UpdatedAt      time.Time      `gorm:"column:updated_at" json:"updated_at"`
	IsDel          bool           `gorm:"column:is_del" json:"is_del"`
	DeletedBy      *int64         `gorm:"column:deleted_by" json:"deleted_by"`
	DeletedAt      gorm.DeletedAt `gorm:"column:deleted_at" json:"deleted_at"`
	IsClosed       bool           `gorm:"column:is_closed" json:"is_closed"`
	ClosedBy       *int64         `gorm:"column:closed_by" json:"closed_by"`
	ClosedAt       time.Time      `gorm:"column:closed_at" json:"closed_at"`
}

func (Bppr) TableName() string {
	return "inv.bppr"
}

type BpprNo struct {
	BpprNo string `gorm:"column:get_no_fn"`
}

type ReturnNo struct {
	ReturnNo string `gorm:"column:get_no_fn"`
}

func (m *Bppr) BeforeCreate(trx *gorm.DB) (err error) {
	var bpprNo BpprNo
	trCode := *m.TrCode
	bpprDateStr := m.BpprDate.Format("2006-01-02")
	bpprDateSubtr := bpprDateStr[2:4]

	queryStr := fmt.Sprintf(`SELECT 
	TRIM(to_char(COALESCE(MAX(TO_NUMBER(SUBSTR(bppr_no,6,5),'99999')),0)+1, '00000')) AS get_no_fn 
	FROM inv.bppr
	WHERE substr(bppr_no,4,2) = '%v' AND cust_id = '%v' AND tr_code = '%v'`, bpprDateSubtr, strings.ToUpper(*m.CustID), strings.ToUpper(trCode))
	err = trx.Raw(queryStr).Scan(&bpprNo).Error
	if err != nil {
		return err
	}

	m.BpprNo = trCode + bpprDateSubtr + bpprNo.BpprNo
	m.CreatedAt = time.Now()
	m.UpdatedAt = time.Now()
	m.UpdatedBy = m.CreatedBy
	return nil
}

func (m *Bppr) BeforeUpdate(trx *gorm.DB) (err error) {
	now := time.Now()

	if *m.DataStatus == 2 {
		trCode := *m.TrCode
		returnDateStr := m.ReturnDate.Format("2006-01-02")
		returnDateSubtr := returnDateStr[2:4]
		var returnNo ReturnNo

		queryStr := fmt.Sprintf(`SELECT 
			TRIM(to_char(COALESCE(MAX(TO_NUMBER(SUBSTR(return_no,6,5),'99999')),0)+1, '00000')) AS get_no_fn 
			FROM inv.bppr
			WHERE substr(return_no,1,3) = '%v' AND substr(return_no,4,2) = '%v' AND cust_id = '%v'`,
			trCode, returnDateSubtr, strings.ToUpper(*m.CustID))
		err = trx.Raw(queryStr).Scan(&returnNo).Error
		if err != nil {
			return err
		}

		newReturnNo := trCode + returnDateSubtr + returnNo.ReturnNo
		trx.Statement.SetColumn("return_no", string(newReturnNo))
	}

	m.UpdatedAt = now
	m.TrCode = nil

	return nil
}

type BpprList struct {
	CustID         string         `gorm:"column:cust_id" json:"cust_id"`
	BpprNo         string         `gorm:"column:bppr_no" json:"bppr_no"`
	BpprDate       *time.Time     `gorm:"column:bppr_date" json:"bppr_date"`
	TrCode         *string        `gorm:"column:tr_code" json:"tr_code"`
	SupID          *int64         `gorm:"column:sup_id" json:"sup_id"`
	SupCode        *string        `gorm:"column:sup_code" json:"sup_code"`
	SupName        *string        `gorm:"column:sup_name" json:"sup_name"`
	WhID           *int64         `gorm:"column:wh_id" json:"wh_id"`
	WhCode         *string        `gorm:"column:wh_code" json:"wh_code"`
	WhName         *string        `gorm:"column:wh_name" json:"wh_name"`
	ItemCdn        *int64         `gorm:"column:item_cdn" json:"item_cdn"`
	ReturnReasonID *int64         `gorm:"column:return_reason_id" json:"return_reason_id"`
	ReturnNo       *string        `gorm:"column:return_no" json:"return_no"`
	ReturnDate     *time.Time     `gorm:"column:return_date" json:"return_date,omitempty"`
	Notes          *string        `gorm:"column:notes" json:"notes"`
	TotEmbInc      *float64       `gorm:"column:tot_emb_inc" json:"tot_emb_inc"`
	TotEmbExc      *float64       `gorm:"column:tot_emb_exc" json:"tot_emb_exc"`
	SubTotal       *float64       `gorm:"column:sub_total" json:"sub_total"`
	Vat            *float64       `gorm:"column:vat" json:"vat"`
	VatValue       *float64       `gorm:"column:vat_value" json:"vat_value"`
	VatLg          *float64       `gorm:"column:vat_lg" json:"vat_lg"`
	VatLgValue     *float64       `gorm:"column:vat_lg_value" json:"vat_lg_value"`
	Total          *float64       `gorm:"column:total" json:"total"`
	VatBg          *float64       `gorm:"column:vat_bg" json:"vat_bg"`
	VatBgValue     *float64       `gorm:"column:vat_bg_value" json:"vat_bg_value"`
	DataStatus     int64          `gorm:"column:data_status" json:"data_status"`
	CreatedBy      *int64         `gorm:"column:created_by" json:"created_by"`
	CreatedAt      time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedBy      *int64         `gorm:"column:updated_by" json:"updated_by"`
	UpdatedByName  *string        `gorm:"column:updated_by_name" json:"updated_by_name"`
	UpdatedAt      time.Time      `gorm:"column:updated_at" json:"updated_at"`
	IsDel          bool           `gorm:"column:is_del" json:"is_del"`
	DeletedBy      *int64         `gorm:"column:deleted_by" json:"deleted_by"`
	DeletedAt      gorm.DeletedAt `gorm:"column:deleted_at" json:"deleted_at"`
	IsClosed       bool           `gorm:"column:is_closed" json:"is_closed"`
	ClosedBy       *int64         `gorm:"column:closed_by" json:"closed_by"`
	ClosedAt       time.Time      `gorm:"column:closed_at" json:"closed_at"`
}

func (BpprList) TableName() string {
	return "inv.bppr"
}
