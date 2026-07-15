package entity

type MCustomerResp struct {
	CustId       string `json:"cust_id"`
	CustName     string `json:"cust_name"`
	ParentCustId string `json:"parent_cust_id"`
}
