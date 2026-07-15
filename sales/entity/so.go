package entity

var payTypeSoName = map[int64]string{
	1: "Cash",
	2: "Check",
	3: "Transfer",
	4: "Credit",
}

type GeneralQueryFilter struct {
	CustId       string
	ParentCustId string
	From         *int64 `query:"from" validate:"required_with=To,omitempty,gte=1000000000"`
	To           *int64 `query:"to" validate:"required_with=From,omitempty,lte=9999999999,gtefield=From"`
	Page         int    `query:"page"`
	Limit        int    `query:"limit" validate:""`
	Query        string `query:"q"`
	Mode         string `query:"mode"`
	Sort         string `query:"sort"`
	IsActive     *int   `query:"is_active"`
}

type CreateSoBody struct {
	CustID        string                   `json:"cust_id"`
	SoDate        *string                  `json:"so_date"`
	SysDate       *string                  `json:"sys_date"`
	OutletID      *int64                   `json:"outlet_id"`
	OutletTaxNo   *string                  `json:"outlet_tax_no"`
	PoNo          *string                  `json:"po_no"`
	VehicleNo     *string                  `json:"vehicle_no"`
	InvoiceNo     *string                  `json:"invoice_no"`
	InvoiceDate   *string                  `json:"invoice_date"`
	DeliveryDate  *string                  `json:"delivery_date"`
	PayType       *int64                   `json:"pay_type"`
	SumNo         *string                  `json:"sum_no"`
	DataSource    *int64                   `json:"data_source"`
	MobileID      *int64                   `json:"mobile_id"`
	SubTotal      *float64                 `json:"sub_total"`
	Disc          *float64                 `json:"disc"`
	DiscValue     *float64                 `json:"disc_value"`
	PromoValue    *float64                 `json:"promo_value"`
	CashDiscValue *float64                 `json:"cash_disc_value"`
	TotDisc1      *float64                 `json:"tot_disc_1"`
	TotDisc2      *float64                 `json:"tot_disc_2"`
	Vat           *float64                 `json:"vat"`
	VatValue      *float64                 `json:"vat_value"`
	TotalInv      float64                  `json:"total_inv"`
	TotalInvIn    *float64                 `json:"total_inv_in"`
	RoundDiff     *float64                 `json:"round_diff"`
	DataStatus    *int64                   `json:"data_status"`
	CreatedBy     *int64                   `json:"created_by"`
	DueDate       *string                  `json:"due_date"`
	Details       CreateSoDetBodyWithGroup `json:"details"`
}
type UpdateSoBody struct {
	CustID        string                   `json:"cust_id"`
	SoDate        *string                  `json:"so_date"`
	SysDate       *string                  `json:"sys_date"`
	OutletID      *int64                   `json:"outlet_id"`
	OutletTaxNo   *string                  `json:"outlet_tax_no"`
	PoNo          *string                  `json:"po_no"`
	VehicleNo     *string                  `json:"vehicle_no"`
	InvoiceNo     *string                  `json:"invoice_no"`
	InvoiceDate   *string                  `json:"invoice_date"`
	DeliveryDate  *string                  `json:"delivery_date"`
	PayType       *int64                   `json:"pay_type"`
	SumNo         *string                  `json:"sum_no"`
	DataSource    *int64                   `json:"data_source"`
	MobileID      *int64                   `json:"mobile_id"`
	SubTotal      *float64                 `json:"sub_total"`
	Disc          *float64                 `json:"disc"`
	DiscValue     *float64                 `json:"disc_value"`
	PromoValue    *float64                 `json:"promo_value"`
	CashDiscValue *float64                 `json:"cash_disc_value"`
	TotDisc1      *float64                 `json:"tot_disc_1"`
	TotDisc2      *float64                 `json:"tot_disc_2"`
	Vat           *float64                 `json:"vat"`
	VatValue      *float64                 `json:"vat_value"`
	TotalInv      float64                  `json:"total_inv"`
	TotalInvIn    *float64                 `json:"total_inv_in"`
	RoundDiff     *float64                 `json:"round_diff"`
	DataStatus    *int64                   `json:"data_status"`
	UpdatedBy     int64                    `json:"updated_by"`
	DueDate       *string                  `json:"due_date"`
	Details       UpdateSoDetBodyWithGroup `json:"details"`
}
type SoResponse struct {
	SoNo          string                 `json:"so_no"`
	SoDate        *string                `json:"so_date"`
	SysDate       *string                `json:"sys_date"`
	OutletID      *int64                 `json:"outlet_id"`
	OutletCode    string                 `json:"outlet_code"`
	OutletName    string                 `json:"outlet_name"`
	OutletTaxNo   *string                `json:"outlet_tax_no"`
	PoNo          *string                `json:"po_no"`
	VehicleNo     *string                `json:"vehicle_no"`
	InvoiceNo     *string                `json:"invoice_no"`
	InvoiceDate   *string                `json:"invoice_date"`
	DeliveryDate  *string                `json:"delivery_date"`
	PayType       *int64                 `json:"pay_type"`
	PayTypeName   string                 `json:"pay_type_name"`
	SumNo         *string                `json:"sum_no"`
	DataSource    *int64                 `json:"data_source"`
	MobileID      *int64                 `json:"mobile_id"`
	SubTotal      *float64               `json:"sub_total"`
	Disc          *float64               `json:"disc"`
	DiscValue     *float64               `json:"disc_value"`
	PromoValue    *float64               `json:"promo_value"`
	CashDiscValue *float64               `json:"cash_disc_value"`
	TotDisc1      *float64               `json:"tot_disc_1"`
	TotDisc2      *float64               `json:"tot_disc_2"`
	Vat           *float64               `json:"vat"`
	VatValue      *float64               `json:"vat_value"`
	TotalInv      float64                `json:"total_inv"`
	TotalInvIn    *float64               `json:"total_inv_in"`
	RoundDiff     *float64               `json:"round_diff"`
	DataStatus    *int64                 `json:"data_status"`
	UpdatedAt     string                 `json:"updated_at"`
	UpdatedByName string                 `json:"updated_by_name"`
	DueDate       *string                `json:"due_date"`
	Details       SoDetResponseWithGroup `json:"details"`
}

