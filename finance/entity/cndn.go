package entity

import "time"

type CreateCndnBody struct {
	CustID              string  `json:"cust_id"`
	CndnNo              string  `json:"cndn_no"`
	CndnDate            *string `json:"cndn_date"`
	OwnerId             int     `json:"owner_id"`
	OutletId            int64   `json:"outlet_id"`
	CndnJenis           string  `json:"cndn_jenis"`
	CndnType            int64   `json:"cndn_type"`
	Amount              float64 `json:"amount"`
	UsedAmount          float64 `json:"used_amount"`
	RemainingAmount     float64 `json:"remaning_amount"`
	LastTransactionDate *string `json:"last_transaction_date"`
	Notes               *string `json:"notes"`
	CreatedBy           *int64  `json:"created_by"`
}
type CndnListResponse struct {
	CndnNo              string    `json:"cndn_no"`
	CndnDate            *string   `json:"cndn_date"`
	OwnerId             int       `json:"owner_id"`
	OwnerName           string    `json:"owner_name"`
	OutletID            *int      `json:"outlet_id"`
	OutletCode          string    `json:"outlet_code"`
	OutletName          *string   `json:"outlet_name"`
	CndnJenis           string    `json:"cndn_jenis"`
	CndnType            string    `json:"cndn_type"`
	Amount              float64   `json:"amount"`
	UsedAmount          float64   `json:"used_amount"`
	RemainingAmount     float64   `json:"remaning_amount"`
	LastTransactionDate *string   `json:"last_transaction_date"`
	Notes               string    `json:"notes"`
	CreatedBy           int64     `json:"created_by"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedBy           int64     `json:"updated_by"`
	UpdatedByName       *string   `json:"updated_by_name"`
	UpdatedAt           time.Time `json:"updated_at"`
}

type CndnListDetailResponse struct {
	CndnNo              string    `json:"cndn_no"`
	CndnDate            *string   `json:"cndn_date"`
	OwnerId             int       `json:"owner_id"`
	OwnerName           string    `json:"owner_name"`
	OutletID            *int      `json:"outlet_id"`
	OutletCode          string    `json:"outlet_code"`
	OutletName          *string   `json:"outlet_name"`
	CndnId              int64     `json:"cndn_id"`
	CndnCode            string    `json:"cndn_code"`
	CndnName            string    `json:"cndn_name"`
	CndnJenis           string    `json:"cndn_jenis"`
	Amount              float64   `json:"amount"`
	UsedAmount          float64   `json:"used_amount"`
	RemainingAmount     float64   `json:"remaning_amount"`
	LastTransactionDate *string   `json:"last_transaction_date"`
	Notes               string    `json:"notes"`
	CreatedBy           int64     `json:"created_by"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedBy           int64     `json:"updated_by"`
	UpdatedByName       *string   `json:"updated_by_name"`
	UpdatedAt           time.Time `json:"updated_at"`
}

type UpdateCndnBody struct {
	CustID              string     `json:"cust_id"`
	CndnNo              string     `json:"cndn_no"`
	CndnDate            *string    `json:"cndn_date"`
	OwnerId             int        `json:"owner_id"`
	OutletId            int64      `json:"outlet_id"`
	CndnJenis           string     `json:"cndn_jenis"`
	CndnType            int64      `json:"cndn_type"`
	Amount              float64    `json:"amount"`
	UsedAmount          float64    `json:"used_amount"`
	RemainingAmount     float64    `json:"remaning_amount"`
	LastTransactionDate *string    `json:"last_transaction_date"`
	Notes               string     `json:"notes"`
	UpdatedAt           *time.Time `json:"updated_at"`
	UpdatedBy           int64      `json:"updated_by"`
	UpdatedByName       string     `json:"updated_by_name"`
}

type DetailCndnParams struct {
	CndnNo string `params:"cndn_no" validate:"required"`
}
type DeleteCndnParams struct {
	CndnNo string `params:"cndn_no" validate:"required"`
}
type UpdateCndnParams struct {
	CndnNo string `params:"cndn_no" validate:"required"`
}

type CndnQueryFilter struct {
	CustId       string
	ParentCustId string
	// From         *int64  `query:"from" validate:"required_with=To,omitempty,gte=1000000000"`
	// To           *int64  `query:"to" validate:"required_with=From,omitempty,lte=9999999999,gtefield=From"`
	From       *string `query:"from"`
	To         *string `query:"to"`
	Page       int     `query:"page"`
	Limit      int     `query:"limit" validate:"required"`
	Query      string  `query:"q"`
	Mode       string  `query:"mode"`
	Sort       string  `query:"sort"`
	OutletId   *int    `query:"outlet_id"`
	SuptId     *int    `query:"sup_id"`
	OwnerId    int     `query:"owner_id"`
	CndnJenis  *string `query:"cndn_jenis"`
	DocumentNo string  `query:"document_no"`
}

var OwnerId = map[int]string{
	1: "Outlet",
	2: "Distributor",
}

var CndnJenis = map[string]string{
	"credit": "Credit",
	"debit":  "Debit",
}

func ConvStatusOwnerId(data map[int]string, param int) string {
	statusString, ok := data[int(param)]
	if !ok {
		statusString = "Unknown"
	}
	return statusString
}

func ConvCndnJenis(data map[string]string, param string) string {
	cndnJnsString, ok := data[string(param)]
	if !ok {
		cndnJnsString = "Unknown"
	}
	return cndnJnsString
}
