package entity

import (
	"time"
)

type MarketResponse struct {
	MarketId      int        `json:"market_id"`
	MarketCode    string     `json:"market_code"`
	MarketName    string     `json:"market_name"`
	IsActive      bool       `json:"is_active"`
	UpdatedBy     *int64     `json:"updated_by"`
	UpdatedAt     *time.Time `json:"updated_at"`
	UpdatedByName *string    `json:"updated_by_name"`
}

type MarketLookupResponse struct {
	MarketId   int    `json:"market_id"`
	MarketCode string `json:"market_code"`
	MarketName string `json:"market_name"`
}

type CreateMarketBody struct {
	CustId     string `json:"cust_id" validate:"required,max=10"`
	CreatedBy  int64  `json:"created_by" validate:"required"`
	MarketCode string `json:"market_code" validate:"required,max=5,alphanumericSpace"`
	MarketName string `json:"market_name" validate:"required,max=40,alphanumericSpace"`
	IsActive   bool   `json:"is_active"`
}

type DetailMarketParams struct {
	MarketId int `params:"market_id" validate:"required"`
}

type UpdateMarketParams struct {
	MarketId int `params:"market_id" validate:"required"`
}

type DeleteMarketParams struct {
	MarketId int `params:"market_id" validate:"required"`
}

type UpdateMarketRequest struct {
	CustId     string `json:"cust_id" validate:"required,max=10"`
	UpdatedBy  int64  `json:"updated_by" validate:"required"`
	MarketCode string `json:"market_code,omitempty" validate:"required,max=5,alphanumericSpace"`
	MarketName string `json:"market_name,omitempty" validate:"max=40,alphanumericSpace"`
	IsActive   *bool  `json:"is_active,omitempty"`
}
