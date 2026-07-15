package model

import (
	"time"
)

type OfficialHierarchy struct {
	CustId        string     `db:"cust_id" json:"cust_id"`
	OfficialType  int        `db:"official_type" json:"official_type"`
	HierarchyCode string     `db:"hierarchy_code" json:"hierarchy_code"`
	IsActive      *bool      `db:"is_active" json:"is_active"`
	CreatedBy     *int64     `db:"created_by,omitempty" json:"created_by"`
	CreatedAt     *time.Time `db:"created_at,omitempty" json:"created_at"`
	UpdatedBy     *int64     `db:"updated_by,omitempty" json:"updated_by"`
	UpdatedAt     *time.Time `db:"updated_at,omitempty" json:"updated_at"`
}

type OfficialHierarchyUpsert struct {
	CustId        *string `sql:"cust_id" json:"cust_id"`
	OfficialType  *int    `sql:"official_type" json:"official_type"`
	HierarchyCode *string `sql:"hierarchy_code" json:"hierarchy_code"`
	IsActive      *bool   `sql:"is_active" json:"is_active"`
	CreatedBy     *int64  `sql:"created_by" json:"created_by"`
	CreatedAt     *string `sql:"created_at" json:"created_at"`
	UpdatedBy     *int64  `sql:"updated_by" json:"updated_by"`
	UpdatedAt     *string `sql:"updated_at" json:"updated_at"`
}

type OfficialHierarchyList struct {
	CustId        string     `db:"cust_id" json:"cust_id"`
	OfficialType  int        `db:"official_type" json:"official_type"`
	HierarchyCode string     `db:"hierarchy_code" json:"hierarchy_code"`
	IsActive      *bool      `db:"is_active" json:"is_active"`
	IsDel         bool       `db:"is_del" json:"is_del"`
	CreatedBy     *int64     `db:"created_by,omitempty" json:"created_by"`
	CreatedAt     *time.Time `db:"created_at,omitempty" json:"created_at"`
	UpdatedBy     *int64     `db:"updated_by,omitempty" json:"updated_by"`
	UpdatedByName *string    `json:"updated_by_name" db:"updated_by_name"`
	UpdatedAt     *time.Time `db:"updated_at,omitempty" json:"updated_at"`
	DeletedBy     *int64     `db:"deleted_by,omitempty" json:"deleted_by"`
	DeletedAt     *time.Time `db:"deleted_at,omitempty" json:"deleted_at"`
}

type OfficialHierarchyUpdate struct {
	OfficialType  *int       `json:"official_type,omitempty" sql:"official_type"`
	HierarchyCode string     `json:"hierarchy_code" sql:"hierarchy_code"`
	IsActive      *bool      `json:"is_active" sql:"is_active"`
	UpdatedAt     *time.Time `json:"updated_at" sql:"updated_at"`
	UpdatedBy     *int64     `json:"updated_by" sql:"updated_by"`
}
