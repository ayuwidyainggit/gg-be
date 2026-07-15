package model

import (
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"
)

type GrBranch struct {
	CustID               string         `gorm:"column:cust_id" json:"cust_id"`
	GrBranchNo           string         `gorm:"column:gr_branch_no" json:"gr_branch_no"`
	GrBranchDate         *time.Time     `gorm:"column:gr_branch_date" json:"gr_branch_date"`
	DeliveryDate         *time.Time     `gorm:"column:delivery_date" json:"delivery_date"`
	DeliveryNo           *string        `gorm:"column:delivery_no" json:"delivery_no"`
	DeliveryFee          *float64       `gorm:"column:delivery_fee" json:"delivery_fee"`
	InvoiceNo            *string        `gorm:"column:invoice_no" json:"invoice_no"`
	InvoiceDate          *time.Time     `gorm:"column:invoice_date" json:"invoice_date"`
	InvoiceNoBranch      *string        `gorm:"column:invoice_no_branch" json:"invoice_no_branch"`
	InvoiceDateBranch    *time.Time     `gorm:"column:invoice_date_branch" json:"invoice_date_branch"`
	InvoiceDueDateBranch *time.Time     `gorm:"column:invoice_due_date_branch" json:"invoice_due_date_branch"`
	VehicleNo            *string        `gorm:"column:vehicle_no" json:"vehicle_no"`
	PoNo                 *string        `gorm:"column:po_no" json:"po_no"`
	SoNo                 *string        `gorm:"column:so_no" json:"so_no"`
	SupID                *int64         `gorm:"column:sup_id" json:"sup_id"`
	WhID                 *int64         `gorm:"column:wh_id" json:"wh_id"`
	Notes                *string        `gorm:"column:notes" json:"notes"`
	DataStatus           *int           `gorm:"column:data_status" json:"data_status"`
	SubTotal             *float64       `gorm:"column:sub_total" json:"sub_total"`
	VatValue             *float64       `gorm:"column:vat_value" json:"vat_value"`
	Total                *float64       `gorm:"column:total" json:"total"`
	CreatedBy            *int64         `gorm:"column:created_by" json:"created_by"`
	CreatedAt            time.Time      `gorm:"column:created_at" json:"created_at,omitempty"`
	UpdatedBy            *int64         `gorm:"column:updated_by" json:"updated_by"`
	UpdatedAt            time.Time      `gorm:"column:updated_at" json:"updated_at,omitempty"`
	IsDel                bool           `gorm:"column:is_del" json:"is_del"`
	DeletedBy            *int64         `gorm:"column:deleted_by" json:"deleted_by"`
	DeletedAt            gorm.DeletedAt `gorm:"column:deleted_at" json:"deleted_at"`
	// IsClosed     bool           `gorm:"column:is_closed" json:"is_closed"`
	// ClosedBy     *int64         `gorm:"column:closed_by" json:"closed_by"`
	// ClosedAt     time.Time      `gorm:"column:closed_at" json:"closed_at"`
}

func (GrBranch) TableName() string {
	return "inv.gr_branch"
}

type GrBranchPrint struct {
	IsPrint    bool      `gorm:"column:is_print" json:"is_print"`
	PrintedBy  int64     `gorm:"column:printed_by" json:"printed_by"`
	PrintedAt  time.Time `gorm:"column:printed_at" json:"printed_at"`
	DataStatus int       `gorm:"column:data_status" json:"data_status"`
}

func (GrBranchPrint) TableName() string {
	return "inv.gr_branch"
}

