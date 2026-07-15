package entity

type VehicleResponse struct {
	VehicleID   int64   `json:"vehicle_id" example:"1"`
	VehicleNo   string  `json:"vehicle_no" example:"G001"`
	VehicleDesc string  `json:"vehicle_desc" example:"Box Double"`
	VehicleType string  `json:"vehicle_type_name" example:"Garuda"`
	Length      float64 `json:"length" example:"0"`
	Width       float64 `json:"width" example:"12"`
	Weight      float64 `json:"weight" example:"11"`
	Height      float64 `json:"height" example:"12"`
	Volume      float64 `json:"volume" example:"1728"`
	DriverID    int64   `json:"driver_id" example:"1"`
	HelperID    int64   `json:"helper_id" example:"2"`
	DriverName  string  `json:"driver_name" example:"John Driver"`
	HelperName  string  `json:"helper_name" example:"George Helper"`
}

type VehicleQueryFilter struct {
	Page         int    `query:"page"`
	Limit        int    `query:"limit" validate:"required"`
	DeliveryDate string `query:"delivery_date"`
	Sort         string `query:"sort"`
}
