package entity

import "time"

type ShipmentResponse struct {
	ShipmentID   int    `json:"id" example:"1"`
	ShipmentNo   string `json:"shipment_no" example:"11024066062"`
	DriverID     int    `json:"driver_id" example:"1"`
	VehicleID    int    `json:"vehicle_id" example:"1"`
	DriverName   string `json:"driver_name" example:"John Driver"`
	VehicleName  string `json:"vehicle_name" example:"Garuda"`
	VehicleType  string `json:"vehicle_type" example:"5"`
	DeliveryDate string `json:"delivery_date" example:"2024-06-11"`
	ShipmentType string `json:"shipment_type" example:"manual"`
	CreatedAt    string `json:"created_at" example:"2024-06-11"`
}

type ShipmentPreviewResponse struct {
	ShipmentNo   string    `json:"shipment_no" example:"11024066062"`
	VehicleNo    string    `json:"vehicle_no" example:"B 323 R"`
	VehicleType  string    `json:"vehicle_type" example:"5"`
	VehicleName  string    `json:"vehicle_name" example:"Garuda"`
	DriverID     int       `json:"driver_id" example:"1"`
	DriverName   string    `json:"driver_name" example:"John Driver"`
	HelperID     int       `json:"helper_id" example:"2"`
	HelperName   string    `json:"helper_name" example:"George Helper"`
	Start        int       `json:"start" example:"1722912443"`
	SalesmanName string    `json:"salesman_name" example:"John Sales"`
	PickListDate time.Time `json:"pick_list_date"`
	PickListNo   string    `json:"pick_list_no"`
	ShipmentInvoiceSummary
	Shipment []ShipmentEntryResponse `json:"shipment"`
}

type ShipmentEntryResponse struct {
	OutletID int                    `json:"outlet_id" example:"34"`
	Outlets  []ShipmentItemResponse `json:"outlets"`
}

type ShipmentItemResponse struct {
	OutletID       int     `json:"outlet_id" example:"34"`
	OutletCode     string  `json:"outlet_code" example:"001001"`
	OutletName     string  `json:"outlet_name" example:"Lapan Belas"`
	OutletStatus   string  `json:"outlet_status" example:"On Progress"`
	ShipmentStatus string  `json:"shipment_status" example:"Delivery"`
	SalesmanName   string  `json:"salesman_name" example:"John Sales"`
	OrderNo        string  `json:"order_no" example:"SO2409100004"`
	InvoiceNo      string  `json:"invoices_no" example:"406438662776815"`
	Sku            string  `json:"sku" example:"0001"`
	ProductID      int     `json:"product_id" example:"1"`
	ProductCode    string  `json:"product_code" example:"PRO5173"`
	DueDate        string  `json:"due_date" example:"2024-12-12"`
	InvoiceDate    string  `json:"invoice_date" example:"2024-12-12"`
	ProductName    string  `json:"product_name" example:"ASSORTED BISCUIT OPP RED 20 275G"`
	ProductStatus  string  `json:"product_status" example:"-"`
	Qty1           int64   `json:"qty1" example:"10"`
	Qty2           int64   `json:"qty2" example:"10"`
	Qty3           int64   `json:"qty3" example:"10"`
	ConvUnit1      int64   `json:"conv_unit1" example:"0"`
	ConvUnit2      int64   `json:"conv_unit2" example:"6"`
	ConvUnit3      int64   `json:"conv_unit3" example:"6"`
	UnitId1        string  `json:"unit_id1" example:"PCS"`
	UnitId2        string  `json:"unit_id2" example:"CTN"`
	UnitId3        string  `json:"unit_id3" example:"CTN"`
	Volume         float64 `json:"volume" example:"0.04"`
	Weight         float64 `json:"weight" example:"0.3578877"`
	SellPrice1     float64 `json:"sell_price1"`
	SellPrice2     float64 `json:"sell_price2"`
	SellPrice3     float64 `json:"sell_price3"`
	ItemCdnName    string  `json:"item_cdn_name"`
	OrderDetailID  int     `json:"order_detail_id"`
}

type ShipmentInvoiceSummary struct {
	TotalBruto      float64 `json:"total_bruto"`
	TotalVolumeDisc float64 `json:"total_volume_disc"`
	TotalPromo      float64 `json:"total_promo"`
	TotalPpn        float64 `json:"total_ppn"`
	TotalNetto      float64 `json:"total_netto"`
}

