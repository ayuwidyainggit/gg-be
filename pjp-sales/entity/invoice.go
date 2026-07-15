package entity

type InvoiceQueryFilter struct {
	InvoiceNo    []string `query:"invoice_no"`
	SalesmanId   []int    `query:"salesman_id"`
	OutletID     []int    `query:"outlet_id"`
	Status       []int    `query:"status"`
	CustId       string
	ParentCustId string
	From         *int64 `query:"from" validate:"required_with=To,omitempty,gte=1000000000"`
	To           *int64 `query:"to" validate:"required_with=From,omitempty,lte=9999999999,gtefield=From"`
	Page         int    `query:"page"`
	Limit        int    `query:"limit"`
	Query        string `query:"q"`
	Mode         string `query:"mode"`
	Sort         string `query:"sort"`
	IsActive     *int   `query:"is_active"`
	IsInvoice    *bool  `query:"is_invoice"`
}

type InvoiceResponse struct {
	SalesmanId        *int64                  `json:"salesman_id"`
	SalesName         string                  `json:"sales_name"`
	WhId              *int64                  `json:"wh_id"`
	WhCode            string                  `json:"wh_code"`
	WhName            string                  `json:"wh_name"`
	WhLatitude        string                  `json:"wh_latitude"`
	WhLongitude       string                  `json:"wh_longitude"`
	OutletID          *int64                  `json:"outlet_id"`
	OutletCode        string                  `json:"outlet_code"`
	OutletName        string                  `json:"outlet_name"`
	OutletAddress     string                  `json:"outlet_address"`
	OutletLatitude    string                  `json:"outlet_latitude"`
	OutletLongitude   string                  `json:"outlet_longitude"`
	DeliveryDate      *string                 `json:"delivery_date"`
	OrderNo           *string                 `json:"order_no"`
	PoNo              *string                 `json:"po_no"`
	VehicleNo         *string                 `json:"vehicle_no"`
	PayType           *int64                  `json:"pay_type"`
	PayTypeName       string                  `json:"pay_type_name"`
	ReffNo            *string                 `json:"reff_no"`
	MobileID          *int64                  `json:"mobile_id"`
	SubTotal          *float64                `json:"sub_total"`
	Disc              *float64                `json:"disc"`
	DiscValue         *float64                `json:"disc_value"`
	PromoValue        *float64                `json:"promo_value"`
	PromoValueFinal   *float64                `json:"promo_value_final"`
	PromoBgValue      *float64                `json:"promo_bg_value"`
	PromoBgValueFinal *float64                `json:"promo_bg_value_final"`
	CashDiscValue     *float64                `json:"cash_disc_value"`
	TotDisc1          *float64                `json:"tot_disc1"`
	TotDisc2          *float64                `json:"tot_disc2"`
	Vat               *float64                `json:"vat"`
	VatValue          *float64                `json:"vat_value"`
	Total             *float64                `json:"total"`
	DataStatus        *int64                  `json:"data_status"`
	DataStatusName    string                  `json:"data_status_name"`
	DataSource        *int64                  `json:"data_source"`
	DueDate           *string                 `json:"due_date"`
	InvoiceNo         *string                 `json:"invoice_no"`
	InvoiceDate       *string                 `json:"invoice_date"`
	Details           InvoiceDetReadWithGroup `json:"details"`
}

func (invoice InvoiceResponse) GeneratePayTypeName() string {
	if invoice.PayType != nil {
		return payTypeName[*invoice.PayType]
	}
	return ""
}

func (invoice InvoiceResponse) GenerateDataStatusName() string {
	if invoice.DataStatus != nil {
		return dataStatusName[*invoice.DataStatus]
	}
	return ""
}

type DetailInvoiceParams struct {
	RoNo string `params:"ro_no" validate:"required"`
}
type DetailInvoiceParamsInv struct {
	InvoiceNo string `params:"invoice_no" validate:"required"`
}
type DeleteInvoiceParams struct {
	RoNo string `params:"ro_no" validate:"required"`
}

type UpdateInvoiceParams struct {
	RoNo string `params:"ro_no" validate:"required"`
}

type InvoiceListResponse struct {
	SalesmanId        *int64               `json:"salesman_id"`
	SalesmanCode      string               `json:"salesman_code"`
	SalesName         string               `json:"sales_name"`
	WhId              *int64               `json:"wh_id"`
	WhCode            string               `json:"wh_code"`
	WhName            string               `json:"wh_name"`
	WhLatitude        string               `json:"wh_latitude"`
	WhLongitude       string               `json:"wh_longitude"`
	OutletID          *int64               `json:"outlet_id"`
	OutletCode        string               `json:"outlet_code"`
	OutletName        string               `json:"outlet_name"`
	OutletAddress     string               `json:"outlet_address"`
	OutletLatitude    string               `json:"outlet_latitude"`
	OutletLongitude   string               `json:"outlet_longitude"`
	DeliveryDate      *string              `json:"delivery_date"`
	OrderNo           string               `json:"order_no"`
	PoNo              *string              `json:"po_no"`
	VehicleNo         *string              `json:"vehicle_no"`
	PayType           *int64               `json:"pay_type"`
	PayTypeName       string               `json:"pay_type_name"`
	ReffNo            *string              `json:"reff_no"`
	MobileID          *int64               `json:"mobile_id"`
	Disc              *float64             `json:"disc"`
	DiscValue         *float64             `json:"disc_value"`
	PromoValue        *float64             `json:"promo_value"`
	PromoValueFinal   *float64             `json:"promo_value_final"`
	PromoBgValue      *float64             `json:"promo_bg_value"`
	PromoBgValueFinal *float64             `json:"promo_bg_value_final"`
	SubTotal          *float64             `json:"sub_total"`
	Vat               *float64             `json:"vat"`
	VatValue          *float64             `json:"vat_value"`
	Total             *float64             `json:"total"`
	DataStatus        *int64               `json:"data_status"`
	DataStatusName    string               `json:"data_status_name"`
	DueDate           *string              `json:"due_date"`
	InvoiceNo         string               `json:"invoice_no"`
	InvoiceDate       string               `json:"invoice_date"`
	TotalWeight       float64              `json:"total_weight"`
	TotalVolume       float64              `json:"total_volume"`
	Details           []InvoiceDetResponse `json:"details,omitempty"`
	IsPrinted         *bool                `json:"is_printed"`
	PrintedBy         *int64               `json:"printed_by"`
	PrintedByName     *string              `json:"printed_by_name"`
	PrintedAt         *string              `json:"printed_at"`
}

