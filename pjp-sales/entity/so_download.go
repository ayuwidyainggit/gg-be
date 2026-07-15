package entity

type SoDownloadQueryFilter struct {
	CustId       string
	ParentCustId string
	StartDate    int64 `query:"start_date" validate:"required,gte=1000000000"`
	EndDate      int64 `query:"end_date" validate:"required,gte=1000000000"`
	SalesmanId   []int64
	ExportBy     string
	ReportID     string
}

type SoDownloadPoRow struct {
	OrderNo                   string   `json:"order_no"`
	PoNo                      string   `json:"po_no"`
	SoNo                      string   `json:"so_no"`
	OrderDate                 string   `json:"order_date"`
	InvoiceDate               string   `json:"invoice_date"`
	InvoiceNo                 string   `json:"invoice_no"`
	OutletCode                string   `json:"outlet_code"`
	OutletName                string   `json:"outlet_name"`
	SalesmanCode              *string  `json:"salesman_code"`
	EmployeeName              string   `json:"employee_name"`
	SupplierCode              string   `json:"supplier_code"`
	SupplierName              string   `json:"supplier_name"`
	ProductCode               string   `json:"product_code"`
	ProductName               string   `json:"product_name"`
	LargestUnit               string   `json:"largest_unit"`
	MiddleUnit                string   `json:"middle_unit"`
	SmallestUnit              string   `json:"smallest_unit"`
	FinalLargestSellingPrice  *float64 `json:"final_largest_selling_price"`
	FinalMiddleSellingPrice   *float64 `json:"final_middle_selling_price"`
	FinalSmallestSellingPrice *float64 `json:"final_smallest_selling_price"`
	LargestSellingPrice       *float64 `json:"largest_selling_price"`
	MiddleSellingPrice        *float64 `json:"middle_selling_price"`
	SmallestSellingPrice      *float64 `json:"smallest_selling_price"`
	LargestQtyOrder           *float64 `json:"largest_qty_order"`
	MiddleQtyOrder            *float64 `json:"middle_qty_order"`
	SmallestQtyOrder          *float64 `json:"smallest_qty_order"`
	GrossSales                *float64 `json:"gross_sales"`
	Promotion                 *float64 `json:"promotion"`
	Discount                  *float64 `json:"discount"`
	NetSales                  *float64 `json:"net_sales"`
	Vat                       *float64 `json:"vat"`
	Gross                     *float64 `json:"gross"`
}

type SoDownloadSoRow struct {
	OrderNo                   string   `json:"order_no"`
	PoNo                      string   `json:"po_no"`
	SoNo                      string   `json:"so_no"`
	OrderDate                 string   `json:"order_date"`
	InvoiceDate               string   `json:"invoice_date"`
	InvoiceNo                 string   `json:"invoice_no"`
	OutletCode                string   `json:"outlet_code"`
	OutletName                string   `json:"outlet_name"`
	SalesmanCode              *string  `json:"salesman_code"`
	EmployeeName              string   `json:"employee_name"`
	SupplierCode              string   `json:"supplier_code"`
	SupplierName              string   `json:"supplier_name"`
	ProductCode               string   `json:"product_code"`
	ProductName               string   `json:"product_name"`
	LargestUnit               string   `json:"largest_unit"`
	MiddleUnit                string   `json:"middle_unit"`
	SmallestUnit              string   `json:"smallest_unit"`
	FinalLargestSellingPrice  *float64 `json:"final_largest_selling_price"`
	FinalMiddleSellingPrice   *float64 `json:"final_middle_selling_price"`
	FinalSmallestSellingPrice *float64 `json:"final_smallest_selling_price"`
	LargestSellingPrice       *float64 `json:"largest_selling_price"`
	MiddleSellingPrice        *float64 `json:"middle_selling_price"`
	SmallestSellingPrice      *float64 `json:"smallest_selling_price"`
	LargestQtyOrder           *float64 `json:"largest_qty_order"`
	MiddleQtyOrder            *float64 `json:"middle_qty_order"`
	SmallestQtyOrder          *float64 `json:"smallest_qty_order"`
	GrossSales                *float64 `json:"gross_sales"`
	Promotion                 *float64 `json:"promotion"`
	Discount                  *float64 `json:"discount"`
	NetSales                  *float64 `json:"net_sales"`
	Vat                       *float64 `json:"vat"`
	Gross                     *float64 `json:"gross"`
}

