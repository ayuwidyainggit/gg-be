package model

import (
	"strings"
	"time"

	"gorm.io/gorm"
)

type ReplenishmentOrder struct {
	CustID            string     `gorm:"column:cust_id;primaryKey" json:"cust_id"`
	ParentCustID      string     `gorm:"-" json:"-"`
	ReplenishmentID   int64      `gorm:"column:replenishment_id;primaryKey;autoIncrement" json:"replenishment_id"`
	ReplenishmentNo   string     `gorm:"column:replenishment_no" json:"replenishment_no"`
	Date              time.Time  `gorm:"column:date" json:"date"`
	DistributorID     *int64     `gorm:"column:distributor_id" json:"distributor_id"`
	SupID             int64      `gorm:"column:sup_id" json:"sup_id"`
	WhID              int64      `gorm:"column:wh_id" json:"wh_id"`
	DeliveryType      string     `gorm:"column:delivery_type" json:"delivery_type"`
	ReplenishmentType string     `gorm:"column:replenishment_type" json:"replenishment_type"`
	SoStartDate       *time.Time `gorm:"column:so_start_date" json:"so_start_date"`
	SoEndDate         *time.Time `gorm:"column:so_end_date" json:"so_end_date"`
	DeliveryDate      *time.Time `gorm:"column:delivery_date" json:"delivery_date"`
	Note              string     `gorm:"column:note" json:"note"`
	Status            int        `gorm:"column:status" json:"status"`
	IsApproval        *bool      `gorm:"column:is_approval" json:"is_approval"`
	ApproveBy         *int64     `gorm:"column:approve_by" json:"approve_by"`
	ApproveAt         *time.Time `gorm:"column:approve_at" json:"approve_at"`
	SoNo              *string    `gorm:"column:so_no" json:"so_no"`
	IsAdditionFrom    bool       `gorm:"column:is_addition_from" json:"is_addition_from"`
	CreatedBy         int64      `gorm:"column:created_by" json:"created_by"`
	CreatedAt         time.Time  `gorm:"column:created_at" json:"created_at"`
	UpdatedBy         *int64     `gorm:"column:updated_by" json:"updated_by"`
	UpdatedAt         *time.Time `gorm:"column:updated_at" json:"updated_at"`
	DeletedBy         *int64     `gorm:"column:deleted_by" json:"deleted_by"`
	DeletedAt         *time.Time `gorm:"column:deleted_at" json:"deleted_at"`
	IsDel             bool       `gorm:"column:is_del" json:"is_del"`
}

func (ReplenishmentOrder) TableName() string {
	return "inv.replenishment_order"
}

func (m *ReplenishmentOrder) BeforeCreate(trx *gorm.DB) (err error) {
	var numberDoc NumberDoc
	trCode := "ARO"
	dateSubtr := m.Date.Format("060102")
	custID := strings.ToUpper(strings.TrimSpace(m.CustID))
	var numberingScopeCustID string
	err = trx.Raw(`
		SELECT CASE
			WHEN EXISTS (
				SELECT 1 FROM mst.m_principal mp
				WHERE UPPER(TRIM(mp.cust_id)) = ?
				  AND mp.is_del = false
			) THEN ?
			ELSE UPPER(TRIM(COALESCE(NULLIF(TRIM(mc.parent_cust_id), ''), mc.cust_id)))
		END
		FROM smc.m_customer mc
		WHERE UPPER(TRIM(mc.cust_id)) = ?
		LIMIT 1
	`, custID, custID, custID).Scan(&numberingScopeCustID).Error
	if err != nil || strings.TrimSpace(numberingScopeCustID) == "" {
		numberingScopeCustID = custID
	}
	numberingScopeCustID = strings.ToUpper(strings.TrimSpace(numberingScopeCustID))

	queryStr := `
		SELECT TRIM(
			to_char(
				COALESCE(MAX(TO_NUMBER(SUBSTR(replenishment_no, 10, 3), '999')), 0) + 1,
				'000'
			)
		) AS next_seq
		FROM inv.replenishment_order
		WHERE substr(replenishment_no, 4, 6) = ?
		  AND (
			UPPER(TRIM(cust_id)) = ?
			OR UPPER(TRIM(cust_id)) IN (
				SELECT UPPER(TRIM(c.cust_id))
				FROM smc.m_customer c
				WHERE UPPER(TRIM(c.parent_cust_id)) = ?
			)
		  )
	`

	err = trx.Raw(queryStr, dateSubtr, numberingScopeCustID, numberingScopeCustID).Scan(&numberDoc).Error
	if err != nil {
		return err
	}

	m.ReplenishmentNo = trCode + dateSubtr + numberDoc.NextSeq
	m.CreatedAt = time.Now()
	now := time.Now()
	m.UpdatedAt = &now
	if m.CreatedBy != 0 {
		m.UpdatedBy = &m.CreatedBy
	}
	// Set default status to 1 (Need Review)
	if m.Status == 0 {
		m.Status = 1
	}
	return nil
}

