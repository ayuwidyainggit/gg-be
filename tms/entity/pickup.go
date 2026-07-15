package entity

type PickUpRequest struct {
	ID          []int `validate:"required,notEmptyIntSlice" json:"id" example:"1218"`
	CurrentTime int64 `validate:"required"  json:"current_time" example:"1727233938570"`
}

type SkipPickUpRequest struct {
	Data []SkipPickUpBody `validate:"required" json:"data"`
}

type SkipPickUpBody struct {
	ID         int    `validate:"required" json:"id" example:"1218"`
	ReasonID   int    `validate:"required"  json:"reason_id" example:"1"`
	ReasonName string `validate:"required"  json:"reason_name" example:"Reject"`
}

type PickupPartialRequest struct {
	Data PickupPartialBody `validate:"required" json:"data"`
}

type PickupPartialBody struct {
	OutletID    int       `validate:"required" json:"outlet_id" example:"88"`
	ShipmentNo  string    `validate:"required" json:"shipment_no" example:"DO202409250178"`
	CurrentTime int64     `validate:"required"  json:"current_time" example:"1727233938570"`
	Products    []Product `validate:"required,dive" json:"products"`
}
