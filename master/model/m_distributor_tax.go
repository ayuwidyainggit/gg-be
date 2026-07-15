package model

type DistributorTax struct {
	CustId              string `db:"cust_id" json:"cust_id"`
	DistributorId       int    `db:"distributor_id" json:"distributor_id"`
	DistributorTaxId    int64  `db:"distributor_tax_id" json:"distributor_tax_id"`
	TaxIdentifierNoType string `db:"tax_identifier_no_type" json:"tax_identifier_no_type"`
	TaxIdentifierNo     string `db:"tax_identifier_no" json:"tax_identifier_no"`
	Nitku               string `db:"nitku" json:"nitku"`
	TaxName             string `db:"tax_name" json:"tax_name"`
	TaxAddress          string `db:"tax_address" json:"tax_address"`
}

type DistributorTaxUpdate struct {
	TaxIdentifierNoType *string `db:"tax_identifier_no_type" json:"tax_identifier_no_type"`
	TaxIdentifierNo     *string `db:"tax_identifier_no" json:"tax_identifier_no"`
	Nitku               *string `db:"nitku" json:"nitku"`
	TaxName             *string `db:"tax_name" json:"tax_name"`
	TaxAddress          *string `db:"tax_address" json:"tax_address"`
}
