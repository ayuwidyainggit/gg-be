package model

import "time"

type Province struct {
	CustId        string     `db:"cust_id" json:"cust_id"`
	ProvinceId    string     `db:"province_id" json:"province_id"`
	Province      string     `db:"province" json:"province"`
	IsActive      bool       `db:"is_active" json:"is_active"`
	CreatedBy     *int64     `db:"created_by,omitempty" json:"created_by"`
	CreatedAt     *time.Time `db:"created_at,omitempty" json:"created_at"`
	UpdatedBy     *int64     `db:"updated_by,omitempty" json:"updated_by"`
	UpdatedAt     *time.Time `db:"updated_at,omitempty" json:"updated_at"`
	UpdatedByName *string    `db:"updated_by_name" json:"updated_by_name"`
}

type ProvinceUpdate struct {
	ProvinceId *string    `json:"province_id,omitempty" sql:"province_id"`
	Province   *string    `json:"province,omitempty" sql:"province"`
	IsActive   *bool      `json:"is_active" sql:"is_active"`
	UpdatedAt  *time.Time `json:"updated_at" sql:"updated_at"`
	UpdatedBy  *int64     `json:"updated_by" sql:"updated_by"`
}