func (invoice InvoiceListResponse) GenerateDataStatusName() string {
	if invoice.DataStatus != nil {
		return dataStatusName[*invoice.DataStatus]
	}
	return ""
}

func (invoice InvoiceListResponse) GeneratePayTypeName() string {
	if invoice.PayType != nil {
		return payTypeName[*invoice.PayType]
	}
	return ""
}

type UpdateInvoiceBody struct {
	CustId        string                    `json:"cust_id"`
	RoNo          string                    `json:"ro_no"`
	RoDate        *string                   `json:"ro_date"`
	ValDate       *string                   `json:"val_date"`
	DueDate       *string                   `json:"due_date"`
	SalesmanId    *int64                    `json:"salesman_id"`
	WhId          *int64                    `json:"wh_id"`
	OutletID      *int64                    `json:"outlet_id"`
	DeliveryDate  *string                   `json:"delivery_date"`
	OrderNo       *string                   `json:"order_no"`
	PoNo          *string                   `json:"po_no"`
	VehicleNo     *string                   `json:"vehicle_no"`
	PayType       *int64                    `json:"pay_type"`
	ReffNo        *string                   `json:"reff_no"`
	MobileID      *int64                    `json:"mobile_id"`
	SubTotal      *float64                  `json:"sub_total"`
	Disc          *float64                  `json:"disc"`
	DiscValue     *float64                  `json:"disc_value"`
	PromoValue    *float64                  `json:"promo_value"`
	CashDiscValue *float64                  `json:"cash_disc_value"`
	TotDisc1      *float64                  `json:"tot_disc1"`
	TotDisc2      *float64                  `json:"tot_disc2"`
	Vat           *float64                  `json:"vat"`
	VatValue      *float64                  `json:"vat_value"`
	Total         *float64                  `json:"total"`
	DataStatus    *int64                    `json:"data_status"`
	CreatedBy     *int64                    `json:"created_by"`
	CreatedAt     *string                   `json:"created_at"`
	UpdatedBy     int64                     `json:"updated_by"`
	InvoiceNo     *string                   `json:"invoice_no"`
	InvoiceDate   *string                   `json:"invoice_date"`
	Details       UpdateInvoiceDetWithGroup `json:"details"`
}

type UpdateInvoiceBodyV1 struct {
	CustId        string                    `json:"cust_id"`
	RoNo          string                    `json:"ro_no"`
	RoDate        *string                   `json:"ro_date"`
	ValDate       *string                   `json:"val_date"`
	DueDate       *string                   `json:"due_date"`
	SalesmanId    *int64                    `json:"salesman_id"`
	WhId          *int64                    `json:"wh_id"`
	OutletID      *int64                    `json:"outlet_id"`
	DeliveryDate  *string                   `json:"delivery_date"`
	OrderNo       *string                   `json:"order_no"`
	PoNo          *string                   `json:"po_no"`
	VehicleNo     *string                   `json:"vehicle_no"`
	PayType       *int64                    `json:"pay_type"`
	ReffNo        *string                   `json:"reff_no"`
	MobileID      *int64                    `json:"mobile_id"`
	SubTotal      *float64                  `json:"sub_total"`
	Disc          *float64                  `json:"disc"`
	DiscValue     *float64                  `json:"disc_value"`
	PromoValue    *float64                  `json:"promo_value"`
	CashDiscValue *float64                  `json:"cash_disc_value"`
	TotDisc1      *float64                  `json:"tot_disc1"`
	TotDisc2      *float64                  `json:"tot_disc2"`
	Vat           *float64                  `json:"vat"`
	VatValue      *float64                  `json:"vat_value"`
	Total         *float64                  `json:"total"`
	DataStatus    *int64                    `json:"data_status"`
	CreatedBy     *int64                    `json:"created_by"`
	CreatedAt     *string                   `json:"created_at"`
	UpdatedBy     int64                     `json:"updated_by"`
	Details       UpdateInvoiceDetWithGroup `json:"details"`
}

type BulkUpdateInvoiceBody struct {
	Orders []UpdateInvoiceBody `json:"orders" validate:"min=1"`
}

type BulkUpdateInvoiceResponse struct {
	Orders []InvoiceResponse `json:"orders"`
}