func (m *ReplenishmentOrder) BeforeUpdate(trx *gorm.DB) (err error) {
	now := time.Now()
	m.UpdatedAt = &now
	return nil
}

type NumberDoc struct {
	NextSeq string `gorm:"column:next_seq"`
}

type ReplenishmentOrderList struct {
	CustID          string     `gorm:"column:cust_id" json:"cust_id"`
	ReplenishmentID int64      `gorm:"column:replenishment_id" json:"replenishment_id"`
	ReplenishmentNo string     `gorm:"column:replenishment_no" json:"replenishment_no"`
	Date            *time.Time `gorm:"column:date" json:"date"`
	DeliveryDate    *time.Time `gorm:"column:delivery_date" json:"delivery_date"`
	SoNo            *string    `gorm:"column:so_no" json:"so_no"`
	SupID           *int64     `gorm:"column:sup_id" json:"sup_id"`
	SupCode         *string    `gorm:"column:sup_code" json:"sup_code"`
	SupName         *string    `gorm:"column:sup_name" json:"sup_name"`
	WhID            *int64     `gorm:"column:wh_id" json:"wh_id"`
	WhCode          *string    `gorm:"column:wh_code" json:"wh_code"`
	WhName          *string    `gorm:"column:wh_name" json:"wh_name"`
	Status          int        `gorm:"column:status" json:"status"`
	StatusName      *string    `gorm:"column:status_name" json:"status_name"`
	DistributorID   *int64     `gorm:"column:distributor_id" json:"distributor_id"`
	DistributorCode *string    `gorm:"column:distributor_code" json:"distributor_code"`
	DistributorName *string    `gorm:"column:distributor_name" json:"distributor_name"`
	Address         *string    `gorm:"column:address" json:"address"`
	CreatedBy       int64      `gorm:"column:created_by" json:"created_by"`
	CreatedByName   *string    `gorm:"column:created_by_name" json:"created_by_name"`
	UpdatedBy       *int64     `gorm:"column:updated_by" json:"updated_by"`
	UpdatedByName   *string    `gorm:"column:updated_by_name" json:"updated_by_name"`
	CreatedAt       time.Time  `gorm:"column:created_at" json:"created_at"`
	UpdatedAt       *time.Time `gorm:"column:updated_at" json:"updated_at"`
	IsDel           bool       `gorm:"column:is_del" json:"is_del"`
}

func (ReplenishmentOrderList) TableName() string {
	return "inv.replenishment_order"
}

