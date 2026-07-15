package model

import "time"

type MTpr struct {
	CustID          string     `db:"cust_id" json:"cust_id"`
	TprID           int64      `db:"tpr_id" json:"tpr_id"`
	TprCode         string     `db:"tpr_code" json:"tpr_code"`
	TprName         *string    `db:"tpr_name" json:"tpr_name"`
	DateStart       *time.Time `db:"date_start" json:"date_start"`
	DateEnd         *time.Time `db:"date_end" json:"date_end"`
	RangeType       *int64     `db:"range_type" json:"range_type"`
	PromoItemType   *int64     `db:"promo_item_type" json:"promo_item_type"`
	IsMultiplePromo *bool      `db:"is_multiple_promo" json:"is_multiple_promo"`
	IsAllOtType     *bool      `db:"is_all_ot_type" json:"is_all_ot_type"`
	IsAllOtGrp      *bool      `db:"is_all_ot_grp" json:"is_all_ot_grp"`
	IsAllSales      *bool      `db:"is_all_sales" json:"is_all_sales"`
	IsAllOt         *bool      `db:"is_all_ot" json:"is_all_ot"`
	IsAllSalesTeam  *bool      `db:"is_all_sales_team" json:"is_all_sales_team"`
	IsAllIndustry   *bool      `db:"is_all_industry" json:"is_all_industry"`
	IsMax           *bool      `db:"is_max" json:"is_max"`
	IsMaxValue      *int64     `db:"is_max_value" json:"is_max_value"`
	Deduction       *float64   `db:"deduction" json:"deduction"`
	MinInvoiceValue *float64   `db:"min_invoice_value" json:"min_invoice_value"`
	IsActive        *bool      `db:"is_active" json:"is_active"`
	CreatedBy       *int64     `db:"created_by" json:"created_by"`
	CreatedAt       *time.Time `db:"created_at" json:"created_at"`
	UpdatedBy       *int64     `db:"updated_by" json:"updated_by"`
	UpdatedAt       *time.Time `db:"updated_at" json:"updated_at"`
	UpdatedByName   *string    `json:"updated_by_name" db:"updated_by_name"`
	IsDel           bool       `db:"is_del" json:"is_del"`
	DeletedBy       *int64     `db:"deleted_by" json:"deleted_by"`
	DeletedAt       *time.Time `db:"deleted_at" json:"deleted_at"`
}
type MTprUpdate struct {
	TprCode         *string    `db:"tpr_code" json:"tpr_code"`
	TprName         *string    `db:"tpr_name" json:"tpr_name"`
	DateStart       *string    `db:"date_start" json:"date_start"`
	DateEnd         *string    `db:"date_end" json:"date_end"`
	RangeType       *int64     `db:"range_type" json:"range_type"`
	PromoItemType   *int64     `db:"promo_item_type" json:"promo_item_type"`
	IsMultiplePromo *bool      `db:"is_multiple_promo" json:"is_multiple_promo"`
	IsAllOtType     *bool      `db:"is_all_ot_type" json:"is_all_ot_type"`
	IsAllOtGrp      *bool      `db:"is_all_ot_grp" json:"is_all_ot_grp"`
	IsAllSales      *bool      `db:"is_all_sales" json:"is_all_sales"`
	IsAllOt         *bool      `db:"is_all_ot" json:"is_all_ot"`
	IsAllSalesTeam  *bool      `db:"is_all_sales_team" json:"is_all_sales_team"`
	IsAllIndustry   *bool      `db:"is_all_industry" json:"is_all_industry"`
	IsMax           *bool      `db:"is_max" json:"is_max"`
	IsMaxValue      *int64     `db:"is_max_value" json:"is_max_value"`
	Deduction       *float64   `db:"deduction" json:"deduction"`
	MinInvoiceValue *float64   `db:"min_invoice_value" json:"min_invoice_value"`
	IsActive        *bool      `db:"is_active" json:"is_active"`
	DeletedBy       *int64     `db:"deleted_by" json:"deleted_by"`
	DeletedAt       *time.Time `db:"deleted_at" json:"deleted_at"`
	UpdatedAt       *time.Time `json:"updated_at" sql:"updated_at"`
	UpdatedBy       *int64     `json:"updated_by" sql:"updated_by"`
}
