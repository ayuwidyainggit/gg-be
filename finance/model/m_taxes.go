package model

import (
	"time"

	"gorm.io/gorm"
)

type MTaxes struct {
	MTaxID                *int64         `gorm:"m_tax_id;default:nextval('acf.m_tax_id_seq'::regclass);not null" json:"m_tax_id"`
	CustID                string         `gorm:"cust_id" json:"cust_id"`
	Year                  int            `gorm:"year" json:"year"`
	TransactionStatusCode string         `gorm:"transaction_status_code" json:"transaction_status_code"`
	SerialCode            string         `gorm:"serial_code" json:"serial_code"`
	SerialFrom            int            `gorm:"serial_from" json:"serial_from"`
	SerialTo              int            `gorm:"serial_to" json:"serial_to"`
	Sequence              int            `gorm:"sequence" json:"sequence"`
	TaxNumberAlert        int            `gorm:"tax_number_alert" json:"tax_number_alert"`
	Status                *int           `gorm:":status" json:"status"`
	CreatedBy             *int64         `gorm:"created_by" json:"created_by"`
	CreatedAt             time.Time      `gorm:"created_at" json:"created_at"`
	UpdatedBy             *int64         `gorm:"updated_by" json:"updated_by"`
	UpdatedAt             time.Time      `gorm:"updated_at" json:"updated_at"`
	IsDel                 *bool          `gorm:"is_del" json:"is_del"`
	DeletedBy             *int64         `gorm:"deleted_by" json:"deleted_by"`
	DeletedAt             gorm.DeletedAt `gorm:"deleted_at" json:"deleted_at"`
	RemainingQty          int            `gorm:"remaining_qty" json:"remaining_qty"`
	TotalTaxNo            int            `gorm:"total_tax_no" json:"total_tax_no"`
	LastGeneratedTax      string         `gorm:"last_generated_tax" json:"last_generated_tax"`
}

func (MTaxes) TableName() string {
	return "acf.m_taxes"
}

type MTaxesRead struct {
	MTaxID                *int64         `gorm:"m_tax_id" json:"m_tax_id"`
	CustID                string         `gorm:"cust_id" json:"cust_id"`
	Year                  int            `gorm:"year" json:"year"`
	TransactionStatusCode string         `gorm:"transaction_status_code" json:"transaction_status_code"`
	SerialCode            string         `gorm:"serial_code" json:"serial_code"`
	From                  int            `gorm:"serial_from" json:"serial_from"`
	To                    int            `gorm:"serial_to" json:"serial_to"`
	Sequence              int            `gorm:"sequence" json:"sequence"`
	TaxNumberAlert        int            `gorm:"tax_number_alert" json:"tax_number_alert"`
	TotalTaxNo            int            `gorm:"total_tax_no" json:"total_tax_no"`
	LastGeneratedTax      string         `gorm:"last_generated_tax" json:"last_generated_tax"`
	RemainingQty          int            `gorm:"remaining_qty" json:"remaining_qty"`
	Status                int            `gorm:"status" json:"status"`
	CreatedBy             *int64         `gorm:"created_by" json:"created_by"`
	CreatedAt             time.Time      `gorm:"created_at" json:"created_at"`
	UpdatedBy             *int64         `gorm:"updated_by" json:"updated_by"`
	UpdatedAt             time.Time      `gorm:"updated_at" json:"updated_at"`
	IsDel                 bool           `gorm:"is_del" json:"is_del"`
	DeletedBy             *int64         `gorm:"deleted_by" json:"deleted_by"`
	DeletedAt             gorm.DeletedAt `grom:"deleted_at" json:"deleted_at"`
}

func (MTaxesRead) TableName() string {
	return "acf.m_taxes"
}
