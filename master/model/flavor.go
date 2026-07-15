package model

import (
	"time"
)

type Flavor struct {
	CustId        string     `db:"cust_id" json:"cust_id"`
	FlavorId      int        `db:"flavor_id" json:"flavor_id"`
	FlavorCode    string     `db:"flavor_code" json:"flavor_code"`
	FlavorName    string     `db:"flavor_name" json:"flavor_name"`
	IsActive      bool       `db:"is_active" json:"is_active"`
	IsDel         bool       `db:"is_del" json:"is_del"`
	CreatedBy     *int64     `db:"created_by,omitempty" json:"created_by"`
	CreatedAt     *time.Time `db:"created_at,omitempty" json:"created_at"`
	UpdatedBy     *int64     `db:"updated_by,omitempty" json:"updated_by"`
	UpdatedAt     *time.Time `db:"updated_at,omitempty" json:"updated_at"`
	DeletedBy     *int64     `db:"deleted_by,omitempty" json:"deleted_by"`
	DeletedAt     *time.Time `db:"deleted_at,omitempty" json:"deleted_at"`
	UpdatedByName *string    `db:"updated_by_name" json:"updated_by_name"`
}

type FlavorUpdate struct {
	FlavorCode *string    `json:"flavor_code,omitempty" sql:"flavor_code"`
	FlavorName *string    `json:"flavor_name,omitempty" sql:"flavor_name"`
	IsActive   *bool      `json:"is_active" sql:"is_active"`
	UpdatedAt  *time.Time `json:"updated_at" sql:"updated_at"`
	UpdatedBy  *int64     `json:"updated_by" sql:"updated_by"`
}
