package entity

/*
type PurchaseReturn struct {
	CustID         string    `json:"cust_id"`
	BpprNo         string    `json:"bppr_no"`
	BpprDate       string    `json:"bppr_date"`
	TrCode         string    `json:"tr_code"`
	SupID          int64     `json:"sup_id"`
	WhID           int64     `json:"wh_id"`
	ItemCdn        int64     `json:"item_cdn"`
	ReturnReasonID int64     `json:"return_reason_id"`
	ReturnNo       string    `json:"return_no"`
	ReturnDate     string    `json:"return_date"`
	Notes          string    `json:"notes"`
	TotEmbInc      float64   `json:"tot_emb_inc"`
	TotEmbExc      float64   `json:"tot_emb_exc"`
	SubTotal       float64   `json:"sub_total"`
	Vat            float64   `json:"vat"`
	VatValue       float64   `json:"vat_value"`
	VatLg          float64   `json:"vat_lg"`
	VatLgValue     float64   `json:"vat_lg_value"`
	Total          float64   `json:"total"`
	VatBg          float64   `json:"vat_bg"`
	VatBgValue     float64   `json:"vat_bg_value"`
	DataStatus     int64     `json:"data_status"`
	CreatedBy      int64     `json:"created_by"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedBy      int64     `json:"updated_by"`
	UpdatedAt      time.Time `json:"updated_at"`
	IsDel          bool      `json:"is_del"`
	DeletedBy      int64     `json:"deleted_by"`
	DeletedAt      time.Time `json:"deleted_at"`
}

type BpprResponse struct {
	BpprNo         string     `json:"bppr_no"`
	BpprDate       string     `json:"bppr_date"`
	TrCode         string     `json:"tr_code"`
	SupID          int64      `json:"sup_id"`
	SupCode        string     `json:"sup_code"`
	SupName        string     `json:"sup_name"`
	WhID           int64      `json:"wh_id"`
	WhCode         string     `json:"wh_code"`
	WhName         string     `json:"wh_name"`
	ItemCdn        int64      `json:"item_cdn"`
	ReturnReasonID int64      `json:"return_reason_id"`
	ReturnNo       string     `json:"return_no"`
	ReturnDate     string     `json:"return_date"`
	Notes          string     `json:"notes"`
	TotEmbInc      float64    `json:"tot_emb_inc"`
	TotEmbExc      float64    `json:"tot_emb_exc"`
	SubTotal       float64    `json:"sub_total"`
	Vat            float64    `json:"vat"`
	VatValue       float64    `json:"vat_value"`
	VatLg          float64    `json:"vat_lg"`
	VatLgValue     float64    `json:"vat_lg_value"`
	Total          float64    `json:"total"`
	VatBg          float64    `json:"vat_bg"`
	VatBgValue     float64    `json:"vat_bg_value"`
	DataStatus     int64      `json:"data_status"`
	UpdatedAt      *time.Time `json:"updated_at"`
	UpdatedByName  string     `json:"updated_by_name"`
	Details        []BpprDet  `json:"details"`
}

type DetailBpprParams struct {
	BpprNo string `params:"bppr_no" validate:"required"`
}
type UpdateBpprParams struct {
	BpprNo string `params:"bppr_no" validate:"required"`
}
type DeleteBpprParams struct {
	BpprNo string `params:"bppr_no" validate:"required"`
}

type CreateBpprBody struct {
	CustID         string    `json:"cust_id"`
	BpprNo         string    `json:"bppr_no"`
	BpprDate       *string   `json:"bppr_date" validate:"required"`
	TrCode         string    `json:"tr_code" validate:"required,oneof='BPR'"`
	SupID          int64     `json:"sup_id" validate:"required"`
	WhID           int64     `json:"wh_id" validate:"required"`
	ItemCdn        int64     `json:"item_cdn" validate:"required,oneof='1' '2' '3'"`
	ReturnReasonID int64     `json:"return_reason_id"`
	ReturnNo       string    `json:"return_no"`
	ReturnDate     *string   `json:"return_date"`
	Notes          string    `json:"notes"`
	TotEmbInc      float64   `json:"tot_emb_inc"`
	TotEmbExc      float64   `json:"tot_emb_exc"`
	SubTotal       float64   `json:"sub_total"`
	Vat            float64   `json:"vat"`
	VatValue       float64   `json:"vat_value"`
	VatLg          float64   `json:"vat_lg"`
	VatLgValue     float64   `json:"vat_lg_value"`
	Total          float64   `json:"total"`
	VatBg          float64   `json:"vat_bg"`
	VatBgValue     float64   `json:"vat_bg_value"`
	DataStatus     int64     `json:"data_status"`
	CreatedBy      int64     `json:"created_by"`
	UpdatedBy      int64     `json:"updated_by"`
	Details        []BpprDet `json:"details"`
}

type UpdateBpprRequest struct {
	CustID         string                 `json:"cust_id"`
	BpprNo         *string                `json:"bppr_no"`
	BpprDate       *string                `json:"bppr_date"`
	TrCode         *string                `json:"tr_code"`
	SupID          *int64                 `json:"sup_id"`
	WhID           *int64                 `json:"wh_id"`
	ItemCdn        *int64                 `json:"item_cdn"`
	ReturnReasonID *int64                 `json:"return_reason_id"`
	ReturnNo       *string                `json:"return_no"`
	ReturnDate     *string                `json:"return_date"`
	Notes          *string                `json:"notes"`
	TotEmbInc      *float64               `json:"tot_emb_inc"`
	TotEmbExc      *float64               `json:"tot_emb_exc"`
	SubTotal       *float64               `json:"sub_total"`
	Vat            *float64               `json:"vat"`
	VatValue       *float64               `json:"vat_value"`
	VatLg          *float64               `json:"vat_lg"`
	VatLgValue     *float64               `json:"vat_lg_value"`
	Total          *float64               `json:"total"`
	VatBg          *float64               `json:"vat_bg"`
	VatBgValue     *float64               `json:"vat_bg_value"`
	DataStatus     *int64                 `json:"data_status"`
	UpdatedBy      int64                  `json:"updated_by"`
	Details        []BpprDetUpdateRequest `json:"details"`
}
*/
