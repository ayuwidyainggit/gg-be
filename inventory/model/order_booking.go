package model

import (
	"time"

	"gorm.io/gorm"
)

type OrderBooking struct {
	CustID             string          `gorm:"cust_id" json:"cust_id"`
	OrderBookingId     *int            `gorm:"column:order_booking_id;autoIncrement" json:"order_booking_id"`
	PoNo               *string         `gorm:"po_no" json:"po_no"`
	SoPo               string          `gorm:"so_po" json:"so_po"`
	SupId              *int64          `gorm:"sup_id" json:"sup_id"`
	StatusOrderBooking int64           `gorm:"status_order_booking" json:"status_order_booking"`
	SubTotal           *float64        `gorm:"sub_total" json:"sub_total"`
	SubTotalAlloc      *float64        `gorm:"sub_total_alloc" json:"sub_total_alloc"`
	Vat                *float64        `gorm:"vat" json:"vat"`
	VatValue           *float64        `gorm:"vat_value" json:"vat_value"`
	VatValueAlloc      *float64        `gorm:"vat_value_alloc" json:"vat_value_alloc"`
	Total              *float64        `gorm:"total" json:"total"`
	TotalAlloc         *float64        `gorm:"total_alloc" json:"total_alloc"`
	CreatedBy          *int64          `gorm:"created_by" json:"created_by"`
	CreatedAt          time.Time       `gorm:"created_at" json:"created_at"`
	UpdatedBy          *int64          `gorm:"updated_by" json:"updated_by"`
	UpdatedAt          time.Time       `gorm:"updated_at" json:"updated_at"`
	IsDel              bool            `gorm:"is_del" json:"is_del"`
	DeletedBy          *int64          `gorm:"deleted_by" json:"deleted_by"`
	DeletedAt          *gorm.DeletedAt `gorm:"deleted_at" json:"deleted_at"`
}

func (OrderBooking) TableName() string {
	return "inv.order_booking"
}

func (m *OrderBooking) BeforeCreate(trx *gorm.DB) (err error) {
	m.CreatedAt = time.Now()
	m.UpdatedAt = time.Now()
	m.UpdatedBy = m.CreatedBy
	return nil
}

type OrderBookingList struct {
	CustID             string          `gorm:"cust_id" json:"cust_id"`
	OrderBookingId     int             `gorm:"order_booking_id" json:"order_booking_id"`
	PoNo               *string         `gorm:"po_no" json:"po_no"`
	GrBranchNo         *string         `gorm:"gr_branch_no" json:"gr_branch_no"`
	SoPo               *string         `gorm:"so_po" json:"so_po"`
	SupId              *int64          `gorm:"sup_id" json:"sup_id"`
	SupName            *string         `gorm:"sup_name" json:"sup_name"`
	SupCode            *string         `gorm:"sup_code" json:"sup_code"`
	CreditLimit        *float64        `gorm:"credit_limit" json:"credit_limit"`
	DistributorId      *int            `gorm:"distributor_id" json:"distributor_id"`
	DistributorName    *string         `gorm:"distributor_name" json:"distributor_name"`
	DistributorCode    *string         `gorm:"distributor_code" json:"distributor_code"`
	DistributorAddress *string         `gorm:"distributor_address" json:"distributor_address"`
	StatusOrderBooking *int64          `gorm:"status_order_booking" json:"status_order_booking"`
	SubTotal           *float64        `gorm:"sub_total" json:"sub_total"`
	SubTotalAlloc      *float64        `gorm:"sub_total_alloc" json:"sub_total_alloc"`
	Vat                *float64        `gorm:"vat" json:"vat"`
	VatValue           *float64        `gorm:"vat_value" json:"vat_value"`
	VatValueAlloc      *float64        `gorm:"vat_value_alloc" json:"vat_value_alloc"`
	DeliveryFee        *float64        `gorm:"delivery_fee" json:"delivery_fee"`
	DeliveryFeeFinal   *float64        `gorm:"delivery_fee_final" json:"delivery_fee_final"`
	SubTotalFinal      *float64        `gorm:"sub_total_final" json:"sub_total_final"`
	VatValueFinal      *float64        `gorm:"vat_value_final" json:"vat_value_final"`
	TotalFinal         *float64        `gorm:"total_final" json:"total_final"`
	Total              *float64        `gorm:"total" json:"total"`
	TotalAlloc         *float64        `gorm:"total_alloc" json:"total_alloc"`
	CreatedBy          *int64          `gorm:"created_by" json:"created_by"`
	CreatedByName      *string         `gorm:"created_by_name" json:"created_by_name"`
	CreatedAt          time.Time       `gorm:"created_at" json:"created_at"`
	UpdatedBy          *int64          `gorm:"updated_by" json:"updated_by"`
	UpdatedByName      *string         `gorm:"updated_by_name" json:"updated_by_name"`
	UpdatedAt          time.Time       `gorm:"updated_at" json:"updated_at"`
	IsDel              bool            `gorm:"is_del" json:"is_del"`
	DeletedBy          *int64          `gorm:"deleted_by" json:"deleted_by"`
	DeletedAt          *gorm.DeletedAt `gorm:"deleted_at" json:"deleted_at"`
}

