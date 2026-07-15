package entity

type CreateVanBsUlBody struct {
	CustID      string                 `json:"cust_id"`
	VanBsUlNo   string                 `json:"van_bs_ul_no"`
	VanBsUlDate *string                `json:"van_bs_ul_date"`
	TrCode      *string                `json:"tr_code" validate:"required,len=3"`
	WhID        *int64                 `json:"wh_id"`
	SalesmanID  *int64                 `json:"salesman_id"`
	RefNo       *string                `json:"ref_no"`
	TotEmbInc   *float64               `json:"tot_emb_inc"`
	TotEmbExc   *float64               `json:"tot_emb_exc"`
	SubTotal    *float64               `json:"sub_total"`
	Vat         *float64               `json:"vat"`
	VatValue    *float64               `json:"vat_value"`
	VatLg       *float64               `json:"vat_lg"`
	VatLgValue  *float64               `json:"vat_lg_value"`
	Total       *float64               `json:"total"`
	VatBg       *float64               `json:"vat_bg"`
	VatBgValue  *float64               `json:"vat_bg_value"`
	DataStatus  *int64                 `json:"data_status"`
	CreatedBy   *int64                 `json:"created_by"`
	Details     []CreateVanBsUlDetBody `json:"details"`
}

type VanBsUlResponse struct {
	VanBsUlNo     string               `json:"van_bs_ul_no"`
	VanBsUlDate   *string              `json:"van_bs_ul_date"`
	TrCode        *string              `json:"tr_code"`
	WhID          *int64               `json:"wh_id"`
	WhCode        string               `json:"wh_code"`
	WhName        string               `json:"wh_name"`
	SalesmanID    *int64               `json:"salesman_id"`
	SalesmanCode  string               `json:"salesman_code"`
	SalesmanName  string               `json:"salesman_name"`
	RefNo         *string              `json:"ref_no"`
	TotEmbInc     *float64             `json:"tot_emb_inc"`
	TotEmbExc     *float64             `json:"tot_emb_exc"`
	SubTotal      *float64             `json:"sub_total"`
	Vat           *float64             `json:"vat"`
	VatValue      *float64             `json:"vat_value"`
	VatLg         *float64             `json:"vat_lg"`
	VatLgValue    *float64             `json:"vat_lg_value"`
	Total         *float64             `json:"total"`
	VatBg         *float64             `json:"vat_bg"`
	VatBgValue    *float64             `json:"vat_bg_value"`
	DataStatus    *int64               `json:"data_status"`
	UpdatedAt     string               `json:"updated_at"`
	UpdatedByName string               `json:"updated_by_name"`
	IsClosed      bool                 `json:"is_closed"`
	ClosedBy      int64                `json:"closed_by"`
	ClosedByName  string               `json:"closed_by_name"`
	ClosedAt      string               `json:"closed_at"`
	Details       []VanBsUlDetResponse `json:"details"`
}
type DetailVanBsUlParams struct {
	VanBsUlNo string `params:"van_bs_ul_no" validate:"required"`
}
type DeleteVanBsUlParams struct {
	VanBsUlNo string `params:"van_bs_ul_no" validate:"required"`
}
type UpdateVanBsUlParams struct {
	VanBsUlNo string `params:"van_bs_ul_no" validate:"required"`
}
type VanBsUlListResponse struct {
	VanBsUlNo     string   `json:"van_bs_ul_no"`
	VanBsUlDate   *string  `json:"van_bs_ul_date"`
	TrCode        *string  `json:"tr_code"`
	WhID          *int64   `json:"wh_id"`
	WhCode        string   `json:"wh_code"`
	WhName        string   `json:"wh_name"`
	SalesmanID    *int64   `json:"salesman_id"`
	SalesmanCode  string   `json:"salesman_code"`
	SalesmanName  string   `json:"salesman_name"`
	RefNo         *string  `json:"ref_no"`
	TotEmbInc     *float64 `json:"tot_emb_inc"`
	TotEmbExc     *float64 `json:"tot_emb_exc"`
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
	IsClosed      bool     `json:"is_closed"`
	ClosedBy      int64    `json:"closed_by"`
	ClosedByName  string   `json:"closed_by_name"`
	ClosedAt      string   `json:"closed_at"`
}
type UpdateVanBsUlBody struct {
	CustID      string                 `json:"cust_id"`
	VanBsUlNo   string                 `json:"van_bs_ul_no"`
	VanBsUlDate *string                `json:"van_bs_ul_date"`
	TrCode      *string                `json:"tr_code"`
	WhID        *int64                 `json:"wh_id"`
	SalesmanID  *int64                 `json:"salesman_id"`
	RefNo       *string                `json:"ref_no"`
	TotEmbInc   *float64               `json:"tot_emb_inc"`
	TotEmbExc   *float64               `json:"tot_emb_exc"`
	SubTotal    *float64               `json:"sub_total"`
	Vat         *float64               `json:"vat"`
	VatValue    *float64               `json:"vat_value"`
	VatLg       *float64               `json:"vat_lg"`
	VatLgValue  *float64               `json:"vat_lg_value"`
	Total       *float64               `json:"total"`
	VatBg       *float64               `json:"vat_bg"`
	VatBgValue  *float64               `json:"vat_bg_value"`
	DataStatus  *int64                 `json:"data_status"`
	CreatedBy   *int64                 `json:"created_by"`
	UpdatedBy   int64                  `json:"updated_by"`
	Details     []UpdateVanBsUlDetBody `json:"details"`
}
