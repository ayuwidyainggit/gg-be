package model

import (
	"time"

	"gorm.io/gorm"
)

type VatExtract struct {
	CustID         string     `gorm:"cust_id" json:"cust_id"`
	VatExtractID   *int64     `gorm:"column:vat_extract_id;default:nextval('acf.vat_extract_id_seq'::regclass);not null" json:"vat_extract_id"`
	VatExtractType int        `gorm:"column:vat_extract_type" json:"vat_extract_type"`
	InvoiceType    string     `gorm:"column:invoice_type" json:"invoice_type"`
	ExtractTotal   int        `gorm:"column:extract_total" json:"extract_total"`
	CreatedBy      int64      `gorm:"column:created_by" json:"created_by"`
	CreatedAt      *time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedBy      int64      `gorm:"column:updated_by" json:"updated_by"`
	UpdatedAt      *time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func (VatExtract) TableName() string {
	return "acf.vat_extracts"
}

func (m *VatExtract) BeforeCreate(trx *gorm.DB) (err error) {
	now := time.Now()
	if m.CreatedAt == nil {
		m.CreatedAt = &now

	}
	m.UpdatedBy = m.CreatedBy
	return
}

type VatExtractDetail struct {
	VatExtractID int64 `gorm:"column:vat_extract_id" json:"vat_extract_id"`
	ReferenceID  int64 `gorm:"column:reference_id" json:"reference_id"`
}

func (VatExtractDetail) TableName() string {
	return "acf.vat_extract_details"
}

type VatExtractList struct {
	VatExtractID   *int64     `gorm:"column:vat_extract_id" json:"vat_extract_id"`
	VatExtractType int        `gorm:"column:vat_extract_type" json:"vat_extract_type"`
	InvoiceType    string     `gorm:"column:invoice_type" json:"invoice_type"`
	ExtractTotal   int        `gorm:"column:extract_total" json:"extract_total"`
	CreatedBy      int64      `gorm:"column:created_by" json:"created_by"`
	CreatedAt      *time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedBy      int64      `gorm:"column:updated_by" json:"updated_by"`
	UpdatedAt      *time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func (VatExtractList) TableName() string {
	return "acf.vat_extracts"
}

type VatExtractDetailList struct {
	ID                 uint           `gorm:"column:account_payable_id;primaryKey" json:"account_payable_id"`
	CustId             string         `gorm:"column:cust_id" json:"cust_id"`
	AccountPayableDate *time.Time     `gorm:"column:account_payable_date" json:"account_payable_date"`
	ApType             string         `gorm:"column:ap_type" json:"ap_type"`
	SupId              *int64         `gorm:"column:sup_id" json:"sup_id"`
	SupName            string         `gorm:"column:sup_name" json:"sup_name"`
	InvoiceNo          string         `gorm:"column:invoice_no;primaryKey" json:"invoice_no"`
	InvoiceDate        *time.Time     `gorm:"column:invoice_date" json:"invoice_date"`
	DocumentNo         string         `gorm:"column:document_no" json:"document_no"`
	TaxInvoiceDate     *time.Time     `gorm:"column:tax_invoice_date" json:"tax_invoice_date"`
	TaxInvoiceNo       *string        `gorm:"column:tax_invoice_no" json:"tax_invoice_no"`
	TaxReturnDate      *time.Time     `gorm:"column:tax_return_date" json:"tax_return_date"`
	TaxReturnNo        *string        `gorm:"column:tax_return_no" json:"tax_return_no"`
	DueDate            *time.Time     `gorm:"column:due_date" json:"due_date"`
	ReturnDate         *time.Time     `gorm:"column:return_date" json:"return_date"`
	Amount             *float64       `gorm:"column:amount" json:"amount"`
	DiscountRp         *float64       `gorm:"column:discount_rp" json:"discount_rp"`
	DiscountPercent    *float64       `gorm:"column:discount_percent" json:"discount_percent"`
	SubTotal           *float64       `gorm:"column:sub_total" json:"sub_total"`
	Vat                *float64       `gorm:"column:vat" json:"vat"`
	VatValue           *float64       `gorm:"column:vat_value" json:"vat_value"`
	VatLg              *float64       `gorm:"column:vat_lg" json:"vat_lg"`
	VatLgValue         *float64       `gorm:"column:vat_lg_value" json:"vat_lg_value"`
	Materai            *float64       `gorm:"column:materai" json:"materai"`
	Total              *float64       `gorm:"column:total" json:"total"`
	CreatedBy          *int64         `gorm:"column:created_by" json:"created_by"`
	CreatedAt          time.Time      `gorm:"column:created_at" json:"created_at"`
	CreatedByName      *string        `gorm:"column:created_by_name" json:"created_by_name"`
	UpdatedBy          *int64         `gorm:"column:updated_by" json:"updated_by"`
	UpdatedAt          time.Time      `gorm:"column:updated_at" json:"updated_at"`
	UpdatedByName      *string        `gorm:"column:updated_by_name" json:"updated_by_name"`
	DeletedBy          *int64         `gorm:"column:deleted_by" json:"deleted_by"`
	DeletedAt          gorm.DeletedAt `gorm:"column:deleted_at" json:"deleted_at"`
	IsDel              bool           `gorm:"column:is_del" json:"is_del"`
	ExtractStatus      string         `gorm:"column:extract_status" json:"extract_status"`
	ExtractedAt        *time.Time     `gorm:"column:extracted_at" json:"extracted_at"`
	Npwp               string         `gorm:"column:npwp" json:"npwp"`
	SupCode            string         `gorm:"column:sup_code" json:"sup_code"`
	Address            string         `gorm:"column:address" json:"address"`
	VatExtractType     int            `gorm:"column:vat_extract_type" json:"vat_extract_type"`
	InvoiceType        string         `gorm:"column:invoice_type" json:"invoice_type"`
}

func (VatExtractDetailList) TableName() string {
	return "acf.vat_extracts"
}