type ReplenishmentOrderDetail struct {
	CustID                string     `gorm:"column:cust_id;primaryKey" json:"cust_id"`
	ReplenishmentDetailID int64      `gorm:"column:replenishment_detail_id;primaryKey;autoIncrement" json:"replenishment_detail_id"`
	ReplenishmentID       int64      `gorm:"column:replenishment_id" json:"replenishment_id"`
	ProID                 int64      `gorm:"column:pro_id" json:"pro_id"`
	ReturnReasonID        *int64     `gorm:"column:return_reason_id" json:"return_reason_id"`
	OrderBookingQty1      float64    `gorm:"column:order_booking_qty1" json:"order_booking_qty1"`
	OrderBookingQty2      float64    `gorm:"column:order_booking_qty2" json:"order_booking_qty2"`
	OrderBookingQty3      float64    `gorm:"column:order_booking_qty3" json:"order_booking_qty3"`
	QtyOrderApproval1     *float64   `gorm:"column:qty_order_approval1" json:"qty_order_approval1"`
	QtyOrderApproval2     *float64   `gorm:"column:qty_order_approval2" json:"qty_order_approval2"`
	QtyOrderApproval3     *float64   `gorm:"column:qty_order_approval3" json:"qty_order_approval3"`
	QtyOrderAllocation1   *float64   `gorm:"column:qty_order_allocation1" json:"qty_order_allocation1"`
	QtyOrderAllocation2   *float64   `gorm:"column:qty_order_allocation2" json:"qty_order_allocation2"`
	QtyOrderAllocation3   *float64   `gorm:"column:qty_order_allocation3" json:"qty_order_allocation3"`
	PurchPrice1           float64    `gorm:"column:purch_price1" json:"purch_price1"`
	PurchPrice2           float64    `gorm:"column:purch_price2" json:"purch_price2"`
	PurchPrice3           float64    `gorm:"column:purch_price3" json:"purch_price3"`
	SapQty3               *float64   `gorm:"column:sap_qty3" json:"sap_qty3,omitempty"`
	SapPurchPrice3        *float64   `gorm:"column:sap_purch_price3" json:"sap_purch_price3,omitempty"`
	EstimatedPrice        float64    `gorm:"column:estimated_price" json:"estimated_price"`
	CreatedBy             int64      `gorm:"column:created_by" json:"created_by"`
	CreatedAt             time.Time  `gorm:"column:created_at" json:"created_at"`
	UpdatedBy             *int64     `gorm:"column:updated_by" json:"updated_by"`
	UpdatedAt             *time.Time `gorm:"column:updated_at" json:"updated_at"`
	DeletedBy             *int64     `gorm:"column:deleted_by" json:"deleted_by"`
	DeletedAt             *time.Time `gorm:"column:deleted_at" json:"deleted_at"`
	IsDel                 bool       `gorm:"column:is_del" json:"is_del"`
}

func (ReplenishmentOrderDetail) TableName() string {
	return "inv.replenishment_order_detail"
}

func (m *ReplenishmentOrderDetail) BeforeCreate(trx *gorm.DB) (err error) {
	m.CreatedAt = time.Now()
	now := time.Now()
	m.UpdatedAt = &now
	if m.CreatedBy != 0 {
		m.UpdatedBy = &m.CreatedBy
	}
	return nil
}

func (m *ReplenishmentOrderDetail) BeforeUpdate(trx *gorm.DB) (err error) {
	now := time.Now()
	m.UpdatedAt = &now
	return nil
}

