package model

import "time"

type RoutePopPermanent struct {
	ID        int       `gorm:"type:int;primary_key" json:"id"`
	Year      int       `gorm:"column:year;null" json:"year"`
	Week      int       `gorm:"column:week;null" json:"week"`
	Date      time.Time `gorm:"column:date;null" json:"date"`
	Day       string    `gorm:"column:day;type:varchar(125);null" json:"day"`
	RouteCode *int      `gorm:"column:route_code;type:int;null" json:"route_code"`
	PjpID     *int      `gorm:"column:pjp_id;type:int;null" json:"pjp_id"`
	PjpCode   *int      `gorm:"column:pjp_code;type:int;null" json:"pjp_code"`
	CustID    string    `gorm:"column:cust_id;type:varchar(125);null" json:"cust_id"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`

	Pjp   *Pjp   `gorm:"foreignKey:pjp_id;references:ID"`
	Route *Route `gorm:"foreignKey:route_code;"`
}

func (RoutePopPermanent) TableName() string {
	return "pjp_principles.route_pop_permanent"
}