type CreateShipmentRequest struct {
	DeliveryDate string         `json:"delivery_date" validate:"required,date"`
	Vehicle      VehicleBody    `json:"vehicles" validate:"required"`
	Shipment     []ShipmentBody `json:"shipments" validate:"required,dive"`
}

type CreateShipmentAutoRequest struct {
	DeliveryDate string         `json:"delivery_date" validate:"required,date"`
	Vehicle      []VehicleBody  `json:"vehicles" validate:"required"`
	Shipment     []ShipmentBody `json:"shipments" validate:"required,dive"`
}

type ShipmentParams struct {
	ShipmentNo string `params:"shipmentNo" validate:"required"`
}

type OrderParams struct {
	OrderNo string `json:"order_no" validate:"required"`
}

type VehicleBody struct {
	VehicleID          int     `json:"vehicle_id" validate:"required"`
	VehicleNo          string  `json:"vehicle_no" validate:"required"`
	VehicleType        string  `json:"vehicle_type" validate:"required"`
	VehicleName        string  `json:"vehicle_name" validate:"required"`
	DriverID           int     `json:"driver_id" validate:"required"`
	DriverName         string  `json:"driver_name" validate:"required"`
	HelperID           int     `json:"helper_id" validate:"required"`
	HelperName         string  `json:"helper_name" validate:"required"`
	Length             float64 `json:"length" validate:"min=0"`
	Width              float64 `json:"width" validate:"min=0"`
	Height             float64 `json:"height" validate:"min=0"`
	Volume             float64 `json:"volume" validate:"min=0"`
	Weight             float64 `json:"weight" validate:"min=0"`
	CustID             string  `json:"cust_id" validate:"required"`
	WarehouseLatitude  any     `json:"warehouse_latitude" validate:"required"`
	WarehouseLongitude any     `json:"warehouse_longitude" validate:"required"`
}

type ShipmentBody struct {
	OrderNo            string   `json:"order_no" validate:"required"`
	Date               string   `json:"date"`
	InvoiceNo          string   `json:"invoice_no"`
	OutletID           int      `json:"outlet_id" validate:"required"`
	OutletCode         string   `json:"outlet_code" validate:"required"`
	OutletAddress      string   `json:"outlet_address" validate:"required"`
	OutletStatus       string   `json:"outlet_status" validate:"required"`
	OutletName         string   `json:"outlet_name" validate:"required"`
	SalesmanID         int      `json:"salesman_id" validate:"required"`
	SalesmanName       string   `json:"salesman_name" validate:"required"`
	Volume             float64  `json:"volume" validate:"required"`
	ProductID          int      `json:"product_id" validate:"required"`
	ProductName        string   `json:"product_name" validate:"required"`
	ProductStatus      string   `json:"product_status" validate:"required"`
	ProductCode        string   `json:"product_code" validate:"required"`
	Sku                string   `json:"sku" validate:"required"`
	Qty1               int64    `json:"qty1" validate:"min=0"`
	Qty2               int64    `json:"qty2" validate:"min=0"`
	Qty3               int64    `json:"qty3" validate:"min=0"`
	ConvUnit1          int64    `json:"conv_unit1" validate:"min=0"`
	ConvUnit2          int64    `json:"conv_unit2" validate:"min=0"`
	ConvUnit3          int64    `json:"conv_unit3" validate:"min=0"`
	UnitId1            string   `json:"unit_id1" validate:"required"`
	UnitId2            string   `json:"unit_id2" validate:"required"`
	UnitId3            string   `json:"unit_id3" validate:"required"`
	Status             string   `json:"status" validate:"required"`
	DeliveryDate       string   `json:"delivery_date"`
	Weight             float64  `json:"weight" validate:"min=0"`
	CustID             string   `json:"cust_id" validate:"required"`
	WarehouseLatitude  any      `json:"warehouse_latitude" validate:"required"`
	WarehouseLongitude any      `json:"warehouse_longitude" validate:"required"`
	OutletLatitude     any      `json:"outlet_latitude" validate:"required"`
	OutletLongitude    any      `json:"outlet_longitude" validate:"required"`
	TotalBruto         *float64 `json:"total_bruto"`
	TotalVolumeDisc    *float64 `json:"total_volume_disc"`
	TotalPromo         *float64 `json:"total_promo"`
	TotalPpn           *float64 `json:"total_ppn"`
	DueDate            string   `json:"due_date"`
	SellPrice1         *float64 `json:"sell_price1"`
	SellPrice2         *float64 `json:"sell_price2"`
	SellPrice3         *float64 `json:"sell_price3"`
	PayTypeName        string   `json:"pay_type_name"`
	TotalNetto         *float64 `json:"total_netto"`
	Vat                float64  `json:"vat" validate:"omitempty,required"`
	VatValue           float64  `json:"vat_value" validate:"omitempty,required"`
	ItemCdnName        *string  `json:"item_cdn_name"`
	OrderDetailID      int      `json:"order_detail_id"`
	// InvoiceDate        string   `json:"invoice_date"`
	// TotalBelumBayar    float64  `json:"total_belum_bayar"`
}

