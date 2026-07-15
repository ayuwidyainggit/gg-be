package entity

import "time"

type CreateGdsBody struct {
	CustID     string             `json:"cust_id"`
	GdsDate    string             `json:"gds_date" validate:"required"`
	TrCode     string             `json:"tr_code" validate:"required,len=3"`
	RefNo      string             `json:"ref_no"`
	WhID       int64              `json:"wh_id" validate:"required"`
	SupID      int64              `json:"sup_id" validate:"required"`
	SubTotal   float64            `json:"sub_total"`
	Vat        float64            `json:"vat"`
	VatValue   float64            `json:"vat_value"`
	VatLg      float64            `json:"vat_lg"`
	VatLgValue float64            `json:"vat_lg_value"`
	Total      float64            `json:"total"`
	VatBg      float64            `json:"vat_bg"`
	VatBgValue float64            `json:"vat_bg_value"`
	TotEmbInc  float64            `json:"tot_emb_inc"`
	TotEmbExc  float64            `json:"tot_emb_exc"`
	Notes      string             `json:"notes"`
	DataStatus int64              `json:"data_status" validate:"required"`
	CreatedBy  int64              `json:"created_by"`
	Details    []CreateGdsDetBody `json:"details" validate:"required,dive"`
}

type GdsResponse struct {
	GdsNo         string           `json:"gds_no"`
	GdsDate       *string          `json:"gds_date"`
	TrCode        *string          `json:"tr_code"`
	RefNo         *string          `json:"ref_no"`
	WhID          *int64           `json:"wh_id"`
	WhCode        string           `json:"wh_code"`
	WhName        string           `json:"wh_name"`
	SupID         *int64           `json:"sup_id"`
	SupCode       string           `json:"sup_code"`
	SubTotal      *float64         `json:"sub_total"`
	Vat           *float64         `json:"vat"`
	VatValue      *float64         `json:"vat_value"`
	VatLg         *float64         `json:"vat_lg"`
	VatLgValue    *float64         `json:"vat_lg_value"`
	Total         *float64         `json:"total"`
	VatBg         *float64         `json:"vat_bg"`
	VatBgValue    *float64         `json:"vat_bg_value"`
	TotEmbInc     *float64         `json:"tot_emb_inc"`
	TotEmbExc     *float64         `json:"tot_emb_exc"`
	Notes         *string          `json:"notes"`
	DataStatus    *int64           `json:"data_status"`
	UpdatedByName string           `json:"updated_by_name"`
	UpdatedAt     *time.Time       `json:"updated_at"`
	IsClosed      bool             `json:"is_closed"`
	ClosedBy      int64            `json:"closed_by"`
	ClosedByName  string           `json:"closed_by_name"`
	ClosedAt      *time.Time       `json:"closed_at"`
	Details       []GdsDetResponse `json:"details"`
}
type DetailGdsParams struct {
	GdsNo string `params:"gds_no" validate:"required"`
}
type DeleteGdsParams struct {
	GdsNo string `params:"gds_no" validate:"required"`
}

type UpdateGdsParams struct {
	GdsNo string `params:"gds_no" validate:"required"`
}
type GdsListResponse struct {
	GdsNo         string     `json:"gds_no"`
	GdsDate       *string    `json:"gds_date"`
	TrCode        *string    `json:"tr_code"`
	RefNo         *string    `json:"ref_no"`
	WhID          *int64     `json:"wh_id"`
	WhCode        string     `json:"wh_code"`
	WhName        string     `json:"wh_name"`
	SupID         *int64     `json:"sup_id"`
	SupCode       string     `json:"sup_code"`
	SupName       string     `json:"sup_name"`
	SubTotal      *float64   `json:"sub_total"`
	Vat           *float64   `json:"vat"`
	VatValue      *float64   `json:"vat_value"`
	VatLg         *float64   `json:"vat_lg"`
	VatLgValue    *float64   `json:"vat_lg_value"`
	Total         *float64   `json:"total"`
	VatBg         *float64   `json:"vat_bg"`
	VatBgValue    *float64   `json:"vat_bg_value"`
	TotEmbInc     *float64   `json:"tot_emb_inc"`
	TotEmbExc     *float64   `json:"tot_emb_exc"`
	Notes         *string    `json:"notes"`
	DataStatus    *int64     `json:"data_status"`
	UpdatedByName string     `json:"updated_by_name"`
	UpdatedAt     *time.Time `json:"updated_at"`
	IsClosed      bool       `json:"is_closed"`
	ClosedBy      int64      `json:"closed_by"`
	ClosedByName  string     `json:"closed_by_name"`
	ClosedAt      *time.Time `json:"closed_at"`
}

type UpdateGdsBody struct {
	GdsNo      string             `json:"gds_no"`
	CustID     string             `json:"cust_id"`
	GdsDate    *string            `json:"gds_date"`
	TrCode     *string            `json:"tr_code"`
	RefNo      *string            `json:"ref_no"`
	WhID       *int64             `json:"wh_id"`
	SupID      *int64             `json:"sup_id"`
	SubTotal   *float64           `json:"sub_total"`
	Vat        *float64           `json:"vat"`
	VatValue   *float64           `json:"vat_value"`
	VatLg      *float64           `json:"vat_lg"`
	VatLgValue *float64           `json:"vat_lg_value"`
	Total      *float64           `json:"total"`
	VatBg      *float64           `json:"vat_bg"`
	VatBgValue *float64           `json:"vat_bg_value"`
	TotEmbInc  *float64           `json:"tot_emb_inc"`
	TotEmbExc  *float64           `json:"tot_emb_exc"`
	Notes      *string            `json:"notes"`
	DataStatus *int64             `json:"data_status"`
	CreatedBy  *int64             `json:"created_by"`
	UpdatedBy  int64              `json:"updated_by"`
	Details    []UpdateGdsDetBody `json:"details"`
}
