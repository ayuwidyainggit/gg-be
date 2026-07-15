package entity

type OutletResponse struct {
	OutletID             int    `json:"outlet_id" example:"34"`
	OutletCode           string `json:"outlet_code" example:"18"`
	OutletName           string `json:"outlet_name" example:"Lapan Belas"`
	OutletAddress        string `json:"outlet_address" example:"address 1"`
	OutletStatus         string `json:"outlet_status" example:"On Progress"`
	TotalProduct         int    `json:"total_product" example:"3"`
	TotalProductDelivery int    `json:"total_product_delivery" example:"3"`
	TotalProductPickup   int    `json:"total_product_pickup" example:"3"`
}

type OutletParams struct {
	DriverId   int    `params:"driverId"   validate:"required"`
	OutletId   int    `params:"outletId"   validate:"required"`
	ShipmentNo string `params:"shipmentNo" validate:"required"`
}