// Model for reading detail with joins
type ReplenishmentOrderDetailRead struct {
	ReplenishmentDetailID int64    `gorm:"column:replenishment_detail_id" json:"replenishment_detail_id"`
	ReplenishmentID       int64    `gorm:"column:replenishment_id" json:"replenishment_id"`
	ReplenishmentNo       string   `gorm:"column:replenishment_no" json:"replenishment_no"`
	ProID                 int64    `gorm:"column:pro_id" json:"pro_id"`
	ProCode               *string  `gorm:"column:pro_code" json:"pro_code"`
	ProName               *string  `gorm:"column:pro_name" json:"pro_name"`
	UnitID1               *string  `gorm:"column:unit_id1" json:"unit_id1"`
	UnitID2               *string  `gorm:"column:unit_id2" json:"unit_id2"`
	UnitID3               *string  `gorm:"column:unit_id3" json:"unit_id3"`
	OrderBookingQty1      float64  `gorm:"column:order_booking_qty1" json:"order_booking_qty1"`
	OrderBookingQty2      float64  `gorm:"column:order_booking_qty2" json:"order_booking_qty2"`
	OrderBookingQty3      float64  `gorm:"column:order_booking_qty3" json:"order_booking_qty3"`
	PurchPrice1           float64  `gorm:"column:purch_price1" json:"purch_price1"`
	PurchPrice2           float64  `gorm:"column:purch_price2" json:"purch_price2"`
	PurchPrice3           float64  `gorm:"column:purch_price3" json:"purch_price3"`
	EstimatedPrice        float64  `gorm:"column:estimated_price" json:"estimated_price"`
	Vat                   *float64 `gorm:"column:vat" json:"vat"`
	QtyOrderAllocation1   *float64 `gorm:"column:qty_order_allocation1" json:"qty_order_allocation1"`
	QtyOrderAllocation2   *float64 `gorm:"column:qty_order_allocation2" json:"qty_order_allocation2"`
	QtyOrderAllocation3   *float64 `gorm:"column:qty_order_allocation3" json:"qty_order_allocation3"`
	QtyOrderApproval1     *float64 `gorm:"column:qty_order_approval1" json:"qty_order_approval1"`
	QtyOrderApproval2     *float64 `gorm:"column:qty_order_approval2" json:"qty_order_approval2"`
	QtyOrderApproval3     *float64 `gorm:"column:qty_order_approval3" json:"qty_order_approval3"`
	QtyFinal1             *float64 `gorm:"column:qty_final1" json:"qty_final1"`
	QtyFinal2             *float64 `gorm:"column:qty_final2" json:"qty_final2"`
	QtyFinal3             *float64 `gorm:"column:qty_final3" json:"qty_final3"`
	StockReceived1        *float64 `gorm:"column:stock_received1" json:"stock_received1"`
	StockReceived2        *float64 `gorm:"column:stock_received2" json:"stock_received2"`
	StockReceived3        *float64 `gorm:"column:stock_received3" json:"stock_received3"`
	// Product master fields
	SafStockQty *float64 `gorm:"column:saf_stock_qty" json:"saf_stock_qty"`
	MinStockQty *float64 `gorm:"column:min_stock_qty" json:"min_stock_qty"`
	// Warehouse stock qty (calculated)
	Qty1 *float64 `gorm:"column:qty1" json:"qty1"`
	Qty2 *float64 `gorm:"column:qty2" json:"qty2"`
	Qty3 *float64 `gorm:"column:qty3" json:"qty3"`
	// In transit stock (sum of order_booking_qty from replenishment orders with status On Delivery)
	InTransitStock1 float64 `gorm:"column:in_transit_stock1" json:"in_transit_stock1"`
	InTransitStock2 float64 `gorm:"column:in_transit_stock2" json:"in_transit_stock2"`
	InTransitStock3 float64 `gorm:"column:in_transit_stock3" json:"in_transit_stock3"`
}

func (ReplenishmentOrderDetailRead) TableName() string {
	return "inv.replenishment_order_detail"
}

// Model for reading header with supplier and warehouse info
type ReplenishmentOrderRead struct {
	CustID            string     `gorm:"column:cust_id" json:"cust_id"`
	ReplenishmentID   int64      `gorm:"column:replenishment_id" json:"replenishment_id"`
	ReplenishmentNo   string     `gorm:"column:replenishment_no" json:"replenishment_no"`
	Date              time.Time  `gorm:"column:date" json:"date"`
	SupID             int64      `gorm:"column:sup_id" json:"sup_id"`
	SupCode           *string    `gorm:"column:sup_code" json:"sup_code"`
	SupName           *string    `gorm:"column:sup_name" json:"sup_name"`
	WhID              int64      `gorm:"column:wh_id" json:"wh_id"`
	WhCode            *string    `gorm:"column:wh_code" json:"wh_code"`
	WhName            *string    `gorm:"column:wh_name" json:"wh_name"`
	DeliveryType      string     `gorm:"column:delivery_type" json:"delivery_type"`
	ReplenishmentType string     `gorm:"column:replenishment_type" json:"replenishment_type"`
	SoStartDate       *time.Time `gorm:"column:so_start_date" json:"so_start_date"`
	SoEndDate         *time.Time `gorm:"column:so_end_date" json:"so_end_date"`
	DeliveryDate      *time.Time `gorm:"column:delivery_date" json:"delivery_date"`
	Note              string     `gorm:"column:note" json:"note"`
	Status            int        `gorm:"column:status" json:"status"`
	StatusName        *string    `gorm:"column:status_name" json:"status_name"`
	SoNo              *string    `gorm:"column:so_no" json:"so_no"`
	CreatedBy         int64      `gorm:"column:created_by" json:"created_by"`
	CreatedAt         time.Time  `gorm:"column:created_at" json:"created_at"`
	UpdatedBy         *int64     `gorm:"column:updated_by" json:"updated_by"`
	UpdatedAt         *time.Time `gorm:"column:updated_at" json:"updated_at"`
	DeletedBy         *int64     `gorm:"column:deleted_by" json:"deleted_by"`
	DeletedAt         *time.Time `gorm:"column:deleted_at" json:"deleted_at"`
	IsDel             bool       `gorm:"column:is_del" json:"is_del"`
	// Distributor fields
	DistributorID   *int64  `gorm:"column:distributor_id" json:"distributor_id"`
	DistributorCode *string `gorm:"column:distributor_code" json:"distributor_code"`
	DistributorName *string `gorm:"column:distributor_name" json:"distributor_name"`
	Address         *string `gorm:"column:address" json:"address"`
	// Delivery fee from GR
	DeliveryFee *float64 `gorm:"column:delivery_fee" json:"delivery_fee"`
}