func (m *GrBranch) BeforeCreate(trx *gorm.DB) (err error) {
	var grBranchNo GrBranchNo
	trCode := "GRB"
	grBranchDateStr := m.GrBranchDate.Format("2006-01-02")
	grBranchDateSubtr := grBranchDateStr[2:4] + grBranchDateStr[5:7] + grBranchDateStr[8:10]
	// log.Println("grBranchDateStr:", grBranchDateStr)
	// log.Println("grBranchDateSubtr:", grBranchDateSubtr)

	queryStr := fmt.Sprintf(`SELECT 
	TRIM(to_char(COALESCE(MAX(TO_NUMBER(SUBSTR(gr_branch_no,10,3),'999')),0)+1, '000')) AS get_no_fn 
	FROM inv.gr_branch
	WHERE substr(gr_branch_no,4,6) = '%v' AND cust_id = '%v'`, grBranchDateSubtr, strings.ToUpper(m.CustID))

	err = trx.Raw(queryStr).Scan(&grBranchNo).Error
	if err != nil {
		return err
	}

	// log.Println("grBranchNo:", grBranchNo.GrBranchNo)

	m.GrBranchNo = trCode + grBranchDateSubtr + grBranchNo.GrBranchNo
	// log.Println("m.GrBranchNo:", m.GrBranchNo)

	if *m.DataStatus == 2 {
		var invoiceNoBranch InvoiceNoGrBranch
		trInvoiceCode := "INVB"
		invoiceDateBranchStr := time.Now().Format("2006-01-02")
		invoiceDateBranchSubtr := invoiceDateBranchStr[2:4] + invoiceDateBranchStr[5:7] + invoiceDateBranchStr[8:10]
		// log.Println("grBranchDateStr:", grBranchDateStr)
		// log.Println("grBranchDateSubtr:", grBranchDateSubtr)

		queryInvoiceStr := fmt.Sprintf(`SELECT 
		TRIM(to_char(COALESCE(MAX(TO_NUMBER(SUBSTR(invoice_no_branch,11,3),'999')),0)+1, '000')) AS get_no_fn 
		FROM inv.gr_branch
		WHERE substr(invoice_no_branch,5,6) = '%v' AND cust_id = '%v'`, invoiceDateBranchSubtr, strings.ToUpper(m.CustID))

		err = trx.Raw(queryInvoiceStr).Scan(&invoiceNoBranch).Error
		if err != nil {
			return err
		}

		invoiceNoGrBranch := trInvoiceCode + invoiceDateBranchSubtr + invoiceNoBranch.InvoiceNoGrBranch
		m.InvoiceNoBranch = &invoiceNoGrBranch

		invoiceDateBranch := time.Now()
		m.InvoiceDateBranch = &invoiceDateBranch

		var termOfPaySupplierGrBranch TermOfPaySupplierGrBranch

		querySupplierStr := fmt.Sprintf(`SELECT 
		pay_term 
		FROM mst.m_supplier
		WHERE sup_id = '%v'`, *m.SupID)

		err = trx.Raw(querySupplierStr).Scan(&termOfPaySupplierGrBranch).Error
		if err != nil {
			return err
		}

		invoiceDueDateBranch := time.Now().AddDate(0, 0, termOfPaySupplierGrBranch.TermOfPay)
		m.InvoiceDueDateBranch = &invoiceDueDateBranch
	}

	m.GrBranchNo = trCode + grBranchDateSubtr + grBranchNo.GrBranchNo
	// log.Println("m.GrBranchNo:", m.GrBranchNo)

	m.CreatedAt = time.Now()
	m.UpdatedAt = time.Now()
	m.UpdatedBy = m.CreatedBy
	return nil
}

func (m *GrBranch) BeforeUpdate(trx *gorm.DB) (err error) {
	// var invoiceNoGrBranch string
	// var grDate time.Time
	if *m.DataStatus == 2 {
		var invoiceNoBranch InvoiceNoGrBranch
		trInvoiceCode := "INVB"
		invoiceDateBranchStr := time.Now().Format("2006-01-02")
		invoiceDateBranchSubtr := invoiceDateBranchStr[2:4] + invoiceDateBranchStr[5:7] + invoiceDateBranchStr[8:10]
		// log.Println("grBranchDateStr:", grBranchDateStr)
		// log.Println("grBranchDateSubtr:", grBranchDateSubtr)

		queryInvoiceStr := fmt.Sprintf(`SELECT 
		TRIM(to_char(COALESCE(MAX(TO_NUMBER(SUBSTR(invoice_no_branch,11,3),'999')),0)+1, '000')) AS get_no_fn 
		FROM inv.gr_branch
		WHERE substr(invoice_no_branch,5,6) = '%v' AND cust_id = '%v'`, invoiceDateBranchSubtr, strings.ToUpper(m.CustID))

		err = trx.Raw(queryInvoiceStr).Scan(&invoiceNoBranch).Error
		if err != nil {
			return err
		}

		invoiceNoGrBranch := trInvoiceCode + invoiceDateBranchSubtr + invoiceNoBranch.InvoiceNoGrBranch
		m.InvoiceNoBranch = &invoiceNoGrBranch

		grDate := time.Now()
		m.InvoiceDateBranch = &grDate

		// log.Println("Invoice No : ", invoiceNoGrBranch)
		// log.Println("Invoice Date : ", grDate.Format("2006-01-02"))
		var termOfPaySupplierGrBranch TermOfPaySupplierGrBranch

		querySupplierStr := fmt.Sprintf(`SELECT 
		pay_term 
		FROM mst.m_supplier
		WHERE sup_id = '%v'`, *m.SupID)

		err = trx.Raw(querySupplierStr).Scan(&termOfPaySupplierGrBranch).Error
		if err != nil {
			return err
		}

		invoiceDueDateBranch := time.Now().AddDate(0, 0, int(termOfPaySupplierGrBranch.TermOfPay))
		m.InvoiceDueDateBranch = &invoiceDueDateBranch
	}

	now := time.Now()
	m.UpdatedAt = now
	// deliveryFee := float64(100000)
	// m.DeliveryFee = &deliveryFee
	// if trx.Statement.Changed("DataStatus") {
	// 	trx.Statement.SetColumn("InvoiceNoBranch", invoiceNoGrBranch)
	// 	trx.Statement.SetColumn("InvoiceDateBranch", grDate.Format("2006-01-02"))
	// }

	// log.Println("Invoice No in Model : ", *m.InvoiceNoBranch)
	// log.Println("Invoice Date in Model : ", *m.InvoiceDateBranch)
	// log.Println("Delivery Fee in Model : ", *m.DeliveryFee)

	return nil
}

