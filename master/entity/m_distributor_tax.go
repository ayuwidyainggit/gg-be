package entity

type DistributorTax struct {
	DistributorTaxId    int64  `json:"distributor_tax_id"`
	DistributorId       int64  `json:"distributor_id"`
	TaxIdentifierNoType string `json:"tax_identifier_no_type"`
	TaxIdentifierNo     string `json:"tax_identifier_no"`
	Nitku               string `json:"nitku"`
	TaxName             string `json:"tax_name"`
	TaxAddress          string `json:"tax_address"`
}

type DistributorTaxUpdate struct {
	CustId              string  `json:"cust_id"`
	DistributorTaxId    *int64  `json:"distributor_tax_id"`
	DistributorId       int64   `json:"distributor_id"`
	TaxIdentifierNoType *string `json:"tax_identifier_no_type"`
	TaxIdentifierNo     *string `json:"tax_identifier_no"`
	Nitku               *string `json:"nitku"`
	TaxName             *string `json:"tax_name"`
	TaxAddress          *string `json:"tax_address"`
}
