package entity

type TravelListResponse struct {
	ArriveAt   int64  `json:"arrive_at"   example:"1719122155729"`
	UnloadAt   int64  `json:"unload_at"   example:"1719122155729"`
	LeaveAt    int64  `json:"leave_at"    example:"1719122155729"`
	UnloadDesc string `json:"unload_desc" example:""`
	PickupAt   int64  `json:"pickup_at"  example:"1719122155729"`
	SkipAt     int64  `json:"skip_at"     example:"1719122155729"`
	OnHold     int64  `json:"on_hold"     example:"1719122155729"`
	ResumeAt   int64  `json:"resume_at"     example:"1719122155729"`
	SkipReason string `json:"skip_reason" example:"outlet closed"`
	InOutlet   bool   `json:"in_outlet"   example:"true"`
}

type TravelListParams struct {
	OutletId   int    `params:"outletId"   validate:"required"`
	ShipmentNo string `params:"shipmentNo" validate:"required"`
}
