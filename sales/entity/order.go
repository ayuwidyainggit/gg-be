package entity

const (
	NEED_REVIEW      = 1
	PROCESSED        = 2
	ON_DELIVERY      = 3
	RECEIVED         = 4
	PARTIAL_RECEIVED = 5
	INVOICING        = 6
	COMPLETED        = 7
	CANCELLED        = 9

	PAY_TYPE_CASH_ON_DELIVERY     = 1
	PAY_TYPE_CASH_BEFORE_DELIVERY = 2
	PAY_TYPE_CREDIT               = 3
)

var dataStatusName = map[int64]string{
	NEED_REVIEW:      "Need Review",
	PROCESSED:        "Processed",
	ON_DELIVERY:      "On Delivery",
	RECEIVED:         "Received",
	PARTIAL_RECEIVED: "Partial Received",
	INVOICING:        "Invoicing",
	COMPLETED:        "Completed",
	CANCELLED:        "Cancelled",
}

var payTypeName = map[int64]string{
	PAY_TYPE_CASH_ON_DELIVERY:     "Cash On Delivery",
	PAY_TYPE_CASH_BEFORE_DELIVERY: "Cash Before Delivery",
	PAY_TYPE_CREDIT:               "Credit",
	//4: "Credit",
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
	IsInvoice    *bool  `query:"is_invoice"`
}

type CreateOrderBody struct {
	RoNo              string                  `json:"ro_no"`
	CustId            string                  `json:"cust_id"`
	ParentCustId      string                  `json:"parent_cust_id"`
	RoDate            *string                 `json:"ro_date"`
	ValDate           *string                 `json:"val_date"`
	DueDate           *string                 `json:"due_date"`
	SalesmanId        int64                   `json:"salesman_id" validate:"required"`
	WhId              *int64                  `json:"wh_id"`
	OutletID          int64                   `json:"outlet_id" validate:"required"`
	OutletAddress1    *string                 `json:"outlet_address1"`
	DeliveryDate      *string                 `json:"delivery_date"`
	OrderNo           *string                 `json:"order_no"`
	PoNo              *string                 `json:"po_no"`
	VehicleNo         *string                 `json:"vehicle_no"`
	PayType           *int64                  `json:"pay_type"`
	ReffNo            *string                 `json:"reff_no"`
	MobileID          *int64                  `json:"mobile_id"`
	SubTotal          *float64                `json:"sub_total"`
	SubTotalFinal     *float64                `json:"sub_total_final"`
	Disc              *float64                `json:"disc"`
	DiscValue         *float64                `json:"disc_value"`
	DiscValueFinal    *float64                `json:"disc_value_final"`
	PromoValue        *float64                `json:"promo_value"`
	PromoValueFinal   *float64                `json:"promo_value_final"`
	PromoBgValue      *float64                `json:"promo_bg_value"`
	PromoBgValueFinal *float64                `json:"promo_bg_value_final"`
	CashDiscValue     *float64                `json:"cash_disc_value"`
	TotDisc1          *float64                `json:"tot_disc1"`
	TotDisc2          *float64                `json:"tot_disc2"`
	Vat               *float64                `json:"vat"`
	VatValue          *float64                `json:"vat_value"`
	VatValueFinal     *float64                `json:"vat_value_final"`
	Total             *float64                `json:"total"`
	TotalFinal        *float64                `json:"total_final"`
	DataStatus        *int                    `json:"data_status"`
	CreatedBy         *int64                  `json:"created_by"`
	DataSource        *int64                  `json:"data_source"`
	Details           OrderDetWithGroup       `json:"details"`
	TrCode            *string                 `json:"tr_code"`
	IsClosed          bool                    `json:"is_closed"`
	Notes             *string                 `json:"notes"`
	InvoiceNo         *string                 `json:"invoice_no"`
	InvoiceDate       *string                 `json:"invoice_date"`
	OrderType         *string                 `json:"order_type" validate:"omitempty,oneof=O C SO"`
	IsSalesMapping    *bool                   `json:"is_sales_mapping"`
	Rewards           []CreateOrderRewardBody `json:"rewards"`
}

type CreateOrderResponse struct {
	RoNo string `json:"ro_no"`
}

