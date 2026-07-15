package entity

import (
	"time"
)

type BeatResponse struct {
	BeatId        int        `json:"beat_id"`
	DistrictCode  string     `json:"district_code"`
	DistrictName  string     `json:"district_name"`
	BeatCode      string     `json:"beat_code"`
	BeatName      string     `json:"beat_name"`
	DistrictId    int        `json:"district_id"`
	IsActive      bool       `json:"is_active"`
	UpdatedBy     *int64     `json:"updated_by"`
	UpdatedAt     *time.Time `json:"updated_at"`
	UpdatedByName *string    `json:"updated_by_name"`
}

type BeatLookupResponse struct {
	BeatId       int    `json:"beat_id"`
	BeatCode     string `json:"beat_code"`
	BeatName     string `json:"beat_name"`
	DistrictId   int    `json:"district_id"`
	DistrictCode string `json:"district_code"`
	DistrictName string `json:"district_name"`
}

type CreateBeatBody struct {
	CustId     string `json:"cust_id" validate:"required,max=10"`
	CreatedBy  int64  `json:"created_by" validate:"required"`
	BeatCode   string `json:"beat_code" validate:"required,max=10,alphanumericSpace"`
	BeatName   string `json:"beat_name" validate:"required,max=150"`
	DistrictId int    `json:"district_id" validate:"required"`
	IsActive   bool   `json:"is_active"`
}

type DetailBeatParams struct {
	BeatId int `params:"beat_id" validate:"required"`
}

type UpdateBeatParams struct {
	BeatId int `params:"beat_id" validate:"required"`
}

type DeleteBeatParams struct {
	BeatId int `params:"beat_id" validate:"required"`
}

type UpdateBeatRequest struct {
	CustId     string `json:"cust_id" validate:"required,max=10"`
	DistrictId int    `json:"district_id" validate:"required"`
	UpdatedBy  int64  `json:"updated_by" validate:"required"`
	BeatCode   string `json:"beat_code,omitempty" validate:"required,max=10,alphanumericSpace"`
	BeatName   string `json:"beat_name,omitempty" validate:"max=150,omitempty"`
	IsActive   *bool  `json:"is_active,omitempty"`
}
