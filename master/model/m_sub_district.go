package model

import "time"

type SubDistrict struct {
	CustId        string     `db:"cust_id" json:"cust_id"`
	SubDistrictId string     `db:"sub_district_id" json:"sub_district_id"`
	SubDistrict   string     `db:"sub_district" json:"sub_district"`
	ProvinceId    string     `db:"province_id" json:"province_id"`
	Province      *string    `db:"province" json:"province"`
	RegencyId     string     `db:"regency_id" json:"regency_id"`
	Regency       *string    `db:"regency" json:"regency"`
	IsActive      bool       `db:"is_active" json:"is_active"`
	CreatedBy     *int64     `db:"created_by,omitempty" json:"created_by"`
	CreatedAt     *time.Time `db:"created_at,omitempty" json:"created_at"`
	UpdatedBy     *int64     `db:"updated_by,omitempty" json:"updated_by"`
	UpdatedAt     *time.Time `db:"updated_at,omitempty" json:"updated_at"`
	UpdatedByName *string    `db:"updated_by_name" json:"updated_by_name"`
}

type SubDistrictUpdate struct {
	SubDistrictId *string    `json:"sub_district_id,omitempty" sql:"sub_district_id"`
	SubDistrict   *string    `json:"sub_district,omitempty" sql:"sub_district"`
	ProvinceId    *string    `json:"province_id,omitempty" sql:"province_id"`
	RegencyId     *string    `json:"regency_id,omitempty" sql:"regency_id"`
	IsActive      *bool      `json:"is_active" sql:"is_active"`
	UpdatedAt     *time.Time `json:"updated_at" sql:"updated_at"`
	UpdatedBy     *int64     `json:"updated_by" sql:"updated_by"`
}
