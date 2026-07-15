package model

import (
	"time"
)

type OutletVisitList struct {
	ID         int       `gorm:"type:int;primary_key" json:"id"`
	Year       int       `gorm:"column:year;type:int;null" json:"year"`
	Week       int       `gorm:"column:week;type:int;null" json:"week"`
	Date       time.Time `gorm:"column:date;null" json:"date"`
	Day        string    `gorm:"column:day;type:varchar(125);null" json:"day"`
	RouteCode  *int      `gorm:"column:route_code;type:int;null" json:"route_code"`
	OutletID   int       `gorm:"column:outlet_id;null" json:"outlet_id"`
	OutletCode string    `gorm:"column:outlet_code;null" json:"outlet_code"`
	PjpID      *int      `gorm:"column:pjp_id;type:int;null" json:"pjp_id"`
	PjpCode    *int      `gorm:"column:pjp_code;type:int;null" json:"pjp_code"`
	Start      *int64    `gorm:"column:start;null" json:"start"`
	Finish     *int64    `gorm:"column:finish;null" json:"finish"`
	SkipAt     *int64    `gorm:"column:skip_at;null" json:"skip_at"`
	LeaveAt    *int64    `gorm:"column:leave_at;null" json:"leave_at"`
	ArriveAt   *int64    `gorm:"column:arrive_at;null" json:"arrive_at"`
	OnHold     *int64    `gorm:"column:on_hold;null" json:"on_hold"`
	ResumeAt   *int64    `gorm:"column:resume_at;null" json:"resume_at"`
	SkipReason *string   `gorm:"column:skip_reason;null" json:"skip_reason"`
	CreatedAt  time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt  time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
	IsPlanned  bool      `gorm:"column:is_planned;true" json:"is_planned"`

	// File upload fields
	PhotoPath        string `gorm:"column:photo_path;type:varchar(500);null" json:"photo_path"`
	Folder           string `gorm:"column:folder;type:varchar(255);null" json:"folder"`
	FileName         string `gorm:"column:file_name;type:varchar(255);null" json:"file_name"`
	FileType         string `gorm:"column:file_type;type:varchar(50);null" json:"file_type"`
	MediaCategory    string `gorm:"column:media_category;type:varchar(20);null" json:"media_category"`
	FileUrl          string `gorm:"column:file_url;type:text;null" json:"file_url"`
	FileSize         *int64 `gorm:"column:file_size;type:bigint;null" json:"file_size"`
	FileBase64       string `gorm:"column:file_base64;type:text;null" json:"file_base64"`
	IsUpdateLocation bool   `gorm:"column:is_update_location;type:boolean;not null;default:false" json:"is_update_location"`
	Latitude         string `gorm:"column:latitude;type:varchar(50);null" json:"latitude"`
	Longitude        string `gorm:"column:longitude;type:varchar(50);null" json:"longitude"`
	DistanceMeter    *int   `gorm:"column:distance_meter;type:integer;comment:jarak antara titik lokasi dengan titik lokasi master (meter)"`
	AllowedRadius    *int   `gorm:"column:allowed_radius;type:integer;default:100;comment:jarak yang diperbolehkan (meter)"`
	LocationStatus   *int   `gorm:"column:location_status;type:smallint;comment:0=OUT_OF_RADIUS,1=IN_RADIUS"`
	SkipInOutlet     bool   `gorm:"column:skip_in_outlet" json:"skip_in_outlet"`
	LeaveLongitude   string `gorm:"column:leave_longitude" json:"leave_longitude"`
	LeaveLatitude    string `gorm:"column:leave_latitude" json:"leave_latitude"`
	DestinationType  string `gorm:"column:destination_type" json:"destination_type"`

	DueDate         string `gorm:"->" json:"due_date"`
	Status          string `gorm:"->" json:"status"`
	OutletName      string `gorm:"->" json:"outlet_name"`
	OutletAddress   string `gorm:"->" json:"outlet_address"`
	OutletLongitude string `gorm:"->" json:"outlet_longitude"`
	OutletLatitude  string `gorm:"->" json:"outlet_latitude"`
	Top             int    `gorm:"->" json:"top"`

	Pjp   *Pjp   `gorm:"foreignKey:pjp_id;references:ID"`
	Route *Route `gorm:"foreignKey:route_code;references:RouteCode"`
}

func (OutletVisitList) TableName() string {
	return "pjp.outlet_visit_list"
}
