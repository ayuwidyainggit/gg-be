package model

import "time"

type PJPDistributorOutletList struct {
	RouteCode     int       `gorm:"column:route_code" json:"route_code"`
	RouteName     string    `gorm:"column:route_name" json:"route_name"`
	Year          int       `gorm:"column:year" json:"year"`
	Week          int       `gorm:"column:week" json:"week"`
	Date          time.Time `gorm:"column:date" json:"date"`
	OutletID      int       `gorm:"column:outlet_id" json:"outlet_id"`
	OutletCode    string    `gorm:"column:outlet_code" json:"outlet_code"`
	OutletName    string    `gorm:"column:outlet_name" json:"outlet_name"`
	OutletStatus  string    `gorm:"column:outlet_status" json:"outlet_status"`
	OutletAddress string    `gorm:"column:outlet_address" json:"outlet_address"`
	Longitude     string    `gorm:"column:longitude" json:"longitude"`
	Latitude      string    `gorm:"column:latitude" json:"latitude"`
}

// TableName overrides the table name for GORM
func (PJPDistributorOutletList) TableName() string {
	return "pjp.route_outlet_history"
}
