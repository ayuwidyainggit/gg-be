package entity

type CreateChequeBody struct {
	CustID       string   `json:"cust_id"`
	ChqTrNo      string   `json:"chq_tr_no"`
	TrCode       *string  `json:"tr_code"`
	ChqTrType    *int     `json:"chq_tr_type"`
	ChqDate      *string  `json:"chq_date"`
	ChqDueDate   *string  `json:"chq_due_date"`
	BankId       *int64   `json:"bank_id"`
	AccountNo    *string  `json:"account_no"`
	ChqAmt       *float64 `json:"chq_amt"`
	ChqUsedAmt   *float64 `json:"chq_used_amt"`
	SalesmanId   *int64   `json:"salesman_id"`
	OutletId     *int64   `json:"outlet_id"`
	Notes        *string  `json:"notes"`
	ClearingDate *string  `json:"clearing_date"`
	ChqStatus    *int     `json:"chq_status"`
	StatusDate   *string  `json:"status_date"`
	IsPosted     bool     `json:"is_posted"`
	CreatedBy    *int64   `json:"created_by"`
}

type ChequeResponse struct {
	ChqTrNo       string   `json:"chq_tr_no"`
	TrCode        *string  `json:"tr_code"`
	ChqNo         *int     `json:"chq_no"`
	ChqTrType     *int     `json:"chq_tr_type"`
	ChqDate       *string  `json:"chq_date"`
	ChqDueDate    *string  `json:"chq_due_date"`
	BankId        *int64   `json:"bank_id"`
	BankCode      string   `json:"bank_code"`
	Bankname      string   `json:"bank_name"`
	AccountNo     *string  `json:"account_no"`
	ChqAmt        *float64 `json:"chq_amt"`
	ChqUsedAmt    *float64 `json:"chq_used_amt"`
	SalesmanId    *int64   `json:"salesman_id"`
	SalesmanCode  string   `json:"salesman_code"`
	SalesmanName  string   `json:"salesman_name"`
	OutletId      *int64   `json:"outlet_id"`
	OutletCode    string   `json:"outlet_code"`
	OutletName    string   `json:"outlet_name"`
	Notes         *string  `json:"notes"`
	ClearingDate  *string  `json:"clearing_date"`
	ChqStatus     *int     `json:"chq_status"`
	StatusDate    *string  `json:"status_date"`
	UpdatedByName string   `json:"updated_by_name"`
	UpdatedAt     string   `json:"updated_at"`
	IsPosted      bool     `json:"is_posted"`
	PostedAt      *string  `json:"posted_at"`
}

type ChequeListResponse struct {
	ChqTrNo       string   `json:"chq_tr_no"`
	TrCode        *string  `json:"tr_code"`
	ChqNo         int      `json:"chq_no"`
	ChqTrType     *int     `json:"chq_tr_type"`
	ChqDate       *string  `json:"chq_date"`
	ChqDueDate    *string  `json:"chq_due_date"`
	BankId        *int64   `json:"bank_id"`
	BankCode      string   `json:"bank_code"`
	Bankname      string   `json:"bank_name"`
	AccountNo     *string  `json:"account_no"`
	ChqAmt        *float64 `json:"chq_amt"`
	ChqUsedAmt    *float64 `json:"chq_used_amt"`
	SalesmanId    *int64   `json:"salesman_id"`
	SalesmanCode  string   `json:"salesman_code"`
	SalesmanName  string   `json:"salesman_name"`
	OutletId      *int64   `json:"outlet_id"`
	OutletCode    string   `json:"outlet_code"`
	OutletName    string   `json:"outlet_name"`
	Notes         *string  `json:"notes"`
	ClearingDate  *string  `json:"clearing_date"`
	ChqStatus     *int     `json:"chq_status"`
	StatusDate    *string  `json:"status_date"`
	UpdatedAt     string   `json:"updated_at"`
	UpdatedByName string   `json:"updated_by_name"`
	IsPosted      *bool    `json:"is_posted"`
}

type UpdateChequeBody struct {
	CustID       string   `json:"cust_id"`
	ChqTrNo      string   `json:"chq_tr_no"`
	TrCode       *string  `json:"tr_code"`
	ChqNo        *int     `json:"chq_no"`
	ChqTrType    *int     `json:"chq_tr_type"`
	ChqDate      *string  `json:"chq_date"`
	ChqDueDate   *string  `json:"chq_due_date"`
	BankId       *int64   `json:"bank_id"`
	AccountNo    *string  `json:"account_no"`
	ChqAmt       *float64 `json:"chq_amt"`
	ChqUsedAmt   *float64 `json:"chq_used_amt"`
	SalesmanId   *int64   `json:"salesman_id"`
	OutletId     *int64   `json:"outlet_id"`
	Notes        *string  `json:"notes"`
	ClearingDate *string  `json:"clearing_date"`
	ChqStatus    *int     `json:"chq_status"`
	StatusDate   *string  `json:"status_date"`
	UpdatedBy    int64    `json:"updated_by"`
	IsPosted     *bool    `json:"is_posted"`
}

type DetailChequeParams struct {
	ChqNo int `params:"chq_no" validate:"required"`
}
type DeleteChequeParams struct {
	ChqNo int `params:"chq_no" validate:"required"`
}
type UpdateChequeParams struct {
	ChqNo int `params:"chq_no" validate:"required"`
}
