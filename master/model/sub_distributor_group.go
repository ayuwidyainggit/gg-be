package model

import "time"

type SubDistributorGroup struct {
	CustID                  string     `db:"cust_id" json:"cust_id"`
	SubDistributorGroupID   int64      `db:"sub_distributor_group_id" json:"sub_distributor_group_id"`
	SubDistributorGroupCode string     `db:"sub_distributor_group_code" json:"sub_distributor_group_code"`
	SubDistributorGroupName string     `db:"sub_distributor_group_name" json:"sub_distributor_group_name"`
	IsActive                bool       `db:"is_active" json:"is_active"`
	CreatedBy               *int64     `db:"created_by" json:"created_by"`
	UpdatedBy               *int64     `db:"updated_by" json:"updated_by"`
	CreatedAt               *time.Time `db:"created_at" json:"created_at"`
	UpdatedAt               *time.Time `db:"updated_at" json:"updated_at"`
	IsDel                   bool       `db:"is_del" json:"is_del"`
	DeletedBy               *int64     `db:"deleted_by,omitempty" json:"deleted_by"`
	DeletedAt               *time.Time `db:"deleted_at,omitempty" json:"deleted_at"`
	UpdatedByName           *string    `db:"updated_by_name" json:"updated_by_name"`
}

type SubDistributorGroupUpdate struct {
	SubDistributorGroupCode *string    `db:"sub_distributor_group_code" json:"sub_distributor_group_code"`
	SubDistributorGroupName *string    `db:"sub_distributor_group_name" json:"sub_distributor_group_name"`
	IsActive                *bool      `db:"is_active" json:"is_active"`
	UpdatedAt               *time.Time `json:"updated_at" sql:"updated_at"`
	UpdatedBy               *int64     `json:"updated_by" sql:"updated_by"`
}
