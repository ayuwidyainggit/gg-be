package entity

import "time"

type ShipmentNoDropdown struct {
	ShipmentNo string `json:"shipment_no" example:"11024066062"`
}

type ProductCodeDropdown struct {
	ProductCode string `json:"product_code" example:"00123456"`
}

type DriverNameDropdown struct {
	DriverName string `json:"driver_name" example:"John Driver"`
}

type ReasonDropdown struct {
	ReasonName string `json:"reason_name" example:"reject"`
}

type OutletDropdown struct {
	OutletCode string `json:"outlet_code" example:"18"`
	OutletName string `json:"outlet_name" example:"Lapan belas"`
}

type OutletStatus struct {
	Planned *int64 `json:"planned"`
	Visited *int64 `json:"visited"`
	Skipped *int64 `json:"skipped"`
}

// type ShipmentReportSummary struct {
// 	DeliveryDate string `json:"delivery_date" example:"2024-06-11"`
// 	ShipmentNo   string `json:"shipment_no" example:"11024066062"`
// 	DriverName   string `json:"driver_name" example:"John Driver"`
// 	Planned      *int64 `json:"planned" example:"10"`
// 	Visited      *int64 `json:"visited" example:"10"`
// 	Skipped      *int64 `json:"skipped" example:"10"`
// 	StartDate    string `json:"start_date" example:"2024-06-11"`
// 	// ArriveAt      *int64 `json:"arrive_at" example:"2024-06-11"`
// 	LeaveAt       string `json:"leave_at" example:"2024-06-11"`
// 	Spent         *int64 `json:"spent" example:"120"`
// 	ETA           *int64 `json:"eta" example:"400"`
// 	Received      *int64 `json:"received" example:"100"`
// 	RejectPartial *int64 `json:"reject_partial" example:"10"`
// 	RejectAll     *int64 `json:"reject_all" example:"0"`
// 	Photo         string `json:"photo"`
// 	Signature     string `json:"signature"`
// }

type Shipment struct {
	ID               uint              `json:"id"`
	DeliveryDate     string            `json:"delivery_date"`
	ShipmentNo       string            `json:"shipment_no"`
	DriverName       string            `json:"driver_name"`
	Start            *int64            `json:"start"`
	Finish           *int64            `json:"finish"`
	ShipmentInvoices []ShipmentInvoice `json:"shipment_invoices"`
}

type ShipmentInvoice struct {
	ID         uint   `json:"id"`
	Status     string `json:"status"`
	ArriveAt   *int64 `json:"arrive_at"`
	LeaveAt    *int64 `json:"leave_at"`
	Qty1       int64  `json:"qty1"`
	QtyReject1 int64  `json:"qty_reject1"`
	QtyReject2 int64  `json:"qty_reject2"`
	Photo      string `json:"photo"`
	Signature  string `json:"signature"`
}

type ShipmentReportSummary struct {
	DeliveryDate  time.Time `json:"delivery_date"`
	ShipmentNo    string    `json:"shipment_no"`
	DriverName    string    `json:"driver_name"`
	StartTime     *int64    `json:"start_time"`
	EndTime       *int64    `json:"end_time"`
	LeaveAt       int64     `json:"leave_at"`
	Planned       int64     `json:"planned"`
	Visited       int64     `json:"visited"`
	Skipped       int64     `json:"skipped"`
	Spent         *int64    `json:"spent"`
	ETA           *int64    `json:"eta"`
	Received      int64     `json:"received"`
	RejectPartial int64     `json:"reject_partial"`
	RejectAll     int64     `json:"reject_all"`
	Photo         string    `json:"photo"`
	Signature     string    `json:"signature"`
}

// TODO Fix Response
type ShipmentReportDetail struct {
	DeliveryDate          string                  `json:"delivery_date" example:"2024-06-11"`
	ShipmentNo            string                  `json:"shipment_no" example:"11024066062"`
	DriverName            string                  `json:"driver_name" example:"John Driver"`
	ShipmentReportDetails []ShipmentReportDetails `json:"details"`
}

type ShipmentReportDetails struct {
	// InvoiceNo      string `json:"invoice_no" example:"INV1243123"`
	// OrderNo        string `json:"order_no" example:"SO20241025009"`
	DocumentNo     string `json:"document_no" example:"SO20241025009"`
	OutletCode     string `json:"outlet_code" example:"18"`
	OutletName     string `json:"outlet_name" example:"Lapan belas"`
	VisitedStatus  string `json:"visited_status" example:"Delivery"`
	StartTime      int    `json:"start_time"`
	ArriveAt       *int64 `json:"arrive_at"`
	LeaveAt        *int64 `json:"leave_at"`
	EndTime        int    `json:"end_time"`
	DriveTime      int    `json:"drive_time"`
	UnloadAt       *int64 `json:"unload_at"`
	Spent          int    `json:"spent"`
	ETA            int    `json:"eta"`
	ReceivedStatus string `json:"received_status"`
	Photo          string `json:"photo"`
	Signature      string `json:"signature"`
}

type ShipmentReportReject struct {
	DeliveryDate                string                        `json:"delivery_date" example:"2024-06-11"`
	ShipmentNo                  string                        `json:"shipment_no" example:"11024066062"`
	DriverName                  string                        `json:"driver_name" example:"John Driver"`
	ShipmentReportRejectDetails []ShipmentReportRejectDetails `json:"details"`
}

type ShipmentReportRejectDetails struct {
	InvoiceNo   string `json:"invoice_no" example:"INV1243123"`
	OutletCode  string `json:"outlet_code" example:"18"`
	OutletName  string `json:"outlet_name" example:"Lapan belas"`
	ProductName string `json:"product_name" example:"ASSORTED BISCUIT OPP RED 20 275G"`
	ProductCode string `json:"product_code" example:"00123456"`
	QtyReject1  int64  `json:"qty_reject_1" example:"10"`
	QtyReject2  int64  `json:"qty_reject_2" example:"10"`
	QtyReject3  int64  `json:"qty_reject_3" example:"10"`
	ReasonName  string `json:"reason_name" example:"reject"`
}

type ShipmentReportQueryFilter struct {
	StartDate  string `json:"start_date"`
	EndDate    string `json:"end_date"`
	ShipmentNo string `json:"shipment_no"`
	DriverName string `json:"driver_name"`
}

type ShipmentReportDetailQueryFilter struct {
	ShipmentReportQueryFilter
	OutletName     string `json:"outlet_name"`
	VisitedStatus  string `json:"visited_status"`
	ReceivedStatus string `json:"received_status"`
}

type ShipmentReportRejectlQueryFilter struct {
	ShipmentReportQueryFilter
	OutletName  string `json:"outlet_name"`
	ProductCode string `json:"product_code"`
	Reason      string `json:"reason"`
}