func (ReplenishmentOrderRead) TableName() string {
	return "inv.replenishment_order"
}

type ReplenishmentFinalRead struct {
	ReplenishmentDetailID int64    `gorm:"column:replenishment_detail_id" json:"replenishment_detail_id"`
	ProID                 int64    `gorm:"column:pro_id" json:"pro_id"`
	ProCode               *string  `gorm:"column:pro_code" json:"pro_code"`
	ProName               *string  `gorm:"column:pro_name" json:"pro_name"`
	UnitID1               *string  `gorm:"column:unit_id1" json:"unit_id1"`
	UnitID2               *string  `gorm:"column:unit_id2" json:"unit_id2"`
	UnitID3               *string  `gorm:"column:unit_id3" json:"unit_id3"`
	PurchPriceDelivery1   float64  `gorm:"column:purch_price_delivery1" json:"purch_price_delivery1"`
	PurchPriceDelivery2   float64  `gorm:"column:purch_price_delivery2" json:"purch_price_delivery2"`
	PurchPriceDelivery3   float64  `gorm:"column:purch_price_delivery3" json:"purch_price_delivery3"`
	FinalOrder1           *float64 `gorm:"column:final_order1" json:"final_order1"`
	FinalOrder2           *float64 `gorm:"column:final_order2" json:"final_order2"`
	FinalOrder3           *float64 `gorm:"column:final_order3" json:"final_order3"`
	GrPrice1              *float64 `gorm:"column:gr_price1" json:"gr_price1"`
	GrPrice2              *float64 `gorm:"column:gr_price2" json:"gr_price2"`
	GrPrice3              *float64 `gorm:"column:gr_price3" json:"gr_price3"`
	StockReceived1        *float64 `gorm:"column:stock_received1" json:"stock_received1"`
	StockReceived2        *float64 `gorm:"column:stock_received2" json:"stock_received2"`
	StockReceived3        *float64 `gorm:"column:stock_received3" json:"stock_received3"`
	Vat                   *float64 `gorm:"column:vat" json:"vat"`
}

type ReplenishmentGoodReceiptRead struct {
	PoNo           string   `gorm:"column:po_no" json:"po_no"`
	ProID          int64    `gorm:"column:pro_id" json:"pro_id"`
	ProCode        *string  `gorm:"column:pro_code" json:"pro_code"`
	ProName        *string  `gorm:"column:pro_name" json:"pro_name"`
	UnitPrice1     float64  `gorm:"column:unit_price1" json:"unit_price1"`
	UnitPrice2     float64  `gorm:"column:unit_price2" json:"unit_price2"`
	UnitPrice3     float64  `gorm:"column:unit_price3" json:"unit_price3"`
	QtyReceived1   float64  `gorm:"column:qty_received1" json:"qty_received1"`
	QtyReceived2   float64  `gorm:"column:qty_received2" json:"qty_received2"`
	QtyReceived3   float64  `gorm:"column:qty_received3" json:"qty_received3"`
	EstimatedPrice float64  `gorm:"column:estimated_price" json:"estimated_price"`
	Vat            *float64 `gorm:"column:vat" json:"vat"`
}

