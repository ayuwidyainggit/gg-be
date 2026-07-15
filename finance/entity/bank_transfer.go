package entity

import "time"

type BankTransferQueryFilter struct {
	GeneralQueryFilter
	BankID    []int    `query:"bank_id"`
	AccountNo []string `query:"account_no"`
}

type CreateBankTransferBody struct {
	OwnerID        int                    `json:"owner_id"`
	SalesmanID     *int                   `json:"salesman_id"`
	Supplier       SupplierDetails        `json:"supplier"`
	Outlet         OutletDetails          `json:"outlet"`
	Bank           BankDetails            `json:"bank"`
	Document       DocumentDetails        `json:"document"`
	Amount         float64                `json:"amount"`
	ProofOfPayment []BankTransferFileBody `json:"proof_of_payment"`
	CustID         string                 `json:"cust_id"`              // Populated from context
	CreatedBy      *int64                 `json:"created_by"`           // Populated from context
	StatusBank     int                    `json:"status_bank_transfer"` // Defaulted in service
}

type SupplierDetails struct {
	SupCode string `json:"sup_code"`
	SupName string `json:"sup_name"`
	SupID   *int   `json:"sup_id,omitempty"` // Added for internal mapping if needed
}

type OutletDetails struct {
	OutletID     int `json:"outlet_id"`
	OutletBankID int `json:"outlet_bank_id"`
}

type BankDetails struct {
	BankID      int    `json:"bank_id"`
	AccountNo   string `json:"account_no"`
	AccountName string `json:"account_name"`
}

type DocumentDetails struct {
	TransferDate string `json:"transfer_date"`
}

type BankTransferFileBody struct {
	FileName      string `json:"file_name"`
	FileURL       string `json:"file_url"`
	FileKey       string `json:"file_key"`
	FileSize      int64  `json:"file_size"`
	MediaCategory string `json:"media_category"`
}

// BankTransferDepositDataItem represents one row of deposit usage for a bank transfer (used_amount query)
type BankTransferDepositDataItem struct {
	DepositNo   string  `json:"deposit_no"`
	DepositDate string  `json:"deposit_date"`
	InvoiceNo   string  `json:"invoice_no"`
	UsedAmount  float64 `json:"used_amount"`
}

type BankTransferResponse struct {
	CustID           string                         `json:"cust_id"`
	BankTransferNo   int                            `json:"bank_transfer_no"`
	DocNoBank        string                         `json:"doc_no_bank"`
	OwnerID          int                            `json:"owner_id"`
	OwnerName        string                         `json:"owner_name"`
	SupplierID       *int                           `json:"sup_id"`
	SupplierName     *string                        `json:"sup_name"`
	SupplierCode     *string                        `json:"sup_code"`
	SalesmanID       *int                           `json:"salesman_id"`
	SalesmanName     *string                        `json:"sales_name"`
	SalesmanCode     *string                        `json:"salesman_code"`
	OutletID         *int                           `json:"outlet_id"`
	OutletName       *string                        `json:"outlet_name"`
	OutletCode       *string                        `json:"outlet_code"`
	BankID           int                            `json:"bank_id"`
	BankName         string                         `json:"bank_name"`
	BankIDCollecting int                            `json:"bank_id_collecting"`
	AccountNo        *string                        `json:"account_no"`
	AccountName      string                         `json:"account_name"`
	TransferDate     *string                        `json:"transfer_date"`
	Amount           float64                        `json:"amount"`
	UsedAmount       float64                        `json:"used_amount"`
	RemainingAmount  float64                        `json:"remaining_amount"`
	StatusBank       int                            `json:"status_bank_transfer"`
	StatusBankText   *string                        `json:"status_bank_transfer_text"`
	CreatedBy        int64                          `json:"created_by"`
	CreatedAt        time.Time                      `json:"created_at"`
	UpdatedBy        int64                          `json:"updated_by"`
	UpdatedAt        time.Time                      `json:"updated_at"`
	DepositData      []BankTransferDepositDataItem  `json:"deposit_data"`
	ProofOfPayment   []BankTransferFileBody         `json:"proof_of_payment"`
}

type UpdateBankTransferBody struct {
	CustID           string               `json:"cust_id"`
	DocNoBank        string               `json:"doc_no_bank"`
	OwnerID          int                  `json:"owner_id"`
	SalesmanID       *int                 `json:"salesman_id"`
	SupplierID       *int                 `json:"sup_id"`
	OutletID         *int                 `json:"outlet_id"`
	BankID           int                  `json:"bank_id"`
	BankIDCollecting int                  `json:"bank_id_collecting"`
	AccountNo        *string              `json:"account_no"`
	AccountName      string               `json:"account_name"`
	TransferDate     *string              `json:"transfer_date"`
	Amount           float64              `json:"amount"`
	StatusBank       int                  `json:"status_bank_transfer"`
	CreatedBy        *int64               `json:"created_by"`
	UpdatedBy        int64                `json:"updated_by"`
	ProofOfPayment   ProofOfPaymentUpdate `json:"proof_of_payment"`
}

type ProofOfPaymentUpdate struct {
	Mode  string                 `json:"mode"`
	Files []BankTransferFileBody `json:"files"`
}

type DetailBankTransferParams struct {
	BankTransferNo int `params:"bank_transfer_no" validate:"required"`
}
type DeleteBankTransferParams struct {
	BankTransferNo int `params:"bank_transfer_no" validate:"required"`
}
type UpdateBankTransferParams struct {
	BankTransferNo int `params:"bank_transfer_no" validate:"required"`
}

type BankLookupBankTransfer struct {
	BankId   int    `json:"bank_id"`
	BankCode string `json:"bank_code"`
	BankName string `json:"bank_name"`
}

type BankAccountLookupBankTransfer struct {
	AccountNo string `json:"account_no"`
}

var StatusBank = map[int]string{
	1: "Rejected",
	2: "Pending",
	3: "Accepted",
}

var OwnerBank = map[int]string{
	1: "Outlet",
	2: "Distributor",
}

func ConvStatusBankTransfer(data map[int]string, param int) string {
	statusString, ok := data[int(param)]
	if !ok {
		statusString = "Unknown"
	}
	return statusString
}