type GrBranchNo struct {
	GrBranchNo string `gorm:"column:get_no_fn"`
}

type GrBranchList struct {
	CustID            string     `gorm:"column:cust_id" json:"cust_id"`
	CustName          string     `gorm:"column:cust_name" json:"cust_name"`
	GrBranchNo        string     `gorm:"column:gr_branch_no" json:"gr_branch_no"`
	GrBranchDate      *time.Time `gorm:"column:gr_branch_date" json:"gr_branch_date"`
	DeliveryNo        *string    `gorm:"column:delivery_no" json:"delivery_no"`
	DeliveryDate      *time.Time `gorm:"column:delivery_date" json:"delivery_date"`
	InvoiceNo         *string    `gorm:"column:invoice_no" json:"invoice_no"`
	InvoiceDate       *time.Time `gorm:"column:invoice_date" json:"invoice_date"`
	InvoiceNoBranch   *string    `gorm:"column:invoice_no_branch" json:"invoice_no_branch"`
	InvoiceDateBranch *time.Time `gorm:"column:invoice_date_branch" json:"invoice_date_branch"`
	PoNo              string     `gorm:"column:po_no" json:"po_no"`
	SoNo              string     `gorm:"column:so_no" json:"so_no"`
	VehicleNo         string     `gorm:"column:vehicle_no" json:"vehicle_no"`
	SupId             *int64     `gorm:"column:sup_id" json:"sup_id"`
	SupCode           *string    `gorm:"column:sup_code" json:"sup_code"`
	SupName           *string    `gorm:"column:sup_name" json:"sup_name"`
	WhId              *int64     `gorm:"column:wh_id" json:"wh_id"`
	WhCode            *string    `gorm:"column:wh_code" json:"wh_code"`
	WhName            *string    `gorm:"column:wh_name" json:"wh_name"`
	UpdatedBy         *int64     `gorm:"column:updated_by" json:"updated_by"`
	UpdatedByName     *string    `gorm:"column:updated_by_name" json:"updated_by_name"`
	UpdatedAt         *time.Time `gorm:"column:updated_at;default:null" json:"updated_at,omitempty"`
	IsPrint           *bool      `gorm:"column:is_print" json:"is_print"`
	PrintedBy         *int64     `gorm:"column:printed_by" json:"printed_by"`
	PrintedByName     *string    `gorm:"column:printed_by_name" json:"printed_by_name"`
	PrintedAt         *time.Time `gorm:"column:printed_at;default:null" json:"printed_at,omitempty"`
	DataStatus        *int64     `gorm:"column:data_status" json:"data_status"`
	DataStatusName    *string    `gorm:"column:data_status_name" json:"data_status_name"`
}

func (GrBranchList) TableName() string {
	return "inv.gr_branch"
}

