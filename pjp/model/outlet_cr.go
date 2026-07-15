package model

import "time"

// OutletCr represents the outlet change request master table
type OutletCr struct {
	CustId     string     `gorm:"column:cust_id;type:varchar(10);not null" json:"cust_id"`
	OutletCrId int64      `gorm:"column:outlet_cr_id;type:bigserial;primaryKey" json:"outlet_cr_id"`
	OutletId   int64      `gorm:"column:outlet_id;type:int8;not null" json:"outlet_id"`
	Source     int        `gorm:"column:source;type:int4;not null;default:1" json:"source"`
	Status     int        `gorm:"column:status;type:int4;not null;default:1" json:"status"`
	CreatedBy  *int64     `gorm:"column:created_by;type:int8" json:"created_by"`
	CreatedAt  time.Time  `gorm:"column:created_at;type:timestamptz(6);not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedBy  *int64     `gorm:"column:updated_by;type:int8" json:"updated_by"`
	UpdatedAt  *time.Time `gorm:"column:updated_at;type:timestamptz(6)" json:"updated_at"`
	ApprovalBy *int64     `gorm:"column:approval_by;type:int8" json:"approval_by"`
	ApprovalAt *time.Time `gorm:"column:approval_at;type:timestamptz(6)" json:"approval_at"`
}

func (OutletCr) TableName() string {
	return "mst.outlet_cr"
}

// OutletCrDet represents the outlet change request detail table
type OutletCrDet struct {
	OutletCrDetId int64   `gorm:"column:outlet_cr_det_id;type:bigserial;primaryKey" json:"outlet_cr_det_id"`
	OutletCrId    int64   `gorm:"column:outlet_cr_id;type:int8;not null" json:"outlet_cr_id"`
	FieldName     string  `gorm:"column:field_name;type:varchar(30);not null" json:"field_name"`
	OldValue      *string `gorm:"column:old_value;type:varchar(225)" json:"old_value"`
	NewValue      *string `gorm:"column:new_value;type:varchar(225)" json:"new_value"`
}

func (OutletCrDet) TableName() string {
	return "mst.outlet_cr_det"
}
