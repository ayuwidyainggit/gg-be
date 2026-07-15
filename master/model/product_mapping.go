package model

import "time"

type ProductMappingListRow struct {
	DistributorID   int64      `db:"distributor_id"`
	DistributorCode string     `db:"distributor_code"`
	DistributorName string     `db:"distributor_name"`
	TotalProduct    int        `db:"total_product"`
	CreatedBy       *int64     `db:"created_by"`
	CreatedByName   *string    `db:"created_by_name"`
	CreatedAt       *time.Time `db:"created_at"`
	UpdatedBy       *int64     `db:"updated_by"`
	UpdatedByName   *string    `db:"updated_by_name"`
	UpdatedAt       *time.Time `db:"updated_at"`
}

type ProductMappingDetailRow struct {
	ProID         int64   `db:"pro_id"`
	ParentProID   int64   `db:"parent_pro_id"`
	ParentProCode string  `db:"parent_pro_code"`
	ParentProName string  `db:"parent_pro_name"`
	ProCode       string  `db:"pro_code"`
	ProName       string  `db:"pro_name"`
	LargestUOM    string  `db:"largest_uom"`
	MiddleUOM     *string `db:"middle_uom"`
	SmallestUOM   *string `db:"smallest_uom"`
}

type ProductMappingProductRow struct {
	ProID           int64  `db:"pro_id"`
	CustID          string `db:"cust_id"`
	DistributorID   int64  `db:"distributor_id"`
	ParentProID     int64  `db:"parent_pro_id"`
	ProCode         string `db:"pro_code"`
	ProName         string `db:"pro_name"`
	IsProductMapping bool  `db:"is_product_mapping"`
}
