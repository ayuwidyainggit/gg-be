package entity

import (
	"time"
)

type SubBeatResponse struct {
	SbeatId       int        `json:"sbeat_id"`
	SbeatCode     string     `json:"sbeat_code"`
	SbeatName     string     `json:"sbeat_name"`
	DistrictCode  string     `json:"district_code"`
	DistrictName  string     `json:"district_name"`
	BeatCode      string     `json:"beat_code"`
	BeatName      string     `json:"beat_name"`
	BeatId        int        `json:"beat_id"`
	DistrictId    int        `json:"district_id"`
	IsActive      bool       `json:"is_active"`
	UpdatedBy     *int64     `json:"updated_by"`
	UpdatedAt     *time.Time `json:"updated_at"`
	UpdatedByName *string    `json:"updated_by_name"`
}

type SubBeatLookupResponse struct {
	SbeatId      int    `json:"sbeat_id"`
	SbeatCode    string `json:"sbeat_code"`
	SbeatName    string `json:"sbeat_name"`
	BeatCode     string `json:"beat_code"`
	BeatName     string `json:"beat_name"`
	BeatId       int    `json:"beat_id"`
	DistrictId   int    `json:"district_id"`
	DistrictCode string `json:"district_code"`
	DistrictName string `json:"district_name"`
}

type CreateSubBeatBody struct {
	CustId     string `json:"cust_id" validate:"required,max=10"`
	CreatedBy  int64  `json:"created_by" validate:"required"`
	SbeatCode  string `json:"sbeat_code" validate:"required,max=5,alphanumericSpace"`
	SbeatName  string `json:"sbeat_name" validate:"required,max=50"`
	BeatId     int    `json:"beat_id" validate:"required"`
	DistrictId int    `json:"district_id" validate:"required"`
	IsActive   bool   `json:"is_active"`
}

type DetailSubBeatParams struct {
	SbeatId int `params:"sbeat_id" validate:"required"`
}

type UpdateSubBeatParams struct {
	SbeatId int `params:"sbeat_id" validate:"required"`
}

type DeleteSubBeatParams struct {
	SbeatId int `params:"sbeat_id" validate:"required"`
}

type UpdateSubBeatRequest struct {
	CustId     string `json:"cust_id" validate:"required,max=10"`
	UpdatedBy  int64  `json:"updated_by" validate:"required"`
	SbeatCode  string `json:"sbeat_code,omitempty" validate:"required,max=5,omitempty,alphanumericSpace"`
	SbeatName  string `json:"sbeat_name,omitempty" validate:"max=50,omitempty"`
	BeatId     int    `json:"beat_id,omitempty" validate:""`
	DistrictId int    `json:"district_id,omitempty" validate:""`
	IsActive   *bool  `json:"is_active,omitempty"`
}
