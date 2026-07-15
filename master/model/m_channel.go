package model

import "time"

type MChannel struct {
	ChannelID     int64      `db:"channel_id" json:"channel_id"`
	CustID        string     `db:"cust_id" json:"cust_id"`
	ChannelCode   string     `db:"channel_code" json:"channel_code"`
	ChannelName   string     `db:"channel_name" json:"channel_name"`
	IsActive      bool       `db:"is_active" json:"is_active"`
	CreatedBy     *int64     `db:"created_by" json:"created_by"`
	UpdatedBy     *int64     `db:"updated_by" json:"updated_by"`
	CreatedAt     *time.Time `db:"created_at" json:"created_at"`
	UpdatedAt     *time.Time `db:"updated_at" json:"updated_at"`
	IsDel         bool       `db:"is_del" json:"is_del"`
	DeletedBy     *int64     `db:"deleted_by,omitempty" json:"deleted_by"`
	DeletedAt     *time.Time `db:"deleted_at,omitempty" json:"deleted_at"`
	UpdatedByName *string    `db:"updated_by_name" json:"updated_by_name"`
}

type MChannelUpdate struct {
	ChannelCode *string    `db:"channel_code" json:"channel_code"`
	ChannelName *string    `db:"channel_name" json:"channel_name"`
	IsActive    *bool      `db:"is_active" json:"is_active"`
	UpdatedAt   *time.Time `json:"updated_at" sql:"updated_at"`
	UpdatedBy   *int64     `json:"updated_by" sql:"updated_by"`
}
