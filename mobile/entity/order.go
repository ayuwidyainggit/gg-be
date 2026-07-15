package entity

var dataStatusName = map[int64]string{
	1: "Need Review",
	2: "Processed",
	3: "On Delivery",
	4: "Received",
	5: "Partial Received,",
	6: "Invoicing,",
	7: "Completed,",
	8: "Cancelled,",
}

var payTypeName = map[int64]string{
	1: "Cash",
	2: "Check",
	3: "Transfer",
	4: "Credit",
}

type OrderQueryFilter struct {
	SalesmanId   []int `query:"salesman_id"`
	OutletID     []int `query:"outlet_id"`
	Status       []int `query:"status"`
	CustId       string
	ParentCustId string
	From         *int64 `query:"from" validate:"required_with=To,omitempty,gte=1000000000"`
	To           *int64 `query:"to" validate:"required_with=From,omitempty,lte=9999999999,gtefield=From"`
	RoFrom       *int64 `query:"ro_date_from" validate:"required_with=RoTo,omitempty,gte=1000000000"`
	RoTo         *int64 `query:"ro_date_to" validate:"required_with=RoFrom,omitempty,lte=9999999999,gtefield=RoFrom"`
	InvoiceFrom  *int64 `query:"inv_date_from" validate:"required_with=InvoiceTo,omitempty,gte=1000000000"`
	InvoiceTo    *int64 `query:"inv_date_to" validate:"required_with=InvoiceFrom,omitempty,lte=9999999999,gtefield=InvoiceFrom"`
	Page         int    `query:"page"`
	Limit        int    `query:"limit" validate:"required"`
	Query        string `query:"q"`
	Mode         string `query:"mode"`
	Sort         string `query:"sort"`
	IsActive     *int   `query:"is_active"`
}

type NoOrderQueryFilter struct {
	SalesmanId   []int `query:"salesman_id"`
	OutletID     []int `query:"outlet_id"`
	Status       []int `query:"status"`
	CustId       string
	ParentCustId string
	From         *int64 `query:"from" validate:"required_with=To,omitempty,gte=1000000000"`
	To           *int64 `query:"to" validate:"required_with=From,omitempty,lte=9999999999,gtefield=From"`
	RoFrom       *int64 `query:"ro_date_from" validate:"required_with=RoTo,omitempty,gte=1000000000"`
	RoTo         *int64 `query:"ro_date_to" validate:"required_with=RoFrom,omitempty,lte=9999999999,gtefield=RoFrom"`
	InvoiceFrom  *int64 `query:"inv_date_from" validate:"required_with=InvoiceTo,omitempty,gte=1000000000"`
	InvoiceTo    *int64 `query:"inv_date_to" validate:"required_with=InvoiceFrom,omitempty,lte=9999999999,gtefield=InvoiceFrom"`
	Page         int    `query:"page"`
	Limit        int    `query:"limit" validate:"required"`
	Query        string `query:"q"`
	Mode         string `query:"mode"`
	Sort         string `query:"sort"`
	IsActive     *int   `query:"is_active"`
}

type CreateOrderBody struct {
	CustId        string            `json:"cust_id"`
	RoDate        *string           `json:"ro_date"`
	ValDate       *string           `json:"val_date"`
	DueDate       *string           `json:"due_date"`
	SalesmanId    int64             `json:"salesman_id" validate:"required"`
	WhId          int64             `json:"wh_id" validate:"required"`
	OutletID      int64             `json:"outlet_id" validate:"required"`
	DeliveryDate  *string           `json:"delivery_date"`
	OrderNo       *string           `json:"order_no"`
	PoNo          *string           `json:"po_no"`
	VehicleNo     *string           `json:"vehicle_no"`
	PayType       *int64            `json:"pay_type"`
	ReffNo        *string           `json:"reff_no"`
	MobileID      *int64            `json:"mobile_id"`
	SubTotal      *float64          `json:"sub_total"`
	Disc          *float64          `json:"disc"`
	DiscValue     *float64          `json:"disc_value"`
	PromoValue    *float64          `json:"promo_value"`
	CashDiscValue *float64          `json:"cash_disc_value"`
	TotDisc1      *float64          `json:"tot_disc1"`
	TotDisc2      *float64          `json:"tot_disc2"`
	Vat           *float64          `json:"vat"`
	VatValue      *float64          `json:"vat_value"`
	Amount        *float64          `json:"amount"`
	Total         *float64          `json:"total"`
	DataStatus    *int64            `json:"data_status"`
	CreatedBy     *int64            `json:"created_by"`
	DataSource    *int64            `json:"data_source"`
	OprType       *string           `json:"opr_type"`
	Details       OrderDetWithGroup `json:"details"`
	TrCode        *string           `json:"tr_code"`
	IsClosed      bool              `json:"is_closed"`
	Notes         *string           `json:"notes"`
	InvoiceNo     *string           `json:"invoice_no"`
	InvoiceDate   *string           `json:"invoice_date"`
}