func (OrderBookingList) TableName() string {
	return "inv.order_booking"
}

type OrderBookingDetail struct {
	CustID         string   `gorm:"cust_id" json:"cust_id"`
	OrderBookingId int      `gorm:"order_booking_id" json:"order_booking_id"`
	ProId          int      `gorm:"pro_id" json:"pro_id"`
	ItemType       int      `gorm:"item_type" json:"item_type"`
	QtyBo          float64  `gorm:"qty_bo" json:"qty_bo"`
	QtyAlloc       float64  `gorm:"qty_alloc" json:"qty_alloc"`
	Qty1           *float64 `gorm:"qty1" json:"qty1"`
	Qty2           *float64 `gorm:"qty2" json:"qty2"`
	Qty3           *float64 `gorm:"qty3" json:"qty3"`
	Qty4           *float64 `gorm:"qty4" json:"qty4"`
	Qty5           *float64 `gorm:"qty5" json:"qty5"`
	Qty1Alloc      *float64 `gorm:"qty1_alloc" json:"qty1_alloc"`
	Qty2Alloc      *float64 `gorm:"qty2_alloc" json:"qty2_alloc"`
	Qty3Alloc      *float64 `gorm:"qty3_alloc" json:"qty3_alloc"`
	Qty4Alloc      *float64 `gorm:"qty4_alloc" json:"qty4_alloc"`
	Qty5Alloc      *float64 `gorm:"qty5_alloc" json:"qty5_alloc"`
	PurchPrice1    *float64 `gorm:"purch_price1" json:"purch_price1"`
	PurchPrice2    *float64 `gorm:"purch_price2" json:"purch_price2"`
	PurchPrice3    *float64 `gorm:"purch_price3" json:"purch_price3"`
	PurchPrice4    *float64 `gorm:"purch_price4" json:"purch_price4"`
	PurchPrice5    *float64 `gorm:"purch_price5" json:"purch_price5"`
	SellPrice1     *float64 `gorm:"sell_price1" json:"sell_price1"`
	SellPrice2     *float64 `gorm:"sell_price2" json:"sell_price2"`
	SellPrice3     *float64 `gorm:"sell_price3" json:"sell_price3"`
	SellPrice4     *float64 `gorm:"sell_price4" json:"sell_price4"`
	SellPrice5     *float64 `gorm:"sell_price5" json:"sell_price5"`
	Amount         *float64 `gorm:"amount" json:"amount"`
	AmountAlloc    *float64 `gorm:"amount_alloc" json:"amount_alloc"`
	Vat            *float64 `gorm:"vat" json:"vat"`
	VatValue       *float64 `gorm:"vat_value" json:"vat_value"`
	VatValueAlloc  *float64 `gorm:"vat_value_alloc" json:"vat_value_alloc"`
	UnitId1        *string  `gorm:"unit_id1" json:"unit_id1"`
	UnitId2        *string  `gorm:"unit_id2" json:"unit_id2"`
	UnitId3        *string  `gorm:"unit_id3" json:"unit_id3"`
	UnitId4        *string  `gorm:"unit_id4" json:"unit_id4"`
	UnitId5        *string  `gorm:"unit_id5" json:"unit_id5"`
	ConvUnit2      *int     `gorm:"conv_unit2" json:"conv_unit2"`
	ConvUnit3      *int     `gorm:"conv_unit3" json:"conv_unit3"`
	ConvUnit4      *int     `gorm:"conv_unit4" json:"conv_unit4"`
	ConvUnit5      *int     `gorm:"conv_unit5" json:"conv_unit5"`
	Notes          *string  `gorm:"notes" json:"notes"`
}

