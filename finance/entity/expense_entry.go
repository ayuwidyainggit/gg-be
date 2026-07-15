package entity

import "time"

type ExpenseEntryQueryFilter struct {
	CustID       string   `validate:"required"`
	ParentCustID string   `json:"-"`
	UserID       int64    `json:"-"`
	Query        string   `query:"q"`
	StartDate    string   `query:"start_date"`
	EndDate      string   `query:"end_date"`
	MinBalance   *float64 `query:"min_balance"`
	CollectorIDs []int64  `json:"-"`
	Page         int      `query:"page" validate:"required"`
	Limit        int      `query:"limit" validate:"required"`
	Sort         string   `query:"sort"`
}

type ExpenseEntryListResponse struct {
	ExpenseID       int64   `json:"expense_id"`
	DocumentNo      string  `json:"document_no"`
	Date            string  `json:"date"`
	ExpenseTypeID   int     `json:"expense_type_id"`
	ExpenseTypeCode string  `json:"expense_type_code"`
	ExpenseTypeName string  `json:"expense_type_name"`
	CollectorID     *int64  `json:"collector_id"`
	CollectorName   string  `json:"collector_name"`
	Balance         float64 `json:"balance"`
	Amount          float64 `json:"amount"`
	Reason          string  `json:"reason"`
	IsClockOut      int     `json:"is_clock_out"`
}

// ExpenseDetailResponse is the response for GET /expense/:expense_id
type ExpenseDetailResponse struct {
	ExpenseID       int64                  `json:"expense_id"`
	Date            string                 `json:"date"`
	DocNo           string                 `json:"doc_no"`
	ExpenseTypeID   int                    `json:"expense_type_id"`
	ExpenseTypeCode string                 `json:"expense_type_code"`
	ExpenseTypeName string                 `json:"expense_type_name"`
	CollectorID     *int64                 `json:"collector_id"`
	CollectorName   string                 `json:"collector_name"`
	Amount          float64                `json:"amount"`
	RemainingAmount float64                `json:"remaining_amount"`
	Note            string                 `json:"note"`
	Files           []ExpenseDetailFile    `json:"files"`
	Deposits        []ExpenseDetailDeposit `json:"deposits"`
}

type ExpenseDetailFile struct {
	FileName      string `json:"file_name"`
	FileType      string `json:"file_type"`
	MediaCategory string `json:"media_category"`
	FileURL       string `json:"file_url"`
	FileSize      int64  `json:"file_size"`
}

type ExpenseDetailDeposit struct {
	DepositExpenseID int64     `json:"deposit_expense_id"`
	DepositID        int64     `json:"deposit_id"`
	UsedAmount       float64   `json:"used_amount"`
	DepositNo        string    `json:"deposit_no"`
	UpdateDate       time.Time `json:"update_date"`
}

// CreateExpenseEntryBody is the request body for POST /acf/v1/expense
type CreateExpenseEntryBody struct {
	ExpenseTypeID int      `json:"expense_type_id" validate:"required"`
	Amount        float64  `json:"amount" validate:"required"`
	CollectorID   *int     `json:"collector_id"`
	Note          *string  `json:"note" validate:"omitempty,max=100"`
	FileURL       []string `json:"file_url"`
}

// UpdateExpenseEntryBody is the request body for PATCH /acf/v1/expense/{expense_id}
type UpdateExpenseEntryBody struct {
	Amount *float64 `json:"amount"`
	Note   *string  `json:"note" validate:"omitempty,max=100"`
}

// CreateExpenseResponseData is the data payload for create expense response
type CreateExpenseResponseData struct {
	ExpenseID  int64  `json:"expense_id"`
	DocumentNo string `json:"document_no"`
}

// UpdateExpenseResponseData is the data payload for update expense response
type UpdateExpenseResponseData struct {
	ExpenseID int64   `json:"expense_id"`
	Amount    float64 `json:"amount"`
	Note      *string `json:"note"`
}

type ExpenseEntryParams struct {
	ExpenseID int64 `params:"expense_id" validate:"required"`
}
