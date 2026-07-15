package model

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

type Order struct {
	CustID        string     `gorm:"cust_id" json:"cust_id"`
	RoNo          string     `gorm:"ro_no" json:"ro_no"`
	SalesmanId    *int64     `gorm:"salesman_id" json:"salesman_id"`
	WhId          *int64     `gorm:"wh_id" json:"wh_id"`
	RoDate        *time.Time `gorm:"ro_date" json:"ro_date"`
	ValDate       *time.Time `gorm:"val_date" json:"val_date"`
	OutletID      *int64     `gorm:"outlet_id" json:"outlet_id"`
	DeliveryDate  *time.Time `gorm:"delivery_date" json:"delivery_date"`
	OrderNo       *string    `gorm:"order_no" json:"order_no"`
	PoNo          *string    `gorm:"po_no" json:"po_no"`
	VehicleNo     *string    `gorm:"vehicle_no" json:"vehicle_no"`
	PayType       *int64     `gorm:"pay_type" json:"pay_type"`
	ReffNo        *string    `gorm:"reff_no" json:"reff_no"`
	MobileID      *int64     `gorm:"mobile_id" json:"mobile_id"`
	SubTotal      *float64   `gorm:"sub_total" json:"sub_total"`
	Disc          *float64   `gorm:"disc" json:"disc"`
	DiscValue     *float64   `gorm:"disc_value" json:"disc_value"`
	PromoValue    *float64   `gorm:"promo_value" json:"promo_value"`
	CashDiscValue *float64   `gorm:"cash_disc_value" json:"cash_disc_value"`
	TotDisc1      *float64   `gorm:"tot_disc1" json:"tot_disc1"`
	TotDisc2      *float64   `gorm:"tot_disc2" json:"tot_disc2"`
	Vat           *float64   `gorm:"vat" json:"vat"`
	VatValue      *float64   `gorm:"vat_value" json:"vat_value"`
	Total         *float64   `gorm:"total" json:"total"`
	DataStatus    *int64     `gorm:"data_status" json:"data_status"`
	DataSource    *int64     `gorm:"data_source" json:"data_source"`

	ValidateStok        bool    `gorm:"validate_stok" json:"validate_stok"`
	ValidateStokMessage *string `gorm:"validate_stok_message" json:"validate_stok_message"`

	TrCode *string `gorm:"tr_code" json:"tr_code"`
	// IsClosed    bool       `gorm:"is_closed" json:"is_closed"`
	// ClosedBy    *int64     `gorm:"closed_by" json:"closed_by"`
	// ClosedAt    *time.Time `gorm:"closed_at" json:"closed_at"`
	DueDate     *time.Time `gorm:"due_date" json:"due_date"`
	Notes       *string    `gorm:"notes" json:"notes"`
	InvoiceNo   *string    `gorm:"invoice_no" json:"invoice_no"`
	InvoiceDate *time.Time `gorm:"invoice_date" json:"invoice_date"`
	OprType     *string    `gorm:"opr_type" json:"opr_type"`

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

type NoOrder struct {
	CustID        string     `gorm:"cust_id" json:"cust_id"`
	SalesmanId    *int64     `gorm:"salesman_id" json:"salesman_id"`
	NoOrderDate   *time.Time `gorm:"no_order_date" json:"no_order_date"`
	OutletId      *int64     `gorm:"outlet_id" json:"outlet_id"`
	TakingOrderId *int64     `gorm:"taking_order_id" json:"taking_order_id"`

	Reason *string `gorm:"reason" json:"reason"`

	CreatedBy *int64          `gorm:"created_by" json:"created_by"`
	CreatedAt time.Time       `gorm:"created_at" json:"created_at"`
	UpdatedBy *int64          `gorm:"updated_by" json:"updated_by"`
	UpdatedAt time.Time       `gorm:"updated_at" json:"updated_at"`
	IsDel     bool            `gorm:"is_del" json:"is_del"`
	DeletedBy *int64          `gorm:"deleted_by" json:"deleted_by"`
	DeletedAt *gorm.DeletedAt `gorm:"deleted_at" json:"deleted_at"`
}

func (NoOrder) TableName() string {
	return "sls.no_order"
}

func (m *NoOrder) BeforeCreate(trx *gorm.DB) (err error) {
	// intTmpsStr := time.Now().UnixNano() / int64(time.Millisecond)
	// m.RoNo = strconv.Itoa(int(intTmpsStr))
	m.CreatedAt = time.Now()
	m.UpdatedAt = time.Now()
	m.UpdatedBy = m.CreatedBy
	return nil
}

type OrderList struct {
	CustID        string     `gorm:"cust_id" json:"cust_id"`
	RoNo          string     `gorm:"ro_no" json:"ro_no"`
	SalesmanId    *int64     `gorm:"salesman_id" json:"salesman_id"`
	SalesName     *string    `gorm:"sales_name" json:"sales_name"`
	WhId          *int64     `gorm:"wh_id" json:"wh_id"`
	RoDate        *time.Time `gorm:"ro_date" json:"ro_date"`
	ValDate       *time.Time `gorm:"val_date" json:"val_date"`
	WhCode        *string    `gorm:"wh_code" json:"wh_code"`
	WhName        *string    `gorm:"wh_name" json:"wh_name"`
	OutletID      *int64     `gorm:"outlet_id" json:"outlet_id"`
	OutletCode    *string    `gorm:"column:outlet_code" json:"outlet_code"`
	OutletName    *string    `gorm:"column:outlet_name" json:"outlet_name"`
	DeliveryDate  *time.Time `gorm:"delivery_date" json:"delivery_date"`
	OrderNo       *string    `gorm:"order_no" json:"order_no"`
	PoNo          *string    `gorm:"po_no" json:"po_no"`
	VehicleNo     *string    `gorm:"vehicle_no" json:"vehicle_no"`
	PayType       *int64     `gorm:"pay_type" json:"pay_type"`
	ReffNo        *string    `gorm:"reff_no" json:"reff_no"`
	MobileID      *int64     `gorm:"mobile_id" json:"mobile_id"`
	SubTotal      *float64   `gorm:"sub_total" json:"sub_total"`
	Disc          *float64   `gorm:"disc" json:"disc"`
	DiscValue     *float64   `gorm:"disc_value" json:"disc_value"`
	PromoValue    *float64   `gorm:"promo_value" json:"promo_value"`
	CashDiscValue *float64   `gorm:"cash_disc_value" json:"cash_disc_value"`
	TotDisc1      *float64   `gorm:"tot_disc1" json:"tot_disc1"`
	TotDisc2      *float64   `gorm:"tot_disc2" json:"tot_disc2"`
	Vat           *float64   `gorm:"vat" json:"vat"`
	VatValue      *float64   `gorm:"vat_value" json:"vat_value"`
	Total         *float64   `gorm:"total" json:"total"`
	DataStatus    *int64     `gorm:"data_status" json:"data_status"`
	DataSource    *int64     `gorm:"data_source" json:"data_source"`

	ValidateStok        bool    `gorm:"validate_stok" json:"validate_stok"`
	ValidateStokMessage *string `gorm:"validate_stok_message" json:"validate_stok_message"`

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
}

func (OrderList) TableName() string {
	return "sls.order"
}

type NoOrderList struct {
	CustID string `gorm:"cust_id" json:"cust_id"`
	// RoNo        string     `gorm:"ro_no" json:"ro_no"`
	NoOrderId        *int32     `gorm:"no_order_id" json:"no_order_id"`
	SalesmanId       *int64     `gorm:"salesman_id" json:"salesman_id"`
	SalesName        *string    `gorm:"sales_name" json:"sales_name"`
	Reason           *string    `gorm:"reason" json:"reason"`
	NoOrderDate      *time.Time `gorm:"no_order_date" json:"no_order_date"`
	OutletID         *int64     `gorm:"outlet_id" json:"outlet_id"`
	OutletCode       *string    `gorm:"column:outlet_code" json:"outlet_code"`
	OutletName       *string    `gorm:"column:outlet_name" json:"outlet_name"`
	TakingOrderId    *int64     `gorm:"taking_order_id" json:"taking_order_id"`
	TakingOrderName  *string    `gorm:"taking_order_name" json:"taking_order_name"`
	TakingOrderImage *string    `gorm:"column:image_url" json:"image_url"`

	CreatedBy     *int64          `gorm:"created_by" json:"created_by"`
	CreatedAt     time.Time       `gorm:"created_at" json:"created_at"`
	UpdatedBy     *int64          `gorm:"updated_by" json:"updated_by"`
	UpdatedAt     time.Time       `gorm:"updated_at" json:"updated_at"`
	UpdatedByName *string         `gorm:"column:updated_by_name" json:"updated_by_name"`
	IsDel         bool            `gorm:"is_del" json:"is_del"`
	DeletedBy     *int64          `gorm:"deleted_by" json:"deleted_by"`
	DeletedAt     *gorm.DeletedAt `gorm:"deleted_at" json:"deleted_at"`
}

func (NoOrderList) TableName() string {
	return "sls.no_order"
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

type SummaryOrder struct {
	TotalSummary float64 `json:"total_summary" db:"total_summary"`
	SalesmanID   int     `json:"salesman_id" db:"salesman_id"`
	CreatedAt    time.Time
}

func (SummaryOrder) TableName() string {
	return "sls.order"
}

type OutletBySalesman struct {
	OutletID   int     `json:"outlet_id" db:"outlet_id"`
	OutletCode *string `json:"outlet_code" db:"outlet_code"`
	OutletName *string `json:"outlet_name" db:"outlet_name"`
}

func (OutletBySalesman) TableName() string {
	return "sls.order"
}

type ProductOrder struct {
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
	Cogs        float64 `json:"cogs" gorm:"column:cogs"`
}

func (ProductOrder) TableName() string {
	return "mst.m_product"
}

type MapProduct map[int64]ProductOrder

func (m MapProduct) SetProduct(id int64, product ProductOrder) {
	m[id] = product
}

func (m MapProduct) GetByID(id int64) (product ProductOrder, err error) {
	val, ok := m[id]

	// If the key exists
	if !ok {
		return product, fmt.Errorf("ProductOrder ID %v Not Found", id)
	}

	return val, nil
}

type OrderListInvoice struct {
	InvoiceNo       string     `json:"invoice_no" db:"invoice_no"`
	InvoiceDate     *time.Time `json:"invoice_date" db:"invoice_date"`
	DueDate         *time.Time `json:"due_date" db:"due_date"`
	OutletID        int64      `json:"outlet_id" db:"outlet_id"`
	OtGrpID         int64      `json:"ot_grp_id" db:"ot_grp_id"`
	SalesmanID      int64      `json:"salesman_id" db:"salesman_id"`
	RoNo            string     `json:"ro_no" db:"ro_no"`
	InvoiceAmount   float64    `json:"invoice_amount" db:"invoice_amount"`
	OutletCode      string     `json:"outlet_code" db:"outlet_code"`
	OutletName      string     `json:"outlet_name" db:"outlet_name"`
	PaidAmount      float64    `json:"paid_amount" db:"paid_amount"`
	RemainingAmount float64    `json:"remaining_amount" db:"remaining_amount"`
}