func (OrderBookingDetail) TableName() string {
	return "inv.order_booking_detail"
}

type OrderBookingDetailRead struct {
	CustID               string   `gorm:"cust_id" json:"cust_id"`
	OrderBookingId       *int     `gorm:"order_booking_id" json:"order_booking_id"`
	OrderBookingDetailId *int     `gorm:"order_booking_detail_id" json:"order_booking_detail_id"`
	ProId                int      `gorm:"pro_id" json:"pro_id"`
	ProCode              string   `gorm:"column:pro_code" json:"pro_code"`
	ProName              string   `gorm:"column:pro_name" json:"pro_name"`
	SupId                *int     `gorm:"column:sup_id" json:"sup_id"`
	SupName              *string  `gorm:"column:sup_name" json:"sup_name"`
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
	MpConvUnit2          *int     `gorm:"column:mconv_unit2" json:"mconv_unit2"`
	MpConvUnit3          *int     `gorm:"column:mconv_unit3" json:"mconv_unit3"`
	ConvUnit2            *int     `gorm:"column:conv_unit2" json:"conv_unit2"`
	ConvUnit3            *int     `gorm:"column:conv_unit3" json:"conv_unit3"`
	ConvUnit4            *int     `gorm:"column:conv_unit4" json:"conv_unit4"`
	ConvUnit5            *int     `gorm:"column:conv_unit5" json:"conv_unit5"`
	Notes                *string  `gorm:"notes" json:"notes"`
}

type OrderBookingDetailFinalRead struct {
	CustID       string   `gorm:"cust_id" json:"cust_id"`
	ProId        int      `gorm:"pro_id" json:"pro_id"`
	ProCode      string   `gorm:"column:pro_code" json:"pro_code"`
	ProName      string   `gorm:"column:pro_name" json:"pro_name"`
	ItemType     int      `gorm:"item_type" json:"item_type"`
	QtyReceived  *float64 `gorm:"qty_received" json:"qty_received"`
	Qty1         *float64 `gorm:"qty1" json:"qty1"`
	Qty2         *float64 `gorm:"qty2" json:"qty2"`
	Qty3         *float64 `gorm:"qty3" json:"qty3"`
	Qty4         *float64 `gorm:"qty4" json:"qty4"`
	Qty5         *float64 `gorm:"qty5" json:"qty5"`
	QtyReceived1 *float64 `gorm:"qty_received1" json:"qty_received1"`
	QtyReceived2 *float64 `gorm:"qty_received2" json:"qty_received2"`
	QtyReceived3 *float64 `gorm:"qty_received3" json:"qty_received3"`
	QtyReceived4 *float64 `gorm:"qty_received4" json:"qty_received4"`
	QtyReceived5 *float64 `gorm:"qty_received5" json:"qty_received5"`
	Amount       *float64 `gorm:"amount" json:"amount"`
	Vat          *float64 `gorm:"vat" json:"vat"`
	VatValue     *float64 `gorm:"vat_value" json:"vat_value"`
	UnitPrice1   *float64 `gorm:"unit_price1" json:"unit_price1"`
	UnitPrice2   *float64 `gorm:"unit_price2" json:"unit_price2"`
	UnitPrice3   *float64 `gorm:"unit_price3" json:"unit_price3"`
	UnitPrice4   *float64 `gorm:"unit_price4" json:"unit_price4"`
	UnitPrice5   *float64 `gorm:"unit_price5" json:"unit_price5"`
	UnitId1      *string  `gorm:"unit_id1" json:"unit_id1"`
	UnitId2      *string  `gorm:"unit_id2" json:"unit_id2"`
	UnitId3      *string  `gorm:"unit_id3" json:"unit_id3"`
	UnitId4      *string  `gorm:"unit_id4" json:"unit_id4"`
	UnitId5      *string  `gorm:"unit_id5" json:"unit_id5"`
	MpConvUnit2  *int     `gorm:"column:mconv_unit2" json:"mconv_unit2"`
	MpConvUnit3  *int     `gorm:"column:mconv_unit3" json:"mconv_unit3"`
	ConvUnit2    *int     `gorm:"column:conv_unit2" json:"conv_unit2"`
	ConvUnit3    *int     `gorm:"column:conv_unit3" json:"conv_unit3"`
	ConvUnit4    *int     `gorm:"column:conv_unit4" json:"conv_unit4"`
	ConvUnit5    *int     `gorm:"column:conv_unit5" json:"conv_unit5"`
	Notes        *string  `gorm:"notes" json:"notes"`
}

