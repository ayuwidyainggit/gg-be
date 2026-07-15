package entity

type CreateArPayBody struct {
	CustID        string               `json:"cust_id"`
	ArPayNo       string               `json:"ar_pay_no"`
	ArPayDate     *string              `json:"ar_pay_date"`
	TrCode        *string              `json:"tr_code"`
	ArNo          *string              `json:"ar_no"`
	SalesmanID    *int64               `json:"salesman_id"`
	CashAmt       *float64             `json:"cash_amt"`
	ChequeAmt     *float64             `json:"cheque_amt"`
	TransferAmt   *float64             `json:"transfer_amt"`
	ReturnAmt     *float64             `json:"return_amt"`
	CndnAmt       *float64             `json:"cndn_amt"`
	DiscAmt       *float64             `json:"disc_amt"`
	DutyStampAmt  *float64             `json:"duty_stamp_amt"`
	TotalAmt      *float64             `json:"total_amt"`
	TotalDiff     *float64             `json:"total_diff"`
	TotalAmtRound *float64             `json:"total_amt_round"`
	DataStatus    *int64               `json:"data_status"`
	IsPosted      bool                 `json:"is_posted"`
	CreatedBy     *int64               `json:"created_by"`
	Details       []CreateArPayDetBody `json:"details"`
}
type ArPayResponse struct {
	ArPayNo       string             `json:"ar_pay_no"`
	ArPayDate     *string            `json:"ar_pay_date"`
	TrCode        *string            `json:"tr_code"`
	ArNo          *string            `json:"ar_no"`
	SalesmanID    *int64             `json:"salesman_id"`
	SalesmanCode  string             `json:"salesman_code"`
	SalesmanName  string             `json:"salesman_name"`
	CashAmt       *float64           `json:"cash_amt"`
	ChequeAmt     *float64           `json:"cheque_amt"`
	TransferAmt   *float64           `json:"transfer_amt"`
	ReturnAmt     *float64           `json:"return_amt"`
	CndnAmt       *float64           `json:"cndn_amt"`
	DiscAmt       *float64           `json:"disc_amt"`
	DutyStampAmt  *float64           `json:"duty_stamp_amt"`
	TotalAmt      *float64           `json:"total_amt"`
	TotalDiff     *float64           `json:"total_diff"`
	TotalAmtRound *float64           `json:"total_amt_round"`
	DataStatus    *int64             `json:"data_status"`
	UpdatedByName string             `json:"updated_by_name"`
	UpdatedAt     string             `json:"updated_at"`
	IsPosted      bool               `json:"is_posted"`
	PostedAt      *string            `json:"posted_at"`
	Details       []ArPayDetResponse `json:"details"`
}

type ArPayListResponse struct {
	ArPayNo       string   `json:"ar_pay_no"`
	ArPayDate     *string  `json:"ar_pay_date"`
	TrCode        *string  `json:"tr_code"`
	ArNo          *string  `json:"ar_no"`
	SalesmanID    *int64   `json:"salesman_id"`
	SalesmanCode  string   `json:"salesman_code"`
	SalesmanName  string   `json:"salesman_name"`
	CashAmt       *float64 `json:"cash_amt"`
	ChequeAmt     *float64 `json:"cheque_amt"`
	TransferAmt   *float64 `json:"transfer_amt"`
	ReturnAmt     *float64 `json:"return_amt"`
	CndnAmt       *float64 `json:"cndn_amt"`
	DiscAmt       *float64 `json:"disc_amt"`
	DutyStampAmt  *float64 `json:"duty_stamp_amt"`
	TotalAmt      *float64 `json:"total_amt"`
	TotalDiff     *float64 `json:"total_diff"`
	TotalAmtRound *float64 `json:"total_amt_round"`
	DataStatus    *int64   `json:"data_status"`
	UpdatedByName string   `json:"updated_by_name"`
	UpdatedAt     string   `json:"updated_at"`
	IsPosted      *bool    `json:"is_posted"`
}
type UpdateArPayBody struct {
	CustID        string               `json:"cust_id"`
	ArPayNo       string               `json:"ar_pay_no"`
	ArPayDate     *string              `json:"ar_pay_date"`
	TrCode        *string              `json:"tr_code"`
	ArNo          *string              `json:"ar_no"`
	SalesmanID    *int64               `json:"salesman_id"`
	CashAmt       *float64             `json:"cash_amt"`
	ChequeAmt     *float64             `json:"cheque_amt"`
	TransferAmt   *float64             `json:"transfer_amt"`
	ReturnAmt     *float64             `json:"return_amt"`
	CndnAmt       *float64             `json:"cndn_amt"`
	DiscAmt       *float64             `json:"disc_amt"`
	DutyStampAmt  *float64             `json:"duty_stamp_amt"`
	TotalAmt      *float64             `json:"total_amt"`
	TotalDiff     *float64             `json:"total_diff"`
	TotalAmtRound *float64             `json:"total_amt_round"`
	DataStatus    *int64               `json:"data_status"`
	UpdatedBy     int64                `json:"updated_by"`
	IsPosted      *bool                `json:"is_posted"`
	Details       []UpdateArPayDetBody `json:"details"`
}

type DetailArPayParams struct {
	ArPayNo string `params:"ar_pay_no" validate:"required"`
}
type DeleteArPayParams struct {
	ArPayNo string `params:"ar_pay_no" validate:"required"`
}
type UpdateArPayParams struct {
	ArPayNo string `params:"ar_pay_no" validate:"required"`
}
