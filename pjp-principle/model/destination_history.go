package model

import "time"

type DestinationHistory struct {
	ID                 int     `gorm:"type:int;primary_key" json:"id"`
	RouteCode          int     `gorm:"column:route_code;type:int;not null" json:"route_code"`
	RouteName          string  `gorm:"column:route_name;type:varchar(125);not null" json:"route_name"`
	DestinationID      int     `gorm:"column:destination_id;type:int" json:"destination_id"`
	DestinationCode    string  `gorm:"column:destination_code;type:varchar(125)" json:"destination_code"`
	DestinationName    string  `gorm:"column:destination_name;type:varchar(125)" json:"destination_name"`
	Longitude          string  `gorm:"column:longitude;type:varchar(125)" json:"longitude"`
	Latitude           string  `gorm:"column:latitude;type:varchar(125)" json:"latitude"`
	DestinationStatus  string  `gorm:"column:destination_status;type:varchar(125)" json:"destination_status"`
	DestinationAddress string  `gorm:"column:destination_address;type:varchar(125);null" json:"destination_address"`
	DestinationType    string  `gorm:"column:destination_type;type:varchar(125);null" json:"destination_type"`
	AvgSalesWeek       float64 `gorm:"column:avg_sales_week;type:numeric(10,2)" json:"avg_sales_week"`
	PjpID              *int    `gorm:"column:pjp_id;type:int;null" json:"pjp_id"`
	PjpCode            *int    `gorm:"column:pjp_code;type:int;null" json:"pjp_code"`
	// Status             string     `gorm:"column:status;type:varchar(125);default:pending" json:"status"`
	VerifiedDate    *time.Time `gorm:"column:verified_date;null" json:"verified_date"`
	OldPjpID        *int       `gorm:"column:old_pjp_id;type:int;null" json:"old_pjp_id"`
	OldPjpCode      *int       `gorm:"column:old_pjp_code;type:int;null" json:"old_pjp_code"`
	OldRouteCode    int        `gorm:"column:old_route_code;type:int;null" json:"old_route_code"`
	OldRouteName    string     `gorm:"column:old_route_name;type:varchar(125);null" json:"old_route_name"`
	CustID          string     `gorm:"column:cust_id;type:varchar(125);null" json:"cust_id"`
	CreatedAt       time.Time  `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt       time.Time  `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
	Year            int        `gorm:"column:year;null" json:"year"`
	Week            int        `gorm:"column:week;null" json:"week"`
	IndexDay        int        `gorm:"column:index_day;null" json:"indexDay"`
	StartWeek       *time.Time `gorm:"column:start_week;null" json:"startWeek"`
	IsInCurrentYear bool       `gorm:"column:is_in_current_year;null" json:"isInCurrentYear"`
	Date            time.Time  `gorm:"column:date;null" json:"date"`
	IsAdditional    bool       `gorm:"column:is_additional" json:"is_additional"`

	//alias column
	OutletVisitListID int `gorm:"->" json:"outlet_visit_list_id"`

	Pjp    *Pjp   `gorm:"foreignKey:pjp_id;references:ID"`
	PjpOld *Pjp   `gorm:"foreignKey:old_pjp_id;references:ID"`
	Route  *Route `gorm:"foreignKey:route_code;references:RouteCode"`
}

func (DestinationHistory) TableName() string {
	return "pjp_principles.destinations_history"
}
