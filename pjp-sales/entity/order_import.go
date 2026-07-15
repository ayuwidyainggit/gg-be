package entity

import "strings"

type OrderImportError struct {
	Row     int    `json:"row"`
	Field   string `json:"field"`
	Message string `json:"message"`
}

type OrderImportRow struct {
	DocumentNo    string `json:"document_no"`
	DocumentDate  string `json:"document_date"`
	OutletCode    string `json:"outlet_code"`
	OutletName    string `json:"outlet_name"`
	SalesmanCode  string `json:"salesman_code"`
	SalesmanName  string `json:"salesman_name"`
	ProCode       string `json:"pro_code"`
	ProName       string `json:"pro_name"`
	Price         string `json:"price"`
	Unit          string `json:"unit"`
	Qty           string `json:"qty"`
	GrossSales    string `json:"gross_sales"`
	Promo         string `json:"promo"`
	Discount      string `json:"discount"`
	PPN           string `json:"ppn"`
	NetSalesIncPPN string `json:"net_sales_inc_ppn"`
}

type OrderImportResult struct {
	StartDate       string   `json:"start_date"`
	EndDate         string   `json:"end_date"`
	NumberOfInvoice int      `json:"number_of_invoice"`
	NumberOfOutlet  int      `json:"number_of_outlet"`
	Amount          float64  `json:"amount"`
	CreatedRoNos    []string `json:"created_ro_nos,omitempty"`
}

type OrderImportSummary struct {
	StartDate       string   `json:"start_date"`
	EndDate         string   `json:"end_date"`
	NumberOfInvoice int      `json:"number_of_invoice"`
	NumberOfOutlet  int      `json:"number_of_outlet"`
	Amount          float64  `json:"amount"`
	FailedReasons   []string `json:"failed_reasons"`
}

type OrderImportFromURLRequest struct {
	URL      string `json:"url"`
	Validate string `json:"validate"`
}

type ImportFailedError struct {
	FailedReasons []string
}

func (e *ImportFailedError) Error() string {
	if len(e.FailedReasons) == 0 {
		return "import failed"
	}
	return strings.Join(e.FailedReasons, "; ")
}