func (OrderBookingDetailFinalRead) TableName() string {
	return "inv.gr_branch_det"
}

func (OrderBookingDetailRead) TableName() string {
	return "inv.order_booking_detail"
}

type OrderBookingDetailApproval struct {
	CustID         string   `gorm:"cust_id" json:"cust_id"`
	OrderBookingId int      `gorm:"order_booking_id" json:"order_booking_id"`
	QtyAlloc       *float64 `gorm:"qty_alloc" json:"qty_alloc"`
	Qty1Alloc      *float64 `gorm:"qty1_alloc" json:"qty1_alloc"`
	Qty2Alloc      *float64 `gorm:"qty2_alloc" json:"qty2_alloc"`
	Qty3Alloc      *float64 `gorm:"qty3_alloc" json:"qty3_alloc"`
	Qty4Alloc      *float64 `gorm:"qty4_alloc" json:"qty4_alloc"`
	Qty5Alloc      *float64 `gorm:"qty5_alloc" json:"qty5_alloc"`
	AmountAlloc    *float64 `gorm:"amount_alloc" json:"amount_alloc"`
	VatValueAlloc  *float64 `gorm:"vat_value_alloc" json:"vat_value_alloc"`
}

func (OrderBookingDetailApproval) TableName() string {
	return "inv.order_booking_detail"
}

type OrderBookingDetailStatus struct {
	CustID             string     `gorm:"cust_id" json:"cust_id"`
	StatusOrderBooking *int64     `gorm:"status_order_booking" json:"status_order_booking"`
	UpdatedBy          *int64     `gorm:"updated_by" json:"updated_by"`
	UpdatedAt          *time.Time `gorm:"updated_at" json:"updated_at"`
}

func (OrderBookingDetailStatus) TableName() string {
	return "inv.order_booking"
}

type OrderBookingDetailStatusApproval struct {
	CustID             string  `gorm:"cust_id" json:"cust_id"`
	StatusOrderBooking int64   `gorm:"status_order_booking" json:"status_order_booking"`
	SubTotalAlloc      float64 `gorm:"sub_total_alloc" json:"sub_total_alloc"`
	VatValueAlloc      float64 `gorm:"vat_value_alloc" json:"vat_value_alloc"`
	TotalAlloc         float64 `gorm:"total_alloc" json:"total_alloc"`
	DeliveryFee        float64 `gorm:"delivery_fee" json:"delivery_fee"`
	TypeApproval       int64   `gorm:"type_approval" json:"type_approval"`
	// UpdatedBy          int64     `gorm:"updated_by" json:"updated_by"`
	// UpdatedAt          time.Time `gorm:"updated_at" json:"updated_at"`
}

func (OrderBookingDetailStatusApproval) TableName() string {
	return "inv.order_booking"
}