func (ReplenishmentFinalRead) TableName() string {
	return "inv.replenishment_order_detail"
}

type ReplenishmentProductList struct {
	ProID           int64   `gorm:"column:pro_id" json:"pro_id"`
	ProCode         string  `gorm:"column:pro_code" json:"pro_code"`
	ProName         string  `gorm:"column:pro_name" json:"pro_name"`
	PurchPrice1     float64 `gorm:"column:purch_price1" json:"purch_price1"`
	PurchPrice2     float64 `gorm:"column:purch_price2" json:"purch_price2"`
	PurchPrice3     float64 `gorm:"column:purch_price3" json:"purch_price3"`
	UnitID1         string  `gorm:"column:unit_id1" json:"unit_id1"`
	UnitID2         string  `gorm:"column:unit_id2" json:"unit_id2"`
	UnitID3         string  `gorm:"column:unit_id3" json:"unit_id3"`
	Vat             float64 `gorm:"column:vat" json:"vat"`
	Qty1            float64 `gorm:"column:qty1" json:"qty1"`
	Qty2            float64 `gorm:"column:qty2" json:"qty2"`
	Qty3            float64 `gorm:"column:qty3" json:"qty3"`
	InTransitStock1 float64 `gorm:"column:in_transit_stock1" json:"in_transit_stock1"`
	InTransitStock2 float64 `gorm:"column:in_transit_stock2" json:"in_transit_stock2"`
	InTransitStock3 float64 `gorm:"column:in_transit_stock3" json:"in_transit_stock3"`
}

// PO List model
type PoList struct {
	ReplenishmentNo   string  `gorm:"column:replenishment_no" json:"replenishment_no"`
	ReplenishmentType string  `gorm:"column:replenishment_type" json:"replenishment_type"`
	WhID              int64   `gorm:"column:wh_id" json:"wh_id"`
	WhCode            *string `gorm:"column:wh_code" json:"wh_code"`
	WhName            *string `gorm:"column:wh_name" json:"wh_name"`
	SupID             int64   `gorm:"column:sup_id" json:"sup_id"`
	SupCode           *string `gorm:"column:sup_code" json:"sup_code"`
	SupName           *string `gorm:"column:sup_name" json:"sup_name"`
}

func (PoList) TableName() string {
	return "inv.replenishment_order"
}

// Replenishment Approval Product model
type ReplenishmentApprovalProduct struct {
	ProID           int64    `gorm:"column:pro_id" json:"pro_id"`
	ProCode         string   `gorm:"column:pro_code" json:"pro_code"`
	ProName         string   `gorm:"column:pro_name" json:"pro_name"`
	Ripening        *float64 `gorm:"column:ripening" json:"ripening"`
	InTransitStock1 float64  `gorm:"column:in_transit_stock1" json:"in_transit_stock1"`
	InTransitStock2 float64  `gorm:"column:in_transit_stock2" json:"in_transit_stock2"`
	InTransitStock3 float64  `gorm:"column:in_transit_stock3" json:"in_transit_stock3"`
	SafStockQty     *float64 `gorm:"column:saf_stock_qty" json:"saf_stock_qty"`
	MinStockQty     *float64 `gorm:"column:min_stock_qty" json:"min_stock_qty"`
	Vat             *float64 `gorm:"column:vat" json:"vat"`
	ConvUnit2       *float64 `gorm:"column:conv_unit2" json:"conv_unit2"`
	ConvUnit3       *float64 `gorm:"column:conv_unit3" json:"conv_unit3"`
	UnitID1         *string  `gorm:"column:unit_id1" json:"unit_id1"`
	UnitID2         *string  `gorm:"column:unit_id2" json:"unit_id2"`
	UnitID3         *string  `gorm:"column:unit_id3" json:"unit_id3"`
	PurchPrice1     *float64 `gorm:"column:purch_price1" json:"purch_price1"`
	PurchPrice2     *float64 `gorm:"column:purch_price2" json:"purch_price2"`
	PurchPrice3     *float64 `gorm:"column:purch_price3" json:"purch_price3"`
	Qty1            *float64 `gorm:"column:qty1" json:"qty1"`
	Qty2            *float64 `gorm:"column:qty2" json:"qty2"`
	Qty3            *float64 `gorm:"column:qty3" json:"qty3"`
	TotalQty        *float64 `gorm:"column:total_qty" json:"total_qty"`
}

