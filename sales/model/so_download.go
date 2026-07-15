package model

import "time"

type SoDownloadPo struct {
	OrderNo          *string    `gorm:"column:order_no" json:"order_no"`
	PoNo             *string    `gorm:"column:po_no" json:"po_no"`
	SoNo             string     `gorm:"column:so_no" json:"so_no"`
	RoDate           *time.Time `gorm:"column:ro_date" json:"ro_date"`
	InvoiceDate      *time.Time `gorm:"column:invoice_date" json:"invoice_date"`
	InvoiceNo        *string    `gorm:"column:invoice_no" json:"invoice_no"`
	OutletCode       *string    `gorm:"column:outlet_code" json:"outlet_code"`
	OutletName       *string    `gorm:"column:outlet_name" json:"outlet_name"`
	SalesmanCode     *string    `gorm:"column:salesman_code" json:"salesman_code"`
	EmployeeName     *string    `gorm:"column:employee_name" json:"employee_name"`
	SupplierCode     *string    `gorm:"column:supplier_code" json:"supplier_code"`
	SupplierName     *string    `gorm:"column:supplier_name" json:"supplier_name"`
	ProductCode      string     `gorm:"column:product_code" json:"product_code"`
	ProductName      string     `gorm:"column:product_name" json:"product_name"`
	UnitId3          *string    `gorm:"column:unit_id3" json:"unit_id3"`
	UnitId2          *string    `gorm:"column:unit_id2" json:"unit_id2"`
	UnitId1          *string    `gorm:"column:unit_id1" json:"unit_id1"`
	SellPriceSystem3 *float64   `gorm:"column:sell_price_system3" json:"sell_price_system3"`
	SellPriceSystem2 *float64   `gorm:"column:sell_price_system2" json:"sell_price_system2"`
	SellPriceSystem1 *float64   `gorm:"column:sell_price_system1" json:"sell_price_system1"`
	SellPricePo3     *float64   `gorm:"column:sell_price_po3" json:"sell_price_po3"`
	SellPricePo2     *float64   `gorm:"column:sell_price_po2" json:"sell_price_po2"`
	SellPricePo1     *float64   `gorm:"column:sell_price_po1" json:"sell_price_po1"`
	QtyPo3           *float64   `gorm:"column:qty_po3" json:"qty_po3"`
	QtyPo2           *float64   `gorm:"column:qty_po2" json:"qty_po2"`
	QtyPo1           *float64   `gorm:"column:qty_po1" json:"qty_po1"`
	VatValueFinal    *float64   `gorm:"column:vat_value_final" json:"vat_value_final"`
	DiscValueFinal   *float64   `gorm:"column:disc_value_final" json:"disc_value_final"`
	Vat              *float64   `gorm:"column:vat" json:"vat"`
}

func (SoDownloadPo) TableName() string {
	return "sls.order_detail"
}

type SoDownloadSo struct {
	OrderNo          *string    `gorm:"column:order_no" json:"order_no"`
	PoNo             *string    `gorm:"column:po_no" json:"po_no"`
	SoNo             string     `gorm:"column:so_no" json:"so_no"`
	RoDate           *time.Time `gorm:"column:ro_date" json:"ro_date"`
	InvoiceDate      *time.Time `gorm:"column:invoice_date" json:"invoice_date"`
	InvoiceNo        *string    `gorm:"column:invoice_no" json:"invoice_no"`
	OutletCode       *string    `gorm:"column:outlet_code" json:"outlet_code"`
	OutletName       *string    `gorm:"column:outlet_name" json:"outlet_name"`
	SalesmanCode     *string    `gorm:"column:salesman_code" json:"salesman_code"`
	EmployeeName     *string    `gorm:"column:employee_name" json:"employee_name"`
	SupplierCode     *string    `gorm:"column:supplier_code" json:"supplier_code"`
	SupplierName     *string    `gorm:"column:supplier_name" json:"supplier_name"`
	ProductCode      string     `gorm:"column:product_code" json:"product_code"`
	ProductName      string     `gorm:"column:product_name" json:"product_name"`
	UnitId3          *string    `gorm:"column:unit_id3" json:"unit_id3"`
	UnitId2          *string    `gorm:"column:unit_id2" json:"unit_id2"`
	UnitId1          *string    `gorm:"column:unit_id1" json:"unit_id1"`
	SellPriceSystem1 *float64   `gorm:"column:sell_price_system1" json:"sell_price_system1"`
	SellPriceSystem2 *float64   `gorm:"column:sell_price_system2" json:"sell_price_system2"`
	SellPriceSystem3 *float64   `gorm:"column:sell_price_system3" json:"sell_price_system3"`
	SellPrice3       *float64   `gorm:"column:sell_price3" json:"sell_price3"`
	SellPrice2       *float64   `gorm:"column:sell_price2" json:"sell_price2"`
	SellPrice1       *float64   `gorm:"column:sell_price1" json:"sell_price1"`
	Qty3             *float64   `gorm:"column:qty3" json:"qty3"`
	Qty2             *float64   `gorm:"column:qty2" json:"qty2"`
	Qty1             *float64   `gorm:"column:qty1" json:"qty1"`
	VatValueFinal    *float64   `gorm:"column:vat_value_final" json:"vat_value_final"`
	DiscValueFinal   *float64   `gorm:"column:disc_value_final" json:"disc_value_final"`
	Vat              *float64   `gorm:"column:vat" json:"vat"`
}

