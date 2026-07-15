package model

import "time"

// ArrivalReport represents the arrival report data for salesman visits
// Table: pjp.arrival_report
type ArrivalReport struct {
	ArrivalId        int64      `gorm:"column:arrival_id;primaryKey;autoIncrement" json:"arrival_id"`
	CustId           string     `gorm:"column:cust_id;type:varchar(50);not null" json:"cust_id"`
	OutletId         int64      `gorm:"column:outlet_id;not null" json:"outlet_id"`
	UserId           int64      `gorm:"column:user_id;not null" json:"user_id"`
	Activity         string     `gorm:"column:activity;type:varchar(50);not null" json:"activity"`
	ArrivalLongitude *string    `gorm:"column:arrival_longitude;type:varchar(50)" json:"arrival_longitude"`
	ArrivalLatitude  *string    `gorm:"column:arrival_latitude;type:varchar(50)" json:"arrival_latitude"`
	OutletLongitude  *string    `gorm:"column:outlet_longitude;type:varchar(50)" json:"outlet_longitude"`
	OutletLatitude   *string    `gorm:"column:outlet_latitude;type:varchar(50)" json:"outlet_latitude"`
	DistanceMeter    *int       `gorm:"column:distance_meter" json:"distance_meter"`
	AllowedRadius    *int       `gorm:"column:allowed_radius" json:"allowed_radius"`
	LocationStatus   *string    `gorm:"column:location_status;type:varchar(50)" json:"location_status"`
	CreatedBy        *string    `gorm:"column:created_by;type:varchar(100)" json:"created_by"`
	CreatedAt        *time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedBy        *string    `gorm:"column:updated_by;type:varchar(100)" json:"updated_by"`
	UpdatedAt        *time.Time `gorm:"column:updated_at" json:"updated_at"`
	DeletedBy        *string    `gorm:"column:deleted_by;type:varchar(100)" json:"deleted_by"`
	DeletedAt        *time.Time `gorm:"column:deleted_at" json:"deleted_at"`
	IsDel            *bool      `gorm:"column:is_del;default:false" json:"is_del"`
}

// TableName returns the table name for ArrivalReport
func (ArrivalReport) TableName() string {
	return "pjp.arrival_report"
}
