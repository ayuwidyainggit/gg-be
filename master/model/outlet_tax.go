package model

type MOutletTax struct {
	CustID            string `db:"cust_id" json:"cust_id"`
	OutletID          int64  `db:"outlet_id" json:"outlet_id"`
	IsEmbBail         bool   `json:"is_emb_bail" db:"is_emb_bail"`
	TaxName           string `json:"tax_name" db:"tax_name"`
	TaxAddr1          string `json:"tax_addr1" db:"tax_addr1"`
	TaxAddr2          string `json:"tax_addr2" db:"tax_addr2"`
	TaxNo             string `json:"tax_no" db:"tax_no"`
	TaxCity           string `json:"tax_city" db:"tax_city"`
	OutletTaxId       *int64 `db:"outlet_tax_id" json:"outlet_tax_id"`
	TaxInvoiceId      int64  `db:"tax_invoice_id" json:"tax_invoice_id"`
	TaxType           string `db:"tax_type" json:"tax_type"`
	Nitku             string `db:"nitku" json:"nitku"`
	AdddressTax       string `db:"address_tax" json:"address_tax"`
	TaxIdentifierType string `json:"tax_identifier_type" db:"tax_identifier_type"`
	TaxIdentifierNo   string `json:"tax_identifier_no" db:"tax_identifier_no"`
}

type MOutletTaxUpdate struct {
	IsEmbBail         *bool   `json:"is_emb_bail" db:"is_emb_bail" sql:"is_emb_bail"`
	TaxName           *string `json:"tax_name" db:"tax_name" sql:"tax_name"`
	TaxAddr1          *string `json:"tax_addr1" db:"tax_addr1" sql:"tax_addr1"`
	TaxAddr2          *string `json:"tax_addr2" db:"tax_addr2" sql:"tax_addr2"`
	TaxCity           *string `json:"tax_city" db:"tax_city" sql:"tax_city"`
	TaxNo             *string `json:"tax_no" db:"tax_no" sql:"tax_no"`
	TaxInvoiceId      *int64  `db:"tax_invoice_id" json:"tax_invoice_id" sql:"tax_invoice_id"`
	TaxType           *string `db:"tax_type" json:"tax_type" sql:"tax_type"`
	Nitku             *string `db:"nitku" json:"nitku" sql:"nitku"`
	AdddressTax       *string `db:"address_tax" json:"address_tax" sql:"address_tax"`
	TaxIdentifierType *string `json:"tax_identifier_type" db:"tax_identifier_type" sql:"tax_identifier_type"`
	TaxIdentifierNo   *string `json:"tax_identifier_no" db:"tax_identifier_no" sql:"tax_identifier_no"`
}
type MOutletTaxRead struct {
	CustID            string  `db:"cust_id" json:"cust_id"`
	OutletID          int64   `db:"outlet_id" json:"outlet_id"`
	IsEmbBail         *bool   `json:"is_emb_bail" db:"is_emb_bail"`
	TaxName           *string `json:"tax_name" db:"tax_name"`
	TaxAddr1          *string `json:"tax_addr1" db:"tax_addr1"`
	TaxAddr2          *string `json:"tax_addr2" db:"tax_addr2"`
	TaxNo             *string `json:"tax_no" db:"tax_no"`
	TaxCity           *string `json:"tax_city" db:"tax_city"`
	OutletTaxId       *int64  `db:"outlet_tax_id" json:"outlet_tax_id"`
	TaxInvoiceId      *int64  `db:"tax_invoice_id" json:"tax_invoice_id"`
	TaxType           string  `db:"tax_type" json:"tax_type"`
	Nitku             string  `db:"nitku" json:"nitku"`
	AdddressTax       string  `db:"address_tax" json:"address_tax"`
	TaxIdentifierType *string `json:"tax_identifier_type" db:"tax_identifier_type"`
	TaxIdentifierNo   *string `json:"tax_identifier_no" db:"tax_identifier_no"`
}