type CreateNoOrderBody struct {
	CustId        string  `json:"cust_id"`
	SalesmanId    int64   `json:"salesman_id" validate:"required"`
	NoOrderDate   *string `json:"no_order_date"`
	OutletId      int64   `json:"outlet_id" validate:"required"`
	TakingOrderId int64   `json:"taking_order_id"`
	Reason        *string `json:"reason"`
	CreatedBy     int64   `json:"created_by"`
	CreatedAt     *string `json:"created_at"`
}

type OrderResponse struct {
	RoNo           string                `json:"ro_no"`
	OrderNo        *string               `json:"order_no"`
	RoDate         *string               `json:"ro_date"`
	ValDate        *string               `json:"val_date"`
	SalesmanId     *int64                `json:"salesman_id"`
	SalesName      string                `json:"sales_name"`
	WhId           *int64                `json:"wh_id"`
	WhCode         string                `json:"wh_code"`
	WhName         string                `json:"wh_name"`
	OutletID       *int64                `json:"outlet_id"`
	OutletCode     string                `json:"outlet_code"`
	OutletName     string                `json:"outlet_name"`
	OutletAddress1 string                `json:"outlet_address1"`
	OutletAddress2 string                `json:"outlet_address2"`
	DeliveryDate   *string               `json:"delivery_date"`
	PoNo           *string               `json:"po_no"`
	VehicleNo      *string               `json:"vehicle_no"`
	PayType        *int64                `json:"pay_type"`
	PayTypeName    string                `json:"pay_type_name"`
	ReffNo         *string               `json:"reff_no"`
	MobileID       *int64                `json:"mobile_id"`
	SubTotal       *float64              `json:"sub_total"`
	Disc           *float64              `json:"disc"`
	DiscValue      *float64              `json:"disc_value"`
	PromoValue     *float64              `json:"promo_value"`
	CashDiscValue  *float64              `json:"cash_disc_value"`
	TotDisc1       *float64              `json:"tot_disc1"`
	TotDisc2       *float64              `json:"tot_disc2"`
	Vat            *float64              `json:"vat"`
	VatValue       *float64              `json:"vat_value"`
	Total          *float64              `json:"total"`
	DataStatus     *int64                `json:"data_status"`
	DataStatusName string                `json:"data_status_name"`
	DataSource     *int64                `json:"data_source"`
	UpdatedAt      string                `json:"updated_at"`
	UpdatedByName  string                `json:"updated_by_name"`
	DueDate        *string               `json:"due_date"`
	Details        OrderDetReadWithGroup `json:"details"`
	DetailsFinal   OrderDetReadWithGroup `json:"details_final"`
	TrCode         *string               `json:"tr_code"`
	IsClosed       bool                  `json:"is_closed"`
	Notes          *string               `json:"notes"`
	InvoiceNo      *string               `json:"invoice_no"`
	InvoiceDate    *string               `json:"invoice_date"`
	IsPrinted      *bool                 `json:"is_printed"`
	PrintedBy      *int64                `json:"printed_by"`
	PrintedByName  *string               `json:"printed_by_name"`
	PrintedAt      *string               `json:"printed_at"`
}

func (ro OrderResponse) GeneratePayTypeName() string {
	if ro.PayType != nil {
		return payTypeName[*ro.PayType]
	}
	return ""
}

func (ro OrderResponse) GenerateDataStatusName() string {
	if ro.DataStatus != nil {
		return dataStatusName[*ro.DataStatus]
	}
	return ""
}

type DetailOrderParams struct {
	RoNo string `params:"ro_no" validate:"required"`
}
type DeleteOrderParams struct {
	RoNo string `params:"ro_no" validate:"required"`
}

type UpdateOrderParams struct {
	RoNo string `params:"ro_no" validate:"required"`
}

