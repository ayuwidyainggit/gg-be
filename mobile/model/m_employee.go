package model

import (
	"time"

	"gorm.io/gorm"
)

type MEmployee struct {
	CustId        string     `gorm:"column:cust_id" json:"cust_id"`
	EmpId         int64      `gorm:"column:emp_id" json:"emp_id"`
	EmpCode       string     `gorm:"column:emp_code" json:"emp_code"`
	EmpName       string     `gorm:"column:emp_name" json:"emp_name"`
	Address       *string    `gorm:"column:address" json:"address"`
	City          *string    `gorm:"column:city" json:"city"`
	LastEducation *string    `gorm:"column:last_education" json:"last_education"`
	PhoneNo       *string    `gorm:"column:phone_no" json:"phone_no"`
	WaNo          *string    `gorm:"column:wa_no" json:"wa_no"`
	Email         *string    `gorm:"column:email" json:"email"`
	EmpTypeId     *string    `gorm:"column:emp_type_id" json:"emp_type_id"`
	EmpTypeName   *string    `gorm:"column:emp_type_name" json:"emp_type_name"`
	EmpGrpId      *int64     `gorm:"column:emp_grp_id" json:"emp_grp_id"`
	EmpGrpCode    *string    `gorm:"column:emp_grp_code" json:"emp_grp_code"`
	EmpGrpName    *string    `gorm:"column:emp_grp_name" json:"emp_grp_name"`
	Dob           *time.Time `gorm:"column:dob" json:"dob"`
	WorkDate      *time.Time `gorm:"column:work_date" json:"work_date"`
	DeviceID      *string    `gorm:"column:device_id" json:"device_id"`
	MacAddress    *string    `gorm:"column:mac_address" json:"mac_address"`
	IsActive      bool       `gorm:"column:is_active" json:"is_active"`
	IsDel         bool       `gorm:"column:is_del" json:"is_del"`
	CreatedBy     *int64     `gorm:"column:created_by,omitempty" json:"created_by"`
	CreatedAt     *time.Time `gorm:"column:created_at,omitempty" json:"created_at"`
	UpdatedBy     *int64     `gorm:"column:updated_by,omitempty" json:"updated_by,omitempty"`
	UpdatedByName *string    `json:"updated_by_name" gorm:"column:updated_by_name"`
	UpdatedAt     *time.Time `gorm:"column:updated_at,omitempty" json:"updated_at,omitempty"`
	DeletedBy     *int64     `gorm:"column:deleted_by,omitempty" json:"deleted_by"`
	DeletedAt     *time.Time `gorm:"column:deleted_at,omitempty" json:"deleted_at"`
}

func (MEmployee) TableName() string {
	return "mst.m_employee"
}

func (m *MEmployee) BeforeUpdate(trx *gorm.DB) (err error) {
	timeNow := time.Now()
	m.UpdatedAt = &timeNow
	return nil
}

type DistributorDetail struct {
	DistributorID   int64  `db:"distributor_id" json:"distributor_id"`
	DistributorCode string `db:"distributor_code" json:"distributor_code"`
	DistributorName string `db:"distributor_name" json:"distributor_name"`
	Address         string `db:"address" json:"address"`
	AreaID          int    `db:"area_id" json:"area_id"`
	AreaCode        string `db:"area_code" json:"area_code"`
	AreaName        string `db:"area_name" json:"area_name"`
	RegionID        int    `db:"region_id" json:"region_id"`
	RegionCode      string `db:"region_code" json:"region_code"`
	RegionName      string `db:"region_name" json:"region_name"`
}

func (DistributorDetail) TableName() string {
	return "mst.m_distributor"
}
