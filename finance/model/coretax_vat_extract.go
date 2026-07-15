package model

import (
	"time"

	"gorm.io/gorm"
)

type CoretaxVatExtract struct {
	CoretaxVatExtractID *int64         `gorm:"column:coretax_vat_extract_id;default:nextval('acf.coretax_vat_extract_id_seq'::regclass);not null" json:"coretax_vat_extract_id"`
	InvoiceType         string         `gorm:"column:invoice_type" json:"invoice_type"`
	ExtractTotal        int            `gorm:"column:extract_total" json:"extract_total"`
	CreatedBy           int64          `gorm:"column:created_by" json:"created_by"`
	CreatedAt           *time.Time     `gorm:"column:created_at" json:"created_at"`
	DeletedBy           *int64         `gorm:"column:deleted_by" json:"deleted_by"`
	DeletedAt           gorm.DeletedAt `gorm:"column:deleted_at" json:"deleted_at"`
}

func (CoretaxVatExtract) TableName() string {
	return "acf.coretax_vat_extracts"
}

func (m *CoretaxVatExtract) BeforeCreate(trx *gorm.DB) (err error) {
	now := time.Now()
	if m.CreatedAt == nil {
		m.CreatedAt = &now

	}
	return
}

type CoretaxVatExtractDetail struct {
	VatExtractID int64  `gorm:"column:coretax_vat_extract_id" json:"coretax_vat_extract_id"`
	ReferenceID  string `gorm:"column:reference_id" json:"reference_id"`
	CustID       string `gorm:"cust_id" json:"cust_id"`
}

func (CoretaxVatExtractDetail) TableName() string {
	return "acf.coretax_vat_extracts_details"
}
