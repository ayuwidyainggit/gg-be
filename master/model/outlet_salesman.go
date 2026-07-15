package model

type MOutletSalesman struct {
	CustID        string `db:"cust_id" json:"cust_id"`
	OutletID      int64  `db:"outlet_id" json:"outlet_id"`
	SalesID       int64  `db:"sales_id" json:"sales_id"`
	W1            bool   `db:"w1" json:"w1"`
	W2            bool   `db:"w2" json:"w2"`
	W3            bool   `db:"w3" json:"w3"`
	W4            bool   `db:"w4" json:"w4"`
	RouteID       int    `db:"route_id" json:"route_id"`
	DayID         int    `db:"day_id" json:"day_id"`
	OutletSalesId *int64 `db:"outlet_sales_id" json:"outlet_sales_id"`
}

type OutletSalesmanList struct {
	CustID   string `db:"cust_id" json:"cust_id"`
	OutletID int    `db:"outlet_id" json:"outlet_id"`
	SalesID  int    `db:"sales_id" json:"sales_id"`
	W1       bool   `db:"w1" json:"w1"`
	W2       bool   `db:"w2" json:"w2"`
	W3       bool   `db:"w3" json:"w3"`
	W4       bool   `db:"w4" json:"w4"`
	RouteID  int    `db:"route_id" json:"route_id"`
	DayID    int    `db:"day_id" json:"day_id"`
}

type MOutletSalesmanUpdate struct {
	SalesID *int64 `db:"sales_id" json:"sales_id"`
	W1      *bool  `db:"w1" json:"w1"`
	W2      *bool  `db:"w2" json:"w2"`
	W3      *bool  `db:"w3" json:"w3"`
	W4      *bool  `db:"w4" json:"w4"`
	RouteID *int   `db:"route_id" json:"route_id"`
	DayID   *int   `db:"day_id" json:"day_id"`
}
type MOutletSalesmanRead struct {
	CustID        string  `db:"cust_id" json:"cust_id"`
	OutletID      int64   `db:"outlet_id" json:"outlet_id"`
	SalesID       int64   `db:"sales_id" json:"sales_id"`
	SalesmanName  *string `db:"sales_name" json:"sales_name"`
	W1            bool    `db:"w1" json:"w1"`
	W2            bool    `db:"w2" json:"w2"`
	W3            bool    `db:"w3" json:"w3"`
	W4            bool    `db:"w4" json:"w4"`
	RouteID       int     `db:"route_id" json:"route_id"`
	DayID         int     `db:"day_id" json:"day_id"`
	DayName       string  `db:"day_name" json:"day_name"`
	OutletSalesId *int64  `db:"outlet_sales_id" json:"outlet_sales_id"`
}
