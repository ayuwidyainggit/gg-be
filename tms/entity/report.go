package entity

type DriverReportResponse struct {
	TotalShipment  int              `json:"total_shipment"`
	TotalTrip      int              `json:"total_trip"`
	TotalFinished  int              `json:"total_finished"`
	TotalSkipped   int              `json:"total_skipped"`
	Progress       float64          `json:"progress"`
	SkippedReasons *[]SkippedReason `json:"skipped_reasons"`
}

type SkippedReason struct {
	SkipReason   string  `json:"reason_name"`
	DeliveryDate string  `json:"-"`
	Count        int     `json:"count"`
	Percentage   float64 `json:"percentage"`
}

type Overview struct {
	Finished int `json:"finished"`
	Skipped  int `json:"skipped"`
}

type DriverReportQueryFilter struct {
	Period   string `query:"period" validate:"required"`
	DriverID int    `query:"driver_id" validate:"required"`
}
