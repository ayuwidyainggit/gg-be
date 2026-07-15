package entity

type CreateVanUlBody struct {
	CustID     string              `json:"cust_id"`
	VanUlNo    string              `json:"van_ul_no"`
	VanUlDate  *string             `json:"van_ul_date" validate:"required"`
	TrCode     *string             `json:"tr_code" validate:"required,len=3"`
	WhID       *int64              `json:"wh_id" validate:"required"`
	WhCode     string              `json:"wh_code"`
	WhName     string              `json:"wh_name"`
	SalesmanID *int64              `json:"salesman_id" validate:"required"`
	Notes      *string             `json:"notes"`
	TotEmbInc  *float64            `json:"tot_emb_inc"`
	TotEmbExc  *float64            `json:"tot_emb_exc"`
	SubTotal   *float64            `json:"sub_total" validate:"required"`
	Vat        *float64            `json:"vat"`
	VatValue   *float64            `json:"vat_value"`
	VatLg      *float64            `json:"vat_lg"`
	VatLgValue *float64            `json:"vat_lg_value"`
	Total      *float64            `json:"total"`
	VatBg      *float64            `json:"vat_bg"`
	VatBgValue *float64            `json:"vat_bg_value"`
	WeekIDOb   int                 `json:"week_id_ob"`
	DataStatus *int64              `json:"data_status"`
	CreatedBy  int64               `json:"created_by"`
	Details    VanUlDetCreateGroup `json:"details"`
}

type VanUlListResponse struct {
	CustID        string   `json:"cust_id"`
	VanUlNo       string   `json:"van_ul_no"`
	VanUlDate     string   `json:"van_ul_date"`
	TrCode        *string  `json:"tr_code"`
	WhID          *int64   `json:"wh_id"`
	WhCode        string   `json:"wh_code"`
	WhName        string   `json:"wh_name"`
	SalesmanID    *int64   `json:"salesman_id"`
	SalesmanName  string   `json:"salesman_name"`
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
	WeekIDOb      int      `json:"week_id_ob"`
	DataStatus    *int64   `json:"data_status"`
	UpdatedAt     string   `json:"updated_at"`
	UpdatedByName string   `json:"updated_by_name"`
	IsClosed      bool     `json:"is_closed"`
	ClosedBy      int64    `json:"closed_by"`
	ClosedByName  string   `json:"closed_by_name"`
	ClosedAt      string   `json:"closed_at"`
}
type VanUlResponse struct {
	CustID        string            `json:"cust_id"`
	VanUlNo       string            `json:"van_ul_no"`
	VanUlDate     string            `json:"van_ul_date"`
	TrCode        *string           `json:"tr_code"`
	WhID          *int64            `json:"wh_id"`
	WhCode        string            `json:"wh_code"`
	WhName        string            `json:"wh_name"`
	SalesmanID    *int64            `json:"salesman_id"`
	SalesmanName  string            `json:"salesman_name"`
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
	WeekIDOb      int               `json:"week_id_ob"`
	DataStatus    *int64            `json:"data_status"`
	UpdatedAt     string            `json:"updated_at"`
	UpdatedByName string            `json:"updated_by_name"`
	IsClosed      bool              `json:"is_closed"`
	ClosedBy      int64             `json:"closed_by"`
	ClosedByName  string            `json:"closed_by_name"`
	ClosedAt      string            `json:"closed_at"`
	Details       VanUlDetReadGroup `json:"details"`
}
type VanUlUpdateBody struct {
	CustID     string              `json:"cust_id"`
	VanUlNo    string              `json:"van_ul_no"`
	VanUlDate  *string             `json:"van_ul_date"`
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
	WeekIDOb   int                 `json:"week_id_ob"`
	DataStatus *int64              `json:"data_status"`
	UpdatedBy  int64               `json:"updated_by"`
	Details    VanUlDetUpdateGroup `json:"details"`
}

type DetailVanUlParams struct {
	VanUlNo string `params:"van_ul_no" validate:"required"`
}
type UpdateVanUlParams struct {
	VanUlNo string `params:"van_ul_no" validate:"required"`
}
