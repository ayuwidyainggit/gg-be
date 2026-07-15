package entity

type CreateConsignmentBody struct {
	CustID     string                 `json:"cust_id"`
	ConsDate   *string                `json:"cons_date"`
	ConsType   *int64                 `json:"cons_type"`
	OutletID   *int64                 `json:"outlet_id"`
	SalesmanID *int64                 `json:"salesman_id"`
	DeliveryNo *string                `json:"delivery_no"`
	WhID       *int64                 `json:"wh_id"`
	Notes      *string                `json:"notes"`
	SubTotal   *float64               `json:"sub_total"`
	Vat        *float64               `json:"vat"`
	VatValue   *float64               `json:"vat_value"`
	VatLg      *float64               `json:"vat_lg"`
	VatLgValue *float64               `json:"vat_lg_value"`
	Total      *float64               `json:"total"`
	VatBg      *float64               `json:"vat_bg"`
	VatBgValue *float64               `json:"vat_bg_value"`
	DataStatus *int64                 `json:"data_status"`
	CreatedBy  *int64                 `json:"created_by"`
	Details    []CreateConsignDetBody `json:"details"`
}

type ConsignmentResponse struct {
	ConsNo        string               `json:"cons_no"`
	ConsDate      *string              `json:"cons_date"`
	ConsType      *int64               `json:"cons_type"`
	OutletID      *int64               `json:"outlet_id"`
	OutletCode    string               `json:"outlet_code"`
	OutletName    string               `json:"outlet_name"`
	SalesmanID    *int64               `json:"salesman_id"`
	SalesmanCode  string               `json:"salesman_code"`
	SalesmanName  string               `json:"salesman_name"`
	DeliveryNo    *string              `json:"delivery_no"`
	WhID          *int64               `json:"wh_id"`
	Notes         *string              `json:"notes"`
	SubTotal      *float64             `json:"sub_total"`
	Vat           *float64             `json:"vat"`
	VatValue      *float64             `json:"vat_value"`
	VatLg         *float64             `json:"vat_lg"`
	VatLgValue    *float64             `json:"vat_lg_value"`
	Total         *float64             `json:"total"`
	VatBg         *float64             `json:"vat_bg"`
	VatBgValue    *float64             `json:"vat_bg_value"`
	DataStatus    *int64               `json:"data_status"`
	CreatedBy     *int64               `json:"created_by"`
	UpdatedAt     string               `json:"updated_at"`
	UpdatedByName string               `json:"updated_by_name"`
	Details       []ConsignDetResponse `json:"details"`
}

type ConsignmentListResponse struct {
	ConsNo        string   `json:"cons_no"`
	ConsDate      *string  `json:"cons_date"`
	ConsType      *int64   `json:"cons_type"`
	OutletID      *int64   `json:"outlet_id"`
	OutletCode    string   `json:"outlet_code"`
	OutletName    string   `json:"outlet_name"`
	SalesmanID    *int64   `json:"salesman_id"`
	SalesmanCode  string   `json:"salesman_code"`
	SalesmanName  string   `json:"salesman_name"`
	DeliveryNo    *string  `json:"delivery_no"`
	WhID          *int64   `json:"wh_id"`
	Notes         *string  `json:"notes"`
	SubTotal      *float64 `json:"sub_total"`
	Vat           *float64 `json:"vat"`
	VatValue      *float64 `json:"vat_value"`
	VatLg         *float64 `json:"vat_lg"`
	VatLgValue    *float64 `json:"vat_lg_value"`
	Total         *float64 `json:"total"`
	VatBg         *float64 `json:"vat_bg"`
	VatBgValue    *float64 `json:"vat_bg_value"`
	DataStatus    *int64   `json:"data_status"`
	UpdatedAt     string   `json:"updated_at"`
	UpdatedByName string   `json:"updated_by_name"`
}
type UpdateConsignmentBody struct {
	ConsNo     string                 `json:"cons_no"`
	CustID     string                 `json:"cust_id"`
	ConsDate   *string                `json:"cons_date"`
	ConsType   *int64                 `json:"cons_type"`
	OutletID   *int64                 `json:"outlet_id"`
	SalesmanID *int64                 `json:"salesman_id"`
	DeliveryNo *string                `json:"delivery_no"`
	WhID       *int64                 `json:"wh_id"`
	Notes      *string                `json:"notes"`
	SubTotal   *float64               `json:"sub_total"`
	Vat        *float64               `json:"vat"`
	VatValue   *float64               `json:"vat_value"`
	VatLg      *float64               `json:"vat_lg"`
	VatLgValue *float64               `json:"vat_lg_value"`
	Total      *float64               `json:"total"`
	VatBg      *float64               `json:"vat_bg"`
	VatBgValue *float64               `json:"vat_bg_value"`
	DataStatus *int64                 `json:"data_status"`
	CreatedBy  *int64                 `json:"created_by"`
	UpdatedBy  int64                  `json:"updated_by"`
	Details    []UpdateConsignDetBody `json:"details"`
}
type DetailConsignmentParams struct {
	ConsNo string `params:"cons_no" validate:"required"`
}

type UpdateConsignmentParams struct {
	ConsNo string `params:"cons_no" validate:"required"`
}