type OrderListResponse struct {
	RoNo           string   `json:"ro_no"`
	RoDate         *string  `json:"ro_date"`
	ValDate        *string  `json:"val_date"`
	SalesmanId     *int64   `json:"salesman_id"`
	SalesName      string   `json:"sales_name"`
	WhId           *int64   `json:"wh_id"`
	WhCode         string   `json:"wh_code"`
	WhName         string   `json:"wh_name"`
	OutletID       *int64   `json:"outlet_id"`
	OutletCode     string   `json:"outlet_code"`
	OutletName     string   `json:"outlet_name"`
	DeliveryDate   *string  `json:"delivery_date"`
	OrderNo        *string  `json:"order_no"`
	PoNo           *string  `json:"po_no"`
	VehicleNo      *string  `json:"vehicle_no"`
	PayType        *int64   `json:"pay_type"`
	PayTypeName    string   `json:"pay_type_name"`
	ReffNo         *string  `json:"reff_no"`
	MobileID       *int64   `json:"mobile_id"`
	SubTotal       *float64 `json:"sub_total"`
	Disc           *float64 `json:"disc"`
	DiscValue      *float64 `json:"disc_value"`
	PromoValue     *float64 `json:"promo_value"`
	CashDiscValue  *float64 `json:"cash_disc_value"`
	TotDisc1       *float64 `json:"tot_disc1"`
	TotDisc2       *float64 `json:"tot_disc2"`
	Vat            *float64 `json:"vat"`
	VatValue       *float64 `json:"vat_value"`
	Total          *float64 `json:"total"`
	DataStatus     *int64   `json:"data_status"`
	DataStatusName string   `json:"data_status_name"`
	UpdatedAt      string   `json:"updated_at"`
	UpdatedByName  string   `json:"updated_by_name"`
	DueDate        *string  `json:"due_date"`
	TrCode         *string  `json:"tr_code"`
	IsClosed       bool     `json:"is_closed"`
	Notes          *string  `json:"notes"`
	InvoiceNo      *string  `json:"invoice_no"`
	InvoiceDate    *string  `json:"invoice_date"`
}

type NoOrderListResponse struct {
	NoOrderId *int32 `json:"no_order_id"`

	SalesmanId *int64 `json:"salesman_id"`
	SalesName  string `json:"sales_name"`

	NoOrderDate *string `json:"no_order_date"`

	OutletID         *int64  `json:"outlet_id"`
	OutletCode       string  `json:"outlet_code"`
	OutletName       string  `json:"outlet_name"`
	TakingOrderId    *int64  `json:"taking_order_id"`
	TakingOrderName  *string `json:"taking_order_name"`
	TakingOrderImage *string `json:"image_url"`

	Reason        *string `json:"reason"`
	UpdatedAt     string  `json:"updated_at"`
	UpdatedByName string  `json:"updated_by_name"`
}

func (ro OrderListResponse) GenerateDataStatusName() string {
	if ro.DataStatus != nil {
		return dataStatusName[*ro.DataStatus]
	}
	return ""
}

func (ro OrderListResponse) GeneratePayTypeName() string {
	if ro.PayType != nil {
		return payTypeName[*ro.PayType]
	}
	return ""
}

type UpdateOrderBody struct {
	CustId        string                  `json:"cust_id"`
	RoNo          string                  `json:"ro_no"`
	RoDate        *string                 `json:"ro_date"`
	ValDate       *string                 `json:"val_date"`
	DueDate       *string                 `json:"due_date"`
	SalesmanId    *int64                  `json:"salesman_id"`
	WhId          *int64                  `json:"wh_id"`
	OutletID      *int64                  `json:"outlet_id"`
	DeliveryDate  *string                 `json:"delivery_date"`
	OrderNo       *string                 `json:"order_no"`
	PoNo          *string                 `json:"po_no"`
	VehicleNo     *string                 `json:"vehicle_no"`
	PayType       *int64                  `json:"pay_type"`
	ReffNo        *string                 `json:"reff_no"`
	MobileID      *int64                  `json:"mobile_id"`
	SubTotal      *float64                `json:"sub_total"`
	Disc          *float64                `json:"disc"`
	DiscValue     *float64                `json:"disc_value"`
	PromoValue    *float64                `json:"promo_value"`
	CashDiscValue *float64                `json:"cash_disc_value"`
	TotDisc1      *float64                `json:"tot_disc1"`
	TotDisc2      *float64                `json:"tot_disc2"`
	Vat           *float64                `json:"vat"`
	VatValue      *float64                `json:"vat_value"`
	Total         *float64                `json:"total"`
	DataStatus    *int64                  `json:"data_status"`
	CreatedBy     *int64                  `json:"created_by"`
	CreatedAt     *string                 `json:"created_at"`
	UpdatedBy     int64                   `json:"updated_by"`
	Details       UpdateOrderDetWithGroup `json:"details"`
	TrCode        *string                 `json:"tr_code"`
	IsClosed      bool                    `json:"is_closed"`
	Notes         *string                 `json:"notes"`
	InvoiceNo     *string                 `json:"invoice_no"`
	InvoiceDate   *string                 `json:"invoice_date"`
}

