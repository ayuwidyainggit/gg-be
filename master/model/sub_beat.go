package model

import (
	"time"
)

type SubBeat struct {
	CustId        string     `db:"cust_id" json:"cust_id"`
	SbeatId       int        `db:"sbeat_id" json:"sbeat_id"`
	SbeatCode     string     `db:"sbeat_code" json:"sbeat_code"`
	SbeatName     string     `db:"sbeat_name" json:"sbeat_name"`
	BeatId        int        `db:"beat_id" json:"beat_id"`
	BeatCode      *string    `db:"beat_code" json:"beat_code"`
	BeatName      *string    `db:"beat_name" json:"beat_name"`
	DistrictId    int        `db:"district_id" json:"district_id"`
	DistrictCode  *string    `db:"district_code" json:"district_code"`
	DistrictName  *string    `db:"district_name" json:"district_name"`
	IsActive      bool       `db:"is_active" json:"is_active"`
	IsDel         bool       `db:"is_del" json:"is_del"`
	CreatedBy     *int64     `db:"created_by,omitempty" json:"created_by"`
	CreatedAt     *time.Time `db:"created_at,omitempty" json:"created_at"`
	UpdatedBy     *int64     `db:"updated_by,omitempty" json:"updated_by"`
	UpdatedByName *string    `db:"updated_by_name" json:"updated_by_name"`
	UpdatedAt     *time.Time `db:"updated_at,omitempty" json:"updated_at"`
	DeletedBy     *int64     `db:"deleted_by,omitempty" json:"deleted_by"`
	DeletedAt     *time.Time `db:"deleted_at,omitempty" json:"deleted_at"`
}

type SubBeatUpdate struct {
	SubBeatCode *string    `json:"sbeat_code,omitempty" sql:"sbeat_code"`
	SubBeatName *string    `json:"sbeat_name,omitempty" sql:"sbeat_name"`
	BeatId      *int       `json:"beat_id" sql:"beat_id"`
	DistrictId  *int       `json:"district_id" sql:"district_id"`
	IsActive    *bool      `json:"is_active" sql:"is_active"`
	UpdatedAt   *time.Time `json:"updated_at" sql:"updated_at"`
	UpdatedBy   *int64     `json:"updated_by" sql:"updated_by"`
}
