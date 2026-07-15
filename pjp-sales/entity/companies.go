package entity

type CompaniesQueryFilter struct {
	CustId       string
	ParentCustId string
	CompanyIds   []string `query:"company_ids"`

	Query string `query:"q"`
}

type CompaniesListResponse struct {
	HeadOffice  bool   `json:"head_office"`
	CompanyID   string `json:"company_id"`
	CompanyCode string `json:"company_code"`
	CompanyName string `json:"company_name"`
}