func (ReplenishmentApprovalProduct) TableName() string {
	return "mst.m_product"
}

// Product GR List model
type ProductGrListDetail struct {
	ReplenishmentDetailID int64    `gorm:"column:replenishment_detail_id" json:"replenishment_detail_id"`
	ProID                 int64    `gorm:"column:pro_id" json:"pro_id"`
	ProCode               *string  `gorm:"column:pro_code" json:"pro_code"`
	ProName               *string  `gorm:"column:pro_name" json:"pro_name"`
	Vat                   *float64 `gorm:"column:vat" json:"vat"`
	QtyOrderApproval1     *float64 `gorm:"column:qty_order_approval1" json:"qty_order_approval1"`
	QtyOrderApproval2     *float64 `gorm:"column:qty_order_approval2" json:"qty_order_approval2"`
	QtyOrderApproval3     *float64 `gorm:"column:qty_order_approval3" json:"qty_order_approval3"`
	UnitID1               *string  `gorm:"column:unit_id1" json:"unit_id1"`
	UnitID2               *string  `gorm:"column:unit_id2" json:"unit_id2"`
	UnitID3               *string  `gorm:"column:unit_id3" json:"unit_id3"`
	PurchPrice1           float64  `gorm:"column:purch_price1" json:"purch_price1"`
	PurchPrice2           float64  `gorm:"column:purch_price2" json:"purch_price2"`
	PurchPrice3           float64  `gorm:"column:purch_price3" json:"purch_price3"`
}

func (ProductGrListDetail) TableName() string {
	return "inv.replenishment_order_detail"
}

type SummarizeReplanishmentRow struct {
	ReplanishmentID       int64    `gorm:"column:replanishment_id"`
	ReplanishmentNo       string   `gorm:"column:replanishment_no"`
	DisributorID          int64    `gorm:"column:disributor_id"`
	DistributorCode       *string  `gorm:"column:distributor_code"`
	DistributorName       *string  `gorm:"column:distributor_name"`
	SupID                 int64    `gorm:"column:sup_id"`
	SupCode               *string  `gorm:"column:sup_code"`
	SupName               *string  `gorm:"column:sup_name"`
	WhID                  int64    `gorm:"column:wh_id"`
	WhCode                *string  `gorm:"column:wh_code"`
	WhName                *string  `gorm:"column:wh_name"`
	ReplanishmentDetailID int64    `gorm:"column:replanishment_detail_id"`
	ProID                 int64    `gorm:"column:pro_id"`
	ProCode               *string  `gorm:"column:pro_code"`
	ProName               *string  `gorm:"column:pro_name"`
	WhStockLarge          *float64 `gorm:"column:wh_stock_large"`
	WhStockMedium         *float64 `gorm:"column:wh_stock_medium"`
	WhStockSmall          *float64 `gorm:"column:wh_stock_small"`
	OptimumQty            *float64 `gorm:"column:optimum_qty"`
	Ripening              *float64 `gorm:"column:ripening"`
	ReturnReasonID        *int64   `gorm:"column:return_reason_id"`
	UnitID1               *string  `gorm:"column:unit_id1"`
	UnitID2               *string  `gorm:"column:unit_id2"`
	UnitID3               *string  `gorm:"column:unit_id3"`
	PurchPrice1           float64  `gorm:"column:purch_price1"`
	PurchPrice2           float64  `gorm:"column:purch_price2"`
	PurchPrice3           float64  `gorm:"column:purch_price3"`
	QtyRo1                float64  `gorm:"column:qty_ro1"`
	QtyRo2                float64  `gorm:"column:qty_ro2"`
	QtyRo3                float64  `gorm:"column:qty_ro3"`
	EstimatedPrice        float64  `gorm:"column:estimated_price"`
}
