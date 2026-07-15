package entity

type LoginSendPick struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SendPickResponse struct {
	ShipmentNo string   `json:"shipment_no" example:"fleet_10_1"`
	Orders     []string `json:"orders" example:"INV1243123"`
}

type MapperSendPick struct {
	ShipmentNumbers       []string
	VehicleIDs            []int
	ShipmentNumbersMapper []string
}

type GenerateResponse struct {
	Message string `json:"message"`
	Data    struct {
		Result []struct {
			Vehicle struct {
				ShipmentNo string `json:"shipment_no"`
			} `json:"vehicle"`
			ItemDetails []struct {
				OrderID string `json:"order_id"`
			} `json:"item_details"`
		} `json:"result"`
		Report struct {
			TotalVehicleUsed    int `json:"total_vehicle_used"`
			TotalItemAssigned   int `json:"total_item_assigned"`
			TotalItemUnassigned int `json:"total_item_unassigned"`
		} `json:"report"`
		Unassigned []struct {
			Reason []struct {
				Description string `json:"description"`
			} `json:"reason"`
		} `json:"unassigned"`
	} `json:"data"`
}
