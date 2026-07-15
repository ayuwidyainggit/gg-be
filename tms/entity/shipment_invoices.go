package entity

type ShipmentInvoicesResponse struct {
	ShipmentID    int    `json:"id" example:"1"`
	ShipmentNo    string `json:"shipment_no" example:"11024066062"`
	OrderNo       string `json:"order_no" example:"OD001"`
	DriverID      int    `json:"driver_id" example:"1"`
	InvoiceNo     string `json:"invoice_no" example:"INV1243123"`
	OutletID      int    `json:"outlet_id" example:"34"`
	OutletCode    string `json:"outlet_code" example:"18"`
	OutletAddress string `json:"outlet_address" example:"address 1"`
	OutletStatus  string `json:"outlet_status" example:"On Progress"`
	OutletName    string `json:"outlet_name" example:"Lapan belas"`
	SalesmanID    int    `json:"salesman_id" example:"1"`
	SalesmanName  string `json:"salesman_name" example:"John Sales"`
	ProductID     int    `json:"product_id" example:"1"`
	ProductName   string `json:"product_name" example:"ASSORTED BISCUIT OPP RED 20 275G"`
	ProductStatus string `json:"product_status" example:"Reject"`
	Status        string `json:"status" example:"Delivery"`
	Qty1          int64  `json:"qty1" example:"10"`
	Qty2          int64  `json:"qty2" example:"10"`
	Qty3          int64  `json:"qty3" example:"10"`
	ConvUnit1     int64  `json:"conv_unit1" example:"0"`
	ConvUnit2     int64  `json:"conv_unit2" example:"6"`
	ConvUnit3     int64  `json:"conv_unit3" example:"6"`
	UnitId1       string `json:"unit_id1" example:"PCS"`
	UnitId2       string `json:"unit_id2" example:"CTN"`
	UnitId3       string `json:"unit_id3" example:"CTN"`
}

type ShipmentInvoicesQueryFilter struct {
	StartDate    string `query:"start_date"`
	EndDate      string `query:"end_date"`
	OutletName   string `query:"outlet_name"`
	OutletID     int    `query:"outlet_id"`
	OutletCode   string `query:"outlet_code"`
	ShipmentNo   string `query:"shipment_no"`
	CustID       string `query:"cust_id"`
	Status       string `query:"status"`
	DriverID     int    `query:"driver_id"`
	DriverName   int    `query:"driver_name"`
	ProductName  string `query:"product_name"`
	ProductID    int    `query:"product_id"`
	SalesmanName string `query:"salesman_name"`
	Sort         string `query:"sort"`
}
