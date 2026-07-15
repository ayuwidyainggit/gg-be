package entity

type SummaryResponse struct {
	StartTime  int64 `json:"start_time" example:"1719122155729"`
	EndTime    int64 `json:"end_time" example:"0"`
	Trip       int   `json:"trip" example:"3"`
	Finished   int   `json:"finished" example:"1"`
	InProgress int   `json:"in_progress" example:"0"`
	Shipment   int   `json:"shipments" example:"10"`
}

type SummaryParams struct {
	DriverId int    `params:"driverId" validate:"required"`
	CustId   string `params:"custId"   validate:"required"`
}

type SummaryDailyResponse struct {
	StartTime int64 `json:"start_time" example:"1719122155729"`
	EndTime   int64 `json:"end_time"   example:"1719122155729"`
}

type SummaryDailyParams struct {
	ShipmentNo string `params:"shipmentNo"   validate:"required"`
	CustId     string `params:"custId"   validate:"required"`
}
