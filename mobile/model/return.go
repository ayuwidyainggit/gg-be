package model

import (
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"
)

/*
	type ReturnOld struct {
		CustID        string          `gorm:"column:cust_id" json:"cust_id"`
		ReturnNo      string          `gorm:"column:return_no;primaryKey" json:"return_no"`
		ReturnDate    *time.Time      `gorm:"column:return_date" json:"return_date"`
		SysDate       *time.Time      `gorm:"column:sys_date" json:"sys_date"`
		ReturnType    *int64          `gorm:"column:return_type" json:"return_type"`
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
		Total         *float64        `gorm:"column:total" json:"total"`
		DataStatus    *int64          `gorm:"column:data_status" json:"data_status"`
		CreatedBy     *int64          `gorm:"column:created_by" json:"created_by"`
		CreatedAt     time.Time       `gorm:"column:created_at" json:"created_at"`
		UpdatedBy     *int64          `gorm:"column:updated_by" json:"updated_by"`
		UpdatedAt     time.Time       `gorm:"column:updated_at" json:"updated_at"`
		IsDel         bool            `gorm:"column:is_del" json:"is_del"`
		DeletedBy     *int64          `gorm:"column:deleted_by" json:"deleted_by"`
		DeletedAt     *gorm.DeletedAt `gorm:"column:deleted_at" json:"deleted_at"`
	}
*/
type ReturnList struct {
	CustID          string     `gorm:"column:cust_id" json:"cust_id"`
	RefferenceNo    *string    `gorm:"column:refference_no" json:"refference_no"`
	ReturnNo        string     `gorm:"column:return_no;primaryKey" json:"return_no"`
	ReturnDate      *time.Time `gorm:"column:return_date" json:"return_date"`
	InvoiceNo       *string    `gorm:"column:invoice_no" json:"invoice_no"`
	InvoiceDate     *time.Time `gorm:"column:invoice_date" json:"invoice_date"`
	SalesmanID      *int64     `gorm:"column:salesman_id" json:"salesman_id"`
	SalesmanCode    *string    `gorm:"column:salesman_code" json:"salesman_code"`
	SalesmanName    *string    `gorm:"column:salesman_name" json:"salesman_name"`
	OutletID        *int64     `gorm:"column:outlet_id" json:"outlet_id"`
	OutletCode      *string    `gorm:"column:outlet_code" json:"outlet_code"`
	OutletName      *string    `gorm:"column:outlet_name" json:"outlet_name"`
	OutletAddress   *string    `gorm:"column:outlet_address" json:"outlet_address"`
	OutletLatitude  *string    `gorm:"column:outlet_latitude" json:"outlet_latitude"`
	OutletLongitude *string    `gorm:"column:outlet_longitude" json:"outlet_longitude"`
	DataStatus      *int64     `gorm:"column:data_status" json:"data_status"`
	CreatedBy       *int64     `gorm:"column:created_by" json:"created_by"`
	CreatedByName   *string    `gorm:"column:created_by_name" json:"created_by_name"`
	CreatedAt       *time.Time `gorm:"column:created_at" json:"created_at"`
	ReviewedBy      *int64     `gorm:"column:reviewed_by" json:"reviewed_by"`
	ReviewedByName  *string    `gorm:"column:reviewed_by_name" json:"reviewed_by_name"`
	ReviewedAt      *time.Time `gorm:"column:reviewed_at" json:"reviewed_at"`
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

func (ReturnList) TableName() string {
	return "sls.return"
}

type ReturnRead struct {
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
	DiscountValue *float64   `gorm:"column:discount_value" json:"discount_value"`
	Vat           *float64   `gorm:"column:vat" json:"vat"`
	VatValue      *float64   `gorm:"column:vat_value" json:"vat_value"`
	SubTotal      *float64   `gorm:"column:sub_total" json:"sub_total"`
	Total         *float64   `gorm:"column:total" json:"total"`
	DataStatus    *int64     `gorm:"column:data_status" json:"data_status"`
}

func (ReturnRead) TableName() string {
	return "sls.return"
}

/*
	type ReturnListOld struct {
		CustID        string          `gorm:"column:cust_id" json:"cust_id"`
		ReturnNo      string          `gorm:"column:return_no;primaryKey" json:"return_no"`
		ReturnDate    *time.Time      `gorm:"column:return_date" json:"return_date"`
		SysDate       *time.Time      `gorm:"column:sys_date" json:"sys_date"`
		ReturnType    *int64          `gorm:"column:return_type" json:"return_type"`
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
		Total         *float64        `gorm:"column:total" json:"total"`
		DataStatus    *int64          `gorm:"column:data_status" json:"data_status"`
		CreatedBy     *int64          `gorm:"column:created_by" json:"created_by"`
		CreatedAt     time.Time       `gorm:"column:created_at" json:"created_at"`
		UpdatedBy     *int64          `gorm:"column:updated_by" json:"updated_by"`
		UpdatedAt     *time.Time      `gorm:"column:updated_at" json:"updated_at"`
		UpdatedByName *string         `gorm:"column:updated_by_name" json:"updated_by_name"`
		IsDel         bool            `gorm:"column:is_del" json:"is_del"`
		DeletedBy     *int64          `gorm:"column:deleted_by" json:"deleted_by"`
		DeletedAt     *gorm.DeletedAt `gorm:"column:deleted_at" json:"deleted_at"`
	}

	func (ReturnListOld) TableName() string {
		return "sls.return"
	}
*/
type SalesmansFilter struct {
	CustId        string          `gorm:"column:cust_id" json:"cust_id"`
	SalesmanId    int             `gorm:"column:salesman_id" json:"salesman_id"`
	SalesmanCode  string          `gorm:"column:salesman_code" json:"salesman_code"`
	SalesmanName  string          `gorm:"column:salesman_name" json:"salesman_name"`
	IsActive      bool            `gorm:"column:is_active" json:"is_active"`
	IsDel         bool            `gorm:"column:is_del" json:"is_del"`
	CreatedBy     *int64          `gorm:"column:created_by" json:"created_by"`
	CreatedAt     *time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedBy     *int64          `gorm:"column:updated_by" json:"updated_by"`
	UpdatedAt     *time.Time      `gorm:"column:updated_at" json:"updated_at"`
	UpdatedByName *string         `gorm:"column:updated_by_name" db:"updated_by_name"`
	DeletedBy     *int64          `gorm:"column:deleted_by" json:"deleted_by"`
	DeletedAt     *gorm.DeletedAt `gorm:"column:deleted_at" json:"deleted_at"`
}

func (SalesmansFilter) TableName() string {
	return "sls.return"
}

type OutletsFilter struct {
	CustId        string          `gorm:"column:cust_id" json:"cust_id"`
	OutletId      int             `gorm:"column:outlet_id" json:"outlet_id"`
	OutletCode    string          `gorm:"column:outlet_code" json:"outlet_code"`
	OutletName    string          `gorm:"column:outlet_name" json:"outlet_name"`
	IsActive      bool            `gorm:"column:is_active" json:"is_active"`
	IsDel         bool            `gorm:"column:is_del" json:"is_del"`
	CreatedBy     *int64          `gorm:"column:created_by" json:"created_by"`
	CreatedAt     *time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedBy     *int64          `gorm:"column:updated_by" json:"updated_by"`
	UpdatedAt     *time.Time      `gorm:"column:updated_at" json:"updated_at"`
	UpdatedByName *string         `gorm:"column:updated_by_name" db:"updated_by_name"`
	DeletedBy     *int64          `gorm:"column:deleted_by" json:"deleted_by"`
	DeletedAt     *gorm.DeletedAt `gorm:"column:deleted_at" json:"deleted_at"`
}

func (OutletsFilter) TableName() string {
	return "sls.return"
}

type ReturnStatusesFilter struct {
	ReturnStatus int `gorm:"column:return_status" json:"return_status"`
}

func (ReturnStatusesFilter) TableName() string {
	return "sls.return"
}

type SalesmansFilterCreate struct {
	CustId        string          `gorm:"column:cust_id" json:"cust_id"`
	SalesmanId    int             `gorm:"column:salesman_id" json:"salesman_id"`
	SalesmanCode  string          `gorm:"column:salesman_code" json:"salesman_code"`
	SalesmanName  string          `gorm:"column:salesman_name" json:"salesman_name"`
	IsActive      bool            `gorm:"column:is_active" json:"is_active"`
	IsDel         bool            `gorm:"column:is_del" json:"is_del"`
	CreatedBy     *int64          `gorm:"column:created_by" json:"created_by"`
	CreatedAt     *time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedBy     *int64          `gorm:"column:updated_by" json:"updated_by"`
	UpdatedAt     *time.Time      `gorm:"column:updated_at" json:"updated_at"`
	UpdatedByName *string         `gorm:"column:updated_by_name" db:"updated_by_name"`
	DeletedBy     *int64          `gorm:"column:deleted_by" json:"deleted_by"`
	DeletedAt     *gorm.DeletedAt `gorm:"column:deleted_at" json:"deleted_at"`
}

func (SalesmansFilterCreate) TableName() string {
	return "sls.order"
}

type OutletsFilterCreate struct {
	CustId        string          `gorm:"column:cust_id" json:"cust_id"`
	OutletId      int             `gorm:"column:outlet_id" json:"outlet_id"`
	OutletCode    string          `gorm:"column:outlet_code" json:"outlet_code"`
	OutletName    string          `gorm:"column:outlet_name" json:"outlet_name"`
	IsActive      bool            `gorm:"column:is_active" json:"is_active"`
	IsDel         bool            `gorm:"column:is_del" json:"is_del"`
	CreatedBy     *int64          `gorm:"column:created_by" json:"created_by"`
	CreatedAt     *time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedBy     *int64          `gorm:"column:updated_by" json:"updated_by"`
	UpdatedAt     *time.Time      `gorm:"column:updated_at" json:"updated_at"`
	UpdatedByName *string         `gorm:"column:updated_by_name" db:"updated_by_name"`
	DeletedBy     *int64          `gorm:"column:deleted_by" json:"deleted_by"`
	DeletedAt     *gorm.DeletedAt `gorm:"column:deleted_at" json:"deleted_at"`
}

func (OutletsFilterCreate) TableName() string {
	return "sls.order"
}

type ProductList struct {
	OrderDetailID *int64     `gorm:"order_detail_id" json:"order_detail_id"`
	InvoiceNo     *string    `gorm:"invoice_no" json:"invoice_no"`
	InvoiceDate   *time.Time `gorm:"invoice_date" json:"invoice_date"`
	ProductId     int64      `gorm:"product_id" json:"product_id"`
	ProductCode   *string    `gorm:"product_code" json:"product_code"`
	ProductName   *string    `gorm:"product_name" json:"product_name"`
	InvoiceQty1   float64    `gorm:"invoice_qty1" json:"invoice_qty1"`
	InvoiceQty2   float64    `gorm:"invoice_qty2" json:"invoice_qty2"`
	InvoiceQty3   float64    `gorm:"invoice_qty3" json:"invoice_qty3"`
	RemainingQty1 float64    `gorm:"remaining_qty1" json:"remaining_qty1"`
	RemainingQty2 float64    `gorm:"remaining_qty2" json:"remaining_qty2"`
	RemainingQty3 float64    `gorm:"remaining_qty3" json:"remaining_qty3"`
	SellPrice1    float64    `gorm:"sell_price1" json:"sell_price1"`
	SellPrice2    float64    `gorm:"sell_price2" json:"sell_price2"`
	SellPrice3    float64    `gorm:"sell_price3" json:"sell_price3"`
	UnitId1       string     `gorm:"unit_id1" json:"unit_id1"`
	UnitId2       string     `gorm:"unit_id2" json:"unit_id2"`
	UnitId3       string     `gorm:"unit_id3" json:"unit_id3"`
	UnitName1     *string    `gorm:"unit_name1" json:"unit_name1"`
	UnitName2     *string    `gorm:"unit_name2" json:"unit_name2"`
	UnitName3     *string    `gorm:"unit_name3" json:"unit_name3"`
	ConvUnit2     *float64   `gorm:"conv_unit2" json:"conv_unit2"`
	ConvUnit3     *float64   `gorm:"conv_unit3" json:"conv_unit3"`
	Vat           *float64   `gorm:"vat" json:"vat"`
	SubTotal      float64    `gorm:"sub_total" json:"sub_total"`
	Total         float64    `gorm:"total" json:"total"`
}

func (ProductList) TableName() string {
	return "sls.order_detail"
}

type ReturnReasonLookup struct {
	CustId           string          `gorm:"column:cust_id" json:"cust_id"`
	ReturnReasonId   int             `gorm:"column:return_reason_id" json:"return_reason_id"`
	ReturnReasonName string          `gorm:"column:return_reason_name" json:"return_reason_name"`
	IsActive         bool            `gorm:"column:is_active" json:"is_active"`
	IsDel            bool            `gorm:"column:is_del" json:"is_del"`
	CreatedBy        *int64          `gorm:"column:created_by" json:"created_by"`
	CreatedAt        *time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedBy        *int64          `gorm:"column:updated_by" json:"updated_by"`
	UpdatedAt        *time.Time      `gorm:"column:updated_at" json:"updated_at"`
	UpdatedByName    *string         `gorm:"column:updated_by_name" db:"updated_by_name"`
	DeletedBy        *int64          `gorm:"column:deleted_by" json:"deleted_by"`
	DeletedAt        *gorm.DeletedAt `gorm:"column:deleted_at" json:"deleted_at"`
}

func (ReturnReasonLookup) TableName() string {
	return "mst.m_return_reason"
}

type WarehouseLookup struct {
	CustId        string          `gorm:"column:cust_id" json:"cust_id"`
	WhId          int             `gorm:"column:wh_id" json:"wh_id"`
	WhCode        string          `gorm:"column:wh_code" json:"wh_code"`
	WhName        string          `gorm:"column:wh_name" json:"wh_name"`
	IsActive      bool            `gorm:"column:is_active" json:"is_active"`
	IsDel         bool            `gorm:"column:is_del" json:"is_del"`
	CreatedBy     *int64          `gorm:"column:created_by" json:"created_by"`
	CreatedAt     *time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedBy     *int64          `gorm:"column:updated_by" json:"updated_by"`
	UpdatedAt     *time.Time      `gorm:"column:updated_at" json:"updated_at"`
	UpdatedByName *string         `gorm:"column:updated_by_name" db:"updated_by_name"`
	DeletedBy     *int64          `gorm:"column:deleted_by" json:"deleted_by"`
	DeletedAt     *gorm.DeletedAt `gorm:"column:deleted_at" json:"deleted_at"`
}

func (WarehouseLookup) TableName() string {
	return "mst.m_warehouse"
}

type ProductsLookup struct {
	ProductId   *int64   `gorm:"product_id" json:"product_id"`
	ProductCode *string  `gorm:"product_code" json:"product_code"`
	ProductName *string  `gorm:"product_name" json:"product_name"`
	SellPrice1  *float64 `gorm:"sell_price1" json:"sell_price1"`
	SellPrice2  *float64 `gorm:"sell_price2" json:"sell_price2"`
	SellPrice3  *float64 `gorm:"sell_price3" json:"sell_price3"`
	UnitId1     *string  `gorm:"unit_id1" json:"unit_id1"`
	UnitId2     *string  `gorm:"unit_id2" json:"unit_id2"`
	UnitId3     *string  `gorm:"unit_id3" json:"unit_id3"`
	UnitName1   *string  `gorm:"unit_name1" json:"unit_name1"`
	UnitName2   *string  `gorm:"unit_name2" json:"unit_name2"`
	UnitName3   *string  `gorm:"unit_name3" json:"unit_name3"`
	ConvUnit2   *float64 `gorm:"conv_unit2" json:"conv_unit2"`
	ConvUnit3   *float64 `gorm:"conv_unit3" json:"conv_unit3"`
	Vat         *float64 `gorm:"vat" json:"vat"`
}

func (ProductsLookup) TableName() string {
	return "mst.m_product"
}

type Return struct {
	CustID       string     `gorm:"column:cust_id" json:"cust_id"`
	RefferenceNo string     `gorm:"column:refference_no" json:"refference_no"`
	ReturnNo     string     `gorm:"column:return_no;primaryKey" json:"return_no"`
	ReturnDate   time.Time  `gorm:"column:return_date" json:"return_date"`
	SalesmanID   int64      `gorm:"column:salesman_id" json:"salesman_id"`
	OutletID     int64      `gorm:"column:outlet_id" json:"outlet_id"`
	InvoiceNo    *string    `gorm:"column:invoice_no" json:"invoice_no"`
	InvoiceDate  *time.Time `gorm:"column:invoice_date" json:"invoice_date"`
	Discount     float64    `gorm:"column:discount" json:"discount"`
	// DiscountValue float64        `gorm:"column:discount_value" json:"discount_value"`
	Vat        float64         `gorm:"column:vat" json:"vat"`
	VatValue   float64         `gorm:"column:vat_value" json:"vat_value"`
	SubTotal   float64         `gorm:"column:sub_total" json:"sub_total"`
	Total      float64         `gorm:"column:total" json:"total"`
	DataStatus int64           `gorm:"column:data_status" json:"data_status"`
	CreatedBy  int64           `gorm:"column:created_by" json:"created_by"`
	CreatedAt  time.Time       `gorm:"column:created_at" json:"created_at"`
	UpdatedBy  int64           `gorm:"column:updated_by" json:"updated_by"`
	UpdatedAt  time.Time       `gorm:"column:updated_at" json:"updated_at"`
	IsReviewed bool            `gorm:"column:is_reviewed" json:"is_reviewed"`
	ReviewedBy *int64          `gorm:"column:reviewed_by" json:"reviewed_by"`
	ReviewedAt *time.Time      `gorm:"column:reviewed_at" json:"reviewed_at"`
	IsDel      bool            `gorm:"column:is_del" json:"is_del"`
	DeletedBy  *int64          `gorm:"column:deleted_by" json:"deleted_by"`
	DeletedAt  *gorm.DeletedAt `gorm:"column:deleted_at" json:"deleted_at"`
}

func (Return) TableName() string {
	return "sls.return"
}

type ReturnNo struct {
	ReturnNo string `gorm:"column:get_no_fn"`
}

func (m *Return) BeforeCreate(trx *gorm.DB) (err error) {
	var returnNo ReturnNo
	trCode := "SR"
	returnDateStr := m.ReturnDate.Format("2006-01-02")
	returnDateSubtr := returnDateStr[2:4] + returnDateStr[5:7] + returnDateStr[8:10]

	queryStr := fmt.Sprintf(`SELECT
	TRIM(to_char(COALESCE(MAX(TO_NUMBER(SUBSTR(return_no,9,4),'9999')),0)+1, '0000')) AS get_no_fn
	FROM sls.return
	WHERE substr(return_no,3,6) = '%v' AND cust_id = '%v'`, returnDateSubtr, strings.ToUpper(m.CustID))
	err = trx.Raw(queryStr).Scan(&returnNo).Error
	if err != nil {
		return err
	}

	m.ReturnNo = trCode + returnDateSubtr + returnNo.ReturnNo
	m.CreatedAt = time.Now()
	m.UpdatedAt = time.Now()
	m.UpdatedBy = m.CreatedBy

	return nil
}

func (m *Return) BeforeUpdate(trx *gorm.DB) (err error) {
	now := time.Now().UTC()
	m.UpdatedAt = now

	return nil
}

type SalesmanRead struct {
	CustId        string          `gorm:"column:cust_id" json:"cust_id"`
	SalesmanId    int             `gorm:"column:salesman_id" json:"salesman_id"`
	SalesmanCode  string          `gorm:"column:salesman_code" json:"salesman_code"`
	SalesmanName  string          `gorm:"column:salesman_name" json:"salesman_name"`
	WhId          int             `gorm:"column:wh_id" json:"wh_id"`
	SalesTeamId   int             `gorm:"column:sales_team_id" json:"sales_team_id"`
	IsActive      bool            `gorm:"column:is_active" json:"is_active"`
	IsDel         bool            `gorm:"column:is_del" json:"is_del"`
	CreatedBy     *int64          `gorm:"column:created_by" json:"created_by"`
	CreatedAt     *time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedBy     *int64          `gorm:"column:updated_by" json:"updated_by"`
	UpdatedAt     *time.Time      `gorm:"column:updated_at" json:"updated_at"`
	UpdatedByName *string         `gorm:"column:updated_by_name" db:"updated_by_name"`
	DeletedBy     *int64          `gorm:"column:deleted_by" json:"deleted_by"`
	DeletedAt     *gorm.DeletedAt `gorm:"column:deleted_at" json:"deleted_at"`
}

func (SalesmanRead) TableName() string {
	return "mst.m_salesman"
}
