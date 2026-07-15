package entity

import "time"

type CheckGiroQueryFilter struct {
	GeneralQueryFilter
	BankID    []int    `query:"bank_id"`
	AccountNo []string `query:"account_no"`
}

type CreateChequeGiroBody struct {
	CustID           string  `json:"cust_id"`
	DocNoCheque      string  `json:"doc_no_cheque"`
	OwnerID          int     `json:"owner_id"`
	SalesmanID       *int    `json:"salesman_id"`
	SupplierID       *int    `json:"sup_id"`
	OutletID         *int    `json:"outlet_id"`
	BankID           int     `json:"bank_id"`
	BankIDCollecting int     `json:"bank_id_collecting"`
	AccountNo        string  `json:"account_no"`
	DocDateCheque    *string `json:"doc_date_cheque"`
	DueDate          *string `json:"due_date"`
	Amount           float64 `json:"amount"`
	StatusCheque     int     `json:"status_cheque"`
	CreatedBy        *int64  `json:"created_by"`
}

type ChequeGiroResponse struct {
	CustID           string     `json:"cust_id"`
	ChequeGiroNo     int        `json:"cheque_giro_no"`
	DocNoCheque      string     `json:"doc_no_cheque"`
	OwnerID          int        `json:"owner_id"`
	OwnerName        string     `json:"owner_name"`
	SupplierID       *int       `json:"sup_id"`
	SupplierName     *int       `json:"sup_name"`
	SalesmanID       *int       `json:"salesman_id"`
	SalesmanName     *string    `json:"sales_name"`
	OutletID         *int       `json:"outlet_id"`
	OutletName       *string    `json:"outlet_name"`
	BankID           int        `json:"bank_id"`
	BankName         string     `json:"bank_name"`
	BankIDCollecting int        `json:"bank_id_collecting"`
	AccountNo        string     `json:"account_no"`
	DocDateCheque    *string    `json:"doc_date_cheque"`
	DueDate          *string    `json:"due_date"`
	Amount           float64    `json:"amount"`
	UsedAmount       float64    `json:"used_amount"`
	RemainingAmount  float64    `json:"remaining_amount"`
	StatusCheque     int        `json:"status_cheque"`
	StatusChequeText *string    `json:"status_cheque_text"`
	ClearingDate     *time.Time `json:"clearing_date"`
	CreatedBy        int64      `json:"created_by"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedBy        int64      `json:"updated_by"`
	UpdatedAt        time.Time  `json:"updated_at"`
}

type UpdateChequeGiroBody struct {
	CustID           string  `json:"cust_id"`
	DocNoCheque      string  `json:"doc_no_cheque"`
	OwnerID          int     `json:"owner_id"`
	SalesmanID       *int    `json:"salesman_id"`
	SupplierID       *int    `json:"sup_id"`
	OutletID         *int    `json:"outlet_id"`
	BankID           int     `json:"bank_id"`
	BankIDCollecting int     `json:"bank_id_collecting"`
	AccountNo        string  `json:"account_no"`
	DocDateCheque    *string `json:"doc_date_cheque"`
	DueDate          *string `json:"due_date"`
	Amount           float64 `json:"amount"`
	StatusCheque     int     `json:"status_cheque"`
	CreatedBy        *int64  `json:"created_by"`
	UpdatedBy        int64   `json:"updated_by"`
}

type DetailChequeGiroParams struct {
	ChequeGiroNo int `params:"cheque_giro_no" validate:"required"`
}
type DeleteChequeGiroParams struct {
	ChequeGiroNo int `params:"cheque_giro_no" validate:"required"`
}
type UpdateChequeGiroParams struct {
	ChequeGiroNo int `params:"cheque_giro_no" validate:"required"`
}

type BankLookup struct {
	BankId   int    `json:"bank_id"`
	BankCode string `json:"bank_code"`
	BankName string `json:"bank_name"`
}

type BankAccountLookup struct {
	AccountNo string `json:"account_no"`
}

var StatusGiro = map[int]string{
	1: "Rejected",
	2: "Pending",
	3: "Accepted",
}

var OwnerGiro = map[int]string{
	1: "Outlet",
	2: "Distributor",
}