type OrderResponse struct {
	RoNo                                  string                `json:"ro_no"`
	OprType                               *string               `json:"opr_type"`
	Source                                *string               `json:"source"`
	IsProformaInv                         *bool                 `json:"is_proforma_inv"`
	OrderNo                               *string               `json:"order_no"`
	RoDate                                *string               `json:"ro_date"`
	ValDate                               *string               `json:"val_date"`
	SalesmanId                            *int64                `json:"salesman_id"`
	SalesmanCode                          *string               `json:"salesman_code"`
	SalesName                             string                `json:"sales_name"`
	WhId                                  *int64                `json:"wh_id"`
	WhCode                                string                `json:"wh_code"`
	WhName                                string                `json:"wh_name"`
	OutletID                              *int64                `json:"outlet_id"`
	OutletCode                            string                `json:"outlet_code"`
	OutletName                            string                `json:"outlet_name"`
	OutletAddress1                        string                `json:"outlet_address1"`
	OutletAddress2                        string                `json:"outlet_address2"`
	InvAddress1                           string                `json:"inv_addr1"`
	InvAddress2                           string                `json:"inv_addr2"`
	DeliveryDate                          *string               `json:"delivery_date"`
	PoNo                                  *string               `json:"po_no"`
	VehicleNo                             *string               `json:"vehicle_no"`
	PayType                               *int64                `json:"pay_type"`
	PayTypeName                           string                `json:"pay_type_name"`
	ReffNo                                *string               `json:"reff_no"`
	MobileID                              *int64                `json:"mobile_id"`
	SubTotal                              *float64              `json:"sub_total"`
	SubTotalFinal                         *float64              `json:"sub_total_final"`
	Disc                                  *float64              `json:"disc"`
	DiscValue                             *float64              `json:"disc_value"`
	DiscValueFinal                        *float64              `json:"disc_value_final"`
	PromoValue                            *float64              `json:"promo_value"`
	PromoValueFinal                       *float64              `json:"promo_value_final"`
	PromoBgValue                          *float64              `json:"promo_bg_value"`
	PromoBgValueFinal                     *float64              `json:"promo_bg_value_final"`
	CashDiscValue                         *float64              `json:"cash_disc_value"`
	TotDisc1                              *float64              `json:"tot_disc1"`
	TotDisc2                              *float64              `json:"tot_disc2"`
	Vat                                   *float64              `json:"vat"`
	VatValue                              *float64              `json:"vat_value"`
	VatValueFinal                         *float64              `json:"vat_value_final"`
	Total                                 *float64              `json:"total"`
	TotalFinal                            *float64              `json:"total_final"`
	DataStatus                            *int64                `json:"data_status"`
	DataStatusName                        string                `json:"data_status_name"`
	DataSource                            *int64                `json:"data_source"`
	UpdatedAt                             string                `json:"updated_at"`
	UpdatedByName                         string                `json:"updated_by_name"`
	DueDate                               *string               `json:"due_date"`
	Details                               OrderDetReadWithGroup `json:"details"`
	DetailsFinal                          OrderDetReadWithGroup `json:"details_final"`
	PurchaseDetails                       OrderDetReadWithGroup `json:"purchase_details"`
	Remarks                               []OrderRewardResponse `json:"remarks"`
	TrCode                                *string               `json:"tr_code"`
	IsClosed                              bool                  `json:"is_closed"`
	Notes                                 *string               `json:"notes"`
	InvoiceNo                             *string               `json:"invoice_no"`
	InvoiceDate                           *string               `json:"invoice_date"`
	IsPrinted                             *bool                 `json:"is_printed"`
	PrintedBy                             *int64                `json:"printed_by"`
	PrintedByName                         *string               `json:"printed_by_name"`
	PrintedAt                             *string               `json:"printed_at"`
	ValidateStok                          bool                  `json:"validate_stok"`
	ValidateStokMessage                   *string               `json:"validate_stok_message"`
	ValidateCreditLimit                   bool                  `json:"validate_credit_limit"`
	ValidateCreditLimitMessage            string                `json:"validate_credit_limit_message"`
	ValidateOverdue                       bool                  `json:"validate_overdue"`
	ValidateOverdueMessage                string                `json:"validate_overdue_message"`
	ValidateOutstanding                   bool                  `json:"validate_outstanding"`
	ValidateOutstandingMessage            string                `json:"validate_outstanding_message"`
	ValidateSummary                       bool                  `json:"validate_summary"`
	CreditLimitType                       *int                  `json:"credit_limit_type"`
	CreditLimitAction                     *int                  `json:"credit_limit_action"`
	CreditLimitActionName                 string                `json:"credit_limit_action_name"`
	SalesInvLimitType                     *int                  `json:"sales_inv_limit_type"`
	SalesInvLimitAction                   *int                  `json:"sales_inv_limit_action"`
	SalesInvLimitActionName               string                `json:"sales_inv_limit_action_name"`
	ObsType                               *int                  `json:"obs_type"`
	ObsLimitAction                        *int                  `json:"obs_limit_action"`
	ObsLimitActionName                    string                `json:"obs_limit_action_name"`
	OrderApprovalRequestID                *int64                `json:"order_approval_request_id"`
	OrderApprovalRequestEmpApprovalStatus *int                  `json:"order_approval_request_emp_approval_status"`
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

type DetailOrderQueryParams struct {
	NoCustID     bool   `query:"no_cust_id"`
	EmpID        *int64 `query:"emp_id"`
	CustIDOrigin string `query:"cust_id_origin"`
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
	OprType        *string  `json:"opr_type"`
	Source         *string  `json:"source"`
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
	OutletAddress1 *string  `json:"outlet_address1"`
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

// PayTypeName returns pay type name by pay type id
func PayTypeName(payType int64) string {
	if name, ok := payTypeName[payType]; ok {
		return name
	}
	return ""
}

type UpdateOrderBody struct {
	CustId            string                  `json:"cust_id"`
	ParentCustId      string                  `json:"parent_cust_id"`
	RoNo              string                  `json:"ro_no"`
	RoDate            *string                 `json:"ro_date"`
	ValDate           *string                 `json:"val_date"`
	DueDate           *string                 `json:"due_date"`
	SalesmanId        *int64                  `json:"salesman_id"`
	WhId              *int64                  `json:"wh_id"`
	OutletID          *int64                  `json:"outlet_id"`
	DeliveryDate      *string                 `json:"delivery_date"`
	OrderNo           *string                 `json:"order_no"`
	PoNo              *string                 `json:"po_no"`
	VehicleNo         *string                 `json:"vehicle_no"`
	PayType           *int64                  `json:"pay_type"`
	ReffNo            *string                 `json:"reff_no"`
	MobileID          *int64                  `json:"mobile_id"`
	SubTotal          *float64                `json:"sub_total"`
	SubTotalFinal     *float64                `json:"sub_total_final"`
	Disc              *float64                `json:"disc"`
	DiscValue         *float64                `json:"disc_value"`
	DiscValueFinal    *float64                `json:"disc_value_final"`
	PromoValue        *float64                `json:"promo_value"`
	PromoValueFinal   *float64                `json:"promo_value_final"`
	PromoBgValue      *float64                `json:"promo_bg_value"`
	PromoBgValueFinal *float64                `json:"promo_bg_value_final"`
	CashDiscValue     *float64                `json:"cash_disc_value"`
	TotDisc1          *float64                `json:"tot_disc1"`
	TotDisc2          *float64                `json:"tot_disc2"`
	Vat               *float64                `json:"vat"`
	VatValue          *float64                `json:"vat_value"`
	VatValueFinal     *float64                `json:"vat_value_final"`
	Total             *float64                `json:"total"`
	TotalFinal        *float64                `json:"total_final"`
	DataStatus        *int                    `json:"data_status"`
	CreatedBy         *int64                  `json:"created_by"`
	CreatedAt         *string                 `json:"created_at"`
	UpdatedBy         int64                   `json:"updated_by"`
	Details           UpdateOrderDetWithGroup `json:"details"`
	TrCode            *string                 `json:"tr_code"`
	IsClosed          bool                    `json:"is_closed"`
	Notes             *string                 `json:"notes"`
	InvoiceNo         *string                 `json:"invoice_no"`
	InvoiceDate       *string                 `json:"invoice_date"`
	Rewards           []CreateOrderRewardBody `json:"rewards"`
}

type UpdateOrderDetailFinal struct {
	CustId            string                       `json:"cust_id"`
	ParentCustId      string                       `json:"parent_cust_id"`
	RoNo              string                       `json:"ro_no"`
	OutletID          *int64                       `json:"outlet_id"`
	SubTotalFinal     *float64                     `json:"sub_total_final"`
	DiscValueFinal    *float64                     `json:"disc_value_final"`
	PromoValueFinal   *float64                     `json:"promo_value_final"`
	PromoBgValueFinal *float64                     `json:"promo_bg_value_final"`
	VatValueFinal     *float64                     `json:"vat_value_final"`
	TotalFinal        *float64                     `json:"total_final"`
	UpdatedBy         int64                        `json:"updated_by"`
	Details           UpdateFinalOrderDetWithGroup `json:"details_final"`
	Rewards           []CreateOrderRewardBody      `json:"rewards"`
}

type CreateConversionBody struct {
	CustId    string `json:"cust_id" validate:"required,max=10"`
	ProductId int64  `json:"pro_id"`
	Qty1      int64  `json:"qty1" validate:"omitempty,numeric"`
	Qty2      int64  `json:"qty2" validate:"omitempty,numeric"`
	Qty3      int64  `json:"qty3" validate:"omitempty,numeric"`
	Qty1Final *int64 `json:"qty1_final" validate:"omitempty,numeric"`
	Qty2Final *int64 `json:"qty2_final" validate:"omitempty,numeric"`
	Qty3Final *int64 `json:"qty3_final" validate:"omitempty,numeric"`
	OutletID  *int64 `json:"outlet_id"`
}

type OrderConversionResponse struct {
	Qty1      *int64   `json:"qty1,omitempty"`
	Qty2      *int64   `json:"qty2,omitempty"`
	Qty3      *int64   `json:"qty3,omitempty"`
	Qty1Final *int64   `json:"qty1_final,omitempty"`
	Qty2Final *int64   `json:"qty2_final,omitempty"`
	Qty3Final *int64   `json:"qty3_final,omitempty"`
	TotalQty  int64    `json:"total_qty"`
	Price     *float64 `json:"price,omitempty"`
	DiscValue *float64 `json:"disc_value,omitempty"`
	VatValue  *float64 `json:"vat_value,omitempty"`
	Total     *float64 `json:"total,omitempty"`
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

type ConsultDiscountOrderBody struct {
	CustId            string                           `json:"cust_id"`
	ParentCustId      string                           `json:"parent_cust_id"`
	SalesmanId        int64                            `json:"salesman_id"`
	WhId              *int64                           `json:"wh_id"`
	OutletID          int64                            `json:"outlet_id"`
	SubTotal          *float64                         `json:"sub_total"`
	SubTotalFinal     *float64                         `json:"sub_total_final"`
	Disc              *float64                         `json:"disc"`
	DiscValue         *float64                         `json:"disc_value"`
	DiscValueFinal    *float64                         `json:"disc_value_final"`
	PromoValue        *float64                         `json:"promo_value"`
	PromoValueFinal   *float64                         `json:"promo_value_final"`
	PromoBgValue      *float64                         `json:"promo_bg_value"`
	PromoBgValueFinal *float64                         `json:"promo_bg_value_final"`
	Vat               *float64                         `json:"vat"`
	VatValue          *float64                         `json:"vat_value"`
	VatValueFinal     *float64                         `json:"vat_value_final"`
	Total             *float64                         `json:"total"`
	TotalFinal        *float64                         `json:"total_final"`
	Details           ConsultDiscountOrderDetWithGroup `json:"details"`
	Rewards           []CreateOrderRewardBody          `json:"rewards"`
}

type OrderMinimumPriceFilter struct {
	SalesmanId   int64 `query:"salesman_id" validate:"required"`
	ProID        int64 `params:"pro_id" validate:"required"`
	CustId       string
	ParentCustId string
}

type OrderMinimumPriceResp struct {
	SalesmanId      int                           `json:"salesman_id"`
	SalesmanCode    string                        `json:"salesman_code"`
	SalesmanName    string                        `json:"salesman_name"`
	AllowInputPrice bool                          `json:"allow_input_price" `
	MinimumPrice    *OrderMinimumPriceSettingResp `json:"minimum_price"`
}

type OrderMinimumPriceSettingResp struct {
	ManageMinimumPriceId     int     `json:"manage_minimum_price_id"`
	BasePrice                int     `json:"base_price"`
	LimitAction              int     `json:"limit_action"`
	Threshold                float64 `json:"threshold"`
	StatusManageMinimumPrice int     `json:"status_manage_minimum_price"`
	ProId                    int     `json:"pro_id"`
	Price1                   float64 `json:"price1"`
	Price2                   float64 `json:"price2"`
	Price3                   float64 `json:"price3"`
	Price4                   float64 `json:"price4"`
	Price5                   float64 `json:"price5"`
	PriceMinimum1            float64 `json:"price1_minimum"`
	PriceMinimum2            float64 `json:"price2_minimum"`
	PriceMinimum3            float64 `json:"price3_minimum"`
	PriceMinimum4            float64 `json:"price4_minimum"`
	PriceMinimum5            float64 `json:"price5_minimum"`
	UnitId1                  string  `json:"unit_id1"`
	UnitId2                  string  `json:"unit_id2"`
	UnitId3                  string  `json:"unit_id3"`
	UnitId4                  string  `json:"unit_id4"`
	UnitId5                  string  `json:"unit_id5"`
	ConvUnit2                int     `json:"conv_unit2"`
	ConvUnit3                int     `json:"conv_unit3"`
	ConvUnit4                int     `json:"conv_unit4"`
	ConvUnit5                int     `json:"conv_unit5"`
}

type ProformaInvoiceQueryFilter struct {
	Page         int    `query:"page"`
	Limit        int    `query:"limit" validate:"required"`
	Sort         string `query:"sort"`
	StartDate    *int64 `query:"start_date" validate:"required"`
	EndDate      *int64 `query:"end_date" validate:"required"`
	SalesmanId   []int  `query:"salesman_id" validate:"required,min=1"`
	OutletID     []int  `query:"outlet_id"`
	CustId       string
	ParentCustId string
}

type ProformaInvoiceListResponse struct {
	PoNo           *string  `json:"po_no"`
	RoNo           string   `json:"ro_no"`
	RoDate         *string  `json:"ro_date"`
	FirstIssueDate *string  `json:"first_issue_date"`
	IsProformaInv  *bool    `json:"is_proforma_inv"`
	ProformaInvNo  *string  `json:"proforma_inv_no"`
	OutletID       *int64   `json:"outlet_id"`
	OutletCode     string   `json:"outlet_code"`
	OutletName     string   `json:"outlet_name"`
	OutletAddress  string   `json:"outlet_address"`
	SalesmanId     *int64   `json:"salesman_id"`
	SalesmanCode   string   `json:"salesman_code"`
	SalesmanName   string   `json:"salesman_name"`
	TotalValue     *float64 `json:"total_value"`
	DataStatus     *int64   `json:"data_status"`
	DataStatusName string   `json:"data_status_name"`
}

type PrintProformaInvoiceRequest struct {
	RoNo []string `json:"ro_no" validate:"required,min=1,dive,required"`
}

type PrintProformaInvoiceResponse struct {
	RoNo             []string                      `json:"ro_no"`
	IsProformaInv    bool                          `json:"is_proforma_inv"`
	NoPo             *string                       `json:"no_po"`
	NoSo             *string                       `json:"no_so"`
	CetakUlang       bool                          `json:"cetak_ulang"`
	SalesmanId       *int64                        `json:"salesman_id"`
	SalesmanName     string                        `json:"salesman_name"`
	Notes            *string                       `json:"notes"`
	NoInvoice        *string                       `json:"no_invoice"`
	TglInvoice       *string                       `json:"tgl_invoice"`
	TglJatuhTempo    *string                       `json:"tgl_jatuh_tempo"`
	TypeBayar        string                        `json:"type_bayar"`
	Source           *string                       `json:"source"`
	OutletCode       *string                       `json:"outlet_code"`
	OutletName       *string                       `json:"outlet_name"`
	Address1         *string                       `json:"address1"`
	ZipCode          *string                       `json:"zip_code"`
	Products         []PrintProformaInvoiceProduct `json:"products"`
	Remark           string                        `json:"remark"`
	Gross            *float64                      `json:"gross"`
	PromotionMoney   *float64                      `json:"promotion_money"`
	PromotionProduct float64                       `json:"promotion_product"`
	Discount         *float64                      `json:"discount"`
	VatValue         *float64                      `json:"vat_value"`
	FakturAmount     *float64                      `json:"faktur_amount"`
}

type PrintProformaInvoiceProduct struct {
	ProductCode string   `json:"product_code"`
	ProductName string   `json:"product_name"`
	Qty1        *float64 `json:"qty1"`
	Qty2        *float64 `json:"qty2"`
	Qty3        *float64 `json:"qty3"`
	UnitId1     *string  `json:"unit_id1"`
	UnitId2     *string  `json:"unit_id2"`
	UnitId3     *string  `json:"unit_id3"`
	SellPrice1  *float64 `json:"sell_price1"`
	SellPrice2  *float64 `json:"sell_price2"`
	SellPrice3  *float64 `json:"sell_price3"`
	Total       *float64 `json:"total"`
	Promo1      float64  `json:"promo_1"`
	Promo2      float64  `json:"promo_2"`
	Promo3      float64  `json:"promo_3"`
	Promo4      float64  `json:"promo_4"`
	Promo5      float64  `json:"promo_5"`
	DiscValue   *float64 `json:"disc_value"`
	NettValue   *float64 `json:"nett_value"`
	Remarks     *string  `json:"remarks"`
}
