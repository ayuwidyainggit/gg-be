package entity

type SummaryQueryFilter struct {
	SalesmanCode string `json:"salesman_code"`
	CustID       string `json:"cust_id"`
	Date         string `json:"date"`
}
