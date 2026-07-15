package model

import "time"

type OutletVisitList struct {
	ID               int64      `gorm:"column:id;primaryKey" json:"id"`
	OutletID         *int64     `gorm:"column:outlet_id" json:"outlet_id"`
	OutletCode       *string    `gorm:"column:outlet_code" json:"outlet_code"`
	OutletName       *string    `gorm:"column:outlet_name" json:"outlet_name"` // Populated from JOIN with mst.m_outlet
	Date             time.Time  `gorm:"column:date" json:"date"`
	ArriveAt         *int64     `gorm:"column:arrive_at" json:"arrive_at"`
	Latitude         *string    `gorm:"column:latitude" json:"latitude"`
	Longitude        *string    `gorm:"column:longitude" json:"longitude"`
	PhotoPath        *string    `gorm:"column:photo_path" json:"photo_path"`
	Folder           *string    `gorm:"column:folder" json:"folder"`
	IsUpdateLocation *bool      `gorm:"column:is_update_location" json:"is_update_location"`
	CreatedAt        time.Time  `gorm:"column:created_at" json:"created_at"`
	UpdatedAt        *time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func (OutletVisitList) TableName() string {
	return "pjp.outlet_visit_list"
}
