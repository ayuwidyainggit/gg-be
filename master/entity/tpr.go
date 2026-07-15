package entity

import "time"

type TprResponse struct {
	TprID           int64     `json:"tpr_id"`
	TprCode         string    `json:"tpr_code"`
	TprName         string    `json:"tpr_name"`
	DateStart       *string   `json:"date_start"`
	DateEnd         *string   `json:"date_end"`
	RangeType       int64     `json:"range_type"`
	PromoItemType   int64     `json:"promo_item_type"`
	IsMultiplePromo bool      `json:"is_multiple_promo"`
	IsAllOtType     bool      `json:"is_all_ot_type"`
	IsAllOtGrp      bool      `json:"is_all_ot_grp"`
	IsAllSales      bool      `json:"is_all_sales"`
	IsAllOt         bool      `json:"is_all_ot"`
	IsAllSalesTeam  bool      `json:"is_all_sales_team"`
	IsAllIndustry   bool      `json:"is_all_industry"`
	IsMax           bool      `json:"is_max"`
	IsMaxValue      int64     `json:"is_max_value"`
	Deduction       float64   `json:"deduction"`
	MinInvoiceValue float64   `json:"min_invoice_value"`
	IsActive        bool      `json:"is_active"`
	UpdatedBy       int64     `json:"updated_by"`
	UpdatedAt       time.Time `json:"updated_at"`
}
type TprListResponse struct {
	TprID           int64     `json:"tpr_id"`
	TprCode         string    `json:"tpr_code"`
	TprName         string    `json:"tpr_name"`
	DateStart       *string   `json:"date_start"`
	DateEnd         *string   `json:"date_end"`
	RangeType       int64     `json:"range_type"`
	PromoItemType   int64     `json:"promo_item_type"`
	IsMultiplePromo bool      `json:"is_multiple_promo"`
	IsAllOtType     bool      `json:"is_all_ot_type"`
	IsAllOtGrp      bool      `json:"is_all_ot_grp"`
	IsAllSales      bool      `json:"is_all_sales"`
	IsAllOt         bool      `json:"is_all_ot"`
	IsAllSalesTeam  bool      `json:"is_all_sales_team"`
	IsAllIndustry   bool      `json:"is_all_industry"`
	IsMax           bool      `json:"is_max"`
	IsMaxValue      int64     `json:"is_max_value"`
	Deduction       float64   `json:"deduction"`
	MinInvoiceValue float64   `json:"min_invoice_value"`
	IsActive        bool      `json:"is_active"`
	UpdatedBy       int64     `json:"updated_by"`
	UpdatedAt       time.Time `json:"updated_at"`
	UpdatedByName   string    `json:"updated_by_name"`
}
type DetailTprParams struct {
	TprId int64 `params:"tpr_id" validate:"required"`
}
type UpdateTprParams struct {
	TprId int64 `params:"tpr_id" validate:"required"`
}

type CreateTprBody struct {
	CustID          string  `json:"cust_id"`
	TprID           int     `json:"tpr_id"`
	TprCode         string  `json:"tpr_code" validate:"required,max=50,alphanumericSpace"`
	TprName         string  `json:"tpr_name" validate:"max=100"`
	DateStart       *string `json:"date_start,omitempty"`
	DateEnd         *string `json:"date_end,omitempty"`
	RangeType       int64   `json:"range_type" validate:"max=5"`
	PromoItemType   int64   `json:"promo_item_type" validate:"max=5"`
	IsMultiplePromo bool    `json:"is_multiple_promo"`
	IsAllOtType     bool    `json:"is_all_ot_type"`
	IsAllOtGrp      bool    `json:"is_all_ot_grp"`
	IsAllSales      bool    `json:"is_all_sales"`
	IsAllOt         bool    `json:"is_all_ot"`
	IsAllSalesTeam  bool    `json:"is_all_sales_team"`
	IsAllIndustry   bool    `json:"is_all_industry"`
	IsMax           bool    `json:"is_max"`
	IsMaxValue      int64   `json:"is_max_value"`
	Deduction       float64 `json:"deduction"`
	MinInvoiceValue float64 `json:"min_invoice_value"`
	IsActive        bool    `json:"is_active"`
	CreatedBy       int64   `json:"created_by"`
}

type UpdateTprRequest struct {
	CustId          string  `json:"cust_id" validate:"required,max=10"`
	TprID           int     `json:"tpr_id"`
	TprCode         string  `json:"tpr_code" validate:"required,max=50,alphanumericSpace"`
	TprName         string  `json:"tpr_name" validate:"max=100"`
	DateStart       *string `json:"date_start"`
	DateEnd         *string `json:"date_end"`
	RangeType       int64   `json:"range_type" validate:"max=5"`
	PromoItemType   int64   `json:"promo_item_type" validate:"max=5"`
	IsMultiplePromo bool    `json:"is_multiple_promo"`
	IsAllOtType     bool    `json:"is_all_ot_type"`
	IsAllOtGrp      bool    `json:"is_all_ot_grp"`
	IsAllSales      bool    `json:"is_all_sales"`
	IsAllOt         bool    `json:"is_all_ot"`
	IsAllSalesTeam  bool    `json:"is_all_sales_team"`
	IsAllIndustry   bool    `json:"is_all_industry"`
	IsMax           bool    `json:"is_max"`
	IsMaxValue      int64   `json:"is_max_value"`
	Deduction       float64 `json:"deduction"`
	MinInvoiceValue float64 `json:"min_invoice_value"`
	IsActive        bool    `json:"is_active"`
	CreatedBy       int64   `json:"created_by"`
	UpdatedBy       int64   `json:"updated_by" validate:"required"`
}

type DeleteTprParams struct {
	TprId int `params:"tpr_id" validate:"required"`
}
