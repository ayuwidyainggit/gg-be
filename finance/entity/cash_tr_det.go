package entity

type CreateCashTrDetBody struct {
	CashTrDetId int      `json:"cash_tr_det_id"`
	CoaId       *int64   `json:"coa_id"`
	Amount      *float64 `json:"amount"`
}

type CashTrDetResponse struct {
	CashTrDetId int      `json:"cash_tr_det_id"`
	CoaId       *int64   `json:"coa_id"`
	CoaCode     string   `json:"coa_code"`
	CoaName     string   `json:"coa_name"`
	Amount      *float64 `json:"amount"`
}

type UpdateCashTrDetBody struct {
	CashTrDetId *int64   `json:"cash_tr_det_id"`
	CoaId       *int64   `json:"coa_id"`
	Amount      *float64 `json:"amount"`
}
