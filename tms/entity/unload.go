package entity

type UnloadRequest struct {
	CurrentTime int64  `json:"current_time" validate:"required"`
	OutletID    int    `json:"outlet_id" validate:"required"`
	ShipmentNo  string `json:"shipment_no" validate:"required"`
}
