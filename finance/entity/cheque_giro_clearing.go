package entity

import "time"

type CheckGiroClearingQueryFilter struct {
	GeneralQueryFilter
	BankID   []int `query:"bank_id"`
	StatusID []int `query:"status_id"`
}

type ChequeGiroClearingResponse struct {
	CustID               string     `json:"cust_id"`
	ChequeGiroNo         int        `json:"cheque_giro_no"`
	DocNoCheque          string     `json:"doc_no_cheque"`
	OwnerID              int        `json:"owner_id"`
	OwnerName            string     `json:"owner_name"`
	SupplierID           *int       `json:"sup_id"`
	SupplierName         *int       `json:"sup_name"`
	SalesmanID           *int       `json:"salesman_id"`
	SalesmanName         *string    `json:"sales_name"`
	OutletID             *int       `json:"outlet_id"`
	OutletCode           *string    `json:"outlet_code"`
	OutletName           *string    `json:"outlet_name"`
	BankID               int        `json:"bank_id"`
	BankName             string     `json:"bank_name"`
	PayingBankName       string     `json:"paying_bank_name"`
	BankIDCollecting     int        `json:"bank_id_collecting"`
	BankIDCollectingName *string    `json:"bank_id_collecting_name"`
	AccountNo            string     `json:"account_no"`
	DocDateCheque        *string    `json:"doc_date_cheque"`
	DueDate              *string    `json:"due_date"`
	Amount               float64    `json:"amount"`
	UsedAmount           float64    `json:"used_amount"`
	RemainingAmount      float64    `json:"remaining_amount"`
	StatusCheque         int        `json:"status_cheque"`
	StatusClearing       *string    `json:"status_cheque_text"`
	ClearingDate         *time.Time `json:"clearing_date"`
	Reason               string     `json:"reason"`
	CashNo               string     `json:"cash_no"`
	CashAmount           float64    `json:"cash_amount"`
	TransferNo           string     `json:"transfer_no"`
	TransferAmount       float64    `json:"transfer_amount"`
	ChequeNo             string     `json:"cheque_no"`
	ChequeAmount         float64    `json:"cheque_amount"`
	CreatedBy            int64      `json:"created_by"`
	CreatedAt            time.Time  `json:"created_at"`
	UpdatedBy            int64      `json:"updated_by"`
	UpdatedAt            time.Time  `json:"updated_at"`
}

type UpdateChequeGiroClearingBody struct {
	CustID           string `json:"cust_id"`
	BankIDCollecting int    `json:"bank_id_collecting"`
	ClearingDate     string `json:"clearing_date"`
	StatusCheque     int    `json:"status_cheque"`
	UpdatedBy        int64  `json:"updated_by"`
}

type UpdateChequeGiroClearingChangeBody struct {
	CustID            string  `json:"cust_id"`
	DocNoCheque       string  `json:"doc_no_cheque"`
	BankIDCollecting  int     `json:"bank_id_collecting"`
	ClearingDate      string  `json:"clearing_date"`
	StatusCheque      int     `json:"status_cheque"`
	Reason            string  `json:"reason"`
	CashNo            string  `json:"cash_no"`
	CashAmount        float64 `json:"cash_amount"`
	TransferNo        string  `json:"transfer_no"`
	TransferBalance   float64 `json:"transfer_balance"`
	TransferAmount    float64 `json:"transfer_amount"`
	ChequeNo          string  `json:"cheque_no"`
	ChequeBalance     float64 `json:"cheque_balance"`
	ChequeAmount      float64 `json:"cheque_amount"`
	CashRemaining     float64
	ChequeRemaining   float64
	TransferRemaining float64
	UpdatedBy         int64 `json:"updated_by"`
}

var FilterStatusGiroClearing = map[int]string{
	0: "All Status",
	1: "Rejected",
	2: "Pending",
	3: "Accepted",
}
