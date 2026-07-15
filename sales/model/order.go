package model

import (
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type Order struct {
	CustID            string          `gorm:"cust_id" json:"cust_id"`
	RoNo              string          `gorm:"ro_no" json:"ro_no"`
	OprType           *string         `gorm:"column:opr_type" json:"opr_type"`
	OrderType         *string         `gorm:"column:order_type" json:"order_type"`
	IsSalesMapping    *bool           `gorm:"column:is_sales_mapping" json:"is_sales_mapping"`
	SalesmanId        *int64          `gorm:"salesman_id" json:"salesman_id"`
	WhId              *int64          `gorm:"wh_id" json:"wh_id"`
	RoDate            *time.Time      `gorm:"ro_date" json:"ro_date"`
	ValDate           *time.Time      `gorm:"val_date" json:"val_date"`
	OutletID          *int64          `gorm:"outlet_id" json:"outlet_id"`
	DeliveryDate      *time.Time      `gorm:"delivery_date" json:"delivery_date"`
	OrderNo           *string         `gorm:"order_no" json:"order_no"`
	PoNo              *string         `gorm:"po_no" json:"po_no"`
	VehicleNo         *string         `gorm:"vehicle_no" json:"vehicle_no"`
	PayType           *int64          `gorm:"pay_type" json:"pay_type"`
	ReffNo            *string         `gorm:"reff_no" json:"reff_no"`
	MobileID          *int64          `gorm:"mobile_id" json:"mobile_id"`
	SubTotal          *float64        `gorm:"sub_total" json:"sub_total"`
	SubTotalFinal     *float64        `gorm:"sub_total_final" json:"sub_total_final"`
	Disc              *float64        `gorm:"disc" json:"disc"`
	DiscValue         *float64        `gorm:"disc_value" json:"disc_value"`
	DiscValueFinal    *float64        `gorm:"disc_value_final" json:"disc_value_final"`
	PromoValue        *float64        `gorm:"promo_value" json:"promo_value"`
	PromoValueFinal   *float64        `gorm:"promo_value_final" json:"promo_value_final"`
	PromoBgValue      *float64        `gorm:"promo_bg_value" json:"promo_bg_value"`
	PromoBgValueFinal *float64        `gorm:"promo_bg_value_final" json:"promo_bg_value_final"`
	PromoRemarksSo    JSONStringArray `gorm:"column:promo_remarks_so;type:jsonb" json:"promo_remarks_so"`
	PromoRemarksFinal JSONStringArray `gorm:"column:promo_remarks_final;type:jsonb" json:"promo_remarks_final"`
	PromoRemarksPo    JSONStringArray `gorm:"column:promo_remarks_po;type:jsonb" json:"promo_remarks_po"`
	CashDiscValue     *float64        `gorm:"cash_disc_value" json:"cash_disc_value"`
	TotDisc1          *float64        `gorm:"tot_disc1" json:"tot_disc1"`
	TotDisc2          *float64        `gorm:"tot_disc2" json:"tot_disc2"`
	Vat               *float64        `gorm:"vat" json:"vat"`
	VatValue          *float64        `gorm:"vat_value" json:"vat_value"`
	VatValueFinal     *float64        `gorm:"vat_value_final" json:"vat_value_final"`
	Total             *float64        `gorm:"total" json:"total"`
	TotalFinal        *float64        `gorm:"total_final" json:"total_final"`
	DataStatus        *int64          `gorm:"data_status" json:"data_status"`
	DataSource        *int64          `gorm:"data_source" json:"data_source"`

	TrCode *string `gorm:"tr_code" json:"tr_code"`
	// IsClosed    bool       `gorm:"is_closed" json:"is_closed"`
	// ClosedBy    *int64     `gorm:"closed_by" json:"closed_by"`
	// ClosedAt    *time.Time `gorm:"closed_at" json:"closed_at"`
	DueDate     *time.Time `gorm:"due_date" json:"due_date"`
	Notes       *string    `gorm:"notes" json:"notes"`
	InvoiceNo   *string    `gorm:"invoice_no" json:"invoice_no"`
	InvoiceDate *time.Time `gorm:"invoice_date" json:"invoice_date"`

	Address1 *string `gorm:"address1" json:"address1"`

	ValidateStok               *bool   `gorm:"validate_stok" json:"validate_stok"`
	ValidateStokMessage        *string `gorm:"validate_stok_message" json:"validate_stok_message"`
	ValidateCreditLimit        *bool   `gorm:"validate_credit_limit" json:"validate_credit_limit"`
	ValidateCreditLimitMessage string  `gorm:"validate_credit_limit_message" json:"validate_credit_limit_message"`
	ValidateCreditLimitValue   float64 `gorm:"validate_credit_limit_value" json:"validate_credit_limit_value"`
	ValidateOverdue            *bool   `gorm:"validate_overdue" json:"validate_overdue"`
	ValidateOverdueMessage     string  `gorm:"validate_overdue_message" json:"validate_overdue_message"`
	ValidateOverdueValue       int     `gorm:"validate_overdue_value" json:"validate_overdue_value"`
	ValidateOutstanding        *bool   `gorm:"validate_outstanding" json:"validate_outstanding"`
	ValidateOutstandingMessage string  `gorm:"validate_outstanding_message" json:"validate_outstanding_message"`
	ValidateOutstandingValue   int     `gorm:"validate_outstanding_value" json:"validate_outstanding_value"`
	ValidateSummary            bool    `gorm:"validate_summary" json:"validate_summary"`

	CreatedBy *int64          `gorm:"created_by" json:"created_by"`
	CreatedAt time.Time       `gorm:"created_at" json:"created_at"`
	UpdatedBy *int64          `gorm:"updated_by" json:"updated_by"`
	UpdatedAt time.Time       `gorm:"updated_at" json:"updated_at"`
	IsDel     bool            `gorm:"is_del" json:"is_del"`
	DeletedBy *int64          `gorm:"deleted_by" json:"deleted_by"`
	DeletedAt *gorm.DeletedAt `gorm:"deleted_at" json:"deleted_at"`
}

func (Order) TableName() string {
	return "sls.order"
}

func (m *Order) BeforeCreate(trx *gorm.DB) (err error) {
	// intTmpsStr := time.Now().UnixNano() / int64(time.Millisecond)
	// m.RoNo = strconv.Itoa(int(intTmpsStr))
	m.CreatedAt = time.Now()
	m.UpdatedAt = time.Now()
	m.UpdatedBy = m.CreatedBy
	return nil
}

type OrderList struct {
	CustID            string          `gorm:"cust_id" json:"cust_id"`
	RoNo              string          `gorm:"ro_no" json:"ro_no"`
	OprType           *string         `gorm:"column:opr_type" json:"opr_type"`
	OrderType         *string         `gorm:"column:order_type" json:"order_type"`
	IsSalesMapping    *bool           `gorm:"column:is_sales_mapping" json:"is_sales_mapping"`
	SalesmanId        *int64          `gorm:"salesman_id" json:"salesman_id"`
	SalesmanCode      *string         `gorm:"salesman_code" json:"salesman_code"`
	SalesName         *string         `gorm:"sales_name" json:"sales_name"`
	WhId              *int64          `gorm:"wh_id" json:"wh_id"`
	RoDate            *time.Time      `gorm:"ro_date" json:"ro_date"`
	ValDate           *time.Time      `gorm:"val_date" json:"val_date"`
	WhCode            *string         `gorm:"wh_code" json:"wh_code"`
	WhName            *string         `gorm:"wh_name" json:"wh_name"`
	OutletID          *int64          `gorm:"outlet_id" json:"outlet_id"`
	OutletCode        *string         `gorm:"column:outlet_code" json:"outlet_code"`
	OutletName        *string         `gorm:"column:outlet_name" json:"outlet_name"`
	OutletAddress1    *string         `gorm:"column:address1" json:"outlet_address1"`
	OutletAddress2    *string         `gorm:"column:address2" json:"outlet_address2"`
	InvAddress1       *string         `gorm:"column:inv_addr1" json:"inv_addr1"`
	InvAddress2       *string         `gorm:"column:inv_addr2" json:"inv_addr2"`
	ZipCode           *string         `gorm:"column:zip_code" json:"zip_code"`
	DeliveryDate      *time.Time      `gorm:"delivery_date" json:"delivery_date"`
	OrderNo           *string         `gorm:"order_no" json:"order_no"`
	PoNo              *string         `gorm:"po_no" json:"po_no"`
	VehicleNo         *string         `gorm:"vehicle_no" json:"vehicle_no"`
	PayType           *int64          `gorm:"pay_type" json:"pay_type"`
	ReffNo            *string         `gorm:"reff_no" json:"reff_no"`
	MobileID          *int64          `gorm:"mobile_id" json:"mobile_id"`
	SubTotal          *float64        `gorm:"sub_total" json:"sub_total"`
	SubTotalFinal     *float64        `gorm:"sub_total_final" json:"sub_total_final"`
	Disc              *float64        `gorm:"disc" json:"disc"`
	DiscValue         *float64        `gorm:"disc_value" json:"disc_value"`
	DiscValueFinal    *float64        `gorm:"disc_value_final" json:"disc_value_final"`
	PromoValue        *float64        `gorm:"promo_value" json:"promo_value"`
	PromoValueFinal   *float64        `gorm:"promo_value_final" json:"promo_value_final"`
	PromoBgValue      *float64        `gorm:"promo_bg_value" json:"promo_bg_value"`
	PromoBgValueFinal *float64        `gorm:"promo_bg_value_final" json:"promo_bg_value_final"`
	PromoRemarksSo    JSONStringArray `gorm:"column:promo_remarks_so;type:jsonb" json:"promo_remarks_so"`
	PromoRemarksFinal JSONStringArray `gorm:"column:promo_remarks_final;type:jsonb" json:"promo_remarks_final"`
	PromoRemarksPo    JSONStringArray `gorm:"column:promo_remarks_po;type:jsonb" json:"promo_remarks_po"`
	CashDiscValue     *float64        `gorm:"cash_disc_value" json:"cash_disc_value"`
	TotDisc1          *float64        `gorm:"tot_disc1" json:"tot_disc1"`
	TotDisc2          *float64        `gorm:"tot_disc2" json:"tot_disc2"`
	Vat               *float64        `gorm:"vat" json:"vat"`
	VatValue          *float64        `gorm:"vat_value" json:"vat_value"`
	VatValueFinal     *float64        `gorm:"vat_value_final" json:"vat_value_final"`
	Total             *float64        `gorm:"total" json:"total"`
	TotalFinal        *float64        `gorm:"total_final" json:"total_final"`
	DataStatus        *int64          `gorm:"data_status" json:"data_status"`
	DataSource        *int64          `gorm:"data_source" json:"data_source"`
	IsProformaInv     *bool           `gorm:"is_proforma_inv" json:"is_proforma_inv"`
	GenerateBy        *int64          `gorm:"generate_by" json:"generate_by"`
	FirstIssueDate    *time.Time      `gorm:"first_issue_date" json:"first_issue_date"`

	TrCode *string `gorm:"tr_code" json:"tr_code"`
	// IsClosed    bool       `gorm:"is_closed" json:"is_closed"`
	// ClosedBy    *int64     `gorm:"closed_by" json:"closed_by"`
	// ClosedAt    *time.Time `gorm:"closed_at" json:"closed_at"`
	DueDate     *time.Time `gorm:"due_date" json:"due_date"`
	Notes       *string    `gorm:"notes" json:"notes"`
	InvoiceNo   *string    `gorm:"invoice_no" json:"invoice_no"`
	InvoiceDate *time.Time `gorm:"invoice_date" json:"invoice_date"`

	CreatedBy     *int64          `gorm:"created_by" json:"created_by"`
	CreatedAt     time.Time       `gorm:"created_at" json:"created_at"`
	UpdatedBy     *int64          `gorm:"updated_by" json:"updated_by"`
	UpdatedAt     time.Time       `gorm:"updated_at" json:"updated_at"`
	UpdatedByName *string         `gorm:"column:updated_by_name" json:"updated_by_name"`
	IsDel         bool            `gorm:"is_del" json:"is_del"`
	DeletedBy     *int64          `gorm:"deleted_by" json:"deleted_by"`
	DeletedAt     *gorm.DeletedAt `gorm:"deleted_at" json:"deleted_at"`

	IsPrinted     bool       `gorm:"column:is_printed" json:"is_printed"`
	PrintedBy     *int64     `gorm:"column:printed_by" json:"printed_by"`
	PrintedByName *string    `gorm:"column:printed_by_name" json:"printed_by_name"`
	PrintedAt     *time.Time `gorm:"column:printed_at" json:"printed_at"`

	ValidateStok               bool    `gorm:"validate_stok" json:"validate_stok"`
	ValidateStokMessage        *string `gorm:"validate_stok_message" json:"validate_stok_message"`
	ValidateCreditLimit        bool    `gorm:"validate_credit_limit" json:"validate_credit_limit"`
	ValidateCreditLimitMessage string  `gorm:"validate_credit_limit_message" json:"validate_credit_limit_message"`
	ValidateCreditLimitValue   float64 `gorm:"validate_credit_limit_value" json:"validate_credit_limit_value"`
	ValidateOverdue            bool    `gorm:"validate_overdue" json:"validate_overdue"`
	ValidateOverdueMessage     string  `gorm:"validate_overdue_message" json:"validate_overdue_message"`
	ValidateOverdueValue       int     `gorm:"validate_overdue_value" json:"validate_overdue_value"`
	ValidateOutstanding        bool    `gorm:"validate_outstanding" json:"validate_outstanding"`
	ValidateOutstandingMessage string  `gorm:"validate_outstanding_message" json:"validate_outstanding_message"`
	ValidateOutstandingValue   int     `gorm:"validate_outstanding_value" json:"validate_outstanding_value"`
	ValidateSummary            bool    `gorm:"validate_summary" json:"validate_summary"`
	CreditLimitType            *int    `json:"credit_limit_type" db:"credit_limit_type"`
	CreditLimitAction          *int    `gorm:"credit_limit_action" json:"credit_limit_action"`
	CreditLimitActionName      string  `gorm:"credit_limit_action_name" json:"credit_limit_action_name"`
	SalesInvLimitType          *int    `json:"sales_inv_limit_type" db:"sales_inv_limit_type"`
	SalesInvLimitAction        *int    `gorm:"sales_inv_limit_action" json:"sales_inv_limit_action"`
	SalesInvLimitActionName    string  `gorm:"sales_inv_limit_action_name" json:"sales_inv_limit_action_name"`
	ObsType                    *int    `json:"obs_type" db:"obs_type"`
	ObsLimitAction             *int    `gorm:"obs_limit_action" json:"obs_limit_action"`
	ObsLimitActionName         string  `gorm:"obs_limit_action_name" json:"obs_limit_action_name"`
	OrderApprovalRequestID     *int64  `gorm:"order_approval_request_id" json:"order_approval_request_id"`
}

func (OrderList) TableName() string {
	return "sls.order"
}

type ProductConversion struct {
	CustId    string  `json:"cust_id" db:"cust_id"`
	ProductId int64   `json:"pro_id" db:"pro_id"`
	ConvUnit2 float32 `json:"conv_unit2" db:"conv_unit2"`
	ConvUnit3 float32 `json:"conv_unit3" db:"conv_unit3"`
	ConvUnit4 float32 `json:"conv_unit4" db:"conv_unit4"`
	ConvUnit5 float32 `json:"conv_unit5" db:"conv_unit5"`
}

func (ProductConversion) TableName() string {
	return "mst.m_product"
}

type OutletBySalesman struct {
	OutletID   int     `json:"outlet_id" db:"outlet_id"`
	OutletCode *string `json:"outlet_code" db:"outlet_code"`
	OutletName *string `json:"outlet_name" db:"outlet_name"`
}

func (OutletBySalesman) TableName() string {
	return "sls.order"
}

type Product struct {
	CustId      string  `json:"cust_id" gorm:"column:cust_id"`
	ProductId   int64   `json:"pro_id" gorm:"column:pro_id"`
	UnitId1     string  `json:"unit_id1" gorm:"column:unit_id1"`
	UnitId2     string  `json:"unit_id2" gorm:"column:unit_id2"`
	UnitId3     string  `json:"unit_id3" gorm:"column:unit_id3"`
	UnitId4     *string `json:"unit_id4" gorm:"column:unit_id4"`
	UnitId5     *string `json:"unit_id5" gorm:"column:unit_id5"`
	ConvUnit2   float32 `json:"conv_unit2" gorm:"column:conv_unit2"`
	ConvUnit3   float32 `json:"conv_unit3" gorm:"column:conv_unit3"`
	ConvUnit4   float32 `json:"conv_unit4" gorm:"column:conv_unit4"`
	ConvUnit5   float32 `json:"conv_unit5" gorm:"column:conv_unit5"`
	PurchPrice1 float64 `json:"purch_price1" gorm:"column:purch_price1"`
	PurchPrice2 float64 `json:"purch_price2" gorm:"column:purch_price2"`
	PurchPrice3 float64 `json:"purch_price3" gorm:"column:purch_price3"`
	PurchPrice4 float64 `json:"purch_price4" gorm:"column:purch_price4"`
	PurchPrice5 float64 `json:"purch_price5" gorm:"column:purch_price5"`
	SellPrice1  float64 `json:"sell_price1" gorm:"column:sell_price1"`
	SellPrice2  float64 `json:"sell_price2" gorm:"column:sell_price2"`
	SellPrice3  float64 `json:"sell_price3" gorm:"column:sell_price3"`
	SellPrice4  float64 `json:"sell_price4" gorm:"column:sell_price4"`
	SellPrice5  float64 `json:"sell_price5" gorm:"column:sell_price5"`
	Vat         float64 `json:"vat" gorm:"column:vat"`
}

func (Product) TableName() string {
	return "mst.m_product"
}

type MapProduct map[int64]Product

func (m MapProduct) SetProduct(id int64, product Product) {
	m[id] = product
}

func (m MapProduct) GetByID(id int64) (product Product, err error) {
	val, ok := m[id]

	// If the key exists
	if !ok {
		return product, errors.New(fmt.Sprintf("Product ID %v Not Found", id))
	}

	return val, nil
}
