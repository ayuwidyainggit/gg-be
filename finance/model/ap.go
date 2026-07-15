package model

import (
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"
)

type Ap struct {
	CustID        string         `gorm:"column:cust_id" json:"cust_id"`
	ApNo          string         `gorm:"column:ap_no;primaryKey" json:"ap_no"`
	ApDate        *time.Time     `gorm:"column:ap_date" json:"ap_date"`
	TrCode        *string        `gorm:"column:tr_code" json:"tr_code"`
	SupID         *int64         `gorm:"column:sup_id" json:"sup_id"`
	InvNo         *string        `gorm:"column:inv_no" json:"inv_no"`
	InvDate       *time.Time     `gorm:"column:inv_date" json:"inv_date"`
	InvDueDate    *time.Time     `gorm:"column:inv_due_date" json:"inv_due_date"`
	TaxInvNo      *string        `gorm:"column:tax_inv_no" json:"tax_inv_no"`
	TaxInvDate    *time.Time     `gorm:"column:tax_inv_date" json:"tax_inv_date"`
	TaxReturnNo   *string        `gorm:"column:tax_return_no" json:"tax_return_no"`
	TaxReturnDate *time.Time     `gorm:"column:tax_return_date" json:"tax_return_date"`
	SubTotal      *float64       `gorm:"column:sub_total" json:"sub_total"`
	MoneyPromo    *float64       `gorm:"column:money_promo" json:"money_promo"`
	InvDisc       *float64       `gorm:"column:inv_disc" json:"inv_disc"`
	InvDiscValue  *float64       `gorm:"column:inv_disc_value" json:"inv_disc_value"`
	SubTotalBtax  *float64       `gorm:"column:sub_total_btax" json:"sub_total_btax"`
	Vat           *float64       `gorm:"column:vat" json:"vat"`
	VatValue      *float64       `gorm:"column:vat_value" json:"vat_value"`
	VatLg         *float64       `gorm:"column:vat_lg" json:"vat_lg"`
	VatLgValue    *float64       `gorm:"column:vat_lg_value" json:"vat_lg_value"`
	VatBg         *float64       `gorm:"column:vat_bg" json:"vat_bg"`
	VatBgValue    *float64       `gorm:"column:vat_bg_value" json:"vat_bg_value"`
	DutyStamp     *float64       `gorm:"column:duty_stamp" json:"duty_stamp"`
	Total         *float64       `gorm:"column:total" json:"total"`
	TotalDiff     *float64       `gorm:"column:total_diff" json:"total_diff"`
	TotalRound    *float64       `gorm:"column:total_round" json:"total_round"`
	ApPaID        *float64       `gorm:"column:ap_paid" json:"ap_paid"`
	DataStatus    *int64         `gorm:"column:data_status" json:"data_status"`
	CreatedBy     *int64         `gorm:"column:created_by" json:"created_by"`
	CreatedAt     time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedBy     *int64         `gorm:"column:updated_by" json:"updated_by"`
	UpdatedAt     *time.Time     `gorm:"column:updated_at" json:"updated_at"`
	IsDel         bool           `gorm:"column:is_del" json:"is_del"`
	DeletedBy     *int64         `gorm:"column:deleted_by" json:"deleted_by"`
	DeletedAt     gorm.DeletedAt `gorm:"column:deleted_at" json:"deleted_at"`
	IsPosted      *bool          `gorm:"column:is_posted" json:"is_posted"`
	PostedAt      *time.Time     `gorm:"column:posted_at" json:"posted_at"`
}

func (Ap) TableName() string {
	return "acf.ap"
}

type ApNo struct {
	ApNo string `gorm:"column:get_no_fn"`
}

func (m *Ap) BeforeCreate(trx *gorm.DB) (err error) {
	var apNo ApNo
	trCode := *m.TrCode
	apDateStr := m.ApDate.Format("2006-01-02")
	apDateSubstr := apDateStr[2:4]

	queryStr := fmt.Sprintf(`SELECT 
	TRIM(to_char(COALESCE(MAX(TO_NUMBER(SUBSTR(ap_no,6,5),'99999')),0)+1, '00000')) AS get_no_fn 
	FROM acf.ap
	WHERE substr(ap_no,4,2) = '%v' AND cust_id = '%v' AND tr_code = '%v'`, apDateSubstr, strings.ToUpper(m.CustID), strings.ToUpper(trCode))
	err = trx.Raw(queryStr).Scan(&apNo).Error
	if err != nil {
		return err
	}

	now := time.Now()
	if m.IsPosted != nil {
		if *m.IsPosted {
			m.PostedAt = &now
		}
	}

	m.ApNo = trCode + apDateSubstr + apNo.ApNo
	m.CreatedAt = now
	m.UpdatedBy = m.CreatedBy
	return nil
}

type ApList struct {
	CustID        string         `gorm:"column:cust_id" json:"cust_id"`
	ApNo          string         `gorm:"column:ap_no;primaryKey" json:"ap_no"`
	ApDate        *time.Time     `gorm:"column:ap_date" json:"ap_date"`
	TrCode        *string        `gorm:"column:tr_code" json:"tr_code"`
	SupID         *int64         `gorm:"column:sup_id" json:"sup_id"`
	SupCode       *string        `gorm:"column:sup_code" json:"sup_code"`
	SupName       *string        `gorm:"column:sup_name" json:"sup_name"`
	InvNo         *string        `gorm:"column:inv_no" json:"inv_no"`
	InvDate       *time.Time     `gorm:"column:inv_date" json:"inv_date"`
	InvDueDate    *time.Time     `gorm:"column:inv_due_date" json:"inv_due_date"`
	TaxInvNo      *string        `gorm:"column:tax_inv_no" json:"tax_inv_no"`
	TaxInvDate    *time.Time     `gorm:"column:tax_inv_date" json:"tax_inv_date"`
	TaxReturnNo   *string        `gorm:"column:tax_return_no" json:"tax_return_no"`
	TaxReturnDate *time.Time     `gorm:"column:tax_return_date" json:"tax_return_date"`
	SubTotal      *float64       `gorm:"column:sub_total" json:"sub_total"`
	MoneyPromo    *float64       `gorm:"column:money_promo" json:"money_promo"`
	InvDisc       *float64       `gorm:"column:inv_disc" json:"inv_disc"`
	InvDiscValue  *float64       `gorm:"column:inv_disc_value" json:"inv_disc_value"`
	SubTotalBtax  *float64       `gorm:"column:sub_total_btax" json:"sub_total_btax"`
	Vat           *float64       `gorm:"column:vat" json:"vat"`
	VatValue      *float64       `gorm:"column:vat_value" json:"vat_value"`
	VatLg         *float64       `gorm:"column:vat_lg" json:"vat_lg"`
	VatLgValue    *float64       `gorm:"column:vat_lg_value" json:"vat_lg_value"`
	VatBg         *float64       `gorm:"column:vat_bg" json:"vat_bg"`
	VatBgValue    *float64       `gorm:"column:vat_bg_value" json:"vat_bg_value"`
	DutyStamp     *float64       `gorm:"column:duty_stamp" json:"duty_stamp"`
	Total         *float64       `gorm:"column:total" json:"total"`
	TotalDiff     *float64       `gorm:"column:total_diff" json:"total_diff"`
	TotalRound    *float64       `gorm:"column:total_round" json:"total_round"`
	ApPaID        *float64       `gorm:"column:ap_paid" json:"ap_paid"`
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

func (ApList) TableName() string {
	return "acf.ap"
}
