package entity

type VisitRequest struct {
	CurrentTime int64  `json:"current_time" validate:"required"`
	DriverID    int    `json:"driver_id" validate:"required"`
	ShipmentNo  string `json:"shipment_no" validate:"required"`
}

type LeaveRequest struct {
	CurrentTime int64  `json:"current_time" validate:"required"`
	OutletID    int    `json:"outlet_id" validate:"required"`
	ShipmentNo  string `json:"shipment_no" validate:"required"`
}

type ArriveRequest struct {
	CurrentTime int64  `json:"current_time" validate:"required"`
	OutletID    int    `json:"outlet_id" validate:"required"`
	ShipmentNo  string `json:"shipment_no" validate:"required"`
}

type SkipRequest struct {
	CurrentTime int64  `json:"current_time" validate:"required"`
	OutletID    int    `json:"outlet_id" validate:"required"`
	ShipmentNo  string `json:"shipment_no" validate:"required"`
	SkipReason  string `json:"skip_reason" validate:"required"`
	InOutlet    *bool  `json:"in_outlet" validate:"required"`
}
