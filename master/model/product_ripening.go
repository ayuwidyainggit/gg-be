package model

import "time"

type ProductRipening struct {
	ID              int64      `db:"id"`
	CustID          string     `db:"cust_id"`
	DistributorID   int64      `db:"distributor_id"`
	ProID           int64      `db:"pro_id"`
	PerYear         int        `db:"per_year"`
	PerID           int        `db:"per_id"`
	WeekID          int        `db:"week_id"`
	SundayQty       int        `db:"sunday_qty"`
	MondayQty       int        `db:"monday_qty"`
	TuesdayQty      int        `db:"tuesday_qty"`
	WednesdayQty    int        `db:"wednesday_qty"`
	ThursdayQty     int        `db:"thursday_qty"`
	FridayQty       int        `db:"friday_qty"`
	SaturdayQty     int        `db:"saturday_qty"`
	CreatedBy       int64      `db:"created_by"`
	CreatedAt       time.Time  `db:"created_at"`
	UpdatedBy       *int64     `db:"updated_by"`
	UpdatedAt       *time.Time `db:"updated_at"`
	DeletedBy       *int64     `db:"deleted_by"`
	DeletedAt       *time.Time `db:"deleted_at"`
	IsDel           bool       `db:"is_del"`
	DistributorCode string     `db:"distributor_code"`
	DistributorName string     `db:"distributor_name"`
	ProductCode     string     `db:"product_code"`
	ProductName     string     `db:"product_name"`
	WeekStart       string     `db:"week_start"`
	WeekEnd         string     `db:"week_end"`
	CreatedByName   string     `db:"created_by_name"`
	UpdatedByName   *string    `db:"updated_by_name"`
}

type ProductRipeningPlanListRow struct {
	ID              int64      `db:"id"`
	CustID          string     `db:"cust_id"`
	DistributorID   int64      `db:"distributor_id"`
	DistributorCode string     `db:"distributor_code"`
	DistributorName string     `db:"distributor_name"`
	PerYear         int        `db:"per_year"`
	PerID           int        `db:"per_id"`
	WeekID          int        `db:"week_id"`
	WeekStart       string     `db:"week_start"`
	WeekEnd         string     `db:"week_end"`
	TotalProduct    int        `db:"total_product"`
	CreatedBy       int64      `db:"created_by"`
	CreatedByName   string     `db:"created_by_name"`
	CreatedAt       time.Time  `db:"created_at"`
	UpdatedBy       *int64     `db:"updated_by"`
	UpdatedByName   *string    `db:"updated_by_name"`
	UpdatedAt       *time.Time `db:"updated_at"`
}

type ProductRipeningWeek struct {
	CustID    string `db:"cust_id"`
	PerYear   int    `db:"per_year"`
	PerID     int    `db:"per_id"`
	WeekID    int    `db:"week_id"`
	WeekStart string `db:"week_start"`
	WeekEnd   string `db:"week_end"`
}

type ProductRipeningAssignedDistributor struct {
	DistributorID     int64  `db:"distributor_id"`
	DistributorCode   string `db:"distributor_code"`
	DistributorName   string `db:"distributor_name"`
	DistributorCustID string `db:"dist_cust_id"`
}
