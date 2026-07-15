package model

import "time"

type Shipment struct {
	ID           int       `json:"id" gorm:"type:int;primary_key"`
	ShipmentNo   string    `json:"shipment_no" gorm:"type:varchar(125)"`
	DriverID     int       `json:"driver_id"`
	DriverName   string    `json:"driver_name"`
	HelperID     int       `json:"helper_id"`
	HelperName   string    `json:"helper_name"`
	VehicleID    int       `json:"vehicle_id" gorm:"type:int"`
	VehicleNo    string    `json:"vehicle_no"`
	VehicleType  string    `json:"vehicle_type" gorm:"type:varchar(125)"`
	VehicleName  string    `json:"vehicle_name"`
	Length       float64   `json:"length"`
	Width        float64   `json:"width"`
	Height       float64   `json:"height"`
	Volume       float64   `json:"volume"`
	DeliveryDate time.Time `json:"delivery_date"`
	Start        *int64    `json:"start"`
	Finish       *int64    `json:"finish"`
	Status       string    `json:"status" gorm:"type:varchar(125);default:Planned"`
	CustID       string    `json:"cust_id"`
	Weight       float64   `json:"weight"`
	ShipmentType string    `json:"shipment_type" gorm:"type:varchar(125)"`
	CreatedAt    time.Time `json:"created_at" gorm:"column:created_at;autoCreateTime"`
	UpdatedAt    time.Time `json:"updated_at" gorm:"column:updated_at;autoUpdateTime"`

	ShipmentInvoices []ShipmentInvoices `gorm:"foreignKey:shipment_no;references:ShipmentNo"`
}

type Tabler interface {
	TableName() string
}

func (Shipment) TableName() string {
	return "tms.shipments"
}

//func (m Shipment) ConvertStart() string {
//	t := time.Unix(0, *m.Start*int64(time.Millisecond))
//	return t.Format(time.RFC3339)
//}
//func (m Shipment) ConvertFinish() string {
//	t := time.Unix(0, *m.Finish*int64(time.Millisecond))
//	return t.Format(time.RFC1123)
//}