type SoDownloadFinalRow struct {
	OrderNo                   string   `json:"order_no"`
	PoNo                      string   `json:"po_no"`
	SoNo                      string   `json:"so_no"`
	OrderDate                 string   `json:"order_date"`
	InvoiceDate               string   `json:"invoice_date"`
	InvoiceNo                 string   `json:"invoice_no"`
	OutletCode                string   `json:"outlet_code"`
	OutletName                string   `json:"outlet_name"`
	SalesmanCode              *string  `json:"salesman_code"`
	EmployeeName              string   `json:"employee_name"`
	SupplierCode              string   `json:"supplier_code"`
	SupplierName              string   `json:"supplier_name"`
	ProductCode               string   `json:"product_code"`
	ProductName               string   `json:"product_name"`
	LargestUnit               string   `json:"largest_unit"`
	MiddleUnit                string   `json:"middle_unit"`
	SmallestUnit              string   `json:"smallest_unit"`
	FinalLargestSellingPrice  *float64 `json:"final_largest_selling_price"`
	FinalMiddleSellingPrice   *float64 `json:"final_middle_selling_price"`
	FinalSmallestSellingPrice *float64 `json:"final_smallest_selling_price"`
	LargestSellingPrice       *float64 `json:"largest_selling_price"`
	MiddleSellingPrice        *float64 `json:"middle_selling_price"`
	SmallestSellingPrice      *float64 `json:"smallest_selling_price"`
	LargestQtyOrder           *float64 `json:"largest_qty_order"`
	MiddleQtyOrder            *float64 `json:"middle_qty_order"`
	SmallestQtyOrder          *float64 `json:"smallest_qty_order"`
	GrossSales                *float64 `json:"gross_sales"`
	Promotion                 *float64 `json:"promotion"`
	Discount                  *float64 `json:"discount"`
	NetSales                  *float64 `json:"net_sales"`
	Vat                       *float64 `json:"vat"`
	Gross                     *float64 `json:"gross"`
}

type SoDownloadQtySummaryRow struct {
	OrderNo          string   `json:"order_no"`
	PoNo             string   `json:"po_no"`
	SoNo             string   `json:"so_no"`
	OrderDate        string   `json:"order_date"`
	InvoiceDate      string   `json:"invoice_date"`
	InvoiceNo        string   `json:"invoice_no"`
	OutletCode       string   `json:"outlet_code"`
	OutletName       string   `json:"outlet_name"`
	SalesmanCode     *string  `json:"salesman_code"`
	EmployeeName     string   `json:"employee_name"`
	SupplierCode     string   `json:"supplier_code"`
	SupplierName     string   `json:"supplier_name"`
	ProductCode      string   `json:"product_code"`
	ProductName      string   `json:"product_name"`
	LargestUnit      string   `json:"largest_unit"`
	MiddleUnit       string   `json:"middle_unit"`
	SmallestUnit     string   `json:"smallest_unit"`
	LargestQtyPo     *float64 `json:"largest_qty_po"`
	MiddleQtyPo      *float64 `json:"middle_qty_po"`
	SmallestQtyPo    *float64 `json:"smallest_qty_po"`
	LargestQtySo     *float64 `json:"largest_qty_so"`
	MiddleQtySo      *float64 `json:"middle_qty_so"`
	SmallestQtySo    *float64 `json:"smallest_qty_so"`
	LargestQtyFinal  *float64 `json:"largest_qty_final"`
	MiddleQtyFinal   *float64 `json:"middle_qty_final"`
	SmallestQtyFinal *float64 `json:"smallest_qty_final"`
}

type SoDownloadResponse struct {
	DataPo     []SoDownloadPoRow         `json:"data_po"`
	DataSo     []SoDownloadSoRow         `json:"data_so"`
	DataFinal  []SoDownloadFinalRow      `json:"data_final"`
	QtySummary []SoDownloadQtySummaryRow `json:"qty_summary"`
}
