package entity

type WhStockQuery struct {
	CustID        string   `json:"cust_id"`
	ParentCustID  string   `json:"parent_cust_id"`
	WhId          int64    `json:"wh_id"`
	ProId         int64    `json:"pro_id"`
	Qty           *float64 `json:"qty"`
	QtyOnOrder    *float64 `json:"qty_on_order"`
	QtyOnShipping *float64 `json:"qty_on_shipping"`
	QtyBs         *float64 `json:"qty_bs"`
	QtyExp        *float64 `json:"qty_exp"`
}
