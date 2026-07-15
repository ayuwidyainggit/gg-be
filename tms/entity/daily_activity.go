package entity

type DailyActivityResponse struct {
	ShipmentID int    `json:"id" example:"1"`
	ShipmentNo string `json:"shipment_no" example:"11024051411"`
	Status     string `json:"status" example:"pending"`
	IsActive   bool   `json:"is_active" example:"true"`
}

type DailyActivityParams struct {
	DriverId int `params:"driverId" validate:"required"`
}
