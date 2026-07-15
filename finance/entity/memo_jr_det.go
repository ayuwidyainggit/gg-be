package entity

type CreateMemoJrDetBody struct {
	CcID   int      `json:"cc_id"`
	CoaID  int      `json:"coa_id"`
	Debit  *float64 `json:"debit"`
	Credit *float64 `json:"credit"`
	Notes  *string  `json:"notes"`
}

type MemoJrDetResponse struct {
	MemoJrDetID int64    `json:"memo_jr_det_id"`
	CcID        int      `json:"cc_id"`
	CoaID       int      `json:"coa_id"`
	CoaCode     string   `json:"coa_code"`
	CoaName     string   `json:"coa_name"`
	Debit       *float64 `json:"debit"`
	Credit      *float64 `json:"credit"`
	Notes       *string  `json:"notes"`
}
type UpdateMemoJrDetBody struct {
	MemoJrDetID *int64   `json:"memo_jr_det_id"`
	CcID        int      `json:"cc_id"`
	CoaID       int      `json:"coa_id"`
	Debit       *float64 `json:"debit"`
	Credit      *float64 `json:"credit"`
	Notes       *string  `json:"notes"`
}
