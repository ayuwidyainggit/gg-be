package model

import "time"

type OutletVisitList struct {
	ID              int       `gorm:"type:int;primary_key" json:"id"`
	Year            int       `gorm:"column:year;type:int;null" json:"year"`
	Week            int       `gorm:"column:week;type:int;null" json:"week"`
	Date            time.Time `gorm:"column:date;null" json:"date"`
	Day             string    `gorm:"column:day;type:varchar(125);null" json:"day"`
	RouteCode       *int      `gorm:"column:route_code;type:int;null" json:"route_code"`
	DestinationID   int       `gorm:"column:outlet_id;null" json:"outlet_id"`
	DestinationCode string    `gorm:"column:outlet_code;null" json:"outlet_code"`
	PjpID           *int      `gorm:"column:pjp_id;type:int;null" json:"pjp_id"`
	PjpCode         *int      `gorm:"column:pjp_code;type:int;null" json:"pjp_code"`
	Start           *int64    `gorm:"column:start;null" json:"start"`
	Finish          *int64    `gorm:"column:finish;null" json:"finish"`
	SkipAt          *int64    `gorm:"column:skip_at;null" json:"skip_at"`
	LeaveAt         *int64    `gorm:"column:leave_at;null" json:"leave_at"`
	ArriveAt        *int64    `gorm:"column:arrive_at;null" json:"arrive_at"`
	OnHold          *int64    `gorm:"column:on_hold;null" json:"on_hold"`
	ResumeAt        *int64    `gorm:"column:resume_at;null" json:"resume_at"`
	SkipReason      *string   `gorm:"column:skip_reason;null" json:"skip_reason"`
	CreatedAt       time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt       time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
	IsPlanned       bool      `gorm:"column:is_planned;true" json:"is_planned"`

	DueDate            string `gorm:"->" json:"due_date"`
	Status             string `gorm:"->" json:"status"`
	DestinationName    string `gorm:"->" json:"outlet_name"`
	DestinationAddress string `gorm:"->" json:"outlet_address"`

	Pjp   *Pjp   `gorm:"foreignKey:pjp_id;references:ID"`
	Route *Route `gorm:"foreignKey:route_code;references:RouteCode"`
}

func (OutletVisitList) TableName() string {
	return "pjp_principles.outlet_visit_list"
}