type GrBranchRead struct {
	CustID               string     `gorm:"column:cust_id" json:"cust_id"`
	CustName             string     `gorm:"column:cust_name" json:"cust_name"`
	GrBranchNo           string     `gorm:"column:gr_branch_no" json:"gr_branch_no"`
	GrBranchDate         *time.Time `gorm:"column:gr_branch_date" json:"gr_branch_date"`
	DeliveryDate         *time.Time `gorm:"column:delivery_date" json:"delivery_date"`
	DeliveryNo           *string    `gorm:"column:delivery_no" json:"delivery_no"`
	DeliveryFee          *float64   `gorm:"column:delivery_fee" json:"delivery_fee"`
	InvoiceDate          *time.Time `gorm:"column:invoice_date" json:"invoice_date"`
	InvoiceNo            *string    `gorm:"column:invoice_no" json:"invoice_no"`
	PoNo                 *string    `gorm:"column:po_no" json:"po_no"`
	SoNo                 *string    `gorm:"column:so_no" json:"so_no"`
	VehicleNo            *string    `gorm:"column:vehicle_no" json:"vehicle_no"`
	SupId                *int64     `gorm:"column:sup_id" json:"sup_id"`
	SupCode              *string    `gorm:"column:sup_code" json:"sup_code"`
	SupName              *string    `gorm:"column:sup_name" json:"sup_name"`
	WhId                 *int64     `gorm:"column:wh_id" json:"wh_id"`
	WhCode               *string    `gorm:"column:wh_code" json:"wh_code"`
	WhName               *string    `gorm:"column:wh_name" json:"wh_name"`
	SubTotal             *float64   `gorm:"column:sub_total" json:"sub_total"`
	VatValue             *float64   `gorm:"column:vat_value" json:"vat_value"`
	Total                *float64   `gorm:"column:total" json:"total"`
	Notes                *string    `gorm:"column:notes" json:"notes"`
	UpdatedBy            *int64     `gorm:"column:updated_by" json:"updated_by"`
	UpdatedByName        *string    `gorm:"column:updated_by_name" json:"updated_by_name"`
	UpdatedAt            *time.Time `gorm:"column:updated_at;default:null" json:"updated_at,omitempty"`
	IsPrint              *bool      `gorm:"column:is_print" json:"is_print"`
	PrintedBy            *int64     `gorm:"column:printed_by" json:"printed_by"`
	PrintedByName        *string    `gorm:"column:printed_by_name" json:"printed_by_name"`
	PrintedAt            *time.Time `gorm:"column:printed_at;default:null" json:"printed_at,omitempty"`
	DataStatus           *int64     `gorm:"column:data_status" json:"data_status"`
	DataStatusName       *string    `gorm:"column:data_status_name" json:"data_status_name"`
	TypeApproval         *int       `gorm:"column:type_approval" json:"type_approval"`
	InvoiceNoBranch      *string    `gorm:"column:invoice_no_branch" json:"invoice_no_branch"`
	InvoiceDateBranch    *time.Time `gorm:"column:invoice_date_branch" json:"invoice_date_branch"`
	InvoiceDueDateBranch *time.Time `gorm:"column:invoice_due_date_branch" json:"invoice_due_date_branch"`
}

func (GrBranchRead) TableName() string {
	return "inv.gr_branch"
}

type GrBranchSupplier struct {
	SupId   *int64  `gorm:"column:sup_id" json:"sup_id"`
	SupCode *string `gorm:"column:sup_code" json:"sup_code"`
	SupName *string `gorm:"column:sup_name" json:"sup_name"`
}

func (GrBranchSupplier) TableName() string {
	return "inv.gr_branch"
}

type GrBranchDistributor struct {
	CustId   *string `gorm:"column:cust_id" json:"cust_id"`
	CustName *string `gorm:"column:cust_name" json:"cust_name"`
}

func (GrBranchDistributor) TableName() string {
	return "inv.gr_branch"
}

type GrBranchWarehouse struct {
	WhID   *int64  `gorm:"column:wh_id" json:"wh_id"`
	WhCode *string `gorm:"column:wh_code" json:"wh_code"`
	WhName *string `gorm:"column:wh_name" json:"wh_name"`
}

func (GrBranchWarehouse) TableName() string {
	return "mst.m_warehouse"
}

