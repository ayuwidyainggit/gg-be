package entity

type CreateOpexTrDetBody struct {
	OpexTrDetID int      `json:"opex_tr_det_id"`
	CcID        *int64   `json:"cc_id"`
	OpexID      *int64   `json:"opex_id"`
	TrDesc      *string  `json:"tr_desc"`
	Amount      *float64 `json:"amount"`
}
type OpexTrDetResponse struct {
	OpexTrDetID int      `json:"opex_tr_det_id"`
	CcID        *int64   `json:"cc_id"`
	OpexID      *int64   `json:"opex_id"`
	OpexCode    string   `json:"opex_code"`
	OpexName    string   `json:"opex_name"`
	TrDesc      *string  `json:"tr_desc"`
	Amount      *float64 `json:"amount"`
}
type UpdateOpexTrDetBody struct {
	OpexTrDetID *int64   `json:"opex_tr_det_id"`
	CcID        *int64   `json:"cc_id"`
	OpexID      *int64   `json:"opex_id"`
	TrDesc      *string  `json:"tr_desc"`
	Amount      *float64 `json:"amount"`
}