type ShipmentQueryFilter struct {
	StartDate    string `query:"start_date"`
	EndDate      string `query:"end_date"`
	DriverID     int    `query:"driver_id"`
	VehicleID    int    `query:"vehicle_id"`
	OutletName   string `query:"outlet_name"`
	DriverName   string `query:"driver_name"`
	ShipmentNo   string `query:"shipment_no"`
	CustID       string `query:"cust_id"`
	DeliveryDate string `query:"delivery_date"`
	Sort         string `query:"sort"`
}

type ShipmentPickList struct {
	TotalInvoice    float64               `json:"total_invoice"`
	TotalBelumBayar float64               `json:"total_belum_bayar"`
	TotalNetto      float64               `json:"total_netto"`
	SipmentInvoice  []ShipmentInvoiceList `json:"invoices"`
}

type ShipmentInvoiceList struct {
	ID                 int        `json:"id"`
	ShipmentNo         *string    `json:"shipment_no"`
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
	ReasonID           *int       `json:"reason_id"`
	ReasonName         *string    `json:"reason_name"`
	Qty1               int64      `json:"qty1"`
	Qty2               int64      `json:"qty2"`
	Qty3               int64      `json:"qty3"`
	QtyReject1         int64      `json:"qty_reject_1" gorm:"column:qty_reject_1"`
	QtyReject2         int64      `json:"qty_reject_2" gorm:"column:qty_reject_2"`
	QtyReject3         int64      `json:"qty_reject_3" gorm:"column:qty_reject_3"`
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
	TotalNetto         *float64   `json:"total_netto"`
	TotalVolumeDisc    *float64   `json:"total_volume_disc"`
	TotalPromo         *float64   `json:"total_promo"`
	TotalPpn           *float64   `json:"total_ppn"`
	DueDate            *time.Time `json:"due_date"`
	SellPrice1         *float64   `json:"sell_price1"`
	SellPrice2         *float64   `json:"sell_price2"`
	SellPrice3         *float64   `json:"sell_price3"`
	PayTypeName        string     `json:"pay_type_name"`
	InvoiceDate        string     `json:"invoice_date"`
	TotalBelumBayar    *float64   `json:"total_belum_bayar"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at"`
	ItemCdnName        *string    `json:"item_cdn_name"`
	NilaiInvoice       float64    `json:"nilai_invoice"`
	OrderDetailID      int        `json:"order_detail_id"`
}

type SubmitShipmentRequest struct {
	ShipmentNo string   `json:"shipment_no" validate:"required"`
	OutletName []string `json:"outlet_name" validate:"required,notEmptyStringSlice"`
	RouteID    []int    `json:"route_id" validate:"required,notEmptyIntSlice"`
}

type DeleteShipmentRequest struct {
	ShipmentNo []string `json:"shipment_no" validate:"required,notEmptyStringSlice"`
}

type UpdateStatusOrder struct {
	Orders []OrderItem `json:"orders"`
}

type OrderItem struct {
	OrderNo string `json:"ro_no"`
	Status  int    `json:"data_status"`
}

type ReturnItem struct {
	OrderNo string `json:"return_no"`
	Status  int    `json:"status"`
}

type UpdateStatusReturn struct {
	Returns []ReturnItem `json:"returns"`
}
