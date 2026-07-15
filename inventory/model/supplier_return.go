package model

import (
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"
)

type SupplierReturn struct {
	CustID             string         `gorm:"column:cust_id" json:"cust_id"`
	SupplierReturnNo   string         `gorm:"column:supplier_return_no" json:"supplier_return_no"`
	SupplierReturnDate *time.Time     `gorm:"column:supplier_return_date" json:"supplier_return_date"`
	InvoiceNo          string         `gorm:"column:invoice_no" json:"invoice_no"`
	SupID              *int64         `gorm:"column:sup_id" json:"sup_id"`
	WhID               *int64         `gorm:"column:wh_id" json:"wh_id"`
	Notes              *string        `gorm:"column:notes" json:"notes"`
	SubTotal           *float64       `gorm:"column:sub_total" json:"sub_total"`
	Total              *float64       `gorm:"column:total" json:"total"`
	DiscountValue      *float64       `gorm:"column:discount_value" json:"discount_value"`
	VatValue           *float64       `gorm:"column:vat_value" json:"vat_value"`
	VatLgValue         *float64       `gorm:"column:vat_lg_value" json:"vat_lg_value"`
	VatBgValue         *float64       `gorm:"column:vat_bg_value" json:"vat_bg_value"`
	CreatedBy          *int64         `gorm:"column:created_by" json:"created_by"`
	CreatedAt          time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedBy          *int64         `gorm:"column:updated_by" json:"updated_by"`
	UpdatedAt          time.Time      `gorm:"column:updated_at" json:"updated_at"`
	IsDel              bool           `gorm:"column:is_del" json:"is_del"`
	DeletedBy          *int64         `gorm:"column:deleted_by" json:"deleted_by"`
	DeletedAt          gorm.DeletedAt `gorm:"column:deleted_at" json:"deleted_at"`
	IsClosed           bool           `gorm:"column:is_closed" json:"is_closed"`
	ClosedBy           *int64         `gorm:"column:closed_by" json:"closed_by"`
	ClosedAt           time.Time      `gorm:"column:closed_at" json:"closed_at"`
}

func (SupplierReturn) TableName() string {
	return "inv.supplier_returns"
}

type SupplierReturnNo struct {
	SupplierReturnNo string `gorm:"column:supplier_return_no_fn"`
}

func (m *SupplierReturn) BeforeCreate(trx *gorm.DB) (err error) {
	var SupplierReturnNo SupplierReturnNo
	trCode := "PR"
	ReturnDateStr := m.SupplierReturnDate.Format("2006-01-02")
	ReturnDateSubtr := ReturnDateStr[2:4] + ReturnDateStr[5:7] + ReturnDateStr[8:10]

	// log.Println("grDateStr:", grDateStr)
	// log.Println("grDateSubtr:", grDateSubtr)

	queryStr := fmt.Sprintf(`SELECT 
	TRIM(to_char(COALESCE(MAX(TO_NUMBER(SUBSTR(supplier_return_no,9,4),'9999')),0)+1, '0000')) AS supplier_return_no_fn 
	FROM inv.supplier_returns
	WHERE substr(supplier_return_no,3,6) = '%v' AND cust_id = '%v'`, ReturnDateSubtr, strings.ToUpper(m.CustID))
	err = trx.Raw(queryStr).Scan(&SupplierReturnNo).Error
	if err != nil {
		return err
	}

	// log.Println("grNo:", grNo.GrNo)

	m.SupplierReturnNo = trCode + ReturnDateSubtr + SupplierReturnNo.SupplierReturnNo
	// log.Println("m.GrNo:", m.GrNo)
	m.CreatedAt = time.Now()
	m.UpdatedAt = time.Now()
	m.UpdatedBy = m.CreatedBy
	return nil
}

type SupplierReturnGet struct {
	CustID             string         `gorm:"column:cust_id" json:"cust_id"`
	SupplierReturnNo   string         `gorm:"column:supplier_return_no" json:"supplier_return_no"`
	GrNO               string         `gorm:"column:gr_no" json:"gr_no"`
	InvoiceNo          *string        `gorm:"column:invoice_no" json:"invoice_no"`
	InvoiceDate        *time.Time     `gorm:"column:invoice_date" json:"invoice_date"`
	TaxInvoiceDate     *time.Time     `gorm:"column:tax_invoice_date"  json:"tax_invoice_date"`
	TaxInvoiceNo       string         `gorm:"column:tax_invoice_no"  json:"tax_invoice_no"`
	DueDate            *time.Time     `gorm:"column:due_date"  json:"due_date"`
	SupplierReturnDate *time.Time     `gorm:"column:supplier_return_date" json:"supplier_return_date"`
	SupID              *int64         `gorm:"column:sup_id" json:"sup_id"`
	SupCode            *string        `gorm:"column:sup_code" json:"sup_code"`
	SupName            *string        `gorm:"column:sup_name" json:"sup_name"`
	WhID               *int64         `gorm:"column:wh_id" json:"wh_id"`
	WhCode             *string        `gorm:"column:wh_code" json:"wh_code"`
	WhName             *string        `gorm:"column:wh_name" json:"wh_name"`
	Notes              *string        `gorm:"column:notes" json:"notes"`
	Vat                float64        `gorm:"column:vat" json:"vat"`
	VatValue           float64        `gorm:"column:vat_value" json:"vat_value"`
	VatLg              float64        `gorm:"column:vat_lg" json:"vat_lg"`
	VatLgValue         float64        `gorm:"column:vat_lg_value" json:"vat_lg_value"`
	VatBg              float64        `gorm:"column:vat_bg" json:"vat_bg"`
	VatBgValue         float64        `gorm:"column:vat_bg_value" json:"vat_bg_value"`
	SubTotal           float64        `gorm:"column:sub_total" json:"sub_total"`
	Total              float64        `gorm:"column:total" json:"total"`
	DiscountValue      float64        `gorm:"column:discount_value" json:"discount_value"`
	DataStatus         int64          `gorm:"column:data_status" json:"status"`
	CreatedBy          *int64         `gorm:"column:created_by" json:"created_by"`
	CreatedAt          time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedBy          *int64         `gorm:"column:updated_by" json:"updated_by"`
	UpdatedByName      *string        `gorm:"column:updated_by_name" json:"updated_by_name"`
	UpdatedAt          time.Time      `gorm:"column:updated_at" json:"updated_at"`
	IsDel              bool           `gorm:"column:is_del" json:"is_del"`
	DeletedBy          *int64         `gorm:"column:deleted_by" json:"deleted_by"`
	DeletedAt          gorm.DeletedAt `gorm:"column:deleted_at" json:"deleted_at"`
	IsClosed           bool           `gorm:"column:is_closed" json:"is_closed"`
	ClosedBy           *int64         `gorm:"column:closed_by" json:"closed_by"`
	ClosedAt           time.Time      `gorm:"column:closed_at" json:"closed_at"`
}

func (SupplierReturnGet) TableName() string {
	return "inv.supplier_returns"
}

type ReturnSuppliers struct {
	SupId   *int64  `gorm:"column:sup_id" json:"sup_id"`
	SupCode *string `gorm:"column:sup_code" json:"sup_code"`
	SupName *string `gorm:"column:sup_name" json:"sup_name"`
}

func (ReturnSuppliers) TableName() string {
	return "inv.supplier_returns"
}