type GrBranchOrderBookingDetail struct {
	CustID               string   `gorm:"cust_id" json:"cust_id"`
	OrderBookingId       int      `gorm:"order_booking_id" json:"order_booking_id"`
	OrderBookingDetailId int      `gorm:"order_booking_detail_id" json:"order_booking_detail_id"`
	ProId                int      `gorm:"pro_id" json:"pro_id"`
	ProCode              string   `gorm:"pro_code" json:"pro_code"`
	ProName              string   `gorm:"pro_name" json:"pro_name"`
	ItemType             int      `gorm:"item_type" json:"item_type"`
	QtyBo                float64  `gorm:"qty_bo" json:"qty_bo"`
	QtyAlloc             float64  `gorm:"qty_alloc" json:"qty_alloc"`
	Qty1                 *float64 `gorm:"qty1" json:"qty1"`
	Qty2                 *float64 `gorm:"qty2" json:"qty2"`
	Qty3                 *float64 `gorm:"qty3" json:"qty3"`
	Qty4                 *float64 `gorm:"qty4" json:"qty4"`
	Qty5                 *float64 `gorm:"qty5" json:"qty5"`
	Qty1Alloc            *float64 `gorm:"qty1_alloc" json:"qty1_alloc"`
	Qty2Alloc            *float64 `gorm:"qty2_alloc" json:"qty2_alloc"`
	Qty3Alloc            *float64 `gorm:"qty3_alloc" json:"qty3_alloc"`
	Qty4Alloc            *float64 `gorm:"qty4_alloc" json:"qty4_alloc"`
	Qty5Alloc            *float64 `gorm:"qty5_alloc" json:"qty5_alloc"`
	PurchPrice1          *float64 `gorm:"purch_price1" json:"purch_price1"`
	PurchPrice2          *float64 `gorm:"purch_price2" json:"purch_price2"`
	PurchPrice3          *float64 `gorm:"purch_price3" json:"purch_price3"`
	PurchPrice4          *float64 `gorm:"purch_price4" json:"purch_price4"`
	PurchPrice5          *float64 `gorm:"purch_price5" json:"purch_price5"`
	SellPrice1           *float64 `gorm:"sell_price1" json:"sell_price1"`
	SellPrice2           *float64 `gorm:"sell_price2" json:"sell_price2"`
	SellPrice3           *float64 `gorm:"sell_price3" json:"sell_price3"`
	SellPrice4           *float64 `gorm:"sell_price4" json:"sell_price4"`
	SellPrice5           *float64 `gorm:"sell_price5" json:"sell_price5"`
	Amount               *float64 `gorm:"amount" json:"amount"`
	AmountAlloc          *float64 `gorm:"amount_alloc" json:"amount_alloc"`
	Vat                  *float64 `gorm:"vat" json:"vat"`
	VatValue             *float64 `gorm:"vat_value" json:"vat_value"`
	VatValueAlloc        *float64 `gorm:"vat_value_alloc" json:"vat_value_alloc"`
	UnitId1              *string  `gorm:"unit_id1" json:"unit_id1"`
	UnitId2              *string  `gorm:"unit_id2" json:"unit_id2"`
	UnitId3              *string  `gorm:"unit_id3" json:"unit_id3"`
	UnitId4              *string  `gorm:"unit_id4" json:"unit_id4"`
	UnitId5              *string  `gorm:"unit_id5" json:"unit_id5"`
	ConvUnit3            *int     `gorm:"conv_unit3" json:"conv_unit3"`
	ConvUnit2            *int     `gorm:"conv_unit2" json:"conv_unit2"`
	Notes                *string  `gorm:"notes" json:"notes"`
}

func (GrBranchOrderBookingDetail) TableName() string {
	return "inv.order_booking_detail"
}

type GrBranchOrderBooking struct {
	OrderBookingId int64    `gorm:"column:order_booking_id" json:"order_booking_id"`
	PoNo           *string  `gorm:"column:po_no" json:"po_no"`
	TypeApproval   *int     `gorm:"column:type_approval" json:"type_approval"`
	SoNo           *string  `gorm:"column:so_no" json:"so_no"`
	SupID          int64    `gorm:"column:sup_id" json:"sup_id"`
	SupCode        string   `gorm:"column:sup_code" json:"sup_code"`
	SupName        string   `gorm:"column:sup_name" json:"sup_name"`
	DeliveryFee    *float64 `gorm:"column:delivery_fee" json:"delivery_fee"`
}

func (GrBranchOrderBooking) TableName() string {
	return "inv.order_booking"
}

type InvoiceNoGrBranch struct {
	InvoiceNoGrBranch string `gorm:"column:get_no_fn"`
}

type TermOfPaySupplierGrBranch struct {
	TermOfPay int `gorm:"column:pay_term"`
}
