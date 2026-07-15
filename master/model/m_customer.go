package model

type MCustomer struct {
	CustId       string `db:"cust_id" json:"cust_id"`
	CustName     string `db:"cust_name" json:"cust_name"`
	ParentCustId string `db:"parent_cust_id" json:"parent_cust_id"`
}