func (SoDownloadSo) TableName() string {
	return "sls.order_detail"
}

type SoDownloadFinal struct {
	OrderNo          *string    `gorm:"column:order_no" json:"order_no"`
	PoNo             *string    `gorm:"column:po_no" json:"po_no"`
	SoNo             string     `gorm:"column:so_no" json:"so_no"`
	RoDate           *time.Time `gorm:"column:ro_date" json:"ro_date"`
	InvoiceDate      *time.Time `gorm:"column:invoice_date" json:"invoice_date"`
	InvoiceNo        *string    `gorm:"column:invoice_no" json:"invoice_no"`
	OutletCode       *string    `gorm:"column:outlet_code" json:"outlet_code"`
	OutletName       *string    `gorm:"column:outlet_name" json:"outlet_name"`
	SalesmanCode     *string    `gorm:"column:salesman_code" json:"salesman_code"`
	EmployeeName     *string    `gorm:"column:employee_name" json:"employee_name"`
	SupplierCode     *string    `gorm:"column:supplier_code" json:"supplier_code"`
	SupplierName     *string    `gorm:"column:supplier_name" json:"supplier_name"`
	ProductCode      string     `gorm:"column:product_code" json:"product_code"`
	ProductName      string     `gorm:"column:product_name" json:"product_name"`
	UnitId3          *string    `gorm:"column:unit_id3" json:"unit_id3"`
	UnitId2          *string    `gorm:"column:unit_id2" json:"unit_id2"`
	UnitId1          *string    `gorm:"column:unit_id1" json:"unit_id1"`
	SellPriceSystem1 *float64   `gorm:"column:sell_price_system1" json:"sell_price_system1"`
	SellPriceSystem2 *float64   `gorm:"column:sell_price_system2" json:"sell_price_system2"`
	SellPriceSystem3 *float64   `gorm:"column:sell_price_system3" json:"sell_price_system3"`
	SellPriceFinal3  *float64   `gorm:"column:sell_price_final3" json:"sell_price_final3"`
	SellPriceFinal2  *float64   `gorm:"column:sell_price_final2" json:"sell_price_final2"`
	SellPriceFinal1  *float64   `gorm:"column:sell_price_final1" json:"sell_price_final1"`
	Qty3Final        *float64   `gorm:"column:qty3_final" json:"qty3_final"`
	Qty2Final        *float64   `gorm:"column:qty2_final" json:"qty2_final"`
	Qty1Final        *float64   `gorm:"column:qty1_final" json:"qty1_final"`
	VatValueFinal    *float64   `gorm:"column:vat_value_final" json:"vat_value_final"`
	DiscValueFinal   *float64   `gorm:"column:disc_value_final" json:"disc_value_final"`
	Vat              *float64   `gorm:"column:vat" json:"vat"`
}

func (SoDownloadFinal) TableName() string {
	return "sls.order_detail"
}

type SoDownloadQtySummary struct {
	OrderNo      *string    `gorm:"column:order_no" json:"order_no"`
	PoNo         *string    `gorm:"column:po_no" json:"po_no"`
	SoNo         string     `gorm:"column:so_no" json:"so_no"`
	RoDate       *time.Time `gorm:"column:ro_date" json:"ro_date"`
	InvoiceDate  *time.Time `gorm:"column:invoice_date" json:"invoice_date"`
	InvoiceNo    *string    `gorm:"column:invoice_no" json:"invoice_no"`
	OutletCode   *string    `gorm:"column:outlet_code" json:"outlet_code"`
	OutletName   *string    `gorm:"column:outlet_name" json:"outlet_name"`
	SalesmanCode *string    `gorm:"column:salesman_code" json:"salesman_code"`
	EmployeeName *string    `gorm:"column:employee_name" json:"employee_name"`
	SupplierCode *string    `gorm:"column:supplier_code" json:"supplier_code"`
	SupplierName *string    `gorm:"column:supplier_name" json:"supplier_name"`
	ProductCode  string     `gorm:"column:product_code" json:"product_code"`
	ProductName  string     `gorm:"column:product_name" json:"product_name"`
	UnitId3      *string    `gorm:"column:unit_id3" json:"unit_id3"`
	UnitId2      *string    `gorm:"column:unit_id2" json:"unit_id2"`
	UnitId1      *string    `gorm:"column:unit_id1" json:"unit_id1"`
	QtyPo3       *float64   `gorm:"column:qty_po3" json:"qty_po3"`
	QtyPo2       *float64   `gorm:"column:qty_po2" json:"qty_po2"`
	QtyPo1       *float64   `gorm:"column:qty_po1" json:"qty_po1"`
	Qty3         *float64   `gorm:"column:qty3" json:"qty3"`
	Qty2         *float64   `gorm:"column:qty2" json:"qty2"`
	Qty1         *float64   `gorm:"column:qty1" json:"qty1"`
	Qty3Final    *float64   `gorm:"column:qty3_final" json:"qty3_final"`
	Qty2Final    *float64   `gorm:"column:qty2_final" json:"qty2_final"`
	Qty1Final    *float64   `gorm:"column:qty1_final" json:"qty1_final"`
}

func (SoDownloadQtySummary) TableName() string {
	return "sls.order_detail"
}
