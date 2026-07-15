package model

import "time"

type Area struct {
	CustId        string     `db:"cust_id" json:"cust_id"`
	AreaId        int        `db:"area_id" json:"area_id"`
	AreaCode      string     `db:"area_code" json:"area_code"`
	AreaName      string     `db:"area_name" json:"area_name"`
	RegionId      int        `db:"region_id" json:"region_id"`
	OfficialId    int        `db:"official_id" json:"official_id"`
	IsActive      bool       `db:"is_active" json:"is_active"`
	IsDel         bool       `db:"is_del" json:"is_del"`
	CreatedBy     *int64     `db:"created_by,omitempty" json:"created_by"`
	CreatedAt     *time.Time `db:"created_at,omitempty" json:"created_at"`
	UpdatedBy     *int64     `db:"updated_by,omitempty" json:"updated_by"`
	UpdatedAt     *time.Time `db:"updated_at,omitempty" json:"updated_at"`
	UpdatedByName *string    `db:"updated_by_name" json:"updated_by_name"`
	DeletedBy     *int64     `db:"deleted_by,omitempty" json:"deleted_by"`
	DeletedAt     *time.Time `db:"deleted_at,omitempty" json:"deleted_at"`
}

type AreaList struct {
	CustId                string     `db:"cust_id" json:"cust_id"`
	AreaId                int        `db:"area_id" json:"area_id"`
	AreaCode              string     `db:"area_code" json:"area_code"`
	AreaName              string     `db:"area_name" json:"area_name"`
	RegionId              int        `db:"region_id" json:"region_id"`
	RegionCode            string     `db:"region_code" json:"region_code"`
	RegionName            string     `db:"region_name" json:"region_name"`
	OfficialId            *int       `db:"official_id" json:"official_id"`
	OfficialType          *int       `db:"official_type" json:"official_type"`
	OfficialHierarchyCode *string    `db:"official_hierarchy_code" json:"official_hierarchy_code"`
	OfficialEmpName       *string    `db:"official_emp_name" json:"official_emp_name"`
	IsActive              bool       `db:"is_active" json:"is_active"`
	IsDel                 bool       `db:"is_del" json:"is_del"`
	CreatedBy             *int64     `db:"created_by,omitempty" json:"created_by"`
	CreatedAt             *time.Time `db:"created_at,omitempty" json:"created_at"`
	UpdatedBy             *int64     `db:"updated_by,omitempty" json:"updated_by"`
	UpdatedAt             *time.Time `db:"updated_at,omitempty" json:"updated_at"`
	UpdatedByName         *string    `db:"updated_by_name" json:"updated_by_name"`
	DeletedBy             *int64     `db:"deleted_by,omitempty" json:"deleted_by"`
	DeletedAt             *time.Time `db:"deleted_at,omitempty" json:"deleted_at"`
}

type AreaUpdate struct {
	AreaCode   *string    `json:"area_code,omitempty" sql:"area_code"`
	AreaName   *string    `json:"area_name,omitempty" sql:"area_name"`
	RegionId   *int       `json:"region_id,omitempty" sql:"region_id"`
	OfficialId *int       `json:"official_id,omitempty" sql:"official_id"`
	IsActive   *bool      `json:"is_active" sql:"is_active"`
	UpdatedAt  *time.Time `json:"updated_at" sql:"updated_at"`
	UpdatedBy  *int64     `json:"updated_by" sql:"updated_by"`
}
