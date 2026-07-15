package model

import "time"

type Taxes struct {
	CustID    string    `gorm:"cust_id" json:"cust_id"`
	TaxesId   *int      `gorm:"column:taxes_id;default:nextval('acf.taxes_taxes_id_seq'::regclass);not null" json:"taxes_id"`
	MTaxId    int64     `gorm:"column:m_tax_id" json:"m_tax_id"`
	TaxNo     string    `gorm:"column:tax_no" json:"tax_no"`
	Status    int       `gorm:"column:status" json:"status"`
	InvoiceNo string    `gorm:"column:invoice_no" json:"invoice_no"`
	CreatedBy *int64    `gorm:"column:created_by" json:"created_by"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedBy *int64    `gorm:"column:updated_by" json:"updated_by"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func (Taxes) TableName() string {
	return "acf.taxes"
}

type TaxesGenerateRead struct {
	CustID      string     `gorm:"cust_id" json:"cust_id"`
	MTaxId      int64      `gorm:"column:m_tax_id" json:"m_tax_id"`
	TaxID       int64      `gorm:"column:taxes_id" json:"taxes_id"`
	TaxNo       string     `gorm:"column:tax_no" json:"tax_no"`
	InvoiceNo   string     `gorm:"column:invoice_no" json:"invoice_no"`
	InvoiceDate *time.Time `gorm:"column:invoice_date" json:"invoice_date"`
	Status      int        `gorm:"column:status" json:"status"`
	CreatedBy   *int64     `gorm:"column:created_by" json:"created_by"`
	CreatedAt   time.Time  `gorm:"column:created_at" json:"created_at"`
	UpdatedBy   *int64     `gorm:"column:updated_by" json:"updated_by"`
	UpdatedAt   time.Time  `gorm:"column:updated_at" json:"updated_at"`
}

func (TaxesGenerateRead) TableName() string {
	return "acf.taxes"
}

type TaxesList struct {
	TaxesId int    `gorm:"column:taxes_id" json:"taxes_id"`
	MTaxId  int    `gorm:"column:m_tax_id" json:"m_tax_id"`
	TaxNo   string `gorm:"column:tax_no" json:"tax_no"`
	Npwp    string `gorm:"column:npwp" json:"npwp"`
	Type    string `gorm:"type" json:"type"`

	RoDate          *time.Time `gorm:"ro_date" json:"ro_date"`
	ReturnDate      *time.Time `gorm:"return_date" json:"return_date"`
	SalesmanId      *int64     `gorm:"salesman_id" json:"salesman_id"`
	SalesmanCode    *string    `gorm:"salesman_code" json:"salesman_code"`
	SalesName       *string    `gorm:"sales_name" json:"sales_name"`
	OutletID        *int64     `gorm:"outlet_id" json:"outlet_id"`
	OutletCode      *string    `gorm:"column:outlet_code" json:"outlet_code"`
	OutletName      *string    `gorm:"column:outlet_name" json:"outlet_name"`
	OutletLatitude  *string    `gorm:"column:outlet_latitude" json:"outlet_latitude"`
	OutletLongitude *string    `gorm:"column:outlet_longitude" json:"outlet_longitude"`
	OutletAddress   *string    `gorm:"column:outlet_address" json:"outlet_address"`
	OrderNo         string     `gorm:"order_no" json:"order_no"`
	ReturnNo        string     `gorm:"return_no" json:"return_no"`
	PoNo            *string    `gorm:"po_no" json:"po_no"`
	PayType         *int64     `gorm:"pay_type" json:"pay_type"`
	MobileID        *int64     `gorm:"mobile_id" json:"mobile_id"`
	SubTotal        *float64   `gorm:"sub_total" json:"sub_total"`
	Vat             *float64   `gorm:"vat" json:"vat"`
	VatValue        *float64   `gorm:"vat_value" json:"vat_value"`
	Total           *float64   `gorm:"total" json:"total"`
	DataStatus      *int64     `gorm:"data_status" json:"data_status"`
	InvoiceNo       *string    `gorm:"invoice_no" json:"invoice_no"`
	InvoiceDate     *time.Time `gorm:"invoice_date" json:"invoice_date"`
	Status          *int       `gorm:"column:status" json:"status"`
}
