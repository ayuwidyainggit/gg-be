package model

import "time"

type RoutePopDaily struct {
	ID          int       `gorm:"type:int;primary_key" json:"id"`
	Year        int       `gorm:"column:year;null" json:"year"`
	Week        int       `gorm:"column:week;type:int;null" json:"week"`
	Date        time.Time `gorm:"column:date;null" json:"date"`
	Day         string    `gorm:"column:day;type:varchar(125);null" json:"day"`
	RouteCode   *int      `gorm:"column:route_code;type:int;null" json:"route_code"`
	PjpID       *int      `gorm:"column:pjp_id;type:int;null" json:"pjp_id"`
	PjpCode     *int      `gorm:"column:pjp_code;type:int;null" json:"pjp_code"`
	Status      string    `gorm:"column:status;type:varchar(125);default:active" json:"status"`
	ParentRoute *int      `gorm:"column:parent_route;type:int;null" json:"parent_route"`
	CustID      string    `gorm:"column:cust_id;type:varchar(125);null" json:"cust_id"`
	CreatedAt   time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`

	Pjp   *Pjp   `gorm:"foreignKey:pjp_id;references:ID"`
	Route *Route `gorm:"foreignKey:route_code;references:RouteCode"`
	// RouteOutletAdditional *RouteOutletAdditional `gorm:"foreignKey:parent_route;references:ParentRoute"`
}

func (RoutePopDaily) TableName() string {
	return "pjp_principles.route_pop_dailies"
}

type RoutePopDailyWithOutlet struct {
	RoutePopDaily
	RouteCode          int
	DestinationCode    string
	DestinationName    string
	Longitude          string
	Latitude           string
	DestinationAddress string
	DestinationStatus  string
	SalesmanName       string
}