type UpdateOrderDetailFinal struct {
	CustId    string                  `json:"cust_id"`
	RoNo      string                  `json:"ro_no"`
	UpdatedBy int64                   `json:"updated_by"`
	Details   UpdateOrderDetWithGroup `json:"details_final"`
}

type CreateConversionBody struct {
	CustId    string `json:"cust_id" validate:"required,max=10"`
	ProductId int64  `json:"pro_id"`
	Qty1      int64  `json:"qty1" validate:"numeric"`
	Qty2      int64  `json:"qty2" validate:"numeric"`
	Qty3      int64  `json:"qty3" validate:"numeric"`
}

type OrderConversionResponse struct {
	Qty1     int64 `json:"qty1"`
	Qty2     int64 `json:"qty2"`
	Qty3     int64 `json:"qty3"`
	TotalQty int64 `json:"total_qty"`
}

type OutletBySalesmanResponse struct {
	OutletID   int     `json:"outlet_id" `
	OutletCode *string `json:"outlet_code" `
	OutletName *string `json:"outlet_name" `
}

type UpdateDataStatusBody struct {
	CustId     string `json:"cust_id"`
	RoNo       string `json:"ro_no"`
	DataStatus *int64 `json:"data_status"`
	UpdatedBy  int64  `json:"updated_by"`
}

type BulkUpdateStatusOrder struct {
	Orders []UpdateDataStatusBody `json:"orders" validate:"min=1"`
}

type OrderDiscountQuery struct {
	CustId       string
	ParentCustId string
	OrderDate    *int64  `query:"order_date" validate:"required_with=To,omitempty,gte=1000000000"`
	OutletID     int     `query:"outlet_id"`
	ProID        int     `query:"pro_id"`
	GrossValue   float64 `query:"subtotal"`
	IsActive     *int    `query:"is_active"`
}

type ResponseSummaryTotal struct {
	SalesmanId   int     `json:"salesman_id"`
	SummaryToday float64 `json:"summary_today"`
	SummaryMonth float64 `json:"summary_month"`
	CreatedAt    string  `json:"order_date"`
}

type SummaryTotalFilter struct {
	CustId       string
	ParentCustId string
	SalesmanId   int    `query:"salesman_id"`
	Date         *int64 `query:"date" validate:"required_with=To,omitempty,gte=1000000000"`
}

type StatisticReportFilter struct {
	EmpID         int    `query:"emp_id" validate:"required"`
	CustID        string `query:"cust_id" validate:"required"`
	Type          string `query:"type" validate:"required"`
	IsDistributor bool
}

type StatisticReportResponse struct {
	Sales              float64             `json:"sales"`
	TotalVisitList     int                 `json:"total_visit_list"`
	TotalVisit         int                 `json:"total_visit"`
	TotalNotVisit      int                 `json:"total_not_visit"`
	VisitPlanned       int                 `json:"visit_planned"`
	VisitNotPlanned    int                 `json:"visit_not_planned"`
	NotVisitPlanned    int                 `json:"not_visit_planned"`
	NotVisitNotPlanned int                 `json:"not_visit_not_planned"`
	TotalBuy           int                 `json:"total_buy"`
	TotalNotBuy        int                 `json:"total_not_buy"`
	NotBuyReasonData   []NotBuyReasonItem  `json:"not_buy_reason_data"`
	SkipReasonData     []SkipReasonItem    `json:"skip_reason_data"`
}

type NotBuyReasonItem struct {
	Reason string `json:"reason"`
	Total  int    `json:"total"`
}

type SkipReasonItem struct {
	SkipReason string `json:"skip_reason"`
	Total      int    `json:"total"`
}

type VisitOverview struct {
	TotalBuy           int `json:"total_buy"`
	TotalNotBuy        int `json:"total_not_buy"`
	TotalVisitList     int `json:"total_visit_list"`
	TotalVisit         int `json:"total_visit"`
	TotalNotVisit      int `json:"total_not_visit"`
	VisitPlanned       int `json:"visit_planned"`
	VisitNotPlanned    int `json:"visit_not_planned"`
	NotVisitPlanned    int `json:"not_visit_planned"`
	NotVisitNotPlanned int `json:"not_visit_not_planned"`
}
