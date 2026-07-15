package entity

import "time"

type CreateChannelBody struct {
	ChannelID   int    `json:"channel_id"`
	CustID      string `json:"cust_id" validate:"required,max=10"`
	ChannelCode string `json:"channel_code" validate:"required,max=4,numeric"`
	ChannelName string `json:"channel_name" validate:"required,max=20,alphanumericSpace"`
	IsActive    bool   `json:"is_active"`
	CreatedBy   int64  `json:"created_by" validate:"required"`
}

type UpdateChannelBody struct {
	CustID      string `json:"cust_id" validate:"required,max=10"`
	ChannelCode string `json:"channel_code" validate:"required,max=4,numeric"`
	ChannelName string `json:"channel_name" validate:"required,max=20,alphanumericSpace"`
	IsActive    bool   `json:"is_active"`
	UpdatedBy   int64  `json:"updated_by" validate:"required"`
}
type ChannelListResponse struct {
	ChannelID     int        `json:"channel_id"`
	CustID        string     `json:"cust_id"`
	ChannelCode   string     `json:"channel_code" `
	ChannelName   string     `json:"channel_name"`
	IsActive      bool       `json:"is_active"`
	UpdatedBy     *int64     `json:"updated_by"`
	UpdatedAt     *time.Time `json:"updated_at"`
	UpdatedByName string     `json:"updated_by_name"`
}

type DetailChannelParams struct {
	ChannelId int `params:"channel_id" validate:"required"`
}

type UpdateChannelParams struct {
	ChannelId int `params:"channel_id" validate:"required"`
}

type DeleteChannelParams struct {
	ChannelId int `params:"channel_id" validate:"required"`
}

type ChannelResponse struct {
	ChannelID     int        `json:"channel_id"`
	CustID        string     `json:"cust_id"`
	ChannelCode   string     `json:"channel_code" `
	ChannelName   string     `json:"channel_name"`
	IsActive      bool       `json:"is_active"`
	UpdatedBy     *int64     `json:"updated_by"`
	UpdatedAt     *time.Time `json:"updated_at"`
	UpdatedByName string     `json:"updated_by_name"`
}

type ChannelLookupResponse struct {
	ChannelID   int    `json:"channel_id"`
	ChannelCode string `json:"channel_code" `
	ChannelName string `json:"channel_name"`
}

type ChannelQueryFilter struct {
	Page      int    `query:"page"`
	Limit     int    `query:"limit" validate:"required"`
	Query     string `query:"q"`
	Mode      string `query:"mode"`
	Sort      string `query:"sort"`
	ChannelID int    `query:"channel_id"`
	IsActive  *int   `query:"is_active"`
}

type ChannelUpdateRequest struct {
	CustID      string `json:"cust_id" validate:"required,max=10"`
	ChannelCode string `json:"channel_code" validate:"required"`
	ChannelName string `json:"channel_name" validate:"required"`
	UpdatedBy   int64  `json:"updated_by"`
	IsActive    *bool  `json:"is_active"`
}
