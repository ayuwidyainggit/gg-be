package entity

type CreateVanLoBody struct {
	CustID       string              `json:"cust_id"`
	VanLoNo      string              `json:"van_lo_no"`
	VanLoDate    *string             `json:"van_lo_date"`
	TrCode       *string             `json:"tr_code" validate:"required,len=3"`
	WhID         *int64              `json:"wh_id"`
	SalesmanID   *int64              `json:"salesman_id"`
	Notes        *string             `json:"notes"`
	TotEmbInc    *float64            `json:"tot_emb_inc"`
	TotEmbExc    *float64            `json:"tot_emb_exc"`
	SubTotal     *float64            `json:"sub_total"`
	Vat          *float64            `json:"vat"`
	VatValue     *float64            `json:"vat_value"`
	VatLg        *float64            `json:"vat_lg"`
	VatLgValue   *float64            `json:"vat_lg_value"`
	Total        *float64            `json:"total"`
	VatBg        *float64            `json:"vat_bg"`
	VatBgValue   *float64            `json:"vat_bg_value"`
	DataStatus   *int64              `json:"data_status"`
	CreatedBy    int64               `json:"created_by"`
	IsDel        *bool               `json:"is_del"`
	Details      VanLoDetCreateGroup `json:"details"`
}

type VanLoListResponse struct {
	CustID        string   `json:"cust_id"`
	VanLoNo       string   `json:"van_lo_no"`
	VanLoDate     *string  `json:"van_lo_date"`
	TrCode        *string  `json:"tr_code"`
	WhID          *int64   `json:"wh_id"`
	SalesmanID    *int64   `json:"salesman_id"`
	Notes         *string  `json:"notes"`
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
	CreatedBy     *int64   `json:"created_by"`
	IsDel         *bool    `json:"is_del"`
	UpdatedAt     string   `json:"updated_at"`
	UpdatedByName string   `json:"updated_by_name"`
	IsClosed      bool     `json:"is_closed"`
	ClosedBy      int64    `json:"closed_by"`
	ClosedByName  string   `json:"closed_by_name"`
	ClosedAt      string   `json:"closed_at"`
}
type VanLoUpdateBody struct {
	CustID     string              `json:"cust_id"`
	VanLoNo    string              `json:"van_lo_no"`
	VanLoDate  *string             `json:"van_lo_date"`
	TrCode     *string             `json:"tr_code"`
	WhID       *int64              `json:"wh_id"`
	SalesmanID *int64              `json:"salesman_id"`
	Notes      *string             `json:"notes"`
	TotEmbInc  *float64            `json:"tot_emb_inc"`
	TotEmbExc  *float64            `json:"tot_emb_exc"`
	SubTotal   *float64            `json:"sub_total"`
	Vat        *float64            `json:"vat"`
	VatValue   *float64            `json:"vat_value"`
	VatLg      *float64            `json:"vat_lg"`
	VatLgValue *float64            `json:"vat_lg_value"`
	Total      *float64            `json:"total"`
	VatBg      *float64            `json:"vat_bg"`
	VatBgValue *float64            `json:"vat_bg_value"`
	DataStatus *int64              `json:"data_status"`
	UpdatedBy  int64               `json:"updated_by"`
	IsDel      *bool               `json:"is_del"`
	Details    VanLoDetUpdateGroup `json:"details"`
}
type VanLoResponse struct {
	CustID        string            `json:"cust_id"`
	VanLoNo       string            `json:"van_lo_no"`
	VanLoDate     string            `json:"van_lo_date"`
	TrCode        *string           `json:"tr_code"`
	WhID          *int64            `json:"wh_id"`
	SalesmanID    *int64            `json:"salesman_id"`
	Notes         *string           `json:"notes"`
	TotEmbInc     *float64          `json:"tot_emb_inc"`
	TotEmbExc     *float64          `json:"tot_emb_exc"`
	SubTotal      *float64          `json:"sub_total"`
	Vat           *float64          `json:"vat"`
	VatValue      *float64          `json:"vat_value"`
	VatLg         *float64          `json:"vat_lg"`
	VatLgValue    *float64          `json:"vat_lg_value"`
	Total         *float64          `json:"total"`
	VatBg         *float64          `json:"vat_bg"`
	VatBgValue    *float64          `json:"vat_bg_value"`
	DataStatus    *int64            `json:"data_status"`
	CreatedBy     *int64            `json:"created_by"`
	IsDel         *bool             `json:"is_del"`
	UpdatedAt     string            `json:"updated_at"`
	UpdatedByName string            `json:"updated_by_name"`
	IsClosed      bool              `json:"is_closed"`
	ClosedBy      int64             `json:"closed_by"`
	ClosedByName  string            `json:"closed_by_name"`
	ClosedAt      string            `json:"closed_at"`
	Details       VanLoDetReadGroup `json:"details"`
}
type DetailVanLoParams struct {
	VanLoNo string `params:"van_lo_no" validate:"required"`
}
type UpdateVanLoParams struct {
	VanLoNo string `params:"van_lo_no" validate:"required"`
}
