package entity

type CreateCashTrBody struct {
	CustId      string                `json:"cust_id"`
	CashTrNo    string                `json:"cash_tr_no"`
	TrCode      *string               `json:"tr_code"`
	CashTrDate  *string               `json:"cash_tr_date"`
	CoaIdTo     *int                  `json:"coa_id_to"`
	Notes       *string               `json:"notes"`
	AccountNo   *string               `json:"account_no"`
	AccountName *string               `json:"account_name"`
	Amount      *float64              `json:"amount"`
	CreatedBy   *int64                `json:"created_by"`
	IsDel       bool                  `json:"is_del"`
	Details     []CreateCashTrDetBody `json:"details"`
}

type CashTrResponse struct {
	CashTrNo      string              `json:"cash_tr_no"`
	TrCode        *string             `json:"tr_code"`
	CashTrDate    *string             `json:"cash_tr_date"`
	CoaIdTo       *int                `json:"coa_id_to"`
	CoaCodeTo     string              `json:"coa_code_to"`
	CoaNameTo     string              `json:"coa_name_to"`
	Notes         *string             `json:"notes"`
	AccountNo     *string             `json:"account_no"`
	AccountName   *string             `json:"account_name"`
	Amount        *float64            `json:"amount"`
	DataStatus    *int64              `json:"data_status"`
	UpdatedByName string              `json:"updated_by_name"`
	UpdatedAt     string              `json:"updated_at"`
	Details       []CashTrDetResponse `json:"details"`
}

type CashTrListResponse struct {
	CashTrNo      string   `json:"cash_tr_no"`
	TrCode        *string  `json:"tr_code"`
	CashTrDate    *string  `json:"cash_tr_date"`
	CoaIdTo       *int64   `json:"coa_id_to"`
	CoaCodeTo     *string  `json:"coa_code_to"`
	CoaNameTo     *string  `json:"coa_name_to"`
	Notes         *string  `json:"notes"`
	AccountNo     *string  `json:"account_no"`
	AccountName   *string  `json:"account_name"`
	Amount        *float64 `json:"amount"`
	UpdatedAt     string   `json:"updated_at"`
	UpdatedByName string   `json:"updated_by_name"`
}

type UpdateCashTrBody struct {
	CustId      string                `json:"cust_id"`
	CashTrNo    string                `json:"cash_tr_no"`
	TrCode      *string               `json:"tr_code"`
	CashTrDate  *string               `json:"cash_tr_date"`
	CoaIdTo     *int                  `json:"coa_id_to"`
	Notes       *string               `json:"notes"`
	AccountNo   *string               `json:"account_no"`
	AccountName *string               `json:"account_name"`
	Amount      *float64              `json:"amount"`
	CreatedBy   *int64                `json:"created_by"`
	UpdatedBy   int64                 `json:"updated_by"`
	Details     []UpdateCashTrDetBody `json:"details"`
}

type DetailCashTrParams struct {
	CashTrNo string `params:"cash_tr_no" validate:"required"`
}

type UpdateCashTrParams struct {
	CashTrNo string `params:"cash_tr_no" validate:"required"`
}
