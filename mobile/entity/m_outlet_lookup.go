package entity

type OutletLookupListQuery struct {
	Page     int    `query:"page"`
	Limit    int    `query:"limit"`
	Order    string `query:"order"`
	CustID   string `query:"cust_id"`
	IsActive *int   `query:"is_active"`
	Search   string `query:"q"`
}
