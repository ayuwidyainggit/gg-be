package entity

import "time"

type CreatedTprLimitBody struct {
	CustId       string  `json:"cust_id"`
	ProId        int     `json:"pro_id"`
	TprType      int64   `json:"tpr_type"`
	DateStart    *string `json:"date_start"`
	DateEnd      *string `json:"date_end"`
	ValueLimit   float64 `json:"value_limit"`
	ValueUsed    float64 `json:"value_used"`
	ValueUsedStr string  `json:"value_used_str"`
	VatType      int64   `json:"vat_type"`
	CreatedBy    int64   `json:"created_by"`
	UpdatedBy    int64   `json:"updated_by"`
}
type TprLimitResponse struct {
	TprLimitId   int       `json:"tpr_limit_id"`
	ProId        int       `json:"pro_id"`
	TprType      int64     `json:"tpr_type"`
	DateStart    string    `json:"date_start"`
	DateEnd      string    `json:"date_end"`
	ValueLimit   float64   `json:"value_limit"`
	ValueUsed    float64   `json:"value_used"`
	ValueUsedStr string    `json:"value_used_str"`
	VatType      int64     `json:"vat_type"`
	UpdatedBy    int64     `json:"updated_by"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type TprLimitListResponse struct {
	TprLimitId    int       `json:"tpr_limit_id"`
	ProId         int       `json:"pro_id"`
	ProCode       string    `json:"pro_code"`
	ProName       string    `json:"pro_name"`
	TprType       int64     `json:"tpr_type"`
	DateStart     string    `json:"date_start"`
	DateEnd       string    `json:"date_end"`
	ValueLimit    float64   `json:"value_limit"`
	ValueUsed     float64   `json:"value_used"`
	ValueUsedStr  string    `json:"value_used_str"`
	VatType       int64     `json:"vat_type"`
	UpdatedAt     time.Time `json:"updated_at"`
	UpdatedByName string    `json:"updated_by_name"`
}

type DetailTprLimitParams struct {
	TprLimitId int `params:"tpr_limit_id" validate:"required"`
}

type UpdateTprLimitParams struct {
	TprLimitId int `params:"tpr_limit_id" validate:"required"`
}

type DeleteTprLimitParams struct {
	TprLimitId int `params:"tpr_limit_id" validate:"required"`
}

type UpdateTprLimitRequest struct {
	CustId       string  `json:"cust_id" validate:"required,max=10"`
	ProId        int     `json:"pro_id"`
	TprType      int64   `json:"tpr_type" validate:"max=5"`
	DateStart    *string `json:"date_start"`
	DateEnd      *string `json:"date_end"`
	ValueLimit   float64 `json:"value_limit"`
	ValueUsed    float64 `json:"value_used"`
	ValueUsedStr string  `json:"value_used_str"`
	VatType      int64   `json:"vat_type" validate:"max=5"`
	UpdatedBy    int64   `json:"updated_by"`
}