func (so SoResponse) GeneratePayTypeName() string {
	if so.PayType != nil {
		return payTypeSoName[*so.PayType]
	}
	return ""
}

type SoListResponse struct {
	SoNo          string   `json:"so_no"`
	SoDate        *string  `json:"so_date"`
	SysDate       *string  `json:"sys_date"`
	OutletID      *int64   `json:"outlet_id"`
	OutletCode    string   `json:"outlet_code"`
	OutletName    string   `json:"outlet_name"`
	OutletTaxNo   *string  `json:"outlet_tax_no"`
	PoNo          *string  `json:"po_no"`
	VehicleNo     *string  `json:"vehicle_no"`
	InvoiceNo     *string  `json:"invoice_no"`
	InvoiceDate   *string  `json:"invoice_date"`
	DeliveryDate  *string  `json:"delivery_date"`
	PayType       *int64   `json:"pay_type"`
	PayTypeName   string   `json:"pay_type_name"`
	SumNo         *string  `json:"sum_no"`
	DataSource    *int64   `json:"data_source"`
	MobileID      *int64   `json:"mobile_id"`
	SubTotal      *float64 `json:"sub_total"`
	Disc          *float64 `json:"disc"`
	DiscValue     *float64 `json:"disc_value"`
	PromoValue    *float64 `json:"promo_value"`
	CashDiscValue *float64 `json:"cash_disc_value"`
	TotDisc1      *float64 `json:"tot_disc_1"`
	TotDisc2      *float64 `json:"tot_disc_2"`
	Vat           *float64 `json:"vat"`
	VatValue      *float64 `json:"vat_value"`
	TotalInv      float64  `json:"total_inv"`
	TotalInvIn    *float64 `json:"total_inv_in"`
	RoundDiff     *float64 `json:"round_diff"`
	DataStatus    *int64   `json:"data_status"`
	UpdatedAt     string   `json:"updated_at"`
	DueDate       *string  `json:"due_date"`
	UpdatedByName string   `json:"updated_by_name"`
}

func (so SoListResponse) GeneratePayTypeName() string {
	if so.PayType != nil {
		return payTypeSoName[*so.PayType]
	}
	return ""
}

type DetailSoParams struct {
	SoNo string `params:"so_no" validate:"required"`
}

type UpdateSoParams struct {
	SoNo string `params:"so_no" validate:"required"`
}
