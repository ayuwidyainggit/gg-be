package model

import (
	"strconv"
	"time"

	"gorm.io/gorm"
)

type So struct {
	CustID        string          `gorm:"column:cust_id" json:"cust_id"`
	SoNo          string          `gorm:"column:so_no" json:"so_no"`
	SoDate        *time.Time      `gorm:"column:so_date" json:"so_date"`
	SysDate       *time.Time      `gorm:"column:sys_date" json:"sys_date"`
	OutletID      *int64          `gorm:"column:outlet_id" json:"outlet_id"`
	OutletTaxNo   *string         `gorm:"column:outlet_tax_no" json:"outlet_tax_no"`
	PoNo          *string         `gorm:"column:po_no" json:"po_no"`
	VehicleNo     *string         `gorm:"column:vehicle_no" json:"vehicle_no"`
	InvoiceNo     *string         `gorm:"column:invoice_no" json:"invoice_no"`
	InvoiceDate   *time.Time      `gorm:"column:invoice_date" json:"invoice_date"`
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
	TotalInv      float64         `gorm:"column:total_inv" json:"total_inv"`
	TotalInvIn    *float64        `gorm:"column:total_inv_in" json:"total_inv_in"`
	RoundDiff     *float64        `gorm:"column:round_diff" json:"round_diff"`
	DataStatus    *int64          `gorm:"column:data_status" json:"data_status"`
	CreatedBy     *int64          `gorm:"column:created_by" json:"created_by"`
	CreatedAt     time.Time       `gorm:"column:created_at" json:"created_at"`
	UpdatedBy     *int64          `gorm:"column:updated_by" json:"updated_by"`
	UpdatedAt     time.Time       `gorm:"column:updated_at" json:"updated_at"`
	IsDel         bool            `gorm:"column:is_del" json:"is_del"`
	DeletedBy     *int64          `gorm:"column:deleted_by" json:"deleted_by"`
	DeletedAt     *gorm.DeletedAt `gorm:"column:deleted_at" json:"deleted_at"`
	DueDate       *time.Time      `gorm:"due_date" json:"due_date"`
}

func (So) TableName() string {
	return "sls.so"
}
func (m *So) BeforeCreate(trx *gorm.DB) (err error) {
	intTmpsStr := time.Now().UnixNano() / int64(time.Millisecond)
	m.SoNo = strconv.Itoa(int(intTmpsStr))
	m.CreatedAt = time.Now()
	m.UpdatedAt = time.Now()
	m.UpdatedBy = m.CreatedBy
	return nil
}

type SoList struct {
	CustID        string          `gorm:"column:cust_id" json:"cust_id"`
	SoNo          string          `gorm:"column:so_no" json:"so_no"`
	SoDate        *time.Time      `gorm:"column:so_date" json:"so_date"`
	SysDate       *time.Time      `gorm:"column:sys_date" json:"sys_date"`
	OutletID      *int64          `gorm:"column:outlet_id" json:"outlet_id"`
	OutletCode    *string         `gorm:"column:outlet_code" json:"outlet_code"`
	OutletName    *string         `gorm:"column:outlet_name" json:"outlet_name"`
	OutletTaxNo   *string         `gorm:"column:outlet_tax_no" json:"outlet_tax_no"`
	PoNo          *string         `gorm:"column:po_no" json:"po_no"`
	VehicleNo     *string         `gorm:"column:vehicle_no" json:"vehicle_no"`
	InvoiceNo     *string         `gorm:"column:invoice_no" json:"invoice_no"`
	InvoiceDate   *time.Time      `gorm:"column:invoice_date" json:"invoice_date"`
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
	TotalInv      float64         `gorm:"column:total_inv" json:"total_inv"`
	TotalInvIn    *float64        `gorm:"column:total_inv_in" json:"total_inv_in"`
	RoundDiff     *float64        `gorm:"column:round_diff" json:"round_diff"`
	DataStatus    *int64          `gorm:"column:data_status" json:"data_status"`
	CreatedBy     *int64          `gorm:"column:created_by" json:"created_by"`
	CreatedAt     time.Time       `gorm:"column:created_at" json:"created_at"`
	UpdatedBy     *int64          `gorm:"column:updated_by" json:"updated_by"`
	UpdatedByName *string         `gorm:"column:updated_by_name" json:"updated_by_name"`
	UpdatedAt     time.Time       `gorm:"column:updated_at" json:"updated_at"`
	IsDel         bool            `gorm:"column:is_del" json:"is_del"`
	DeletedBy     *int64          `gorm:"column:deleted_by" json:"deleted_by"`
	DeletedAt     *gorm.DeletedAt `gorm:"column:deleted_at" json:"deleted_at"`
	DueDate       *time.Time      `gorm:"due_date" json:"due_date"`
}

func (SoList) TableName() string {
	return "sls.so"
}
