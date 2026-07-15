package entity

type ValidateOrderQueryFilter struct {
	CustID            string  `query:"cust_id" json:"cust_id,omitempty"`
	ParentCustID      string  `query:"parent_cust_id" json:"parent_cust_id,omitempty"`
	Date              string  `query:"date" json:"date" validate:"omitempty,len=10,yyyyMmDdDate"`
	Page              int     `query:"page" json:"page"`
	Limit             int     `query:"limit" json:"limit" validate:""`
	Query             string  `query:"q" json:"q"`
	Sort              string  `query:"sort" json:"sort" validate:"oneof='pro_id:asc' 'pro_id:desc' 'pro_code:asc' 'pro_code:desc' 'pro_name:asc' 'pro_name:desc'"`
	ProID             []int64 `query:"pro_id" json:"pro_id"`
	WhID              []int   `query:"wh_id" json:"wh_id" validate:"required,min=1"`
	OutletID          int64   `json:"outlet_id" validate:"required"`
	ActiveProductOnly string  `query:"active_product_only" json:"active_product_only" validate:"omitempty,oneof='true' 'false'"`
}

type ValidateOrderBody struct {
	CustID            string             `json:"cust_id,omitempty"`
	ParentCustID      string             `json:"parent_cust_id,omitempty"`
	OutletID          int64              `json:"outlet_id" validate:"required"`
	WhID              int                `json:"wh_id" validate:"required,min=1"`
	Total             float64            `json:"total"`
	ProStok           []ProductsValidate `json:"product"`
	Date              string             `json:"date"`
	ProID             []int64            `json:"pro_id"`
	ActiveProductOnly string             `json:"active_product_only" validate:"omitempty,oneof='true' 'false'"`
	Page              int                `json:"page"`
	Limit             int                `json:"limit" validate:""`
	Query             string             `json:"q"`
	Sort              string             `json:"sort"`
}

type ValidateOrderDetailBody struct {
	CustID       string `json:"cust_id,omitempty"`
	ParentCustID string `json:"parent_cust_id,omitempty"`
	Limit        int    `json:"limit" validate:""`
	Sort         string `json:"sort"`
	Date         string `json:"date"`
	OutletID     int64  `json:"outlet_id" validate:"required"`
}

type ProductsValidate struct {
	CustId     string `json:"cust_id" db:"cust_id"`
	ProductId  int64  `json:"pro_id" db:"pro_id"`
	Qty1       int64  `json:"qty1"`
	Qty2       int64  `json:"qty2"`
	Qty3       int64  `json:"qty3"`
	QtyChange1 int64  `json:"qty_change1"`
	QtyChange2 int64  `json:"qty_change2"`
	QtyChange3 int64  `json:"qty_change3"`
}

type StockReport struct {
	ProID       int64   `json:"pro_id"`
	ProCode     string  `json:"pro_code"`
	ProName     string  `json:"pro_name"`
	UnitId1     string  `json:"unit_id1"`
	UnitId2     string  `json:"unit_id2"`
	UnitId3     string  `json:"unit_id3"`
	PurchPrice1 float64 `json:"purch_price1"`
	PurchPrice2 float64 `json:"purch_price2"`
	PurchPrice3 float64 `json:"purch_price3"`
	SellPrice1  float64 `json:"sell_price1"`
	SellPrice2  float64 `json:"sell_price2"`
	SellPrice3  float64 `json:"sell_price3"`
	TotalQty    float64 `json:"total_qty"`
	Qty1        float64 `json:"qty1"`
	Qty2        float64 `json:"qty2"`
	Qty3        float64 `json:"qty3"`
	ConvUnit2   int     `json:"conv_unit2"`
	ConvUnit3   int     `json:"conv_unit3"`
	IsActive    bool    `json:"is_active"`
	Vat         float64 `json:"vat"`
	VatLgPurch  float64 `json:"vat_lg_purch"`
	VatLgSell   float64 `json:"vat_lg_sell"`
}

type ArListResponse struct {
	InvoiceNo         *string `json:"invoice_no"`
	InvoiceDate       *string `json:"invoice_date"`
	DueDate           *string `json:"due_date"`
	InvoiceAmount     float64 `json:"invoice_amount"`
	PaidAmount        float64 `json:"paid_amount"`
	RemainingAmount   float64 `json:"remaining_amount"`
	SalesmanId        int64   `json:"salesman_id"`
	SalesmanCode      *string `json:"salesman_code"`
	SalesmanName      *string `json:"salesman_name"`
	OutletID          int64   `json:"outlet_id"`
	OutletCode        *string `json:"outlet_code"`
	OutletName        *string `json:"outlet_name"`
	InvoiceStatus     int64   `json:"invoice_status"`
	InvoiceStatusName string  `json:"invoice_status_name"`
	DueDateStatus     *int64  `json:"due_date_status"`
	DueDateStatusName *string `json:"due_date_status_name"`
	Aging             *int64  `json:"aging"`
}

type ValidateResponse struct {
	Validate1Success  bool    `json:"validate1_success"`
	Validate1         string  `json:"validate1_message"`
	Validate2Success  bool    `json:"validate2_success"`
	Validate2         string  `json:"validate2_message"`
	Validate2value    float64 `json:"validate2_value"`
	Validate3Success  bool    `json:"validate3_success"`
	Validate3         string  `json:"validate3_message"`
	Validate3Value    int     `json:"validate3_value"`
	Validate4Success  bool    `json:"validate4_success"`
	Validate4         string  `json:"validate4_message"`
	Validate4Value    int     `json:"validate4_value"`
	IsSuccessValidate bool    `json:"validate_summary_success"`
}

type ValidateDetailResponse struct {
	DuedateInvoive     []ArListResponse `json:"due_date_invoice"`
	OutstandingInvoice []ArListResponse `json:"outstanding_invoice"`
}
