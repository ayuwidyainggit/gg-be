package entity

type MOutletSalesman struct {
	CustID   string `json:"cust_id"`
	OutletId int    `json:"outlet_id"`
	SalesId  *int   `json:"sales_id"`
	W1       bool   `json:"w1"`
	W2       bool   `json:"w2"`
	W3       bool   `json:"w3"`
	W4       bool   `json:"w4"`
	RouteId  *int   `json:"route_id"`
	DayId    *int   `json:"day_id"`
}

type CreateOutletSalesmanBody struct {
	SalesId *int `json:"sales_id"`
	W1      bool `json:"w1"`
	W2      bool `json:"w2"`
	W3      bool `json:"w3"`
	W4      bool `json:"w4"`
	RouteId *int `json:"route_id"`
	DayId   *int `json:"day_id"`
}

type UpdateOutletSalesmanBody struct {
	OutletSalesId *int64 `json:"outlet_sales_id"`
	SalesId       *int   `json:"sales_id"`
	W1            *bool  `json:"w1"`
	W2            *bool  `json:"w2"`
	W3            *bool  `json:"w3"`
	W4            *bool  `json:"w4"`
	RouteId       *int   `json:"route_id"`
	DayId         *int   `json:"day_id"`
}
