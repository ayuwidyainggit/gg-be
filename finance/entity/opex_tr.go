package entity

type CreateOpexTrBody struct {
	CustID     string                `json:"cust_id"`
	OpexTrNo   string                `json:"opex_tr_no"`
	OpexTrDate *string               `json:"opex_tr_date"`
	TrCode     *string               `json:"tr_code"`
	Notes      *string               `json:"notes"`
	TotAmount  *float64              `json:"tot_amount"`
	DataStatus *int64                `json:"data_status"`
	CreatedBy  *int64                `json:"created_by"`
	Details    []CreateOpexTrDetBody `json:"details"`
}
type OpexTrResponse struct {
	OpexTrNo      string              `json:"opex_tr_no"`
	OpexTrDate    *string             `json:"opex_tr_date"`
	TrCode        *string             `json:"tr_code"`
	Notes         *string             `json:"notes"`
	TotAmount     *float64            `json:"tot_amount"`
	DataStatus    *int64              `json:"data_status"`
	UpdatedByName string              `json:"updated_by_name"`
	UpdatedAt     string              `json:"updated_at"`
	Details       []OpexTrDetResponse `json:"details"`
}
type OpexTrListResponse struct {
	OpexTrNo      string   `json:"opex_tr_no"`
	OpexTrDate    *string  `json:"opex_tr_date"`
	TrCode        *string  `json:"tr_code"`
	Notes         *string  `json:"notes"`
	TotAmount     *float64 `json:"tot_amount"`
	DataStatus    *int64   `json:"data_status"`
	UpdatedByName string   `json:"updated_by_name"`
	UpdatedAt     string   `json:"updated_at"`
}
type UpdateOpexTrBody struct {
	CustID     string                `json:"cust_id"`
	OpexTrNo   string                `json:"opex_tr_no"`
	OpexTrDate *string               `json:"opex_tr_date"`
	TrCode     *string               `json:"tr_code"`
	Notes      *string               `json:"notes"`
	TotAmount  *float64              `json:"tot_amount"`
	DataStatus *int64                `json:"data_status"`
	CreatedBy  *int64                `json:"created_by"`
	UpdatedBy  int64                 `json:"updated_by"`
	Details    []UpdateOpexTrDetBody `json:"details"`
}
type DetailOpextrParams struct {
	OpexTrNo string `params:"opex_tr_no" validate:"required"`
}

type UpdateOpextrParams struct {
	OpexTrNo string `params:"opex_tr_no" validate:"required"`
}
