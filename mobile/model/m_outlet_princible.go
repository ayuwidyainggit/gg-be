package model

import "time"

type PJPPrincipleDestinationList struct {
	RouteCode          int       `gorm:"column:route_code" json:"route_code"`
	RouteName          string    `gorm:"column:route_name" json:"route_name"`
	Year               int       `gorm:"column:year" json:"year"`
	Week               int       `gorm:"column:week" json:"week"`
	Date               time.Time `gorm:"column:date" json:"date"`
	DestinationID      int       `gorm:"column:destination_id" json:"outlet_id"`
	DestinationCode    string    `gorm:"column:destination_code" json:"outlet_code"`
	DestinationName    string    `gorm:"column:destination_name" json:"outlet_name"`
	DestinationStatus  string    `gorm:"column:destination_status" json:"outlet_status"`
	DestinationAddress string    `gorm:"column:destination_address" json:"outlet_address"`
	Longitude          string    `gorm:"column:longitude" json:"longitude"`
	Latitude           string    `gorm:"column:latitude" json:"latitude"`
}

// TableName overrides the table name for GORM
func (PJPPrincipleDestinationList) TableName() string {
	return "pjp_principles.destinations_history"
}
