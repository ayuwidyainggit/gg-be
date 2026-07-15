package entity

import "time"

type DiscProductResponse struct {
	DiscId    int        `json:"disc_id"`
	DiscCode  string     `json:"disc_code"`
	DiscName  string     `json:"disc_name"`
	ProId     int        `json:"pro_id"`
	ProCode   string     `json:"pro_code"`
	ProName   string     `json:"pro_name"`
	MinQty    float64    `json:"min_qty"`
	MinQtyStr string     `json:"min_qty_str"`
	MaxQty    float64    `json:"max_qty"`
	MaxQtyStr string     `json:"max_qty_str"`
	DiscPerc  int        `json:"disc_perc"`
	IsActive  bool       `json:"is_active"`
	UpdatedBy *int64     `json:"updated_by"`
	UpdatedAt *time.Time `json:"updated_at"`
}

type DiscProductListResponse struct {
	DiscId        int        `json:"disc_id"`
	DiscCode      string     `json:"disc_code"`
	DiscName      string     `json:"disc_name"`
	ProId         int        `json:"pro_id"`
	ProCode       string     `json:"pro_code"`
	ProName       string     `json:"pro_name"`
	MinQty        float64    `json:"min_qty"`
	MinQtyStr     string     `json:"min_qty_str"`
	MaxQty        float64    `json:"max_qty"`
	MaxQtyStr     string     `json:"max_qty_str"`
	DiscPerc      int        `json:"disc_perc"`
	IsActive      bool       `json:"is_active"`
	UpdatedBy     *int64     `json:"updated_by"`
	UpdatedAt     *time.Time `json:"updated_at"`
	UpdatedByName string     `json:"updated_by_name"`
}

type CreateDiscProductBody struct {
	CustId    string  `json:"cust_id"`
	DiscId    int     `json:"disc_id"`
	ProId     int     `json:"pro_id"`
	MinQty    float64 `json:"min_qty"`
	MinQtyStr string  `json:"min_qty_str" validate:"required,qtystr"`
	MaxQty    float64 `json:"max_qty"`
	MaxQtyStr string  `json:"max_qty_str" validate:"required,qtystr"`
	DiscPerc  int     `json:"disc_perc"`
	IsActive  bool    `json:"is_active"`
	CreatedBy int64   `json:"created_by" validate:"required"`
	UpdatedBy int64   `json:"updated_by"`
}

type DetailDiscProductParams struct {
	DiscId int `params:"disc_id" validate:"required"`
	ProId  int `params:"pro_id" validate:"required"`
}

type UpdateDiscProductParams struct {
	DiscId int `params:"disc_id" validate:"required"`
	ProId  int `params:"pro_id" validate:"required"`
}

type DeleteDiscProductParams struct {
	DiscId int `params:"disc_id" validate:"required"`
	ProId  int `params:"pro_id" validate:"required"`
}

type UpdateDiscProductRequest struct {
	CustId    string  `json:"cust_id"`
	DiscId    int     `json:"disc_id" validate:"required"`
	ProId     int     `json:"pro_id" validate:"required"`
	MinQty    float64 `json:"min_qty"`
	MinQtyStr string  `json:"min_qty_str"`
	MaxQty    float64 `json:"max_qty"`
	MaxQtyStr string  `json:"max_qty_str"`
	DiscPerc  int     `json:"disc_perc"`
	IsActive  *bool   `json:"is_active,omitempty"`
	UpdatedBy int64   `json:"updated_by" validate:"required"`
}
