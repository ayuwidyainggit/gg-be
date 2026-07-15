package model

import "time"

type RoutePrinciple struct {
	ID             int       `gorm:"type:int;primary_key" json:"id"`
	RouteCode      int       `gorm:"column:route_code;type:int;uniqueIndex;not null" json:"route_code"`
	RouteName      string    `gorm:"column:route_name;type:varchar(125);unique;not null" json:"route_name"`
	IsAssign       bool      `gorm:"column:is_assign;type:bool;default:false" json:"is_assign"`
	IsAssignPjp    bool      `gorm:"->" json:"is_assign_pjp"`
	OutletID       int       `gorm:"->" json:"outlet_id"`
	OutletCode     string    `gorm:"->" json:"outlet_code"`
	OutletName     string    `gorm:"->" json:"outlet_name"`
	Longitude      string    `gorm:"->" json:"longitude"`
	Latitude       string    `gorm:"->" json:"latitude"`
	OutletStatus   string    `gorm:"->" json:"outlet_status"`
	OutletAddress  string    `gorm:"->" json:"outlet_address"`
	AvgSalesWeek   float64   `gorm:"->" json:"avg_sales_week"`
	TotalOutlet    int       `gorm:"->" json:"total_outlet"`
	Status         string    `gorm:"->" json:"status"`
	RoutePopStatus string    `gorm:"->" json:"route_pop_status"`
	PjpID          int       `gorm:"column:pjp_id;type:int;null" json:"pjp_id"`
	CustID         string    `gorm:"column:cust_id;type:varchar(125);null" json:"cust_id"`
	CreatedAt      time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt      time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
	IsPjpOld       bool      `gorm:"column:is_pjp_old" json:"is_pjp_old"`
	Sequence       int       `gorm:"column:sequence;type:int;null" json:"sequence"`
}

type DatasetRoutePrinciple struct {
	ID          int       `json:"id"`
	RouteCode   int       `json:"route_code"`
	RouteName   string    `json:"route_name"`
	IsAssign    bool      `json:"is_assign"`
	IsAssignPjp bool      `json:"is_assign_pjp"`
	CustID      string    `json:"cust_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (RoutePrinciple) TableName() string {
	return "pjp_principles.routes"
}
