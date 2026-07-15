package model

import (
	"errors"
	"fmt"
	"time"
)

type StockReturnList struct {
	CustID         string     `gorm:"column:cust_id" json:"cust_id"`
	RefferenceNo   *string    `gorm:"column:refference_no" json:"refference_no"`
	ReturnNo       string     `gorm:"column:return_no;primaryKey" json:"return_no"`
	InvoiceNo      *string    `gorm:"column:invoice_no" json:"invoice_no"`
	InvoiceDate    *time.Time `gorm:"column:invoice_date" json:"invoice_date"`
	SalesmanID     *int64     `gorm:"column:salesman_id" json:"salesman_id"`
	SalesmanCode   *string    `gorm:"column:salesman_code" json:"salesman_code"`
	SalesmanName   *string    `gorm:"column:salesman_name" json:"salesman_name"`
	OutletID       *int64     `gorm:"column:outlet_id" json:"outlet_id"`
	OutletCode     *string    `gorm:"column:outlet_code" json:"outlet_code"`
	OutletName     *string    `gorm:"column:outlet_name" json:"outlet_name"`
	DataStatus     *int64     `gorm:"column:data_status" json:"data_status"`
	CreatedBy      *int64     `gorm:"column:created_by" json:"created_by"`
	CreatedByName  *string    `gorm:"column:created_by_name" json:"created_by_name"`
	CreatedAt      *time.Time `gorm:"column:created_at" json:"created_at"`
	ReviewedBy     *int64     `gorm:"column:reviewed_by" json:"reviewed_by"`
	ReviewedByName *string    `gorm:"column:reviewed_by_name" json:"reviewed_by_name"`
	ReviewedAt     *time.Time `gorm:"column:reviewed_at" json:"reviewed_at"`
	/*
		UpdatedBy     *int64     `gorm:"column:updated_by" json:"updated_by"`
		UpdatedAt     *time.Time `gorm:"column:updated_at" json:"updated_at"`
		UpdatedByName *string    `gorm:"column:updated_by_name" json:"updated_by_name"`
		SysDate       *time.Time      `gorm:"column:sys_date" json:"sys_date"`
		ReturnType    *int64          `gorm:"column:return_type" json:"return_type"`
		OutletTaxNo   *string         `gorm:"column:outlet_tax_no" json:"outlet_tax_no"`
		PoNo          *string         `gorm:"column:po_no" json:"po_no"`
		VehicleNo     *string         `gorm:"column:vehicle_no" json:"vehicle_no"`
		DeliveryDate  *time.Time      `gorm:"column:delivery_date" json:"delivery_date"`
		PayType       *int64          `gorm:"column:pay_type" json:"pay_type"`
		SumNo         *string         `gorm:"column:sum_no" json:"sum_no"`
		DataSource    *int64          `gorm:"column:data_source" json:"data_source"`
		MobileID      *int64          `gorm:"column:mobile_id" json:"mobile_id"`
		SubTotal      *float64        `gorm:"column:sub_total" json:"sub_total"`
		Disc          *float64        `gorm:"column:disc" json:"disc"`
		DiscValue     *float64        `gorm:"column:disc_value" json:"disc_value"`
		PromoValue    *float64        `gorm:"column:promo_value" json:"promo_value"`
		CashDiscValue *float64        `gorm:"column:cash_disc_value" json:"cash_disc_value"`
		TotDisc1      *float64        `gorm:"column:tot_disc1" json:"tot_disc_1"`
		TotDisc2      *float64        `gorm:"column:tot_disc2" json:"tot_disc_2"`
		Vat           *float64        `gorm:"column:vat" json:"vat"`
		VatValue      *float64        `gorm:"column:vat_value" json:"vat_value"`
		Total         *float64        `gorm:"column:total" json:"total"`
		IsDel         bool            `gorm:"column:is_del" json:"is_del"`
		DeletedBy     *int64          `gorm:"column:deleted_by" json:"deleted_by"`
		DeletedAt     *gorm.DeletedAt `gorm:"column:deleted_at" json:"deleted_at"`
	*/
}

func (StockReturnList) TableName() string {
	return "sls.return"
}

type StockReturnRead struct {
	CustID        string     `gorm:"column:cust_id" json:"cust_id"`
	RefferenceNo  *string    `gorm:"column:refference_no" json:"refference_no"`
	ReturnNo      string     `gorm:"column:return_no;primaryKey" json:"return_no"`
	ReturnDate    *time.Time `gorm:"column:return_date" json:"return_date"`
	InvoiceNo     *string    `gorm:"column:invoice_no" json:"invoice_no"`
	InvoiceDate   *time.Time `gorm:"column:invoice_date" json:"invoice_date"`
	SalesmanID    *int64     `gorm:"column:salesman_id" json:"salesman_id"`
	SalesmanCode  *string    `gorm:"column:salesman_code" json:"salesman_code"`
	SalesmanName  *string    `gorm:"column:salesman_name" json:"salesman_name"`
	OutletID      *int64     `gorm:"column:outlet_id" json:"outlet_id"`
	OutletCode    *string    `gorm:"column:outlet_code" json:"outlet_code"`
	OutletName    *string    `gorm:"column:outlet_name" json:"outlet_name"`
	TprCashValue  *float64   `gorm:"column:tpr_cash_value" json:"tpr_cash_value"`
	TprItemValue  *float64   `gorm:"column:tpr_item_value" json:"tpr_item_value"`
	Discount      *float64   `gorm:"column:discount" json:"discount"`
	DiscountValue *float64   `gorm:"column:disc_value" json:"disc_value"`
	Vat           *float64   `gorm:"column:vat" json:"vat"`
	VatValue      *float64   `gorm:"column:vat_value" json:"vat_value"`
	SubTotal      *float64   `gorm:"column:sub_total" json:"sub_total"`
	Total         *float64   `gorm:"column:total" json:"total"`
	DataStatus    *int64     `gorm:"column:data_status" json:"data_status"`
}

func (StockReturnRead) TableName() string {
	return "sls.return"
}

type MapStockReturn map[string]StockReturnRead

func (m MapStockReturn) Set(id string, stockReturn StockReturnRead) {
	m[id] = stockReturn
}

func (m MapStockReturn) GetByID(id string) (stockReturn StockReturnRead, err error) {
	val, ok := m[id]

	// If the key exists
	if !ok {
		return stockReturn, errors.New(fmt.Sprintf("stockReturn no %v Not Found", id))
	}

	return val, nil
}
