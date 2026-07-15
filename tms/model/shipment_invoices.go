package model

import "time"

type ShipmentInvoices struct {
	ID                 int        `json:"id"`
	ShipmentNo         *string    `json:"shipment_no" gorm:"type:varchar(125)"`
	OrderNo            *string    `json:"order_no"`
	InvoiceNo          *string    `json:"invoice_no"`
	OutletID           int        `json:"outlet_id"`
	OutletCode         string     `json:"outlet_code"`
	OutletName         string     `json:"outlet_name"`
	OutletAddress      string     `json:"outlet_address"`
	OutletStatus       string     `json:"outlet_status"`
	SalesmanID         int        `json:"salesman_id"`
	SalesmanName       string     `json:"salesman_name"`
	ProductID          int        `json:"product_id"`
	ProductName        string     `json:"product_name"`
	ProductStatus      string     `json:"product_status"`
	ProductCode        string     `json:"product_code"`
	DeliveryDate       time.Time  `json:"delivery_date"`
	Sku                string     `json:"sku"`
	Volume             float64    `json:"volume"`
	Weight             float64    `json:"weight"`
	CustID             string     `json:"cust_id"`
	RouteID            *int       `json:"route_id"'`
	ArriveAt           *int64     `json:"arrive_at"`
	UnloadAt           *int64     `json:"unload_at"`
	LeaveAt            *int64     `json:"leave_at"`
	SkipAt             *int64     `json:"skip_at"`
	PickupAt           *int64     `json:"pickup_at"`
	ReasonID           *int       `json:"reason_id"`
	ReasonName         *string    `json:"reason_name"`
	Qty1               int64      `json:"qty1"`
	Qty2               int64      `json:"qty2"`
	Qty3               int64      `json:"qty3"`
	QtyReject1         *int64     `json:"qty_reject_1" gorm:"column:qty_reject_1"`
	QtyReject2         *int64     `json:"qty_reject_2" gorm:"column:qty_reject_2"`
	QtyReject3         *int64     `json:"qty_reject_3" gorm:"column:qty_reject_3"`
	ConvUnit1          int64      `json:"conv_unit1"`
	ConvUnit2          int64      `json:"conv_unit2"`
	ConvUnit3          int64      `json:"conv_unit3"`
	UnitId1            string     `json:"unit_id1"`
	UnitId2            string     `json:"unit_id2"`
	UnitId3            string     `json:"unit_id3"`
	WarehouseLatitude  string     `json:"warehouse_latitude"`
	WarehouseLongitude string     `json:"warehouse_longitude"`
	OutletLatitude     string     `json:"outlet_latitude"`
	OutletLongitude    string     `json:"outlet_longitude"`
	Status             string     `json:"status"`
	SkipReason         string     `json:"skip_reason"`
	InOutlet           bool       `json:"in_outlet"`
	Photo              string     `json:"photo"`
	Signature          string     `json:"signature"`
	TotalBruto         *float64   `json:"total_bruto"`
	TotalVolumeDisc    *float64   `json:"total_volume_disc"`
	TotalPromo         *float64   `json:"total_promo"`
	TotalPpn           *float64   `json:"total_ppn"`
	DueDate            *time.Time `json:"due_date"`
	SellPrice1         *float64   `json:"sell_price1"`
	SellPrice2         *float64   `json:"sell_price2"`
	SellPrice3         *float64   `json:"sell_price3"`
	PayTypeName        string     `json:"pay_type_name"`
	InvoiceDate        *time.Time `json:"invoice_date"`
	ItemCdnName        *string    `json:"item_cdn_name"`
	TotalBelumBayar    *float64   `json:"total_belum_bayar"`
	TotalNetto         *float64   `json:"total_netto"`
	Vat                float64    `json:"vat"`
	VatValue           float64    `json:"vat_value"`
	OrderDetailID      int        `json:"order_detail_id"`
	OnHold             *int64     `json:"on_hold"`
	ResumeAt           *int64     `json:"resume_at"`
	CreatedAt          time.Time  `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt          time.Time  `gorm:"column:updated_at;autoUpdateTime"`

	// alias column
	//DriverId int `gorm:"->" json:"-"`

	Shipment *Shipment `gorm:"foreignKey:shipment_no;references:ShipmentNo"`
}

func (ShipmentInvoices) TableName() string {
	return "tms.shipment_invoices"
}
