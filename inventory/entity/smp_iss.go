package entity

type CreateSampleIssueBody struct {
	CustID     string                     `json:"cust_id"`
	SmpIssNo   string                     `json:"smp_iss_no"`
	SmpIssDate string                     `json:"smp_iss_date"`
	TrCode     string                     `json:"tr_code" validate:"required,len=3"`
	WhID       int64                      `json:"wh_id" validate:"required"`
	CndnID     int64                      `json:"cndn_id" validate:"required"`
	OutletID   int64                      `json:"outlet_id" validate:"required"`
	Notes      string                     `json:"notes"`
	SubTotal   float64                    `json:"sub_total" validate:"required"`
	Vat        float64                    `json:"vat"`
	VatValue   float64                    `json:"vat_value"`
	VatLg      float64                    `json:"vat_lg"`
	VatLgValue float64                    `json:"vat_lg_value"`
	Total      float64                    `json:"total" validate:"required"`
	VatBg      float64                    `json:"vat_bg"`
	VatBgValue float64                    `json:"vat_bg_value"`
	DataStatus int64                      `json:"data_status"`
	CreatedBy  int64                      `json:"created_by"`
	Details    []CreateSampleIssueDetBody `json:"details" validate:"required,dive"`
}

type SampleIssueResponse struct {
	SmpIssNo      string               `json:"smp_iss_no"`
	SmpIssDate    string               `json:"smp_iss_date"`
	TrCode        string               `json:"tr_code"`
	WhID          int64                `json:"wh_id"`
	WhCode        string               `json:"wh_code"`
	WhName        string               `json:"wh_name"`
	CndnID        int64                `json:"cndn_id"`
	CndnCode      string               `json:"cndn_code"`
	CndnName      string               `json:"cndn_name"`
	OutletID      int64                `json:"outlet_id"`
	Notes         string               `json:"notes"`
	SubTotal      float64              `json:"sub_total"`
	Vat           float64              `json:"vat"`
	VatValue      float64              `json:"vat_value"`
	VatLg         float64              `json:"vat_lg"`
	VatLgValue    float64              `json:"vat_lg_value"`
	Total         float64              `json:"total"`
	VatBg         float64              `json:"vat_bg"`
	VatBgValue    float64              `json:"vat_bg_value"`
	DataStatus    int64                `json:"data_status"`
	UpdatedAt     string               `json:"updated_at"`
	UpdatedByName string               `json:"updated_by_name"`
	IsClosed      bool                 `json:"is_closed"`
	ClosedBy      int64                `json:"closed_by"`
	ClosedByName  string               `json:"closed_by_name"`
	ClosedAt      string               `json:"closed_at"`
	Details       []SampleIssueDetResp `json:"details"`
}
type SampleIssueListResponse struct {
	SmpIssNo      string  `json:"smp_iss_no"`
	SmpIssDate    string  `json:"smp_iss_date"`
	TrCode        string  `json:"tr_code"`
	WhID          int64   `json:"wh_id"`
	WhCode        string  `json:"wh_code"`
	WhName        string  `json:"wh_name"`
	CndnID        int64   `json:"cndn_id"`
	CndnCode      string  `json:"cndn_code"`
	CndnName      string  `json:"cndn_name"`
	OutletID      int64   `json:"outlet_id"`
	Notes         string  `json:"notes"`
	SubTotal      float64 `json:"sub_total"`
	Vat           float64 `json:"vat"`
	VatValue      float64 `json:"vat_value"`
	VatLg         float64 `json:"vat_lg"`
	VatLgValue    float64 `json:"vat_lg_value"`
	Total         float64 `json:"total"`
	VatBg         float64 `json:"vat_bg"`
	VatBgValue    float64 `json:"vat_bg_value"`
	DataStatus    int64   `json:"data_status"`
	UpdatedAt     string  `json:"updated_at"`
	UpdatedByName string  `json:"updated_by_name"`
	IsClosed      bool    `json:"is_closed"`
	ClosedBy      int64   `json:"closed_by"`
	ClosedByName  string  `json:"closed_by_name"`
	ClosedAt      string  `json:"closed_at"`
}
type DetailSmpIssueParams struct {
	SmpIssNo string `params:"smp_iss_no" validate:"required"`
}
type UpdateSmpIssueParams struct {
	SmpIssNo string `params:"smp_iss_no" validate:"required"`
}
type UpdateSampleIssueBody struct {
	CustID       string                     `json:"cust_id"`
	ParentCustID string                     `json:"parent_cust_id"`
	SmpIssNo     string                     `json:"smp_iss_no"`
	SmpIssDate   *string                    `json:"smp_iss_date"`
	TrCode       *string                    `json:"tr_code"`
	WhID         *int64                     `json:"wh_id"`
	CndnID       *int64                     `json:"cndn_id"`
	OutletID     *int64                     `json:"outlet_id"`
	Notes        *string                    `json:"notes"`
	SubTotal     *float64                   `json:"sub_total"`
	Vat          *float64                   `json:"vat"`
	VatValue     *float64                   `json:"vat_value"`
	VatLg        *float64                   `json:"vat_lg"`
	VatLgValue   *float64                   `json:"vat_lg_value"`
	Total        *float64                   `json:"total"`
	VatBg        *float64                   `json:"vat_bg"`
	VatBgValue   *float64                   `json:"vat_bg_value"`
	DataStatus   *int64                     `json:"data_status"`
	CreatedBy    *int64                     `json:"created_by"`
	UpdatedBy    int64                      `json:"updated_by"`
	Details      []UpdateSampleIssueDetBody `json:"details"`
}
